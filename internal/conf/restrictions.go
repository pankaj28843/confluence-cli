package conf

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

const defaultRestrictionExpand = "restrictions.user,restrictions.group"

// ContentRestrictionsByOperation is the documented content restriction-by-operation
// surface. Server/DC wraps the operation map in "restrictions"; Cloud may return
// operation keys at the top level, so UnmarshalJSON accepts both shapes.
type ContentRestrictionsByOperation struct {
	Restrictions map[string]ContentRestriction `json:"restrictions,omitempty"`
	Links        map[string]any                `json:"_links,omitempty"`
	GetLinks     map[string]any                `json:"get_links,omitempty"`
}

// ContentRestriction is one read/update restriction entry.
type ContentRestriction struct {
	Operation     string              `json:"operation,omitempty"`
	Restrictions  RestrictionSubjects `json:"restrictions,omitempty"`
	Content       map[string]any      `json:"content,omitempty"`
	Links         map[string]any      `json:"_links,omitempty"`
	GetLinks      map[string]any      `json:"get_links,omitempty"`
	Expandable    map[string]any      `json:"_expandable,omitempty"`
	GetExpandable map[string]any      `json:"get_expandable,omitempty"`
	LastModified  string              `json:"lastModificationDate,omitempty"`
}

// RestrictionSubjects holds the user and group subject pages for one operation.
type RestrictionSubjects struct {
	User  RestrictionSubjectPage `json:"user,omitempty"`
	Group RestrictionSubjectPage `json:"group,omitempty"`
}

// RestrictionSubjectPage is intentionally loose because Cloud and Server/DC
// expose different user/group fields while keeping the paging envelope stable.
type RestrictionSubjectPage struct {
	Results    []map[string]any `json:"results,omitempty"`
	Start      int              `json:"start,omitempty"`
	Limit      int              `json:"limit,omitempty"`
	Size       int              `json:"size,omitempty"`
	TotalSize  int64            `json:"totalSize,omitempty"`
	TotalCount int64            `json:"totalCount,omitempty"`
	Links      map[string]any   `json:"_links,omitempty"`
	Page       map[string]any   `json:"pageRequest,omitempty"`
	NextCursor string           `json:"nextCursor,omitempty"`
	PrevCursor string           `json:"prevCursor,omitempty"`
}

func (r *ContentRestrictionsByOperation) UnmarshalJSON(data []byte) error {
	var wire struct {
		Restrictions map[string]ContentRestriction `json:"restrictions"`
		Links        map[string]any                `json:"_links"`
		GetLinks     map[string]any                `json:"get_links"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	r.Restrictions = wire.Restrictions
	r.Links = wire.Links
	r.GetLinks = wire.GetLinks

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if r.Restrictions == nil {
		r.Restrictions = make(map[string]ContentRestriction)
	}
	for key, msg := range raw {
		if isRestrictionMetadataKey(key) {
			continue
		}
		restriction, ok := decodeOperationRestriction(key, msg)
		if ok {
			r.Restrictions[key] = restriction
		}
	}
	return nil
}

// Operations returns restrictions sorted by operation name for stable text output.
func (r ContentRestrictionsByOperation) Operations() []ContentRestriction {
	out := make([]ContentRestriction, 0, len(r.Restrictions))
	keys := make([]string, 0, len(r.Restrictions))
	for key := range r.Restrictions {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		restriction := r.Restrictions[key]
		if restriction.Operation == "" {
			restriction.Operation = key
		}
		out = append(out, restriction)
	}
	return out
}

func (r ContentRestriction) UserCount() int {
	return r.Restrictions.User.Count()
}

func (r ContentRestriction) GroupCount() int {
	return r.Restrictions.Group.Count()
}

func (p RestrictionSubjectPage) Count() int {
	maxInt := int64(^uint(0) >> 1)
	switch {
	case len(p.Results) > 0:
		return len(p.Results)
	case p.Size > 0:
		return p.Size
	case p.TotalSize > 0 && p.TotalSize <= maxInt:
		return int(p.TotalSize)
	case p.TotalCount > 0 && p.TotalCount <= maxInt:
		return int(p.TotalCount)
	default:
		return 0
	}
}

// ListContentRestrictions returns content restrictions grouped by operation.
func ListContentRestrictions(ctx context.Context, c *client.Client, contentID string) (*ContentRestrictionsByOperation, error) {
	contentID = strings.TrimSpace(contentID)
	if contentID == "" {
		return nil, fmt.Errorf("ListContentRestrictions: content id is required")
	}
	if c == nil {
		return nil, fmt.Errorf("ListContentRestrictions: client is required")
	}
	params := url.Values{"expand": {defaultRestrictionExpand}}
	path := "/rest/api/content/" + url.PathEscape(contentID) + "/restriction/byOperation"
	data, _, err := c.Get(ctx, path, params)
	if err != nil {
		return nil, err
	}
	var out ContentRestrictionsByOperation
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parse content restrictions: %w", err)
	}
	if out.Restrictions == nil {
		out.Restrictions = make(map[string]ContentRestriction)
	}
	return &out, nil
}

// GetContentRestrictionForOperation returns restrictions for one operation.
func GetContentRestrictionForOperation(ctx context.Context, c *client.Client, contentID, operation string, limit int) (*ContentRestriction, error) {
	contentID = strings.TrimSpace(contentID)
	if contentID == "" {
		return nil, fmt.Errorf("GetContentRestrictionForOperation: content id is required")
	}
	operation, err := normalizeRestrictionOperation(operation)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, fmt.Errorf("GetContentRestrictionForOperation: client is required")
	}
	limit = clampLimit(limit)
	params := url.Values{
		"expand": {defaultRestrictionExpand},
		"start":  {"0"},
		"limit":  {strconv.Itoa(limit)},
	}
	path := "/rest/api/content/" + url.PathEscape(contentID) + "/restriction/byOperation/" + url.PathEscape(operation)
	data, _, err := c.Get(ctx, path, params)
	if err != nil {
		return nil, err
	}
	var out ContentRestriction
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parse content restriction: %w", err)
	}
	if out.Operation == "" {
		out.Operation = operation
	}
	return &out, nil
}

func normalizeRestrictionOperation(operation string) (string, error) {
	operation = strings.ToLower(strings.TrimSpace(operation))
	switch operation {
	case "read", "update":
		return operation, nil
	default:
		return "", fmt.Errorf("unsupported restriction operation %q (use read or update)", operation)
	}
}

func isRestrictionMetadataKey(key string) bool {
	switch key {
	case "restrictions", "_links", "get_links", "_expandable", "get_expandable", "links", "start", "limit", "size", "results", "restrictionsHash":
		return true
	default:
		return strings.HasPrefix(key, "_")
	}
}

func decodeOperationRestriction(key string, msg json.RawMessage) (ContentRestriction, bool) {
	var restriction ContentRestriction
	if err := json.Unmarshal(msg, &restriction); err != nil {
		return ContentRestriction{}, false
	}
	if restriction.Operation == "" {
		var wrapped struct {
			OperationType ContentRestriction `json:"operationType"`
		}
		if err := json.Unmarshal(msg, &wrapped); err == nil && wrapped.OperationType.Operation != "" {
			restriction = wrapped.OperationType
		}
	}
	if restriction.Operation == "" && restriction.hasSubjects() {
		restriction.Operation = key
	}
	return restriction, restriction.Operation != ""
}

func (r ContentRestriction) hasSubjects() bool {
	return len(r.Restrictions.User.Results) > 0 ||
		len(r.Restrictions.Group.Results) > 0 ||
		r.Restrictions.User.Size > 0 ||
		r.Restrictions.Group.Size > 0 ||
		r.Restrictions.User.TotalSize > 0 ||
		r.Restrictions.Group.TotalSize > 0 ||
		r.Restrictions.User.TotalCount > 0 ||
		r.Restrictions.Group.TotalCount > 0
}
