package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	"fundchain/x/milestones/types"
)

// AppendProject stores a new project assigning the next id and returns it.
func (k Keeper) AppendProject(ctx context.Context, p types.Project) (uint64, error) {
	id, err := k.nextProjectID(ctx)
	if err != nil {
		return 0, err
	}
	p.Id = id
	if err := k.Projects.Set(ctx, id, p); err != nil {
		return 0, err
	}
	return id, nil
}

// GetProject returns a project by id and a boolean if found.
func (k Keeper) GetProject(ctx context.Context, id uint64) (types.Project, bool, error) {
	p, err := k.Projects.Get(ctx, id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return types.Project{}, false, nil
		}
		return types.Project{}, false, err
	}
	return p, true, nil
}

// SetProject overwrites a project at its id.
func (k Keeper) SetProject(ctx context.Context, p types.Project) error {
	return k.Projects.Set(ctx, p.Id, p)
}

// IterateProjects walks over all projects and calls cb. If cb returns true, the walk stops.
func (k Keeper) IterateProjects(ctx context.Context, cb func(id uint64, p types.Project) (stop bool)) error {
	return k.Projects.Walk(ctx, nil, func(key uint64, value types.Project) (bool, error) {
		return cb(key, value), nil
	})
}
