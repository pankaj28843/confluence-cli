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

// VersionListOptions controls read-only content version listing.
type VersionListOptions struct {
	Limit      int
	Sort       string
	BodyFormat string
}

// Version is the shared subset of Cloud v2 DetailedVersion/Version and
// Server/Data Center content version responses.
type Version struct {
	Number              int         `json:"number,omitempty"`
	AuthorID            string      `json:"authorId,omitempty"`
	Message             string      `json:"message,omitempty"`
	CreatedAt           string      `json:"createdAt,omitempty"`
	MinorEdit           bool        `json:"minorEdit,omitempty"`
	ContentTypeModified bool        `json:"contentTypeModified,omitempty"`
	Collaborators       []string    `json:"collaborators,omitempty"`
	PrevVersion         int         `json:"prevVersion,omitempty"`
	NextVersion         int         `json:"nextVersion,omitempty"`
	Page                *Content    `json:"page,omitempty"`
	BlogPost            *Content    `json:"blogpost,omitempty"`
	Attachment          *Attachment `json:"attachment,omitempty"`
	Comment             *Comment    `json:"comment,omitempty"`
	Custom              *Content    `json:"custom,omitempty"`
	By                  struct {
		DisplayName string `json:"displayName,omitempty"`
		PublicName  string `json:"publicName,omitempty"`
		Username    string `json:"username,omitempty"`
		AccountID   string `json:"accountId,omitempty"`
	} `json:"by,omitempty"`
	When string `json:"when,omitempty"`
}

// ListVersions returns version rows for a page, blogpost, attachment, footer
// comment, inline comment, custom content, or Server/DC content id.
func ListVersions(ctx context.Context, c *client.Client, target, id string, opts VersionListOptions) ([]Version, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("ListVersions: id is required")
	}
	opts.Limit = clampLimit(opts.Limit)
	path, err := versionListPath(c, target, id)
	if err != nil {
		return nil, err
	}
	params := versionListParams(c, target, opts)
	out := make([]Version, 0, opts.Limit)
	for len(out) < opts.Limit {
		data, headers, err := c.Get(ctx, path, params)
		if err != nil {
			return nil, err
		}
		var page struct {
			Results []Version `json:"results"`
			Links   struct {
				Next string `json:"next,omitempty"`
			} `json:"_links"`
		}
		if err := json.Unmarshal(data, &page); err != nil {
			return nil, fmt.Errorf("parse versions: %w", err)
		}
		for _, version := range page.Results {
			out = append(out, version)
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

// GetVersion returns one version record for a page, blogpost, attachment,
// footer comment, inline comment, custom content, or Server/DC content id.
func GetVersion(ctx context.Context, c *client.Client, target, id string, number int) (*Version, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("GetVersion: id is required")
	}
	if number <= 0 {
		return nil, fmt.Errorf("GetVersion: version number must be positive")
	}
	path, err := versionDetailPath(c, target, id, number)
	if err != nil {
		return nil, err
	}
	data, _, err := c.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}
	var out Version
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parse version: %w", err)
	}
	return &out, nil
}

func versionListParams(c *client.Client, target string, opts VersionListOptions) url.Values {
	params := url.Values{"limit": {strconv.Itoa(opts.Limit)}}
	if isCloud(c) {
		if opts.Sort != "" {
			params.Set("sort", opts.Sort)
		}
		if opts.BodyFormat != "" && cloudVersionBodyFormatTarget(target) {
			params.Set("body-format", opts.BodyFormat)
		}
	}
	return params
}

func versionListPath(c *client.Client, target, id string) (string, error) {
	if isCloud(c) {
		prefix, err := cloudVersionPrefix(target)
		if err != nil {
			return "", err
		}
		return "/api/v2/" + prefix + "/" + url.PathEscape(id) + "/versions", nil
	}
	if !serverContentVersionTarget(target) {
		return "", fmt.Errorf("unsupported Server/Data Center version target %q", target)
	}
	return "/rest/api/content/" + url.PathEscape(id) + "/version", nil
}

func versionDetailPath(c *client.Client, target, id string, number int) (string, error) {
	if isCloud(c) {
		prefix, err := cloudVersionPrefix(target)
		if err != nil {
			return "", err
		}
		return "/api/v2/" + prefix + "/" + url.PathEscape(id) + "/versions/" + strconv.Itoa(number), nil
	}
	if !serverContentVersionTarget(target) {
		return "", fmt.Errorf("unsupported Server/Data Center version target %q", target)
	}
	return "/rest/api/content/" + url.PathEscape(id) + "/version/" + strconv.Itoa(number), nil
}

func cloudVersionPrefix(target string) (string, error) {
	switch normalizeEntityTarget(target) {
	case "page":
		return "pages", nil
	case "blogpost":
		return "blogposts", nil
	case "attachment":
		return "attachments", nil
	case "footer-comment":
		return "footer-comments", nil
	case "inline-comment":
		return "inline-comments", nil
	case "custom-content":
		return "custom-content", nil
	default:
		return "", fmt.Errorf("unsupported version target %q", target)
	}
}

func cloudVersionBodyFormatTarget(target string) bool {
	switch normalizeEntityTarget(target) {
	case "page", "blogpost", "footer-comment", "inline-comment", "custom-content":
		return true
	default:
		return false
	}
}

func serverContentVersionTarget(target string) bool {
	switch normalizeEntityTarget(target) {
	case "content", "page", "blogpost", "footer-comment", "inline-comment", "comment":
		return true
	default:
		return false
	}
}
