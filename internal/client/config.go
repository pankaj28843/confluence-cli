// Package client is the HTTP client for Atlassian Confluence (Server/DC + Cloud).
package client

import (
	"errors"
	"os"
	"strings"
)

// Flavor distinguishes Confluence deployment types.
type Flavor int

const (
	// FlavorServer is Confluence Server / Data Center (Bearer PAT, /rest/api/...).
	FlavorServer Flavor = iota
	// FlavorCloud is Confluence Cloud (Basic email:token, /wiki/rest/api + /wiki/api/v2).
	FlavorCloud
)

func (f Flavor) String() string {
	if f == FlavorCloud {
		return "cloud"
	}
	return "server"
}

// DefaultUserAgent is sent with every HTTP request.
const DefaultUserAgent = "confluence-cli/0.1 (+https://github.com/pankaj28843/confluence-cli)"

// ErrMissingURL indicates CONFLUENCE_URL is unset.
var ErrMissingURL = errors.New("CONFLUENCE_URL is required (e.g. https://wiki.example.com or https://example.atlassian.net/wiki)")

// ErrMissingServerAuth indicates Server/DC auth is incomplete.
var ErrMissingServerAuth = errors.New("Server/DC auth requires CONFLUENCE_PAT (or CONFLUENCE_PERSONAL_ACCESS_TOKEN)")

// ErrMissingCloudAuth indicates Cloud auth is incomplete.
var ErrMissingCloudAuth = errors.New("Cloud auth requires CONFLUENCE_EMAIL and CONFLUENCE_API_TOKEN")

// Config collects env-driven configuration.
type Config struct {
	BaseURL      string
	Flavor       Flavor
	PAT          string // Server/DC
	Email        string // Cloud
	APIToken     string // Cloud
	DefaultSpace string
	Debug        bool
}

// DetectFlavor picks cloud vs server based on URL hints. Override via env.
func DetectFlavor(rawURL string) Flavor {
	u := strings.ToLower(rawURL)
	if strings.Contains(u, ".atlassian.net") {
		return FlavorCloud
	}
	// A /wiki suffix is another Cloud hint (custom domain fronting atlassian.net).
	if strings.HasSuffix(strings.TrimRight(u, "/"), "/wiki") {
		return FlavorCloud
	}
	return FlavorServer
}

// FromEnv loads Config from canonical env vars and validates per flavor.
func FromEnv() (Config, error) {
	cfg := Config{
		BaseURL:      strings.TrimRight(os.Getenv("CONFLUENCE_URL"), "/"),
		PAT:          firstNonEmpty(os.Getenv("CONFLUENCE_PAT"), os.Getenv("CONFLUENCE_PERSONAL_ACCESS_TOKEN")),
		Email:        os.Getenv("CONFLUENCE_EMAIL"),
		APIToken:     os.Getenv("CONFLUENCE_API_TOKEN"),
		DefaultSpace: os.Getenv("CONFLUENCE_DEFAULT_SPACE"),
	}
	if cfg.BaseURL == "" {
		return cfg, ErrMissingURL
	}

	switch strings.ToLower(os.Getenv("CONFLUENCE_FLAVOR")) {
	case "server", "dc", "datacenter":
		cfg.Flavor = FlavorServer
	case "cloud":
		cfg.Flavor = FlavorCloud
	default:
		cfg.Flavor = DetectFlavor(cfg.BaseURL)
	}

	switch cfg.Flavor {
	case FlavorCloud:
		if cfg.Email == "" || cfg.APIToken == "" {
			return cfg, ErrMissingCloudAuth
		}
	default:
		if cfg.PAT == "" {
			return cfg, ErrMissingServerAuth
		}
	}
	return cfg, nil
}

func firstNonEmpty(ss ...string) string {
	for _, s := range ss {
		if s != "" {
			return s
		}
	}
	return ""
}
