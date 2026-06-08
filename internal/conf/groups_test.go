package conf

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

func TestListGroupsCloudPaginatesAndFiltersAccessType(t *testing.T) {
	var seen []string
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		seen = append(seen, r.URL.String())
		if r.URL.Path != "/wiki/rest/api/group" {
			t.Fatalf("path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("accessType"); got != "admin" {
			t.Fatalf("accessType: %q", got)
		}
		if got := r.URL.Query().Get("limit"); got != "2" {
			t.Fatalf("limit: %q", got)
		}
		switch r.URL.Query().Get("start") {
		case "0":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []Group{{ID: "grp-1", Name: "engineering", Type: "group", UsageType: "global"}},
				"_links":  map[string]any{"next": "/wiki/rest/api/group?start=1&limit=2&accessType=admin"},
			})
		case "1":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []Group{{ID: "grp-2", Name: "platform", Type: "group"}},
				"_links":  map[string]any{},
			})
		default:
			t.Fatalf("start: %q", r.URL.Query().Get("start"))
		}
	})

	got, err := ListGroupsWithOptions(context.Background(), c, GroupListOptions{Limit: 2, AccessType: "admin"})
	if err != nil {
		t.Fatalf("ListGroupsWithOptions: %v", err)
	}
	if len(got) != 2 || got[0].ID != "grp-1" || got[1].Name != "platform" {
		t.Fatalf("groups = %+v", got)
	}
	if len(seen) != 2 {
		t.Fatalf("requests = %v", seen)
	}
}

func TestGetGroupUsesFlavorSpecificLookup(t *testing.T) {
	_, cloud := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/wiki/rest/api/group/by-id" {
			t.Fatalf("cloud path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("id"); got != "grp-1" {
			t.Fatalf("cloud id: %q", got)
		}
		_ = json.NewEncoder(w).Encode(Group{ID: "grp-1", Name: "engineering", Type: "group"})
	})
	group, err := GetGroup(context.Background(), cloud, GroupLookupOptions{ID: "grp-1"})
	if err != nil {
		t.Fatalf("GetGroup cloud: %v", err)
	}
	if group.ID != "grp-1" || group.Name != "engineering" {
		t.Fatalf("cloud group = %+v", group)
	}

	_, server := testClientWithFlavor(t, client.FlavorServer, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/group/engineering" {
			t.Fatalf("server path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("expand"); got != "members" {
			t.Fatalf("server expand: %q", got)
		}
		_ = json.NewEncoder(w).Encode(Group{Name: "engineering", Type: "group"})
	})
	group, err = GetGroup(context.Background(), server, GroupLookupOptions{Name: "engineering", Expand: "members"})
	if err != nil {
		t.Fatalf("GetGroup server: %v", err)
	}
	if group.Name != "engineering" {
		t.Fatalf("server group = %+v", group)
	}
}

func TestListGroupMembersUsesDocumentedFlavorRoutes(t *testing.T) {
	_, cloud := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/wiki/rest/api/group/grp-1/membersByGroupId" {
			t.Fatalf("cloud path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("start") != "0" || q.Get("limit") != "2" || q.Get("shouldReturnTotalSize") != "true" {
			t.Fatalf("cloud query: %s", r.URL.RawQuery)
		}
		if got := q.Get("expand"); got != "operations,personalSpace" {
			t.Fatalf("cloud expand: %q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"results": []User{{AccountID: "acct-1", DisplayName: "Ada Lovelace"}},
			"_links":  map[string]any{},
		})
	})
	users, err := ListGroupMembersWithOptions(context.Background(), cloud, GroupMemberOptions{
		GroupID:               "grp-1",
		Limit:                 2,
		Expand:                []string{"operations", "personalSpace"},
		ShouldReturnTotalSize: true,
	})
	if err != nil {
		t.Fatalf("ListGroupMembersWithOptions cloud: %v", err)
	}
	if len(users) != 1 || users[0].AccountID != "acct-1" {
		t.Fatalf("cloud users = %+v", users)
	}

	_, server := testClientWithFlavor(t, client.FlavorServer, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/group/engineering/member" {
			t.Fatalf("server path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("start") != "0" || q.Get("limit") != "2" || q.Get("expand") != "operations" {
			t.Fatalf("server query: %s", r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"results": []User{{Username: "ada", DisplayName: "Ada Lovelace"}},
			"_links":  map[string]any{},
		})
	})
	users, err = ListGroupMembersWithOptions(context.Background(), server, GroupMemberOptions{
		GroupName: "engineering",
		Limit:     2,
		Expand:    []string{"operations"},
	})
	if err != nil {
		t.Fatalf("ListGroupMembersWithOptions server: %v", err)
	}
	if len(users) != 1 || users[0].Username != "ada" {
		t.Fatalf("server users = %+v", users)
	}
}

func TestPickGroupsCloudUsesPicker(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/wiki/rest/api/group/picker" {
			t.Fatalf("path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("query") != "eng" || q.Get("start") != "0" || q.Get("limit") != "3" || q.Get("shouldReturnTotalSize") != "true" {
			t.Fatalf("query: %s", r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"results": []Group{{ID: "grp-1", Name: "engineering"}},
			"_links":  map[string]any{},
		})
	})
	got, err := PickGroups(context.Background(), c, "eng", GroupPickerOptions{Limit: 3, ShouldReturnTotalSize: true})
	if err != nil {
		t.Fatalf("PickGroups: %v", err)
	}
	if len(got) != 1 || got[0].Name != "engineering" {
		t.Fatalf("groups = %+v", got)
	}
}

func TestGroupRelationsServerOnlyUseDocumentedRoutes(t *testing.T) {
	paths := []string{
		"/rest/api/group/engineering/groupmember",
		"/rest/api/group/engineering/groupparent",
		"/rest/api/group/engineering/groupancestor",
	}
	var index int
	_, c := testClientWithFlavor(t, client.FlavorServer, func(w http.ResponseWriter, r *http.Request) {
		if index >= len(paths) {
			t.Fatalf("unexpected request %s", r.URL.String())
		}
		if r.URL.Path != paths[index] {
			t.Fatalf("path %d: %s", index, r.URL.Path)
		}
		if r.URL.Query().Get("start") != "0" || r.URL.Query().Get("limit") != "4" || r.URL.Query().Get("expand") != "members" {
			t.Fatalf("query %d: %s", index, r.URL.RawQuery)
		}
		index++
		_ = json.NewEncoder(w).Encode(map[string]any{
			"results": []Group{{Name: "child"}},
			"_links":  map[string]any{},
		})
	})

	opts := GroupRelationOptions{GroupName: "engineering", Limit: 4, Expand: "members"}
	if got, err := ListGroupChildGroups(context.Background(), c, opts); err != nil || len(got) != 1 {
		t.Fatalf("ListGroupChildGroups: %v %+v", err, got)
	}
	if got, err := ListGroupParents(context.Background(), c, opts); err != nil || len(got) != 1 {
		t.Fatalf("ListGroupParents: %v %+v", err, got)
	}
	if got, err := ListGroupAncestors(context.Background(), c, opts); err != nil || len(got) != 1 {
		t.Fatalf("ListGroupAncestors: %v %+v", err, got)
	}
	if index != len(paths) {
		t.Fatalf("requests: %d", index)
	}
}

func TestGroupHelpersRejectUnsupportedInputs(t *testing.T) {
	_, cloud := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected cloud request: %s", r.URL.String())
	})
	_, server := testClientWithFlavor(t, client.FlavorServer, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected server request: %s", r.URL.String())
	})

	cases := []struct {
		name string
		err  error
	}{
		{name: "cloud group view requires id", err: firstGroupError(GetGroup(context.Background(), cloud, GroupLookupOptions{Name: "engineering"}))},
		{name: "server group view requires name", err: firstGroupError(GetGroup(context.Background(), server, GroupLookupOptions{ID: "grp-1"}))},
		{name: "cloud child groups unsupported", err: firstGroupSliceError(ListGroupChildGroups(context.Background(), cloud, GroupRelationOptions{GroupName: "engineering"}))},
		{name: "server picker unsupported", err: firstGroupSliceError(PickGroups(context.Background(), server, "eng", GroupPickerOptions{}))},
	}
	for _, tc := range cases {
		if tc.err == nil || !strings.Contains(tc.err.Error(), "unsupported") && !strings.Contains(tc.err.Error(), "required") {
			t.Fatalf("%s error = %v", tc.name, tc.err)
		}
	}
}

func firstGroupError(_ *Group, err error) error {
	return err
}

func firstGroupSliceError(_ []Group, err error) error {
	return err
}
