package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"fundchain/x/milestones/types"
)

func (q queryServer) Project(ctx context.Context, req *types.QueryProjectRequest) (*types.QueryProjectResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id must be > 0")
	}

	p, found, err := q.k.GetProject(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !found {
		return nil, status.Error(codes.NotFound, "project not found")
	}

	return &types.QueryProjectResponse{Project: p}, nil
}
