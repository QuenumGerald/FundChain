package cli

import (
	"strconv"

	"github.com/spf13/cobra"

	client "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"fundchain/x/milestones/types"
)

func NewVoteProjectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote-project [projectId] [support:true|false]",
		Short: "Vote for or against a project",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			projectID := args[0]
			support, err := strconv.ParseBool(args[1])
			if err != nil {
				return err
			}

			msg := &types.MsgVoteProject{
				Creator:   clientCtx.GetFromAddress().String(),
				ProjectId: projectID,
				Support:   support,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
