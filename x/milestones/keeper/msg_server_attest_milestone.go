package keeper

import (
    "context"
    "strconv"

    "fundchain/x/milestones/types"

    errorsmod "cosmossdk.io/errors"
    sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) AttestMilestone(ctx context.Context, msg *types.MsgAttestMilestone) (*types.MsgAttestMilestoneResponse, error) {
    if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
        return nil, errorsmod.Wrap(err, "invalid authority address")
    }

    if msg.ProjectId == "" {
        return nil, errorsmod.Wrap(types.ErrInvalidParam, "project_id cannot be empty")
    }
    if msg.MilestoneHash == "" {
        return nil, errorsmod.Wrap(types.ErrInvalidParam, "milestone_hash cannot be empty")
    }
    id, err := strconv.ParseUint(msg.ProjectId, 10, 64)
    if err != nil {
        return nil, errorsmod.Wrap(types.ErrInvalidParam, "invalid project_id")
    }

    if _, found, err := k.GetProject(ctx, id); err != nil {
        return nil, err
    } else if !found {
        return nil, errorsmod.Wrapf(types.ErrNotFound, "project %d not found", id)
    }

    m := types.Milestone{
        ProjectId: id,
        Hash:      msg.MilestoneHash,
        Attesters: []string{msg.Creator},
    }
    if _, err := k.AppendMilestone(ctx, id, m); err != nil {
        return nil, err
    }

    // emit event
    sdkCtx := sdk.UnwrapSDKContext(ctx)
    sdkCtx.EventManager().EmitEvent(
        sdk.NewEvent(
            types.EventAttestMilestone,
            sdk.NewAttribute(types.AttrProjectID, strconv.FormatUint(id, 10)),
            sdk.NewAttribute(types.AttrHash, msg.MilestoneHash),
            sdk.NewAttribute(types.AttrOwner, msg.Creator),
        ),
    )

    return &types.MsgAttestMilestoneResponse{}, nil
}
