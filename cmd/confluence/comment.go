package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pankaj28843/confluence-cli/internal/conf"
)

func commentCmdReal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "comment",
		Short: "Comments (list inline/footer/resolved)",
		Long: `Comment operations.

Examples:
  confluence comment list --page 12345
  confluence comment list --page 12345 --locations footer,inline
  confluence comment list --page 12345 --json --limit 50`,
	}
	var page string
	var limit int
	var locations []string
	list := &cobra.Command{
		Use:   "list",
		Short: "List comments (footer, inline, resolved by default)",
		Long: `List comments on a content id.

Examples:
  confluence comment list --page 12345
  confluence comment list --page 12345 --locations footer
  confluence comment list --page 12345 --json --limit 100`,
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
			cs, err := conf.ListComments(ctx, c, page, locations, limit)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(cs)
			}
			for _, co := range cs {
				resolved := ""
				if co.Resolved {
					resolved = " [resolved]"
				}
				w.Text("[%s] %s (%s)%s\n", co.Location, co.Author, co.Date, resolved)
				if co.InlineOriginalSelection != "" {
					w.Text("  > %s\n", co.InlineOriginalSelection)
				}
				w.Text("  %s\n\n", co.Body)
			}
			return nil
		},
	}
	list.Flags().StringVar(&page, "page", "", "Content id (required)")
	list.Flags().StringSliceVar(&locations, "locations", nil, "Subset of footer,inline,resolved")
	list.Flags().IntVar(&limit, "limit", 100, "Max comments (hard cap 200)")
	cmd.AddCommand(list)
	return cmd
}
