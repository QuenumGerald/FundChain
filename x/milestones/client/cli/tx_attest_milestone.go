package cli

import (
	"github.com/spf13/cobra"

	client "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"fundchain/x/milestones/types"
)

func NewAttestMilestoneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attest-milestone [projectId] [hash]",
		Short: "Attest a project milestone",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgAttestMilestone{
				Creator:       clientCtx.GetFromAddress().String(),
				ProjectId:     args[0],
				MilestoneHash: args[1],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
