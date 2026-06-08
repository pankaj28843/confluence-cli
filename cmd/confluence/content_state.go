package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pankaj28843/confluence-cli/internal/conf"
)

func contentStateCmdReal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "content-state",
		Short: "Cloud content state reads",
		Long: `Cloud content state read operations.

Content states are documented in Confluence Cloud REST API v1. Server/Data
Center does not expose the same content-state REST group in the current
official REST OpenAPI, so these typed commands are Cloud-only.

Examples:
  confluence content-state current 12345
  confluence content-state available 12345 --json
  confluence content-state content ENG --state-id 1 --limit 25`,
	}
	cmd.AddCommand(contentStateCurrentCmd())
	cmd.AddCommand(contentStateAvailableCmd())
	cmd.AddCommand(contentStateCustomCmd())
	cmd.AddCommand(contentStateSpaceCmd())
	cmd.AddCommand(contentStateSettingsCmd())
	cmd.AddCommand(contentStateContentCmd())
	return cmd
}

func contentStateCurrentCmd() *cobra.Command {
	var status string
	cmd := &cobra.Command{
		Use:   "current <content-id>",
		Short: "Show current Cloud content state",
		Long: `Show the content state attached to the draft, current, or archived
version of a Cloud content item.

Examples:
  confluence content-state current 12345
  confluence content-state current 12345 --status draft --json`,
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
			state, err := conf.GetContentState(ctx, c, args[0], status)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(state)
			}
			printContentStateResponse(w, *state)
			return nil
		},
	}
	cmd.Flags().StringVar(&status, "status", "current", "Cloud content status: current, draft, or archived")
	return cmd
}

func contentStateAvailableCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "available <content-id>",
		Short: "List Cloud states available for content",
		Long: `List Cloud content states that are available for one content item.

Examples:
  confluence content-state available 12345
  confluence content-state available 12345 --json`,
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
			states, err := conf.ListAvailableContentStates(ctx, c, args[0])
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(states)
			}
			printAvailableContentStates(w, *states)
			return nil
		},
	}
	return cmd
}

func contentStateCustomCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "custom",
		Short: "List Cloud custom content states",
		Long: `List Cloud custom content states created by the authenticated user.

Examples:
  confluence content-state custom
  confluence content-state custom --json`,
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
			states, err := conf.ListCustomContentStates(ctx, c)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(states)
			}
			printContentStates(w, "custom", states)
			return nil
		},
	}
	return cmd
}

func contentStateSpaceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "space <space-key>",
		Short: "List suggested Cloud states for a space",
		Long: `List Cloud content states suggested in one space.

Examples:
  confluence content-state space ENG
  confluence content-state space ENG --json`,
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
			states, err := conf.ListSpaceContentStates(ctx, c, args[0])
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(states)
			}
			printContentStates(w, "space", states)
			return nil
		},
	}
	return cmd
}

func contentStateSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "settings <space-key>",
		Short: "Show Cloud content-state settings for a space",
		Long: `Show Cloud space settings for content states.

Examples:
  confluence content-state settings ENG
  confluence content-state settings ENG --json`,
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
			settings, err := conf.GetContentStateSettings(ctx, c, args[0])
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(settings)
			}
			printContentStateSettings(w, *settings)
			return nil
		},
	}
	return cmd
}

func contentStateContentCmd() *cobra.Command {
	var stateID int64
	var limit int
	var start int
	var expand []string
	cmd := &cobra.Command{
		Use:   "content <space-key>",
		Short: "List Cloud content with a given state",
		Long: `List Cloud content in one space that has the provided content state.

Examples:
  confluence content-state content ENG --state-id 1
  confluence content-state content ENG --state-id 1 --expand space --expand version --json
  confluence content-state content ENG --state-id 0 --start 25 --limit 25`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			if !cmd.Flags().Changed("state-id") {
				return fmt.Errorf("--state-id is required")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			content, err := conf.ListContentWithState(ctx, c, conf.ContentWithStateOptions{
				SpaceKey: args[0],
				StateID:  &stateID,
				Expand:   expand,
				Start:    start,
				Limit:    limit,
			})
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(content)
			}
			printStateContent(w, content)
			return nil
		},
	}
	cmd.Flags().Int64Var(&stateID, "state-id", 0, "Content state id; required, and 0 is valid for Cloud default states")
	cmd.Flags().IntVar(&limit, "limit", 25, "Max content rows (Cloud endpoint hard cap 100)")
	cmd.Flags().IntVar(&start, "start", 0, "Zero-based result offset")
	cmd.Flags().StringSliceVar(&expand, "expand", nil, "Expand value; repeatable or comma-separated")
	return cmd
}

type contentStateTextWriter interface {
	Text(format string, args ...any)
}

func printContentStateResponse(w contentStateTextWriter, state conf.ContentStateResponse) {
	if state.State == nil {
		w.Text("none\t-\t-\t%s\n", firstNonEmpty(state.LastUpdated, "-"))
		return
	}
	w.Text("%d\t%s\t%s\t%s\n", state.State.ID, firstNonEmpty(state.State.Name, "-"), firstNonEmpty(state.State.Color, "-"), firstNonEmpty(state.LastUpdated, "-"))
}

func printAvailableContentStates(w contentStateTextWriter, states conf.AvailableContentStates) {
	printContentStates(w, "space", states.SpaceContentStates)
	printContentStates(w, "custom", states.CustomContentStates)
}

func printContentStates(w contentStateTextWriter, kind string, states []conf.ContentState) {
	for _, state := range states {
		w.Text("%s\t%d\t%s\t%s\n", kind, state.ID, firstNonEmpty(state.Name, "-"), firstNonEmpty(state.Color, "-"))
	}
}

func printContentStateSettings(w contentStateTextWriter, settings conf.ContentStateSettings) {
	w.Text("allowed=%t\tcustom=%t\tspace=%t\n", settings.ContentStatesAllowed, settings.CustomContentStatesAllowed, settings.SpaceContentStatesAllowed)
	printContentStates(w, "space", settings.SpaceContentStates)
}

func printStateContent(w contentStateTextWriter, content []conf.Content) {
	for _, item := range content {
		w.Text("%s\t%s\t%s\t%s\n", item.ID, firstNonEmpty(item.Type, "-"), firstNonEmpty(item.Status, "-"), firstNonEmpty(item.Title, "-"))
	}
}
