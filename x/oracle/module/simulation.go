package oracle

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"realfin/testutil/sample"
	oraclesimulation "realfin/x/oracle/simulation"
	"realfin/x/oracle/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	oracleGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		PriceMap: []types.Price{{Creator: sample.AccAddress(),
			Symbol: "0",
		}, {Creator: sample.AccAddress(),
			Symbol: "1",
		}}}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&oracleGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgCreatePrice          = "op_weight_msg_oracle"
		defaultWeightMsgCreatePrice int = 100
	)

	var weightMsgCreatePrice int
	simState.AppParams.GetOrGenerate(opWeightMsgCreatePrice, &weightMsgCreatePrice, nil,
		func(_ *rand.Rand) {
			weightMsgCreatePrice = defaultWeightMsgCreatePrice
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreatePrice,
		oraclesimulation.SimulateMsgCreatePrice(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdatePrice          = "op_weight_msg_oracle"
		defaultWeightMsgUpdatePrice int = 100
	)

	var weightMsgUpdatePrice int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdatePrice, &weightMsgUpdatePrice, nil,
		func(_ *rand.Rand) {
			weightMsgUpdatePrice = defaultWeightMsgUpdatePrice
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdatePrice,
		oraclesimulation.SimulateMsgUpdatePrice(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgDeletePrice          = "op_weight_msg_oracle"
		defaultWeightMsgDeletePrice int = 100
	)

	var weightMsgDeletePrice int
	simState.AppParams.GetOrGenerate(opWeightMsgDeletePrice, &weightMsgDeletePrice, nil,
		func(_ *rand.Rand) {
			weightMsgDeletePrice = defaultWeightMsgDeletePrice
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeletePrice,
		oraclesimulation.SimulateMsgDeletePrice(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
