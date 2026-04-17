package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pankaj28843/confluence-cli/internal/conf"
)

func pageCmdReal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "page",
		Short: "Pages (view, search, children, ancestors, history, versions, update)",
		Long: `Page operations.

Examples:
  confluence page view 12345 --markdown
  confluence page search --cql "type=page AND space=ENG"
  confluence page children 12345 --recursive`,
	}
	cmd.AddCommand(pageViewCmd())
	cmd.AddCommand(pageSearchCmd())
	cmd.AddCommand(pageChildrenCmd())
	cmd.AddCommand(pageAncestorsCmd())
	cmd.AddCommand(pageHistoryCmd())
	cmd.AddCommand(pageVersionsCmd())
	cmd.AddCommand(pageUpdateCmd()) // implemented in writes.go (Phase 7)
	return cmd
}

func pageViewCmd() *cobra.Command {
	var markdown, rawStorage bool
	var expand string
	cmd := &cobra.Command{
		Use:   "view <id>",
		Short: "Fetch a page (or other content) by id",
		Long: `Fetch content by id. --markdown renders the storage body to Markdown;
--raw-storage emits the raw Confluence storage XML/HTML.

Examples:
  confluence page view 12345 --markdown
  confluence page view 12345 --json
  confluence page view 12345 --raw-storage`,
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
			p, err := conf.GetContent(ctx, c, args[0], expand)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(p)
			}
			if markdown || rawStorage {
				w.Text("# %s\n\n", p.Title)
				w.Text("**URL:** %s\n", p.AbsoluteURL())
				if p.Space.Key != "" {
					w.Text("**Space:** %s\n", p.Space.Key)
				}
				if p.Version.Number > 0 {
					w.Text("**Version:** %d (%s)\n", p.Version.Number, p.Version.When)
				}
				if len(p.Metadata.Labels.Results) > 0 {
					w.Text("**Labels:** ")
					for i, l := range p.Metadata.Labels.Results {
						if i > 0 {
							w.Text(", ")
						}
						w.Text("%s", l.Name)
					}
					w.Text("\n")
				}
				w.Text("\n---\n\n")
				w.Text("%s\n", p.RenderMarkdown(rawStorage))
				return nil
			}
			w.Text("%s\t%s\t%s\n", p.ID, p.Type, p.Title)
			w.Text("  space: %s\n  url:   %s\n", p.Space.Key, p.AbsoluteURL())
			return nil
		},
	}
	cmd.Flags().BoolVar(&markdown, "markdown", false, "Render body as Markdown (default output)")
	cmd.Flags().BoolVar(&rawStorage, "raw-storage", false, "Emit the raw storage-format XML")
	cmd.Flags().StringVar(&expand, "expand", "", "Override the expand= parameter (default: "+conf.DefaultExpand+")")
	return cmd
}

func pageSearchCmd() *cobra.Command {
	var cql, space, title string
	var limit int
	cmd := &cobra.Command{
		Use:   "search",
		Short: "CQL-powered page search",
		Long: `Search pages using CQL. Convenience flags build CQL for you; --cql overrides
everything. Default CQL if no flag: 'type=page'.

Examples:
  confluence page search --cql "type=page AND space=ENG AND text ~ 'deploy'"
  confluence page search --space ENG --title "Release notes"
  confluence page search --cql "type=page" --limit 10 --json`,
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
			effective := cql
			if effective == "" {
				effective = buildCQL(space, title)
			}
			hits, err := conf.SearchCQL(ctx, c, effective, limit, "")
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(hits)
			}
			for _, h := range hits {
				w.Text("%s\t%s\t%s\n", h.ID, h.Space.Key, h.Title)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&cql, "cql", "", "CQL expression (overrides --space/--title)")
	cmd.Flags().StringVar(&space, "space", "", "Space key filter (CQL helper)")
	cmd.Flags().StringVar(&title, "title", "", "Title contains (CQL helper)")
	cmd.Flags().IntVar(&limit, "limit", 25, "Max results (hard cap 200)")
	return cmd
}

func pageChildrenCmd() *cobra.Command {
	var childType string
	var limit int
	var recursive bool
	cmd := &cobra.Command{
		Use:   "children <id>",
		Short: "List children of a page",
		Long: `List child content under a parent id. Default childType is 'page'.
--recursive traverses breadth-first.

Examples:
  confluence page children 12345
  confluence page children 12345 --recursive --json
  confluence page children 12345 --type comment`,
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
			var collected []conf.Content
			queue := []string{args[0]}
			visited := map[string]bool{}
			for len(queue) > 0 {
				id := queue[0]
				queue = queue[1:]
				if visited[id] {
					continue
				}
				visited[id] = true
				kids, err := conf.GetChildren(ctx, c, id, childType, limit)
				if err != nil {
					return err
				}
				collected = append(collected, kids...)
				if recursive && childType == "page" {
					for _, k := range kids {
						queue = append(queue, k.ID)
					}
				}
			}
			if w.IsJSON() {
				return w.JSON(collected)
			}
			for _, k := range collected {
				w.Text("%s\t%s\t%s\n", k.ID, k.Type, k.Title)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&childType, "type", "page", "Child type: page | comment | attachment")
	cmd.Flags().IntVar(&limit, "limit", 50, "Max children per parent (hard cap 200)")
	cmd.Flags().BoolVar(&recursive, "recursive", false, "Walk the full page tree")
	return cmd
}

func pageAncestorsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ancestors <id>",
		Short: "List ancestor pages of a content id",
		Long: `Walk parents up to the space root.

Examples:
  confluence page ancestors 12345
  confluence page ancestors 12345 --json`,
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
			ancs, err := conf.GetAncestors(ctx, c, args[0])
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(ancs)
			}
			for _, a := range ancs {
				w.Text("%s\t%s\t%s\n", a.ID, a.Type, a.Title)
			}
			return nil
		},
	}
}

func pageHistoryCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "history <id>",
		Short: "Fetch page history",
		Long: `Fetch /rest/api/content/{id}/history.

Examples:
  confluence page history 12345 --json`,
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
			data, err := conf.GetHistory(ctx, c, args[0])
			if err != nil {
				return err
			}
			_, _ = w.Out.Write(data)
			return nil
		},
	}
}

func pageVersionsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "versions <id>",
		Short: "List version records for a page",
		Long: `List /rest/api/content/{id}/version.

Examples:
  confluence page versions 12345 --json`,
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
			data, err := conf.ListVersions(ctx, c, args[0])
			if err != nil {
				return err
			}
			_, _ = w.Out.Write(data)
			return nil
		},
	}
}

func buildCQL(space, title string) string {
	out := "type=page"
	if space != "" {
		out += fmt.Sprintf(` AND space="%s"`, space)
	}
	if title != "" {
		out += fmt.Sprintf(` AND title~"%s"`, title)
	}
	return out
}

// pageUpdateCmd is defined in writes.go (Phase 7). Provide a stub here so the
// page group compiles during Phase 3.
func pageUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update <id>",
		Short: "Update a page (Phase 7)",
		Long: `Update an existing page's title and/or body.

Examples:
  confluence page update 12345 --title "New Title" --body-format wiki --body-file body.txt`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error { return notImplemented("page update") },
	}
}
