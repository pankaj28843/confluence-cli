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

// ContentTreeEntityOptions controls Cloud v2 content-tree entity reads.
type ContentTreeEntityOptions struct {
	IncludeCollaborators  bool
	IncludeDirectChildren bool
	IncludeOperations     bool
	IncludeProperties     bool
}

// ContentTreeEntity is the shared Cloud v2 shape for modern content-tree
// entities: databases, folders, whiteboards, and Smart Link embeds.
type ContentTreeEntity struct {
	ID         string      `json:"id,omitempty"`
	Type       string      `json:"type,omitempty"`
	Status     string      `json:"status,omitempty"`
	Title      string      `json:"title,omitempty"`
	ParentID   string      `json:"parentId,omitempty"`
	ParentType string      `json:"parentType,omitempty"`
	Position   int         `json:"position,omitempty"`
	AuthorID   string      `json:"authorId,omitempty"`
	OwnerID    string      `json:"ownerId,omitempty"`
	CreatedAt  string      `json:"createdAt,omitempty"`
	EmbedURL   string      `json:"embedUrl,omitempty"`
	SpaceID    string      `json:"spaceId,omitempty"`
	Version    TreeVersion `json:"version,omitempty"`
	Links      struct {
		Base  string `json:"base,omitempty"`
		WebUI string `json:"webui,omitempty"`
		Self  string `json:"self,omitempty"`
	} `json:"_links,omitempty"`
	Collaborators  any         `json:"collaborators,omitempty"`
	DirectChildren any         `json:"directChildren,omitempty"`
	Operations     []Operation `json:"operations,omitempty"`
	Properties     any         `json:"properties,omitempty"`
}

// TreeVersion is the version subdocument used by Cloud v2 content-tree entities.
type TreeVersion struct {
	CreatedAt string `json:"createdAt,omitempty"`
	Message   string `json:"message,omitempty"`
	Number    int    `json:"number,omitempty"`
	MinorEdit bool   `json:"minorEdit,omitempty"`
	AuthorID  string `json:"authorId,omitempty"`
}

// GetContentTreeEntity returns one Cloud v2 content-tree entity by id.
func GetContentTreeEntity(ctx context.Context, c *client.Client, kind, id string, opts ContentTreeEntityOptions) (*ContentTreeEntity, error) {
	if !isCloud(c) {
		return nil, fmt.Errorf("GetContentTreeEntity: Confluence Cloud only")
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("GetContentTreeEntity: id is required")
	}
	path, err := contentTreeEntityPath(kind, id)
	if err != nil {
		return nil, err
	}
	data, _, err := c.Get(ctx, path, contentTreeEntityParams(opts))
	if err != nil {
		return nil, err
	}
	var out ContentTreeEntity
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parse content-tree entity: %w", err)
	}
	return &out, nil
}

// ListContentTreeDirectChildren lists direct children under one Cloud v2
// database, folder, whiteboard, or Smart Link embed.
func ListContentTreeDirectChildren(ctx context.Context, c *client.Client, kind, id string, opts DirectChildrenOptions) ([]Content, error) {
	if !isCloud(c) {
		return nil, fmt.Errorf("ListContentTreeDirectChildren: Confluence Cloud only")
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("ListContentTreeDirectChildren: id is required")
	}
	prefix, err := contentTreeEntityPlural(kind)
	if err != nil {
		return nil, err
	}
	opts.Limit = clampLimit(opts.Limit)
	opts.Types = normalizeContentTypes(opts.Types)
	path := "/api/v2/" + prefix + "/" + url.PathEscape(id) + "/direct-children"
	return listCloudDirectChildren(ctx, c, path, opts)
}

func contentTreeEntityPath(kind, id string) (string, error) {
	plural, err := contentTreeEntityPlural(kind)
	if err != nil {
		return "", err
	}
	return "/api/v2/" + plural + "/" + url.PathEscape(id), nil
}

func contentTreeEntityPlural(kind string) (string, error) {
	switch normalizeEntityTarget(kind) {
	case "database":
		return "databases", nil
	case "folder":
		return "folders", nil
	case "whiteboard":
		return "whiteboards", nil
	case "embed":
		return "embeds", nil
	default:
		return "", fmt.Errorf("unsupported content-tree entity kind %q", kind)
	}
}

func contentTreeEntityParams(opts ContentTreeEntityOptions) url.Values {
	params := url.Values{}
	setBoolParam(params, "include-collaborators", opts.IncludeCollaborators)
	setBoolParam(params, "include-direct-children", opts.IncludeDirectChildren)
	setBoolParam(params, "include-operations", opts.IncludeOperations)
	setBoolParam(params, "include-properties", opts.IncludeProperties)
	return params
}

func setBoolParam(params url.Values, key string, value bool) {
	if value {
		params.Set(key, "true")
	}
}

func listCloudDirectChildren(ctx context.Context, c *client.Client, path string, opts DirectChildrenOptions) ([]Content, error) {
	params := url.Values{"limit": {strconv.Itoa(opts.Limit)}}
	if opts.Sort != "" {
		params.Set("sort", opts.Sort)
	}
	out := make([]Content, 0, opts.Limit)
	for len(out) < opts.Limit {
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
			return nil, fmt.Errorf("parse direct children: %w", err)
		}
		for _, child := range page.Results {
			if !contentTypeAllowed(opts.Types, child.Type) {
				continue
			}
			out = append(out, child)
			if len(out) == opts.Limit {
				break
			}
		}

		next := nextPageURL(headers, page.Links.Next)
		if next == "" || len(page.Results) == 0 || len(out) == opts.Limit {
			break
		}
		path, params, err = nextPageRequest(c, next)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}
