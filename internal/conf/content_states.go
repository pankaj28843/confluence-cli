package conf

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

const contentStateContentLimit = 100

// ContentState is the documented Cloud content state shape.
type ContentState struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

// ContentStateResponse is the current state attached to a content item.
type ContentStateResponse struct {
	State       *ContentState `json:"contentState"`
	LastUpdated string        `json:"lastUpdated,omitempty"`
}

// AvailableContentStates groups states that can be applied to one content item.
type AvailableContentStates struct {
	SpaceContentStates  []ContentState `json:"spaceContentStates"`
	CustomContentStates []ContentState `json:"customContentStates"`
}

// ContentStateSettings describes space-level content-state configuration.
type ContentStateSettings struct {
	ContentStatesAllowed       bool           `json:"contentStatesAllowed"`
	CustomContentStatesAllowed bool           `json:"customContentStatesAllowed"`
	SpaceContentStatesAllowed  bool           `json:"spaceContentStatesAllowed"`
	SpaceContentStates         []ContentState `json:"spaceContentStates,omitempty"`
}

// ContentWithStateOptions controls content-state content listing.
type ContentWithStateOptions struct {
	SpaceKey string
	StateID  *int64
	Expand   []string
	Start    int
	Limit    int
}

// GetContentState returns the current Cloud content state for content.
func GetContentState(ctx context.Context, c *client.Client, id, status string) (*ContentStateResponse, error) {
	if !isCloud(c) {
		return nil, fmt.Errorf("content states are only supported on Confluence Cloud")
	}
	if id == "" {
		return nil, fmt.Errorf("GetContentState: content id is required")
	}
	if err := validateContentStateStatus(status); err != nil {
		return nil, err
	}
	params := url.Values{}
	if status != "" {
		params.Set("status", status)
	}
	data, _, err := c.Get(ctx, "/rest/api/content/"+url.PathEscape(id)+"/state", params)
	if err != nil {
		return nil, err
	}
	var out ContentStateResponse
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parse content state: %w", err)
	}
	return &out, nil
}

// ListAvailableContentStates returns states available for one Cloud content item.
func ListAvailableContentStates(ctx context.Context, c *client.Client, id string) (*AvailableContentStates, error) {
	if !isCloud(c) {
		return nil, fmt.Errorf("content states are only supported on Confluence Cloud")
	}
	if id == "" {
		return nil, fmt.Errorf("ListAvailableContentStates: content id is required")
	}
	data, _, err := c.Get(ctx, "/rest/api/content/"+url.PathEscape(id)+"/state/available", nil)
	if err != nil {
		return nil, err
	}
	var out AvailableContentStates
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parse available content states: %w", err)
	}
	return &out, nil
}

// ListCustomContentStates returns custom states created by the Cloud user.
func ListCustomContentStates(ctx context.Context, c *client.Client) ([]ContentState, error) {
	if !isCloud(c) {
		return nil, fmt.Errorf("content states are only supported on Confluence Cloud")
	}
	data, _, err := c.Get(ctx, "/rest/api/content-states", nil)
	if err != nil {
		return nil, err
	}
	var out []ContentState
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parse custom content states: %w", err)
	}
	return out, nil
}

// ListSpaceContentStates returns suggested states for a Cloud space.
func ListSpaceContentStates(ctx context.Context, c *client.Client, spaceKey string) ([]ContentState, error) {
	if !isCloud(c) {
		return nil, fmt.Errorf("content states are only supported on Confluence Cloud")
	}
	if spaceKey == "" {
		return nil, fmt.Errorf("ListSpaceContentStates: space key is required")
	}
	data, _, err := c.Get(ctx, "/rest/api/space/"+url.PathEscape(spaceKey)+"/state", nil)
	if err != nil {
		return nil, err
	}
	var out []ContentState
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parse space content states: %w", err)
	}
	return out, nil
}

// GetContentStateSettings returns space content-state settings for Cloud.
func GetContentStateSettings(ctx context.Context, c *client.Client, spaceKey string) (*ContentStateSettings, error) {
	if !isCloud(c) {
		return nil, fmt.Errorf("content states are only supported on Confluence Cloud")
	}
	if spaceKey == "" {
		return nil, fmt.Errorf("GetContentStateSettings: space key is required")
	}
	data, _, err := c.Get(ctx, "/rest/api/space/"+url.PathEscape(spaceKey)+"/state/settings", nil)
	if err != nil {
		return nil, err
	}
	var out ContentStateSettings
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parse content state settings: %w", err)
	}
	return &out, nil
}

// ListContentWithState returns content in a Cloud space with a given state.
func ListContentWithState(ctx context.Context, c *client.Client, opts ContentWithStateOptions) ([]Content, error) {
	if !isCloud(c) {
		return nil, fmt.Errorf("content states are only supported on Confluence Cloud")
	}
	if opts.SpaceKey == "" {
		return nil, fmt.Errorf("ListContentWithState: space key is required")
	}
	if opts.StateID == nil {
		return nil, fmt.Errorf("ListContentWithState: state id is required")
	}
	limit := clampContentStateContentLimit(opts.Limit)
	start := opts.Start
	if start < 0 {
		start = 0
	}
	path := "/rest/api/space/" + url.PathEscape(opts.SpaceKey) + "/state/content"
	params := url.Values{
		"state-id": {strconv.FormatInt(*opts.StateID, 10)},
		"start":    {strconv.Itoa(start)},
		"limit":    {strconv.Itoa(limit)},
	}
	addValues(params, "expand", opts.Expand)
	return collectContentWithState(ctx, c, path, params, limit)
}

func collectContentWithState(ctx context.Context, c *client.Client, path string, params url.Values, limit int) ([]Content, error) {
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
			return nil, fmt.Errorf("parse content with state: %w", err)
		}
		for _, content := range page.Results {
			out = append(out, content)
			if len(out) == limit {
				break
			}
		}
		next := nextPageURL(headers, page.Links.Next)
		if next == "" || len(page.Results) == 0 || len(out) == limit {
			break
		}
		nextPath, nextParams, err := nextPageRequest(c, next)
		if err != nil {
			return nil, err
		}
		path, params = nextPath, nextParams
	}
	return out, nil
}

func validateContentStateStatus(status string) error {
	switch status {
	case "", "current", "draft", "archived":
		return nil
	default:
		return fmt.Errorf("content state status must be current, draft, or archived")
	}
}

func clampContentStateContentLimit(limit int) int {
	if limit <= 0 {
		return 25
	}
	if limit > contentStateContentLimit {
		return contentStateContentLimit
	}
	return limit
}
