// Package conf implements Confluence resource-level operations.
package conf

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

// CurrentUser is the shared subset of /rest/api/user/current for both flavors.
// Server/DC fills {Username, DisplayName, UserKey}; Cloud fills {AccountID,
// PublicName}. Both populate Email (Cloud) / UniqueName when available.
type CurrentUser struct {
	Username    string `json:"username,omitempty"`
	UserKey     string `json:"userKey,omitempty"`
	AccountID   string `json:"accountId,omitempty"`
	PublicName  string `json:"publicName,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Email       string `json:"email,omitempty"`
	Type        string `json:"type,omitempty"`
}

// GetCurrentUser fetches /rest/api/user/current (same path on both flavors).
func GetCurrentUser(ctx context.Context, c *client.Client) (*CurrentUser, error) {
	data, _, err := c.Get(ctx, "/rest/api/user/current", nil)
	if err != nil {
		return nil, err
	}
	var u CurrentUser
	if err := json.Unmarshal(data, &u); err != nil {
		return nil, fmt.Errorf("parse current user: %w", err)
	}
	return &u, nil
}

// Label reflects how the current user shows up in text.
func (u *CurrentUser) Label() string {
	switch {
	case u.DisplayName != "":
		return u.DisplayName
	case u.PublicName != "":
		return u.PublicName
	case u.Username != "":
		return u.Username
	default:
		return u.AccountID
	}
}

// ID returns the primary identifier for this flavor.
func (u *CurrentUser) ID() string {
	if u.Username != "" {
		return u.Username
	}
	return u.AccountID
}
