package oracle

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"realfin/x/oracle/types"
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
					RpcMethod: "ListPrice",
					Use:       "list-price",
					Short:     "List all price",
				},
				{
					RpcMethod:      "GetPrice",
					Use:            "get-price [id]",
					Short:          "Gets a price",
					Alias:          []string{"show-price"},
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
					RpcMethod:      "CreatePrice",
					Use:            "create-price [symbol] [rate] [name] [description]",
					Short:          "Create a new price",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "symbol"}, {ProtoField: "rate"}, {ProtoField: "name"}, {ProtoField: "description"}},
				},
				{
					RpcMethod:      "UpdatePrice",
					Use:            "update-price [symbol] [rate] [name] [description]",
					Short:          "Update price",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "symbol"}, {ProtoField: "rate"}, {ProtoField: "name"}, {ProtoField: "description"}},
				},
				{
					RpcMethod:      "DeletePrice",
					Use:            "delete-price [symbol]",
					Short:          "Delete price",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "symbol"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
