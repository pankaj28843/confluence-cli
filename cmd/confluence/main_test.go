package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func captureHelp(t *testing.T, path []string) string {
	t.Helper()
	root := &cobra.Command{Use: "confluence"}
	root.PersistentFlags().Bool("json", false, "")
	root.PersistentFlags().String("jq", "", "")
	root.PersistentFlags().String("template", "", "")
	root.PersistentFlags().Bool("timing", false, "")
	root.PersistentFlags().Bool("debug", false, "")
	root.AddCommand(doctorCmd(), spaceCmd(), pageCmd(), attachmentCmd(), labelCmd(), commentCmd(), userCmd(), groupCmd(), watcherCmd(), restrictionCmd(), searchCmd(), apiCmd())

	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs(append(path, "--help"))
	_ = root.Execute()
	return buf.String()
}

// TestExamplesHelpBlocks asserts every leaf verb's help text contains an
// "Examples:" block.
func TestExamplesHelpBlocks(t *testing.T) {
	cases := [][]string{
		{"doctor"},
		{"space", "list"},
		{"space", "view"},
		{"page", "view"},
		{"page", "search"},
		{"page", "children"},
		{"page", "ancestors"},
		{"page", "history"},
		{"page", "versions"},
		{"page", "create"},
		{"page", "update"},
		{"page", "publish"},
		{"page", "url"},
		{"page", "screenshot"},
		{"attachment", "list"},
		{"attachment", "download"},
		{"attachment", "upload"},
		{"attachment", "replace"},
		{"attachment", "delete"},
		{"label", "list"},
		{"label", "add"},
		{"label", "remove"},
		{"comment", "list"},
		{"user", "current"},
		{"user", "view"},
		{"user", "search"},
		{"group", "list"},
		{"group", "members"},
		{"watcher", "content"},
		{"watcher", "space"},
		{"restriction", "list"},
		{"search", "content"},
		{"search", "spaces"},
		{"search", "users"},
		{"search", "attachments"},
		{"search", "all"},
		{"api"},
	}
	for _, c := range cases {
		help := captureHelp(t, c)
		if !strings.Contains(help, "Examples:") {
			t.Errorf("%s --help missing 'Examples:' block", strings.Join(c, " "))
		}
	}
}

// TestJSONFlagReachesAllCommands asserts --json is surfaced on list/search/view.
func TestJSONFlagReachesAllCommands(t *testing.T) {
	cases := [][]string{
		{"space", "list"},
		{"page", "search"},
		{"attachment", "list"},
		{"label", "list"},
		{"comment", "list"},
		{"user", "search"},
		{"search", "all"},
		{"api"},
	}
	for _, c := range cases {
		help := captureHelp(t, c)
		if !strings.Contains(help, "--json") {
			t.Errorf("%s --help does not surface --json", strings.Join(c, " "))
		}
	}
}
