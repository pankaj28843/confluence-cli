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

// Content is the minimised /rest/api/content{/:id,/search} shape. Handles pages,
// blogposts, comments, and attachments.
type Content struct {
	ID     string `json:"id"`
	Type   string `json:"type,omitempty"` // page | blogpost | comment | attachment
	Status string `json:"status,omitempty"`
	Title  string `json:"title,omitempty"`
	Space  struct {
		Key  string `json:"key,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"space,omitempty"`
	Version struct {
		Number int    `json:"number,omitempty"`
		When   string `json:"when,omitempty"`
		By     struct {
			DisplayName string `json:"displayName,omitempty"`
			Email       string `json:"email,omitempty"`
			AccountID   string `json:"accountId,omitempty"`
			Username    string `json:"username,omitempty"`
		} `json:"by,omitempty"`
	} `json:"version,omitempty"`
	Body struct {
		Storage struct {
			Value          string `json:"value,omitempty"`
			Representation string `json:"representation,omitempty"`
		} `json:"storage,omitempty"`
		View struct {
			Value string `json:"value,omitempty"`
		} `json:"view,omitempty"`
	} `json:"body,omitempty"`
	Metadata struct {
		Labels struct {
			Results []Label `json:"results,omitempty"`
		} `json:"labels,omitempty"`
	} `json:"metadata,omitempty"`
	Ancestors []Content `json:"ancestors,omitempty"`
	Links     struct {
		WebUI string `json:"webui,omitempty"`
		Base  string `json:"base,omitempty"`
		Self  string `json:"self,omitempty"`
	} `json:"_links,omitempty"`
	Excerpt string `json:"excerpt,omitempty"`
}

// Label is a content label.
type Label struct {
	Prefix string `json:"prefix,omitempty"`
	Name   string `json:"name"`
	ID     string `json:"id,omitempty"`
	Label  string `json:"label,omitempty"`
}

// DefaultExpand is the set of expansions most commands use.
const DefaultExpand = "body.storage,version,space,ancestors,metadata.labels"

// GetContent fetches /rest/api/content/{id}. expand defaults to DefaultExpand.
func GetContent(ctx context.Context, c *client.Client, id, expand string) (*Content, error) {
	if expand == "" {
		expand = DefaultExpand
	}
	data, _, err := c.Get(ctx, "/rest/api/content/"+url.PathEscape(id), url.Values{"expand": {expand}})
	if err != nil {
		return nil, err
	}
	var out Content
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parse content: %w", err)
	}
	return &out, nil
}

// SearchCQL runs an advanced search via /rest/api/content/search (typed
// results) or /rest/api/search (generic). Using /content/search here for
// page-focused queries.
func SearchCQL(ctx context.Context, c *client.Client, cql string, limit int, expand string) ([]Content, error) {
	if limit <= 0 {
		limit = 25
	}
	if limit > 200 {
		limit = 200
	}
	if expand == "" {
		expand = "version,space,metadata.labels"
	}
	params := url.Values{
		"cql":    {cql},
		"limit":  {strconv.Itoa(limit)},
		"expand": {expand},
	}
	data, _, err := c.Get(ctx, "/rest/api/content/search", params)
	if err != nil {
		return nil, err
	}
	var page struct {
		Results []Content `json:"results"`
	}
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parse content search: %w", err)
	}
	return page.Results, nil
}

// SearchGeneric runs /rest/api/search?cql=..., which returns mixed-type rows
// (pages, blogposts, users, spaces, attachments) with `_searchResult` excerpts.
func SearchGeneric(ctx context.Context, c *client.Client, cql string, limit int) ([]SearchHit, error) {
	if limit <= 0 {
		limit = 25
	}
	if limit > 200 {
		limit = 200
	}
	params := url.Values{"cql": {cql}, "limit": {strconv.Itoa(limit)}}
	data, _, err := c.Get(ctx, "/rest/api/search", params)
	if err != nil {
		return nil, err
	}
	var page struct {
		Results []SearchHit `json:"results"`
	}
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parse search: %w", err)
	}
	return page.Results, nil
}

// SearchHit is one row from /rest/api/search — mixed entity types.
type SearchHit struct {
	Content      *Content `json:"content,omitempty"`
	Title        string   `json:"title,omitempty"`
	Excerpt      string   `json:"excerpt,omitempty"`
	URL          string   `json:"url,omitempty"`
	EntityType   string   `json:"entityType,omitempty"` // content | space | user
	Space        *Space   `json:"space,omitempty"`
	User         *User    `json:"user,omitempty"`
	LastModified string   `json:"lastModified,omitempty"`
}

// User is a minimised user record used by search hits + user endpoints.
type User struct {
	AccountID   string `json:"accountId,omitempty"`
	UserKey     string `json:"userKey,omitempty"`
	Username    string `json:"username,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	PublicName  string `json:"publicName,omitempty"`
	Email       string `json:"email,omitempty"`
	Type        string `json:"type,omitempty"`
}

// GetChildren returns children of a content id (default type: page).
func GetChildren(ctx context.Context, c *client.Client, id, childType string, limit int) ([]Content, error) {
	if childType == "" {
		childType = "page"
	}
	if limit <= 0 {
		limit = 50
	}
	params := url.Values{
		"limit":  {strconv.Itoa(limit)},
		"expand": {"version,space"},
	}
	path := "/rest/api/content/" + url.PathEscape(id) + "/child/" + url.PathEscape(childType)
	data, _, err := c.Get(ctx, path, params)
	if err != nil {
		return nil, err
	}
	var page struct {
		Results []Content `json:"results"`
	}
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parse children: %w", err)
	}
	return page.Results, nil
}

// GetAncestors walks parents of a content id.
func GetAncestors(ctx context.Context, c *client.Client, id string) ([]Content, error) {
	// `ancestors` is an expand on the content itself.
	cnt, err := GetContent(ctx, c, id, "ancestors")
	if err != nil {
		return nil, err
	}
	return cnt.Ancestors, nil
}

// GetHistory fetches /rest/api/content/{id}/history.
func GetHistory(ctx context.Context, c *client.Client, id string) ([]byte, error) {
	data, _, err := c.Get(ctx, "/rest/api/content/"+url.PathEscape(id)+"/history", nil)
	return data, err
}

// ListVersions fetches /rest/api/content/{id}/version.
func ListVersions(ctx context.Context, c *client.Client, id string) ([]byte, error) {
	data, _, err := c.Get(ctx, "/rest/api/content/"+url.PathEscape(id)+"/version", nil)
	return data, err
}

// RenderMarkdown returns the page body converted to Markdown, or the raw
// storage value if storage rendering is requested.
func (p *Content) RenderMarkdown(rawStorage bool) string {
	if rawStorage {
		return p.Body.Storage.Value
	}
	body := p.Body.Storage.Value
	if body == "" {
		body = p.Body.View.Value
	}
	return HTMLToMarkdown(body)
}

// AbsoluteURL combines _links.base + _links.webui. Useful in JSON output.
func (p *Content) AbsoluteURL() string {
	if p.Links.Base == "" {
		return p.Links.WebUI
	}
	return strings.TrimRight(p.Links.Base, "/") + p.Links.WebUI
}
