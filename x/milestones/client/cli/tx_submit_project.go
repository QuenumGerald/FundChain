package cli

import (
	"github.com/spf13/cobra"

	client "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"fundchain/x/milestones/types"
)

func NewSubmitProjectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-project [title] [budget] [ipfsHash]",
		Short: "Submit a new project proposal",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress().String()
			title := args[0]
			budget := args[1]
			hash := args[2]

			msg := types.NewMsgSubmitProject(from, title, budget, hash)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
