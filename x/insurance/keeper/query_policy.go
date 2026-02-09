package keeper

import (
	"context"
	"errors"

	"realfin/x/insurance/types"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ListPolicy(ctx context.Context, req *types.QueryAllPolicyRequest) (*types.QueryAllPolicyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	policies, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.Policy,
		req.Pagination,
		func(_ string, value types.Policy) (types.Policy, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllPolicyResponse{Policy: policies, Pagination: pageRes}, nil
}

func (q queryServer) GetPolicy(ctx context.Context, req *types.QueryGetPolicyRequest) (*types.QueryGetPolicyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.Policy.Get(ctx, req.PolicyId)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetPolicyResponse{Policy: val}, nil
}
