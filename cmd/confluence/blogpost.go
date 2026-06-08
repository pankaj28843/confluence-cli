package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pankaj28843/confluence-cli/internal/conf"
)

func blogpostCmdReal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "blogpost",
		Short: "Blog posts (list, view, create, update, delete, purge)",
		Long: `Blog post operations.

Examples:
  confluence blogpost list --space ENG
  confluence blogpost view 12345 --markdown
  confluence blogpost create --space ENG --title "Weekly Update" --body-file body.html`,
	}
	cmd.AddCommand(blogpostListCmd())
	cmd.AddCommand(blogpostViewCmd())
	cmd.AddCommand(blogpostCreateCmd())
	cmd.AddCommand(blogpostUpdateCmd())
	cmd.AddCommand(blogpostDeleteCmd())
	cmd.AddCommand(blogpostPurgeCmd())
	return cmd
}

func blogpostListCmd() *cobra.Command {
	var space, labelID, title, status, postingDay string
	var limit int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List blog posts",
		Long: `List blog posts. On Cloud, --space resolves a space key to its v2 space id.
On Server/Data Center, --posting-day maps to the documented content resource
postingDay filter.

Examples:
  confluence blogpost list --space ENG
  confluence blogpost list --space ENG --title "Weekly" --limit 10 --json
  confluence blogpost list --posting-day 2026-06-08`,
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
			posts, err := conf.ListBlogPosts(ctx, c, conf.BlogPostListOptions{
				SpaceKey:   space,
				LabelID:    labelID,
				Title:      title,
				Status:     status,
				PostingDay: postingDay,
				Limit:      limit,
			})
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(posts)
			}
			for _, p := range posts {
				w.Text("%s\t%s\t%s\n", p.ID, p.Space.Key, p.Title)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&space, "space", "", "Space key filter")
	cmd.Flags().StringVar(&labelID, "label-id", "", "Cloud only: label id filter")
	cmd.Flags().StringVar(&title, "title", "", "Exact title filter")
	cmd.Flags().StringVar(&status, "status", "", "Status filter, e.g. current,draft,trashed")
	cmd.Flags().StringVar(&postingDay, "posting-day", "", "Server/Data Center only: posting day YYYY-MM-DD")
	cmd.Flags().IntVar(&limit, "limit", 25, "Max results (hard cap 200)")
	return cmd
}

func blogpostViewCmd() *cobra.Command {
	var markdown, rawStorage, bodyOnly bool
	cmd := &cobra.Command{
		Use:   "view <id>",
		Short: "Fetch a blog post by id",
		Long: `Fetch a blog post by id. --markdown renders the storage body to Markdown;
--raw-storage emits raw Confluence storage XML/HTML. --body-only emits only the
storage body for edit-and-reupload workflows.

Examples:
  confluence blogpost view 12345 --markdown
  confluence blogpost view 12345 --json
  confluence blogpost view 12345 --body-only > body.html`,
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
			post, err := conf.GetBlogPost(ctx, c, args[0])
			if err != nil {
				return err
			}
			if bodyOnly {
				w.Text("%s", post.Body.Storage.Value)
				return nil
			}
			if w.IsJSON() {
				return w.JSON(post)
			}
			if markdown || rawStorage {
				w.Text("# %s\n\n", post.Title)
				if u := post.AbsoluteURL(); u != "" {
					w.Text("**URL:** %s\n", u)
				}
				if post.Space.Key != "" {
					w.Text("**Space:** %s\n", post.Space.Key)
				}
				if post.Version.Number > 0 {
					w.Text("**Version:** %d (%s)\n", post.Version.Number, post.Version.When)
				}
				w.Text("\n---\n\n")
				w.Text("%s\n", post.RenderMarkdown(rawStorage))
				return nil
			}
			w.Text("%s\t%s\t%s\n", post.ID, post.Type, post.Title)
			w.Text("  space: %s\n  url:   %s\n", post.Space.Key, post.AbsoluteURL())
			return nil
		},
	}
	cmd.Flags().BoolVar(&markdown, "markdown", false, "Render body as Markdown (default output)")
	cmd.Flags().BoolVar(&rawStorage, "raw-storage", false, "Emit the raw storage-format XML")
	cmd.Flags().BoolVar(&bodyOnly, "body-only", false, "Emit only raw storage-format XML/HTML")
	return cmd
}

func blogpostCreateCmd() *cobra.Command {
	var space, title, bodyFormat, bodyFile, bodyInline string
	var draft bool
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new blog post",
		Long: `Create a new blog post with a storage-format body. On Cloud this uses
the v2 blogposts endpoint; on Server/Data Center this uses the content resource
with type=blogpost.

Examples:
  confluence blogpost create --space ENG --title "Weekly Update" --body-file body.html
  confluence blogpost create --space ENG --title "Draft" --draft --body "<p>Hello</p>"
  echo "<p>Hello</p>" | confluence blogpost create --space ENG --title "Hello" --body-file -`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()

			if space == "" || title == "" {
				return fmt.Errorf("--space and --title are required")
			}
			body, err := resolveBody(bodyFile, bodyInline)
			if err != nil {
				return err
			}
			if body == "" {
				return fmt.Errorf("--body-file or --body is required")
			}
			status := "current"
			if draft {
				status = "draft"
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			created, err := conf.CreateBlogPost(ctx, c, conf.BlogPostInput{
				SpaceKey: space, Title: title, BodyFormat: bodyFormat, BodyValue: body, Status: status,
			})
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(created)
			}
			w.Text("created blogpost %s", created.ID)
			if created.Version.Number > 0 {
				w.Text(" (v%d)", created.Version.Number)
			}
			if u := created.AbsoluteURL(); u != "" {
				w.Text("\n%s", u)
			}
			w.Text("\n")
			return nil
		},
	}
	cmd.Flags().StringVar(&space, "space", "", "Space key (required)")
	cmd.Flags().StringVar(&title, "title", "", "Blog post title (required)")
	cmd.Flags().StringVar(&bodyFormat, "body-format", "storage", "Body format: storage | wiki | view")
	cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to body file, or '-' for stdin")
	cmd.Flags().StringVar(&bodyInline, "body", "", "Inline body string")
	cmd.Flags().BoolVar(&draft, "draft", false, "Create as draft")
	return cmd
}

func blogpostUpdateCmd() *cobra.Command {
	var title, bodyFormat, bodyFile, bodyInline, newVersion string
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an existing blog post",
		Long: `Update a blog post title, body, or both. The command fetches the current
blog post version and body, then writes version.number + 1 unless --version is
supplied.

Examples:
  confluence blogpost update 12345 --title "New Title"
  confluence blogpost update 12345 --body-file body.html
  echo "<p>Hello</p>" | confluence blogpost update 12345 --body-file -`,
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
			if title == "" && body == "" {
				return fmt.Errorf("at least one of --title, --body-file, or --body is required")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			post, err := conf.GetBlogPost(ctx, c, args[0])
			if err != nil {
				return err
			}
			if title == "" {
				title = post.Title
			}
			if body == "" {
				body = post.Body.Storage.Value
			}
			ver := post.Version.Number
			if newVersion != "" {
				if _, err := fmt.Sscanf(newVersion, "%d", &ver); err != nil || ver < 0 {
					return fmt.Errorf("invalid --version %q", newVersion)
				}
			}
			updated, err := conf.UpdateBlogPost(ctx, c, conf.BlogPostInput{
				ID: args[0], Title: title, BodyFormat: bodyFormat, BodyValue: body, VersionNumber: ver,
			})
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(updated)
			}
			w.Text("updated blogpost %s (v%d)\n", updated.ID, updated.Version.Number)
			return nil
		},
	}
	cmd.Flags().StringVar(&title, "title", "", "New title (keeps existing if omitted)")
	cmd.Flags().StringVar(&bodyFormat, "body-format", "storage", "Body format: storage | wiki | view")
	cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to body file, or '-' for stdin")
	cmd.Flags().StringVar(&bodyInline, "body", "", "Inline body string")
	cmd.Flags().StringVar(&newVersion, "version", "", "Explicit current version number (auto-fetched if omitted)")
	return cmd
}

func blogpostDeleteCmd() *cobra.Command {
	var force, draft bool
	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Move a blog post to trash",
		Long: `Move a blog post to trash. On Cloud, --draft deletes a draft blog post;
discarded drafts are permanently deleted by Confluence and are not sent to trash.

Examples:
  confluence blogpost delete 12345
  confluence blogpost delete 12345 --force
  confluence blogpost delete 12345 --draft --force --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()

			if !force && !confirmDelete(args[0], "blogpost") {
				return fmt.Errorf("delete cancelled")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			if err := conf.DeleteBlogPost(ctx, c, args[0], conf.BlogPostDeleteOptions{Draft: draft}); err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(map[string]any{"deleted": true, "id": args[0], "draft": draft})
			}
			if draft {
				w.Text("deleted draft blogpost %s\n", args[0])
				return nil
			}
			w.Text("moved blogpost %s to trash\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "Do not prompt for confirmation")
	cmd.Flags().BoolVar(&draft, "draft", false, "Cloud only: delete a draft blog post")
	return cmd
}

func blogpostPurgeCmd() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "purge <id>",
		Short: "Permanently delete a trashed blog post",
		Long: `Permanently delete a trashed blog post. Cloud requires the blog post to
already be in trash. Server/Data Center sends status=trashed to purge trashable
content.

Examples:
  confluence blogpost purge 12345
  confluence blogpost purge 12345 --force
  confluence blogpost purge 12345 --force --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()

			if !force && !confirmDelete(args[0], "blogpost") {
				return fmt.Errorf("delete cancelled")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			if err := conf.DeleteBlogPost(ctx, c, args[0], conf.BlogPostDeleteOptions{Purge: true}); err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(map[string]any{"purged": true, "id": args[0]})
			}
			w.Text("purged blogpost %s\n", args[0])
			return nil
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "Do not prompt for confirmation")
	return cmd
}
