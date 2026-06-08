package conf

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

// ListLabels fetches labels for a page/content id.
func ListLabels(ctx context.Context, c *client.Client, contentID string, limit int) ([]Label, error) {
	if isCloud(c) {
		return listLabelsCloudV2(ctx, c, contentID, limit)
	}
	return listLabelsServerV1(ctx, c, contentID, limit)
}

func listLabelsServerV1(ctx context.Context, c *client.Client, contentID string, limit int) ([]Label, error) {
	if limit <= 0 {
		limit = 25
	}
	if limit > 200 {
		limit = 200
	}
	data, _, err := c.Get(ctx, "/rest/api/content/"+url.PathEscape(contentID)+"/label", url.Values{"limit": {strconv.Itoa(limit)}})
	if err != nil {
		return nil, err
	}
	var page struct {
		Results []Label `json:"results"`
	}
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parse labels: %w", err)
	}
	return page.Results, nil
}

func listLabelsCloudV2(ctx context.Context, c *client.Client, contentID string, limit int) ([]Label, error) {
	if limit <= 0 {
		limit = 25
	}
	if limit > 200 {
		limit = 200
	}

	path := "/api/v2/pages/" + url.PathEscape(contentID) + "/labels"
	params := url.Values{"limit": {strconv.Itoa(limit)}}
	out := make([]Label, 0, limit)

	for len(out) < limit {
		data, hdrs, err := c.Get(ctx, path, params)
		if err != nil {
			return nil, err
		}
		var page struct {
			Results []Label `json:"results"`
			Links   struct {
				Next string `json:"next,omitempty"`
			} `json:"_links"`
		}
		if err := json.Unmarshal(data, &page); err != nil {
			return nil, fmt.Errorf("parse labels: %w", err)
		}
		for _, label := range page.Results {
			out = append(out, label)
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

// AddLabels POSTs one or more labels to a content id.
func AddLabels(ctx context.Context, c *client.Client, contentID string, names []string) ([]Label, error) {
	body := make([]map[string]string, 0, len(names))
	for _, n := range names {
		body = append(body, map[string]string{"prefix": "global", "name": n})
	}
	data, _, err := c.Post(ctx, "/rest/api/content/"+url.PathEscape(contentID)+"/label", nil, body)
	if err != nil {
		return nil, err
	}
	var page struct {
		Results []Label `json:"results"`
	}
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parse add-labels response: %w", err)
	}
	return page.Results, nil
}

// RemoveLabel DELETEs one label by name.
func RemoveLabel(ctx context.Context, c *client.Client, contentID, name string) error {
	params := url.Values{"name": {name}}
	_, _, err := c.Delete(ctx, "/rest/api/content/"+url.PathEscape(contentID)+"/label", params)
	return err
}
