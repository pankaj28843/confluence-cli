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

// Group is one row of a Confluence group collection.
type Group struct {
	ID        string         `json:"id,omitempty"`
	Name      string         `json:"name"`
	Type      string         `json:"type,omitempty"` // "group"
	UsageType string         `json:"usageType,omitempty"`
	ManagedBy string         `json:"managedBy,omitempty"`
	Links     map[string]any `json:"_links,omitempty"`
}

// GroupListOptions controls group listing.
type GroupListOptions struct {
	Limit      int
	AccessType string
}

// GroupLookupOptions identifies one group.
type GroupLookupOptions struct {
	ID     string
	Name   string
	Expand string
}

// GroupMemberOptions controls group member listing.
type GroupMemberOptions struct {
	GroupID               string
	GroupName             string
	Limit                 int
	Expand                []string
	ShouldReturnTotalSize bool
}

// GroupPickerOptions controls Cloud group picker searches.
type GroupPickerOptions struct {
	Limit                 int
	ShouldReturnTotalSize bool
}

// GroupRelationOptions controls Server/DC group hierarchy listing.
type GroupRelationOptions struct {
	GroupName string
	Limit     int
	Expand    string
}

// ListGroups pages the documented group collection endpoint.
func ListGroups(ctx context.Context, c *client.Client, limit int) ([]Group, error) {
	return ListGroupsWithOptions(ctx, c, GroupListOptions{Limit: limit})
}

// ListGroupsWithOptions pages /rest/api/group.
func ListGroupsWithOptions(ctx context.Context, c *client.Client, opts GroupListOptions) ([]Group, error) {
	if opts.AccessType != "" && !isCloud(c) {
		return nil, fmt.Errorf("unsupported Server/Data Center group access type filter")
	}
	limit := clampLimit(opts.Limit)
	params := url.Values{
		"start": {"0"},
		"limit": {strconv.Itoa(limit)},
	}
	if opts.AccessType != "" {
		params.Set("accessType", opts.AccessType)
	}
	return collectGroupPages(ctx, c, "/rest/api/group", params, limit, "parse groups")
}

// GetGroup fetches one group using Cloud group id or Server/DC group name.
func GetGroup(ctx context.Context, c *client.Client, opts GroupLookupOptions) (*Group, error) {
	var path string
	params := url.Values{}
	if isCloud(c) {
		if opts.ID == "" {
			return nil, fmt.Errorf("Cloud group id is required")
		}
		path = "/rest/api/group/by-id"
		params.Set("id", opts.ID)
	} else {
		if opts.Name == "" {
			return nil, fmt.Errorf("Server/Data Center group name is required")
		}
		path = "/rest/api/group/" + url.PathEscape(opts.Name)
		if opts.Expand != "" {
			params.Set("expand", opts.Expand)
		}
	}
	data, _, err := c.Get(ctx, path, params)
	if err != nil {
		return nil, err
	}
	var group Group
	if err := json.Unmarshal(data, &group); err != nil {
		return nil, fmt.Errorf("parse group: %w", err)
	}
	return &group, nil
}

// ListGroupMembers lists members of a Server/DC group by name.
func ListGroupMembers(ctx context.Context, c *client.Client, groupName string, limit int) ([]User, error) {
	return ListGroupMembersWithOptions(ctx, c, GroupMemberOptions{GroupName: groupName, Limit: limit})
}

// ListGroupMembersWithOptions lists group members using the documented
// flavor-specific group-member endpoint.
func ListGroupMembersWithOptions(ctx context.Context, c *client.Client, opts GroupMemberOptions) ([]User, error) {
	limit := clampLimit(opts.Limit)
	params := url.Values{
		"start": {"0"},
		"limit": {strconv.Itoa(limit)},
	}
	if len(opts.Expand) > 0 {
		params.Set("expand", strings.Join(opts.Expand, ","))
	}
	if opts.ShouldReturnTotalSize {
		params.Set("shouldReturnTotalSize", "true")
	}

	var path string
	if isCloud(c) {
		if opts.GroupID == "" {
			return nil, fmt.Errorf("Cloud group id is required")
		}
		path = "/rest/api/group/" + url.PathEscape(opts.GroupID) + "/membersByGroupId"
	} else {
		if opts.GroupName == "" {
			return nil, fmt.Errorf("Server/Data Center group name is required")
		}
		params.Del("shouldReturnTotalSize")
		path = "/rest/api/group/" + url.PathEscape(opts.GroupName) + "/member"
	}
	return collectUserPages(ctx, c, path, params, limit, "parse group members")
}

// PickGroups searches Cloud groups using the group picker endpoint.
func PickGroups(ctx context.Context, c *client.Client, query string, opts GroupPickerOptions) ([]Group, error) {
	if !isCloud(c) {
		return nil, fmt.Errorf("unsupported Server/Data Center group picker")
	}
	if query == "" {
		return nil, fmt.Errorf("group picker query is required")
	}
	limit := clampLimit(opts.Limit)
	params := url.Values{
		"query": {query},
		"start": {"0"},
		"limit": {strconv.Itoa(limit)},
	}
	if opts.ShouldReturnTotalSize {
		params.Set("shouldReturnTotalSize", "true")
	}
	return collectGroupPages(ctx, c, "/rest/api/group/picker", params, limit, "parse group picker")
}

// ListGroupChildGroups lists direct child groups for a Server/DC group.
func ListGroupChildGroups(ctx context.Context, c *client.Client, opts GroupRelationOptions) ([]Group, error) {
	return listGroupRelation(ctx, c, opts, "groupmember")
}

// ListGroupParents lists direct parent groups for a Server/DC group.
func ListGroupParents(ctx context.Context, c *client.Client, opts GroupRelationOptions) ([]Group, error) {
	return listGroupRelation(ctx, c, opts, "groupparent")
}

// ListGroupAncestors lists ancestor groups for a Server/DC group.
func ListGroupAncestors(ctx context.Context, c *client.Client, opts GroupRelationOptions) ([]Group, error) {
	return listGroupRelation(ctx, c, opts, "groupancestor")
}

func listGroupRelation(ctx context.Context, c *client.Client, opts GroupRelationOptions, relation string) ([]Group, error) {
	if isCloud(c) {
		return nil, fmt.Errorf("unsupported Cloud group relation %q", relation)
	}
	if opts.GroupName == "" {
		return nil, fmt.Errorf("Server/Data Center group name is required")
	}
	limit := clampLimit(opts.Limit)
	params := url.Values{
		"start": {"0"},
		"limit": {strconv.Itoa(limit)},
	}
	if opts.Expand != "" {
		params.Set("expand", opts.Expand)
	}
	path := "/rest/api/group/" + url.PathEscape(opts.GroupName) + "/" + relation
	return collectGroupPages(ctx, c, path, params, limit, "parse group relations")
}

func collectGroupPages(ctx context.Context, c *client.Client, path string, params url.Values, limit int, parseContext string) ([]Group, error) {
	out := make([]Group, 0, limit)
	for len(out) < limit {
		data, headers, err := c.Get(ctx, path, params)
		if err != nil {
			return nil, err
		}
		var page struct {
			Results []Group `json:"results"`
			Links   struct {
				Next string `json:"next,omitempty"`
			} `json:"_links"`
		}
		if err := json.Unmarshal(data, &page); err != nil {
			return nil, fmt.Errorf("%s: %w", parseContext, err)
		}
		for _, group := range page.Results {
			out = append(out, group)
			if len(out) == limit {
				break
			}
		}
		next := nextPageURL(headers, page.Links.Next)
		if next == "" || len(page.Results) == 0 || len(out) == limit {
			break
		}
		var nextErr error
		path, params, nextErr = nextPageRequest(c, next)
		if nextErr != nil {
			return nil, nextErr
		}
	}
	return out, nil
}

func collectUserPages(ctx context.Context, c *client.Client, path string, params url.Values, limit int, parseContext string) ([]User, error) {
	out := make([]User, 0, limit)
	for len(out) < limit {
		data, headers, err := c.Get(ctx, path, params)
		if err != nil {
			return nil, err
		}
		var page struct {
			Results []User `json:"results"`
			Links   struct {
				Next string `json:"next,omitempty"`
			} `json:"_links"`
		}
		if err := json.Unmarshal(data, &page); err != nil {
			return nil, fmt.Errorf("%s: %w", parseContext, err)
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
		var nextErr error
		path, params, nextErr = nextPageRequest(c, next)
		if nextErr != nil {
			return nil, nextErr
		}
	}
	return out, nil
}
