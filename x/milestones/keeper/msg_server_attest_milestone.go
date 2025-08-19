package keeper

import (
	"context"

	"fundchain/x/milestones/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) AttestMilestone(ctx context.Context, msg *types.MsgAttestMilestone) (*types.MsgAttestMilestoneResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgAttestMilestoneResponse{}, nil
}
