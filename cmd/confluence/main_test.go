package main

import (
	"bytes"
	"strings"
	"testing"
)

func captureHelp(t *testing.T, path []string) string {
	t.Helper()
	root := newRootCommand()

	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs(append(path, "--help"))
	if err := root.Execute(); err != nil {
		t.Fatalf("%s --help: %v\n%s", strings.Join(path, " "), err, buf.String())
	}
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
		{"page", "direct-children"},
		{"page", "descendants"},
		{"page", "ancestors"},
		{"page", "history"},
		{"page", "versions"},
		{"page", "create"},
		{"page", "update"},
		{"page", "publish"},
		{"page", "delete"},
		{"page", "purge"},
		{"page", "url"},
		{"page", "screenshot"},
		{"blogpost", "list"},
		{"blogpost", "view"},
		{"blogpost", "create"},
		{"blogpost", "update"},
		{"blogpost", "delete"},
		{"blogpost", "purge"},
		{"attachment", "list"},
		{"attachment", "download"},
		{"attachment", "upload"},
		{"attachment", "replace"},
		{"attachment", "delete"},
		{"label", "list"},
		{"label", "space"},
		{"label", "search"},
		{"label", "recent"},
		{"label", "related"},
		{"label", "add"},
		{"label", "remove"},
		{"property", "content", "list"},
		{"property", "content", "get"},
		{"property", "content", "set"},
		{"property", "content", "delete"},
		{"property", "space", "list"},
		{"property", "space", "get"},
		{"property", "space", "set"},
		{"property", "space", "delete"},
		{"task", "list"},
		{"task", "view"},
		{"task", "complete"},
		{"task", "reopen"},
		{"task", "long", "list"},
		{"task", "long", "view"},
		{"operation", "list"},
		{"like", "count"},
		{"like", "users"},
		{"body", "convert"},
		{"macro", "body"},
		{"template", "list"},
		{"template", "blueprint", "list"},
		{"template", "view"},
		{"docs", "markdown"},
		{"comment", "list"},
		{"comment", "add"},
		{"comment", "update"},
		{"comment", "delete"},
		{"user", "current"},
		{"user", "view"},
		{"user", "search"},
		{"group", "list"},
		{"group", "members"},
		{"watcher", "content"},
		{"watcher", "space"},
		{"watcher", "status"},
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
		{"page", "direct-children"},
		{"page", "descendants"},
		{"blogpost", "list"},
		{"attachment", "list"},
		{"label", "list"},
		{"label", "space"},
		{"label", "search"},
		{"label", "recent"},
		{"label", "related"},
		{"property", "content", "list"},
		{"property", "space", "list"},
		{"task", "list"},
		{"task", "long", "list"},
		{"operation", "list"},
		{"like", "count"},
		{"like", "users"},
		{"body", "convert"},
		{"macro", "body"},
		{"template", "list"},
		{"template", "blueprint", "list"},
		{"template", "view"},
		{"comment", "list"},
		{"comment", "add"},
		{"comment", "update"},
		{"comment", "delete"},
		{"user", "search"},
		{"watcher", "content"},
		{"watcher", "space"},
		{"watcher", "status"},
		{"restriction", "list"},
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

func TestBodyCommandExists(t *testing.T) {
	root := newRootCommand()
	path := []string{"body", "convert"}
	cmd, remaining, err := root.Find(path)
	if err != nil {
		t.Fatalf("%s: %v", strings.Join(path, " "), err)
	}
	if len(remaining) != 0 {
		t.Fatalf("%s resolved with remaining args %v; command path=%s", strings.Join(path, " "), remaining, cmd.CommandPath())
	}
	if !strings.HasSuffix(cmd.CommandPath(), strings.Join(path, " ")) {
		t.Fatalf("%s resolved to %s", strings.Join(path, " "), cmd.CommandPath())
	}
}

func TestMacroCommandExists(t *testing.T) {
	root := newRootCommand()
	path := []string{"macro", "body"}
	cmd, remaining, err := root.Find(path)
	if err != nil {
		t.Fatalf("%s: %v", strings.Join(path, " "), err)
	}
	if len(remaining) != 0 {
		t.Fatalf("%s resolved with remaining args %v; command path=%s", strings.Join(path, " "), remaining, cmd.CommandPath())
	}
	if !strings.HasSuffix(cmd.CommandPath(), strings.Join(path, " ")) {
		t.Fatalf("%s resolved to %s", strings.Join(path, " "), cmd.CommandPath())
	}
}

func TestTemplateCommandsExist(t *testing.T) {
	root := newRootCommand()
	for _, path := range [][]string{
		{"template", "list"},
		{"template", "blueprint", "list"},
		{"template", "view"},
	} {
		cmd, remaining, err := root.Find(path)
		if err != nil {
			t.Fatalf("%s: %v", strings.Join(path, " "), err)
		}
		if len(remaining) != 0 {
			t.Fatalf("%s resolved with remaining args %v; command path=%s", strings.Join(path, " "), remaining, cmd.CommandPath())
		}
		if !strings.HasSuffix(cmd.CommandPath(), strings.Join(path, " ")) {
			t.Fatalf("%s resolved to %s", strings.Join(path, " "), cmd.CommandPath())
		}
	}
}

func TestCommentMutationCommandsExist(t *testing.T) {
	root := newRootCommand()
	for _, path := range [][]string{
		{"comment", "add"},
		{"comment", "update"},
		{"comment", "delete"},
	} {
		cmd, remaining, err := root.Find(path)
		if err != nil {
			t.Fatalf("%s: %v", strings.Join(path, " "), err)
		}
		if len(remaining) != 0 {
			t.Fatalf("%s resolved with remaining args %v; command path=%s", strings.Join(path, " "), remaining, cmd.CommandPath())
		}
		if !strings.HasSuffix(cmd.CommandPath(), strings.Join(path, " ")) {
			t.Fatalf("%s resolved to %s", strings.Join(path, " "), cmd.CommandPath())
		}
	}
}

func TestPageDescendantsCommandExists(t *testing.T) {
	root := newRootCommand()
	path := []string{"page", "descendants"}
	cmd, remaining, err := root.Find(path)
	if err != nil {
		t.Fatalf("%s: %v", strings.Join(path, " "), err)
	}
	if len(remaining) != 0 {
		t.Fatalf("%s resolved with remaining args %v; command path=%s", strings.Join(path, " "), remaining, cmd.CommandPath())
	}
	if !strings.HasSuffix(cmd.CommandPath(), strings.Join(path, " ")) {
		t.Fatalf("%s resolved to %s", strings.Join(path, " "), cmd.CommandPath())
	}
}

func TestPageDirectChildrenCommandExists(t *testing.T) {
	root := newRootCommand()
	path := []string{"page", "direct-children"}
	cmd, remaining, err := root.Find(path)
	if err != nil {
		t.Fatalf("%s: %v", strings.Join(path, " "), err)
	}
	if len(remaining) != 0 {
		t.Fatalf("%s resolved with remaining args %v; command path=%s", strings.Join(path, " "), remaining, cmd.CommandPath())
	}
	if !strings.HasSuffix(cmd.CommandPath(), strings.Join(path, " ")) {
		t.Fatalf("%s resolved to %s", strings.Join(path, " "), cmd.CommandPath())
	}
}

func TestOperationAndLikeCommandsExist(t *testing.T) {
	root := newRootCommand()
	for _, path := range [][]string{
		{"operation", "list"},
		{"like", "count"},
		{"like", "users"},
	} {
		cmd, remaining, err := root.Find(path)
		if err != nil {
			t.Fatalf("%s: %v", strings.Join(path, " "), err)
		}
		if len(remaining) != 0 {
			t.Fatalf("%s resolved with remaining args %v; command path=%s", strings.Join(path, " "), remaining, cmd.CommandPath())
		}
		if !strings.HasSuffix(cmd.CommandPath(), strings.Join(path, " ")) {
			t.Fatalf("%s resolved to %s", strings.Join(path, " "), cmd.CommandPath())
		}
	}
}

func TestRestrictionListHasTypedReadFlags(t *testing.T) {
	help := captureHelp(t, []string{"restriction", "list"})
	for _, want := range []string{"--operation", "--limit"} {
		if !strings.Contains(help, want) {
			t.Fatalf("restriction list --help missing %s:\n%s", want, help)
		}
	}
}

func TestWatcherCommandsHaveTypedReadFlags(t *testing.T) {
	for _, path := range [][]string{
		{"watcher", "content"},
		{"watcher", "space"},
	} {
		help := captureHelp(t, path)
		if !strings.Contains(help, "--limit") {
			t.Fatalf("%s --help missing --limit:\n%s", strings.Join(path, " "), help)
		}
	}

	help := captureHelp(t, []string{"watcher", "status"})
	for _, want := range []string{"--page", "--space", "--account-id", "--content-type"} {
		if !strings.Contains(help, want) {
			t.Fatalf("watcher status --help missing %s:\n%s", want, help)
		}
	}
}

func TestLabelCommandsHaveTypedReadFlags(t *testing.T) {
	help := captureHelp(t, []string{"label", "list"})
	for _, want := range []string{"--page", "--blogpost", "--attachment", "--custom-content", "--prefix", "--limit"} {
		if !strings.Contains(help, want) {
			t.Fatalf("label list --help missing %s:\n%s", want, help)
		}
	}

	help = captureHelp(t, []string{"label", "space"})
	for _, want := range []string{"--space", "--prefix", "--limit"} {
		if !strings.Contains(help, want) {
			t.Fatalf("label space --help missing %s:\n%s", want, help)
		}
	}

	help = captureHelp(t, []string{"label", "search"})
	for _, want := range []string{"--label-id", "--prefix", "--sort", "--limit"} {
		if !strings.Contains(help, want) {
			t.Fatalf("label search --help missing %s:\n%s", want, help)
		}
	}

	help = captureHelp(t, []string{"label", "related"})
	for _, want := range []string{"--label", "--space", "--limit"} {
		if !strings.Contains(help, want) {
			t.Fatalf("label related --help missing %s:\n%s", want, help)
		}
	}
}
