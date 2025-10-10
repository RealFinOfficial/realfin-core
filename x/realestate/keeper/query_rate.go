package keeper

import (
	"context"
	"errors"

	"realfin/x/realestate/types"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ListRate(ctx context.Context, req *types.QueryAllRateRequest) (*types.QueryAllRateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	rates, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.Rate,
		req.Pagination,
		func(_ string, value types.Rate) (types.Rate, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllRateResponse{Rate: rates, Pagination: pageRes}, nil
}

func (q queryServer) GetRate(ctx context.Context, req *types.QueryGetRateRequest) (*types.QueryGetRateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.Rate.Get(ctx, req.Symbol)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetRateResponse{Rate: val}, nil
}
