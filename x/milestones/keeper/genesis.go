package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	"fundchain/x/milestones/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	// set params
	if err := k.Params.Set(ctx, genState.Params); err != nil {
		return err
	}

	// set projects and compute max id
	var maxID uint64
	for _, p := range genState.Projects {
		if err := k.Projects.Set(ctx, p.Id, p); err != nil {
			return err
		}
		if p.Id > maxID {
			maxID = p.Id
		}
	}

	// set milestones
	for _, me := range genState.Milestones {
		// ensure milestone has correct project id
		m := me.Milestone
		m.ProjectId = me.ProjectId
		if err := k.Milestones.Set(ctx, collections.Join(me.ProjectId, me.Index), m); err != nil {
			return err
		}
	}

	// set sequence from next_project_id if provided, else from maxID
	var nextID uint64 = genState.NextProjectId
	if nextID == 0 {
		// if no explicit next, assume next is maxID+1 (Sequence expects current value)
		if maxID > 0 {
			nextID = maxID + 1
		} else {
			nextID = 0
		}
	}
	// Sequence.Set sets current value; Next() returns current+1, so set to nextID-1 safely
	if nextID > 0 {
		if err := k.ProjectSeq.Set(ctx, nextID-1); err != nil {
			return err
		}
	} else {
		// ensure sequence is initialized to default (0)
		if err := k.ProjectSeq.Set(ctx, 0); err != nil {
			return err
		}
	}

	return nil
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	gs := types.DefaultGenesis()

	// params
	params, err := k.Params.Get(ctx)
	if err != nil {
        // When exporting before InitGenesis has ever set Params, collection may be empty.
        // Default to module defaults instead of failing the export.
        if errors.Is(err, collections.ErrNotFound) {
            params = types.DefaultParams()
        } else {
            return nil, err
        }
    }
    gs.Params = params

	// next project id from sequence
	// Peek returns current value; Next would increment. We want next id, so peek+1
	cur, err := k.ProjectSeq.Peek(ctx)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return nil, err
		}
		cur = 0
	}
	gs.NextProjectId = cur + 1

	// projects
	if err := k.Projects.Walk(ctx, nil, func(id uint64, p types.Project) (bool, error) {
		gs.Projects = append(gs.Projects, p)
		return false, nil
	}); err != nil {
		return nil, err
	}

	// milestones
	if err := k.Milestones.Walk(ctx, nil, func(key collections.Pair[uint64, uint64], m types.Milestone) (bool, error) {
		gs.Milestones = append(gs.Milestones, types.MilestoneEntry{
			ProjectId: key.K1(),
			Index:     key.K2(),
			Milestone: m,
		})
		return false, nil
	}); err != nil {
		return nil, err
	}

	return gs, nil
}
