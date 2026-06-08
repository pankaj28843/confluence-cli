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

// WatcherPage is a paginated list of content or space watchers.
type WatcherPage struct {
	Results    []WatchRecord  `json:"results,omitempty"`
	TotalCount int64          `json:"totalCount,omitempty"`
	Start      int            `json:"start,omitempty"`
	Limit      int            `json:"limit,omitempty"`
	Size       int            `json:"size,omitempty"`
	Links      map[string]any `json:"_links,omitempty"`
}

// WatchRecord normalizes Cloud watch records and Server/DC user rows.
type WatchRecord struct {
	Type      string `json:"type,omitempty"`
	ContentID int64  `json:"contentId,omitempty"`
	SpaceKey  string `json:"spaceKey,omitempty"`
	LabelName string `json:"labelName,omitempty"`
	Prefix    string `json:"prefix,omitempty"`
	Watcher   User   `json:"watcher,omitempty"`
}

// WatchStatusOptions controls user-specific watch-status checks.
type WatchStatusOptions struct {
	AccountID   string
	UserKey     string
	Username    string
	ContentType string
}

// WatchStatus is the normalized result of a user watch-status check.
type WatchStatus struct {
	Watching bool `json:"watching"`
}

func (w *WatchRecord) UnmarshalJSON(data []byte) error {
	var wrapped struct {
		Type      string `json:"type"`
		ContentID int64  `json:"contentId"`
		SpaceKey  string `json:"spaceKey"`
		LabelName string `json:"labelName"`
		Prefix    string `json:"prefix"`
		Watcher   User   `json:"watcher"`
	}
	if err := json.Unmarshal(data, &wrapped); err != nil {
		return err
	}
	if watchUserLabel(wrapped.Watcher) != "" || wrapped.Type != "" || wrapped.ContentID != 0 || wrapped.SpaceKey != "" || wrapped.LabelName != "" {
		w.Type = wrapped.Type
		w.ContentID = wrapped.ContentID
		w.SpaceKey = wrapped.SpaceKey
		w.LabelName = wrapped.LabelName
		w.Prefix = wrapped.Prefix
		w.Watcher = wrapped.Watcher
		return nil
	}

	var user User
	if err := json.Unmarshal(data, &user); err != nil {
		return err
	}
	w.Watcher = user
	return nil
}

func watchUserLabel(u User) string {
	switch {
	case u.DisplayName != "":
		return u.DisplayName
	case u.PublicName != "":
		return u.PublicName
	case u.Username != "":
		return u.Username
	case u.AccountID != "":
		return u.AccountID
	case u.UserKey != "":
		return u.UserKey
	default:
		return ""
	}
}

func (s *WatchStatus) UnmarshalJSON(data []byte) error {
	var watching bool
	if err := json.Unmarshal(data, &watching); err == nil {
		s.Watching = watching
		return nil
	}
	var object struct {
		Watching bool `json:"watching"`
	}
	if err := json.Unmarshal(data, &object); err != nil {
		return err
	}
	s.Watching = object.Watching
	return nil
}

// GetContentWatchers returns raw bytes from the flavor-specific content watcher
// route. Prefer ListContentWatchers for typed command behavior.
func GetContentWatchers(ctx context.Context, c *client.Client, contentID string) ([]byte, error) {
	contentID = strings.TrimSpace(contentID)
	if contentID == "" {
		return nil, fmt.Errorf("GetContentWatchers: content id is required")
	}
	if c == nil {
		return nil, fmt.Errorf("GetContentWatchers: client is required")
	}
	data, _, err := c.Get(ctx, contentWatchersPath(c, contentID), nil)
	return data, err
}

// GetSpaceWatchers returns raw bytes from the flavor-specific space watcher
// route. Prefer ListSpaceWatchers for typed command behavior.
func GetSpaceWatchers(ctx context.Context, c *client.Client, spaceKey string) ([]byte, error) {
	spaceKey = strings.TrimSpace(spaceKey)
	if spaceKey == "" {
		return nil, fmt.Errorf("GetSpaceWatchers: space key is required")
	}
	if c == nil {
		return nil, fmt.Errorf("GetSpaceWatchers: client is required")
	}
	data, _, err := c.Get(ctx, spaceWatchersPath(c, spaceKey), nil)
	return data, err
}

// ListContentWatchers returns users watching content. Cloud uses the v1
// notification watch route; Server/DC uses the documented content watchers route.
func ListContentWatchers(ctx context.Context, c *client.Client, contentID string, limit int) (*WatcherPage, error) {
	contentID = strings.TrimSpace(contentID)
	if contentID == "" {
		return nil, fmt.Errorf("ListContentWatchers: content id is required")
	}
	if c == nil {
		return nil, fmt.Errorf("ListContentWatchers: client is required")
	}
	params := watchPageParams(limit)
	data, _, err := c.Get(ctx, contentWatchersPath(c, contentID), params)
	if err != nil {
		return nil, err
	}
	var out WatcherPage
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parse content watchers: %w", err)
	}
	return &out, nil
}

// ListSpaceWatchers returns users watching a space.
func ListSpaceWatchers(ctx context.Context, c *client.Client, spaceKey string, limit int) (*WatcherPage, error) {
	spaceKey = strings.TrimSpace(spaceKey)
	if spaceKey == "" {
		return nil, fmt.Errorf("ListSpaceWatchers: space key is required")
	}
	if c == nil {
		return nil, fmt.Errorf("ListSpaceWatchers: client is required")
	}
	params := watchPageParams(limit)
	data, _, err := c.Get(ctx, spaceWatchersPath(c, spaceKey), params)
	if err != nil {
		return nil, err
	}
	var out WatcherPage
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parse space watchers: %w", err)
	}
	return &out, nil
}

// GetContentWatchStatus returns whether a user watches a content id.
func GetContentWatchStatus(ctx context.Context, c *client.Client, contentID string, opts WatchStatusOptions) (*WatchStatus, error) {
	contentID = strings.TrimSpace(contentID)
	if contentID == "" {
		return nil, fmt.Errorf("GetContentWatchStatus: content id is required")
	}
	if c == nil {
		return nil, fmt.Errorf("GetContentWatchStatus: client is required")
	}
	data, _, err := c.Get(ctx, "/rest/api/user/watch/content/"+url.PathEscape(contentID), watchStatusParams(opts, false))
	if err != nil {
		return nil, err
	}
	var out WatchStatus
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parse content watch status: %w", err)
	}
	return &out, nil
}

// GetSpaceWatchStatus returns whether a user watches a space.
func GetSpaceWatchStatus(ctx context.Context, c *client.Client, spaceKey string, opts WatchStatusOptions) (*WatchStatus, error) {
	spaceKey = strings.TrimSpace(spaceKey)
	if spaceKey == "" {
		return nil, fmt.Errorf("GetSpaceWatchStatus: space key is required")
	}
	if c == nil {
		return nil, fmt.Errorf("GetSpaceWatchStatus: client is required")
	}
	data, _, err := c.Get(ctx, "/rest/api/user/watch/space/"+url.PathEscape(spaceKey), watchStatusParams(opts, true))
	if err != nil {
		return nil, err
	}
	var out WatchStatus
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parse space watch status: %w", err)
	}
	return &out, nil
}

func contentWatchersPath(c *client.Client, contentID string) string {
	if isCloud(c) {
		return "/rest/api/content/" + url.PathEscape(contentID) + "/notification/child-created"
	}
	return "/rest/api/content/" + url.PathEscape(contentID) + "/watchers"
}

func spaceWatchersPath(c *client.Client, spaceKey string) string {
	if isCloud(c) {
		return "/rest/api/space/" + url.PathEscape(spaceKey) + "/watch"
	}
	return "/rest/api/space/" + url.PathEscape(spaceKey) + "/watchers"
}

func watchPageParams(limit int) url.Values {
	limit = clampLimit(limit)
	return url.Values{
		"start": {"0"},
		"limit": {strconv.Itoa(limit)},
	}
}

func watchStatusParams(opts WatchStatusOptions, includeContentType bool) url.Values {
	params := url.Values{}
	if opts.UserKey != "" {
		params.Set("key", opts.UserKey)
	}
	if opts.Username != "" {
		params.Set("username", opts.Username)
	}
	if opts.AccountID != "" {
		params.Set("accountId", opts.AccountID)
	}
	if includeContentType && opts.ContentType != "" {
		params.Set("contentType", opts.ContentType)
	}
	return params
}

// GetContentRestrictions returns ACEs for a content id.
func GetContentRestrictions(ctx context.Context, c *client.Client, contentID string) ([]byte, error) {
	data, _, err := c.Get(ctx, "/rest/api/content/"+url.PathEscape(contentID)+"/restriction", url.Values{"expand": {"restrictions.user,restrictions.group"}})
	return data, err
}
