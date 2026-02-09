package insurance

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"realfin/testutil/sample"
	insurancesimulation "realfin/x/insurance/simulation"
	"realfin/x/insurance/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	insuranceGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		PolicyMap: []types.Policy{{Creator: sample.AccAddress(),
			PolicyId: "0",
		}, {Creator: sample.AccAddress(),
			PolicyId: "1",
		}}}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&insuranceGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgCreatePolicy          = "op_weight_msg_insurance"
		defaultWeightMsgCreatePolicy int = 100
	)

	var weightMsgCreatePolicy int
	simState.AppParams.GetOrGenerate(opWeightMsgCreatePolicy, &weightMsgCreatePolicy, nil,
		func(_ *rand.Rand) {
			weightMsgCreatePolicy = defaultWeightMsgCreatePolicy
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreatePolicy,
		insurancesimulation.SimulateMsgCreatePolicy(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdatePolicy          = "op_weight_msg_insurance"
		defaultWeightMsgUpdatePolicy int = 100
	)

	var weightMsgUpdatePolicy int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdatePolicy, &weightMsgUpdatePolicy, nil,
		func(_ *rand.Rand) {
			weightMsgUpdatePolicy = defaultWeightMsgUpdatePolicy
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdatePolicy,
		insurancesimulation.SimulateMsgUpdatePolicy(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgDeletePolicy          = "op_weight_msg_insurance"
		defaultWeightMsgDeletePolicy int = 100
	)

	var weightMsgDeletePolicy int
	simState.AppParams.GetOrGenerate(opWeightMsgDeletePolicy, &weightMsgDeletePolicy, nil,
		func(_ *rand.Rand) {
			weightMsgDeletePolicy = defaultWeightMsgDeletePolicy
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeletePolicy,
		insurancesimulation.SimulateMsgDeletePolicy(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
