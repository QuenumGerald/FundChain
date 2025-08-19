package cli

import (
    "strconv"

    "github.com/spf13/cobra"

    "github.com/cosmos/cosmos-sdk/client"
    "github.com/cosmos/cosmos-sdk/client/flags"
    sdkquery "github.com/cosmos/cosmos-sdk/types/query"
    "fundchain/x/milestones/types"
)

func NewMilestonesCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "milestones [projectId]",
        Short: "List milestones for project [--limit --offset]",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            clientCtx, err := client.GetClientQueryContext(cmd)
            if err != nil {
                return err
            }
            projectID, err := strconv.ParseUint(args[0], 10, 64)
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
            resp, err := qc.ProjectMilestones(cmd.Context(), &types.QueryProjectMilestonesRequest{
                ProjectId:  projectID,
                Pagination: &sdkquery.PageRequest{Limit: limit, Offset: offset},
            })
            if err != nil {
                return err
            }
            return clientCtx.PrintProto(resp)
        },
    }

    flags.AddQueryFlagsToCmd(cmd)
    cmd.Flags().String("limit", "", "pagination limit")
    cmd.Flags().String("offset", "", "pagination offset")
    return cmd
}
