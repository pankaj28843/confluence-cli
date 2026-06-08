package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pankaj28843/confluence-cli/internal/conf"
	"github.com/pankaj28843/confluence-cli/internal/output"
)

func auditCmdReal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "audit",
		Short: "Audit log reads",
		Long: `Audit log read operations.

Examples:
  confluence audit list --limit 25 --json
  confluence audit since --number 7 --unit DAYS --search group
  confluence audit retention --json`,
	}
	cmd.AddCommand(auditListCmd())
	cmd.AddCommand(auditSinceCmd())
	cmd.AddCommand(auditRetentionCmd())
	return cmd
}

func auditListCmd() *cobra.Command {
	var limit int
	var startDate string
	var endDate string
	var search string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List audit records",
		Long: `List audit records. Cloud supports date and search filters; Server/Data Center uses the documented deprecated read endpoint.

Examples:
  confluence audit list
  confluence audit list --start-date 1700000000000 --end-date 1700100000000 --search space --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			c, err := newClient()
			if err != nil {
				return err
			}
			records, err := conf.ListAuditRecords(ctx, c, conf.AuditListOptions{
				StartDate:    startDate,
				EndDate:      endDate,
				SearchString: search,
				Limit:        limit,
			})
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(records)
			}
			printAuditRecords(w, records)
			return nil
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 25, "Max records (hard cap 200)")
	cmd.Flags().StringVar(&startDate, "start-date", "", "Cloud start date as epoch milliseconds")
	cmd.Flags().StringVar(&endDate, "end-date", "", "Cloud end date as epoch milliseconds")
	cmd.Flags().StringVar(&search, "search", "", "Cloud audit search string")
	return cmd
}

func auditSinceCmd() *cobra.Command {
	var limit int
	var number int64
	var unit string
	var search string
	cmd := &cobra.Command{
		Use:   "since",
		Short: "List recent Cloud audit records",
		Long: `List Cloud audit records for a time period back from the current date.

Examples:
  confluence audit since --number 3 --unit MONTHS
  confluence audit since --number 7 --unit DAYS --search group --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			if number <= 0 {
				return fmt.Errorf("--number is required")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			records, err := conf.ListAuditRecords(ctx, c, conf.AuditListOptions{
				SinceNumber:  number,
				SinceUnit:    unit,
				SearchString: search,
				Limit:        limit,
			})
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(records)
			}
			printAuditRecords(w, records)
			return nil
		},
	}
	cmd.Flags().Int64Var(&number, "number", 0, "Cloud time period number")
	cmd.Flags().StringVar(&unit, "unit", "MONTHS", "Cloud time period unit, e.g. DAYS, WEEKS, MONTHS")
	cmd.Flags().IntVar(&limit, "limit", 25, "Max records (hard cap 200)")
	cmd.Flags().StringVar(&search, "search", "", "Cloud audit search string")
	return cmd
}

func auditRetentionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "retention",
		Short: "Show Cloud audit retention period",
		Long: `Show the Cloud audit retention period.

Examples:
  confluence audit retention
  confluence audit retention --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			c, err := newClient()
			if err != nil {
				return err
			}
			retention, err := conf.GetAuditRetention(ctx, c)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(retention)
			}
			w.Text("%d\t%s\n", retention.Number, retention.Units)
			return nil
		},
	}
	return cmd
}

func printAuditRecords(w *output.Writer, records []conf.AuditRecord) {
	for _, record := range records {
		w.Text("%d\t%s\t%s\t%s\n",
			record.CreationDate,
			firstNonEmpty(record.Category, "-"),
			firstNonEmpty(record.Author.DisplayName, record.Author.PublicName, record.Author.Username, record.Author.AccountID, "-"),
			firstNonEmpty(record.Summary, "-"))
	}
}
