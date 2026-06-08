package main

import (
	"github.com/spf13/cobra"

	"github.com/pankaj28843/confluence-cli/internal/conf"
)

func settingsCmdReal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "settings",
		Short: "Cloud settings reads",
		Long: `Cloud settings read operations.

Settings and space-settings reads are documented in Confluence Cloud REST API
v1. Server/Data Center does not expose the same settings REST groups in the
current official REST reference, so these typed commands are Cloud-only.

Examples:
  confluence settings system-info --json
  confluence settings lookandfeel --space ENG
  confluence settings space ENG`,
	}
	cmd.AddCommand(settingsSystemInfoCmd())
	cmd.AddCommand(settingsLookAndFeelCmd())
	cmd.AddCommand(settingsSpaceCmd())
	return cmd
}

func settingsSystemInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "system-info",
		Short: "Show Cloud system information",
		Long: `Show Confluence Cloud system information for the current tenant.

Examples:
  confluence settings system-info
  confluence settings system-info --json`,
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
			info, err := conf.GetSystemInfo(ctx, c)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(info)
			}
			printSystemInfo(w, *info)
			return nil
		},
	}
	return cmd
}

func settingsLookAndFeelCmd() *cobra.Command {
	var space string
	cmd := &cobra.Command{
		Use:   "lookandfeel",
		Short: "Show Cloud look-and-feel settings",
		Long: `Show global or space-specific Confluence Cloud look-and-feel settings.

Examples:
  confluence settings lookandfeel
  confluence settings lookandfeel --space ENG --json`,
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
			settings, err := conf.GetLookAndFeelSettings(ctx, c, space)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(settings)
			}
			printLookAndFeelSettings(w, *settings)
			return nil
		},
	}
	cmd.Flags().StringVar(&space, "space", "", "Space key; omit for global look-and-feel settings")
	return cmd
}

func settingsSpaceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "space <space-key>",
		Short: "Show Cloud space settings",
		Long: `Show settings for one Confluence Cloud space.

Examples:
  confluence settings space ENG
  confluence settings space ENG --json`,
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
			settings, err := conf.GetSpaceSettings(ctx, c, args[0])
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(settings)
			}
			printSpaceSettings(w, *settings)
			return nil
		},
	}
	return cmd
}

type settingsTextWriter interface {
	Text(format string, args ...any)
}

func printSystemInfo(w settingsTextWriter, info conf.SystemInfo) {
	w.Text("%s\t%s\t%s\t%s\t%s\n",
		firstNonEmpty(info.CloudID, "-"),
		firstNonEmpty(info.SiteTitle, "-"),
		firstNonEmpty(info.Edition, "-"),
		firstNonEmpty(info.DefaultLocale, "-"),
		firstNonEmpty(info.DefaultTimeZone, "-"))
}

func printLookAndFeelSettings(w settingsTextWriter, settings conf.LookAndFeelSettings) {
	w.Text("headings=%s\tlinks=%s\tborders=%s\n",
		firstNonEmpty(settings.Headings.Color, "-"),
		firstNonEmpty(settings.Links.Color, "-"),
		firstNonEmpty(settings.BordersAndDividers.Color, "-"))
}

func printSpaceSettings(w settingsTextWriter, settings conf.SpaceSettings) {
	w.Text("%s\trouteOverride=%t\tpageEditor=%s\tblogpostEditor=%s\tdefaultEditor=%s\n",
		firstNonEmpty(settings.SpaceKey, "-"),
		settings.RouteOverrideEnabled,
		firstNonEmpty(settings.Editor.Page, "-"),
		firstNonEmpty(settings.Editor.Blogpost, "-"),
		firstNonEmpty(settings.Editor.Default, "-"))
}
