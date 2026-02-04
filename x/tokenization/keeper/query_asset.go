package keeper

import (
	"context"
	"errors"

	"realfin/x/tokenization/types"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ListAsset(ctx context.Context, req *types.QueryAllAssetRequest) (*types.QueryAllAssetResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	assets, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.Asset,
		req.Pagination,
		func(_ string, value types.Asset) (types.Asset, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllAssetResponse{Asset: assets, Pagination: pageRes}, nil
}

func (q queryServer) GetAsset(ctx context.Context, req *types.QueryGetAssetRequest) (*types.QueryGetAssetResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.Asset.Get(ctx, req.Symbol)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetAssetResponse{Asset: val}, nil
}
