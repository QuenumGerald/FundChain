package keeper

import (
	"context"
 
 	"cosmossdk.io/collections"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"fundchain/x/milestones/types"
)

// Milestones lists milestones for a given project
func (q queryServer) Milestones(ctx context.Context, req *types.QueryMilestonesRequest) (*types.QueryMilestonesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	milestones := []types.Milestone{}
	// NOTE: For now we walk all milestones and filter by project_id. Can be optimized with a ranged iterator.
	if err := q.k.Milestones.Walk(ctx, nil, func(_ collections.Pair[uint64, string], m types.Milestone) (bool, error) {
		// keyPair is collections.Pair[uint64,string], but to avoid importing generics reflection,
		// rely on milestone value's ProjectId for filtering.
		if m.ProjectId == req.ProjectId {
			milestones = append(milestones, m)
		}
		return false, nil
	}); err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryMilestonesResponse{Milestones: milestones}, nil
}
