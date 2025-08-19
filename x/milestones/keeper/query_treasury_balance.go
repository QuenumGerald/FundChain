package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"fundchain/x/milestones/types"
)

func (q queryServer) TreasuryBalance(ctx context.Context, _ *types.QueryTreasuryBalanceRequest) (*types.QueryTreasuryBalanceResponse, error) {
	// load params for treasury and denom
	params, err := q.k.GetParams(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if params.Treasury == "" || params.Denom == "" {
		return nil, status.Error(codes.FailedPrecondition, "treasury or denom parameter not set")
	}

	treasuryAddr, err := sdk.AccAddressFromBech32(params.Treasury)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid treasury bech32 address")
	}

	coins := q.k.bankKeeper.SpendableCoins(ctx, treasuryAddr)
	amt := coins.AmountOf(params.Denom)
	coin := sdk.NewCoin(params.Denom, amt)

	return &types.QueryTreasuryBalanceResponse{
		Treasury: params.Treasury,
		Balance:  coin,
	}, nil
}
