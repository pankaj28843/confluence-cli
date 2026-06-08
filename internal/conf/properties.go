package conf

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

// Property is a Confluence content or space content-property row.
type Property struct {
	ID      string `json:"id,omitempty"`
	Key     string `json:"key"`
	Value   any    `json:"value,omitempty"`
	Version struct {
		Number    int    `json:"number,omitempty"`
		CreatedAt string `json:"createdAt,omitempty"`
		When      string `json:"when,omitempty"`
		Message   string `json:"message,omitempty"`
	} `json:"version,omitempty"`
}

// ListContentProperties lists content properties for a page/content id.
func ListContentProperties(ctx context.Context, c *client.Client, contentID, key string, limit int) ([]Property, error) {
	if contentID == "" {
		return nil, fmt.Errorf("ListContentProperties: content ID is required")
	}
	if isCloud(c) {
		return listCloudProperties(ctx, c, cloudContentPropertiesPath(contentID), key, limit)
	}
	return listServerProperties(ctx, c, serverContentPropertiesPath(contentID), key, limit)
}

// GetContentProperty returns one content property by key.
func GetContentProperty(ctx context.Context, c *client.Client, contentID, key string) (*Property, error) {
	if contentID == "" {
		return nil, fmt.Errorf("GetContentProperty: content ID is required")
	}
	if key == "" {
		return nil, fmt.Errorf("GetContentProperty: key is required")
	}
	if isCloud(c) {
		return getCloudPropertyByKey(ctx, c, cloudContentPropertiesPath(contentID), key)
	}
	return getServerProperty(ctx, c, serverContentPropertyPath(contentID, key))
}

// SetContentProperty creates or updates one content property by key.
func SetContentProperty(ctx context.Context, c *client.Client, contentID, key string, value any) (*Property, error) {
	if contentID == "" {
		return nil, fmt.Errorf("SetContentProperty: content ID is required")
	}
	if key == "" {
		return nil, fmt.Errorf("SetContentProperty: key is required")
	}
	if isCloud(c) {
		return setCloudProperty(ctx, c, cloudContentPropertiesPath(contentID), key, value)
	}
	return setServerProperty(ctx, c, serverContentPropertiesPath(contentID), serverContentPropertyPath(contentID, key), key, value)
}

// DeleteContentProperty deletes one content property by key.
func DeleteContentProperty(ctx context.Context, c *client.Client, contentID, key string) error {
	if contentID == "" {
		return fmt.Errorf("DeleteContentProperty: content ID is required")
	}
	if key == "" {
		return fmt.Errorf("DeleteContentProperty: key is required")
	}
	if isCloud(c) {
		return deleteCloudProperty(ctx, c, cloudContentPropertiesPath(contentID), key)
	}
	_, _, err := c.Delete(ctx, serverContentPropertyPath(contentID, key), nil)
	return err
}

// ListSpaceProperties lists space properties for a space key.
func ListSpaceProperties(ctx context.Context, c *client.Client, spaceKey, key string, limit int) ([]Property, error) {
	if spaceKey == "" {
		return nil, fmt.Errorf("ListSpaceProperties: space key is required")
	}
	if isCloud(c) {
		path, err := cloudSpacePropertiesPath(ctx, c, spaceKey)
		if err != nil {
			return nil, err
		}
		return listCloudProperties(ctx, c, path, key, limit)
	}
	return listServerProperties(ctx, c, serverSpacePropertiesPath(spaceKey), key, limit)
}

// GetSpaceProperty returns one space property by key.
func GetSpaceProperty(ctx context.Context, c *client.Client, spaceKey, key string) (*Property, error) {
	if spaceKey == "" {
		return nil, fmt.Errorf("GetSpaceProperty: space key is required")
	}
	if key == "" {
		return nil, fmt.Errorf("GetSpaceProperty: key is required")
	}
	if isCloud(c) {
		path, err := cloudSpacePropertiesPath(ctx, c, spaceKey)
		if err != nil {
			return nil, err
		}
		return getCloudPropertyByKey(ctx, c, path, key)
	}
	return getServerProperty(ctx, c, serverSpacePropertyPath(spaceKey, key))
}

// SetSpaceProperty creates or updates one space property by key.
func SetSpaceProperty(ctx context.Context, c *client.Client, spaceKey, key string, value any) (*Property, error) {
	if spaceKey == "" {
		return nil, fmt.Errorf("SetSpaceProperty: space key is required")
	}
	if key == "" {
		return nil, fmt.Errorf("SetSpaceProperty: key is required")
	}
	if isCloud(c) {
		path, err := cloudSpacePropertiesPath(ctx, c, spaceKey)
		if err != nil {
			return nil, err
		}
		return setCloudProperty(ctx, c, path, key, value)
	}
	return setServerProperty(ctx, c, serverSpacePropertiesPath(spaceKey), serverSpacePropertyPath(spaceKey, key), key, value)
}

// DeleteSpaceProperty deletes one space property by key.
func DeleteSpaceProperty(ctx context.Context, c *client.Client, spaceKey, key string) error {
	if spaceKey == "" {
		return fmt.Errorf("DeleteSpaceProperty: space key is required")
	}
	if key == "" {
		return fmt.Errorf("DeleteSpaceProperty: key is required")
	}
	if isCloud(c) {
		path, err := cloudSpacePropertiesPath(ctx, c, spaceKey)
		if err != nil {
			return err
		}
		return deleteCloudProperty(ctx, c, path, key)
	}
	_, _, err := c.Delete(ctx, serverSpacePropertyPath(spaceKey, key), nil)
	return err
}

func listServerProperties(ctx context.Context, c *client.Client, collectionPath, key string, limit int) ([]Property, error) {
	limit = normalizePropertyLimit(limit)
	if key != "" {
		itemPath := collectionPath + "/" + url.PathEscape(key)
		prop, err := getServerProperty(ctx, c, itemPath)
		if isStatus(err, http.StatusNotFound) {
			return []Property{}, nil
		}
		if err != nil {
			return nil, err
		}
		return []Property{*prop}, nil
	}

	out := make([]Property, 0, limit)
	start := 0
	for len(out) < limit {
		remaining := limit - len(out)
		if remaining > 200 {
			remaining = 200
		}
		params := url.Values{
			"expand": {"version"},
			"limit":  {strconv.Itoa(remaining)},
			"start":  {strconv.Itoa(start)},
		}
		data, _, err := c.Get(ctx, collectionPath, params)
		if err != nil {
			return nil, err
		}
		var page struct {
			Results []Property `json:"results"`
			Size    int        `json:"size"`
			Links   struct {
				Next string `json:"next,omitempty"`
			} `json:"_links"`
		}
		if err := json.Unmarshal(data, &page); err != nil {
			return nil, fmt.Errorf("parse properties: %w", err)
		}
		for _, prop := range page.Results {
			out = append(out, prop)
			if len(out) == limit {
				break
			}
		}
		if page.Links.Next == "" || len(page.Results) == 0 || len(out) == limit {
			break
		}
		if page.Size > 0 {
			start += page.Size
		} else {
			start += len(page.Results)
		}
	}
	return out, nil
}

func listCloudProperties(ctx context.Context, c *client.Client, path, key string, limit int) ([]Property, error) {
	limit = normalizePropertyLimit(limit)
	out := make([]Property, 0, limit)

	params := url.Values{}
	remaining := limit
	if remaining > 200 {
		remaining = 200
	}
	params.Set("limit", strconv.Itoa(remaining))
	if key != "" {
		params.Set("key", key)
	}

	for len(out) < limit {
		data, hdrs, err := c.Get(ctx, path, params)
		if err != nil {
			return nil, err
		}
		var page struct {
			Results []Property `json:"results"`
			Links   struct {
				Next string `json:"next,omitempty"`
			} `json:"_links"`
		}
		if err := json.Unmarshal(data, &page); err != nil {
			return nil, fmt.Errorf("parse properties: %w", err)
		}
		for _, prop := range page.Results {
			out = append(out, prop)
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

func getServerProperty(ctx context.Context, c *client.Client, path string) (*Property, error) {
	data, _, err := c.Get(ctx, path, url.Values{"expand": {"version"}})
	if err != nil {
		return nil, err
	}
	return parseProperty(data, "property")
}

func getCloudPropertyByKey(ctx context.Context, c *client.Client, collectionPath, key string) (*Property, error) {
	prop, err := resolveCloudPropertyByKey(ctx, c, collectionPath, key)
	if err != nil {
		return nil, err
	}
	if prop == nil {
		return nil, fmt.Errorf("property %q not found", key)
	}
	return getCloudPropertyByID(ctx, c, collectionPath, prop.ID)
}

func getCloudPropertyByID(ctx context.Context, c *client.Client, collectionPath, id string) (*Property, error) {
	if id == "" {
		return nil, fmt.Errorf("cloud property id is required")
	}
	data, _, err := c.Get(ctx, collectionPath+"/"+url.PathEscape(id), nil)
	if err != nil {
		return nil, err
	}
	return parseProperty(data, "property")
}

func setServerProperty(ctx context.Context, c *client.Client, collectionPath, itemPath, key string, value any) (*Property, error) {
	existing, err := getServerProperty(ctx, c, itemPath)
	if err != nil && !isStatus(err, http.StatusNotFound) {
		return nil, err
	}

	body := map[string]any{
		"key":   key,
		"value": value,
	}
	if existing == nil {
		data, _, err := c.Post(ctx, collectionPath, nil, body)
		if err != nil {
			return nil, err
		}
		return parseProperty(data, "property")
	}

	if existing.ID != "" {
		body["id"] = existing.ID
	}
	body["version"] = map[string]any{"number": nextPropertyVersion(existing.Version.Number)}
	data, _, err := c.Put(ctx, itemPath, nil, body)
	if err != nil {
		return nil, err
	}
	return parseProperty(data, "property")
}

func setCloudProperty(ctx context.Context, c *client.Client, collectionPath, key string, value any) (*Property, error) {
	existing, err := resolveCloudPropertyByKey(ctx, c, collectionPath, key)
	if err != nil {
		return nil, err
	}

	body := map[string]any{
		"key":   key,
		"value": value,
	}
	if existing == nil {
		data, _, err := c.Post(ctx, collectionPath, nil, body)
		if err != nil {
			return nil, err
		}
		return parseProperty(data, "property")
	}
	if existing.ID == "" {
		return nil, fmt.Errorf("property %q has no id", key)
	}
	body["version"] = map[string]any{"number": nextPropertyVersion(existing.Version.Number)}
	data, _, err := c.Put(ctx, collectionPath+"/"+url.PathEscape(existing.ID), nil, body)
	if err != nil {
		return nil, err
	}
	return parseProperty(data, "property")
}

func deleteCloudProperty(ctx context.Context, c *client.Client, collectionPath, key string) error {
	prop, err := resolveCloudPropertyByKey(ctx, c, collectionPath, key)
	if err != nil {
		return err
	}
	if prop == nil {
		return fmt.Errorf("property %q not found", key)
	}
	if prop.ID == "" {
		return fmt.Errorf("property %q has no id", key)
	}
	_, _, err = c.Delete(ctx, collectionPath+"/"+url.PathEscape(prop.ID), nil)
	return err
}

func resolveCloudPropertyByKey(ctx context.Context, c *client.Client, collectionPath, key string) (*Property, error) {
	props, err := listCloudProperties(ctx, c, collectionPath, key, 1)
	if err != nil {
		return nil, err
	}
	if len(props) == 0 {
		return nil, nil
	}
	for i := range props {
		if props[i].Key == key {
			return &props[i], nil
		}
	}
	return &props[0], nil
}

func parseProperty(data []byte, what string) (*Property, error) {
	var out Property
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parse %s: %w", what, err)
	}
	return &out, nil
}

func cloudSpacePropertiesPath(ctx context.Context, c *client.Client, spaceKey string) (string, error) {
	space, err := getCloudSpaceV2(ctx, c, spaceKey)
	if err != nil {
		return "", err
	}
	if space.ID == "" {
		return "", fmt.Errorf("space %q has no id", spaceKey)
	}
	return "/api/v2/spaces/" + url.PathEscape(space.ID) + "/properties", nil
}

func serverContentPropertiesPath(contentID string) string {
	return "/rest/api/content/" + url.PathEscape(contentID) + "/property"
}

func serverContentPropertyPath(contentID, key string) string {
	return serverContentPropertiesPath(contentID) + "/" + url.PathEscape(key)
}

func cloudContentPropertiesPath(contentID string) string {
	return "/api/v2/pages/" + url.PathEscape(contentID) + "/properties"
}

func serverSpacePropertiesPath(spaceKey string) string {
	return "/rest/api/space/" + url.PathEscape(spaceKey) + "/property"
}

func serverSpacePropertyPath(spaceKey, key string) string {
	return serverSpacePropertiesPath(spaceKey) + "/" + url.PathEscape(key)
}

func normalizePropertyLimit(limit int) int {
	if limit <= 0 {
		return 25
	}
	if limit > 200 {
		return 200
	}
	return limit
}

func nextPropertyVersion(current int) int {
	if current <= 0 {
		return 1
	}
	return current + 1
}

func isStatus(err error, status int) bool {
	var apiErr *client.APIError
	return errors.As(err, &apiErr) && apiErr.Status == status
}
