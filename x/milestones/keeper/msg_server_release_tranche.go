package keeper

import (
    "context"
    "strconv"

    "fundchain/x/milestones/types"

    errorsmod "cosmossdk.io/errors"
)

func (k msgServer) ReleaseTranche(ctx context.Context, msg *types.MsgReleaseTranche) (*types.MsgReleaseTrancheResponse, error) {
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

    // TODO: transfer tranche amount from treasury to project creator once denom and treasury account are defined

    // increment tranche count and update status
    project.Tranche++
    if project.Tranche >= 3 {
        project.Status = "completed"
    }

    if err := k.SetProject(ctx, project); err != nil {
        return nil, err
    }

    return &types.MsgReleaseTrancheResponse{}, nil
}
