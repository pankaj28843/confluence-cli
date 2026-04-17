package main

import (
	"github.com/spf13/cobra"

	"github.com/pankaj28843/confluence-cli/internal/conf"
)

func groupCmdReal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "group",
		Short: "Groups (list, members)",
		Long: `Group operations.

Examples:
  confluence group list
  confluence group members engineering`,
	}
	var limit int
	list := &cobra.Command{
		Use:   "list",
		Short: "List groups",
		Long: `List groups.

Examples:
  confluence group list
  confluence group list --json --limit 200`,
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
			gs, err := conf.ListGroups(ctx, c, limit)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(gs)
			}
			for _, g := range gs {
				w.Text("%s\n", g.Name)
			}
			return nil
		},
	}
	list.Flags().IntVar(&limit, "limit", 50, "Max groups (hard cap 200)")
	cmd.AddCommand(list)

	members := &cobra.Command{
		Use:   "members <name>",
		Short: "List members of a group",
		Long: `List group members.

Examples:
  confluence group members engineering
  confluence group members engineering --json --limit 200`,
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
			us, err := conf.ListGroupMembers(ctx, c, args[0], limit)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(us)
			}
			for _, u := range us {
				w.Text("%s\t%s\t%s\n",
					firstNonEmpty(u.Username, u.AccountID, u.UserKey),
					u.Email,
					firstNonEmpty(u.DisplayName, u.PublicName))
			}
			return nil
		},
	}
	members.Flags().IntVar(&limit, "limit", 100, "Max members (hard cap 200)")
	cmd.AddCommand(members)
	return cmd
}
