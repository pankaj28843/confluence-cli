package conf

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

func TestListContentRestrictionsServerUsesByOperationEndpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorServer, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/rest/api/content/12345/restriction/byOperation" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		if got := r.URL.Query().Get("expand"); got != "restrictions.user,restrictions.group" {
			t.Fatalf("expand = %q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"restrictions": map[string]any{
				"read": map[string]any{
					"operation": "read",
					"restrictions": map[string]any{
						"user": map[string]any{
							"results": []map[string]any{{"displayName": "Alice Example"}},
							"start":   0,
							"limit":   25,
							"size":    1,
						},
						"group": map[string]any{
							"results": []map[string]any{{"name": "eng"}, {"name": "ops"}},
							"start":   0,
							"limit":   25,
							"size":    2,
						},
					},
				},
				"update": map[string]any{
					"operation": "update",
					"restrictions": map[string]any{
						"user":  map[string]any{"results": []map[string]any{}, "size": 0},
						"group": map[string]any{"results": []map[string]any{}, "size": 0},
					},
				},
			},
			"get_links": map[string]any{"self": "/rest/api/content/12345/restriction/byOperation"},
		})
	})

	got, err := ListContentRestrictions(context.Background(), c, "12345")
	if err != nil {
		t.Fatalf("ListContentRestrictions: %v", err)
	}
	if len(got.Restrictions) != 2 {
		t.Fatalf("restrictions = %+v", got.Restrictions)
	}
	read := got.Restrictions["read"]
	if read.Operation != "read" || read.UserCount() != 1 || read.GroupCount() != 2 {
		t.Fatalf("read restriction = %+v", read)
	}
}

func TestListContentRestrictionsCloudUsesByOperationEndpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/wiki/rest/api/content/12345/restriction/byOperation" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		if got := r.URL.Query().Get("expand"); got != "restrictions.user,restrictions.group" {
			t.Fatalf("expand = %q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"read": map[string]any{
				"operation": "read",
				"restrictions": map[string]any{
					"user":  map[string]any{"results": []map[string]any{{"accountId": "abc"}}, "size": 1},
					"group": map[string]any{"results": []map[string]any{}, "size": 0},
				},
			},
			"update": map[string]any{
				"operation": "update",
				"restrictions": map[string]any{
					"user":  map[string]any{"results": []map[string]any{}, "size": 0},
					"group": map[string]any{"results": []map[string]any{}, "size": 0},
				},
			},
			"_links": map[string]any{"self": "/wiki/rest/api/content/12345/restriction/byOperation"},
		})
	})

	got, err := ListContentRestrictions(context.Background(), c, "12345")
	if err != nil {
		t.Fatalf("ListContentRestrictions: %v", err)
	}
	if got.Restrictions["read"].UserCount() != 1 || got.Restrictions["update"].GroupCount() != 0 {
		t.Fatalf("restrictions = %+v", got.Restrictions)
	}
}

func TestGetContentRestrictionForOperationUsesOperationEndpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/wiki/rest/api/content/12345/restriction/byOperation/read" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		if got := r.URL.Query().Get("start"); got != "0" {
			t.Fatalf("start = %q", got)
		}
		if got := r.URL.Query().Get("limit"); got != "7" {
			t.Fatalf("limit = %q", got)
		}
		if got := r.URL.Query().Get("expand"); got != "restrictions.user,restrictions.group" {
			t.Fatalf("expand = %q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"operation": "read",
			"restrictions": map[string]any{
				"user":  map[string]any{"results": []map[string]any{{"accountId": "abc"}}, "size": 1},
				"group": map[string]any{"results": []map[string]any{}, "size": 0},
			},
		})
	})

	got, err := GetContentRestrictionForOperation(context.Background(), c, "12345", "read", 7)
	if err != nil {
		t.Fatalf("GetContentRestrictionForOperation: %v", err)
	}
	if got.Operation != "read" || got.UserCount() != 1 {
		t.Fatalf("restriction = %+v", got)
	}
}

func TestRestrictionHelpersValidateRequiredInputs(t *testing.T) {
	if _, err := ListContentRestrictions(context.Background(), nil, ""); err == nil {
		t.Fatal("ListContentRestrictions should require content id")
	}
	if _, err := GetContentRestrictionForOperation(context.Background(), nil, "12345", "administer", 25); err == nil {
		t.Fatal("GetContentRestrictionForOperation should reject unsupported operations")
	}
}
