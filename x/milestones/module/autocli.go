package milestones

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"fundchain/x/milestones/types"
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
					RpcMethod:      "SubmitProject",
					Use:            "submit-project [title] [budget] [ipfs-hash]",
					Short:          "Send a submit-project tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "title"}, {ProtoField: "budget"}, {ProtoField: "ipfs_hash"}},
				},
				{
					RpcMethod:      "VoteProject",
					Use:            "vote-project [project-id] [support]",
					Short:          "Send a vote-project tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "project_id"}, {ProtoField: "support"}},
				},
				{
					RpcMethod:      "AttestMilestone",
					Use:            "attest-milestone [project-id] [milestone-hash]",
					Short:          "Send a attest-milestone tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "project_id"}, {ProtoField: "milestone_hash"}},
				},
				{
					RpcMethod:      "ReleaseTranche",
					Use:            "release-tranche [project-id]",
					Short:          "Send a release-tranche tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "project_id"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
