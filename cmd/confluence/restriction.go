package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pankaj28843/confluence-cli/internal/conf"
)

func restrictionCmdReal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restriction",
		Short: "Content restrictions (list)",
		Long: `Content restriction helpers.

List restrictions grouped by operation, or inspect the read/update restrictions
for one content id.

Examples:
  confluence restriction list --page 12345
  confluence restriction list --page 12345 --operation read --json`,
	}
	var page string
	var operation string
	var limit int
	list := &cobra.Command{
		Use:   "list",
		Short: "List read/update restrictions on a content id",
		Long: `List read/update restrictions on a content id.

Examples:
  confluence restriction list --page 12345
  confluence restriction list --page 12345 --operation read
  confluence restriction list --page 12345 --operation update --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			if page == "" {
				return fmt.Errorf("--page is required")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			if operation != "" {
				restriction, err := conf.GetContentRestrictionForOperation(ctx, c, page, operation, limit)
				if err != nil {
					return err
				}
				if w.IsJSON() {
					return w.JSON(restriction)
				}
				w.Text("%s\tusers=%d\tgroups=%d\n", restriction.Operation, restriction.UserCount(), restriction.GroupCount())
				return nil
			}
			resp, err := conf.ListContentRestrictions(ctx, c, page)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(resp)
			}
			for _, restriction := range resp.Operations() {
				w.Text("%s\tusers=%d\tgroups=%d\n", restriction.Operation, restriction.UserCount(), restriction.GroupCount())
			}
			return nil
		},
	}
	list.Flags().StringVar(&page, "page", "", "Content id (required)")
	list.Flags().StringVar(&operation, "operation", "", "Restriction operation: read or update")
	list.Flags().IntVar(&limit, "limit", 25, "Maximum users/groups returned for --operation")
	cmd.AddCommand(list)
	return cmd
}
