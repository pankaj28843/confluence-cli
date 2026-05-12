package conf

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

// CreatePageInput collects fields for POST /rest/api/content.
type CreatePageInput struct {
	SpaceKey   string
	Title      string
	BodyFormat string // storage | wiki | view
	BodyValue  string
	ParentID   string
	Type       string // page | blogpost; defaults to page
}

// CreatePage creates a page or blogpost with the supplied body.
func CreatePage(ctx context.Context, c *client.Client, in CreatePageInput) (*Content, error) {
	if in.SpaceKey == "" {
		return nil, fmt.Errorf("CreatePage: space key is required")
	}
	if in.Title == "" {
		return nil, fmt.Errorf("CreatePage: title is required")
	}
	if in.BodyFormat == "" {
		in.BodyFormat = "storage"
	}
	if in.Type == "" {
		in.Type = "page"
	}
	body := map[string]any{
		"type":  in.Type,
		"title": in.Title,
		"space": map[string]any{"key": in.SpaceKey},
		"body": map[string]any{
			in.BodyFormat: map[string]any{
				"value":          in.BodyValue,
				"representation": in.BodyFormat,
			},
		},
	}
	if in.ParentID != "" {
		body["ancestors"] = []map[string]string{{"id": in.ParentID}}
	}
	data, _, err := c.Post(ctx, "/rest/api/content", nil, body)
	if err != nil {
		return nil, err
	}
	var out Content
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parse create response: %w", err)
	}
	return &out, nil
}

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

// DeleteContent deletes any Confluence content entity, including attachments.
func DeleteContent(ctx context.Context, c *client.Client, id string) error {
	if id == "" {
		return fmt.Errorf("DeleteContent: ID is required")
	}
	_, _, err := c.Delete(ctx, "/rest/api/content/"+url.PathEscape(id), nil)
	return err
}

// UploadAttachment creates an attachment or updates an existing attachment with
// the same filename on the page.
func UploadAttachment(ctx context.Context, c *client.Client, contentID, filename string, r io.Reader, comment string) ([]Attachment, error) {
	return PutAttachment(ctx, c, contentID, filename, r, comment)
}

// PutAttachment creates an attachment or updates an existing attachment with the
// same filename on the page.
func PutAttachment(ctx context.Context, c *client.Client, contentID, filename string, r io.Reader, comment string) ([]Attachment, error) {
	if contentID == "" {
		return nil, fmt.Errorf("PutAttachment: content ID is required")
	}
	if filename == "" {
		return nil, fmt.Errorf("PutAttachment: filename is required")
	}
	existing, err := FindAttachmentByTitle(ctx, c, contentID, filepath.Base(filename))
	if err != nil {
		return nil, err
	}
	path := "/rest/api/content/" + url.PathEscape(contentID) + "/child/attachment"
	if existing != nil {
		path += "/" + url.PathEscape(existing.ID) + "/data"
	}
	return postAttachmentData(ctx, c, path, filename, r, comment)
}

// FindAttachmentByTitle returns the first attachment with exactly matching title.
func FindAttachmentByTitle(ctx context.Context, c *client.Client, contentID, title string) (*Attachment, error) {
	atts, err := ListAttachments(ctx, c, contentID, 200)
	if err != nil {
		return nil, err
	}
	for i := range atts {
		if atts[i].Title == title {
			return &atts[i], nil
		}
	}
	return nil, nil
}

func postAttachmentData(ctx context.Context, c *client.Client, path, filename string, r io.Reader, comment string) ([]Attachment, error) {
	pr, pw := io.Pipe()
	mw := multipart.NewWriter(pw)
	errc := make(chan error, 1)
	go func() {
		defer close(errc)
		part, err := mw.CreateFormFile("file", filepath.Base(filename))
		if err != nil {
			errc <- err
			_ = pw.CloseWithError(err)
			return
		}
		if _, err := io.Copy(part, r); err != nil {
			errc <- err
			_ = pw.CloseWithError(err)
			return
		}
		if comment != "" {
			if err := mw.WriteField("comment", comment); err != nil {
				errc <- err
				_ = pw.CloseWithError(err)
				return
			}
		}
		if err := mw.Close(); err != nil {
			errc <- err
			_ = pw.CloseWithError(err)
			return
		}
		errc <- pw.Close()
	}()

	data, _, err := c.PostRawReader(ctx, path, nil, pr, mw.FormDataContentType(), map[string]string{
		"X-Atlassian-Token": "no-check",
	})
	if writeErr := <-errc; writeErr != nil && err == nil {
		err = writeErr
	}
	if err != nil {
		return nil, err
	}
	return parseAttachmentResponse(data)
}

func parseAttachmentResponse(data []byte) ([]Attachment, error) {
	var page struct {
		Results []Attachment `json:"results"`
	}
	if err := json.Unmarshal(data, &page); err == nil && page.Results != nil {
		return page.Results, nil
	}
	var single Attachment
	if err := json.Unmarshal(data, &single); err == nil && single.ID != "" {
		return []Attachment{single}, nil
	}
	return nil, fmt.Errorf("parse attachment response: body=%s", strings.TrimSpace(string(data)))
}
