package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/pankaj28843/confluence-cli/internal/conf"
)

func bodyCmdReal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "body",
		Short: "Content body utilities",
		Long: `Content body utility operations.

Examples:
  confluence body convert --from storage --to view --value '<p>Hello</p>'
  confluence body convert --to export_view --value @body-storage.xml --json`,
	}
	cmd.AddCommand(bodyConvertCmd())
	return cmd
}

func bodyConvertCmd() *cobra.Command {
	var from, to, valueFlag, spaceContext, contentContext, embeddedRender string
	var expand []string
	var noCache bool
	var pollAttempts int
	var pollInterval time.Duration
	cmd := &cobra.Command{
		Use:   "convert",
		Short: "Convert a content body representation",
		Long: `Convert a Confluence content body between documented representations.

Cloud uses the asynchronous content-body conversion API and polls for the
result by default. Use --poll-attempts 0 to return only the Cloud async id.

Examples:
  confluence body convert --from storage --to view --value '<p>Hello</p>'
  confluence body convert --from storage --to export_view --value @body.xml --json
  confluence body convert --to view --value @- --expand webresource.uris.css --expand webresource.uris.js`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			if from == "" || to == "" || valueFlag == "" {
				return fmt.Errorf("--from, --to, and --value are required")
			}
			raw, err := readDataArg(valueFlag)
			if err != nil {
				return err
			}
			if len(raw) == 0 {
				return fmt.Errorf("--value is empty")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			input := conf.BodyConversionInput{
				From:                  from,
				To:                    to,
				Value:                 string(raw),
				Expand:                expand,
				SpaceKeyContext:       spaceContext,
				ContentIDContext:      contentContext,
				EmbeddedContentRender: embeddedRender,
				CloudPollAttempts:     pollAttempts,
				CloudPollInterval:     pollInterval,
			}
			if noCache {
				allowCache := false
				input.AllowCache = &allowCache
			}
			got, err := conf.ConvertBody(ctx, c, input)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(got)
			}
			printBodyConversion(w, *got)
			return nil
		},
	}
	cmd.Flags().StringVar(&from, "from", "storage", "Source representation")
	cmd.Flags().StringVar(&to, "to", "view", "Target representation")
	cmd.Flags().StringVar(&valueFlag, "value", "", "Body value, @file, or @- (required)")
	cmd.Flags().StringSliceVar(&expand, "expand", nil, "Expand value; repeatable or comma-separated")
	cmd.Flags().StringVar(&spaceContext, "space-context", "", "Cloud space key context for permission-sensitive conversion")
	cmd.Flags().StringVar(&contentContext, "content-context", "", "Cloud content id context for permission-sensitive conversion")
	cmd.Flags().StringVar(&embeddedRender, "embedded-render", "", "Cloud embeddedContentRender value, for example current")
	cmd.Flags().BoolVar(&noCache, "no-cache", false, "Cloud only: queue conversion with allowCache=false")
	cmd.Flags().IntVar(&pollAttempts, "poll-attempts", 10, "Cloud poll attempts after queueing; 0 returns async id only")
	cmd.Flags().DurationVar(&pollInterval, "poll-interval", 500*time.Millisecond, "Cloud poll interval")
	return cmd
}

type bodyTextWriter interface {
	Text(format string, args ...interface{})
}

func printBodyConversion(w bodyTextWriter, got conf.BodyConversion) {
	if got.Value != "" {
		w.Text("%s\n", got.Value)
		return
	}
	if got.AsyncID != "" {
		if got.Status != "" {
			w.Text("%s\t%s\n", got.AsyncID, got.Status)
			return
		}
		w.Text("%s\n", got.AsyncID)
		return
	}
	if got.Status != "" {
		w.Text("%s\n", got.Status)
	}
}
