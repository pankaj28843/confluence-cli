package conf

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

// Group is one row of /rest/api/group.
type Group struct {
	Name string `json:"name"`
	Type string `json:"type,omitempty"` // "group"
}

// ListGroups pages /rest/api/group.
func ListGroups(ctx context.Context, c *client.Client, limit int) ([]Group, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	params := url.Values{"limit": {strconv.Itoa(limit)}}
	data, _, err := c.Get(ctx, "/rest/api/group", params)
	if err != nil {
		return nil, err
	}
	var page struct {
		Results []Group `json:"results"`
	}
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parse groups: %w", err)
	}
	return page.Results, nil
}

// ListGroupMembers lists members of a group by name.
func ListGroupMembers(ctx context.Context, c *client.Client, groupName string, limit int) ([]User, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 200 {
		limit = 200
	}
	params := url.Values{"limit": {strconv.Itoa(limit)}}
	data, _, err := c.Get(ctx, "/rest/api/group/"+url.PathEscape(groupName)+"/member", params)
	if err != nil {
		return nil, err
	}
	var page struct {
		Results []User `json:"results"`
	}
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parse group members: %w", err)
	}
	return page.Results, nil
}
