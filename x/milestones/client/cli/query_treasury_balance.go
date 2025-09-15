package cli

import (
	"github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"fundchain/x/milestones/types"
)

func NewTreasuryBalanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "treasury-balance",
		Short: "Show the module treasury balance",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			qc := types.NewQueryClient(clientCtx)
			resp, err := qc.TreasuryBalance(cmd.Context(), &types.QueryTreasuryBalanceRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
