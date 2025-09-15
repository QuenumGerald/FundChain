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

			reviewers, err := cmd.Flags().GetStringSlice("reviewer")
			if err != nil {
				return err
			}
			threshold, err := cmd.Flags().GetUint32("threshold")
			if err != nil {
				return err
			}

			msg := &types.MsgSubmitProject{
				Creator:         from,
				Title:           title,
				Budget:          budget,
				IpfsHash:        hash,
				Reviewers:       reviewers,
				AttestThreshold: threshold,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().StringSlice("reviewer", nil, "address allowed to attest; repeat for multiple (required)")
	cmd.Flags().Uint32("threshold", 1, "number of attestations required to validate milestones")
	return cmd
}
