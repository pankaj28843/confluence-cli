package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pankaj28843/confluence-cli/internal/conf"
)

func userCmdReal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Users (current, view, search, bulk)",
		Long: `User operations.

Examples:
  confluence user current
  confluence user view --username alice                  # Server/DC
  confluence user view --account-id 557058:abc...         # Cloud
  confluence user search "Jane Smith"
  confluence user bulk --account-id 557058:abc --json`,
	}
	cmd.AddCommand(userCurrentCmd())
	cmd.AddCommand(userViewCmd())
	cmd.AddCommand(userSearchCmd())
	cmd.AddCommand(userBulkCmd())
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
		Short: "Show a user by username / key / account id",
		Long: `Show a user. Pick one selector. On Server/DC use --username or --key;
on Cloud use --account-id.

Examples:
  confluence user view --username alice
  confluence user view --account-id 557058:abc
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
				return fmt.Errorf("provide one of --username, --key, --account-id")
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
	cmd.Flags().StringVar(&accountID, "account-id", "", "Cloud account id")
	cmd.Flags().StringVar(&accountID, "accountId", "", "Deprecated alias for --account-id")
	_ = cmd.Flags().MarkHidden("accountId")
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

func userBulkCmd() *cobra.Command {
	var accountIDs []string
	cmd := &cobra.Command{
		Use:   "bulk",
		Short: "Show Cloud users by account id",
		Long: `Show Cloud users by account id using the documented v2 users-bulk endpoint.

Examples:
  confluence user bulk --account-id 557058:abc
  confluence user bulk --account-id 557058:abc --account-id 557058:def --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			if len(accountIDs) == 0 {
				return fmt.Errorf("--account-id is required")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			users, err := conf.BulkGetUsers(ctx, c, accountIDs)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(users)
			}
			for _, u := range users {
				w.Text("%s\t%s\t%s\n",
					firstNonEmpty(u.AccountID, u.Username, u.UserKey),
					u.Email,
					firstNonEmpty(u.DisplayName, u.PublicName))
			}
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "Cloud account id; repeatable or comma-separated")
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
