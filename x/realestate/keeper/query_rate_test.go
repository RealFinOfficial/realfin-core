package keeper_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"realfin/x/realestate/keeper"
	"realfin/x/realestate/types"
)

func createNRate(keeper keeper.Keeper, ctx context.Context, n int) []types.Rate {
	items := make([]types.Rate, n)
	for i := range items {
		items[i].Symbol = strconv.Itoa(i)
		items[i].Rate = uint64(i)
		items[i].Name = strconv.Itoa(i)
		items[i].Description = strconv.Itoa(i)
		_ = keeper.Rate.Set(ctx, items[i].Symbol, items[i])
	}
	return items
}

func TestRateQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNRate(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetRateRequest
		response *types.QueryGetRateResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetRateRequest{
				Symbol: msgs[0].Symbol,
			},
			response: &types.QueryGetRateResponse{Rate: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetRateRequest{
				Symbol: msgs[1].Symbol,
			},
			response: &types.QueryGetRateResponse{Rate: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetRateRequest{
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
			response, err := qs.GetRate(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestRateQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNRate(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllRateRequest {
		return &types.QueryAllRateRequest{
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
			resp, err := qs.ListRate(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Rate), step)
			require.Subset(t, msgs, resp.Rate)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListRate(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Rate), step)
			require.Subset(t, msgs, resp.Rate)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListRate(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.Rate)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListRate(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
