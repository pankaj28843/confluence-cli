package conf

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

// ListLabels fetches /rest/api/content/{id}/label.
func ListLabels(ctx context.Context, c *client.Client, contentID string) ([]Label, error) {
	data, _, err := c.Get(ctx, "/rest/api/content/"+url.PathEscape(contentID)+"/label", nil)
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
