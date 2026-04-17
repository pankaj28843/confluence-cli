package conf

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

// UpdatePageInput collects fields for a PUT /rest/api/content/{id}.
type UpdatePageInput struct {
	ID            string
	Title         string
	BodyFormat    string // storage | wiki | view
	BodyValue     string
	VersionNumber int
	Type          string // page | blogpost; defaults to page
}

// UpdatePage bumps the version number and writes new title / body.
func UpdatePage(ctx context.Context, c *client.Client, in UpdatePageInput) (*Content, error) {
	if in.ID == "" {
		return nil, fmt.Errorf("UpdatePage: ID is required")
	}
	if in.BodyFormat == "" {
		in.BodyFormat = "storage"
	}
	if in.Type == "" {
		in.Type = "page"
	}
	body := map[string]any{
		"version": map[string]any{"number": in.VersionNumber + 1},
		"type":    in.Type,
	}
	if in.Title != "" {
		body["title"] = in.Title
	}
	if in.BodyValue != "" {
		body["body"] = map[string]any{
			in.BodyFormat: map[string]any{
				"value":          in.BodyValue,
				"representation": in.BodyFormat,
			},
		}
	}
	raw, _, err := c.Put(ctx, "/rest/api/content/"+url.PathEscape(in.ID), nil, body)
	if err != nil {
		return nil, err
	}
	var out Content
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("parse update response: %w", err)
	}
	return &out, nil
}

// UploadAttachment POSTs a multipart form to
// /rest/api/content/{id}/child/attachment with X-Atlassian-Token: no-check.
// `comment` is optional. Existing files with the same filename become a new
// version; to update use /child/attachment/{attId}/data instead (out of scope v1).
func UploadAttachment(ctx context.Context, c *client.Client, contentID, filename string, r io.Reader, comment string) ([]Attachment, error) {
	if filename == "" {
		return nil, fmt.Errorf("UploadAttachment: filename is required")
	}
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)

	// "file" field — must be named "file" per Confluence API.
	fh := make(textproto.MIMEHeader)
	fh.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename=%q`, filepath.Base(filename)))
	fh.Set("Content-Type", "application/octet-stream")
	part, err := mw.CreatePart(fh)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, r); err != nil {
		return nil, err
	}
	if comment != "" {
		if err := mw.WriteField("comment", comment); err != nil {
			return nil, err
		}
	}
	if err := mw.Close(); err != nil {
		return nil, err
	}

	path := "/rest/api/content/" + url.PathEscape(contentID) + "/child/attachment"
	data, _, err := c.PostRaw(ctx, path, nil, buf.Bytes(), mw.FormDataContentType(), map[string]string{
		"X-Atlassian-Token": "no-check",
	})
	if err != nil {
		return nil, err
	}
	// Response shape is {results: [attachments]}.
	var page struct {
		Results []Attachment `json:"results"`
	}
	if err := json.Unmarshal(data, &page); err != nil {
		// Server/DC sometimes returns a single Attachment instead of an envelope.
		var single Attachment
		if errSingle := json.Unmarshal(data, &single); errSingle == nil && single.ID != "" {
			return []Attachment{single}, nil
		}
		return nil, fmt.Errorf("parse upload response: %w (body=%s)", err, strings.TrimSpace(string(data)))
	}
	return page.Results, nil
}
