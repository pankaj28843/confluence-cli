package conf

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

// Space is a minimised /rest/api/space row.
type Space struct {
	ID          int64  `json:"id,omitempty"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Type        string `json:"type,omitempty"`   // global | personal
	Status      string `json:"status,omitempty"` // current | archived
	Description struct {
		Plain struct {
			Value string `json:"value,omitempty"`
		} `json:"plain,omitempty"`
	} `json:"description,omitempty"`
	Links struct {
		WebUI string `json:"webui,omitempty"`
	} `json:"_links,omitempty"`
}

// SpaceFilter narrows the list.
type SpaceFilter struct {
	Type   string // global | personal
	Status string // current | archived
	Limit  int
}

// ListSpaces pages /rest/api/space. Follows _links.next until limit is reached.
func ListSpaces(ctx context.Context, c *client.Client, f SpaceFilter) ([]Space, error) {
	if f.Limit <= 0 {
		f.Limit = 25
	}
	out := make([]Space, 0, f.Limit)
	start := 0
	for len(out) < f.Limit {
		params := url.Values{}
		remaining := f.Limit - len(out)
		if remaining > 200 {
			remaining = 200
		}
		params.Set("limit", strconv.Itoa(remaining))
		params.Set("start", strconv.Itoa(start))
		if f.Type != "" {
			params.Set("type", f.Type)
		}
		if f.Status != "" {
			params.Set("status", f.Status)
		}
		data, _, err := c.Get(ctx, "/rest/api/space", params)
		if err != nil {
			return nil, err
		}
		var page struct {
			Results []Space `json:"results"`
			Start   int     `json:"start"`
			Limit   int     `json:"limit"`
			Size    int     `json:"size"`
			Links   struct {
				Next string `json:"next,omitempty"`
			} `json:"_links"`
		}
		if err := json.Unmarshal(data, &page); err != nil {
			return nil, fmt.Errorf("parse spaces: %w", err)
		}
		out = append(out, page.Results...)
		if page.Links.Next == "" || len(page.Results) == 0 {
			break
		}
		start += page.Size
	}
	if len(out) > f.Limit {
		out = out[:f.Limit]
	}
	return out, nil
}

// GetSpace fetches /rest/api/space/{key}.
func GetSpace(ctx context.Context, c *client.Client, key string) (*Space, error) {
	data, _, err := c.Get(ctx, "/rest/api/space/"+url.PathEscape(key), url.Values{"expand": {"description.plain,homepage"}})
	if err != nil {
		return nil, err
	}
	var s Space
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parse space: %w", err)
	}
	return &s, nil
}
