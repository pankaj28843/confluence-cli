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
		Long: `Restriction operations. Output is passed through verbatim — ACE shape is
Confluence-version specific.

Examples:
  confluence restriction list --page 12345 --json`,
	}
	var page string
	list := &cobra.Command{
		Use:   "list",
		Short: "List read/update restrictions on a content id",
		Long: `List restrictions.

Examples:
  confluence restriction list --page 12345 --json`,
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
			data, err := conf.GetContentRestrictions(ctx, c, page)
			if err != nil {
				return err
			}
			_, _ = w.Out.Write(data)
			return nil
		},
	}
	list.Flags().StringVar(&page, "page", "", "Content id (required)")
	cmd.AddCommand(list)
	return cmd
}
