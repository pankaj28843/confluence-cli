package conf

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

func TestGetContent(t *testing.T) {
	_, c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/rest/api/content/12345") {
			t.Fatalf("path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("expand"); got != DefaultExpand {
			t.Fatalf("expand: %s", got)
		}
		cnt := Content{ID: "12345", Title: "Demo", Type: "page"}
		cnt.Body.Storage.Value = "<h1>Hello</h1><p>world</p>"
		_ = json.NewEncoder(w).Encode(cnt)
	})
	got, err := GetContent(context.Background(), c, "12345", "")
	if err != nil || got.ID != "12345" {
		t.Fatalf("GetContent: %v %+v", err, got)
	}
	if md := got.RenderMarkdown(false); !strings.Contains(md, "# Hello") {
		t.Fatalf("markdown conversion: %q", md)
	}
}

func TestSearchCQL(t *testing.T) {
	_, c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/content/search" {
			t.Fatalf("path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("cql"); got != "type=page" {
			t.Fatalf("cql: %s", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"results": []Content{{ID: "1", Title: "Hit"}}})
	})
	got, err := SearchCQL(context.Background(), c, "type=page", 3, "")
	if err != nil || len(got) != 1 {
		t.Fatalf("SearchCQL: %v %+v", err, got)
	}
}

func TestHTMLToMarkdown(t *testing.T) {
	// Block-level elements are preserved; inline emphasis/links inside a <p>
	// are flattened to plain text — matches the ported converter semantics.
	in := `<h1>Title</h1><p>Body text.</p><ul><li>one</li><li>two</li></ul><pre><code>fn()</code></pre>`
	md := HTMLToMarkdown(in)
	for _, want := range []string{"# Title", "Body text.", "- one", "- two", "```", "fn()"} {
		if !strings.Contains(md, want) {
			t.Errorf("missing %q in %q", want, md)
		}
	}
}

func TestGetChildren(t *testing.T) {
	_, c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/rest/api/content/12345/child/page") {
			t.Fatalf("path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"results": []Content{{ID: "c1", Title: "child"}}})
	})
	got, err := GetChildren(context.Background(), c, "12345", "", 10)
	if err != nil || len(got) != 1 {
		t.Fatalf("GetChildren: %v %+v", err, got)
	}
}

func TestGetChildrenCloudPagesUsesV2DirectChildrenAndPagination(t *testing.T) {
	requests := make([]string, 0, 2)
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.Method+" "+r.URL.RequestURI())
		w.Header().Set("Content-Type", "application/json")

		switch len(requests) {
		case 1:
			if r.Method != http.MethodGet || r.URL.Path != "/wiki/api/v2/pages/12345/direct-children" {
				t.Fatalf("first request: %s %s", r.Method, r.URL.RequestURI())
			}
			if got := r.URL.Query().Get("limit"); got != "2" {
				t.Fatalf("limit = %q, want 2", got)
			}
			w.Header().Set("Link", "<http://"+r.Host+"/wiki/api/v2/pages/12345/direct-children?cursor=abc&limit=2>; rel=\"next\"")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{
					{
						"id":            "db1",
						"type":          "database",
						"status":        "current",
						"title":         "Database",
						"spaceId":       "1001",
						"childPosition": 1,
					},
					{
						"id":            "p1",
						"type":          "page",
						"status":        "current",
						"title":         "Page One",
						"spaceId":       "1001",
						"childPosition": 2,
					},
				},
				"_links": map[string]any{},
			})
		case 2:
			if r.Method != http.MethodGet || r.URL.RequestURI() != "/wiki/api/v2/pages/12345/direct-children?cursor=abc&limit=2" {
				t.Fatalf("second request: %s %s", r.Method, r.URL.RequestURI())
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{
					"id":            "p2",
					"type":          "page",
					"status":        "current",
					"title":         "Page Two",
					"spaceId":       "1001",
					"childPosition": 3,
				}},
				"_links": map[string]any{},
			})
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
	})

	got, err := GetChildren(context.Background(), c, "12345", "page", 2)
	if err != nil {
		t.Fatalf("GetChildren: %v", err)
	}
	if len(requests) != 2 {
		t.Fatalf("requests = %v, want 2 paginated requests", requests)
	}
	if len(got) != 2 {
		t.Fatalf("len(children) = %d, want 2 (%+v)", len(got), got)
	}
	if got[0].ID != "p1" || got[0].Type != "page" || got[0].SpaceID != "1001" || got[0].ChildPosition != 2 {
		t.Fatalf("unexpected first page child: %+v", got[0])
	}
	if got[1].ID != "p2" || got[1].Type != "page" || got[1].ChildPosition != 3 {
		t.Fatalf("unexpected second page child: %+v", got[1])
	}
}

func TestListDirectChildrenCloudUsesV2EndpointAndPagination(t *testing.T) {
	requests := make([]string, 0, 2)
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.Method+" "+r.URL.RequestURI())
		w.Header().Set("Content-Type", "application/json")

		switch len(requests) {
		case 1:
			if r.Method != http.MethodGet || r.URL.Path != "/wiki/api/v2/pages/12345/direct-children" {
				t.Fatalf("first request: %s %s", r.Method, r.URL.RequestURI())
			}
			if got := r.URL.Query().Get("limit"); got != "3" {
				t.Fatalf("limit = %q, want 3", got)
			}
			w.Header().Set("Link", "<http://"+r.Host+"/wiki/api/v2/pages/12345/direct-children?cursor=abc&limit=3>; rel=\"next\"")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{
					{
						"id":            "db1",
						"type":          "database",
						"status":        "current",
						"title":         "Database",
						"spaceId":       "1001",
						"childPosition": 1,
					},
					{
						"id":            "p1",
						"type":          "page",
						"status":        "current",
						"title":         "Page One",
						"spaceId":       "1001",
						"childPosition": 2,
					},
				},
				"_links": map[string]any{},
			})
		case 2:
			if r.Method != http.MethodGet || r.URL.RequestURI() != "/wiki/api/v2/pages/12345/direct-children?cursor=abc&limit=3" {
				t.Fatalf("second request: %s %s", r.Method, r.URL.RequestURI())
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{
					"id":            "wb1",
					"type":          "whiteboard",
					"status":        "current",
					"title":         "Whiteboard",
					"spaceId":       "1001",
					"childPosition": 3,
				}},
				"_links": map[string]any{},
			})
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
	})

	got, err := ListDirectChildren(context.Background(), c, "12345", DirectChildrenOptions{Limit: 3})
	if err != nil {
		t.Fatalf("ListDirectChildren: %v", err)
	}
	if len(requests) != 2 {
		t.Fatalf("requests = %v, want 2 paginated requests", requests)
	}
	if len(got) != 3 {
		t.Fatalf("len(children) = %d, want 3 (%+v)", len(got), got)
	}
	if got[0].ID != "db1" || got[0].Type != "database" || got[0].ChildPosition != 1 {
		t.Fatalf("unexpected first child: %+v", got[0])
	}
	if got[2].ID != "wb1" || got[2].Type != "whiteboard" || got[2].ChildPosition != 3 {
		t.Fatalf("unexpected third child: %+v", got[2])
	}
}

func TestListDirectChildrenCloudFiltersTypesAcrossPages(t *testing.T) {
	requests := make([]string, 0, 2)
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.Method+" "+r.URL.RequestURI())
		w.Header().Set("Content-Type", "application/json")

		switch len(requests) {
		case 1:
			w.Header().Set("Link", "<http://"+r.Host+"/wiki/api/v2/pages/12345/direct-children?cursor=abc&limit=1>; rel=\"next\"")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{
					"id":            "db1",
					"type":          "database",
					"title":         "Database",
					"childPosition": 1,
				}},
				"_links": map[string]any{},
			})
		case 2:
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{
					"id":            "p1",
					"type":          "page",
					"title":         "Page One",
					"childPosition": 2,
				}},
				"_links": map[string]any{},
			})
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
	})

	got, err := ListDirectChildren(context.Background(), c, "12345", DirectChildrenOptions{
		Limit: 1,
		Types: []string{"page"},
	})
	if err != nil {
		t.Fatalf("ListDirectChildren: %v", err)
	}
	if len(requests) != 2 {
		t.Fatalf("requests = %v, want pagination until a matching page child is found", requests)
	}
	if len(got) != 1 || got[0].ID != "p1" || got[0].Type != "page" {
		t.Fatalf("children = %+v, want page p1 only", got)
	}
}

func TestListDirectChildrenServerUsesChildEndpointWithExpandTypes(t *testing.T) {
	_, c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/rest/api/content/12345/child" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		q := r.URL.Query()
		if got := q.Get("limit"); got != "2" {
			t.Fatalf("limit = %q, want 2", got)
		}
		if got := q.Get("expand"); got != "page,comment" {
			t.Fatalf("expand = %q, want page,comment", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"results": []map[string]any{
				{
					"id":     "p1",
					"type":   "page",
					"status": "current",
					"title":  "Page One",
				},
				{
					"id":     "c1",
					"type":   "comment",
					"status": "current",
					"title":  "Comment One",
				},
			},
			"_links": map[string]any{},
		})
	})

	got, err := ListDirectChildren(context.Background(), c, "12345", DirectChildrenOptions{
		Limit: 2,
		Types: []string{"page", "comment"},
	})
	if err != nil {
		t.Fatalf("ListDirectChildren: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len(children) = %d, want 2 (%+v)", len(got), got)
	}
	if got[0].ID != "p1" || got[0].Type != "page" || got[1].ID != "c1" || got[1].Type != "comment" {
		t.Fatalf("unexpected children: %+v", got)
	}
}

func TestListDirectChildrenRequiresID(t *testing.T) {
	_, c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
	})

	_, err := ListDirectChildren(context.Background(), c, "", DirectChildrenOptions{Limit: 1})
	if err == nil || !strings.Contains(err.Error(), "ID is required") {
		t.Fatalf("ListDirectChildren error = %v, want ID required", err)
	}
}

func TestListPageDescendantsCloudUsesV2EndpointAndPagination(t *testing.T) {
	requests := make([]string, 0, 2)
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.Method+" "+r.URL.RequestURI())
		w.Header().Set("Content-Type", "application/json")

		switch len(requests) {
		case 1:
			if r.Method != http.MethodGet || r.URL.Path != "/wiki/api/v2/pages/12345/descendants" {
				t.Fatalf("first request: %s %s", r.Method, r.URL.RequestURI())
			}
			q := r.URL.Query()
			if got := q.Get("limit"); got != "2" {
				t.Fatalf("limit = %q, want 2", got)
			}
			if got := q.Get("depth"); got != "3" {
				t.Fatalf("depth = %q, want 3", got)
			}
			w.Header().Set("Link", "<http://"+r.Host+"/wiki/api/v2/pages/12345/descendants?cursor=abc>; rel=\"next\"")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{
					"id":            "c1",
					"type":          "page",
					"status":        "current",
					"title":         "Child",
					"parentId":      "12345",
					"depth":         1,
					"childPosition": 10,
				}},
				"_links": map[string]any{},
			})
		case 2:
			if r.Method != http.MethodGet || r.URL.RequestURI() != "/wiki/api/v2/pages/12345/descendants?cursor=abc" {
				t.Fatalf("second request: %s %s", r.Method, r.URL.RequestURI())
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{
					"id":            "g1",
					"type":          "page",
					"status":        "current",
					"title":         "Grandchild",
					"parentId":      "c1",
					"depth":         2,
					"childPosition": 20,
				}},
				"_links": map[string]any{},
			})
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
	})

	got, err := ListPageDescendants(context.Background(), c, "12345", PageDescendantsOptions{Limit: 2, Depth: 3})
	if err != nil {
		t.Fatalf("ListPageDescendants: %v", err)
	}
	if len(requests) != 2 {
		t.Fatalf("requests = %v, want 2 paginated requests", requests)
	}
	if len(got) != 2 {
		t.Fatalf("len(descendants) = %d, want 2 (%+v)", len(got), got)
	}
	if got[0].ID != "c1" || got[0].ParentID != "12345" || got[0].Depth != 1 || got[0].ChildPosition != 10 {
		t.Fatalf("unexpected first descendant: %+v", got[0])
	}
	if got[1].ID != "g1" || got[1].ParentID != "c1" || got[1].Depth != 2 {
		t.Fatalf("unexpected second descendant: %+v", got[1])
	}
}

func TestListPageDescendantsServerWalksChildPagesRecursively(t *testing.T) {
	requests := make([]string, 0, 2)
	_, c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.Method+" "+r.URL.RequestURI())
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want GET", r.Method)
		}
		q := r.URL.Query()
		if got := q.Get("limit"); got != "5" {
			t.Fatalf("limit = %q, want 5", got)
		}
		if got := q.Get("expand"); got != "version,space" {
			t.Fatalf("expand = %q, want version,space", got)
		}

		switch r.URL.Path {
		case "/rest/api/content/12345/child/page":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{
					"id":     "c1",
					"type":   "page",
					"status": "current",
					"title":  "Child",
					"space":  map[string]any{"key": "ENG"},
				}},
			})
		case "/rest/api/content/c1/child/page":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{
					"id":     "g1",
					"type":   "page",
					"status": "current",
					"title":  "Grandchild",
					"space":  map[string]any{"key": "ENG"},
				}},
			})
		default:
			t.Fatalf("path: %s", r.URL.Path)
		}
	})

	got, err := ListPageDescendants(context.Background(), c, "12345", PageDescendantsOptions{Limit: 5, Depth: 2})
	if err != nil {
		t.Fatalf("ListPageDescendants: %v", err)
	}
	if len(requests) != 2 {
		t.Fatalf("requests = %v, want root and child traversal", requests)
	}
	if len(got) != 2 {
		t.Fatalf("len(descendants) = %d, want 2 (%+v)", len(got), got)
	}
	if got[0].ID != "c1" || got[0].ParentID != "12345" || got[0].Depth != 1 {
		t.Fatalf("unexpected child: %+v", got[0])
	}
	if got[1].ID != "g1" || got[1].ParentID != "c1" || got[1].Depth != 2 {
		t.Fatalf("unexpected grandchild: %+v", got[1])
	}
}
