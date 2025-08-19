package keeper

import (
    "context"
    "strconv"

    "fundchain/x/milestones/types"

    errorsmod "cosmossdk.io/errors"
    sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) VoteProject(ctx context.Context, msg *types.MsgVoteProject) (*types.MsgVoteProjectResponse, error) {
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

    p, found, err := k.GetProject(ctx, id)
    if err != nil {
        return nil, err
    }
    if !found {
        return nil, errorsmod.Wrapf(types.ErrNotFound, "project %d not found", id)
    }

    if msg.Support {
        p.VYes++
    } else {
        p.VNo++
    }

    if p.VYes > p.VNo {
        p.Status = "accepted"
    }

    if err := k.SetProject(ctx, p); err != nil {
        return nil, err
    }

    // emit event
    sdkCtx := sdk.UnwrapSDKContext(ctx)
    sdkCtx.EventManager().EmitEvent(
        sdk.NewEvent(
            types.EventVoteProject,
            sdk.NewAttribute(types.AttrProjectID, strconv.FormatUint(p.Id, 10)),
            sdk.NewAttribute(types.AttrOwner, msg.Creator),
        ),
    )

    return &types.MsgVoteProjectResponse{}, nil
}
