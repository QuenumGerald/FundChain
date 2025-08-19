package cli

import (
	"github.com/spf13/cobra"

	client "github.com/cosmos/cosmos-sdk/client"
)

// NewTxCmd returns the root tx command for the milestones module.
func NewTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "milestones",
		Short:                      "Milestones transactions",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE: func(cmd *cobra.Command, args []string) error {
			return client.ValidateCmd(cmd, args)
		},
	}

	cmd.AddCommand(NewSubmitProjectCmd())
	cmd.AddCommand(NewVoteProjectCmd())
	cmd.AddCommand(NewAttestMilestoneCmd())
	cmd.AddCommand(NewReleaseTrancheCmd())

	return cmd
}
