package cli

import (
	"github.com/spf13/cobra"

	client "github.com/cosmos/cosmos-sdk/client"
)

// NewQueryCmd returns the root query command for milestones module.
func NewQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "milestones",
		Short:                      "Querying commands for the milestones module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE: func(cmd *cobra.Command, args []string) error {
			return client.ValidateCmd(cmd, args)
		},
	}

	cmd.AddCommand(NewParamsCmd())
	cmd.AddCommand(NewProjectsCmd())
	cmd.AddCommand(NewProjectCmd())
	cmd.AddCommand(NewMilestonesCmd())
	cmd.AddCommand(NewTreasuryBalanceCmd())

	return cmd
}
