package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"fundchain/x/milestones/types"
)

// Project returns a single project by ID
func (q queryServer) Project(ctx context.Context, req *types.QueryProjectRequest) (*types.QueryProjectResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	p, found, err := q.k.GetProject(ctx, req.ProjectId)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	if !found {
		return nil, status.Error(codes.NotFound, "project not found")
	}

	return &types.QueryProjectResponse{Project: p}, nil
}
