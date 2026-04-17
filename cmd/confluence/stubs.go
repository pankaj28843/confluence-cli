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
func searchCmd() *cobra.Command      { return searchCmdReal() }
func apiCmd() *cobra.Command         { return apiCmdReal() }

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
