package tokenization

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"realfin/testutil/sample"
	tokenizationsimulation "realfin/x/tokenization/simulation"
	"realfin/x/tokenization/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	tokenizationGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		AssetMap: []types.Asset{{Creator: sample.AccAddress(),
			Symbol: "0",
		}, {Creator: sample.AccAddress(),
			Symbol: "1",
		}}}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&tokenizationGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgCreateAsset          = "op_weight_msg_tokenization"
		defaultWeightMsgCreateAsset int = 100
	)

	var weightMsgCreateAsset int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateAsset, &weightMsgCreateAsset, nil,
		func(_ *rand.Rand) {
			weightMsgCreateAsset = defaultWeightMsgCreateAsset
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateAsset,
		tokenizationsimulation.SimulateMsgCreateAsset(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdateAsset          = "op_weight_msg_tokenization"
		defaultWeightMsgUpdateAsset int = 100
	)

	var weightMsgUpdateAsset int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateAsset, &weightMsgUpdateAsset, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateAsset = defaultWeightMsgUpdateAsset
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateAsset,
		tokenizationsimulation.SimulateMsgUpdateAsset(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgDeleteAsset          = "op_weight_msg_tokenization"
		defaultWeightMsgDeleteAsset int = 100
	)

	var weightMsgDeleteAsset int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteAsset, &weightMsgDeleteAsset, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteAsset = defaultWeightMsgDeleteAsset
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteAsset,
		tokenizationsimulation.SimulateMsgDeleteAsset(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
