package conf

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

func TestListLabelsServerUsesV1ContentLabelEndpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorServer, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/content/12345/label" {
			t.Fatalf("path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"results": []Label{{Name: "needs-review", Prefix: "global"}},
		})
	})

	got, err := ListLabels(context.Background(), c, "12345", 2)
	if err != nil {
		t.Fatalf("ListLabels: %v", err)
	}
	if len(got) != 1 || got[0].Name != "needs-review" {
		t.Fatalf("unexpected labels: %+v", got)
	}
}

func TestListLabelsCloudUsesV2PageLabelEndpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/wiki/api/v2/pages/12345/labels" {
			t.Fatalf("path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("limit"); got != "2" {
			t.Fatalf("limit: %q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"results": []Label{{ID: "label1", Name: "cloud", Prefix: "global"}},
			"_links":  map[string]any{},
		})
	})

	got, err := ListLabels(context.Background(), c, "12345", 2)
	if err != nil {
		t.Fatalf("ListLabels: %v", err)
	}
	if len(got) != 1 || got[0].ID != "label1" || got[0].Name != "cloud" {
		t.Fatalf("unexpected labels: %+v", got)
	}
}

func TestListLabelsCloudFollowsV2Pagination(t *testing.T) {
	requests := make([]string, 0, 2)
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.URL.RequestURI())
		switch len(requests) {
		case 1:
			if r.URL.RequestURI() != "/wiki/api/v2/pages/12345/labels?limit=2" {
				t.Fatalf("first request: %s", r.URL.RequestURI())
			}
			w.Header().Set("Link", "<http://"+r.Host+"/wiki/api/v2/pages/12345/labels?cursor=abc>; rel=\"next\"")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []Label{{ID: "label1", Name: "first", Prefix: "global"}},
				"_links":  map[string]any{},
			})
		case 2:
			if r.URL.RequestURI() != "/wiki/api/v2/pages/12345/labels?cursor=abc" {
				t.Fatalf("second request: %s", r.URL.RequestURI())
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []Label{{ID: "label2", Name: "second", Prefix: "global"}},
				"_links":  map[string]any{},
			})
		default:
			t.Fatalf("unexpected extra request: %s", r.URL.RequestURI())
		}
	})

	got, err := ListLabels(context.Background(), c, "12345", 2)
	if err != nil {
		t.Fatalf("ListLabels: %v", err)
	}
	if len(requests) != 2 {
		t.Fatalf("requests = %d, want 2 (%v)", len(requests), requests)
	}
	if len(got) != 2 || got[0].Name != "first" || got[1].Name != "second" {
		t.Fatalf("unexpected labels: %+v", got)
	}
}

func TestListLabelsCloudSupportsDocumentedV2Targets(t *testing.T) {
	cases := []struct {
		name     string
		target   LabelTarget
		wantPath string
	}{
		{name: "page", target: LabelTarget{Type: "page", ID: "12345"}, wantPath: "/wiki/api/v2/pages/12345/labels"},
		{name: "blogpost", target: LabelTarget{Type: "blogpost", ID: "67890"}, wantPath: "/wiki/api/v2/blogposts/67890/labels"},
		{name: "attachment", target: LabelTarget{Type: "attachment", ID: "att1"}, wantPath: "/wiki/api/v2/attachments/att1/labels"},
		{name: "custom content", target: LabelTarget{Type: "custom-content", ID: "custom1"}, wantPath: "/wiki/api/v2/custom-content/custom1/labels"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != tc.wantPath {
					t.Fatalf("path: %s", r.URL.Path)
				}
				if got := r.URL.Query().Get("prefix"); got != "global" {
					t.Fatalf("prefix: %q", got)
				}
				_ = json.NewEncoder(w).Encode(map[string]any{
					"results": []Label{{ID: "label1", Name: "targeted", Prefix: "global"}},
					"_links":  map[string]any{},
				})
			})

			got, err := ListTargetLabels(context.Background(), c, tc.target, LabelListOptions{Limit: 2, Prefix: "global"})
			if err != nil {
				t.Fatalf("ListTargetLabels: %v", err)
			}
			if len(got) != 1 || got[0].Name != "targeted" {
				t.Fatalf("unexpected labels: %+v", got)
			}
		})
	}
}

func TestListTargetLabelsServerUsesContentLabelEndpointForAnyContentTarget(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorServer, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/content/12345/label" {
			t.Fatalf("path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"results": []Label{{Name: "server-content", Prefix: "global"}},
		})
	})

	got, err := ListTargetLabels(context.Background(), c, LabelTarget{Type: "blogpost", ID: "12345"}, LabelListOptions{Limit: 2})
	if err != nil {
		t.Fatalf("ListTargetLabels: %v", err)
	}
	if len(got) != 1 || got[0].Name != "server-content" {
		t.Fatalf("unexpected labels: %+v", got)
	}
}

func TestListSpaceLabelsCloudUsesContentLabelEndpointAndResolvesSpaceKey(t *testing.T) {
	requests := make([]string, 0, 2)
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.URL.RequestURI())
		switch r.URL.Path {
		case "/wiki/api/v2/spaces":
			if got := r.URL.Query().Get("keys"); got != "ENG" {
				t.Fatalf("keys: %q", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{"id": "42", "key": "ENG", "name": "Engineering"}},
			})
		case "/wiki/api/v2/spaces/42/content/labels":
			if got := r.URL.Query().Get("limit"); got != "2" {
				t.Fatalf("limit: %q", got)
			}
			if got := r.URL.Query().Get("prefix"); got != "global" {
				t.Fatalf("prefix: %q", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []Label{{ID: "label1", Name: "incident", Prefix: "global"}},
				"_links":  map[string]any{},
			})
		default:
			t.Fatalf("path: %s", r.URL.Path)
		}
	})

	got, err := ListSpaceLabels(context.Background(), c, "ENG", LabelListOptions{Limit: 2, Prefix: "global"})
	if err != nil {
		t.Fatalf("ListSpaceLabels: %v", err)
	}
	if len(requests) != 2 {
		t.Fatalf("requests = %d, want 2 (%v)", len(requests), requests)
	}
	if len(got) != 1 || got[0].Name != "incident" {
		t.Fatalf("unexpected labels: %+v", got)
	}
}

func TestListSpaceLabelsCloudCanReadSpaceEntityLabels(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/wiki/api/v2/spaces":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{"id": "42", "key": "ENG", "name": "Engineering"}},
			})
		case "/wiki/api/v2/spaces/42/labels":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []Label{{ID: "label1", Name: "space-label", Prefix: "global"}},
				"_links":  map[string]any{},
			})
		default:
			t.Fatalf("path: %s", r.URL.Path)
		}
	})

	got, err := ListSpaceLabels(context.Background(), c, "ENG", LabelListOptions{Limit: 2, Scope: "space"})
	if err != nil {
		t.Fatalf("ListSpaceLabels: %v", err)
	}
	if len(got) != 1 || got[0].Name != "space-label" {
		t.Fatalf("unexpected labels: %+v", got)
	}
}

func TestListSpaceLabelsServerUsesSpaceLabelEndpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorServer, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/space/ENG/labels" {
			t.Fatalf("path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("limit"); got != "2" {
			t.Fatalf("limit: %q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"results": []Label{{Name: "server-space", Prefix: "global"}},
		})
	})

	got, err := ListSpaceLabels(context.Background(), c, "ENG", LabelListOptions{Limit: 2})
	if err != nil {
		t.Fatalf("ListSpaceLabels: %v", err)
	}
	if len(got) != 1 || got[0].Name != "server-space" {
		t.Fatalf("unexpected labels: %+v", got)
	}
}

func TestSearchLabelsCloudUsesV2CatalogFilters(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/wiki/api/v2/labels" {
			t.Fatalf("path: %s", r.URL.Path)
		}
		if got := r.URL.Query()["label-id"]; len(got) != 2 || got[0] != "1" || got[1] != "2" {
			t.Fatalf("label-id: %#v", got)
		}
		if got := r.URL.Query()["prefix"]; len(got) != 2 || got[0] != "global" || got[1] != "my" {
			t.Fatalf("prefix: %#v", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"results": []Label{{ID: "1", Name: "cloud", Prefix: "global"}},
			"_links":  map[string]any{},
		})
	})

	got, err := SearchLabels(context.Background(), c, LabelSearchOptions{
		Limit:    2,
		LabelIDs: []string{"1", "2"},
		Prefixes: []string{"global", "my"},
	})
	if err != nil {
		t.Fatalf("SearchLabels: %v", err)
	}
	if len(got) != 1 || got[0].ID != "1" {
		t.Fatalf("unexpected labels: %+v", got)
	}
}

func TestSearchLabelsServerUnsupported(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorServer, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s", r.URL.RequestURI())
	})

	_, err := SearchLabels(context.Background(), c, LabelSearchOptions{Limit: 2})
	if err == nil || !strings.Contains(err.Error(), "Confluence Cloud only") {
		t.Fatalf("err = %v, want Cloud-only error", err)
	}
}

func TestListRecentLabelsServerUsesRecentEndpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorServer, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/label/recent" {
			t.Fatalf("path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("limit"); got != "2" {
			t.Fatalf("limit: %q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"results": []Label{{Name: "recent", Prefix: "global"}},
		})
	})

	got, err := ListRecentLabels(context.Background(), c, 2)
	if err != nil {
		t.Fatalf("ListRecentLabels: %v", err)
	}
	if len(got) != 1 || got[0].Name != "recent" {
		t.Fatalf("unexpected labels: %+v", got)
	}
}

func TestListRelatedLabelsServerUsesGlobalAndSpaceEndpoints(t *testing.T) {
	cases := []struct {
		name     string
		opts     LabelRelatedOptions
		wantPath string
	}{
		{name: "global", opts: LabelRelatedOptions{Label: "incident", Limit: 2}, wantPath: "/rest/api/label/incident/related"},
		{name: "space", opts: LabelRelatedOptions{SpaceKey: "ENG", Label: "incident", Limit: 2}, wantPath: "/rest/api/space/ENG/labels/incident/related"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, c := testClientWithFlavor(t, client.FlavorServer, func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != tc.wantPath {
					t.Fatalf("path: %s", r.URL.Path)
				}
				if got := r.URL.Query().Get("limit"); got != "2" {
					t.Fatalf("limit: %q", got)
				}
				_ = json.NewEncoder(w).Encode(map[string]any{
					"results": []Label{{Name: "related", Prefix: "global"}},
				})
			})

			got, err := ListRelatedLabels(context.Background(), c, tc.opts)
			if err != nil {
				t.Fatalf("ListRelatedLabels: %v", err)
			}
			if len(got) != 1 || got[0].Name != "related" {
				t.Fatalf("unexpected labels: %+v", got)
			}
		})
	}
}
