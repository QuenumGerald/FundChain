package keeper

import (
	"context"
	"fmt"
	"strconv"

	"fundchain/x/milestones/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) SubmitProject(ctx context.Context, msg *types.MsgSubmitProject) (*types.MsgSubmitProjectResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// basic validation
	if msg.Title == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidParam, "title cannot be empty")
	}
	if msg.Budget == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidParam, "budget cannot be empty")
	}
	budget, err := strconv.ParseUint(msg.Budget, 10, 64)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidParam, fmt.Sprintf("invalid budget: %v", err))
	}
	if budget == 0 {
		return nil, types.ErrZeroAmount
	}
	if msg.IpfsHash == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidParam, "ipfs hash cannot be empty")
	}

	// reviewers & threshold validation
	if msg.AttestThreshold == 0 {
		return nil, errorsmod.Wrap(types.ErrInvalidParam, "attest_threshold must be >= 1")
	}
	if len(msg.Reviewers) == 0 {
		return nil, errorsmod.Wrap(types.ErrInvalidParam, "reviewers cannot be empty")
	}
	if int(msg.AttestThreshold) > len(msg.Reviewers) {
		return nil, errorsmod.Wrap(types.ErrInvalidParam, "attest_threshold cannot exceed reviewers length")
	}
	for _, r := range msg.Reviewers {
		if _, err := k.addressCodec.StringToBytes(r); err != nil {
			return nil, errorsmod.Wrapf(types.ErrInvalidParam, "invalid reviewer address: %s", r)
		}
	}

	// create project
	p := types.Project{
		Title:    msg.Title,
		Budget:   budget,
		IpfsHash: msg.IpfsHash,
		VYes:     0,
		VNo:      0,
		Tranche:  0,
		Status:   "submitted",
		Owner:    msg.Creator,
		Reviewers: msg.Reviewers,
		AttestThreshold: msg.AttestThreshold,
	}

	id, err := k.AppendProject(ctx, p)
	if err != nil {
		return nil, err
	}

	// emit event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventSubmitProject,
			sdk.NewAttribute(types.AttrProjectID, strconv.FormatUint(id, 10)),
			sdk.NewAttribute(types.AttrOwner, p.Owner),
			sdk.NewAttribute(types.AttrHash, p.IpfsHash),
		),
	)

	return &types.MsgSubmitProjectResponse{Id: id}, nil
}
