package keeper_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"realfin/x/tokenization/keeper"
	"realfin/x/tokenization/types"
)

func createNAsset(keeper keeper.Keeper, ctx context.Context, n int) []types.Asset {
	items := make([]types.Asset, n)
	for i := range items {
		items[i].Symbol = strconv.Itoa(i)
		items[i].Name = strconv.Itoa(i)
		items[i].Description = strconv.Itoa(i)
		items[i].AssetType = "real_estate"
		items[i].Metadata = strconv.Itoa(i)
		_ = keeper.Asset.Set(ctx, items[i].Symbol, items[i])
	}
	return items
}

func TestAssetQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNAsset(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetAssetRequest
		response *types.QueryGetAssetResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetAssetRequest{
				Symbol: msgs[0].Symbol,
			},
			response: &types.QueryGetAssetResponse{Asset: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetAssetRequest{
				Symbol: msgs[1].Symbol,
			},
			response: &types.QueryGetAssetResponse{Asset: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetAssetRequest{
				Symbol: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := qs.GetAsset(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestAssetQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNAsset(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllAssetRequest {
		return &types.QueryAllAssetRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListAsset(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Asset), step)
			require.Subset(t, msgs, resp.Asset)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListAsset(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Asset), step)
			require.Subset(t, msgs, resp.Asset)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListAsset(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.Asset)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListAsset(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
