package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
)

// nextProjectID returns the next project id using Peek+1 semantics and persists it.
// This guarantees IDs start at 1 and are never zero.
func (k Keeper) nextProjectID(ctx context.Context) (uint64, error) {
	cur, err := k.ProjectSeq.Peek(ctx)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return 0, err
		}
		cur = 0
	}
	next := cur + 1
	if err := k.ProjectSeq.Set(ctx, next); err != nil {
		return 0, err
	}
	return next, nil
}
