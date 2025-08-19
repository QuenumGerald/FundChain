package cli

import (
    "strconv"

    "github.com/spf13/cobra"

    "github.com/cosmos/cosmos-sdk/client"
    "github.com/cosmos/cosmos-sdk/client/flags"
    sdkquery "github.com/cosmos/cosmos-sdk/types/query"

    "fundchain/x/milestones/types"
)

func NewProjectsCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "projects",
        Short: "List projects [--limit --offset]",
        Args:  cobra.NoArgs,
        RunE: func(cmd *cobra.Command, args []string) error {
            clientCtx, err := client.GetClientQueryContext(cmd)
            if err != nil {
                return err
            }

            limitStr, _ := cmd.Flags().GetString("limit")
            offsetStr, _ := cmd.Flags().GetString("offset")
            var limit, offset uint64
            if limitStr != "" {
                if v, err := strconv.ParseUint(limitStr, 10, 64); err == nil {
                    limit = v
                } else {
                    return err
                }
            }
            if offsetStr != "" {
                if v, err := strconv.ParseUint(offsetStr, 10, 64); err == nil {
                    offset = v
                } else {
                    return err
                }
            }

            qc := types.NewQueryClient(clientCtx)
            resp, err := qc.Projects(cmd.Context(), &types.QueryProjectsRequest{
                Pagination: &sdkquery.PageRequest{Limit: limit, Offset: offset},
            })
            if err != nil {
                return err
            }
            return clientCtx.PrintProto(resp)
        },
    }

    // Pagination flags compatibility (no-op until implemented)
    flags.AddQueryFlagsToCmd(cmd)
    cmd.Flags().String("limit", "", "pagination limit")
    cmd.Flags().String("offset", "", "pagination offset")

    return cmd
}
