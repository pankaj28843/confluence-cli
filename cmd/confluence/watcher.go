package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pankaj28843/confluence-cli/internal/conf"
)

func watcherCmdReal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watcher",
		Short: "Watchers (content, space)",
		Long: `Watcher operations.

Examples:
  confluence watcher content --page 12345
  confluence watcher space --space ENG`,
	}
	cmd.AddCommand(watcherContentCmd())
	cmd.AddCommand(watcherSpaceCmd())
	return cmd
}

func watcherContentCmd() *cobra.Command {
	var page string
	cmd := &cobra.Command{
		Use:   "content",
		Short: "Watchers subscribed to a content id (raw passthrough)",
		Long: `Show watcher/notification records for a content id. The shape varies
between Confluence versions, so output is passed through verbatim.

Examples:
  confluence watcher content --page 12345 --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			if page == "" {
				return fmt.Errorf("--page is required")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := conf.GetContentWatchers(ctx, c, page)
			if err != nil {
				return err
			}
			_, _ = w.Out.Write(data)
			return nil
		},
	}
	cmd.Flags().StringVar(&page, "page", "", "Content id (required)")
	return cmd
}

func watcherSpaceCmd() *cobra.Command {
	var space string
	cmd := &cobra.Command{
		Use:   "space",
		Short: "Watchers subscribed to a space (raw passthrough)",
		Long: `Show watcher records for a space key.

Examples:
  confluence watcher space --space ENG --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			if space == "" {
				return fmt.Errorf("--space is required")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			data, err := conf.GetSpaceWatchers(ctx, c, space)
			if err != nil {
				return err
			}
			_, _ = w.Out.Write(data)
			return nil
		},
	}
	cmd.Flags().StringVar(&space, "space", "", "Space key (required)")
	return cmd
}
