package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pankaj28843/confluence-cli/internal/conf"
)

func pageCreateCmd() *cobra.Command {
	var space, title, parent, bodyFormat, bodyFile, bodyInline string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new page",
		Long: `Create a new page with a storage-format body. Pass --parent to create the
page below an existing page.

Examples:
  confluence page create --space ENG --title "Runbook" --body-file body.html
  confluence page create --space ENG --title "Child" --parent 12345 --body-file body.html
  echo "<p>Hello</p>" | confluence page create --space ENG --title "Hello" --body-file -`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()

			if space == "" || title == "" {
				return fmt.Errorf("--space and --title are required")
			}
			body, err := resolveBody(bodyFile, bodyInline)
			if err != nil {
				return err
			}
			if body == "" {
				return fmt.Errorf("--body-file or --body is required")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			created, err := conf.CreatePage(ctx, c, conf.CreatePageInput{
				SpaceKey: space, Title: title, BodyFormat: bodyFormat, BodyValue: body, ParentID: parent,
			})
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(created)
			}
			w.Text("created %s", created.ID)
			if created.Version.Number > 0 {
				w.Text(" (v%d)", created.Version.Number)
			}
			if u := created.AbsoluteURL(); u != "" {
				w.Text("\n%s", u)
			}
			w.Text("\n")
			return nil
		},
	}
	cmd.Flags().StringVar(&space, "space", "", "Space key (required)")
	cmd.Flags().StringVar(&title, "title", "", "Page title (required)")
	cmd.Flags().StringVar(&parent, "parent", "", "Parent page id")
	cmd.Flags().StringVar(&bodyFormat, "body-format", "storage", "Body format: storage | wiki | view")
	cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to body file, or '-' for stdin")
	cmd.Flags().StringVar(&bodyInline, "body", "", "Inline body string")
	return cmd
}

func pageUpdateCmd() *cobra.Command {
	var title, bodyFormat, bodyFile, bodyInline, newVersion string
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an existing page (title and/or body)",
		Long: `Update a page title, body, or both. The command fetches the current page
version and writes version.number + 1 unless --version is supplied.

Examples:
  confluence page update 12345 --title "New Title"
  confluence page update 12345 --body-file body.html
  echo "<p>Hello</p>" | confluence page update 12345 --body-format storage --body-file -`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()

			body, err := resolveBody(bodyFile, bodyInline)
			if err != nil {
				return err
			}
			if title == "" && body == "" {
				return fmt.Errorf("at least one of --title, --body-file, or --body is required")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			page, err := conf.GetContent(ctx, c, args[0], "body.storage,version")
			if err != nil {
				return err
			}
			if title == "" {
				title = page.Title
			}
			ver := page.Version.Number
			if newVersion != "" {
				if _, err := fmt.Sscanf(newVersion, "%d", &ver); err != nil || ver < 0 {
					return fmt.Errorf("invalid --version %q", newVersion)
				}
			}
			out, err := conf.UpdatePage(ctx, c, conf.UpdatePageInput{
				ID: args[0], Title: title, BodyFormat: bodyFormat, BodyValue: body, VersionNumber: ver,
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
	cmd.Flags().StringVar(&bodyInline, "body", "", "Inline body string")
	cmd.Flags().StringVar(&newVersion, "version", "", "Explicit current version number (auto-fetched if omitted)")
	return cmd
}

func pagePublishCmd() *cobra.Command {
	var bodyFile, bodyFormat, title string
	var attaches []string
	cmd := &cobra.Command{
		Use:   "publish <id>",
		Short: "Upload attachments, then update page body",
		Long: `Publish a page body and any referenced attachments in order. Attachments are
created or updated by filename before the page body is PUT with an incremented
version.

Examples:
  confluence page publish 12345 --body-file page.html --attach hld.png
  confluence page publish 12345 --body-file page.html --attach hld.png --attach flow.png
  confluence page publish 12345 --title "Runbook" --body-format storage --body-file page.html`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()

			if bodyFile == "" {
				return fmt.Errorf("--body-file is required")
			}
			body, err := readBodyFile(bodyFile)
			if err != nil {
				return err
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			var uploaded []conf.Attachment
			for _, path := range attaches {
				f, err := os.Open(path)
				if err != nil {
					return fmt.Errorf("open attachment %s: %w", path, err)
				}
				atts, err := conf.PutAttachment(ctx, c, args[0], path, f, "")
				closeErr := f.Close()
				if err != nil {
					return fmt.Errorf("upload attachment %s: %w", path, err)
				}
				if closeErr != nil {
					return fmt.Errorf("close attachment %s: %w", path, closeErr)
				}
				uploaded = append(uploaded, atts...)
			}
			page, err := conf.GetContent(ctx, c, args[0], "body.storage,version")
			if err != nil {
				return err
			}
			if title == "" {
				title = page.Title
			}
			updated, err := conf.UpdatePage(ctx, c, conf.UpdatePageInput{
				ID: args[0], Title: title, BodyFormat: bodyFormat, BodyValue: body, VersionNumber: page.Version.Number,
			})
			if err != nil {
				return err
			}
			out := struct {
				Page        *conf.Content     `json:"page"`
				Attachments []conf.Attachment `json:"attachments"`
				Version     int               `json:"version"`
			}{updated, uploaded, updated.Version.Number}
			if w.IsJSON() {
				return w.JSON(out)
			}
			w.Text("published %s (v%d), attachments=%d\n", updated.ID, updated.Version.Number, len(uploaded))
			return nil
		},
	}
	cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to body file, or '-' for stdin (required)")
	cmd.Flags().StringVar(&bodyFormat, "body-format", "storage", "Body format: storage | wiki | view")
	cmd.Flags().StringVar(&title, "title", "", "New title (keeps existing if omitted)")
	cmd.Flags().StringArrayVar(&attaches, "attach", nil, "Attachment file to create or update; repeatable")
	return cmd
}

func attachmentUploadCmd() *cobra.Command { return attachmentPutCmd("upload") }

func attachmentReplaceCmd() *cobra.Command { return attachmentPutCmd("replace") }

func attachmentPutCmd(verb string) *cobra.Command {
	var page, file, fileName, comment string
	short := "Upload an attachment to a page"
	if verb == "replace" {
		short = "Create or replace an attachment on a page"
	}
	cmd := &cobra.Command{
		Use:   verb,
		Short: short,
		Long: `Upload an attachment. If the filename already exists on the page, the
existing attachment data is updated as a new version.

Examples:
  confluence attachment ` + verb + ` --page 12345 --file ./report.pdf
  confluence attachment ` + verb + ` --page 12345 --file ./logo.png --comment "v2 logo"
  cat report.pdf | confluence attachment ` + verb + ` --page 12345 --file - --file-name report.pdf`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()

			if page == "" || file == "" {
				return fmt.Errorf("--page and --file are required")
			}
			reader, filename, cleanup, err := attachmentReader(file, fileName)
			if err != nil {
				return err
			}
			defer cleanup()
			c, err := newClient()
			if err != nil {
				return err
			}
			atts, err := conf.PutAttachment(ctx, c, page, filename, reader, comment)
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
	cmd.Flags().StringVar(&file, "file", "", "Local file path, or '-' for stdin (required)")
	cmd.Flags().StringVar(&fileName, "file-name", "", "Filename to use when --file is '-'")
	cmd.Flags().StringVar(&comment, "comment", "", "Optional attachment comment")
	return cmd
}

func attachmentDeleteCmd() *cobra.Command {
	var id, page, name string
	var force bool
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an attachment by id or by page/name",
		Long: `Delete an attachment content entity. Pass --id directly, or pass --page and
--name to resolve the attachment id first. A confirmation prompt is shown unless
--force is supplied.

Examples:
  confluence attachment delete --id 1884909332 --force
  confluence attachment delete --page 12345 --name old-diagram.png
  confluence attachment delete --page 12345 --name old-diagram.png --force --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()

			if id == "" && (page == "" || name == "") {
				return fmt.Errorf("pass --id, or pass both --page and --name")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			title := name
			if id == "" {
				match, err := conf.FindAttachmentByTitle(ctx, c, page, name)
				if err != nil {
					return err
				}
				if match == nil {
					return fmt.Errorf("no attachment named %q on page %s", name, page)
				}
				id = match.ID
				title = match.Title
			}
			if !force && !confirmDelete(id, title) {
				return fmt.Errorf("delete cancelled")
			}
			if err := conf.DeleteContent(ctx, c, id); err != nil {
				return err
			}
			out := struct {
				Deleted bool   `json:"deleted"`
				ID      string `json:"id"`
				Title   string `json:"title,omitempty"`
			}{true, id, title}
			if w.IsJSON() {
				return w.JSON(out)
			}
			w.Text("deleted %s", id)
			if title != "" {
				w.Text("\t%s", title)
			}
			w.Text("\n")
			return nil
		},
	}
	cmd.Flags().StringVar(&id, "id", "", "Attachment content id")
	cmd.Flags().StringVar(&page, "page", "", "Page id used with --name")
	cmd.Flags().StringVar(&name, "name", "", "Attachment filename used with --page")
	cmd.Flags().BoolVar(&force, "force", false, "Delete without confirmation")
	return cmd
}

func runPageScreenshot(ctx context.Context, pageURL, pageID, out string, newTab bool) error {
	if _, err := exec.LookPath("cdp"); err != nil {
		return fmt.Errorf("cdp is required for page screenshot: %w", err)
	}
	openArgs := []string{"open", pageURL, "--new-tab=" + fmt.Sprint(newTab)}
	if !newTab {
		openArgs = append(openArgs, "--url-contains", pageID)
	}
	if data, err := exec.CommandContext(ctx, "cdp", openArgs...).CombinedOutput(); err != nil {
		return fmt.Errorf("cdp open: %w: %s", err, strings.TrimSpace(string(data)))
	}
	if data, err := exec.CommandContext(ctx, "cdp", "wait", "selector", "body", "--url-contains", pageID).CombinedOutput(); err != nil {
		return fmt.Errorf("cdp wait: %w: %s", err, strings.TrimSpace(string(data)))
	}
	if data, err := exec.CommandContext(ctx, "cdp", "screenshot", "--url-contains", pageID, "--out", out, "--full-page").CombinedOutput(); err != nil {
		return fmt.Errorf("cdp screenshot: %w: %s", err, strings.TrimSpace(string(data)))
	}
	return nil
}

func resolveBody(bodyFile, bodyInline string) (string, error) {
	if bodyFile != "" && bodyInline != "" {
		return "", fmt.Errorf("use only one of --body-file or --body")
	}
	if bodyFile == "" {
		return bodyInline, nil
	}
	return readBodyFile(bodyFile)
}

func readBodyFile(path string) (string, error) {
	var data []byte
	var err error
	if path == "-" {
		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(path)
	}
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}
	return string(data), nil
}

func attachmentReader(file, fileName string) (io.Reader, string, func(), error) {
	if file == "-" {
		if fileName == "" {
			return nil, "", func() {}, fmt.Errorf("--file - requires --file-name")
		}
		return os.Stdin, fileName, func() {}, nil
	}
	f, err := os.Open(file)
	if err != nil {
		return nil, "", func() {}, err
	}
	name := fileName
	if name == "" {
		name = filepath.Base(file)
	}
	return f, name, func() { _ = f.Close() }, nil
}

func confirmDelete(id, title string) bool {
	fmt.Fprintf(os.Stderr, "Delete attachment %s", id)
	if title != "" {
		fmt.Fprintf(os.Stderr, " (%s)", title)
	}
	fmt.Fprint(os.Stderr, "? Type 'delete' to confirm: ")
	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	return err == nil && strings.TrimSpace(line) == "delete"
}
