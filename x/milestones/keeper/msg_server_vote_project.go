package keeper

import (
	"context"
	"strconv"

	"fundchain/x/milestones/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) VoteProject(ctx context.Context, msg *types.MsgVoteProject) (*types.MsgVoteProjectResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	id, err := strconv.ParseUint(msg.ProjectId, 10, 64)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid project_id")
	}

	project, found, err := k.GetProject(ctx, id)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errorsmod.Wrapf(types.ErrNotFound, "project %d", id)
	}

	if msg.Support {
		project.VYes++
	} else {
		project.VNo++
	}

	if project.VYes > project.VNo {
		project.Status = "accepted"
	}
	if err := k.SetProject(ctx, project); err != nil {
		return nil, err
	}

	return &types.MsgVoteProjectResponse{}, nil
}
