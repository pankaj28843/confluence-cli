package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pankaj28843/confluence-cli/internal/conf"
)

func labelCmdReal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "label",
		Short: "Content and space labels",
		Long: `Label operations.

Examples:
  confluence label list --page 12345 --json
  confluence label space --space ENG --limit 50
  confluence label search --prefix global --json
  confluence label add --page 12345 --label needs-review,shipped
  confluence label remove --page 12345 --label needs-review`,
	}
	cmd.AddCommand(labelListCmd())
	cmd.AddCommand(labelSpaceCmd())
	cmd.AddCommand(labelSearchCmd())
	cmd.AddCommand(labelRecentCmd())
	cmd.AddCommand(labelRelatedCmd())
	cmd.AddCommand(labelAddCmd())
	cmd.AddCommand(labelRemoveCmd())
	return cmd
}

func labelListCmd() *cobra.Command {
	var target entityTargetFlags
	var limit int
	var prefix string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List labels on a content target",
		Long: `List labels.

Examples:
  confluence label list --page 12345 --json
  confluence label list --blogpost 67890 --prefix global
  confluence label list --attachment att123 --limit 100`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			selected, err := selectedLabelTarget(target)
			if err != nil {
				return err
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			labels, err := conf.ListTargetLabels(ctx, c, conf.LabelTarget{Type: selected.typ, ID: selected.id}, conf.LabelListOptions{Limit: limit, Prefix: prefix})
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(labels)
			}
			return writeLabels(w, labels)
		},
	}
	addLabelTargetFlags(cmd, &target)
	cmd.Flags().IntVar(&limit, "limit", 25, "Max labels (hard cap 200)")
	cmd.Flags().StringVar(&prefix, "prefix", "", "Label prefix filter")
	return cmd
}

func labelSpaceCmd() *cobra.Command {
	var space, prefix, scope string
	var limit int
	cmd := &cobra.Command{
		Use:   "space",
		Short: "List labels used in a space",
		Long: `List labels used in a space.

Examples:
  confluence label space --space ENG --json
  confluence label space --space ENG --prefix global --limit 100
  confluence label space --space ENG --scope space --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			if space == "" {
				return fmt.Errorf("--space is required")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			labels, err := conf.ListSpaceLabels(ctx, c, space, conf.LabelListOptions{Limit: limit, Prefix: prefix, Scope: scope})
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(labels)
			}
			return writeLabels(w, labels)
		},
	}
	cmd.Flags().StringVar(&space, "space", "", "Space key or Cloud space id (required)")
	cmd.Flags().StringVar(&prefix, "prefix", "", "Label prefix filter")
	cmd.Flags().StringVar(&scope, "scope", "content", "Cloud only: content or space")
	cmd.Flags().IntVar(&limit, "limit", 25, "Max labels (hard cap 200)")
	return cmd
}

func labelSearchCmd() *cobra.Command {
	var labelIDs []string
	var prefixes []string
	var sort string
	var limit int
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search the Cloud label catalog",
		Long: `Search the Cloud label catalog.

Examples:
  confluence label search --json
  confluence label search --prefix global --limit 100
  confluence label search --label-id 123 --json`,
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
			labels, err := conf.SearchLabels(ctx, c, conf.LabelSearchOptions{
				Limit:    limit,
				LabelIDs: labelIDs,
				Prefixes: prefixes,
				Sort:     sort,
			})
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(labels)
			}
			return writeLabels(w, labels)
		},
	}
	cmd.Flags().StringSliceVar(&labelIDs, "label-id", nil, "Cloud label id filter; repeatable")
	cmd.Flags().StringSliceVar(&prefixes, "prefix", nil, "Cloud label prefix filter; repeatable")
	cmd.Flags().StringVar(&sort, "sort", "", "Cloud sort expression")
	cmd.Flags().IntVar(&limit, "limit", 25, "Max labels (hard cap 200)")
	return cmd
}

func labelRecentCmd() *cobra.Command {
	var limit int
	cmd := &cobra.Command{
		Use:   "recent",
		Short: "List recently used Server/Data Center labels",
		Long: `List recently used Server/Data Center labels.

Examples:
  confluence label recent --json
  confluence label recent --limit 100`,
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
			labels, err := conf.ListRecentLabels(ctx, c, limit)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(labels)
			}
			return writeLabels(w, labels)
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 25, "Max labels (hard cap 200)")
	return cmd
}

func labelRelatedCmd() *cobra.Command {
	var label, space string
	var limit int
	cmd := &cobra.Command{
		Use:   "related",
		Short: "List related Server/Data Center labels",
		Long: `List related Server/Data Center labels.

Examples:
  confluence label related --label incident --json
  confluence label related --space ENG --label incident --limit 100`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			if label == "" {
				return fmt.Errorf("--label is required")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			labels, err := conf.ListRelatedLabels(ctx, c, conf.LabelRelatedOptions{SpaceKey: space, Label: label, Limit: limit})
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(labels)
			}
			return writeLabels(w, labels)
		},
	}
	cmd.Flags().StringVar(&label, "label", "", "Label name (required)")
	cmd.Flags().StringVar(&space, "space", "", "Server/Data Center space key")
	cmd.Flags().IntVar(&limit, "limit", 25, "Max labels (hard cap 200)")
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

type labelTextWriter interface {
	Text(format string, args ...any)
}

func writeLabels(w labelTextWriter, labels []conf.Label) error {
	for _, l := range labels {
		w.Text("%s\t%s\n", l.Prefix, l.Name)
	}
	return nil
}
