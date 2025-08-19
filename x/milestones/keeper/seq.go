package keeper

import "context"

// nextProjectID increments and returns the next project id
func (k Keeper) nextProjectID(ctx context.Context) (uint64, error) {
	return k.ProjectSeq.Next(ctx)
}
