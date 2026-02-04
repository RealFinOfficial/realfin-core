package tokenization

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"realfin/x/tokenization/types"
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
					RpcMethod: "ListAsset",
					Use:       "list-asset",
					Short:     "List all asset",
				},
				{
					RpcMethod:      "GetAsset",
					Use:            "get-asset [id]",
					Short:          "Gets an asset",
					Alias:          []string{"show-asset"},
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
					RpcMethod:      "CreateAsset",
					Use:            "create-asset [symbol] [name] [description] [asset_type] [metadata]",
					Short:          "Create a new asset",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "symbol"}, {ProtoField: "name"}, {ProtoField: "description"}, {ProtoField: "asset_type"}, {ProtoField: "metadata"}},
				},
				{
					RpcMethod:      "UpdateAsset",
					Use:            "update-asset [symbol] [name] [description] [asset_type] [metadata]",
					Short:          "Update asset",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "symbol"}, {ProtoField: "name"}, {ProtoField: "description"}, {ProtoField: "asset_type"}, {ProtoField: "metadata"}},
				},
				{
					RpcMethod:      "DeleteAsset",
					Use:            "delete-asset [symbol]",
					Short:          "Delete asset",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "symbol"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
