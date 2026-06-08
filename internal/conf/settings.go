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

// ThemeListOptions controls Cloud theme listing.
type ThemeListOptions struct {
	Start int
	Limit int
}

// Theme is a Confluence Cloud theme.
type Theme struct {
	ThemeKey    string         `json:"themeKey,omitempty"`
	Name        string         `json:"name,omitempty"`
	Description string         `json:"description,omitempty"`
	Icon        map[string]any `json:"icon,omitempty"`
	Links       map[string]any `json:"_links,omitempty"`
}

// SystemInfo is the documented Cloud system information response.
type SystemInfo struct {
	CloudID         string `json:"cloudId,omitempty"`
	CommitHash      string `json:"commitHash,omitempty"`
	BaseURL         string `json:"baseUrl,omitempty"`
	FallbackBaseURL string `json:"fallbackBaseUrl,omitempty"`
	Edition         string `json:"edition,omitempty"`
	SiteTitle       string `json:"siteTitle,omitempty"`
	DefaultLocale   string `json:"defaultLocale,omitempty"`
	DefaultTimeZone string `json:"defaultTimeZone,omitempty"`
	MicrosPerimeter string `json:"microsPerimeter,omitempty"`
}

// ColorSetting is one documented look-and-feel color entry.
type ColorSetting struct {
	Color string `json:"color,omitempty"`
}

// LookAndFeelSettings is the documented Cloud look-and-feel response.
type LookAndFeelSettings struct {
	Headings           ColorSetting   `json:"headings,omitempty"`
	Links              ColorSetting   `json:"links,omitempty"`
	Menus              map[string]any `json:"menus,omitempty"`
	Header             map[string]any `json:"header,omitempty"`
	HorizontalHeader   map[string]any `json:"horizontalHeader,omitempty"`
	Content            map[string]any `json:"content,omitempty"`
	BordersAndDividers ColorSetting   `json:"bordersAndDividers,omitempty"`
	SpaceReference     map[string]any `json:"spaceReference,omitempty"`
}

// SpaceSettings is the documented Cloud space settings response.
type SpaceSettings struct {
	RouteOverrideEnabled bool `json:"routeOverrideEnabled"`
	Editor               struct {
		Page     string `json:"page,omitempty"`
		Blogpost string `json:"blogpost,omitempty"`
		Default  string `json:"default,omitempty"`
	} `json:"editor,omitempty"`
	SpaceKey string         `json:"spaceKey,omitempty"`
	Links    map[string]any `json:"_links,omitempty"`
}

// ListThemes lists available Cloud themes.
func ListThemes(ctx context.Context, c *client.Client, opts ThemeListOptions) ([]Theme, error) {
	if !isCloud(c) {
		return nil, fmt.Errorf("themes are only supported on Confluence Cloud")
	}
	limit := clampLimit(opts.Limit)
	start := opts.Start
	if start < 0 {
		start = 0
	}
	path := "/rest/api/settings/theme"
	params := url.Values{
		"start": {strconv.Itoa(start)},
		"limit": {strconv.Itoa(limit)},
	}
	out := make([]Theme, 0, limit)
	for len(out) < limit {
		data, _, err := c.Get(ctx, path, params)
		if err != nil {
			return nil, err
		}
		var page struct {
			Results []Theme `json:"results"`
			Links   struct {
				Next string `json:"next,omitempty"`
			} `json:"_links"`
		}
		if err := json.Unmarshal(data, &page); err != nil {
			return nil, fmt.Errorf("parse themes: %w", err)
		}
		for _, theme := range page.Results {
			out = append(out, theme)
			if len(out) == limit {
				break
			}
		}
		if page.Links.Next == "" || len(page.Results) == 0 || len(out) == limit {
			break
		}
		nextPath, nextParams, err := nextPageRequest(c, page.Links.Next)
		if err != nil {
			return nil, err
		}
		path, params = nextPath, nextParams
	}
	return out, nil
}

// GetGlobalTheme returns the selected global Cloud theme.
func GetGlobalTheme(ctx context.Context, c *client.Client) (*Theme, error) {
	if !isCloud(c) {
		return nil, fmt.Errorf("themes are only supported on Confluence Cloud")
	}
	return getTheme(ctx, c, "/rest/api/settings/theme/selected")
}

// GetTheme returns one Cloud theme by key.
func GetTheme(ctx context.Context, c *client.Client, themeKey string) (*Theme, error) {
	if !isCloud(c) {
		return nil, fmt.Errorf("themes are only supported on Confluence Cloud")
	}
	themeKey = strings.TrimSpace(themeKey)
	if themeKey == "" {
		return nil, fmt.Errorf("GetTheme: theme key is required")
	}
	return getTheme(ctx, c, "/rest/api/settings/theme/"+url.PathEscape(themeKey))
}

// GetSpaceTheme returns the selected theme for one Cloud space.
func GetSpaceTheme(ctx context.Context, c *client.Client, spaceKey string) (*Theme, error) {
	if !isCloud(c) {
		return nil, fmt.Errorf("themes are only supported on Confluence Cloud")
	}
	spaceKey = strings.TrimSpace(spaceKey)
	if spaceKey == "" {
		return nil, fmt.Errorf("GetSpaceTheme: space key is required")
	}
	return getTheme(ctx, c, "/rest/api/space/"+url.PathEscape(spaceKey)+"/theme")
}

// GetSystemInfo returns Cloud system information.
func GetSystemInfo(ctx context.Context, c *client.Client) (*SystemInfo, error) {
	if !isCloud(c) {
		return nil, fmt.Errorf("settings are only supported on Confluence Cloud")
	}
	data, _, err := c.Get(ctx, "/rest/api/settings/systemInfo", nil)
	if err != nil {
		return nil, err
	}
	var out SystemInfo
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parse system info: %w", err)
	}
	return &out, nil
}

// GetLookAndFeelSettings returns global or space-specific Cloud look-and-feel settings.
func GetLookAndFeelSettings(ctx context.Context, c *client.Client, spaceKey string) (*LookAndFeelSettings, error) {
	if !isCloud(c) {
		return nil, fmt.Errorf("settings are only supported on Confluence Cloud")
	}
	params := url.Values{}
	if spaceKey = strings.TrimSpace(spaceKey); spaceKey != "" {
		params.Set("spaceKey", spaceKey)
	}
	data, _, err := c.Get(ctx, "/rest/api/settings/lookandfeel", params)
	if err != nil {
		return nil, err
	}
	var out LookAndFeelSettings
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parse look and feel settings: %w", err)
	}
	return &out, nil
}

// GetSpaceSettings returns Cloud settings for one space.
func GetSpaceSettings(ctx context.Context, c *client.Client, spaceKey string) (*SpaceSettings, error) {
	if !isCloud(c) {
		return nil, fmt.Errorf("settings are only supported on Confluence Cloud")
	}
	spaceKey = strings.TrimSpace(spaceKey)
	if spaceKey == "" {
		return nil, fmt.Errorf("GetSpaceSettings: space key is required")
	}
	data, _, err := c.Get(ctx, "/rest/api/space/"+url.PathEscape(spaceKey)+"/settings", nil)
	if err != nil {
		return nil, err
	}
	var out SpaceSettings
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parse space settings: %w", err)
	}
	return &out, nil
}

func getTheme(ctx context.Context, c *client.Client, path string) (*Theme, error) {
	data, _, err := c.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}
	var out Theme
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parse theme: %w", err)
	}
	return &out, nil
}
