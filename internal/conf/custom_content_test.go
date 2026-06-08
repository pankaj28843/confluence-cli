package conf

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

func TestListCustomContentCloudUsesDocumentedEndpointAndPagination(t *testing.T) {
	requests := make([]string, 0, 2)
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.URL.RequestURI())
		switch len(requests) {
		case 1:
			if r.Method != http.MethodGet || r.URL.Path != "/wiki/api/v2/custom-content" {
				t.Fatalf("request = %s %s", r.Method, r.URL.Path)
			}
			q := r.URL.Query()
			if got := q.Get("type"); got != "ac:example" {
				t.Fatalf("type = %q", got)
			}
			if got := q.Get("limit"); got != "2" {
				t.Fatalf("limit = %q", got)
			}
			if got := q.Get("sort"); got != "-created-date" {
				t.Fatalf("sort = %q", got)
			}
			if got := q.Get("body-format"); got != "storage" {
				t.Fatalf("body-format = %q", got)
			}
			if got := q["id"]; !reflect.DeepEqual(got, []string{"777"}) {
				t.Fatalf("id = %v", got)
			}
			if got := q["space-id"]; !reflect.DeepEqual(got, []string{"100"}) {
				t.Fatalf("space-id = %v", got)
			}
			w.Header().Set("Link", "<http://"+r.Host+"/wiki/api/v2/custom-content?cursor=abc&limit=2&type=ac%3Aexample>; rel=\"next\"")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []Content{{ID: "cc1", Type: "ac:example", Title: "Custom one", SpaceID: "100"}},
				"_links":  map[string]any{},
			})
		case 2:
			if r.URL.RequestURI() != "/wiki/api/v2/custom-content?cursor=abc&limit=2&type=ac%3Aexample" {
				t.Fatalf("second request: %s", r.URL.RequestURI())
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []Content{{ID: "cc2", Type: "ac:example", Title: "Custom two"}},
				"_links":  map[string]any{},
			})
		default:
			t.Fatalf("unexpected extra request: %s", r.URL.RequestURI())
		}
	})

	got, err := ListCustomContent(context.Background(), c, CustomContentListOptions{
		Type:       "ac:example",
		IDs:        []string{"777"},
		SpaceIDs:   []string{"100"},
		Sort:       "-created-date",
		BodyFormat: "storage",
		Limit:      2,
	})
	if err != nil {
		t.Fatalf("ListCustomContent: %v", err)
	}
	if len(requests) != 2 {
		t.Fatalf("requests = %d, want 2", len(requests))
	}
	if len(got) != 2 || got[0].ID != "cc1" || got[1].ID != "cc2" {
		t.Fatalf("unexpected custom content: %+v", got)
	}
}

func TestListCustomContentCloudUsesContainerEndpoints(t *testing.T) {
	cases := []struct {
		name          string
		containerType string
		wantPath      string
		sort          string
	}{
		{name: "page", containerType: "page", wantPath: "/wiki/api/v2/pages/12345/custom-content", sort: "-created-date"},
		{name: "blogpost", containerType: "blogpost", wantPath: "/wiki/api/v2/blogposts/12345/custom-content", sort: "created-date"},
		{name: "space", containerType: "space", wantPath: "/wiki/api/v2/spaces/12345/custom-content"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet || r.URL.Path != tc.wantPath {
					t.Fatalf("request = %s %s", r.Method, r.URL.Path)
				}
				if got := r.URL.Query().Get("type"); got != "ac:example" {
					t.Fatalf("type = %q", got)
				}
				if got := r.URL.Query().Get("limit"); got != "1" {
					t.Fatalf("limit = %q", got)
				}
				if got := r.URL.Query().Get("body-format"); got != "storage" {
					t.Fatalf("body-format = %q", got)
				}
				if got := r.URL.Query().Get("sort"); got != tc.sort {
					t.Fatalf("sort = %q", got)
				}
				_ = json.NewEncoder(w).Encode(map[string]any{
					"results": []Content{{ID: "cc1", Type: "ac:example", Title: "Custom one"}},
					"_links":  map[string]any{},
				})
			})

			got, err := ListCustomContent(context.Background(), c, CustomContentListOptions{
				Type:          "ac:example",
				ContainerType: tc.containerType,
				ContainerID:   "12345",
				Sort:          tc.sort,
				BodyFormat:    "storage",
				Limit:         1,
			})
			if err != nil {
				t.Fatalf("ListCustomContent: %v", err)
			}
			if len(got) != 1 || got[0].ID != "cc1" {
				t.Fatalf("unexpected custom content: %+v", got)
			}
		})
	}
}

func TestGetCustomContentCloudUsesDocumentedEndpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/wiki/api/v2/custom-content/12345" {
			t.Fatalf("request = %s %s", r.Method, r.URL.Path)
		}
		for _, key := range []string{"include-labels", "include-properties", "include-operations", "include-versions", "include-version", "include-collaborators"} {
			if got := r.URL.Query().Get(key); got != "true" {
				t.Fatalf("%s = %q", key, got)
			}
		}
		if got := r.URL.Query().Get("body-format"); got != "storage" {
			t.Fatalf("body-format = %q", got)
		}
		if got := r.URL.Query().Get("version"); got != "7" {
			t.Fatalf("version = %q", got)
		}
		_ = json.NewEncoder(w).Encode(Content{
			ID:      "12345",
			Type:    "ac:example",
			Status:  "current",
			Title:   "Custom one",
			SpaceID: "100",
		})
	})

	got, err := GetCustomContent(context.Background(), c, "12345", CustomContentGetOptions{
		BodyFormat:           "storage",
		Version:              7,
		IncludeLabels:        true,
		IncludeProperties:    true,
		IncludeOperations:    true,
		IncludeVersions:      true,
		IncludeVersion:       true,
		IncludeCollaborators: true,
	})
	if err != nil {
		t.Fatalf("GetCustomContent: %v", err)
	}
	if got.ID != "12345" || got.Title != "Custom one" {
		t.Fatalf("unexpected custom content: %+v", got)
	}
}

func TestListCustomContentChildrenCloudUsesChildrenEndpoint(t *testing.T) {
	requests := make([]string, 0, 2)
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.URL.RequestURI())
		switch len(requests) {
		case 1:
			if r.URL.RequestURI() != "/wiki/api/v2/custom-content/12345/children?limit=2&sort=title" {
				t.Fatalf("first request: %s", r.URL.RequestURI())
			}
			w.Header().Set("Link", "<http://"+r.Host+"/wiki/api/v2/custom-content/12345/children?cursor=abc&limit=2>; rel=\"next\"")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []Content{{ID: "child1", Type: "ac:keep", Title: "Keep"}},
				"_links":  map[string]any{},
			})
		case 2:
			if r.URL.RequestURI() != "/wiki/api/v2/custom-content/12345/children?cursor=abc&limit=2" {
				t.Fatalf("second request: %s", r.URL.RequestURI())
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []Content{{ID: "child2", Type: "ac:skip", Title: "Skip"}},
				"_links":  map[string]any{},
			})
		default:
			t.Fatalf("unexpected request: %s", r.URL.RequestURI())
		}
	})

	got, err := ListCustomContentChildren(context.Background(), c, "12345", DirectChildrenOptions{
		Limit: 2,
		Sort:  "title",
		Types: []string{"ac:keep"},
	})
	if err != nil {
		t.Fatalf("ListCustomContentChildren: %v", err)
	}
	if len(requests) != 2 {
		t.Fatalf("requests = %d, want 2", len(requests))
	}
	if len(got) != 1 || got[0].ID != "child1" {
		t.Fatalf("unexpected children: %+v", got)
	}
}

func TestCustomContentHelpersRejectUnsupportedInputs(t *testing.T) {
	_, server := testClientWithFlavor(t, client.FlavorServer, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s", r.URL.RequestURI())
	})
	_, cloud := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s", r.URL.RequestURI())
	})

	if _, err := ListCustomContent(context.Background(), server, CustomContentListOptions{Type: "ac:example"}); err == nil || !strings.Contains(err.Error(), "Cloud") {
		t.Fatalf("server list error = %v", err)
	}
	if _, err := GetCustomContent(context.Background(), server, "12345", CustomContentGetOptions{}); err == nil || !strings.Contains(err.Error(), "Cloud") {
		t.Fatalf("server get error = %v", err)
	}
	if _, err := ListCustomContentChildren(context.Background(), server, "12345", DirectChildrenOptions{}); err == nil || !strings.Contains(err.Error(), "Cloud") {
		t.Fatalf("server children error = %v", err)
	}
	if _, err := ListCustomContent(context.Background(), cloud, CustomContentListOptions{}); err == nil || !strings.Contains(err.Error(), "type") {
		t.Fatalf("missing type error = %v", err)
	}
	if _, err := ListCustomContent(context.Background(), cloud, CustomContentListOptions{Type: "ac:example", ContainerType: "page"}); err == nil || !strings.Contains(err.Error(), "container id") {
		t.Fatalf("missing container id error = %v", err)
	}
	if _, err := ListCustomContent(context.Background(), cloud, CustomContentListOptions{Type: "ac:example", ContainerType: "database", ContainerID: "12345"}); err == nil || !strings.Contains(err.Error(), "unsupported") {
		t.Fatalf("unsupported container error = %v", err)
	}
	if _, err := GetCustomContent(context.Background(), cloud, "", CustomContentGetOptions{}); err == nil || !strings.Contains(err.Error(), "id") {
		t.Fatalf("missing get id error = %v", err)
	}
	if _, err := ListCustomContentChildren(context.Background(), cloud, "", DirectChildrenOptions{}); err == nil || !strings.Contains(err.Error(), "id") {
		t.Fatalf("missing children id error = %v", err)
	}
}
