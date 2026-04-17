package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pankaj28843/confluence-cli/internal/conf"
)

func labelCmdReal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "label",
		Short: "Content labels (list, add, remove)",
		Long: `Label operations.

Examples:
  confluence label list --page 12345
  confluence label add --page 12345 --label needs-review,shipped
  confluence label remove --page 12345 --label needs-review`,
	}
	cmd.AddCommand(labelListCmd())
	cmd.AddCommand(labelAddCmd())
	cmd.AddCommand(labelRemoveCmd())
	return cmd
}

func labelListCmd() *cobra.Command {
	var page string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List labels on a content id",
		Long: `List labels.

Examples:
  confluence label list --page 12345 --json`,
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
			labels, err := conf.ListLabels(ctx, c, page)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(labels)
			}
			for _, l := range labels {
				w.Text("%s\t%s\n", l.Prefix, l.Name)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&page, "page", "", "Content id (required)")
	return cmd
}

func labelAddCmd() *cobra.Command {
	var page string
	var names []string
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add one or more labels to a content id",
		Long: `Add labels. Pass multiple --label flags or one comma-separated list.

Examples:
  confluence label add --page 12345 --label needs-review
  confluence label add --page 12345 --label review,shipped,v1`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			if page == "" || len(names) == 0 {
				return fmt.Errorf("--page and at least one --label are required")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			out, err := conf.AddLabels(ctx, c, page, names)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(out)
			}
			for _, l := range out {
				w.Text("%s\t%s\n", l.Prefix, l.Name)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&page, "page", "", "Content id (required)")
	cmd.Flags().StringSliceVar(&names, "label", nil, "Label name(s); repeatable")
	return cmd
}

func labelRemoveCmd() *cobra.Command {
	var page, name string
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove one label from a content id",
		Long: `Remove a label.

Examples:
  confluence label remove --page 12345 --label needs-review`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			if page == "" || name == "" {
				return fmt.Errorf("--page and --label are required")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			if err := conf.RemoveLabel(ctx, c, page, name); err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(map[string]any{"ok": true, "page": page, "label": name})
			}
			w.Text("removed label %q from page %s\n", name, page)
			return nil
		},
	}
	cmd.Flags().StringVar(&page, "page", "", "Content id (required)")
	cmd.Flags().StringVar(&name, "label", "", "Label name to remove (required)")
	return cmd
}
