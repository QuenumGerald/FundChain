package cli

import (
    "strconv"

    "github.com/spf13/cobra"

    "github.com/cosmos/cosmos-sdk/client"
    "github.com/cosmos/cosmos-sdk/client/flags"
    "fundchain/x/milestones/types"
)

func NewProjectCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "project [id]",
        Short: "Show a project by id",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            clientCtx, err := client.GetClientQueryContext(cmd)
            if err != nil {
                return err
            }
            id, err := strconv.ParseUint(args[0], 10, 64)
            if err != nil {
                return err
            }
            qc := types.NewQueryClient(clientCtx)
            resp, err := qc.Project(cmd.Context(), &types.QueryProjectRequest{Id: id})
            if err != nil {
                return err
            }
            return clientCtx.PrintProto(resp)
        },
    }

    flags.AddQueryFlagsToCmd(cmd)
    return cmd
}
