package keeper

import (
	"context"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"fundchain/x/milestones/types"
)

// TreasuryBalance returns the spendable balance for the milestones module account
func (q queryServer) TreasuryBalance(ctx context.Context, req *types.QueryTreasuryBalanceRequest) (*types.QueryTreasuryBalanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	moduleAddr := authtypes.NewModuleAddress(types.ModuleName)
	coins := q.k.bankKeeper.SpendableCoins(ctx, moduleAddr)
	return &types.QueryTreasuryBalanceResponse{Balance: coins}, nil
}
