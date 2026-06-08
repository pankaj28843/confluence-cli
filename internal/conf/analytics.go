package conf

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

// AnalyticsCount is the documented Cloud analytics count response.
type AnalyticsCount struct {
	ID    int64 `json:"id"`
	Count int64 `json:"count"`
}

// AnalyticsOptions controls Cloud analytics reads.
type AnalyticsOptions struct {
	FromDate string
}

// GetContentViewCount returns the total Cloud view count for one content item.
func GetContentViewCount(ctx context.Context, c *client.Client, contentID string, opts AnalyticsOptions) (*AnalyticsCount, error) {
	return getContentAnalyticsCount(ctx, c, contentID, "views", opts)
}

// GetContentViewerCount returns the total Cloud distinct viewer count for one content item.
func GetContentViewerCount(ctx context.Context, c *client.Client, contentID string, opts AnalyticsOptions) (*AnalyticsCount, error) {
	return getContentAnalyticsCount(ctx, c, contentID, "viewers", opts)
}

func getContentAnalyticsCount(ctx context.Context, c *client.Client, contentID, kind string, opts AnalyticsOptions) (*AnalyticsCount, error) {
	if !isCloud(c) {
		return nil, fmt.Errorf("analytics are only supported on Confluence Cloud")
	}
	contentID = strings.TrimSpace(contentID)
	if contentID == "" {
		return nil, fmt.Errorf("GetContentAnalyticsCount: content id is required")
	}
	params := url.Values{}
	if fromDate := strings.TrimSpace(opts.FromDate); fromDate != "" {
		params.Set("fromDate", fromDate)
	}
	data, _, err := c.Get(ctx, "/rest/api/analytics/content/"+url.PathEscape(contentID)+"/"+kind, params)
	if err != nil {
		return nil, err
	}
	var out AnalyticsCount
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parse analytics count: %w", err)
	}
	return &out, nil
}
