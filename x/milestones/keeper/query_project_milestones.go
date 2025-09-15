package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"fundchain/x/milestones/types"
)

func (q queryServer) ProjectMilestones(ctx context.Context, req *types.QueryProjectMilestonesRequest) (*types.QueryProjectMilestonesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	if req.ProjectId == 0 {
		return nil, status.Error(codes.InvalidArgument, "project_id must be > 0")
	}

	// ensure project exists
	_, found, err := q.k.GetProject(ctx, req.ProjectId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !found {
		return nil, status.Error(codes.NotFound, "project not found")
	}

	ms, err := q.k.ListMilestones(ctx, req.ProjectId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// pagination
	offset := uint64(0)
	limit := uint64(len(ms))
	if req.Pagination != nil {
		offset = req.Pagination.Offset
		if req.Pagination.Limit > 0 {
			limit = req.Pagination.Limit
		}
	}
	end := offset + limit
	if offset > uint64(len(ms)) {
		offset = uint64(len(ms))
	}
	if end > uint64(len(ms)) {
		end = uint64(len(ms))
	}
	items := ms[offset:end]

	return &types.QueryProjectMilestonesResponse{Milestones: items}, nil
}
