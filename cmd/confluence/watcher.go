package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pankaj28843/confluence-cli/internal/conf"
)

func watcherCmdReal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watcher",
		Short: "Watchers (content, space, status)",
		Long: `Watcher read helpers.

Examples:
  confluence watcher content --page 12345
  confluence watcher space --space ENG
  confluence watcher status --page 12345 --json`,
	}
	cmd.AddCommand(watcherContentCmd())
	cmd.AddCommand(watcherSpaceCmd())
	cmd.AddCommand(watcherStatusCmd())
	return cmd
}

func watcherContentCmd() *cobra.Command {
	var page string
	var limit int
	cmd := &cobra.Command{
		Use:   "content",
		Short: "List watchers subscribed to a content id",
		Long: `List watchers subscribed to a content id.

Examples:
  confluence watcher content --page 12345
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
			watchers, err := conf.ListContentWatchers(ctx, c, page, limit)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(watchers)
			}
			for _, watcher := range watchers.Results {
				w.Text("%s\t%s\n", watcherSubject(watcher), watcher.Type)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&page, "page", "", "Content id (required)")
	cmd.Flags().IntVar(&limit, "limit", 25, "Maximum watchers returned")
	return cmd
}

func watcherSpaceCmd() *cobra.Command {
	var space string
	var limit int
	cmd := &cobra.Command{
		Use:   "space",
		Short: "List watchers subscribed to a space",
		Long: `List watchers subscribed to a space key.

Examples:
  confluence watcher space --space ENG
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
			watchers, err := conf.ListSpaceWatchers(ctx, c, space, limit)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(watchers)
			}
			for _, watcher := range watchers.Results {
				w.Text("%s\t%s\n", watcherSubject(watcher), watcher.Type)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&space, "space", "", "Space key (required)")
	cmd.Flags().IntVar(&limit, "limit", 25, "Maximum watchers returned")
	return cmd
}

func watcherStatusCmd() *cobra.Command {
	var page string
	var space string
	var accountID string
	var userKey string
	var username string
	var contentType string
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show whether a user watches content or a space",
		Long: `Show whether the current or specified user watches a content id or
space key.

Examples:
  confluence watcher status --page 12345
  confluence watcher status --space ENG --content-type blogpost --json
  confluence watcher status --page 12345 --account-id abc123`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if (page == "" && space == "") || (page != "" && space != "") {
				return fmt.Errorf("set exactly one of --page or --space")
			}
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			c, err := newClient()
			if err != nil {
				return err
			}
			opts := conf.WatchStatusOptions{
				AccountID:   accountID,
				UserKey:     userKey,
				Username:    username,
				ContentType: contentType,
			}
			var status *conf.WatchStatus
			if page != "" {
				status, err = conf.GetContentWatchStatus(ctx, c, page, opts)
			} else {
				status, err = conf.GetSpaceWatchStatus(ctx, c, space, opts)
			}
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(status)
			}
			w.Text("watching=%t\n", status.Watching)
			return nil
		},
	}
	cmd.Flags().StringVar(&page, "page", "", "Content id")
	cmd.Flags().StringVar(&space, "space", "", "Space key")
	cmd.Flags().StringVar(&accountID, "account-id", "", "Cloud account id to check")
	cmd.Flags().StringVar(&userKey, "user-key", "", "Server/DC user key to check")
	cmd.Flags().StringVar(&username, "username", "", "Server/DC username to check")
	cmd.Flags().StringVar(&contentType, "content-type", "", "Space watch content type, such as page or blogpost")
	return cmd
}

func watcherSubject(watcher conf.WatchRecord) string {
	user := watcher.Watcher
	switch {
	case user.DisplayName != "":
		return user.DisplayName
	case user.PublicName != "":
		return user.PublicName
	case user.Username != "":
		return user.Username
	case user.AccountID != "":
		return user.AccountID
	case user.UserKey != "":
		return user.UserKey
	default:
		return "-"
	}
}
