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
	if isCloud(c) {
		return listAttachmentsCloudV2(ctx, c, contentID, limit)
	}
	return listAttachmentsServerV1(ctx, c, contentID, limit)
}

func listAttachmentsServerV1(ctx context.Context, c *client.Client, contentID string, limit int) ([]Attachment, error) {
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

// DownloadAttachment fetches raw attachment bytes. Server/DC uses the returned
// download link. Cloud v2 returns a browser-facing downloadLink, so Cloud uses
// the v1 raw-content route that accepts Basic email:token auth.
func DownloadAttachment(ctx context.Context, c *client.Client, contentID string, attachment Attachment) ([]byte, error) {
	downloadPath := attachment.Links.Download
	if isCloud(c) {
		if contentID == "" {
			return nil, fmt.Errorf("DownloadAttachment: content ID is required")
		}
		if attachment.ID == "" {
			return nil, fmt.Errorf("DownloadAttachment: attachment ID is required")
		}
		downloadPath = "/rest/api/content/" + url.PathEscape(contentID) + "/child/attachment/" + url.PathEscape(attachment.ID) + "/download"
	}
	if downloadPath == "" {
		return nil, fmt.Errorf("DownloadAttachment: download path is required")
	}
	data, _, err := c.Get(ctx, downloadPath, nil)
	return data, err
}
