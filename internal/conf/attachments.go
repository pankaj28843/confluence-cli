package conf

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

// Attachment is one row of /rest/api/content/{id}/child/attachment.
type Attachment struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Type    string `json:"type,omitempty"` // always "attachment"
	Status  string `json:"status,omitempty"`
	Version struct {
		Number int    `json:"number,omitempty"`
		When   string `json:"when,omitempty"`
	} `json:"version,omitempty"`
	Metadata struct {
		MediaType string `json:"mediaType,omitempty"`
		Comment   string `json:"comment,omitempty"`
	} `json:"metadata,omitempty"`
	Extensions struct {
		FileSize  int64  `json:"fileSize,omitempty"`
		MediaType string `json:"mediaType,omitempty"`
	} `json:"extensions,omitempty"`
	Links struct {
		Download string `json:"download,omitempty"`
		WebUI    string `json:"webui,omitempty"`
		Base     string `json:"base,omitempty"`
	} `json:"_links,omitempty"`
}

// ListAttachments pages /rest/api/content/{id}/child/attachment.
func ListAttachments(ctx context.Context, c *client.Client, contentID string, limit int) ([]Attachment, error) {
	if limit <= 0 {
		limit = 50
	}
	params := url.Values{
		"limit":  {strconv.Itoa(limit)},
		"expand": {"version,metadata,extensions"},
	}
	data, _, err := c.Get(ctx, "/rest/api/content/"+url.PathEscape(contentID)+"/child/attachment", params)
	if err != nil {
		return nil, err
	}
	var page struct {
		Results []Attachment `json:"results"`
	}
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parse attachments: %w", err)
	}
	return page.Results, nil
}

// DownloadURL returns the absolute download URL for an attachment.
func (a *Attachment) DownloadURL() string {
	if a.Links.Base == "" {
		return a.Links.Download
	}
	return a.Links.Base + a.Links.Download
}

// DownloadAttachment fetches the raw bytes of an attachment via its download
// link. Server/DC and Cloud both serve this on the same Authorization.
func DownloadAttachment(ctx context.Context, c *client.Client, downloadPath string) ([]byte, error) {
	if downloadPath == "" {
		return nil, fmt.Errorf("empty download path")
	}
	data, _, err := c.Get(ctx, downloadPath, nil)
	return data, err
}
