package milestones

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	milestonessimulation "fundchain/x/milestones/simulation"
	"fundchain/x/milestones/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	milestonesGenesis := types.GenesisState{
		Params: types.DefaultParams(),
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&milestonesGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgSubmitProject          = "op_weight_msg_milestones"
		defaultWeightMsgSubmitProject int = 100
	)

	var weightMsgSubmitProject int
	simState.AppParams.GetOrGenerate(opWeightMsgSubmitProject, &weightMsgSubmitProject, nil,
		func(_ *rand.Rand) {
			weightMsgSubmitProject = defaultWeightMsgSubmitProject
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSubmitProject,
		milestonessimulation.SimulateMsgSubmitProject(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgVoteProject          = "op_weight_msg_milestones"
		defaultWeightMsgVoteProject int = 100
	)

	var weightMsgVoteProject int
	simState.AppParams.GetOrGenerate(opWeightMsgVoteProject, &weightMsgVoteProject, nil,
		func(_ *rand.Rand) {
			weightMsgVoteProject = defaultWeightMsgVoteProject
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgVoteProject,
		milestonessimulation.SimulateMsgVoteProject(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgAttestMilestone          = "op_weight_msg_milestones"
		defaultWeightMsgAttestMilestone int = 100
	)

	var weightMsgAttestMilestone int
	simState.AppParams.GetOrGenerate(opWeightMsgAttestMilestone, &weightMsgAttestMilestone, nil,
		func(_ *rand.Rand) {
			weightMsgAttestMilestone = defaultWeightMsgAttestMilestone
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAttestMilestone,
		milestonessimulation.SimulateMsgAttestMilestone(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgReleaseTranche          = "op_weight_msg_milestones"
		defaultWeightMsgReleaseTranche int = 100
	)

	var weightMsgReleaseTranche int
	simState.AppParams.GetOrGenerate(opWeightMsgReleaseTranche, &weightMsgReleaseTranche, nil,
		func(_ *rand.Rand) {
			weightMsgReleaseTranche = defaultWeightMsgReleaseTranche
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgReleaseTranche,
		milestonessimulation.SimulateMsgReleaseTranche(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
