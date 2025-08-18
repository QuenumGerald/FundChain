package keeper

import (
	"context"
	"strconv"

	"fundchain/x/milestones/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) AttestMilestone(ctx context.Context, msg *types.MsgAttestMilestone) (*types.MsgAttestMilestoneResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	id, err := strconv.ParseUint(msg.ProjectId, 10, 64)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid project_id")
	}
	// optional existence check
	if _, found, err := k.GetProject(ctx, id); err != nil {
		return nil, err
	} else if !found {
		return nil, errorsmod.Wrapf(types.ErrNotFound, "project %d", id)
	}

	milestone := types.Milestone{
		ProjectId: id,
		Hash:      msg.MilestoneHash,
		Attesters: []string{msg.Creator},
	}
	if err := k.AppendMilestone(ctx, milestone); err != nil {
		return nil, err
	}

	return &types.MsgAttestMilestoneResponse{}, nil
}
