package keeper

import (
    "context"
    "strconv"

    "fundchain/x/milestones/types"

    "cosmossdk.io/collections"
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

    p, found, err := k.GetProject(ctx, id)
    if err != nil {
        return nil, err
    } else if !found {
        return nil, errorsmod.Wrapf(types.ErrNotFound, "project %d not found", id)
    }

    // enforce reviewers allowlist
    isReviewer := false
    for _, r := range p.Reviewers {
        if r == msg.Creator {
            isReviewer = true
            break
        }
    }
    if !isReviewer {
        return nil, errorsmod.Wrap(types.ErrUnauthorized, "attester not in project reviewers")
    }

    // find existing milestone by hash; update attesters uniquely or create new
    ms, err := k.ListMilestones(ctx, id)
    if err != nil {
        return nil, err
    }
    updated := false
    for idx, m := range ms {
        if m.Hash == msg.MilestoneHash {
            // check duplicate attestation by same creator
            for _, a := range m.Attesters {
                if a == msg.Creator {
                    return nil, errorsmod.Wrap(types.ErrInvalidParam, "duplicate attestation for this milestone by the same address")
                }
            }
            m.Attesters = append(m.Attesters, msg.Creator)
            // persist back at same index
            // We need the index key; ListMilestones loses index, so walk to find key
            // Re-walk and update by matching hash again
            // Simpler: perform a walk and update in place
            // Since collections API lacks direct update by value without key, do a second walk
            err := k.Milestones.Walk(ctx, collections.NewPrefixedPairRange[uint64, uint64](id), func(key collections.Pair[uint64, uint64], val types.Milestone) (bool, error) {
                if val.Hash == msg.MilestoneHash {
                    return true, k.Milestones.Set(ctx, key, m)
                }
                return false, nil
            })
            if err != nil {
                return nil, err
            }
            updated = true
            _ = idx // idx unused after refactor
            break
        }
    }
    if !updated {
        m := types.Milestone{
            ProjectId: id,
            Hash:      msg.MilestoneHash,
            Attesters: []string{msg.Creator},
        }
        if _, err := k.AppendMilestone(ctx, id, m); err != nil {
            return nil, err
        }
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
