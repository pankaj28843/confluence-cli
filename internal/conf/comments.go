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

// Comment is the flattened comment shape we expose to the CLI.
type Comment struct {
	ID                      string `json:"id"`
	Author                  string `json:"author,omitempty"`
	Date                    string `json:"date,omitempty"`
	Body                    string `json:"body,omitempty"`
	Location                string `json:"location,omitempty"` // inline | footer | resolved
	Resolved                bool   `json:"resolved,omitempty"`
	InlineOriginalSelection string `json:"inlineOriginalSelection,omitempty"`
}

// ListComments fetches comments under a content id, optionally filtered by
// locations (footer, inline, resolved). Locations is repeatable, so we use
// url.Values.Add.
func ListComments(ctx context.Context, c *client.Client, contentID string, locations []string, limit int) ([]Comment, error) {
	if limit <= 0 {
		limit = 100
	}
	params := url.Values{
		"limit":  {strconv.Itoa(limit)},
		"expand": {"body.view,version,extensions.inlineProperties,extensions.resolution"},
	}
	if len(locations) == 0 {
		locations = []string{"footer", "inline", "resolved"}
	}
	for _, loc := range locations {
		params.Add("location", loc)
	}
	data, _, err := c.Get(ctx, "/rest/api/content/"+url.PathEscape(contentID)+"/child/comment", params)
	if err != nil {
		return nil, err
	}

	var raw struct {
		Results []struct {
			ID      string `json:"id"`
			Type    string `json:"type,omitempty"`
			Status  string `json:"status,omitempty"`
			Version struct {
				When string `json:"when,omitempty"`
				By   struct {
					DisplayName string `json:"displayName,omitempty"`
					Username    string `json:"username,omitempty"`
					PublicName  string `json:"publicName,omitempty"`
				} `json:"by,omitempty"`
			} `json:"version,omitempty"`
			Body struct {
				View struct {
					Value string `json:"value,omitempty"`
				} `json:"view,omitempty"`
			} `json:"body,omitempty"`
			Extensions struct {
				Location         string `json:"location,omitempty"`
				InlineProperties struct {
					MarkerRef         string `json:"markerRef,omitempty"`
					OriginalSelection string `json:"originalSelection,omitempty"`
				} `json:"inlineProperties,omitempty"`
				Resolution struct {
					Status string `json:"status,omitempty"`
				} `json:"resolution,omitempty"`
			} `json:"extensions,omitempty"`
		} `json:"results"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse comments: %w", err)
	}

	out := make([]Comment, 0, len(raw.Results))
	for _, r := range raw.Results {
		author := r.Version.By.DisplayName
		if author == "" {
			author = r.Version.By.PublicName
		}
		if author == "" {
			author = r.Version.By.Username
		}
		loc := r.Extensions.Location
		if loc == "" {
			loc = "footer"
		}
		out = append(out, Comment{
			ID:                      r.ID,
			Author:                  author,
			Date:                    r.Version.When,
			Body:                    HTMLToMarkdown(strings.TrimSpace(r.Body.View.Value)),
			Location:                loc,
			Resolved:                strings.EqualFold(r.Extensions.Resolution.Status, "resolved"),
			InlineOriginalSelection: r.Extensions.InlineProperties.OriginalSelection,
		})
	}
	return out, nil
}
