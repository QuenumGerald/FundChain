package keeper

import (
    "context"
    "fmt"
    "strconv"

    "fundchain/x/milestones/types"

    errorsmod "cosmossdk.io/errors"
)

func (k msgServer) VoteProject(ctx context.Context, msg *types.MsgVoteProject) (*types.MsgVoteProjectResponse, error) {
    if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
        return nil, errorsmod.Wrap(err, "invalid authority address")
    }

    if msg.ProjectId == "" {
        return nil, fmt.Errorf("project_id cannot be empty")
    }
    id, err := strconv.ParseUint(msg.ProjectId, 10, 64)
    if err != nil {
        return nil, fmt.Errorf("invalid project_id: %w", err)
    }

    p, found, err := k.GetProject(ctx, id)
    if err != nil {
        return nil, err
    }
    if !found {
        return nil, fmt.Errorf("project %d not found", id)
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

    return &types.MsgVoteProjectResponse{}, nil
}
