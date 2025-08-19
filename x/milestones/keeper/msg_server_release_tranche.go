package keeper

import (
	"context"
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
		return nil, errorsmod.Wrap(types.ErrInvalidParam, "project_id cannot be empty")
	}
	id, err := strconv.ParseUint(msg.ProjectId, 10, 64)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidParam, "invalid project_id")
	}

	// load project
	p, found, err := k.GetProject(ctx, id)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errorsmod.Wrapf(types.ErrNotFound, "project %d not found", id)
	}

	// only allow release when accepted and tranche < 3
	if p.Status != "accepted" && p.Status != "completed" {
		return nil, errorsmod.Wrap(types.ErrUnauthorized, "project status must be accepted or completed to release tranche")
	}
	if p.Tranche >= 3 {
		return nil, errorsmod.Wrap(types.ErrInvalidParam, "all tranches already released")
	}

	// params
	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, err
	}
	if params.Treasury == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidParam, "treasury param is not set")
	}
	if params.Denom == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidParam, "denom param is not set")
	}

	treasuryAddr, err := sdk.AccAddressFromBech32(params.Treasury)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidParam, "invalid treasury address")
	}
	ownerAddr, err := sdk.AccAddressFromBech32(p.Owner)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidParam, "invalid project owner address")
	}

	// require threshold: at least one milestone with enough unique attestations
	ms, err := k.ListMilestones(ctx, id)
	if err != nil {
		return nil, err
	}
	meets := false
	for _, m := range ms {
		// dedupe in case of any future storage issues (should already be unique)
		seen := make(map[string]struct{})
		for _, a := range m.Attesters {
			seen[a] = struct{}{}
		}
		if uint32(len(seen)) >= p.AttestThreshold {
			meets = true
			break
		}
	}
	if !meets {
		return nil, errorsmod.Wrap(types.ErrUnauthorized, "milestone threshold not met")
	}

	// compute tranche amount: 1/3 of budget; last tranche gets remainder
	base := p.Budget / 3
	remainder := p.Budget - base*3
	amountU64 := base
	if p.Tranche == 2 { // third tranche (0,1,2)
		amountU64 += remainder
	}
	if amountU64 == 0 {
		return nil, types.ErrZeroAmount
	}
	coin := sdk.NewCoin(params.Denom, sdkmath.NewIntFromUint64(amountU64))

	// send coins from treasury to owner
	if err := k.bankKeeper.SendCoins(ctx, treasuryAddr, ownerAddr, sdk.NewCoins(coin)); err != nil {
		return nil, errorsmod.Wrap(types.ErrUnauthorized, err.Error())
	}

	// update project tranche and possibly status
	p.Tranche += 1
	if p.Tranche >= 3 {
		p.Status = "completed"
	}
	if err := k.SetProject(ctx, p); err != nil {
		return nil, err
	}

	// emit event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventReleaseTranche,
			sdk.NewAttribute(types.AttrProjectID, strconv.FormatUint(p.Id, 10)),
			sdk.NewAttribute(types.AttrAmount, coin.Amount.String()),
			sdk.NewAttribute(types.AttrDenom, coin.Denom),
			sdk.NewAttribute(types.AttrOwner, p.Owner),
		),
	)

	return &types.MsgReleaseTrancheResponse{}, nil
}
