package keeper_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	storetypes "cosmossdk.io/store/types"
	"github.com/stretchr/testify/require"

	"cosmossdk.io/core/address"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"fundchain/testutil/sample"
	"fundchain/x/milestones/keeper"
	module "fundchain/x/milestones/module"
	"fundchain/x/milestones/types"
)

// mockBankKeeper is a simple in-memory bank keeper for tests
type mockBankKeeper struct{
	balances map[string]sdk.Coins
}

func newMockBankKeeper() *mockBankKeeper { return &mockBankKeeper{balances: make(map[string]sdk.Coins)} }

func (m *mockBankKeeper) SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins {
	return m.balances[addr.String()]
}
func (m *mockBankKeeper) SendCoins(ctx context.Context, from, to sdk.AccAddress, coins sdk.Coins) error {
	fromStr, toStr := from.String(), to.String()
	fromBal := m.balances[fromStr]
	if !fromBal.IsAllGTE(coins) {
		return fmt.Errorf("insufficient funds")
	}
	m.balances[fromStr] = fromBal.Sub(coins...)
	m.balances[toStr] = m.balances[toStr].Add(coins...)
	return nil
}

// test fixture with mock bank
func initFixtureWithBank(t *testing.T) (context.Context, keeper.Keeper, address.Codec, *mockBankKeeper) {
	t.Helper()
	encCfg := moduletestutil.MakeTestEncodingConfig(module.AppModule{})
	addressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	storeService := runtime.NewKVStoreService(storeKey)
	ctx := testutil.DefaultContextWithDB(t, storeKey, storetypes.NewTransientStoreKey("transient_test")).Ctx

	authority := authtypes.NewModuleAddress(types.GovModuleName)
	bk := newMockBankKeeper()
	k := keeper.NewKeeper(storeService, encCfg.Codec, addressCodec, []byte(authority), bk)
	require.NoError(t, k.Params.Set(ctx, types.DefaultParams()))
	return ctx, k, addressCodec, bk
}

func TestMsgServer_TableDriven(t *testing.T) {
	ctx, k, addrCodec, bank := initFixtureWithBank(t)
	msgSrv := keeper.NewMsgServerImpl(k)

	// setup params
	params, err := k.GetParams(ctx)
	require.NoError(t, err)
	params.Denom = "ufund"
	treasury := sample.AccAddress()
	params.Treasury = treasury
	require.NoError(t, k.SetParams(ctx, params))

	// fund treasury
	treasuryAddr, err := addrCodec.StringToBytes(treasury)
	require.NoError(t, err)
	bank.balances[sdk.AccAddress(treasuryAddr).String()] = sdk.NewCoins(sdk.NewInt64Coin("ufund", 1_000_000))

	// Setup: submit a project and vote to accept it so later subtests can run standalone
	var projectID uint64
	creator := sample.AccAddress()
	res, err := msgSrv.SubmitProject(ctx, &types.MsgSubmitProject{
		Creator:  creator,
		Title:    "Proj A",
		Budget:   "900",
		IpfsHash: "Qm123",
	})
	require.NoError(t, err)
	projectID = res.Id
	// verify project exists and status submitted
	p, found, err := k.GetProject(ctx, projectID)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, projectID, p.Id)
	require.Equal(t, "submitted", p.Status)
	// vote yes -> status accepted
	_, err = msgSrv.VoteProject(ctx, &types.MsgVoteProject{
		Creator:   sample.AccAddress(),
		ProjectId: strconv.FormatUint(projectID, 10),
		Support:   true,
	})
	require.NoError(t, err)
	p, found, err = k.GetProject(ctx, projectID)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, "accepted", p.Status)

	t.Run("attest -> milestone appended", func(t *testing.T) {
		creator := sample.AccAddress()
		_, err := msgSrv.AttestMilestone(ctx, &types.MsgAttestMilestone{
			Creator:       creator,
			ProjectId:     strconv.FormatUint(projectID, 10),
			MilestoneHash: "hash-1",
		})
		require.NoError(t, err)
		// query milestones via query server
		q := keeper.NewQueryServerImpl(k)
		require.NotZero(t, projectID)
		resp, err := q.ProjectMilestones(ctx, &types.QueryProjectMilestonesRequest{ProjectId: projectID, Pagination: &query.PageRequest{Limit: 10}})
		require.NoError(t, err)
		require.Len(t, resp.Milestones, 1)
		require.Equal(t, "hash-1", resp.Milestones[0].Hash)
	})

	t.Run("release called 3x -> bank balances change, tranche=3, status=completed", func(t *testing.T) {
		owner := sample.AccAddress()
		// set project owner to an account we control
		p, found, err := k.GetProject(ctx, projectID)
		require.NoError(t, err)
		require.True(t, found)
		p.Owner = owner
		require.NoError(t, k.SetProject(ctx, p))

		ownerAddr, err := addrCodec.StringToBytes(owner)
		require.NoError(t, err)
		ownerBech := sdk.AccAddress(ownerAddr).String()
		bank.balances[ownerBech] = sdk.NewCoins() // zero to start

		// call release 3 times
		for i := 0; i < 3; i++ {
			_, err := msgSrv.ReleaseTranche(ctx, &types.MsgReleaseTranche{Creator: sample.AccAddress(), ProjectId: strconv.FormatUint(projectID, 10)})
			require.NoError(t, err)
		}

		// check project updated
		p2, found, err := k.GetProject(ctx, projectID)
		require.NoError(t, err)
		require.True(t, found)
		require.Equal(t, uint64(3), p2.Tranche)
		require.Equal(t, "completed", p2.Status)

		// check balances
		params, _ := k.GetParams(ctx)
		tr := sdk.MustAccAddressFromBech32(params.Treasury)
		trBal := bank.balances[tr.String()].AmountOf(params.Denom)
		ownerBal := bank.balances[ownerBech].AmountOf(params.Denom)

		// budget was 900, so total out = 900
		require.Equal(t, int64(1_000_000-900), trBal.Int64())
		require.Equal(t, int64(900), ownerBal.Int64())
	})
}
