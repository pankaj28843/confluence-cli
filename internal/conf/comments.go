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
	PageID                  string `json:"pageId,omitempty"`
	BlogPostID              string `json:"blogPostId,omitempty"`
	ParentCommentID         string `json:"parentCommentId,omitempty"`
	Author                  string `json:"author,omitempty"`
	Date                    string `json:"date,omitempty"`
	Body                    string `json:"body,omitempty"`
	Location                string `json:"location,omitempty"` // inline | footer | resolved
	Resolved                bool   `json:"resolved,omitempty"`
	VersionNumber           int    `json:"versionNumber,omitempty"`
	InlineOriginalSelection string `json:"inlineOriginalSelection,omitempty"`
}

// CommentInput controls typed footer comment create/update operations.
type CommentInput struct {
	ID              string
	PageID          string
	BlogPostID      string
	ParentCommentID string
	BodyFormat      string
	BodyValue       string
	VersionNumber   int
}

// ListComments fetches comments under a content id, optionally filtered by
// locations (footer, inline, resolved). Locations is repeatable, so we use
// url.Values.Add.
func ListComments(ctx context.Context, c *client.Client, contentID string, locations []string, limit int) ([]Comment, error) {
	if isCloud(c) {
		return listCommentsCloudV2(ctx, c, contentID, locations, limit)
	}
	return listCommentsServerV1(ctx, c, contentID, locations, limit)
}

func listCommentsServerV1(ctx context.Context, c *client.Client, contentID string, locations []string, limit int) ([]Comment, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 200 {
		limit = 200
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
				Number int    `json:"number,omitempty"`
				When   string `json:"when,omitempty"`
				By     struct {
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
			VersionNumber:           r.Version.Number,
			InlineOriginalSelection: r.Extensions.InlineProperties.OriginalSelection,
		})
	}
	return out, nil
}

type cloudCommentV2 struct {
	ID               string `json:"id"`
	Status           string `json:"status,omitempty"`
	Title            string `json:"title,omitempty"`
	PageID           string `json:"pageId,omitempty"`
	BlogPostID       string `json:"blogPostId,omitempty"`
	AttachmentID     string `json:"attachmentId,omitempty"`
	ParentCommentID  string `json:"parentCommentId,omitempty"`
	ResolutionStatus string `json:"resolutionStatus,omitempty"`
	Version          struct {
		Number    int    `json:"number,omitempty"`
		CreatedAt string `json:"createdAt,omitempty"`
		AuthorID  string `json:"authorId,omitempty"`
	} `json:"version,omitempty"`
	Body struct {
		Storage struct {
			Value string `json:"value,omitempty"`
		} `json:"storage,omitempty"`
	} `json:"body,omitempty"`
	Properties struct {
		InlineMarkerRef         string `json:"inlineMarkerRef,omitempty"`
		InlineOriginalSelection string `json:"inlineOriginalSelection,omitempty"`
	} `json:"properties,omitempty"`
	Links struct {
		WebUI string `json:"webui,omitempty"`
	} `json:"_links,omitempty"`
}

func listCommentsCloudV2(ctx context.Context, c *client.Client, contentID string, locations []string, limit int) ([]Comment, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 200 {
		limit = 200
	}

	wantFooter, wantInline, wantResolved := commentLocationSet(locations)
	out := make([]Comment, 0, limit)
	if wantFooter {
		params := cloudCommentParams(limit)
		comments, err := listCloudCommentPageV2(ctx, c, "/api/v2/pages/"+url.PathEscape(contentID)+"/footer-comments", params, limit-len(out), "footer")
		if err != nil {
			return nil, err
		}
		out = append(out, comments...)
	}
	if len(out) < limit && wantInline {
		params := cloudCommentParams(limit - len(out))
		if !wantResolved {
			params.Add("resolution-status", "open")
		}
		comments, err := listCloudCommentPageV2(ctx, c, "/api/v2/pages/"+url.PathEscape(contentID)+"/inline-comments", params, limit-len(out), "inline")
		if err != nil {
			return nil, err
		}
		out = append(out, comments...)
	}
	if len(out) < limit && wantResolved && !wantInline {
		params := cloudCommentParams(limit - len(out))
		params.Add("resolution-status", "resolved")
		comments, err := listCloudCommentPageV2(ctx, c, "/api/v2/pages/"+url.PathEscape(contentID)+"/inline-comments", params, limit-len(out), "resolved")
		if err != nil {
			return nil, err
		}
		out = append(out, comments...)
	}
	return out, nil
}

func commentLocationSet(locations []string) (wantFooter, wantInline, wantResolved bool) {
	if len(locations) == 0 {
		return true, true, true
	}
	for _, loc := range locations {
		switch strings.ToLower(strings.TrimSpace(loc)) {
		case "footer":
			wantFooter = true
		case "inline":
			wantInline = true
		case "resolved":
			wantResolved = true
		}
	}
	return wantFooter, wantInline, wantResolved
}

func cloudCommentParams(limit int) url.Values {
	if limit <= 0 {
		limit = 25
	}
	if limit > 200 {
		limit = 200
	}
	return url.Values{
		"body-format": {"STORAGE"},
		"limit":       {strconv.Itoa(limit)},
	}
}

func listCloudCommentPageV2(ctx context.Context, c *client.Client, path string, params url.Values, limit int, location string) ([]Comment, error) {
	out := make([]Comment, 0, limit)
	for len(out) < limit {
		data, hdrs, err := c.Get(ctx, path, params)
		if err != nil {
			return nil, err
		}
		var page struct {
			Results []cloudCommentV2 `json:"results"`
			Links   struct {
				Next string `json:"next,omitempty"`
			} `json:"_links"`
		}
		if err := json.Unmarshal(data, &page); err != nil {
			return nil, fmt.Errorf("parse comments: %w", err)
		}
		for _, raw := range page.Results {
			out = append(out, normalizeCloudCommentV2(raw, location))
			if len(out) == limit {
				break
			}
		}
		next := nextPageURL(hdrs, page.Links.Next)
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

func normalizeCloudCommentV2(raw cloudCommentV2, location string) Comment {
	resolved := strings.EqualFold(raw.ResolutionStatus, "resolved")
	if resolved {
		location = "resolved"
	}
	return Comment{
		ID:                      raw.ID,
		PageID:                  raw.PageID,
		BlogPostID:              raw.BlogPostID,
		ParentCommentID:         raw.ParentCommentID,
		Author:                  raw.Version.AuthorID,
		Date:                    raw.Version.CreatedAt,
		Body:                    HTMLToMarkdown(strings.TrimSpace(raw.Body.Storage.Value)),
		Location:                location,
		Resolved:                resolved,
		VersionNumber:           raw.Version.Number,
		InlineOriginalSelection: raw.Properties.InlineOriginalSelection,
	}
}

// CreateFooterComment creates a footer comment on a page/blogpost or replies to
// an existing footer comment on Cloud.
func CreateFooterComment(ctx context.Context, c *client.Client, in CommentInput) (*Comment, error) {
	if in.BodyValue == "" {
		return nil, fmt.Errorf("CreateFooterComment: body is required")
	}
	if in.BodyFormat == "" {
		in.BodyFormat = "storage"
	}
	if isCloud(c) {
		if err := validateCloudFooterCommentTarget(in); err != nil {
			return nil, err
		}
		return createFooterCommentCloudV2(ctx, c, in)
	}
	if in.ParentCommentID != "" {
		return nil, fmt.Errorf("CreateFooterComment: replies are only supported on Confluence Cloud")
	}
	if err := validateServerFooterCommentTarget(in); err != nil {
		return nil, err
	}
	return createFooterCommentServerV1(ctx, c, in)
}

func createFooterCommentCloudV2(ctx context.Context, c *client.Client, in CommentInput) (*Comment, error) {
	body := map[string]any{
		"body": cloudCommentBody(in.BodyFormat, in.BodyValue),
	}
	if in.PageID != "" {
		body["pageId"] = in.PageID
	}
	if in.BlogPostID != "" {
		body["blogPostId"] = in.BlogPostID
	}
	if in.ParentCommentID != "" {
		body["parentCommentId"] = in.ParentCommentID
	}
	raw, _, err := c.Post(ctx, "/api/v2/footer-comments", nil, body)
	if err != nil {
		return nil, err
	}
	return parseCloudFooterComment(raw, "create footer comment")
}

func createFooterCommentServerV1(ctx context.Context, c *client.Client, in CommentInput) (*Comment, error) {
	containerID, containerType := serverCommentContainer(in)
	body := map[string]any{
		"type":      "comment",
		"container": map[string]any{"id": containerID, "type": containerType},
		"body":      serverCommentBody(in.BodyFormat, in.BodyValue),
	}
	raw, _, err := c.Post(ctx, "/rest/api/content", nil, body)
	if err != nil {
		return nil, err
	}
	return parseServerFooterComment(raw, "create footer comment")
}

// GetFooterComment fetches one footer comment with storage body and version.
func GetFooterComment(ctx context.Context, c *client.Client, id string) (*Comment, error) {
	if id == "" {
		return nil, fmt.Errorf("GetFooterComment: ID is required")
	}
	if isCloud(c) {
		params := url.Values{
			"body-format":     {"STORAGE"},
			"include-version": {"true"},
		}
		raw, _, err := c.Get(ctx, "/api/v2/footer-comments/"+url.PathEscape(id), params)
		if err != nil {
			return nil, err
		}
		return parseCloudFooterComment(raw, "footer comment")
	}
	cnt, err := GetContent(ctx, c, id, "body.storage,version")
	if err != nil {
		return nil, err
	}
	return normalizeServerComment(*cnt), nil
}

// UpdateFooterComment updates a footer comment body using version.number + 1.
func UpdateFooterComment(ctx context.Context, c *client.Client, in CommentInput) (*Comment, error) {
	if in.ID == "" {
		return nil, fmt.Errorf("UpdateFooterComment: ID is required")
	}
	if in.BodyValue == "" {
		return nil, fmt.Errorf("UpdateFooterComment: body is required")
	}
	if in.BodyFormat == "" {
		in.BodyFormat = "storage"
	}
	if isCloud(c) {
		return updateFooterCommentCloudV2(ctx, c, in)
	}
	return updateFooterCommentServerV1(ctx, c, in)
}

func updateFooterCommentCloudV2(ctx context.Context, c *client.Client, in CommentInput) (*Comment, error) {
	body := map[string]any{
		"version": map[string]any{"number": in.VersionNumber + 1},
		"body":    cloudCommentBody(in.BodyFormat, in.BodyValue),
	}
	raw, _, err := c.Put(ctx, "/api/v2/footer-comments/"+url.PathEscape(in.ID), nil, body)
	if err != nil {
		return nil, err
	}
	return parseCloudFooterComment(raw, "update footer comment")
}

func updateFooterCommentServerV1(ctx context.Context, c *client.Client, in CommentInput) (*Comment, error) {
	body := map[string]any{
		"type":    "comment",
		"version": map[string]any{"number": in.VersionNumber + 1},
		"body":    serverCommentBody(in.BodyFormat, in.BodyValue),
	}
	raw, _, err := c.Put(ctx, "/rest/api/content/"+url.PathEscape(in.ID), nil, body)
	if err != nil {
		return nil, err
	}
	return parseServerFooterComment(raw, "update footer comment")
}

// DeleteFooterComment permanently deletes a footer comment.
func DeleteFooterComment(ctx context.Context, c *client.Client, id string) error {
	if id == "" {
		return fmt.Errorf("DeleteFooterComment: ID is required")
	}
	if isCloud(c) {
		_, _, err := c.Delete(ctx, "/api/v2/footer-comments/"+url.PathEscape(id), nil)
		return err
	}
	_, _, err := c.Delete(ctx, "/rest/api/content/"+url.PathEscape(id), nil)
	return err
}

func validateCloudFooterCommentTarget(in CommentInput) error {
	count := 0
	for _, id := range []string{in.PageID, in.BlogPostID, in.ParentCommentID} {
		if id != "" {
			count++
		}
	}
	if count != 1 {
		return fmt.Errorf("CreateFooterComment: pass exactly one of --page, --blogpost, or --parent")
	}
	return nil
}

func validateServerFooterCommentTarget(in CommentInput) error {
	count := 0
	for _, id := range []string{in.PageID, in.BlogPostID} {
		if id != "" {
			count++
		}
	}
	if count != 1 {
		return fmt.Errorf("CreateFooterComment: pass exactly one of --page or --blogpost")
	}
	return nil
}

func serverCommentContainer(in CommentInput) (id, contentType string) {
	if in.BlogPostID != "" {
		return in.BlogPostID, "blogpost"
	}
	return in.PageID, "page"
}

func cloudCommentBody(format, value string) map[string]any {
	return map[string]any{
		"representation": format,
		"value":          value,
	}
}

func serverCommentBody(format, value string) map[string]any {
	return map[string]any{
		format: map[string]any{
			"value":          value,
			"representation": format,
		},
	}
}

func parseCloudFooterComment(data []byte, verb string) (*Comment, error) {
	var raw cloudCommentV2
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse %s response: %w", verb, err)
	}
	out := normalizeCloudCommentV2(raw, "footer")
	return &out, nil
}

func parseServerFooterComment(data []byte, verb string) (*Comment, error) {
	var raw Content
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse %s response: %w", verb, err)
	}
	return normalizeServerComment(raw), nil
}

func normalizeServerComment(raw Content) *Comment {
	author := raw.Version.By.DisplayName
	if author == "" {
		author = raw.Version.By.Username
	}
	body := raw.Body.Storage.Value
	if body == "" {
		body = raw.Body.View.Value
	}
	return &Comment{
		ID:            raw.ID,
		Author:        author,
		Date:          raw.Version.When,
		Body:          HTMLToMarkdown(strings.TrimSpace(body)),
		Location:      "footer",
		VersionNumber: raw.Version.Number,
	}
}
