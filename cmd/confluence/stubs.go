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

// Command-group factories. Each phase replaces its group with the full
// implementation via a *Real() function and edits this file.
//
// Live groups (edit entries here to point at the real cmd):
func spaceCmd() *cobra.Command       { return spaceCmdReal() }
func pageCmd() *cobra.Command        { return pageCmdReal() }
func attachmentCmd() *cobra.Command  { return attachmentCmdReal() }
func labelCmd() *cobra.Command       { return labelCmdReal() }
func commentCmd() *cobra.Command     { return commentCmdReal() }
func userCmd() *cobra.Command        { return userCmdReal() }
func groupCmd() *cobra.Command       { return groupCmdReal() }
func watcherCmd() *cobra.Command     { return watcherCmdReal() }
func restrictionCmd() *cobra.Command { return restrictionCmdReal() }

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

// Still stubbed:

func searchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search content, spaces, users, attachments, or all",
		Long: `Search Confluence.

Examples:
  confluence search content "release"
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
		Long: `Issue a raw REST call.

Examples:
  confluence api /rest/api/user/current`,
		Args: cobra.ExactArgs(1),
		RunE: func(*cobra.Command, []string) error { return notImplemented("api") },
	}
}
