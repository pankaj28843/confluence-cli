package client

import (
	"errors"
	"testing"
)

func TestFromEnvMissingURL(t *testing.T) {
	t.Setenv("CONFLUENCE_URL", "")
	t.Setenv("CONFLUENCE_PAT", "x")
	if _, err := FromEnv(); !errors.Is(err, ErrMissingURL) {
		t.Fatalf("want ErrMissingURL, got %v", err)
	}
}

func TestFromEnvServerMissingPAT(t *testing.T) {
	t.Setenv("CONFLUENCE_URL", "https://wiki.example.com")
	t.Setenv("CONFLUENCE_PAT", "")
	t.Setenv("CONFLUENCE_PERSONAL_ACCESS_TOKEN", "")
	if _, err := FromEnv(); !errors.Is(err, ErrMissingServerAuth) {
		t.Fatalf("want ErrMissingServerAuth, got %v", err)
	}
}

func TestFromEnvCloudMissingAuth(t *testing.T) {
	t.Setenv("CONFLUENCE_URL", "https://example.atlassian.net/wiki")
	t.Setenv("CONFLUENCE_EMAIL", "")
	t.Setenv("CONFLUENCE_API_TOKEN", "")
	if _, err := FromEnv(); !errors.Is(err, ErrMissingCloudAuth) {
		t.Fatalf("want ErrMissingCloudAuth, got %v", err)
	}
}

func TestFromEnvPATAlias(t *testing.T) {
	t.Setenv("CONFLUENCE_URL", "https://wiki.example.com")
	t.Setenv("CONFLUENCE_PAT", "")
	t.Setenv("CONFLUENCE_PERSONAL_ACCESS_TOKEN", "alias-pat")
	cfg, err := FromEnv()
	if err != nil {
		t.Fatalf("FromEnv: %v", err)
	}
	if cfg.PAT != "alias-pat" || cfg.Flavor != FlavorServer {
		t.Fatalf("alias or flavor wrong: %+v", cfg)
	}
}

func TestFromEnvFlavorOverride(t *testing.T) {
	t.Setenv("CONFLUENCE_URL", "https://wiki.example.com")
	t.Setenv("CONFLUENCE_FLAVOR", "cloud")
	t.Setenv("CONFLUENCE_EMAIL", "e@example.com")
	t.Setenv("CONFLUENCE_API_TOKEN", "t")
	cfg, err := FromEnv()
	if err != nil {
		t.Fatalf("FromEnv: %v", err)
	}
	if cfg.Flavor != FlavorCloud {
		t.Fatalf("flavor override ignored: %v", cfg.Flavor)
	}
}
