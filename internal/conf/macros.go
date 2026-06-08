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

// MacroLookup identifies a macro instance on a historical content version.
type MacroLookup struct {
	ContentID string
	Version   int
	MacroID   string
	Hash      string
}

// MacroInstance is the documented macro body response.
type MacroInstance struct {
	Name       string         `json:"name,omitempty"`
	Body       string         `json:"body,omitempty"`
	Parameters map[string]any `json:"parameters,omitempty"`
	Links      map[string]any `json:"_links,omitempty"`
}

// GetMacroBody fetches a macro body in storage format for a content version.
func GetMacroBody(ctx context.Context, c *client.Client, in MacroLookup) (*MacroInstance, error) {
	contentID := strings.TrimSpace(in.ContentID)
	macroID := strings.TrimSpace(in.MacroID)
	hash := strings.TrimSpace(in.Hash)
	if contentID == "" {
		return nil, fmt.Errorf("GetMacroBody: content id is required")
	}
	if in.Version <= 0 {
		return nil, fmt.Errorf("GetMacroBody: version must be greater than zero")
	}
	if (macroID == "") == (hash == "") {
		return nil, fmt.Errorf("GetMacroBody: exactly one of macro id or hash is required")
	}
	if hash != "" && isCloud(c) {
		return nil, fmt.Errorf("GetMacroBody: hash lookup is only documented on Confluence Server/Data Center; use macro id for Cloud")
	}

	selectorKind := "id"
	selector := macroID
	if hash != "" {
		selectorKind = "hash"
		selector = hash
	}
	path := "/rest/api/content/" +
		url.PathEscape(contentID) +
		"/history/" +
		strconv.Itoa(in.Version) +
		"/macro/" +
		selectorKind +
		"/" +
		url.PathEscape(selector)

	data, _, err := c.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}
	var out MacroInstance
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parse macro body: %w", err)
	}
	return &out, nil
}
