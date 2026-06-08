package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pankaj28843/confluence-cli/internal/client"
	"github.com/pankaj28843/confluence-cli/internal/conf"
	"github.com/pankaj28843/confluence-cli/internal/output"
)

func groupCmdReal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "group",
		Short: "Groups (list, view, picker, members, hierarchy)",
		Long: `Group operations.

Examples:
  confluence group list
  confluence group view engineering
  confluence group view --id 11111111-2222-3333-4444-555555555555
  confluence group picker eng
  confluence group members engineering`,
	}
	cmd.AddCommand(groupListCmd())
	cmd.AddCommand(groupViewCmd())
	cmd.AddCommand(groupPickerCmd())
	cmd.AddCommand(groupMembersCmd())
	cmd.AddCommand(groupRelationCmd("children", "List Server/DC child groups", conf.ListGroupChildGroups))
	cmd.AddCommand(groupRelationCmd("parents", "List Server/DC parent groups", conf.ListGroupParents))
	cmd.AddCommand(groupRelationCmd("ancestors", "List Server/DC ancestor groups", conf.ListGroupAncestors))
	return cmd
}

func groupListCmd() *cobra.Command {
	var limit int
	var accessType string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List groups",
		Long: `List groups.

Examples:
  confluence group list
  confluence group list --json --limit 200
  confluence group list --access-type admin --json`,
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
			gs, err := conf.ListGroupsWithOptions(ctx, c, conf.GroupListOptions{
				Limit:      limit,
				AccessType: accessType,
			})
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
	cmd.Flags().IntVar(&limit, "limit", 25, "Max groups (hard cap 200)")
	cmd.Flags().StringVar(&accessType, "access-type", "", "Cloud access type filter: user, admin, or site-admin")
	return cmd
}

func groupViewCmd() *cobra.Command {
	var id string
	var expand string
	cmd := &cobra.Command{
		Use:   "view [name]",
		Short: "Show a group",
		Long: `Show a group. On Server/DC pass the group name argument. On Cloud pass --id.

Examples:
  confluence group view engineering
  confluence group view engineering --expand members
  confluence group view --id 11111111-2222-3333-4444-555555555555 --json`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			name := ""
			if len(args) > 0 {
				name = args[0]
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			group, err := conf.GetGroup(ctx, c, conf.GroupLookupOptions{
				ID:     id,
				Name:   name,
				Expand: expand,
			})
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(group)
			}
			printGroup(w, *group)
			return nil
		},
	}
	cmd.Flags().StringVar(&id, "id", "", "Cloud group id")
	cmd.Flags().StringVar(&expand, "expand", "", "Server/DC expand value")
	return cmd
}

func groupPickerCmd() *cobra.Command {
	var limit int
	var totalSize bool
	cmd := &cobra.Command{
		Use:   "picker <query>",
		Short: "Search Cloud groups by picker query",
		Long: `Search Cloud groups using the documented group picker endpoint.

Examples:
  confluence group picker eng
  confluence group picker eng --limit 50 --total-size --json`,
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
			gs, err := conf.PickGroups(ctx, c, args[0], conf.GroupPickerOptions{
				Limit:                 limit,
				ShouldReturnTotalSize: totalSize,
			})
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(gs)
			}
			printGroups(w, gs)
			return nil
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 25, "Max groups (hard cap 200)")
	cmd.Flags().BoolVar(&totalSize, "total-size", false, "Ask Cloud to include total size metadata")
	return cmd
}

func groupMembersCmd() *cobra.Command {
	var id string
	var limit int
	var expand []string
	var totalSize bool
	cmd := &cobra.Command{
		Use:   "members [name]",
		Short: "List members of a group",
		Long: `List group members. On Server/DC pass the group name argument. On Cloud pass --id.

Examples:
  confluence group members engineering
  confluence group members engineering --json --limit 200
  confluence group members --id 11111111-2222-3333-4444-555555555555 --expand personalSpace --json`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			name := ""
			if len(args) > 0 {
				name = args[0]
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			us, err := conf.ListGroupMembersWithOptions(ctx, c, conf.GroupMemberOptions{
				GroupID:               id,
				GroupName:             name,
				Limit:                 limit,
				Expand:                expand,
				ShouldReturnTotalSize: totalSize,
			})
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
	cmd.Flags().StringVar(&id, "id", "", "Cloud group id")
	cmd.Flags().IntVar(&limit, "limit", 25, "Max members (hard cap 200)")
	cmd.Flags().StringSliceVar(&expand, "expand", nil, "Expand value; repeatable or comma-separated")
	cmd.Flags().BoolVar(&totalSize, "total-size", false, "Ask Cloud to include total size metadata")
	return cmd
}

func groupRelationCmd(use, short string, run func(ctx context.Context, c *client.Client, opts conf.GroupRelationOptions) ([]conf.Group, error)) *cobra.Command {
	var limit int
	var expand string
	cmd := &cobra.Command{
		Use:   use + " <name>",
		Short: short,
		Long: fmt.Sprintf(`%s for a Server/Data Center group.

Examples:
  confluence group %s engineering
  confluence group %s engineering --expand members --limit 50 --json`, short, use, use),
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
			gs, err := run(ctx, c, conf.GroupRelationOptions{
				GroupName: args[0],
				Limit:     limit,
				Expand:    expand,
			})
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(gs)
			}
			printGroups(w, gs)
			return nil
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 25, "Max groups (hard cap 200)")
	cmd.Flags().StringVar(&expand, "expand", "", "Server/DC expand value")
	return cmd
}

func printGroups(w *output.Writer, groups []conf.Group) {
	for _, group := range groups {
		printGroup(w, group)
	}
}

func printGroup(w *output.Writer, group conf.Group) {
	w.Text("%s\t%s\t%s\n",
		firstNonEmpty(group.ID, "-"),
		group.Name,
		firstNonEmpty(group.Type, group.UsageType, "-"))
}
