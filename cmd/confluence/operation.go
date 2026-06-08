package main

import (
	"github.com/spf13/cobra"

	"github.com/pankaj28843/confluence-cli/internal/conf"
)

func operationCmdReal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "operation",
		Short: "Permitted operations (list)",
		Long: `Permitted operation helpers.

Cloud supports pages, blog posts, attachments, spaces, comments, and newer
content-tree entities. Server/Data Center supports content ids through the
documented operations expansion.

Examples:
  confluence operation list --page 12345
  confluence operation list --space ENG --json`,
	}
	cmd.AddCommand(operationListCmd())
	return cmd
}

func operationListCmd() *cobra.Command {
	var target entityTargetFlags
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List permitted operations for one entity",
		Long: `List permitted operations for one Confluence entity.

Examples:
  confluence operation list --page 12345
  confluence operation list --blogpost 67890 --json
  confluence operation list --space ENG`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			selected, err := selectedOperationTarget(target)
			if err != nil {
				return err
			}
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			c, err := newClient()
			if err != nil {
				return err
			}
			ops, err := conf.ListOperations(ctx, c, selected.typ, selected.id)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(ops)
			}
			for _, op := range ops {
				w.Text("%s\t%s\n", op.Operation, op.TargetType)
			}
			return nil
		},
	}
	addOperationTargetFlags(cmd, &target)
	return cmd
}
