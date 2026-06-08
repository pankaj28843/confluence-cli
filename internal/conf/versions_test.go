package conf

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

func TestListVersionsCloudUsesV2DocumentedEndpoints(t *testing.T) {
	cases := []struct {
		name     string
		target   string
		wantPath string
		wantBody bool
	}{
		{name: "page", target: "page", wantPath: "/wiki/api/v2/pages/12345/versions", wantBody: true},
		{name: "blogpost", target: "blogpost", wantPath: "/wiki/api/v2/blogposts/12345/versions", wantBody: true},
		{name: "attachment", target: "attachment", wantPath: "/wiki/api/v2/attachments/12345/versions"},
		{name: "footer comment", target: "footer-comment", wantPath: "/wiki/api/v2/footer-comments/12345/versions", wantBody: true},
		{name: "inline comment", target: "inline-comment", wantPath: "/wiki/api/v2/inline-comments/12345/versions", wantBody: true},
		{name: "custom content", target: "custom-content", wantPath: "/wiki/api/v2/custom-content/12345/versions", wantBody: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet || r.URL.Path != tc.wantPath {
					t.Fatalf("request = %s %s", r.Method, r.URL.Path)
				}
				if got := r.URL.Query().Get("limit"); got != "2" {
					t.Fatalf("limit = %q", got)
				}
				if got := r.URL.Query().Get("sort"); got != "-modified-date" {
					t.Fatalf("sort = %q", got)
				}
				wantBodyFormat := ""
				if tc.wantBody {
					wantBodyFormat = "storage"
				}
				if got := r.URL.Query().Get("body-format"); got != wantBodyFormat {
					t.Fatalf("body-format = %q", got)
				}
				_ = json.NewEncoder(w).Encode(map[string]any{
					"results": []Version{{
						Number:    7,
						AuthorID:  "acct-1",
						Message:   "edited",
						CreatedAt: "2026-06-08T10:00:00Z",
					}},
					"_links": map[string]any{},
				})
			})

			got, err := ListVersions(context.Background(), c, tc.target, "12345", VersionListOptions{
				Limit:      2,
				Sort:       "-modified-date",
				BodyFormat: "storage",
			})
			if err != nil {
				t.Fatalf("ListVersions: %v", err)
			}
			if len(got) != 1 || got[0].Number != 7 || got[0].AuthorID != "acct-1" {
				t.Fatalf("unexpected versions: %+v", got)
			}
		})
	}
}

func TestListVersionsCloudPaginatesWithLinkHeader(t *testing.T) {
	requests := make([]string, 0, 2)
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.URL.RequestURI())
		switch len(requests) {
		case 1:
			if r.URL.RequestURI() != "/wiki/api/v2/pages/12345/versions?limit=2" {
				t.Fatalf("first request: %s", r.URL.RequestURI())
			}
			w.Header().Set("Link", "<http://"+r.Host+"/wiki/api/v2/pages/12345/versions?cursor=abc&limit=2>; rel=\"next\"")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []Version{{Number: 1}},
				"_links":  map[string]any{},
			})
		case 2:
			if r.URL.RequestURI() != "/wiki/api/v2/pages/12345/versions?cursor=abc&limit=2" {
				t.Fatalf("second request: %s", r.URL.RequestURI())
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []Version{{Number: 2}},
				"_links":  map[string]any{},
			})
		default:
			t.Fatalf("unexpected request: %s", r.URL.RequestURI())
		}
	})

	got, err := ListVersions(context.Background(), c, "page", "12345", VersionListOptions{Limit: 2})
	if err != nil {
		t.Fatalf("ListVersions: %v", err)
	}
	if len(requests) != 2 {
		t.Fatalf("requests = %d, want 2", len(requests))
	}
	if len(got) != 2 || got[0].Number != 1 || got[1].Number != 2 {
		t.Fatalf("unexpected versions: %+v", got)
	}
}

func TestGetVersionCloudUsesV2DocumentedEndpoints(t *testing.T) {
	cases := []struct {
		name     string
		target   string
		wantPath string
	}{
		{name: "page", target: "page", wantPath: "/wiki/api/v2/pages/12345/versions/7"},
		{name: "blogpost", target: "blogpost", wantPath: "/wiki/api/v2/blogposts/12345/versions/7"},
		{name: "attachment", target: "attachment", wantPath: "/wiki/api/v2/attachments/12345/versions/7"},
		{name: "footer comment", target: "footer-comment", wantPath: "/wiki/api/v2/footer-comments/12345/versions/7"},
		{name: "inline comment", target: "inline-comment", wantPath: "/wiki/api/v2/inline-comments/12345/versions/7"},
		{name: "custom content", target: "custom-content", wantPath: "/wiki/api/v2/custom-content/12345/versions/7"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet || r.URL.Path != tc.wantPath {
					t.Fatalf("request = %s %s", r.Method, r.URL.Path)
				}
				_ = json.NewEncoder(w).Encode(Version{
					Number:              7,
					AuthorID:            "acct-1",
					Message:             "edited",
					ContentTypeModified: true,
					PrevVersion:         6,
					NextVersion:         8,
				})
			})

			got, err := GetVersion(context.Background(), c, tc.target, "12345", 7)
			if err != nil {
				t.Fatalf("GetVersion: %v", err)
			}
			if got.Number != 7 || !got.ContentTypeModified || got.PrevVersion != 6 || got.NextVersion != 8 {
				t.Fatalf("unexpected version: %+v", got)
			}
		})
	}
}

func TestVersionsServerUsesContentVersionRoutes(t *testing.T) {
	requests := make([]string, 0, 2)
	_, c := testClientWithFlavor(t, client.FlavorServer, func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.URL.RequestURI())
		switch len(requests) {
		case 1:
			if r.Method != http.MethodGet || r.URL.Path != "/rest/api/content/12345/version" {
				t.Fatalf("list request = %s %s", r.Method, r.URL.Path)
			}
			if got := r.URL.Query().Get("limit"); got != "2" {
				t.Fatalf("limit = %q", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"results": []Version{{Number: 1}}})
		case 2:
			if r.Method != http.MethodGet || r.URL.Path != "/rest/api/content/12345/version/1" {
				t.Fatalf("detail request = %s %s", r.Method, r.URL.Path)
			}
			_ = json.NewEncoder(w).Encode(Version{Number: 1, Message: "server"})
		default:
			t.Fatalf("unexpected request: %s", r.URL.RequestURI())
		}
	})

	list, err := ListVersions(context.Background(), c, "page", "12345", VersionListOptions{Limit: 2})
	if err != nil {
		t.Fatalf("ListVersions server: %v", err)
	}
	if len(list) != 1 || list[0].Number != 1 {
		t.Fatalf("unexpected server versions: %+v", list)
	}
	detail, err := GetVersion(context.Background(), c, "page", "12345", 1)
	if err != nil {
		t.Fatalf("GetVersion server: %v", err)
	}
	if detail.Message != "server" {
		t.Fatalf("unexpected server detail: %+v", detail)
	}
}

func TestVersionHelpersRejectUnsupportedInputs(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s", r.URL.RequestURI())
	})

	if _, err := ListVersions(context.Background(), c, "page", "", VersionListOptions{}); err == nil || !strings.Contains(err.Error(), "id") {
		t.Fatalf("missing id error = %v", err)
	}
	if _, err := GetVersion(context.Background(), c, "page", "12345", 0); err == nil || !strings.Contains(err.Error(), "version") {
		t.Fatalf("missing version error = %v", err)
	}
	if _, err := ListVersions(context.Background(), c, "database", "12345", VersionListOptions{}); err == nil || !strings.Contains(err.Error(), "unsupported") {
		t.Fatalf("unsupported target error = %v", err)
	}
}
