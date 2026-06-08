package main

import (
	"github.com/spf13/cobra"

	"github.com/pankaj28843/confluence-cli/internal/conf"
)

func analyticsCmdReal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "analytics",
		Short: "Cloud content analytics reads",
		Long: `Cloud content analytics read operations.

Analytics view and viewer counts are documented in Confluence Cloud REST API v1.
Server/Data Center does not expose the same Analytics REST group in the current
official REST reference, so these typed commands are Cloud-only.

Examples:
  confluence analytics views 12345 --json
  confluence analytics viewers 12345 --from-date YYYY-MM-DDTHH:MM:SS.sssZ`,
	}
	cmd.AddCommand(analyticsViewsCmd())
	cmd.AddCommand(analyticsViewersCmd())
	return cmd
}

func analyticsViewsCmd() *cobra.Command {
	var fromDate string
	cmd := &cobra.Command{
		Use:   "views <content-id>",
		Short: "Show Cloud content view count",
		Long: `Show the total number of Cloud views for one content item.

Examples:
  confluence analytics views 12345
  confluence analytics views 12345 --from-date YYYY-MM-DDTHH:MM:SS.sssZ --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			c, err := newClient()
			if err != nil {
				return err
			}
			count, err := conf.GetContentViewCount(ctx, c, args[0], conf.AnalyticsOptions{FromDate: fromDate})
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(count)
			}
			printAnalyticsCount(w, *count)
			return nil
		},
	}
	cmd.Flags().StringVar(&fromDate, "from-date", "", "Cloud analytics start date, e.g. YYYY-MM-DDTHH:MM:SS.sssZ")
	return cmd
}

func analyticsViewersCmd() *cobra.Command {
	var fromDate string
	cmd := &cobra.Command{
		Use:   "viewers <content-id>",
		Short: "Show Cloud distinct viewer count",
		Long: `Show the total number of distinct Cloud viewers for one content item.

Examples:
  confluence analytics viewers 12345
  confluence analytics viewers 12345 --from-date YYYY-MM-DDTHH:MM:SS.sssZ --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			c, err := newClient()
			if err != nil {
				return err
			}
			count, err := conf.GetContentViewerCount(ctx, c, args[0], conf.AnalyticsOptions{FromDate: fromDate})
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(count)
			}
			printAnalyticsCount(w, *count)
			return nil
		},
	}
	cmd.Flags().StringVar(&fromDate, "from-date", "", "Cloud analytics start date, e.g. YYYY-MM-DDTHH:MM:SS.sssZ")
	return cmd
}

type analyticsTextWriter interface {
	Text(format string, args ...any)
}

func printAnalyticsCount(w analyticsTextWriter, count conf.AnalyticsCount) {
	w.Text("%d\t%d\n", count.ID, count.Count)
}
