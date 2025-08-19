package cli

import (
	"github.com/spf13/cobra"

	client "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"fundchain/x/milestones/types"
)

func NewParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Show module parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			qc := types.NewQueryClient(clientCtx)
			resp, err := qc.Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
