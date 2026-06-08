package conf

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

func testClient(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *client.Client) {
	t.Helper()
	return testClientWithFlavor(t, client.FlavorServer, handler)
}

func testClientWithFlavor(t *testing.T, flavor client.Flavor, handler http.HandlerFunc) (*httptest.Server, *client.Client) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	baseURL := srv.URL
	cfg := client.Config{BaseURL: baseURL, Flavor: flavor, PAT: "x"}
	if flavor == client.FlavorCloud {
		cfg.BaseURL = srv.URL + "/wiki"
		cfg.Email = "x@example.com"
		cfg.APIToken = "cloud-token"
	}
	c, err := client.New(cfg)
	if err != nil {
		t.Fatalf("client.New: %v", err)
	}
	c.HTTPClient = &http.Client{Timeout: 5 * time.Second}
	return srv, c
}

func TestListSpacesServerUsesV1SpaceEndpoint(t *testing.T) {
	_, c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/rest/api/space") {
			t.Fatalf("path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("start"); got != "0" {
			t.Fatalf("start: %q", got)
		}
		if got := r.URL.Query().Get("limit"); got != "5" {
			t.Fatalf("limit: %q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"results": []Space{{Key: "ENG", Name: "Engineering"}, {Key: "OPS", Name: "Operations"}},
			"size":    2,
			"_links":  map[string]any{},
		})
	})
	got, err := ListSpaces(context.Background(), c, SpaceFilter{Limit: 5})
	if err != nil || len(got) != 2 {
		t.Fatalf("ListSpaces: %v %+v", err, got)
	}
}

func TestListSpacesCloudUsesV2SpaceEndpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/wiki/api/v2/spaces" {
			t.Fatalf("path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("start"); got != "" {
			t.Fatalf("cloud v2 should not send start, got %q", got)
		}
		if got := r.URL.Query().Get("limit"); got != "5" {
			t.Fatalf("limit: %q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"results": []map[string]any{{
				"id":     "1001",
				"key":    "ENG",
				"name":   "Engineering",
				"type":   "global",
				"status": "current",
				"description": map[string]any{
					"plain": map[string]any{"value": "Engineering docs"},
				},
				"_links": map[string]any{"webui": "/spaces/ENG"},
			}},
			"_links": map[string]any{},
		})
	})

	got, err := ListSpaces(context.Background(), c, SpaceFilter{Limit: 5})
	if err != nil {
		t.Fatalf("ListSpaces: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len(spaces) = %d, want 1", len(got))
	}
	if got[0].Key != "ENG" || got[0].Name != "Engineering" || got[0].Status != "current" {
		t.Fatalf("unexpected space: %+v", got[0])
	}
}

func TestListSpacesCloudFollowsV2Pagination(t *testing.T) {
	requests := make([]string, 0, 2)
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.URL.RequestURI())
		switch len(requests) {
		case 1:
			if r.URL.Path != "/wiki/api/v2/spaces" {
				t.Fatalf("first path: %s", r.URL.Path)
			}
			if got := r.URL.Query().Get("limit"); got != "2" {
				t.Fatalf("first limit: %q", got)
			}
			if got := r.URL.Query().Get("cursor"); got != "" {
				t.Fatalf("first cursor: %q", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{
					"id":     "1001",
					"key":    "ENG",
					"name":   "Engineering",
					"type":   "global",
					"status": "current",
				}},
				"_links": map[string]any{"next": "/api/v2/spaces?cursor=abc"},
			})
		case 2:
			if r.URL.RequestURI() != "/wiki/api/v2/spaces?cursor=abc" {
				t.Fatalf("second request: %s", r.URL.RequestURI())
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{
					"id":     "1002",
					"key":    "OPS",
					"name":   "Operations",
					"type":   "global",
					"status": "current",
				}},
				"_links": map[string]any{},
			})
		default:
			t.Fatalf("unexpected extra request: %s", r.URL.RequestURI())
		}
	})

	got, err := ListSpaces(context.Background(), c, SpaceFilter{Limit: 2})
	if err != nil {
		t.Fatalf("ListSpaces: %v", err)
	}
	if len(requests) != 2 {
		t.Fatalf("requests = %d, want 2 (%v)", len(requests), requests)
	}
	if len(got) != 2 || got[0].Key != "ENG" || got[1].Key != "OPS" {
		t.Fatalf("unexpected spaces: %+v", got)
	}
}

func TestGetSpace(t *testing.T) {
	_, c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/space/ENG" {
			t.Fatalf("path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(Space{Key: "ENG", Name: "Engineering"})
	})
	s, err := GetSpace(context.Background(), c, "ENG")
	if err != nil || s.Name != "Engineering" {
		t.Fatalf("GetSpace: %v %+v", err, s)
	}
}

func TestGetSpaceCloudUsesV2LookupAndNormalizesFields(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/wiki/api/v2/spaces" {
			t.Fatalf("path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("keys"); got != "ENG" {
			t.Fatalf("keys: %q", got)
		}
		if got := r.URL.Query().Get("limit"); got != "1" {
			t.Fatalf("limit: %q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"results": []map[string]any{{
				"id":     "1001",
				"key":    "ENG",
				"name":   "Engineering",
				"type":   "global",
				"status": "current",
				"description": map[string]any{
					"plain": map[string]any{"value": "Engineering docs"},
				},
				"_links": map[string]any{"webui": "/spaces/ENG"},
			}},
			"_links": map[string]any{},
		})
	})

	s, err := GetSpace(context.Background(), c, "ENG")
	if err != nil {
		t.Fatalf("GetSpace: %v", err)
	}
	if s.Key != "ENG" || s.Name != "Engineering" || s.Status != "current" {
		t.Fatalf("unexpected space: %+v", s)
	}
}

func TestListSpacesCloudFollowsAbsoluteNextLink(t *testing.T) {
	requests := make([]string, 0, 2)
	var base string
	srv, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.URL.RequestURI())
		switch len(requests) {
		case 1:
			base = "http://" + r.Host
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{
					"id":   "1001",
					"key":  "ENG",
					"name": "Engineering",
				}},
				"_links": map[string]any{"next": base + "/wiki/api/v2/spaces?cursor=xyz"},
			})
		case 2:
			if r.URL.RequestURI() != "/wiki/api/v2/spaces?cursor=xyz" {
				t.Fatalf("second request: %s", r.URL.RequestURI())
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{
					"id":   "1002",
					"key":  "OPS",
					"name": "Operations",
				}},
				"_links": map[string]any{},
			})
		default:
			t.Fatalf("unexpected extra request: %s", r.URL.RequestURI())
		}
	})
	defer srv.Close()

	got, err := ListSpaces(context.Background(), c, SpaceFilter{Limit: 2})
	if err != nil {
		t.Fatalf("ListSpaces: %v", err)
	}
	if len(requests) != 2 {
		t.Fatalf("requests = %d, want 2", len(requests))
	}
	if len(got) != 2 {
		t.Fatalf("len(spaces) = %d, want 2", len(got))
	}
}

func TestCloudSpaceNextLinkFixtureParses(t *testing.T) {
	u, err := url.Parse("/wiki/api/v2/spaces?cursor=abc")
	if err != nil {
		t.Fatalf("url.Parse: %v", err)
	}
	if got := u.Query().Get("cursor"); got != "abc" {
		t.Fatalf("cursor: %q", got)
	}
}
