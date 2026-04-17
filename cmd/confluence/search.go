package main

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"github.com/pankaj28843/confluence-cli/internal/conf"
)

func searchCmdReal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search content, spaces, users, attachments, or all",
		Long: `Search Confluence. Five sub-verbs — content, spaces, users, attachments, all.
'all' fans out to the first four in parallel and merges via reciprocal-rank fusion.

Examples:
  confluence search content "release"
  confluence search spaces "engineering"
  confluence search users "Jane Smith"
  confluence search attachments "logo"
  confluence search all "release process" --json`,
	}
	cmd.AddCommand(searchContentCmd())
	cmd.AddCommand(searchSpacesCmd())
	cmd.AddCommand(searchUsersCmd())
	cmd.AddCommand(searchAttachmentsCmd())
	cmd.AddCommand(searchAllCmd())
	return cmd
}

func searchContentCmd() *cobra.Command {
	var space string
	var limit int
	cmd := &cobra.Command{
		Use:   "content <query>",
		Short: "Search pages (type=page by default) via CQL text match",
		Long: `Search pages via /rest/api/content/search with
cql = 'type=page AND text ~ "<query>"' (plus optional --space).

Examples:
  confluence search content "release" --limit 10
  confluence search content "deploy" --space ENG --json`,
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
			cql := fmt.Sprintf(`type=page AND text ~ "%s"`, args[0])
			if space != "" {
				cql += fmt.Sprintf(` AND space="%s"`, space)
			}
			hits, err := conf.SearchCQL(ctx, c, cql, limit, "")
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
	cmd.Flags().StringVar(&space, "space", "", "Optional space key filter")
	cmd.Flags().IntVar(&limit, "limit", 25, "Max results (hard cap 200)")
	return cmd
}

func searchSpacesCmd() *cobra.Command {
	var limit int
	cmd := &cobra.Command{
		Use:   "spaces <query>",
		Short: "Search spaces via CQL type=space",
		Long: `Search spaces via /rest/api/search with
cql = 'type=space AND (title ~ "Q" OR text ~ "Q")'.

Examples:
  confluence search spaces "engineering" --json`,
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
			cql := fmt.Sprintf(`type=space AND (title ~ "%s" OR text ~ "%s")`, args[0], args[0])
			hits, err := conf.SearchGeneric(ctx, c, cql, limit)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(hits)
			}
			for _, h := range hits {
				name := h.Title
				if h.Space != nil && name == "" {
					name = h.Space.Name
				}
				w.Text("%s\n", name)
			}
			return nil
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 25, "Max results (hard cap 200)")
	return cmd
}

func searchUsersCmd() *cobra.Command {
	var limit int
	cmd := &cobra.Command{
		Use:   "users <query>",
		Short: "Search users via CQL user.fullname~",
		Long: `Search users via /rest/api/search with
cql = 'type=user AND user.fullname ~ "<query>"'.

Examples:
  confluence search users "Jane Smith" --json`,
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

func searchAttachmentsCmd() *cobra.Command {
	var limit int
	var space string
	cmd := &cobra.Command{
		Use:   "attachments <query>",
		Short: "Search attachments via CQL type=attachment",
		Long: `Search attachments via /rest/api/content/search with
cql = 'type=attachment AND (title ~ "Q" OR text ~ "Q")'.

Examples:
  confluence search attachments "report.pdf" --json
  confluence search attachments "logo" --space ENG`,
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
			cql := fmt.Sprintf(`type=attachment AND (title ~ "%s" OR text ~ "%s")`, args[0], args[0])
			if space != "" {
				cql += fmt.Sprintf(` AND space="%s"`, space)
			}
			hits, err := conf.SearchCQL(ctx, c, cql, limit, "version,space")
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
	cmd.Flags().StringVar(&space, "space", "", "Optional space key filter")
	cmd.Flags().IntVar(&limit, "limit", 25, "Max results (hard cap 200)")
	return cmd
}

// UnifiedResult is the stable JSON shape for `search all`.
type UnifiedResult struct {
	Kind    string      `json:"kind"` // page | space | user | attachment
	Title   string      `json:"title"`
	URL     string      `json:"url,omitempty"`
	Snippet string      `json:"snippet,omitempty"`
	Space   string      `json:"space,omitempty"`
	Score   float64     `json:"score"`
	Source  interface{} `json:"source"`
}

func searchAllCmd() *cobra.Command {
	var space string
	var limit int
	var branchTimeout time.Duration
	cmd := &cobra.Command{
		Use:   "all <query>",
		Short: "Unified search — pages + spaces + users + attachments in parallel",
		Long: `Fan-out page + space + user + attachment search in parallel, merge via
reciprocal-rank fusion (k=60). Partial failures log a warning to stderr and
the command continues with the healthy branches.

Output schema (stable):
  [{"kind":"page|space|user|attachment","title":"...","url":"...",
    "snippet":"...","score":0.0123,"source":{...raw...}}, ...]

Examples:
  confluence search all "release process" --limit 20 --json
  confluence search all "deploy" --space ENG --branch-timeout 800ms
  confluence search all "<q>" --jq '.[] | {kind, title, space}'`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			parent, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			c, err := newClient()
			if err != nil {
				return err
			}

			branchCtx := func() (context.Context, context.CancelFunc) {
				if branchTimeout <= 0 {
					return parent, func() {}
				}
				return context.WithTimeout(parent, branchTimeout)
			}

			type branch struct {
				kind string
				fn   func() ([]UnifiedResult, error)
			}
			branches := []branch{
				{"page", func() ([]UnifiedResult, error) {
					ctx, cb := branchCtx()
					defer cb()
					cql := fmt.Sprintf(`type=page AND text ~ "%s"`, args[0])
					if space != "" {
						cql += fmt.Sprintf(` AND space="%s"`, space)
					}
					hits, err := conf.SearchCQL(ctx, c, cql, limit, "")
					if err != nil {
						return nil, err
					}
					out := make([]UnifiedResult, 0, len(hits))
					for _, h := range hits {
						out = append(out, UnifiedResult{
							Kind: "page", Title: h.Title, URL: h.AbsoluteURL(),
							Space: h.Space.Key, Source: h,
						})
					}
					return out, nil
				}},
				{"space", func() ([]UnifiedResult, error) {
					ctx, cb := branchCtx()
					defer cb()
					cql := fmt.Sprintf(`type=space AND (title ~ "%s" OR text ~ "%s")`, args[0], args[0])
					hits, err := conf.SearchGeneric(ctx, c, cql, limit)
					if err != nil {
						return nil, err
					}
					out := make([]UnifiedResult, 0, len(hits))
					for _, h := range hits {
						name := h.Title
						spaceKey := ""
						if h.Space != nil {
							if name == "" {
								name = h.Space.Name
							}
							spaceKey = h.Space.Key
						}
						out = append(out, UnifiedResult{
							Kind: "space", Title: name, URL: h.URL,
							Snippet: conf.StripHTML(h.Excerpt),
							Space:   spaceKey, Source: h,
						})
					}
					return out, nil
				}},
				{"user", func() ([]UnifiedResult, error) {
					ctx, cb := branchCtx()
					defer cb()
					hits, err := conf.SearchUsers(ctx, c, args[0], limit)
					if err != nil {
						return nil, err
					}
					out := make([]UnifiedResult, 0, len(hits))
					for _, h := range hits {
						if h.User == nil {
							continue
						}
						u := h.User
						out = append(out, UnifiedResult{
							Kind: "user", Title: firstNonEmpty(u.DisplayName, u.PublicName, u.Username), URL: h.URL,
							Snippet: conf.StripHTML(h.Excerpt), Source: h,
						})
					}
					return out, nil
				}},
				{"attachment", func() ([]UnifiedResult, error) {
					ctx, cb := branchCtx()
					defer cb()
					cql := fmt.Sprintf(`type=attachment AND (title ~ "%s" OR text ~ "%s")`, args[0], args[0])
					if space != "" {
						cql += fmt.Sprintf(` AND space="%s"`, space)
					}
					hits, err := conf.SearchCQL(ctx, c, cql, limit, "version,space")
					if err != nil {
						return nil, err
					}
					out := make([]UnifiedResult, 0, len(hits))
					for _, h := range hits {
						out = append(out, UnifiedResult{
							Kind: "attachment", Title: h.Title, URL: h.AbsoluteURL(),
							Space: h.Space.Key, Source: h,
						})
					}
					return out, nil
				}},
			}

			var wg sync.WaitGroup
			results := make([][]UnifiedResult, len(branches))
			errs := make([]error, len(branches))
			for i, b := range branches {
				wg.Add(1)
				go func(i int, b branch) {
					defer wg.Done()
					r, err := b.fn()
					results[i] = r
					errs[i] = err
				}(i, b)
			}
			wg.Wait()

			for i, err := range errs {
				if err != nil {
					w.Warn("warning: %s search failed: %s\n", branches[i].kind, err)
				}
			}

			merged := mergeRRF(results, 60)
			if w.IsJSON() {
				return w.JSON(merged)
			}
			for _, r := range merged {
				w.Text("[%s] %s\t(%.4f)\n", r.Kind, r.Title, r.Score)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&space, "space", "", "Optional space key filter for page/attachment branches")
	cmd.Flags().IntVar(&limit, "limit", 25, "Max results per branch (hard cap 200)")
	cmd.Flags().DurationVar(&branchTimeout, "branch-timeout", 1200*time.Millisecond,
		"Cancel a branch after this duration (0 disables — parent context only)")
	return cmd
}

// mergeRRF combines per-branch ranked lists using reciprocal-rank fusion.
func mergeRRF(branches [][]UnifiedResult, k int) []UnifiedResult {
	var out []UnifiedResult
	for _, lst := range branches {
		for rank, item := range lst {
			item.Score = 1.0 / float64(k+rank+1)
			out = append(out, item)
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Score > out[j].Score })
	return out
}
