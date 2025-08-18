package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"fundchain/x/milestones/types"
)

// Projects lists all projects
func (q queryServer) Projects(ctx context.Context, req *types.QueryProjectsRequest) (*types.QueryProjectsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	projects := []types.Project{}
	if err := q.k.Projects.Walk(ctx, nil, func(_ uint64, p types.Project) (bool, error) {
		projects = append(projects, p)
		return false, nil
	}); err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryProjectsResponse{Projects: projects}, nil
}
