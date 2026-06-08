package conf

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

type cloudAttachmentV2 struct {
	ID           string `json:"id"`
	Status       string `json:"status,omitempty"`
	Title        string `json:"title"`
	MediaType    string `json:"mediaType,omitempty"`
	Comment      string `json:"comment,omitempty"`
	FileSize     int64  `json:"fileSize,omitempty"`
	DownloadLink string `json:"downloadLink,omitempty"`
	WebUILink    string `json:"webuiLink,omitempty"`
	Version      struct {
		CreatedAt string `json:"createdAt,omitempty"`
		Number    int    `json:"number,omitempty"`
	} `json:"version,omitempty"`
	Links struct {
		WebUI    string `json:"webui,omitempty"`
		Download string `json:"download,omitempty"`
	} `json:"_links,omitempty"`
}

func listAttachmentsCloudV2(ctx context.Context, c *client.Client, contentID string, limit int) ([]Attachment, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}

	path := "/api/v2/pages/" + url.PathEscape(contentID) + "/attachments"
	params := url.Values{"limit": {strconv.Itoa(limit)}}
	out := make([]Attachment, 0, limit)

	for len(out) < limit {
		data, hdrs, err := c.Get(ctx, path, params)
		if err != nil {
			return nil, err
		}
		var page struct {
			Results []cloudAttachmentV2 `json:"results"`
			Links   struct {
				Next string `json:"next,omitempty"`
				Base string `json:"base,omitempty"`
			} `json:"_links"`
		}
		if err := json.Unmarshal(data, &page); err != nil {
			return nil, fmt.Errorf("parse attachments: %w", err)
		}
		for _, raw := range page.Results {
			out = append(out, normalizeCloudAttachmentV2(raw, page.Links.Base))
			if len(out) == limit {
				break
			}
		}

		next := nextPageURL(hdrs, page.Links.Next)
		if next == "" || len(page.Results) == 0 || len(out) == limit {
			break
		}
		path, params, err = nextPageRequest(c, next)
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

func normalizeCloudAttachmentV2(in cloudAttachmentV2, base string) Attachment {
	var out Attachment
	out.ID = in.ID
	out.Title = in.Title
	out.Type = "attachment"
	out.Status = in.Status
	out.Version.Number = in.Version.Number
	out.Version.When = in.Version.CreatedAt
	out.Metadata.MediaType = in.MediaType
	out.Metadata.Comment = in.Comment
	out.Extensions.MediaType = in.MediaType
	out.Extensions.FileSize = in.FileSize
	out.Links.Base = base
	out.Links.WebUI = in.Links.WebUI
	out.Links.Download = in.Links.Download
	if out.Links.WebUI == "" {
		out.Links.WebUI = in.WebUILink
	}
	if out.Links.Download == "" {
		out.Links.Download = in.DownloadLink
	}
	return out
}
