package conf

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

func TestListOperationsCloudUsesV2Endpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/wiki/api/v2/pages/12345/operations" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"operations": []map[string]any{
				{"operation": "read", "targetType": "page"},
				{"operation": "update", "targetType": "page"},
			},
		})
	})

	got, err := ListOperations(context.Background(), c, "page", "12345")
	if err != nil {
		t.Fatalf("ListOperations: %v", err)
	}
	if len(got) != 2 || got[0].Operation != "read" || got[1].TargetType != "page" {
		t.Fatalf("unexpected operations: %+v", got)
	}
}

func TestListOperationsCloudResolvesSpaceKey(t *testing.T) {
	var requests []string
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.Method+" "+r.URL.RequestURI())
		switch len(requests) {
		case 1:
			if r.Method != http.MethodGet || r.URL.Path != "/wiki/api/v2/spaces" {
				t.Fatalf("space lookup: %s %s", r.Method, r.URL.RequestURI())
			}
			if got := r.URL.Query().Get("keys"); got != "ENG" {
				t.Fatalf("keys = %q, want ENG", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{"id": "1001", "key": "ENG", "name": "Engineering"}},
			})
		case 2:
			if r.Method != http.MethodGet || r.URL.Path != "/wiki/api/v2/spaces/1001/operations" {
				t.Fatalf("operations request: %s %s", r.Method, r.URL.RequestURI())
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"operations": []map[string]any{{"operation": "read", "targetType": "space"}},
			})
		default:
			t.Fatalf("unexpected extra request: %s", r.URL.RequestURI())
		}
	})

	got, err := ListOperations(context.Background(), c, "space", "ENG")
	if err != nil {
		t.Fatalf("ListOperations: %v", err)
	}
	if len(got) != 1 || got[0].TargetType != "space" {
		t.Fatalf("unexpected operations: %+v", got)
	}
	if len(requests) != 2 {
		t.Fatalf("requests = %v, want lookup then operations", requests)
	}
}

func TestListOperationsServerUsesContentExpand(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorServer, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/rest/api/content/12345" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		if got := r.URL.Query().Get("expand"); got != "operations" {
			t.Fatalf("expand = %q, want operations", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":   "12345",
			"type": "page",
			"operations": []map[string]any{
				{"operation": "read", "targetType": "page"},
				{"operation": "update", "targetType": "page"},
			},
		})
	})

	got, err := ListOperations(context.Background(), c, "page", "12345")
	if err != nil {
		t.Fatalf("ListOperations: %v", err)
	}
	if len(got) != 2 || got[1].Operation != "update" {
		t.Fatalf("unexpected operations: %+v", got)
	}
}

func TestGetLikeCountCloudUsesV2Endpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/wiki/api/v2/pages/12345/likes/count" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"count": 16})
	})

	got, err := GetLikeCount(context.Background(), c, "page", "12345")
	if err != nil {
		t.Fatalf("GetLikeCount: %v", err)
	}
	if got.Count != 16 {
		t.Fatalf("count = %d, want 16", got.Count)
	}
}

func TestListLikeUsersCloudFollowsLinkPagination(t *testing.T) {
	var requests []string
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.Method+" "+r.URL.RequestURI())
		switch len(requests) {
		case 1:
			if r.Method != http.MethodGet || r.URL.Path != "/wiki/api/v2/pages/12345/likes/users" {
				t.Fatalf("first request: %s %s", r.Method, r.URL.RequestURI())
			}
			if got := r.URL.Query().Get("limit"); got != "2" {
				t.Fatalf("first limit = %q, want 2", got)
			}
			w.Header().Set("Link", `</wiki/api/v2/pages/12345/likes/users?cursor=abc&limit=2>; rel="next"`)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{"accountId": "a-1"}},
				"_links":  map[string]any{"next": "/api/v2/pages/12345/likes/users?cursor=abc&limit=2"},
			})
		case 2:
			if r.Method != http.MethodGet || r.URL.Path != "/wiki/api/v2/pages/12345/likes/users" {
				t.Fatalf("second request: %s %s", r.Method, r.URL.RequestURI())
			}
			if got := r.URL.Query().Get("cursor"); got != "abc" {
				t.Fatalf("cursor = %q, want abc", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{"accountId": "a-2"}},
				"_links":  map[string]any{},
			})
		default:
			t.Fatalf("unexpected extra request: %s", r.URL.RequestURI())
		}
	})

	got, err := ListLikeUsers(context.Background(), c, "page", "12345", 2)
	if err != nil {
		t.Fatalf("ListLikeUsers: %v", err)
	}
	if len(got) != 2 || got[0].AccountID != "a-1" || got[1].AccountID != "a-2" {
		t.Fatalf("unexpected users: %+v", got)
	}
}

func TestGetLikeCountServerUnsupported(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorServer, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
	})

	if _, err := GetLikeCount(context.Background(), c, "page", "12345"); err == nil {
		t.Fatal("GetLikeCount succeeded on Server/Data Center, want unsupported error")
	}
}
