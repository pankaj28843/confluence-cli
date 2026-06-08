package main

import (
	"github.com/spf13/cobra"

	"github.com/pankaj28843/confluence-cli/internal/conf"
)

func likeCmdReal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "like",
		Short: "Cloud likes (count, users)",
		Long: `Cloud like helpers for pages, blog posts, footer comments, and inline comments.

Examples:
  confluence like count --page 12345
  confluence like users --blogpost 67890 --limit 50 --json`,
	}
	cmd.AddCommand(likeCountCmd())
	cmd.AddCommand(likeUsersCmd())
	return cmd
}

func likeCountCmd() *cobra.Command {
	var target entityTargetFlags
	cmd := &cobra.Command{
		Use:   "count",
		Short: "Show Cloud like count for one entity",
		Long: `Show the Cloud like count for one page, blog post, footer comment, or inline comment.

Examples:
  confluence like count --page 12345
  confluence like count --footer-comment 67890 --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			selected, err := selectedLikeTarget(target)
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
			count, err := conf.GetLikeCount(ctx, c, selected.typ, selected.id)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(count)
			}
			w.Text("%d\n", count.Count)
			return nil
		},
	}
	addLikeTargetFlags(cmd, &target)
	return cmd
}

func likeUsersCmd() *cobra.Command {
	var target entityTargetFlags
	var limit int
	cmd := &cobra.Command{
		Use:   "users",
		Short: "List Cloud account ids that liked one entity",
		Long: `List Cloud account ids that liked one page, blog post, footer comment, or inline comment.

Examples:
  confluence like users --page 12345
  confluence like users --inline-comment 67890 --limit 50 --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			selected, err := selectedLikeTarget(target)
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
			users, err := conf.ListLikeUsers(ctx, c, selected.typ, selected.id, limit)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(users)
			}
			for _, user := range users {
				w.Text("%s\n", user.AccountID)
			}
			return nil
		},
	}
	addLikeTargetFlags(cmd, &target)
	cmd.Flags().IntVar(&limit, "limit", 25, "Max users (hard cap 200)")
	return cmd
}
