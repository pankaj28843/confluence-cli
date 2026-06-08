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

// LabelTarget identifies the entity whose labels should be read.
type LabelTarget struct {
	Type string
	ID   string
}

// LabelListOptions controls label list requests.
type LabelListOptions struct {
	Limit  int
	Prefix string
	Sort   string
	Scope  string // content or space; only used by ListSpaceLabels on Cloud
}

// LabelSearchOptions controls Cloud v2 label catalog reads.
type LabelSearchOptions struct {
	Limit    int
	LabelIDs []string
	Prefixes []string
	Sort     string
}

// LabelRelatedOptions controls Server/Data Center related-label reads.
type LabelRelatedOptions struct {
	SpaceKey string
	Label    string
	Limit    int
}

// ListLabels fetches labels for a page/content id.
func ListLabels(ctx context.Context, c *client.Client, contentID string, limit int) ([]Label, error) {
	return ListTargetLabels(ctx, c, LabelTarget{Type: "page", ID: contentID}, LabelListOptions{Limit: limit})
}

// ListTargetLabels fetches labels for a supported entity target.
func ListTargetLabels(ctx context.Context, c *client.Client, target LabelTarget, opts LabelListOptions) ([]Label, error) {
	if target.ID == "" {
		return nil, fmt.Errorf("ListTargetLabels: target id is required")
	}
	target.Type = strings.TrimSpace(strings.ToLower(target.Type))
	if target.Type == "" {
		target.Type = "page"
	}

	if isCloud(c) {
		path, ok := cloudLabelPath(target)
		if !ok {
			return nil, fmt.Errorf("ListTargetLabels: unsupported Cloud label target %q", target.Type)
		}
		return listLabelsCloudV2(ctx, c, path, opts)
	}
	return listLabelsServerV1(ctx, c, "/rest/api/content/"+url.PathEscape(target.ID)+"/label", opts)
}

func cloudLabelPath(target LabelTarget) (string, bool) {
	id := url.PathEscape(target.ID)
	switch target.Type {
	case "page":
		return "/api/v2/pages/" + id + "/labels", true
	case "blogpost", "blog-post":
		return "/api/v2/blogposts/" + id + "/labels", true
	case "attachment":
		return "/api/v2/attachments/" + id + "/labels", true
	case "custom-content":
		return "/api/v2/custom-content/" + id + "/labels", true
	default:
		return "", false
	}
}

func listLabelsServerV1(ctx context.Context, c *client.Client, path string, opts LabelListOptions) ([]Label, error) {
	limit := clampLimit(opts.Limit)
	params := labelListParams(opts, limit)
	params.Set("start", "0")

	data, _, err := c.Get(ctx, path, params)
	if err != nil {
		return nil, err
	}
	labels, err := parseLabelResults(data)
	if err != nil {
		return nil, err
	}
	if len(labels) > limit {
		labels = labels[:limit]
	}
	return labels, nil
}

func listLabelsCloudV2(ctx context.Context, c *client.Client, path string, opts LabelListOptions) ([]Label, error) {
	limit := clampLimit(opts.Limit)
	params := labelListParams(opts, limit)
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

func labelListParams(opts LabelListOptions, limit int) url.Values {
	params := url.Values{"limit": {strconv.Itoa(limit)}}
	if opts.Prefix != "" {
		params.Set("prefix", opts.Prefix)
	}
	if opts.Sort != "" {
		params.Set("sort", opts.Sort)
	}
	return params
}

// ListSpaceLabels fetches labels used by content in a space. On Cloud, Scope
// may be "space" to read labels attached to the space entity itself.
func ListSpaceLabels(ctx context.Context, c *client.Client, spaceKey string, opts LabelListOptions) ([]Label, error) {
	if spaceKey == "" {
		return nil, fmt.Errorf("ListSpaceLabels: --space is required")
	}
	scope := strings.TrimSpace(strings.ToLower(opts.Scope))
	if scope == "" {
		scope = "content"
	}

	if isCloud(c) {
		spaceID, err := resolveCloudSpaceID(ctx, c, spaceKey)
		if err != nil {
			return nil, err
		}
		path := "/api/v2/spaces/" + url.PathEscape(spaceID) + "/content/labels"
		if scope == "space" {
			path = "/api/v2/spaces/" + url.PathEscape(spaceID) + "/labels"
		}
		if scope != "content" && scope != "space" {
			return nil, fmt.Errorf("ListSpaceLabels: unsupported scope %q", opts.Scope)
		}
		return listLabelsCloudV2(ctx, c, path, opts)
	}
	if scope != "content" {
		return nil, fmt.Errorf("ListSpaceLabels: scope %q is Confluence Cloud only", opts.Scope)
	}
	return listLabelsServerV1(ctx, c, "/rest/api/space/"+url.PathEscape(spaceKey)+"/labels", opts)
}

func resolveCloudSpaceID(ctx context.Context, c *client.Client, keyOrID string) (string, error) {
	if _, err := strconv.ParseInt(keyOrID, 10, 64); err == nil {
		return keyOrID, nil
	}
	space, err := getCloudSpaceV2(ctx, c, keyOrID)
	if err != nil {
		return "", err
	}
	if space.ID == "" {
		return "", fmt.Errorf("space %q has no Cloud id", keyOrID)
	}
	return space.ID, nil
}

// SearchLabels lists Cloud v2 labels with optional id and prefix filters.
func SearchLabels(ctx context.Context, c *client.Client, opts LabelSearchOptions) ([]Label, error) {
	if !isCloud(c) {
		return nil, fmt.Errorf("SearchLabels: Confluence Cloud only")
	}
	limit := clampLimit(opts.Limit)
	params := url.Values{"limit": {strconv.Itoa(limit)}}
	for _, id := range opts.LabelIDs {
		if strings.TrimSpace(id) != "" {
			params.Add("label-id", strings.TrimSpace(id))
		}
	}
	for _, prefix := range opts.Prefixes {
		if strings.TrimSpace(prefix) != "" {
			params.Add("prefix", strings.TrimSpace(prefix))
		}
	}
	if opts.Sort != "" {
		params.Set("sort", opts.Sort)
	}
	return listLabelsCloudV2WithParams(ctx, c, "/api/v2/labels", params, limit)
}

func listLabelsCloudV2WithParams(ctx context.Context, c *client.Client, path string, params url.Values, limit int) ([]Label, error) {
	out := make([]Label, 0, limit)
	for len(out) < limit {
		data, hdrs, err := c.Get(ctx, path, params)
		if err != nil {
			return nil, err
		}
		labels, next, err := parseCloudLabelPage(data)
		if err != nil {
			return nil, err
		}
		for _, label := range labels {
			out = append(out, label)
			if len(out) == limit {
				break
			}
		}
		next = nextPageURL(hdrs, next)
		if next == "" || len(labels) == 0 || len(out) == limit {
			break
		}
		path, params, err = nextPageRequest(c, next)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

// ListRecentLabels fetches recently used Server/Data Center labels.
func ListRecentLabels(ctx context.Context, c *client.Client, limit int) ([]Label, error) {
	if isCloud(c) {
		return nil, fmt.Errorf("ListRecentLabels: Confluence Server/Data Center only")
	}
	return listLabelsServerV1(ctx, c, "/rest/api/label/recent", LabelListOptions{Limit: limit})
}

// ListRelatedLabels fetches Server/Data Center labels related to another label.
func ListRelatedLabels(ctx context.Context, c *client.Client, opts LabelRelatedOptions) ([]Label, error) {
	if isCloud(c) {
		return nil, fmt.Errorf("ListRelatedLabels: Confluence Server/Data Center only")
	}
	if opts.Label == "" {
		return nil, fmt.Errorf("ListRelatedLabels: --label is required")
	}
	path := "/rest/api/label/" + url.PathEscape(opts.Label) + "/related"
	if opts.SpaceKey != "" {
		path = "/rest/api/space/" + url.PathEscape(opts.SpaceKey) + "/labels/" + url.PathEscape(opts.Label) + "/related"
	}
	return listLabelsServerV1(ctx, c, path, LabelListOptions{Limit: opts.Limit})
}

func parseLabelResults(data []byte) ([]Label, error) {
	var page struct {
		Results []Label `json:"results"`
	}
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parse labels: %w", err)
	}
	if page.Results != nil {
		return page.Results, nil
	}
	var single Label
	if err := json.Unmarshal(data, &single); err != nil {
		return nil, fmt.Errorf("parse labels: %w", err)
	}
	if single.Name == "" && single.ID == "" && single.Label == "" {
		return nil, nil
	}
	return []Label{single}, nil
}

func parseCloudLabelPage(data []byte) ([]Label, string, error) {
	var page struct {
		Results []Label `json:"results"`
		Links   struct {
			Next string `json:"next,omitempty"`
		} `json:"_links"`
	}
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, "", fmt.Errorf("parse labels: %w", err)
	}
	return page.Results, page.Links.Next, nil
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
