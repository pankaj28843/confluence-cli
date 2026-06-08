package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

func docsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "docs",
		Short: "Generate CLI reference documentation",
		Long: `Generate CLI reference documentation from Cobra command metadata.

Examples:
  confluence docs markdown --out docs/cli`,
	}
	cmd.AddCommand(docsMarkdownCmd())
	return cmd
}

func docsMarkdownCmd() *cobra.Command {
	var out string
	cmd := &cobra.Command{
		Use:   "markdown",
		Short: "Generate Markdown CLI reference files",
		Long: `Generate one timestamp-free Markdown file per command. The output is stable
and suitable for docs sites, search indexes, and LLM context.

Examples:
  confluence docs markdown --out docs/cli
  confluence docs markdown --out /tmp/confluence-cli-docs`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			w := getWriter()
			defer w.Finish()
			if out == "" {
				return fmt.Errorf("--out is required")
			}
			root := newRootCommand()
			if err := generateMarkdownDocs(root, out); err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(map[string]any{"ok": true, "out": out})
			}
			w.Text("wrote CLI docs to %s\n", out)
			return nil
		},
	}
	cmd.Flags().StringVar(&out, "out", "", "Output directory (required)")
	return cmd
}

func generateMarkdownDocs(root *cobra.Command, out string) error {
	if root == nil {
		return fmt.Errorf("root command is required")
	}
	if out == "" {
		return fmt.Errorf("output directory is required")
	}
	if err := os.MkdirAll(out, 0o755); err != nil {
		return fmt.Errorf("create docs dir %s: %w", out, err)
	}
	root.DisableAutoGenTag = true

	commands := collectDocCommands(root)
	for _, cmd := range commands {
		if err := writeMarkdownDoc(cmd, out); err != nil {
			return err
		}
	}
	return nil
}

func collectDocCommands(root *cobra.Command) []*cobra.Command {
	var out []*cobra.Command
	var walk func(*cobra.Command)
	walk = func(cmd *cobra.Command) {
		if cmd == nil || cmd.Hidden {
			return
		}
		out = append(out, cmd)
		children := cmd.Commands()
		sort.SliceStable(children, func(i, j int) bool {
			return children[i].Name() < children[j].Name()
		})
		for _, child := range children {
			if child.IsAvailableCommand() || child.HasAvailableSubCommands() {
				walk(child)
			}
		}
	}
	walk(root)
	return out
}

func writeMarkdownDoc(cmd *cobra.Command, out string) error {
	name := commandPath(cmd)
	filename := strings.ReplaceAll(name, " ", "_") + ".md"
	path := filepath.Join(out, filename)

	var b bytes.Buffer
	fmt.Fprintf(&b, "# %s\n\n", name)
	if cmd.Short != "" {
		fmt.Fprintf(&b, "%s\n\n", cmd.Short)
	}
	if cmd.Long != "" {
		writeLongSections(&b, cmd.Long)
	}
	if cmd.Example != "" {
		writeSection(&b, "Examples", cmd.Example)
	}
	fmt.Fprintf(&b, "## Usage\n\n```text\n%s\n```\n\n", cmd.UseLine())
	if cmd.HasAvailableSubCommands() {
		writeSubcommands(&b, cmd)
	}
	writeFlagSection(&b, "Options", cmd.LocalFlags().FlagUsages())
	writeFlagSection(&b, "Inherited Options", cmd.InheritedFlags().FlagUsages())
	if parent := cmd.Parent(); parent != nil {
		fmt.Fprintf(&b, "## See Also\n\n- [%s](%s.md)\n", commandPath(parent), strings.ReplaceAll(commandPath(parent), " ", "_"))
	}

	if err := os.WriteFile(path, b.Bytes(), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

func writeSection(b *bytes.Buffer, title, text string) {
	fmt.Fprintf(b, "## %s\n\n%s\n\n", title, strings.TrimSpace(text))
}

func writeLongSections(b *bytes.Buffer, text string) {
	type section struct {
		title string
		lines []string
	}

	var sections []section
	current := section{title: "Synopsis"}
	for _, line := range strings.Split(strings.TrimSpace(text), "\n") {
		trimmed := strings.TrimSpace(line)
		if isDocHeading(line) {
			if strings.TrimSpace(strings.Join(current.lines, "\n")) != "" {
				sections = append(sections, current)
			}
			current = section{title: strings.TrimSuffix(trimmed, ":")}
			continue
		}
		current.lines = append(current.lines, line)
	}
	if strings.TrimSpace(strings.Join(current.lines, "\n")) != "" {
		sections = append(sections, current)
	}
	for _, section := range sections {
		writeSection(b, section.title, strings.Join(section.lines, "\n"))
	}
}

func isDocHeading(line string) bool {
	trimmed := strings.TrimSpace(line)
	return line == trimmed && strings.HasSuffix(trimmed, ":") && len(trimmed) > 1
}

func writeSubcommands(b *bytes.Buffer, cmd *cobra.Command) {
	fmt.Fprintln(b, "## Commands")
	fmt.Fprintln(b)
	children := cmd.Commands()
	sort.SliceStable(children, func(i, j int) bool {
		return children[i].Name() < children[j].Name()
	})
	for _, child := range children {
		if !child.IsAvailableCommand() && !child.HasAvailableSubCommands() {
			continue
		}
		fmt.Fprintf(b, "- `%s` - %s\n", child.Name(), child.Short)
	}
	fmt.Fprintln(b)
}

func writeFlagSection(b *bytes.Buffer, title, usages string) {
	usages = strings.TrimRight(usages, "\n")
	if strings.TrimSpace(usages) == "" {
		return
	}
	fmt.Fprintf(b, "## %s\n\n```text\n%s\n```\n\n", title, usages)
}

func commandPath(cmd *cobra.Command) string {
	if cmd == nil {
		return ""
	}
	return cmd.CommandPath()
}
