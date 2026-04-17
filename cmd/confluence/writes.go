package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/pankaj28843/confluence-cli/internal/conf"
)

// pageUpdateCmdReal replaces the Phase-3 stub for `confluence page update`.
func pageUpdateCmdReal() *cobra.Command {
	var title, bodyFormat, bodyFile, bodyInline, newVersion string
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an existing page (title and/or body)",
		Long: `Update a page — title, body, or both. Supply the body via --body-file
(path, or '-' for stdin) or --body (inline string). Body format defaults to
'storage'; 'wiki' is also common.

Version handling:
  If --version is provided, it is used verbatim. Otherwise we GET the page
  first, read its current version.number, and PUT version.number + 1.

Examples:
  confluence page update 12345 --title "New Title"
  echo "<p>Hello</p>" | confluence page update 12345 --body-format storage --body-file -
  confluence page update 12345 --body-format wiki --body-file body.txt --version 4`,
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

			body := bodyInline
			if bodyFile != "" {
				var data []byte
				if bodyFile == "-" {
					data, err = io.ReadAll(os.Stdin)
				} else {
					data, err = os.ReadFile(bodyFile)
				}
				if err != nil {
					return fmt.Errorf("read body: %w", err)
				}
				body = string(data)
			}

			// Resolve version: GET the page if --version wasn't specified.
			ver := 0
			if newVersion != "" {
				v, err := strconv.Atoi(newVersion)
				if err != nil || v < 0 {
					return fmt.Errorf("invalid --version %q", newVersion)
				}
				ver = v
			} else {
				page, err := conf.GetContent(ctx, c, args[0], "version")
				if err != nil {
					return err
				}
				ver = page.Version.Number
				// Use existing title/type if caller didn't override.
				if title == "" {
					title = page.Title
				}
			}

			out, err := conf.UpdatePage(ctx, c, conf.UpdatePageInput{
				ID:            args[0],
				Title:         title,
				BodyFormat:    bodyFormat,
				BodyValue:     body,
				VersionNumber: ver,
			})
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(out)
			}
			w.Text("updated %s (v%d)\n", out.ID, out.Version.Number)
			return nil
		},
	}
	cmd.Flags().StringVar(&title, "title", "", "New title (keeps existing if omitted)")
	cmd.Flags().StringVar(&bodyFormat, "body-format", "storage", "Body format: storage | wiki | view")
	cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to body file, or '-' for stdin")
	cmd.Flags().StringVar(&bodyInline, "body", "", "Inline body string (overrides --body-file)")
	cmd.Flags().StringVar(&newVersion, "version", "", "Explicit version number (auto-fetched if omitted)")
	return cmd
}

// attachmentUploadCmdReal replaces the Phase-4 stub for `confluence attachment upload`.
func attachmentUploadCmdReal() *cobra.Command {
	var page, file, comment string
	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload an attachment to a page (multipart + X-Atlassian-Token)",
		Long: `Upload a new attachment. If the filename already exists, Confluence creates
a new version of the existing attachment.

Examples:
  confluence attachment upload --page 12345 --file ./report.pdf
  confluence attachment upload --page 12345 --file ./logo.png --comment "v2 logo"
  cat report.pdf | confluence attachment upload --page 12345 --file - --comment "piped"`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()

			if page == "" || file == "" {
				return fmt.Errorf("--page and --file are required")
			}
			c, err := newClient()
			if err != nil {
				return err
			}

			var reader io.Reader
			var filename string
			if file == "-" {
				reader = os.Stdin
				filename = "stdin"
				return fmt.Errorf("stdin upload requires --file <name>; pass a filename after '-' via --file-name")
			}
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer f.Close()
			reader = f
			filename = file

			atts, err := conf.UploadAttachment(ctx, c, page, filename, reader, comment)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(atts)
			}
			for _, a := range atts {
				w.Text("%s\t%s\t%d bytes\n", a.ID, a.Title, a.Extensions.FileSize)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&page, "page", "", "Content id to attach to (required)")
	cmd.Flags().StringVar(&file, "file", "", "Local file path (required)")
	cmd.Flags().StringVar(&comment, "comment", "", "Optional attachment comment")
	return cmd
}
