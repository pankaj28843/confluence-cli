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

// Operation is one permitted Confluence operation on an entity.
type Operation struct {
	Operation  string `json:"operation,omitempty"`
	TargetType string `json:"targetType,omitempty"`
}

// LikeCount is the Cloud v2 like-count response.
type LikeCount struct {
	Count int `json:"count"`
}

// LikeUser is one account id returned by Cloud v2 like-user listing.
type LikeUser struct {
	AccountID string `json:"accountId,omitempty"`
}

// ListOperations returns permitted operations for a target entity.
func ListOperations(ctx context.Context, c *client.Client, targetType, id string) ([]Operation, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("ListOperations: id is required")
	}
	if isCloud(c) {
		path, err := cloudOperationPath(ctx, c, targetType, id)
		if err != nil {
			return nil, err
		}
		return getOperations(ctx, c, path)
	}
	if !isServerContentOperationTarget(targetType) {
		return nil, fmt.Errorf("ListOperations: %s operations are only supported on Confluence Cloud", targetType)
	}
	data, _, err := c.Get(ctx, "/rest/api/content/"+url.PathEscape(id), url.Values{"expand": {"operations"}})
	if err != nil {
		return nil, err
	}
	var out struct {
		Operations []Operation `json:"operations"`
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parse operations: %w", err)
	}
	return out.Operations, nil
}

// GetLikeCount returns a Cloud v2 like count for a supported target entity.
func GetLikeCount(ctx context.Context, c *client.Client, targetType, id string) (*LikeCount, error) {
	if !isCloud(c) {
		return nil, fmt.Errorf("likes are only supported on Confluence Cloud")
	}
	path, err := cloudLikePath(targetType, id)
	if err != nil {
		return nil, err
	}
	data, _, err := c.Get(ctx, path+"/likes/count", nil)
	if err != nil {
		return nil, err
	}
	var out LikeCount
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parse like count: %w", err)
	}
	return &out, nil
}

// ListLikeUsers returns Cloud v2 account ids that liked a supported target.
func ListLikeUsers(ctx context.Context, c *client.Client, targetType, id string, limit int) ([]LikeUser, error) {
	if !isCloud(c) {
		return nil, fmt.Errorf("likes are only supported on Confluence Cloud")
	}
	path, err := cloudLikePath(targetType, id)
	if err != nil {
		return nil, err
	}
	limit = clampLimit(limit)
	params := url.Values{"limit": {strconv.Itoa(limit)}}
	path += "/likes/users"
	out := make([]LikeUser, 0, limit)

	for len(out) < limit {
		data, headers, err := c.Get(ctx, path, params)
		if err != nil {
			return nil, err
		}
		var page struct {
			Results []LikeUser `json:"results"`
			Links   struct {
				Next string `json:"next,omitempty"`
			} `json:"_links"`
		}
		if err := json.Unmarshal(data, &page); err != nil {
			return nil, fmt.Errorf("parse like users: %w", err)
		}
		for _, user := range page.Results {
			out = append(out, user)
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

func getOperations(ctx context.Context, c *client.Client, path string) ([]Operation, error) {
	data, _, err := c.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}
	var out struct {
		Operations []Operation `json:"operations"`
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parse operations: %w", err)
	}
	return out.Operations, nil
}

func cloudOperationPath(ctx context.Context, c *client.Client, targetType, id string) (string, error) {
	id = strings.TrimSpace(id)
	switch normalizeEntityTarget(targetType) {
	case "attachment":
		return "/api/v2/attachments/" + url.PathEscape(id) + "/operations", nil
	case "blogpost":
		return "/api/v2/blogposts/" + url.PathEscape(id) + "/operations", nil
	case "custom-content":
		return "/api/v2/custom-content/" + url.PathEscape(id) + "/operations", nil
	case "database":
		return "/api/v2/databases/" + url.PathEscape(id) + "/operations", nil
	case "embed":
		return "/api/v2/embeds/" + url.PathEscape(id) + "/operations", nil
	case "folder":
		return "/api/v2/folders/" + url.PathEscape(id) + "/operations", nil
	case "footer-comment":
		return "/api/v2/footer-comments/" + url.PathEscape(id) + "/operations", nil
	case "inline-comment":
		return "/api/v2/inline-comments/" + url.PathEscape(id) + "/operations", nil
	case "page":
		return "/api/v2/pages/" + url.PathEscape(id) + "/operations", nil
	case "space":
		spaceID := id
		if !isDecimalID(id) {
			space, err := getCloudSpaceV2(ctx, c, id)
			if err != nil {
				return "", err
			}
			spaceID = space.ID
		}
		if spaceID == "" {
			return "", fmt.Errorf("space %q has no id", id)
		}
		return "/api/v2/spaces/" + url.PathEscape(spaceID) + "/operations", nil
	case "whiteboard":
		return "/api/v2/whiteboards/" + url.PathEscape(id) + "/operations", nil
	default:
		return "", fmt.Errorf("unsupported operation target type %q", targetType)
	}
}

func cloudLikePath(targetType, id string) (string, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return "", fmt.Errorf("like target id is required")
	}
	switch normalizeEntityTarget(targetType) {
	case "blogpost":
		return "/api/v2/blogposts/" + url.PathEscape(id), nil
	case "footer-comment":
		return "/api/v2/footer-comments/" + url.PathEscape(id), nil
	case "inline-comment":
		return "/api/v2/inline-comments/" + url.PathEscape(id), nil
	case "page":
		return "/api/v2/pages/" + url.PathEscape(id), nil
	default:
		return "", fmt.Errorf("unsupported like target type %q", targetType)
	}
}

func isServerContentOperationTarget(targetType string) bool {
	switch normalizeEntityTarget(targetType) {
	case "attachment", "blogpost", "comment", "content", "page":
		return true
	default:
		return false
	}
}

func normalizeEntityTarget(targetType string) string {
	targetType = strings.ToLower(strings.TrimSpace(targetType))
	targetType = strings.ReplaceAll(targetType, "_", "-")
	switch targetType {
	case "blog-post", "blog":
		return "blogpost"
	case "comment-footer":
		return "footer-comment"
	case "comment-inline":
		return "inline-comment"
	case "embed", "smart-link", "smartlink":
		return "embed"
	case "space-id":
		return "space"
	default:
		return targetType
	}
}

func isDecimalID(id string) bool {
	if id == "" {
		return false
	}
	for _, r := range id {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
