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

// SpacePermissionAssignment is a granted permission on a Confluence space.
// Cloud returns principal/operation.key, while Server/DC returns
// subject/operation.operationKey.
type SpacePermissionAssignment struct {
	ID        string                   `json:"id,omitempty"`
	Principal SpacePermissionPrincipal `json:"principal,omitempty"`
	Operation SpacePermissionOperation `json:"operation,omitempty"`
	Subject   SpacePermissionSubject   `json:"subject,omitempty"`
	SpaceKey  string                   `json:"spaceKey,omitempty"`
	SpaceID   int64                    `json:"spaceId,omitempty"`
}

// SpacePermissionPrincipal identifies a Cloud v2 permission principal.
type SpacePermissionPrincipal struct {
	Type string `json:"type,omitempty"`
	ID   string `json:"id,omitempty"`
}

// SpacePermissionSubject identifies a Server/DC permission subject.
type SpacePermissionSubject struct {
	DisplayName string `json:"displayName,omitempty"`
	Type        string `json:"type,omitempty"`
	ID          string `json:"id,omitempty"`
	UserKey     string `json:"userKey,omitempty"`
	GroupName   string `json:"groupName,omitempty"`
	AccountID   string `json:"accountId,omitempty"`
}

// SpacePermissionOperation identifies the operation allowed by a permission.
type SpacePermissionOperation struct {
	TargetType   string `json:"targetType,omitempty"`
	OperationKey string `json:"operationKey,omitempty"`
	Key          string `json:"key,omitempty"`
}

// Name returns the flavor-neutral operation key.
func (o SpacePermissionOperation) Name() string {
	if o.OperationKey != "" {
		return o.OperationKey
	}
	return o.Key
}

// SpacePermissionDefinition describes an available Cloud space permission.
type SpacePermissionDefinition struct {
	ID                    string   `json:"id,omitempty"`
	DisplayName           string   `json:"displayName,omitempty"`
	Description           string   `json:"description,omitempty"`
	RequiredPermissionIDs []string `json:"requiredPermissionIds,omitempty"`
}

// SpacePermissionSubjectSelector selects one Server/DC subject-specific read.
type SpacePermissionSubjectSelector struct {
	Anonymous bool
	GroupName string
	UserKey   string
}

// ListSpacePermissionAssignments lists permissions granted on a space.
func ListSpacePermissionAssignments(ctx context.Context, c *client.Client, space string, limit int) ([]SpacePermissionAssignment, error) {
	space = strings.TrimSpace(space)
	if space == "" {
		return nil, fmt.Errorf("ListSpacePermissionAssignments: space is required")
	}
	if c == nil {
		return nil, fmt.Errorf("ListSpacePermissionAssignments: client is required")
	}
	limit = clampLimit(limit)
	if isCloud(c) {
		spaceID, err := resolveCloudSpaceID(ctx, c, space)
		if err != nil {
			return nil, err
		}
		path := "/api/v2/spaces/" + url.PathEscape(spaceID) + "/permissions"
		return listCloudSpacePermissionAssignments(ctx, c, path, url.Values{"limit": {strconv.Itoa(limit)}}, limit)
	}

	data, _, err := c.Get(ctx, "/rest/api/space/"+url.PathEscape(space)+"/permissions", nil)
	if err != nil {
		return nil, err
	}
	assignments, err := parseServerSpacePermissionAssignments(data)
	if err != nil {
		return nil, err
	}
	return capSpacePermissionAssignments(assignments, limit), nil
}

// ListAvailableSpacePermissions lists the Cloud v2 available permission catalog.
func ListAvailableSpacePermissions(ctx context.Context, c *client.Client, limit int) ([]SpacePermissionDefinition, error) {
	if c == nil {
		return nil, fmt.Errorf("ListAvailableSpacePermissions: client is required")
	}
	if !isCloud(c) {
		return nil, fmt.Errorf("ListAvailableSpacePermissions: Confluence Cloud only")
	}
	limit = clampLimit(limit)
	return listCloudSpacePermissionDefinitions(ctx, c, "/api/v2/space-permissions", url.Values{"limit": {strconv.Itoa(limit)}}, limit)
}

// ListSpacePermissionsForSubject lists Server/DC permissions for one subject in a space.
func ListSpacePermissionsForSubject(ctx context.Context, c *client.Client, space string, selector SpacePermissionSubjectSelector, limit int) ([]SpacePermissionAssignment, error) {
	space = strings.TrimSpace(space)
	if space == "" {
		return nil, fmt.Errorf("ListSpacePermissionsForSubject: space is required")
	}
	if c == nil {
		return nil, fmt.Errorf("ListSpacePermissionsForSubject: client is required")
	}
	segment, err := selector.serverPathSegment()
	if err != nil {
		return nil, err
	}
	if isCloud(c) {
		return nil, fmt.Errorf("ListSpacePermissionsForSubject: Confluence Server/Data Center only")
	}
	limit = clampLimit(limit)
	path := "/rest/api/space/" + url.PathEscape(space) + "/permissions/" + segment
	data, _, err := c.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}
	assignments, err := parseServerSpacePermissionAssignments(data)
	if err != nil {
		return nil, err
	}
	return capSpacePermissionAssignments(assignments, limit), nil
}

func (s SpacePermissionSubjectSelector) serverPathSegment() (string, error) {
	group := strings.TrimSpace(s.GroupName)
	userKey := strings.TrimSpace(s.UserKey)
	count := 0
	if s.Anonymous {
		count++
	}
	if group != "" {
		count++
	}
	if userKey != "" {
		count++
	}
	if count != 1 {
		return "", fmt.Errorf("SpacePermissionSubjectSelector: select exactly one of anonymous, group, or user key")
	}
	switch {
	case s.Anonymous:
		return "anonymous", nil
	case group != "":
		return "group/" + url.PathEscape(group), nil
	default:
		return "user/" + url.PathEscape(userKey), nil
	}
}

func listCloudSpacePermissionAssignments(ctx context.Context, c *client.Client, path string, params url.Values, limit int) ([]SpacePermissionAssignment, error) {
	out := make([]SpacePermissionAssignment, 0, limit)
	for len(out) < limit {
		data, headers, err := c.Get(ctx, path, params)
		if err != nil {
			return nil, err
		}
		var page struct {
			Results []SpacePermissionAssignment `json:"results"`
			Links   struct {
				Next string `json:"next,omitempty"`
			} `json:"_links"`
		}
		if err := json.Unmarshal(data, &page); err != nil {
			return nil, fmt.Errorf("parse space permission assignments: %w", err)
		}
		for _, assignment := range page.Results {
			out = append(out, assignment)
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

func listCloudSpacePermissionDefinitions(ctx context.Context, c *client.Client, path string, params url.Values, limit int) ([]SpacePermissionDefinition, error) {
	out := make([]SpacePermissionDefinition, 0, limit)
	for len(out) < limit {
		data, headers, err := c.Get(ctx, path, params)
		if err != nil {
			return nil, err
		}
		var page struct {
			Results []SpacePermissionDefinition `json:"results"`
			Links   struct {
				Next string `json:"next,omitempty"`
			} `json:"_links"`
		}
		if err := json.Unmarshal(data, &page); err != nil {
			return nil, fmt.Errorf("parse available space permissions: %w", err)
		}
		for _, permission := range page.Results {
			out = append(out, permission)
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

func parseServerSpacePermissionAssignments(data []byte) ([]SpacePermissionAssignment, error) {
	var assignments []SpacePermissionAssignment
	if err := json.Unmarshal(data, &assignments); err == nil {
		return assignments, nil
	}

	var page struct {
		Results []SpacePermissionAssignment `json:"results"`
	}
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parse space permission assignments: %w", err)
	}
	return page.Results, nil
}

func capSpacePermissionAssignments(assignments []SpacePermissionAssignment, limit int) []SpacePermissionAssignment {
	if len(assignments) > limit {
		return assignments[:limit]
	}
	return assignments
}
