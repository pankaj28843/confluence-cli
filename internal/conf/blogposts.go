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

// BlogPostListOptions controls typed blogpost listing.
type BlogPostListOptions struct {
	SpaceKey   string
	LabelID    string
	Title      string
	Status     string
	PostingDay string
	Limit      int
}

// BlogPostInput controls typed blogpost create/update operations.
type BlogPostInput struct {
	ID            string
	SpaceKey      string
	Title         string
	BodyFormat    string
	BodyValue     string
	Status        string
	VersionNumber int
}

// BlogPostDeleteOptions controls blogpost trash/delete behavior.
type BlogPostDeleteOptions struct {
	Purge bool
	Draft bool
}

// ListBlogPosts lists blog posts with flavor-appropriate endpoints.
func ListBlogPosts(ctx context.Context, c *client.Client, opts BlogPostListOptions) ([]Content, error) {
	opts.Limit = clampLimit(opts.Limit)
	if isCloud(c) {
		return listBlogPostsCloudV2(ctx, c, opts)
	}
	return listBlogPostsServerV1(ctx, c, opts)
}

func listBlogPostsCloudV2(ctx context.Context, c *client.Client, opts BlogPostListOptions) ([]Content, error) {
	if opts.PostingDay != "" {
		return nil, fmt.Errorf("ListBlogPosts: --posting-day is only supported on Server/Data Center")
	}
	if opts.LabelID != "" && opts.Title != "" {
		return nil, fmt.Errorf("ListBlogPosts: --title is not supported with --label-id on Confluence Cloud")
	}

	path := "/api/v2/blogposts"
	params := url.Values{
		"limit":       {strconv.Itoa(opts.Limit)},
		"body-format": {"storage"},
	}
	var spaceKey string
	var spaceID string
	if opts.SpaceKey != "" {
		space, err := getCloudSpaceV2(ctx, c, opts.SpaceKey)
		if err != nil {
			return nil, err
		}
		if space.ID == "" {
			return nil, fmt.Errorf("ListBlogPosts: cloud space %q has no id", opts.SpaceKey)
		}
		spaceKey = space.Key
		spaceID = space.ID
		path = "/api/v2/spaces/" + url.PathEscape(space.ID) + "/blogposts"
	}
	if opts.LabelID != "" {
		path = "/api/v2/labels/" + url.PathEscape(opts.LabelID) + "/blogposts"
		if spaceID != "" {
			params.Set("space-id", spaceID)
		}
	}
	if opts.Title != "" {
		params.Set("title", opts.Title)
	}
	addCSVValues(params, "status", opts.Status)
	return collectBlogPostPages(ctx, c, path, params, opts.Limit, spaceKey)
}

func collectBlogPostPages(ctx context.Context, c *client.Client, path string, params url.Values, limit int, spaceKey string) ([]Content, error) {
	out := make([]Content, 0, limit)
	for len(out) < limit {
		data, headers, err := c.Get(ctx, path, params)
		if err != nil {
			return nil, err
		}
		var page struct {
			Results []Content `json:"results"`
			Links   struct {
				Next string `json:"next,omitempty"`
			} `json:"_links"`
		}
		if err := json.Unmarshal(data, &page); err != nil {
			return nil, fmt.Errorf("parse blogposts: %w", err)
		}
		for i := range page.Results {
			normalizeBlogPost(&page.Results[i], spaceKey)
			out = append(out, page.Results[i])
			if len(out) == limit {
				break
			}
		}
		next := nextPageURL(headers, page.Links.Next)
		if next == "" || len(page.Results) == 0 || len(out) == limit {
			break
		}
		nextPath, nextParams, err := nextPageRequest(c, next)
		if err != nil {
			return nil, err
		}
		path, params = nextPath, nextParams
	}
	return out, nil
}

func listBlogPostsServerV1(ctx context.Context, c *client.Client, opts BlogPostListOptions) ([]Content, error) {
	if opts.LabelID != "" {
		return nil, fmt.Errorf("ListBlogPosts: --label-id is only supported on Confluence Cloud")
	}
	path := "/rest/api/content"
	params := url.Values{
		"type":   {"blogpost"},
		"limit":  {strconv.Itoa(opts.Limit)},
		"expand": {DefaultExpand},
	}
	if opts.SpaceKey != "" {
		params.Set("spaceKey", opts.SpaceKey)
	}
	if opts.PostingDay != "" {
		params.Set("postingDay", opts.PostingDay)
	}
	if opts.Title != "" {
		params.Set("title", opts.Title)
	}
	if opts.Status != "" {
		params.Set("status", opts.Status)
	}
	return collectBlogPostPages(ctx, c, path, params, opts.Limit, opts.SpaceKey)
}

// GetBlogPost fetches a blog post by id.
func GetBlogPost(ctx context.Context, c *client.Client, id string) (*Content, error) {
	if id == "" {
		return nil, fmt.Errorf("GetBlogPost: ID is required")
	}
	if isCloud(c) {
		params := url.Values{
			"body-format":        {"storage"},
			"include-labels":     {"true"},
			"include-properties": {"true"},
			"include-version":    {"true"},
		}
		data, _, err := c.Get(ctx, "/api/v2/blogposts/"+url.PathEscape(id), params)
		if err != nil {
			return nil, err
		}
		out, err := parseContentResponse(data, "blogpost")
		if err != nil {
			return nil, err
		}
		normalizeBlogPost(out, "")
		return out, nil
	}
	out, err := GetContent(ctx, c, id, DefaultExpand)
	if err != nil {
		return nil, err
	}
	normalizeBlogPost(out, "")
	return out, nil
}

// CreateBlogPost creates a blog post.
func CreateBlogPost(ctx context.Context, c *client.Client, in BlogPostInput) (*Content, error) {
	if in.SpaceKey == "" {
		return nil, fmt.Errorf("CreateBlogPost: space key is required")
	}
	if in.Title == "" {
		return nil, fmt.Errorf("CreateBlogPost: title is required")
	}
	if in.BodyValue == "" {
		return nil, fmt.Errorf("CreateBlogPost: body is required")
	}
	if in.BodyFormat == "" {
		in.BodyFormat = "storage"
	}
	if in.Status == "" {
		in.Status = "current"
	}
	if isCloud(c) {
		return createBlogPostCloudV2(ctx, c, in)
	}
	if in.Status != "current" {
		return nil, fmt.Errorf("CreateBlogPost: draft creation is only supported on Confluence Cloud")
	}
	return CreatePage(ctx, c, CreatePageInput{
		SpaceKey:   in.SpaceKey,
		Title:      in.Title,
		BodyFormat: in.BodyFormat,
		BodyValue:  in.BodyValue,
		Type:       "blogpost",
	})
}

func createBlogPostCloudV2(ctx context.Context, c *client.Client, in BlogPostInput) (*Content, error) {
	if in.BodyFormat != "storage" {
		return nil, fmt.Errorf("CreateBlogPost: Confluence Cloud v2 blog posts only support storage body format")
	}
	space, err := getCloudSpaceV2(ctx, c, in.SpaceKey)
	if err != nil {
		return nil, err
	}
	if space.ID == "" {
		return nil, fmt.Errorf("CreateBlogPost: cloud space %q has no id", in.SpaceKey)
	}
	body := map[string]any{
		"spaceId": space.ID,
		"status":  in.Status,
		"title":   in.Title,
		"body": map[string]any{
			"representation": in.BodyFormat,
			"value":          in.BodyValue,
		},
	}
	raw, _, err := c.Post(ctx, "/api/v2/blogposts", nil, body)
	if err != nil {
		return nil, err
	}
	out, err := parseContentResponse(raw, "create blogpost")
	if err != nil {
		return nil, err
	}
	normalizeBlogPost(out, space.Key)
	return out, nil
}

// UpdateBlogPost updates a blog post.
func UpdateBlogPost(ctx context.Context, c *client.Client, in BlogPostInput) (*Content, error) {
	if in.ID == "" {
		return nil, fmt.Errorf("UpdateBlogPost: ID is required")
	}
	if in.Title == "" {
		return nil, fmt.Errorf("UpdateBlogPost: title is required")
	}
	if in.BodyValue == "" {
		return nil, fmt.Errorf("UpdateBlogPost: body is required")
	}
	if in.BodyFormat == "" {
		in.BodyFormat = "storage"
	}
	if in.Status == "" {
		in.Status = "current"
	}
	if isCloud(c) {
		return updateBlogPostCloudV2(ctx, c, in)
	}
	return UpdatePage(ctx, c, UpdatePageInput{
		ID:            in.ID,
		Title:         in.Title,
		BodyFormat:    in.BodyFormat,
		BodyValue:     in.BodyValue,
		VersionNumber: in.VersionNumber,
		Type:          "blogpost",
	})
}

func updateBlogPostCloudV2(ctx context.Context, c *client.Client, in BlogPostInput) (*Content, error) {
	if in.BodyFormat != "storage" {
		return nil, fmt.Errorf("UpdateBlogPost: Confluence Cloud v2 blog posts only support storage body format")
	}
	body := map[string]any{
		"id":     in.ID,
		"status": in.Status,
		"title":  in.Title,
		"body": map[string]any{
			"representation": in.BodyFormat,
			"value":          in.BodyValue,
		},
		"version": map[string]any{
			"number": in.VersionNumber + 1,
		},
	}
	raw, _, err := c.Put(ctx, "/api/v2/blogposts/"+url.PathEscape(in.ID), nil, body)
	if err != nil {
		return nil, err
	}
	out, err := parseContentResponse(raw, "update blogpost")
	if err != nil {
		return nil, err
	}
	normalizeBlogPost(out, "")
	return out, nil
}

// DeleteBlogPost trashes or purges a blog post using the flavor-appropriate endpoint.
func DeleteBlogPost(ctx context.Context, c *client.Client, id string, opts BlogPostDeleteOptions) error {
	if id == "" {
		return fmt.Errorf("DeleteBlogPost: ID is required")
	}
	if isCloud(c) {
		params := url.Values{}
		if opts.Purge {
			params.Set("purge", "true")
		}
		if opts.Draft {
			params.Set("draft", "true")
		}
		_, _, err := c.Delete(ctx, "/api/v2/blogposts/"+url.PathEscape(id), params)
		return err
	}
	if opts.Draft {
		return fmt.Errorf("DeleteBlogPost: draft deletion is only supported on Confluence Cloud")
	}
	params := url.Values{}
	if opts.Purge {
		params.Set("status", "trashed")
	}
	_, _, err := c.Delete(ctx, "/rest/api/content/"+url.PathEscape(id), params)
	return err
}

func normalizeBlogPost(out *Content, spaceKey string) {
	if out.Type == "" {
		out.Type = "blogpost"
	}
	if out.Space.Key == "" {
		out.Space.Key = spaceKey
	}
	if out.Body.Storage.Representation == "" && out.Body.Storage.Value != "" {
		out.Body.Storage.Representation = "storage"
	}
}

func clampLimit(limit int) int {
	if limit <= 0 {
		return 25
	}
	if limit > 200 {
		return 200
	}
	return limit
}

func addCSVValues(params url.Values, key, raw string) {
	for _, value := range strings.Split(raw, ",") {
		value = strings.TrimSpace(value)
		if value != "" {
			params.Add(key, value)
		}
	}
}
