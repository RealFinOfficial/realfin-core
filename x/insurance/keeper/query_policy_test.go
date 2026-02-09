package keeper_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"realfin/x/insurance/keeper"
	"realfin/x/insurance/types"
)

func createNPolicy(keeper keeper.Keeper, ctx context.Context, n int) []types.Policy {
	items := make([]types.Policy, n)
	for i := range items {
		items[i].PolicyId = strconv.Itoa(i)
		items[i].AssetSymbol = strconv.Itoa(i)
		items[i].Provider = strconv.Itoa(i)
		items[i].CoverageType = "full"
		items[i].CoveragePercentage = "100"
		_ = keeper.Policy.Set(ctx, items[i].PolicyId, items[i])
	}
	return items
}

func TestPolicyQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNPolicy(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetPolicyRequest
		response *types.QueryGetPolicyResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetPolicyRequest{
				PolicyId: msgs[0].PolicyId,
			},
			response: &types.QueryGetPolicyResponse{Policy: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetPolicyRequest{
				PolicyId: msgs[1].PolicyId,
			},
			response: &types.QueryGetPolicyResponse{Policy: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetPolicyRequest{
				PolicyId: strconv.Itoa(100000),
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
			response, err := qs.GetPolicy(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestPolicyQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNPolicy(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllPolicyRequest {
		return &types.QueryAllPolicyRequest{
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
			resp, err := qs.ListPolicy(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Policy), step)
			require.Subset(t, msgs, resp.Policy)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListPolicy(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Policy), step)
			require.Subset(t, msgs, resp.Policy)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListPolicy(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.Policy)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListPolicy(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
