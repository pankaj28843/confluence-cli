package main

import (
	"github.com/spf13/cobra"

	"github.com/pankaj28843/confluence-cli/internal/conf"
)

func themeCmdReal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "theme",
		Short: "Cloud theme reads",
		Long: `Cloud theme read operations.

Theme reads are documented in Confluence Cloud REST API v1. Server/Data Center
does not expose the same theme REST group in the current official REST
reference, so these typed commands are Cloud-only.

Examples:
  confluence theme list --limit 10
  confluence theme global --json
  confluence theme space ENG`,
	}
	cmd.AddCommand(themeListCmd())
	cmd.AddCommand(themeGlobalCmd())
	cmd.AddCommand(themeViewCmd())
	cmd.AddCommand(themeSpaceCmd())
	return cmd
}

func themeListCmd() *cobra.Command {
	var start int
	var limit int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Cloud themes",
		Long: `List available Confluence Cloud themes.

Examples:
  confluence theme list
  confluence theme list --start 25 --limit 25 --json`,
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
			themes, err := conf.ListThemes(ctx, c, conf.ThemeListOptions{Start: start, Limit: limit})
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(themes)
			}
			printThemes(w, themes)
			return nil
		},
	}
	cmd.Flags().IntVar(&start, "start", 0, "Zero-based result offset")
	cmd.Flags().IntVar(&limit, "limit", 25, "Max themes (hard cap 200)")
	return cmd
}

func themeGlobalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "global",
		Short: "Show selected Cloud global theme",
		Long: `Show the selected global theme for the Confluence Cloud site.

Examples:
  confluence theme global
  confluence theme global --json`,
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
			theme, err := conf.GetGlobalTheme(ctx, c)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(theme)
			}
			printTheme(w, *theme)
			return nil
		},
	}
	return cmd
}

func themeViewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view <theme-key>",
		Short: "Show one Cloud theme",
		Long: `Show one Confluence Cloud theme by key.

Examples:
  confluence theme view global
  confluence theme view global --json`,
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
			theme, err := conf.GetTheme(ctx, c, args[0])
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(theme)
			}
			printTheme(w, *theme)
			return nil
		},
	}
	return cmd
}

func themeSpaceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "space <space-key>",
		Short: "Show selected Cloud space theme",
		Long: `Show the selected theme for one Confluence Cloud space.

Examples:
  confluence theme space ENG
  confluence theme space ENG --json`,
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
			theme, err := conf.GetSpaceTheme(ctx, c, args[0])
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(theme)
			}
			printTheme(w, *theme)
			return nil
		},
	}
	return cmd
}

type themeTextWriter interface {
	Text(format string, args ...any)
}

func printThemes(w themeTextWriter, themes []conf.Theme) {
	for _, theme := range themes {
		printTheme(w, theme)
	}
}

func printTheme(w themeTextWriter, theme conf.Theme) {
	w.Text("%s\t%s\t%s\n",
		firstNonEmpty(theme.ThemeKey, "-"),
		firstNonEmpty(theme.Name, "-"),
		firstNonEmpty(theme.Description, "-"))
}
