package keeper

import (
	"context"
	"fmt"
	"strconv"

	"fundchain/x/milestones/types"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) ReleaseTranche(ctx context.Context, msg *types.MsgReleaseTranche) (*types.MsgReleaseTrancheResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	if msg.ProjectId == "" {
		return nil, fmt.Errorf("project_id cannot be empty")
	}
	id, err := strconv.ParseUint(msg.ProjectId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid project_id: %w", err)
	}

	// load project
	p, found, err := k.GetProject(ctx, id)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("project %d not found", id)
	}

	// only allow release when accepted and tranche < 3
	if p.Status != "accepted" && p.Status != "completed" {
		return nil, fmt.Errorf("project status must be accepted or completed to release tranche; got %s", p.Status)
	}
	if p.Tranche >= 3 {
		return nil, fmt.Errorf("all tranches already released")
	}

	// params
	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, err
	}
	if params.Treasury == "" {
		return nil, fmt.Errorf("treasury param is not set")
	}
	if params.Denom == "" {
		return nil, fmt.Errorf("denom param is not set")
	}

	treasuryAddr, err := sdk.AccAddressFromBech32(params.Treasury)
	if err != nil {
		return nil, fmt.Errorf("invalid treasury address: %w", err)
	}
	ownerAddr, err := sdk.AccAddressFromBech32(p.Owner)
	if err != nil {
		return nil, fmt.Errorf("invalid project owner address: %w", err)
	}

	// compute tranche amount: 1/3 of budget; last tranche gets remainder
	base := p.Budget / 3
	remainder := p.Budget - base*3
	amountU64 := base
	if p.Tranche == 2 { // third tranche (0,1,2)
		amountU64 += remainder
	}
	if amountU64 == 0 {
		return nil, fmt.Errorf("computed tranche amount is zero")
	}
	coin := sdk.NewCoin(params.Denom, sdkmath.NewIntFromUint64(amountU64))

	// send coins from treasury to owner
	if err := k.bankKeeper.SendCoins(ctx, treasuryAddr, ownerAddr, sdk.NewCoins(coin)); err != nil {
		return nil, fmt.Errorf("send coins failed: %w", err)
	}

	// update project tranche and possibly status
	p.Tranche += 1
	if p.Tranche >= 3 {
		p.Status = "completed"
	}
	if err := k.SetProject(ctx, p); err != nil {
		return nil, err
	}

	return &types.MsgReleaseTrancheResponse{}, nil
}
