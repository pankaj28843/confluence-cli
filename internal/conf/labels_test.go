package conf

import (
	"context"
	"encoding/json"
	"net/http"
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
