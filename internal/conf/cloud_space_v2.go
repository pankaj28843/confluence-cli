package conf

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

type cloudSpaceV2 struct {
	ID          string `json:"id"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Type        string `json:"type,omitempty"`
	Status      string `json:"status,omitempty"`
	Description *struct {
		Plain *struct {
			Value string `json:"value,omitempty"`
		} `json:"plain,omitempty"`
	} `json:"description,omitempty"`
	Links struct {
		WebUI string `json:"webui,omitempty"`
	} `json:"_links,omitempty"`
}

func listSpacesCloudV2(ctx context.Context, c *client.Client, f SpaceFilter) ([]Space, error) {
	if f.Limit <= 0 {
		f.Limit = 25
	}

	out := make([]Space, 0, f.Limit)
	path := "/api/v2/spaces"
	params := url.Values{}

	remaining := f.Limit
	if remaining > 200 {
		remaining = 200
	}
	params.Set("limit", strconv.Itoa(remaining))
	if f.Type != "" {
		params.Set("type", f.Type)
	}
	if f.Status != "" {
		params.Set("status", f.Status)
	}

	for len(out) < f.Limit {
		data, _, err := c.Get(ctx, path, params)
		if err != nil {
			return nil, err
		}

		var page struct {
			Results []cloudSpaceV2 `json:"results"`
			Links   struct {
				Next string `json:"next,omitempty"`
			} `json:"_links"`
		}
		if err := json.Unmarshal(data, &page); err != nil {
			return nil, fmt.Errorf("parse spaces: %w", err)
		}

		for _, raw := range page.Results {
			out = append(out, normalizeCloudSpaceV2(raw))
			if len(out) == f.Limit {
				break
			}
		}
		if page.Links.Next == "" || len(page.Results) == 0 || len(out) == f.Limit {
			break
		}

		path, params, err = nextPageRequest(c, page.Links.Next)
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

func getSpaceCloudV2(ctx context.Context, c *client.Client, key string) (*Space, error) {
	raw, err := getCloudSpaceV2(ctx, c, key)
	if err != nil {
		return nil, err
	}
	s := normalizeCloudSpaceV2(*raw)
	return &s, nil
}

func getCloudSpaceV2(ctx context.Context, c *client.Client, key string) (*cloudSpaceV2, error) {
	params := url.Values{
		"keys":  {key},
		"limit": {"1"},
	}
	data, _, err := c.Get(ctx, "/api/v2/spaces", params)
	if err != nil {
		return nil, err
	}

	var page struct {
		Results []cloudSpaceV2 `json:"results"`
	}
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parse space: %w", err)
	}
	if len(page.Results) == 0 {
		return nil, fmt.Errorf("space %q not found", key)
	}

	for _, raw := range page.Results {
		if strings.EqualFold(raw.Key, key) {
			return &raw, nil
		}
	}

	return &page.Results[0], nil
}

func normalizeCloudSpaceV2(in cloudSpaceV2) Space {
	var out Space
	if in.ID != "" {
		if id, err := strconv.ParseInt(in.ID, 10, 64); err == nil {
			out.ID = id
		}
	}
	out.Key = in.Key
	out.Name = in.Name
	out.Type = in.Type
	out.Status = in.Status
	out.Links.WebUI = in.Links.WebUI
	if in.Description != nil && in.Description.Plain != nil {
		out.Description.Plain.Value = in.Description.Plain.Value
	}
	return out
}
