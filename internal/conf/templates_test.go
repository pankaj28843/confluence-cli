package conf

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

func TestListContentTemplatesCloudUsesV1EndpointAndPagination(t *testing.T) {
	requests := make([]string, 0, 2)
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.URL.RequestURI())
		w.Header().Set("Content-Type", "application/json")

		switch len(requests) {
		case 1:
			if r.Method != http.MethodGet || r.URL.Path != "/wiki/rest/api/template/page" {
				t.Fatalf("first request: %s %s", r.Method, r.URL.RequestURI())
			}
			q := r.URL.Query()
			if got := q.Get("spaceKey"); got != "ENG" {
				t.Fatalf("spaceKey = %q, want ENG", got)
			}
			if got := q.Get("limit"); got != "2" {
				t.Fatalf("limit = %q, want 2", got)
			}
			if got := q.Get("start"); got != "0" {
				t.Fatalf("start = %q, want 0", got)
			}
			if got := q["expand"]; len(got) != 2 || got[0] != "body.storage" || got[1] != "space" {
				t.Fatalf("expand = %v, want [body.storage space]", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{
					"templateId":   "tpl1",
					"name":         "Incident report",
					"description":  "Capture impact",
					"templateType": "page",
					"body": map[string]any{
						"storage": map[string]any{
							"value":          "<p>Impact</p>",
							"representation": "storage",
						},
					},
				}},
				"size":  1,
				"start": 0,
				"limit": 2,
				"_links": map[string]any{
					"next": "/wiki/rest/api/template/page?spaceKey=ENG&start=1&limit=2",
				},
			})
		case 2:
			if r.Method != http.MethodGet || r.URL.Path != "/wiki/rest/api/template/page" {
				t.Fatalf("second request: %s %s", r.Method, r.URL.RequestURI())
			}
			q := r.URL.Query()
			if got := q.Get("spaceKey"); got != "ENG" {
				t.Fatalf("second spaceKey = %q, want ENG", got)
			}
			if got := q.Get("start"); got != "1" {
				t.Fatalf("second start = %q, want 1", got)
			}
			if got := q.Get("limit"); got != "2" {
				t.Fatalf("second limit = %q, want 2", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{
					"templateId":   "tpl2",
					"name":         "Decision",
					"templateType": "page",
				}},
				"size":   1,
				"_links": map[string]any{},
			})
		default:
			t.Fatalf("unexpected extra request: %s", r.URL.RequestURI())
		}
	})

	got, err := ListContentTemplates(context.Background(), c, TemplateListOptions{
		SpaceKey: "ENG",
		Expand:   []string{"body.storage", "space"},
		Limit:    2,
	})
	if err != nil {
		t.Fatalf("ListContentTemplates: %v", err)
	}
	if len(requests) != 2 {
		t.Fatalf("requests = %v, want two paginated requests", requests)
	}
	if len(got) != 2 {
		t.Fatalf("len(templates) = %d, want 2 (%+v)", len(got), got)
	}
	if got[0].TemplateID != "tpl1" || got[0].Name != "Incident report" || got[0].Body.Storage.Value != "<p>Impact</p>" {
		t.Fatalf("unexpected first template: %+v", got[0])
	}
	if got[1].TemplateID != "tpl2" || got[1].Name != "Decision" {
		t.Fatalf("unexpected second template: %+v", got[1])
	}
}

func TestListBlueprintTemplatesCloudUsesBlueprintEndpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/wiki/rest/api/template/blueprint" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		q := r.URL.Query()
		if got := q.Get("spaceKey"); got != "ENG" {
			t.Fatalf("spaceKey = %q, want ENG", got)
		}
		if got := q.Get("limit"); got != "1" {
			t.Fatalf("limit = %q, want 1", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"results": []map[string]any{{
				"templateId":           "bp1",
				"name":                 "Product requirements",
				"referencingBlueprint": "com.example.blueprint",
				"templateType":         "page",
				"originalTemplate": map[string]any{
					"pluginKey": "com.example",
					"moduleKey": "requirements",
				},
			}},
			"size":   1,
			"_links": map[string]any{},
		})
	})

	got, err := ListBlueprintTemplates(context.Background(), c, TemplateListOptions{SpaceKey: "ENG", Limit: 1})
	if err != nil {
		t.Fatalf("ListBlueprintTemplates: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len(templates) = %d, want 1", len(got))
	}
	if got[0].TemplateID != "bp1" || got[0].OriginalTemplate.ModuleKey != "requirements" || got[0].ReferencingBlueprint == "" {
		t.Fatalf("unexpected blueprint template: %+v", got[0])
	}
}

func TestGetContentTemplateCloudUsesV1Endpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/wiki/rest/api/template/tpl1" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		if got := r.URL.Query()["expand"]; len(got) != 1 || got[0] != "body.storage" {
			t.Fatalf("expand = %v, want [body.storage]", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"templateId":   "tpl1",
			"name":         "Incident report",
			"description":  "Capture impact",
			"templateType": "page",
			"body": map[string]any{
				"storage": map[string]any{
					"value":          "<p>Impact</p>",
					"representation": "storage",
				},
			},
			"labels": []map[string]any{{"name": "incident", "label": "incident"}},
		})
	})

	got, err := GetContentTemplate(context.Background(), c, "tpl1", []string{"body.storage"})
	if err != nil {
		t.Fatalf("GetContentTemplate: %v", err)
	}
	if got.TemplateID != "tpl1" || got.Name != "Incident report" || got.Body.Storage.Value != "<p>Impact</p>" {
		t.Fatalf("unexpected template: %+v", got)
	}
	if len(got.Labels) != 1 || got.Labels[0].Name != "incident" {
		t.Fatalf("labels = %+v, want incident", got.Labels)
	}
}

func TestTemplatesRejectServerFlavor(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorServer, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected server request: %s %s", r.Method, r.URL.RequestURI())
	})

	if _, err := ListContentTemplates(context.Background(), c, TemplateListOptions{Limit: 1}); err == nil || !strings.Contains(err.Error(), "only supported on Confluence Cloud") {
		t.Fatalf("ListContentTemplates error = %v, want Cloud-only error", err)
	}
	if _, err := ListBlueprintTemplates(context.Background(), c, TemplateListOptions{Limit: 1}); err == nil || !strings.Contains(err.Error(), "only supported on Confluence Cloud") {
		t.Fatalf("ListBlueprintTemplates error = %v, want Cloud-only error", err)
	}
	if _, err := GetContentTemplate(context.Background(), c, "tpl1", nil); err == nil || !strings.Contains(err.Error(), "only supported on Confluence Cloud") {
		t.Fatalf("GetContentTemplate error = %v, want Cloud-only error", err)
	}
}
