package keeper

import (
	"context"
	"strconv"

	"fundchain/x/milestones/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) SubmitProject(ctx context.Context, msg *types.MsgSubmitProject) (*types.MsgSubmitProjectResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// parse budget (stored as uint64 in state, comes as string in msg)
	budget, err := strconv.ParseUint(msg.Budget, 10, 64)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid budget")
	}

	id, err := k.AppendProject(ctx, types.Project{
		Title:    msg.Title,
		Budget:   budget,
		IpfsHash: msg.IpfsHash,
		Status:   "submitted",
	})
	if err != nil {
		return nil, err
	}

	return &types.MsgSubmitProjectResponse{ProjectId: id}, nil
}
