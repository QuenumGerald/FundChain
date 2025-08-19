package keeper

import (
	"context"

	"fundchain/x/milestones/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) ReleaseTranche(ctx context.Context, msg *types.MsgReleaseTranche) (*types.MsgReleaseTrancheResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgReleaseTrancheResponse{}, nil
}
