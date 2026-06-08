package main

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/pankaj28843/confluence-cli/internal/conf"
)

func commentCmdReal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "comment",
		Short: "Comments (list/versions/version/add/update/delete)",
		Long: `Comment operations.

Examples:
  confluence comment list --page 12345
  confluence comment add --page 12345 --body "<p>Looks good.</p>"
  confluence comment update 998877 --body-file comment.html
  confluence comment delete 998877 --force
  confluence comment list --page 12345 --locations footer,inline
  confluence comment list --page 12345 --json --limit 50`,
	}
	var page string
	var limit int
	var locations []string
	list := &cobra.Command{
		Use:   "list",
		Short: "List comments (footer, inline, resolved by default)",
		Long: `List comments on a content id.

Examples:
  confluence comment list --page 12345
  confluence comment list --page 12345 --locations footer
  confluence comment list --page 12345 --json --limit 100`,
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
			cs, err := conf.ListComments(ctx, c, page, locations, limit)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(cs)
			}
			for _, co := range cs {
				resolved := ""
				if co.Resolved {
					resolved = " [resolved]"
				}
				w.Text("[%s] %s (%s)%s\n", co.Location, co.Author, co.Date, resolved)
				if co.InlineOriginalSelection != "" {
					w.Text("  > %s\n", co.InlineOriginalSelection)
				}
				w.Text("  %s\n\n", co.Body)
			}
			return nil
		},
	}
	list.Flags().StringVar(&page, "page", "", "Content id (required)")
	list.Flags().StringSliceVar(&locations, "locations", nil, "Subset of footer,inline,resolved")
	list.Flags().IntVar(&limit, "limit", 100, "Max comments (hard cap 200)")
	cmd.AddCommand(list)
	cmd.AddCommand(commentVersionsCmd())
	cmd.AddCommand(commentVersionCmd())
	cmd.AddCommand(commentAddCmd())
	cmd.AddCommand(commentUpdateCmd())
	cmd.AddCommand(commentDeleteCmd())
	return cmd
}

func commentVersionsCmd() *cobra.Command {
	return versionListReadCmd("comment", "footer-comment", true, true)
}

func commentVersionCmd() *cobra.Command {
	return versionDetailReadCmd("comment", "footer-comment", true)
}

func commentAddCmd() *cobra.Command {
	var page, blogpost, parent, bodyFormat, bodyFile, bodyInline string
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a footer comment",
		Long: `Add a footer comment to a page or blog post. On Confluence Cloud, --parent
creates a reply to an existing footer comment.

Examples:
  confluence comment add --page 12345 --body "<p>Looks good.</p>"
  confluence comment add --blogpost 2001 --body-file comment.html
  echo "<p>Reply</p>" | confluence comment add --parent 998877 --body-file - --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()

			body, err := resolveBody(bodyFile, bodyInline)
			if err != nil {
				return err
			}
			if body == "" {
				return fmt.Errorf("--body-file or --body is required")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			out, err := conf.CreateFooterComment(ctx, c, conf.CommentInput{
				PageID:          page,
				BlogPostID:      blogpost,
				ParentCommentID: parent,
				BodyFormat:      bodyFormat,
				BodyValue:       body,
			})
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(out)
			}
			if parent != "" {
				w.Text("created reply %s", out.ID)
			} else {
				w.Text("created comment %s", out.ID)
			}
			if out.VersionNumber > 0 {
				w.Text(" (v%d)", out.VersionNumber)
			}
			w.Text("\n")
			return nil
		},
	}
	cmd.Flags().StringVar(&page, "page", "", "Page id to comment on")
	cmd.Flags().StringVar(&blogpost, "blogpost", "", "Blog post id to comment on")
	cmd.Flags().StringVar(&parent, "parent", "", "Cloud only: parent footer comment id for a reply")
	cmd.Flags().StringVar(&bodyFormat, "body-format", "storage", "Body format: storage")
	cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to body file, or '-' for stdin")
	cmd.Flags().StringVar(&bodyInline, "body", "", "Inline body string")
	return cmd
}

func commentUpdateCmd() *cobra.Command {
	var bodyFormat, bodyFile, bodyInline, version string
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a footer comment body",
		Long: `Update a footer comment body. The command fetches the current comment version
and writes version.number + 1 unless --version is supplied.

Examples:
  confluence comment update 998877 --body "<p>Updated.</p>"
  confluence comment update 998877 --body-file comment.html
  echo "<p>Updated.</p>" | confluence comment update 998877 --body-file - --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()

			body, err := resolveBody(bodyFile, bodyInline)
			if err != nil {
				return err
			}
			if body == "" {
				return fmt.Errorf("--body-file or --body is required")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			ver := 0
			if version != "" {
				ver, err = strconv.Atoi(version)
				if err != nil || ver < 0 {
					return fmt.Errorf("invalid --version %q", version)
				}
			} else {
				current, err := conf.GetFooterComment(ctx, c, args[0])
				if err != nil {
					return err
				}
				ver = current.VersionNumber
			}
			out, err := conf.UpdateFooterComment(ctx, c, conf.CommentInput{
				ID:            args[0],
				BodyFormat:    bodyFormat,
				BodyValue:     body,
				VersionNumber: ver,
			})
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(out)
			}
			w.Text("updated comment %s", out.ID)
			if out.VersionNumber > 0 {
				w.Text(" (v%d)", out.VersionNumber)
			}
			w.Text("\n")
			return nil
		},
	}
	cmd.Flags().StringVar(&bodyFormat, "body-format", "storage", "Body format: storage")
	cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to body file, or '-' for stdin")
	cmd.Flags().StringVar(&bodyInline, "body", "", "Inline body string")
	cmd.Flags().StringVar(&version, "version", "", "Explicit current version number (auto-fetched if omitted)")
	return cmd
}

func commentDeleteCmd() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a footer comment",
		Long: `Delete a footer comment permanently. A confirmation prompt is shown unless
--force is supplied.

Examples:
  confluence comment delete 998877
  confluence comment delete 998877 --force
  confluence comment delete 998877 --force --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()

			if !force && !confirmDelete(args[0], "comment") {
				return fmt.Errorf("delete cancelled")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			if err := conf.DeleteFooterComment(ctx, c, args[0]); err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(map[string]any{"deleted": true, "id": args[0]})
			}
			w.Text("deleted comment %s\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "Do not prompt for confirmation")
	return cmd
}
