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
		{"page", "version"},
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
		{"blogpost", "versions"},
		{"blogpost", "version"},
		{"attachment", "list"},
		{"attachment", "download"},
		{"attachment", "upload"},
		{"attachment", "replace"},
		{"attachment", "delete"},
		{"attachment", "versions"},
		{"attachment", "version"},
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
		{"audit", "list"},
		{"audit", "since"},
		{"audit", "retention"},
		{"database", "view"},
		{"database", "children"},
		{"folder", "view"},
		{"folder", "children"},
		{"whiteboard", "view"},
		{"whiteboard", "children"},
		{"smart-link", "view"},
		{"smart-link", "children"},
		{"custom-content", "list"},
		{"custom-content", "page"},
		{"custom-content", "blogpost"},
		{"custom-content", "space"},
		{"custom-content", "view"},
		{"custom-content", "children"},
		{"custom-content", "versions"},
		{"custom-content", "version"},
		{"macro", "body"},
		{"template", "list"},
		{"template", "blueprint", "list"},
		{"template", "view"},
		{"docs", "markdown"},
		{"comment", "list"},
		{"comment", "add"},
		{"comment", "update"},
		{"comment", "delete"},
		{"comment", "versions"},
		{"comment", "version"},
		{"user", "current"},
		{"user", "view"},
		{"user", "search"},
		{"user", "bulk"},
		{"group", "list"},
		{"group", "view"},
		{"group", "picker"},
		{"group", "members"},
		{"group", "children"},
		{"group", "parents"},
		{"group", "ancestors"},
		{"watcher", "content"},
		{"watcher", "space"},
		{"watcher", "status"},
		{"permission", "space", "list"},
		{"permission", "space", "available"},
		{"permission", "space", "subject"},
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
		{"page", "versions"},
		{"page", "version"},
		{"blogpost", "list"},
		{"blogpost", "versions"},
		{"blogpost", "version"},
		{"attachment", "list"},
		{"attachment", "versions"},
		{"attachment", "version"},
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
		{"audit", "list"},
		{"audit", "since"},
		{"audit", "retention"},
		{"database", "view"},
		{"database", "children"},
		{"folder", "view"},
		{"folder", "children"},
		{"whiteboard", "view"},
		{"whiteboard", "children"},
		{"smart-link", "view"},
		{"smart-link", "children"},
		{"custom-content", "list"},
		{"custom-content", "page"},
		{"custom-content", "blogpost"},
		{"custom-content", "space"},
		{"custom-content", "view"},
		{"custom-content", "children"},
		{"custom-content", "versions"},
		{"custom-content", "version"},
		{"macro", "body"},
		{"template", "list"},
		{"template", "blueprint", "list"},
		{"template", "view"},
		{"comment", "list"},
		{"comment", "add"},
		{"comment", "update"},
		{"comment", "delete"},
		{"comment", "versions"},
		{"comment", "version"},
		{"user", "bulk"},
		{"user", "search"},
		{"group", "list"},
		{"group", "view"},
		{"group", "picker"},
		{"group", "members"},
		{"group", "children"},
		{"group", "parents"},
		{"group", "ancestors"},
		{"watcher", "content"},
		{"watcher", "space"},
		{"watcher", "status"},
		{"permission", "space", "list"},
		{"permission", "space", "available"},
		{"permission", "space", "subject"},
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

func TestUserAndGroupCommandsHaveTypedReadFlags(t *testing.T) {
	help := captureHelp(t, []string{"user", "view"})
	for _, want := range []string{"--username", "--key", "--account-id"} {
		if !strings.Contains(help, want) {
			t.Fatalf("user view --help missing %s:\n%s", want, help)
		}
	}

	help = captureHelp(t, []string{"user", "bulk"})
	for _, want := range []string{"--account-id"} {
		if !strings.Contains(help, want) {
			t.Fatalf("user bulk --help missing %s:\n%s", want, help)
		}
	}

	help = captureHelp(t, []string{"group", "list"})
	for _, want := range []string{"--limit", "--access-type"} {
		if !strings.Contains(help, want) {
			t.Fatalf("group list --help missing %s:\n%s", want, help)
		}
	}

	help = captureHelp(t, []string{"group", "view"})
	for _, want := range []string{"--id", "--expand"} {
		if !strings.Contains(help, want) {
			t.Fatalf("group view --help missing %s:\n%s", want, help)
		}
	}

	help = captureHelp(t, []string{"group", "members"})
	for _, want := range []string{"--id", "--limit", "--expand", "--total-size"} {
		if !strings.Contains(help, want) {
			t.Fatalf("group members --help missing %s:\n%s", want, help)
		}
	}

	help = captureHelp(t, []string{"group", "picker"})
	for _, want := range []string{"--limit", "--total-size"} {
		if !strings.Contains(help, want) {
			t.Fatalf("group picker --help missing %s:\n%s", want, help)
		}
	}

	for _, path := range [][]string{
		{"group", "children"},
		{"group", "parents"},
		{"group", "ancestors"},
	} {
		help = captureHelp(t, path)
		for _, want := range []string{"--limit", "--expand"} {
			if !strings.Contains(help, want) {
				t.Fatalf("%s --help missing %s:\n%s", strings.Join(path, " "), want, help)
			}
		}
	}
}

func TestAuditCommandsHaveTypedReadFlags(t *testing.T) {
	help := captureHelp(t, []string{"audit", "list"})
	for _, want := range []string{"--limit", "--start-date", "--end-date", "--search"} {
		if !strings.Contains(help, want) {
			t.Fatalf("audit list --help missing %s:\n%s", want, help)
		}
	}

	help = captureHelp(t, []string{"audit", "since"})
	for _, want := range []string{"--number", "--unit", "--limit", "--search"} {
		if !strings.Contains(help, want) {
			t.Fatalf("audit since --help missing %s:\n%s", want, help)
		}
	}

	help = captureHelp(t, []string{"audit", "retention"})
	if !strings.Contains(help, "Cloud") {
		t.Fatalf("audit retention --help missing Cloud scope:\n%s", help)
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

func TestPermissionSpaceCommandsHaveTypedReadFlags(t *testing.T) {
	help := captureHelp(t, []string{"permission", "space", "list"})
	for _, want := range []string{"--space", "--limit"} {
		if !strings.Contains(help, want) {
			t.Fatalf("permission space list --help missing %s:\n%s", want, help)
		}
	}

	help = captureHelp(t, []string{"permission", "space", "available"})
	if !strings.Contains(help, "--limit") {
		t.Fatalf("permission space available --help missing --limit:\n%s", help)
	}

	help = captureHelp(t, []string{"permission", "space", "subject"})
	for _, want := range []string{"--space", "--anonymous", "--group", "--user-key", "--limit"} {
		if !strings.Contains(help, want) {
			t.Fatalf("permission space subject --help missing %s:\n%s", want, help)
		}
	}
}

func TestContentTreeCommandsHaveTypedReadFlags(t *testing.T) {
	for _, group := range []string{"database", "folder", "whiteboard", "smart-link"} {
		help := captureHelp(t, []string{group, "view"})
		for _, want := range []string{"--include-collaborators", "--include-direct-children", "--include-operations", "--include-properties"} {
			if !strings.Contains(help, want) {
				t.Fatalf("%s view --help missing %s:\n%s", group, want, help)
			}
		}

		help = captureHelp(t, []string{group, "children"})
		for _, want := range []string{"--limit", "--sort", "--type"} {
			if !strings.Contains(help, want) {
				t.Fatalf("%s children --help missing %s:\n%s", group, want, help)
			}
		}
	}
}

func TestVersionCommandsHaveTypedReadFlags(t *testing.T) {
	for _, group := range []string{"page", "blogpost", "custom-content"} {
		help := captureHelp(t, []string{group, "versions"})
		for _, want := range []string{"--limit", "--sort", "--body-format"} {
			if !strings.Contains(help, want) {
				t.Fatalf("%s versions --help missing %s:\n%s", group, want, help)
			}
		}
		if help := captureHelp(t, []string{group, "version"}); !strings.Contains(help, "<number>") {
			t.Fatalf("%s version --help missing version number usage:\n%s", group, help)
		}
	}
	help := captureHelp(t, []string{"attachment", "versions"})
	for _, want := range []string{"--limit", "--sort"} {
		if !strings.Contains(help, want) {
			t.Fatalf("attachment versions --help missing %s:\n%s", want, help)
		}
	}
	if strings.Contains(help, "--body-format") {
		t.Fatalf("attachment versions --help should not expose undocumented --body-format:\n%s", help)
	}
	if help := captureHelp(t, []string{"attachment", "version"}); !strings.Contains(help, "<number>") {
		t.Fatalf("attachment version --help missing version number usage:\n%s", help)
	}

	help = captureHelp(t, []string{"comment", "versions"})
	for _, want := range []string{"--limit", "--sort", "--body-format", "--location"} {
		if !strings.Contains(help, want) {
			t.Fatalf("comment versions --help missing %s:\n%s", want, help)
		}
	}
	help = captureHelp(t, []string{"comment", "version"})
	for _, want := range []string{"<number>", "--location"} {
		if !strings.Contains(help, want) {
			t.Fatalf("comment version --help missing %s:\n%s", want, help)
		}
	}
}

func TestCustomContentCommandsHaveTypedReadFlags(t *testing.T) {
	help := captureHelp(t, []string{"custom-content", "list"})
	for _, want := range []string{"--type", "--id", "--space-id", "--limit", "--sort", "--body-format"} {
		if !strings.Contains(help, want) {
			t.Fatalf("custom-content list --help missing %s:\n%s", want, help)
		}
	}

	for _, leaf := range []string{"page", "blogpost"} {
		help = captureHelp(t, []string{"custom-content", leaf})
		for _, want := range []string{"--type", "--limit", "--sort", "--body-format"} {
			if !strings.Contains(help, want) {
				t.Fatalf("custom-content %s --help missing %s:\n%s", leaf, want, help)
			}
		}
	}

	help = captureHelp(t, []string{"custom-content", "space"})
	for _, want := range []string{"--type", "--limit", "--body-format"} {
		if !strings.Contains(help, want) {
			t.Fatalf("custom-content space --help missing %s:\n%s", want, help)
		}
	}
	if strings.Contains(help, "--sort") {
		t.Fatalf("custom-content space --help should not expose undocumented --sort:\n%s", help)
	}

	help = captureHelp(t, []string{"custom-content", "view"})
	for _, want := range []string{"--body-format", "--version", "--include-labels", "--include-properties", "--include-operations", "--include-versions", "--include-version", "--include-collaborators"} {
		if !strings.Contains(help, want) {
			t.Fatalf("custom-content view --help missing %s:\n%s", want, help)
		}
	}

	help = captureHelp(t, []string{"custom-content", "children"})
	for _, want := range []string{"--limit", "--sort", "--type"} {
		if !strings.Contains(help, want) {
			t.Fatalf("custom-content children --help missing %s:\n%s", want, help)
		}
	}
}
