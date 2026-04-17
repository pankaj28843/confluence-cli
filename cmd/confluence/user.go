package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pankaj28843/confluence-cli/internal/conf"
)

func userCmdReal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Users (current, view, search)",
		Long: `User operations.

Examples:
  confluence user current
  confluence user view --username alice                  # Server/DC
  confluence user view --accountId 557058:abc...          # Cloud
  confluence user search "Jane Smith"`,
	}
	cmd.AddCommand(userCurrentCmd())
	cmd.AddCommand(userViewCmd())
	cmd.AddCommand(userSearchCmd())
	return cmd
}

func userCurrentCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "current",
		Short: "Show the authenticated user",
		Long: `Show the authenticated user (shortcut for GET /rest/api/user/current).

Examples:
  confluence user current
  confluence user current --json`,
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
			u, err := conf.GetCurrentUser(ctx, c)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(u)
			}
			w.Text("%s\t%s\n", u.ID(), u.Label())
			return nil
		},
	}
}

func userViewCmd() *cobra.Command {
	var username, key, accountID string
	cmd := &cobra.Command{
		Use:   "view",
		Short: "Show a user by username / key / accountId",
		Long: `Show a user. Pick one selector. On Server/DC use --username or --key;
on Cloud use --accountId.

Examples:
  confluence user view --username alice
  confluence user view --accountId 557058:abc
  confluence user view --key u1234`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()

			var sel, val string
			switch {
			case username != "":
				sel, val = "username", username
			case key != "":
				sel, val = "key", key
			case accountID != "":
				sel, val = "accountId", accountID
			default:
				return fmt.Errorf("provide one of --username, --key, --accountId")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			u, err := conf.GetUser(ctx, c, sel, val)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(u)
			}
			w.Text("%s\t%s\t%s\t%s\n",
				firstNonEmpty(u.Username, u.AccountID, u.UserKey),
				u.Type,
				u.Email,
				firstNonEmpty(u.DisplayName, u.PublicName))
			return nil
		},
	}
	cmd.Flags().StringVar(&username, "username", "", "Server/DC username")
	cmd.Flags().StringVar(&key, "key", "", "Server/DC user key")
	cmd.Flags().StringVar(&accountID, "accountId", "", "Cloud accountId")
	return cmd
}

func userSearchCmd() *cobra.Command {
	var limit int
	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search users by full name (CQL user.fullname~)",
		Long: `Search users using user.fullname~"<query>" CQL.

Examples:
  confluence user search "Jane Smith"
  confluence user search "Jane" --limit 50 --json`,
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
			hits, err := conf.SearchUsers(ctx, c, args[0], limit)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(hits)
			}
			for _, h := range hits {
				if h.User == nil {
					continue
				}
				u := h.User
				w.Text("%s\t%s\t%s\n",
					firstNonEmpty(u.Username, u.AccountID, u.UserKey),
					u.Email,
					firstNonEmpty(u.DisplayName, u.PublicName))
			}
			return nil
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 25, "Max results (hard cap 200)")
	return cmd
}

func firstNonEmpty(ss ...string) string {
	for _, s := range ss {
		if s != "" {
			return s
		}
	}
	return ""
}
