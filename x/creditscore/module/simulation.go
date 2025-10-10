package creditscore

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"realfin/testutil/sample"
	creditscoresimulation "realfin/x/creditscore/simulation"
	"realfin/x/creditscore/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	creditscoreGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		RateMap: []types.Rate{{Creator: sample.AccAddress(),
			Symbol: "0",
		}, {Creator: sample.AccAddress(),
			Symbol: "1",
		}}}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&creditscoreGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgCreateRate          = "op_weight_msg_creditscore"
		defaultWeightMsgCreateRate int = 100
	)

	var weightMsgCreateRate int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateRate, &weightMsgCreateRate, nil,
		func(_ *rand.Rand) {
			weightMsgCreateRate = defaultWeightMsgCreateRate
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateRate,
		creditscoresimulation.SimulateMsgCreateRate(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdateRate          = "op_weight_msg_creditscore"
		defaultWeightMsgUpdateRate int = 100
	)

	var weightMsgUpdateRate int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateRate, &weightMsgUpdateRate, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateRate = defaultWeightMsgUpdateRate
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateRate,
		creditscoresimulation.SimulateMsgUpdateRate(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgDeleteRate          = "op_weight_msg_creditscore"
		defaultWeightMsgDeleteRate int = 100
	)

	var weightMsgDeleteRate int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteRate, &weightMsgDeleteRate, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteRate = defaultWeightMsgDeleteRate
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteRate,
		creditscoresimulation.SimulateMsgDeleteRate(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
