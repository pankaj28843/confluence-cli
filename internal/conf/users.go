package conf

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

// GetUser fetches /rest/api/user.
// Server/DC accepts ?username= or ?key=; Cloud accepts ?accountId=.
// The caller chooses which param via selector. If selector is "auto", we guess
// based on flavor.
func GetUser(ctx context.Context, c *client.Client, selector, value string) (*User, error) {
	if selector == "" || selector == "auto" {
		if c.Flavor == client.FlavorCloud {
			selector = "accountId"
		} else {
			selector = "username"
		}
	}
	params := url.Values{selector: {value}}
	data, _, err := c.Get(ctx, "/rest/api/user", params)
	if err != nil {
		return nil, err
	}
	var u User
	if err := json.Unmarshal(data, &u); err != nil {
		return nil, fmt.Errorf("parse user: %w", err)
	}
	return &u, nil
}

// SearchUsers runs a CQL-like user search via /rest/api/search?cql=user.fullname~"Q".
// Returns SearchHit rows with .User populated.
func SearchUsers(ctx context.Context, c *client.Client, query string, limit int) ([]SearchHit, error) {
	if limit <= 0 {
		limit = 25
	}
	if limit > 200 {
		limit = 200
	}
	cql := fmt.Sprintf(`type=user AND user.fullname ~ "%s"`, query)
	params := url.Values{"cql": {cql}, "limit": {strconv.Itoa(limit)}}
	data, _, err := c.Get(ctx, "/rest/api/search", params)
	if err != nil {
		return nil, err
	}
	var page struct {
		Results []SearchHit `json:"results"`
	}
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parse user search: %w", err)
	}
	return page.Results, nil
}
