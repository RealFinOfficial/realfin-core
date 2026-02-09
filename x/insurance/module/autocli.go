package insurance

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"realfin/x/insurance/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "ListPolicy",
					Use:       "list-policy",
					Short:     "List all policy",
				},
				{
					RpcMethod:      "GetPolicy",
					Use:            "get-policy [id]",
					Short:          "Gets a policy",
					Alias:          []string{"show-policy"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "policy_id"}},
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "CreatePolicy",
					Use:            "create-policy [policy_id] [asset_symbol] [provider] [coverage_type] [coverage_percentage]",
					Short:          "Create a new policy",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "policy_id"}, {ProtoField: "asset_symbol"}, {ProtoField: "provider"}, {ProtoField: "coverage_type"}, {ProtoField: "coverage_percentage"}},
				},
				{
					RpcMethod:      "UpdatePolicy",
					Use:            "update-policy [policy_id] [asset_symbol] [provider] [coverage_type] [coverage_percentage]",
					Short:          "Update policy",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "policy_id"}, {ProtoField: "asset_symbol"}, {ProtoField: "provider"}, {ProtoField: "coverage_type"}, {ProtoField: "coverage_percentage"}},
				},
				{
					RpcMethod:      "DeletePolicy",
					Use:            "delete-policy [policy_id]",
					Short:          "Delete policy",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "policy_id"}},
				},
			},
		},
	}
}
