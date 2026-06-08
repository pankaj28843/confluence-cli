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

// CustomContentListOptions controls Cloud v2 custom-content listing.
type CustomContentListOptions struct {
	Type          string
	IDs           []string
	SpaceIDs      []string
	ContainerType string
	ContainerID   string
	Limit         int
	Sort          string
	BodyFormat    string
}

// CustomContentGetOptions controls Cloud v2 custom-content detail reads.
type CustomContentGetOptions struct {
	BodyFormat           string
	Version              int
	IncludeLabels        bool
	IncludeProperties    bool
	IncludeOperations    bool
	IncludeVersions      bool
	IncludeVersion       bool
	IncludeCollaborators bool
}

// ListCustomContent returns Cloud v2 custom content by type, either globally
// or inside a documented page, blogpost, or space container.
func ListCustomContent(ctx context.Context, c *client.Client, opts CustomContentListOptions) ([]Content, error) {
	if !isCloud(c) {
		return nil, fmt.Errorf("ListCustomContent: Confluence Cloud only")
	}
	opts.Type = strings.TrimSpace(opts.Type)
	if opts.Type == "" {
		return nil, fmt.Errorf("ListCustomContent: type is required")
	}
	opts.ContainerType = normalizeEntityTarget(opts.ContainerType)
	opts.ContainerID = strings.TrimSpace(opts.ContainerID)
	if opts.ContainerType != "" && opts.ContainerID == "" {
		return nil, fmt.Errorf("ListCustomContent: container id is required")
	}
	if opts.ContainerType == "" && opts.ContainerID != "" {
		return nil, fmt.Errorf("ListCustomContent: container type is required")
	}
	if opts.ContainerType != "" && (len(opts.IDs) > 0 || len(opts.SpaceIDs) > 0) {
		return nil, fmt.Errorf("ListCustomContent: id and space-id filters are only supported by the global endpoint")
	}
	if opts.ContainerType == "space" && opts.Sort != "" {
		return nil, fmt.Errorf("ListCustomContent: sort is not documented for space-scoped custom content")
	}

	opts.Limit = clampLimit(opts.Limit)
	path, err := customContentListPath(opts.ContainerType, opts.ContainerID)
	if err != nil {
		return nil, err
	}
	params := customContentListParams(opts)
	return collectCustomContentPages(ctx, c, path, params, opts.Limit)
}

// GetCustomContent returns one Cloud v2 custom-content record by id.
func GetCustomContent(ctx context.Context, c *client.Client, id string, opts CustomContentGetOptions) (*Content, error) {
	if !isCloud(c) {
		return nil, fmt.Errorf("GetCustomContent: Confluence Cloud only")
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("GetCustomContent: id is required")
	}
	if opts.Version < 0 {
		return nil, fmt.Errorf("GetCustomContent: version must be positive")
	}
	data, _, err := c.Get(ctx, "/api/v2/custom-content/"+url.PathEscape(id), customContentGetParams(opts))
	if err != nil {
		return nil, err
	}
	var out Content
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parse custom content: %w", err)
	}
	return &out, nil
}

// ListCustomContentChildren returns Cloud v2 child custom-content records for
// one custom-content id.
func ListCustomContentChildren(ctx context.Context, c *client.Client, id string, opts DirectChildrenOptions) ([]Content, error) {
	if !isCloud(c) {
		return nil, fmt.Errorf("ListCustomContentChildren: Confluence Cloud only")
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("ListCustomContentChildren: id is required")
	}
	opts.Limit = clampLimit(opts.Limit)
	opts.Types = normalizeContentTypes(opts.Types)
	path := "/api/v2/custom-content/" + url.PathEscape(id) + "/children"
	return listCloudDirectChildren(ctx, c, path, opts)
}

func customContentListPath(containerType, containerID string) (string, error) {
	if containerType == "" {
		return "/api/v2/custom-content", nil
	}
	switch normalizeEntityTarget(containerType) {
	case "page":
		return "/api/v2/pages/" + url.PathEscape(containerID) + "/custom-content", nil
	case "blogpost":
		return "/api/v2/blogposts/" + url.PathEscape(containerID) + "/custom-content", nil
	case "space":
		return "/api/v2/spaces/" + url.PathEscape(containerID) + "/custom-content", nil
	default:
		return "", fmt.Errorf("unsupported custom-content container %q", containerType)
	}
}

func customContentListParams(opts CustomContentListOptions) url.Values {
	params := url.Values{
		"type":  {opts.Type},
		"limit": {strconv.Itoa(opts.Limit)},
	}
	if opts.Sort != "" {
		params.Set("sort", opts.Sort)
	}
	if opts.BodyFormat != "" {
		params.Set("body-format", opts.BodyFormat)
	}
	for _, id := range normalizedValues(opts.IDs) {
		params.Add("id", id)
	}
	for _, id := range normalizedValues(opts.SpaceIDs) {
		params.Add("space-id", id)
	}
	return params
}

func customContentGetParams(opts CustomContentGetOptions) url.Values {
	params := url.Values{}
	if opts.BodyFormat != "" {
		params.Set("body-format", opts.BodyFormat)
	}
	if opts.Version > 0 {
		params.Set("version", strconv.Itoa(opts.Version))
	}
	setBoolParam(params, "include-labels", opts.IncludeLabels)
	setBoolParam(params, "include-properties", opts.IncludeProperties)
	setBoolParam(params, "include-operations", opts.IncludeOperations)
	setBoolParam(params, "include-versions", opts.IncludeVersions)
	setBoolParam(params, "include-version", opts.IncludeVersion)
	setBoolParam(params, "include-collaborators", opts.IncludeCollaborators)
	return params
}

func collectCustomContentPages(ctx context.Context, c *client.Client, path string, params url.Values, limit int) ([]Content, error) {
	out := make([]Content, 0, limit)
	for len(out) < limit {
		data, headers, err := c.Get(ctx, path, params)
		if err != nil {
			return nil, err
		}
		var page struct {
			Results []Content `json:"results"`
			Links   struct {
				Next string `json:"next,omitempty"`
			} `json:"_links"`
		}
		if err := json.Unmarshal(data, &page); err != nil {
			return nil, fmt.Errorf("parse custom content: %w", err)
		}
		for _, item := range page.Results {
			out = append(out, item)
			if len(out) == limit {
				break
			}
		}
		next := nextPageURL(headers, page.Links.Next)
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

func normalizedValues(values []string) []string {
	out := make([]string, 0, len(values))
	for _, raw := range values {
		for _, part := range strings.Split(raw, ",") {
			value := strings.TrimSpace(part)
			if value != "" {
				out = append(out, value)
			}
		}
	}
	return out
}
