package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/pankaj28843/confluence-cli/internal/conf"
)

func attachmentCmdReal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attachment",
		Short: "Attachments (list, download, upload)",
		Long: `Attachment operations.

Examples:
  confluence attachment list --page 12345
  confluence attachment download --page 12345 --name logo.png --output ./logo.png
  confluence attachment upload --page 12345 --file ./report.pdf`,
	}
	cmd.AddCommand(attachmentListCmd())
	cmd.AddCommand(attachmentDownloadCmd())
	cmd.AddCommand(attachmentUploadCmd()) // implemented in writes.go (Phase 7)
	return cmd
}

func attachmentListCmd() *cobra.Command {
	var page string
	var limit int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List attachments on a page",
		Long: `List attachments.

Examples:
  confluence attachment list --page 12345
  confluence attachment list --page 12345 --json --limit 100`,
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
			atts, err := conf.ListAttachments(ctx, c, page, limit)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(atts)
			}
			for _, a := range atts {
				w.Text("%s\t%s\t%d bytes\t%s\n", a.ID, a.Extensions.MediaType, a.Extensions.FileSize, a.Title)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&page, "page", "", "Content id (required)")
	cmd.Flags().IntVar(&limit, "limit", 50, "Max attachments (hard cap 200)")
	return cmd
}

func attachmentDownloadCmd() *cobra.Command {
	var page, name, output string
	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download an attachment by name",
		Long: `Download an attachment. --output file writes to disk; omit for stdout.

Examples:
  confluence attachment download --page 12345 --name logo.png --output ./logo.png
  confluence attachment download --page 12345 --name report.pdf > report.pdf`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			if page == "" || name == "" {
				return fmt.Errorf("--page and --name are required")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			atts, err := conf.ListAttachments(ctx, c, page, 200)
			if err != nil {
				return err
			}
			var match *conf.Attachment
			for i := range atts {
				if atts[i].Title == name {
					match = &atts[i]
					break
				}
			}
			if match == nil {
				return fmt.Errorf("no attachment named %q on page %s", name, page)
			}
			data, err := conf.DownloadAttachment(ctx, c, match.Links.Download)
			if err != nil {
				return err
			}
			if output == "" {
				_, err = w.Out.Write(data)
				return err
			}
			return os.WriteFile(output, data, 0o644)
		},
	}
	cmd.Flags().StringVar(&page, "page", "", "Content id (required)")
	cmd.Flags().StringVar(&name, "name", "", "Attachment file name (required)")
	cmd.Flags().StringVar(&output, "output", "", "Write to file (default: stdout)")
	return cmd
}

// attachmentUploadCmd is implemented in writes.go (Phase 7). Placeholder here.
func attachmentUploadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "upload",
		Short: "Upload an attachment (Phase 7)",
		Long: `Upload a new attachment to a page.

Examples:
  confluence attachment upload --page 12345 --file ./report.pdf`,
		RunE: func(*cobra.Command, []string) error { return notImplemented("attachment upload") },
	}
}
