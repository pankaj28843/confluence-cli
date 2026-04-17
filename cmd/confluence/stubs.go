package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// runDoctor is replaced by init() in doctor.go (Phase 2).
var runDoctor = func(cmd *cobra.Command, args []string) error {
	return notImplemented("doctor")
}

func notImplemented(what string) error {
	return fmt.Errorf("%s: not yet implemented", what)
}

// All command-group placeholders. Each phase replaces its group with real
// implementations. Keeping stubs here means cmd/confluence/main.go compiles
// from day one.

func doctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Verify environment, auth, and flavor detection",
		Long: `Verify that:
  1. Required environment variables are set
  2. Flavor is detected (server | cloud)
  3. The auth credentials are valid (GET /rest/api/user/current)

Exits 0 on success, 2 on user-fixable config errors, 1 on unexpected errors.

Examples:
  confluence doctor
  confluence doctor --json`,
		Args: cobra.NoArgs,
		RunE: runDoctor,
	}
}

func spaceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "space",
		Short: "Spaces (list, view)",
		Long: `Space operations.

Examples:
  confluence space list
  confluence space view ENG`,
	}
	cmd.AddCommand(&cobra.Command{Use: "list", Short: "List spaces", Long: "List spaces.\n\nExamples:\n  confluence space list --json", RunE: func(*cobra.Command, []string) error { return notImplemented("space list") }})
	cmd.AddCommand(&cobra.Command{Use: "view <key>", Short: "Show one space", Long: "Show one space by key.\n\nExamples:\n  confluence space view ENG --json", Args: cobra.ExactArgs(1), RunE: func(*cobra.Command, []string) error { return notImplemented("space view") }})
	return cmd
}

func pageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "page",
		Short: "Pages (view, search, children, ancestors, history, update)",
		Long: `Page operations.

Examples:
  confluence page view 12345 --markdown
  confluence page search --cql "type=page AND space=ENG"
  confluence page children 12345 --recursive`,
	}
	for _, verb := range []string{"view <id>", "search", "children <id>", "ancestors <id>", "history <id>", "update <id>"} {
		v := verb
		cmd.AddCommand(&cobra.Command{
			Use:   v,
			Short: "Stub (Phase 3/7)",
			Long:  "Stub.\n\nExamples:\n  confluence page " + v,
			RunE:  func(*cobra.Command, []string) error { return notImplemented("page " + v) },
		})
	}
	return cmd
}

func attachmentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attachment",
		Short: "Attachments (list, download, upload)",
		Long: `Attachment operations.

Examples:
  confluence attachment list --page 12345
  confluence attachment download --page 12345 --name logo.png`,
	}
	for _, verb := range []string{"list", "download", "upload"} {
		v := verb
		cmd.AddCommand(&cobra.Command{Use: v, Short: "Stub (Phase 4/7)", Long: "Stub.\n\nExamples:\n  confluence attachment " + v, RunE: func(*cobra.Command, []string) error { return notImplemented("attachment " + v) }})
	}
	return cmd
}

func labelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "label",
		Short: "Content labels (list, add, remove)",
		Long: `Label operations.

Examples:
  confluence label list --page 12345
  confluence label add --page 12345 --label needs-review`,
	}
	for _, verb := range []string{"list", "add", "remove"} {
		v := verb
		cmd.AddCommand(&cobra.Command{Use: v, Short: "Stub (Phase 4/7)", Long: "Stub.\n\nExamples:\n  confluence label " + v, RunE: func(*cobra.Command, []string) error { return notImplemented("label " + v) }})
	}
	return cmd
}

func commentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "comment",
		Short: "Comments (list inline/footer/resolved)",
		Long: `Comment operations.

Examples:
  confluence comment list --page 12345
  confluence comment list --page 12345 --locations footer,inline`,
	}
	cmd.AddCommand(&cobra.Command{Use: "list", Short: "Stub (Phase 4)", Long: "Stub.\n\nExamples:\n  confluence comment list --page 12345", RunE: func(*cobra.Command, []string) error { return notImplemented("comment list") }})
	return cmd
}

func userCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Users (current, view, search)",
		Long: `User operations.

Examples:
  confluence user current
  confluence user view alice`,
	}
	for _, verb := range []string{"current", "view <keyOrName>", "search"} {
		v := verb
		cmd.AddCommand(&cobra.Command{Use: v, Short: "Stub (Phase 5)", Long: "Stub.\n\nExamples:\n  confluence user " + v, RunE: func(*cobra.Command, []string) error { return notImplemented("user " + v) }})
	}
	return cmd
}

func groupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "group",
		Short: "Groups (list, members)",
		Long: `Group operations.

Examples:
  confluence group list
  confluence group members engineering`,
	}
	for _, verb := range []string{"list", "members <name>"} {
		v := verb
		cmd.AddCommand(&cobra.Command{Use: v, Short: "Stub (Phase 5)", Long: "Stub.\n\nExamples:\n  confluence group " + v, RunE: func(*cobra.Command, []string) error { return notImplemented("group " + v) }})
	}
	return cmd
}

func watcherCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watcher",
		Short: "Watchers (content, space)",
		Long: `Watcher operations.

Examples:
  confluence watcher content --page 12345
  confluence watcher space --space ENG`,
	}
	for _, verb := range []string{"content", "space"} {
		v := verb
		cmd.AddCommand(&cobra.Command{Use: v, Short: "Stub (Phase 5)", Long: "Stub.\n\nExamples:\n  confluence watcher " + v, RunE: func(*cobra.Command, []string) error { return notImplemented("watcher " + v) }})
	}
	return cmd
}

func restrictionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restriction",
		Short: "Content restrictions (list)",
		Long: `Restriction operations.

Examples:
  confluence restriction list --page 12345`,
	}
	cmd.AddCommand(&cobra.Command{Use: "list", Short: "Stub (Phase 5)", Long: "Stub.\n\nExamples:\n  confluence restriction list --page 12345", RunE: func(*cobra.Command, []string) error { return notImplemented("restriction list") }})
	return cmd
}

func searchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search content, spaces, users, attachments, or all",
		Long: `Search Confluence. Five sub-verbs — content, spaces, users, attachments, all.
'all' fans out to the first four in parallel and merges via reciprocal-rank fusion.

Examples:
  confluence search content "release"
  confluence search spaces "engineering"
  confluence search all "release process" --json`,
	}
	for _, verb := range []string{"content <query>", "spaces <query>", "users <query>", "attachments <query>", "all <query>"} {
		v := verb
		cmd.AddCommand(&cobra.Command{Use: v, Short: "Stub (Phase 6)", Long: "Stub.\n\nExamples:\n  confluence search " + v, Args: cobra.ExactArgs(1), RunE: func(*cobra.Command, []string) error { return notImplemented("search " + v) }})
	}
	return cmd
}

func apiCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "api <path>",
		Short: "Call any Confluence REST endpoint (escape hatch)",
		Long: `Issue a raw REST call. Auto-routes to v1 or v2 on Cloud based on the path.

Examples:
  confluence api /rest/api/user/current
  confluence api /wiki/api/v2/spaces --param 'limit=10'`,
		Args: cobra.ExactArgs(1),
		RunE: func(*cobra.Command, []string) error { return notImplemented("api") },
	}
}
