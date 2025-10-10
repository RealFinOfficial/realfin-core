package realestate

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"realfin/x/realestate/types"
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
					RpcMethod: "ListRate",
					Use:       "list-rate",
					Short:     "List all rate",
				},
				{
					RpcMethod:      "GetRate",
					Use:            "get-rate [id]",
					Short:          "Gets a rate",
					Alias:          []string{"show-rate"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "symbol"}},
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "CreateRate",
					Use:            "create-rate [symbol] [rate] [name] [description]",
					Short:          "Create a new rate",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "symbol"}, {ProtoField: "rate"}, {ProtoField: "name"}, {ProtoField: "description"}},
				},
				{
					RpcMethod:      "UpdateRate",
					Use:            "update-rate [symbol] [rate] [name] [description]",
					Short:          "Update rate",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "symbol"}, {ProtoField: "rate"}, {ProtoField: "name"}, {ProtoField: "description"}},
				},
				{
					RpcMethod:      "DeleteRate",
					Use:            "delete-rate [symbol]",
					Short:          "Delete rate",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "symbol"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
