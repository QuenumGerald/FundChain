package keeper

import (
    "context"
    "fmt"
    "strconv"

    "fundchain/x/milestones/types"

    errorsmod "cosmossdk.io/errors"
)

func (k msgServer) AttestMilestone(ctx context.Context, msg *types.MsgAttestMilestone) (*types.MsgAttestMilestoneResponse, error) {
    if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
        return nil, errorsmod.Wrap(err, "invalid authority address")
    }

    if msg.ProjectId == "" {
        return nil, fmt.Errorf("project_id cannot be empty")
    }
    if msg.MilestoneHash == "" {
        return nil, fmt.Errorf("milestone_hash cannot be empty")
    }
    id, err := strconv.ParseUint(msg.ProjectId, 10, 64)
    if err != nil {
        return nil, fmt.Errorf("invalid project_id: %w", err)
    }

    if _, found, err := k.GetProject(ctx, id); err != nil {
        return nil, err
    } else if !found {
        return nil, fmt.Errorf("project %d not found", id)
    }

    m := types.Milestone{
        ProjectId: id,
        Hash:      msg.MilestoneHash,
        Attesters: []string{msg.Creator},
    }
    if _, err := k.AppendMilestone(ctx, id, m); err != nil {
        return nil, err
    }

    return &types.MsgAttestMilestoneResponse{}, nil
}
