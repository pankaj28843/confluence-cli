package conf

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

func TestGetContentTreeEntityCloudUsesDocumentedEndpoints(t *testing.T) {
	cases := []struct {
		name     string
		kind     string
		wantPath string
	}{
		{name: "database", kind: "database", wantPath: "/wiki/api/v2/databases/12345"},
		{name: "folder", kind: "folder", wantPath: "/wiki/api/v2/folders/12345"},
		{name: "whiteboard", kind: "whiteboard", wantPath: "/wiki/api/v2/whiteboards/12345"},
		{name: "smart link", kind: "smart-link", wantPath: "/wiki/api/v2/embeds/12345"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet || r.URL.Path != tc.wantPath {
					t.Fatalf("request = %s %s", r.Method, r.URL.Path)
				}
				for _, key := range []string{"include-collaborators", "include-direct-children", "include-operations", "include-properties"} {
					if got := r.URL.Query().Get(key); got != "true" {
						t.Fatalf("%s: %q", key, got)
					}
				}
				_ = json.NewEncoder(w).Encode(ContentTreeEntity{
					ID:       "12345",
					Type:     tc.kind,
					Status:   "current",
					Title:    "Example",
					SpaceID:  "42",
					EmbedURL: "https://example.com/card",
				})
			})

			got, err := GetContentTreeEntity(context.Background(), c, tc.kind, "12345", ContentTreeEntityOptions{
				IncludeCollaborators:  true,
				IncludeDirectChildren: true,
				IncludeOperations:     true,
				IncludeProperties:     true,
			})
			if err != nil {
				t.Fatalf("GetContentTreeEntity: %v", err)
			}
			if got.ID != "12345" || got.Title != "Example" {
				t.Fatalf("unexpected entity: %+v", got)
			}
		})
	}
}

func TestListContentTreeDirectChildrenCloudUsesDocumentedEndpoints(t *testing.T) {
	cases := []struct {
		name     string
		kind     string
		wantPath string
	}{
		{name: "database", kind: "database", wantPath: "/wiki/api/v2/databases/12345/direct-children"},
		{name: "folder", kind: "folder", wantPath: "/wiki/api/v2/folders/12345/direct-children"},
		{name: "whiteboard", kind: "whiteboard", wantPath: "/wiki/api/v2/whiteboards/12345/direct-children"},
		{name: "smart link", kind: "smart-link", wantPath: "/wiki/api/v2/embeds/12345/direct-children"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet || r.URL.Path != tc.wantPath {
					t.Fatalf("request = %s %s", r.Method, r.URL.Path)
				}
				if got := r.URL.Query().Get("limit"); got != "2" {
					t.Fatalf("limit: %q", got)
				}
				if got := r.URL.Query().Get("sort"); got != "position" {
					t.Fatalf("sort: %q", got)
				}
				_ = json.NewEncoder(w).Encode(map[string]any{
					"results": []Content{{ID: "child1", Type: "database", Title: "Database child", ChildPosition: 1}},
					"_links":  map[string]any{},
				})
			})

			got, err := ListContentTreeDirectChildren(context.Background(), c, tc.kind, "12345", DirectChildrenOptions{Limit: 2, Sort: "position"})
			if err != nil {
				t.Fatalf("ListContentTreeDirectChildren: %v", err)
			}
			if len(got) != 1 || got[0].ID != "child1" {
				t.Fatalf("unexpected children: %+v", got)
			}
		})
	}
}

func TestListContentTreeDirectChildrenCloudPaginatesAndFiltersTypes(t *testing.T) {
	requests := make([]string, 0, 2)
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.URL.RequestURI())
		switch len(requests) {
		case 1:
			if r.URL.RequestURI() != "/wiki/api/v2/folders/12345/direct-children?limit=2" {
				t.Fatalf("first request: %s", r.URL.RequestURI())
			}
			w.Header().Set("Link", "<http://"+r.Host+"/wiki/api/v2/folders/12345/direct-children?cursor=abc&limit=2>; rel=\"next\"")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []Content{{ID: "page1", Type: "page", Title: "Page child"}},
				"_links":  map[string]any{},
			})
		case 2:
			if r.URL.RequestURI() != "/wiki/api/v2/folders/12345/direct-children?cursor=abc&limit=2" {
				t.Fatalf("second request: %s", r.URL.RequestURI())
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []Content{
					{ID: "folder1", Type: "folder", Title: "Folder child"},
					{ID: "db1", Type: "database", Title: "Database child"},
				},
				"_links": map[string]any{},
			})
		default:
			t.Fatalf("unexpected extra request: %s", r.URL.RequestURI())
		}
	})

	got, err := ListContentTreeDirectChildren(context.Background(), c, "folder", "12345", DirectChildrenOptions{
		Limit: 2,
		Types: []string{"database"},
	})
	if err != nil {
		t.Fatalf("ListContentTreeDirectChildren: %v", err)
	}
	if len(requests) != 2 {
		t.Fatalf("requests = %d, want 2 (%v)", len(requests), requests)
	}
	if len(got) != 1 || got[0].ID != "db1" {
		t.Fatalf("unexpected filtered children: %+v", got)
	}
}

func TestContentTreeEntityHelpersRejectUnsupportedInputs(t *testing.T) {
	_, server := testClientWithFlavor(t, client.FlavorServer, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s", r.URL.RequestURI())
	})
	_, cloud := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s", r.URL.RequestURI())
	})

	if _, err := GetContentTreeEntity(context.Background(), server, "database", "12345", ContentTreeEntityOptions{}); err == nil || !strings.Contains(err.Error(), "Cloud") {
		t.Fatalf("server get error = %v", err)
	}
	if _, err := ListContentTreeDirectChildren(context.Background(), server, "folder", "12345", DirectChildrenOptions{}); err == nil || !strings.Contains(err.Error(), "Cloud") {
		t.Fatalf("server children error = %v", err)
	}
	if _, err := GetContentTreeEntity(context.Background(), cloud, "database", "", ContentTreeEntityOptions{}); err == nil || !strings.Contains(err.Error(), "id") {
		t.Fatalf("missing id error = %v", err)
	}
	if _, err := GetContentTreeEntity(context.Background(), cloud, "page", "12345", ContentTreeEntityOptions{}); err == nil || !strings.Contains(err.Error(), "unsupported") {
		t.Fatalf("unsupported kind error = %v", err)
	}
}
