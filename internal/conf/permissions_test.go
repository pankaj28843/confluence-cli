package conf

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

func TestListSpacePermissionAssignmentsCloudUsesV2EndpointAndResolvesSpaceKey(t *testing.T) {
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
		case "/wiki/api/v2/spaces/42/permissions":
			if got := r.URL.Query().Get("limit"); got != "2" {
				t.Fatalf("limit: %q", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []SpacePermissionAssignment{
					{ID: "perm-1", Principal: SpacePermissionPrincipal{Type: "user", ID: "acct-1"}, Operation: SpacePermissionOperation{Key: "read", TargetType: "page"}},
				},
				"_links": map[string]any{},
			})
		default:
			t.Fatalf("path: %s", r.URL.Path)
		}
	})

	got, err := ListSpacePermissionAssignments(context.Background(), c, "ENG", 2)
	if err != nil {
		t.Fatalf("ListSpacePermissionAssignments: %v", err)
	}
	if len(requests) != 2 {
		t.Fatalf("requests = %d, want 2 (%v)", len(requests), requests)
	}
	if len(got) != 1 || got[0].Principal.ID != "acct-1" || got[0].Operation.Name() != "read" {
		t.Fatalf("unexpected permissions: %+v", got)
	}
}

func TestListSpacePermissionAssignmentsServerUsesSpacePermissionsEndpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorServer, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/space/ENG/permissions" {
			t.Fatalf("path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("limit"); got != "" {
			t.Fatalf("server list should not send limit query, got %q", got)
		}
		_ = json.NewEncoder(w).Encode([]SpacePermissionAssignment{
			{Subject: SpacePermissionSubject{Type: "user", DisplayName: "Ada Lovelace"}, Operation: SpacePermissionOperation{OperationKey: "read", TargetType: "space"}, SpaceKey: "ENG", SpaceID: 42},
			{Subject: SpacePermissionSubject{Type: "group", DisplayName: "engineering"}, Operation: SpacePermissionOperation{OperationKey: "administer", TargetType: "space"}, SpaceKey: "ENG", SpaceID: 42},
		})
	})

	got, err := ListSpacePermissionAssignments(context.Background(), c, "ENG", 1)
	if err != nil {
		t.Fatalf("ListSpacePermissionAssignments: %v", err)
	}
	if len(got) != 1 || got[0].Subject.DisplayName != "Ada Lovelace" || got[0].Operation.Name() != "read" {
		t.Fatalf("unexpected permissions: %+v", got)
	}
}

func TestListAvailableSpacePermissionsCloudUsesV2Catalog(t *testing.T) {
	requests := make([]string, 0, 2)
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.URL.RequestURI())
		switch len(requests) {
		case 1:
			if r.URL.RequestURI() != "/wiki/api/v2/space-permissions?limit=2" {
				t.Fatalf("first request: %s", r.URL.RequestURI())
			}
			w.Header().Set("Link", "<http://"+r.Host+"/wiki/api/v2/space-permissions?cursor=abc>; rel=\"next\"")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []SpacePermissionDefinition{{ID: "read-space", DisplayName: "View space"}},
				"_links":  map[string]any{},
			})
		case 2:
			if r.URL.RequestURI() != "/wiki/api/v2/space-permissions?cursor=abc" {
				t.Fatalf("second request: %s", r.URL.RequestURI())
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []SpacePermissionDefinition{{ID: "update-page", DisplayName: "Update page"}},
				"_links":  map[string]any{},
			})
		default:
			t.Fatalf("unexpected extra request: %s", r.URL.RequestURI())
		}
	})

	got, err := ListAvailableSpacePermissions(context.Background(), c, 2)
	if err != nil {
		t.Fatalf("ListAvailableSpacePermissions: %v", err)
	}
	if len(requests) != 2 {
		t.Fatalf("requests = %d, want 2 (%v)", len(requests), requests)
	}
	if len(got) != 2 || got[0].ID != "read-space" || got[1].ID != "update-page" {
		t.Fatalf("unexpected permissions: %+v", got)
	}
}

func TestListSpacePermissionsForSubjectServerRoutesBySelector(t *testing.T) {
	cases := []struct {
		name     string
		selector SpacePermissionSubjectSelector
		wantPath string
	}{
		{name: "anonymous", selector: SpacePermissionSubjectSelector{Anonymous: true}, wantPath: "/rest/api/space/ENG/permissions/anonymous"},
		{name: "group", selector: SpacePermissionSubjectSelector{GroupName: "eng team"}, wantPath: "/rest/api/space/ENG/permissions/group/eng%20team"},
		{name: "user key", selector: SpacePermissionSubjectSelector{UserKey: "ada"}, wantPath: "/rest/api/space/ENG/permissions/user/ada"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, c := testClientWithFlavor(t, client.FlavorServer, func(w http.ResponseWriter, r *http.Request) {
				if r.URL.EscapedPath() != tc.wantPath {
					t.Fatalf("path: %s (escaped %s)", r.URL.Path, r.URL.EscapedPath())
				}
				_ = json.NewEncoder(w).Encode([]SpacePermissionAssignment{
					{Subject: SpacePermissionSubject{Type: "user", DisplayName: "Ada Lovelace"}, Operation: SpacePermissionOperation{OperationKey: "read", TargetType: "space"}},
				})
			})

			got, err := ListSpacePermissionsForSubject(context.Background(), c, "ENG", tc.selector, 25)
			if err != nil {
				t.Fatalf("ListSpacePermissionsForSubject: %v", err)
			}
			if len(got) != 1 || got[0].Operation.Name() != "read" {
				t.Fatalf("unexpected permissions: %+v", got)
			}
		})
	}
}

func TestSpacePermissionHelpersValidateInputs(t *testing.T) {
	_, server := testClientWithFlavor(t, client.FlavorServer, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s", r.URL.RequestURI())
	})
	_, cloud := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s", r.URL.RequestURI())
	})

	if _, err := ListSpacePermissionAssignments(context.Background(), server, "", 25); err == nil || !strings.Contains(err.Error(), "space") {
		t.Fatalf("missing space error = %v", err)
	}
	if _, err := ListAvailableSpacePermissions(context.Background(), server, 25); err == nil || !strings.Contains(err.Error(), "Cloud") {
		t.Fatalf("server available error = %v", err)
	}
	if _, err := ListSpacePermissionsForSubject(context.Background(), cloud, "ENG", SpacePermissionSubjectSelector{Anonymous: true}, 25); err == nil || !strings.Contains(err.Error(), "Server") {
		t.Fatalf("cloud subject error = %v", err)
	}
	if _, err := ListSpacePermissionsForSubject(context.Background(), server, "ENG", SpacePermissionSubjectSelector{}, 25); err == nil || !strings.Contains(err.Error(), "exactly one") {
		t.Fatalf("missing selector error = %v", err)
	}
	if _, err := ListSpacePermissionsForSubject(context.Background(), server, "ENG", SpacePermissionSubjectSelector{Anonymous: true, GroupName: "eng"}, 25); err == nil || !strings.Contains(err.Error(), "exactly one") {
		t.Fatalf("multiple selector error = %v", err)
	}
}
