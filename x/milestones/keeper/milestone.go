package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"fundchain/x/milestones/types"
)

// AppendMilestone appends a milestone for a given project, assigning the next index.
// The milestone key is (projectId, index).
func (k Keeper) AppendMilestone(ctx context.Context, projectID uint64, m types.Milestone) (uint64, error) {
	// compute next index by scanning existing milestones for projectID
	var nextIdx uint64
	rng := collections.NewPrefixedPairRange[uint64, uint64](projectID)
	if err := k.Milestones.Walk(ctx, rng, func(key collections.Pair[uint64, uint64], _ types.Milestone) (bool, error) {
		if key.K2() >= nextIdx {
			nextIdx = key.K2() + 1
		}
		return false, nil
	}); err != nil {
		return 0, err
	}

	m.ProjectId = projectID
	if err := k.Milestones.Set(ctx, collections.Join(projectID, nextIdx), m); err != nil {
		return 0, err
	}
	return nextIdx, nil
}

// ListMilestones returns all milestones for a given projectID in key order (by index).
func (k Keeper) ListMilestones(ctx context.Context, projectID uint64) ([]types.Milestone, error) {
	ms := []types.Milestone{}
	rng := collections.NewPrefixedPairRange[uint64, uint64](projectID)
	if err := k.Milestones.Walk(ctx, rng, func(_ collections.Pair[uint64, uint64], value types.Milestone) (bool, error) {
		ms = append(ms, value)
		return false, nil
	}); err != nil {
		return nil, err
	}
	return ms, nil
}
