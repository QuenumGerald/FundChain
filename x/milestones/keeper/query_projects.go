package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"fundchain/x/milestones/types"
)

func (q queryServer) Projects(ctx context.Context, req *types.QueryProjectsRequest) (*types.QueryProjectsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// collect all projects
	var all []types.Project
	if err := q.k.IterateProjects(ctx, func(_ uint64, p types.Project) (stop bool) {
		all = append(all, p)
		return false
	}); err != nil {
		return nil, err
	}

	// simple offset/limit pagination
	offset := uint64(0)
	limit := uint64(len(all))
	if req.Pagination != nil {
		offset = req.Pagination.Offset
		if req.Pagination.Limit > 0 {
			limit = req.Pagination.Limit
		}
	}
	end := offset + limit
	if offset > uint64(len(all)) {
		offset = uint64(len(all))
	}
	if end > uint64(len(all)) {
		end = uint64(len(all))
	}
	items := all[offset:end]

	resp := &types.QueryProjectsResponse{
		Projects: items,
	}
	return resp, nil
}
