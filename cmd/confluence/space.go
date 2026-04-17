package main

import (
	"github.com/spf13/cobra"

	"github.com/pankaj28843/confluence-cli/internal/conf"
)

func spaceCmdReal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "space",
		Short: "Spaces (list, view)",
		Long: `Space operations.

Examples:
  confluence space list
  confluence space list --type global --status current --json
  confluence space view ENG`,
	}
	cmd.AddCommand(spaceListCmd())
	cmd.AddCommand(spaceViewCmd())
	return cmd
}

func spaceListCmd() *cobra.Command {
	var typ, status string
	var limit int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List spaces in the site",
		Long: `List spaces with optional filters.

Examples:
  confluence space list --json --limit 100
  confluence space list --type global --status current`,
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
			spaces, err := conf.ListSpaces(ctx, c, conf.SpaceFilter{Type: typ, Status: status, Limit: limit})
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(spaces)
			}
			for _, s := range spaces {
				w.Text("%s\t%s\t%s\n", s.Key, s.Type, s.Name)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&typ, "type", "", "Filter: global | personal")
	cmd.Flags().StringVar(&status, "status", "", "Filter: current | archived")
	cmd.Flags().IntVar(&limit, "limit", 25, "Max spaces to return (hard cap 200)")
	return cmd
}

func spaceViewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "view <key>",
		Short: "Show one space",
		Long: `Show a space by key.

Examples:
  confluence space view ENG
  confluence space view ENG --json`,
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
			s, err := conf.GetSpace(ctx, c, args[0])
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(s)
			}
			w.Text("%s  %s\n  type: %s\n  status: %s\n", s.Key, s.Name, s.Type, s.Status)
			if s.Description.Plain.Value != "" {
				w.Text("  description: %s\n", s.Description.Plain.Value)
			}
			return nil
		},
	}
}
