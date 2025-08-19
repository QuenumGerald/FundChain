package keeper

import (
	"context"
	"fmt"
	"strconv"

	"fundchain/x/milestones/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) SubmitProject(ctx context.Context, msg *types.MsgSubmitProject) (*types.MsgSubmitProjectResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// basic validation
	if msg.Title == "" {
		return nil, fmt.Errorf("title cannot be empty")
	}
	if msg.Budget == "" {
		return nil, fmt.Errorf("budget cannot be empty")
	}
	budget, err := strconv.ParseUint(msg.Budget, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid budget: %w", err)
	}
	if budget == 0 {
		return nil, fmt.Errorf("budget must be > 0")
	}
	if msg.IpfsHash == "" {
		return nil, fmt.Errorf("ipfs hash cannot be empty")
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
	}

	if _, err := k.AppendProject(ctx, p); err != nil {
		return nil, err
	}

	return &types.MsgSubmitProjectResponse{}, nil
}
