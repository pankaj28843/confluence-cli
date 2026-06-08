package conf

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

func TestGetContentStateCloudUsesV1Endpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/wiki/rest/api/content/123/state" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		if got := r.URL.Query().Get("status"); got != "draft" {
			t.Fatalf("status = %q, want draft", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"contentState": map[string]any{"id": 7, "name": "Ready", "color": "#36B37E"},
			"lastUpdated":  "2026-06-08T10:00:00.000Z",
		})
	})

	got, err := GetContentState(context.Background(), c, "123", "draft")
	if err != nil {
		t.Fatalf("GetContentState: %v", err)
	}
	if got.State == nil || got.State.ID != 7 || got.State.Name != "Ready" || got.LastUpdated == "" {
		t.Fatalf("state = %+v", got)
	}
}

func TestListAvailableContentStatesCloud(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/wiki/rest/api/content/123/state/available" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"spaceContentStates":  []map[string]any{{"id": 1, "name": "In review", "color": "#FFAB00"}},
			"customContentStates": []map[string]any{{"id": 22, "name": "Legal review", "color": "#6554C0"}},
		})
	})

	got, err := ListAvailableContentStates(context.Background(), c, "123")
	if err != nil {
		t.Fatalf("ListAvailableContentStates: %v", err)
	}
	if len(got.SpaceContentStates) != 1 || got.SpaceContentStates[0].Name != "In review" {
		t.Fatalf("space states = %+v", got.SpaceContentStates)
	}
	if len(got.CustomContentStates) != 1 || got.CustomContentStates[0].ID != 22 {
		t.Fatalf("custom states = %+v", got.CustomContentStates)
	}
}

func TestListCustomContentStatesCloud(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/wiki/rest/api/content-states" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		_ = json.NewEncoder(w).Encode([]map[string]any{
			{"id": 30, "name": "My state", "color": "#0052CC"},
		})
	})

	got, err := ListCustomContentStates(context.Background(), c)
	if err != nil {
		t.Fatalf("ListCustomContentStates: %v", err)
	}
	if len(got) != 1 || got[0].ID != 30 || got[0].Name != "My state" {
		t.Fatalf("states = %+v", got)
	}
}

func TestListSpaceContentStatesCloud(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/wiki/rest/api/space/ENG/state" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		_ = json.NewEncoder(w).Encode([]map[string]any{
			{"id": 1, "name": "In progress", "color": "#0052CC"},
		})
	})

	got, err := ListSpaceContentStates(context.Background(), c, "ENG")
	if err != nil {
		t.Fatalf("ListSpaceContentStates: %v", err)
	}
	if len(got) != 1 || got[0].Name != "In progress" {
		t.Fatalf("states = %+v", got)
	}
}

func TestGetContentStateSettingsCloud(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/wiki/rest/api/space/ENG/state/settings" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"contentStatesAllowed":       true,
			"customContentStatesAllowed": false,
			"spaceContentStatesAllowed":  true,
			"spaceContentStates":         []map[string]any{{"id": 2, "name": "Archived", "color": "#6B778C"}},
		})
	})

	got, err := GetContentStateSettings(context.Background(), c, "ENG")
	if err != nil {
		t.Fatalf("GetContentStateSettings: %v", err)
	}
	if !got.ContentStatesAllowed || got.CustomContentStatesAllowed || !got.SpaceContentStatesAllowed {
		t.Fatalf("settings = %+v", got)
	}
	if len(got.SpaceContentStates) != 1 || got.SpaceContentStates[0].ID != 2 {
		t.Fatalf("space states = %+v", got.SpaceContentStates)
	}
}

func TestListContentWithStateCloudPaginates(t *testing.T) {
	requests := make([]string, 0, 2)
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.URL.RequestURI())
		if r.Method != http.MethodGet || r.URL.Path != "/wiki/rest/api/space/ENG/state/content" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		q := r.URL.Query()
		if got := q.Get("state-id"); got != "7" {
			t.Fatalf("state-id = %q, want 7", got)
		}
		if got := q.Get("limit"); got != "2" {
			t.Fatalf("limit = %q, want 2", got)
		}
		if got := q["expand"]; len(got) != 2 || got[0] != "space" || got[1] != "version" {
			t.Fatalf("expand = %v, want [space version]", got)
		}
		switch q.Get("start") {
		case "0":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{"id": "p1", "type": "page", "title": "First"}},
				"_links":  map[string]any{"next": "/wiki/rest/api/space/ENG/state/content?state-id=7&start=1&limit=2&expand=space&expand=version"},
			})
		case "1":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{"id": "p2", "type": "page", "title": "Second"}},
				"_links":  map[string]any{},
			})
		default:
			t.Fatalf("start = %q", q.Get("start"))
		}
	})

	got, err := ListContentWithState(context.Background(), c, ContentWithStateOptions{
		SpaceKey: "ENG",
		StateID:  int64Ptr(7),
		Expand:   []string{"space", "version"},
		Limit:    2,
	})
	if err != nil {
		t.Fatalf("ListContentWithState: %v", err)
	}
	if len(requests) != 2 {
		t.Fatalf("requests = %v, want two paginated requests", requests)
	}
	if len(got) != 2 || got[0].ID != "p1" || got[1].Title != "Second" {
		t.Fatalf("content = %+v", got)
	}
}

func TestContentStatesRejectServerFlavorAndMissingInputs(t *testing.T) {
	_, server := testClientWithFlavor(t, client.FlavorServer, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected server request: %s %s", r.Method, r.URL.RequestURI())
	})
	_, cloud := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected cloud request: %s %s", r.Method, r.URL.RequestURI())
	})

	cases := []struct {
		name string
		err  error
		want string
	}{
		{name: "server current", err: firstContentStateError(GetContentState(context.Background(), server, "123", "")), want: "Confluence Cloud"},
		{name: "server available", err: firstAvailableContentStateError(ListAvailableContentStates(context.Background(), server, "123")), want: "Confluence Cloud"},
		{name: "server custom", err: firstContentStatesError(ListCustomContentStates(context.Background(), server)), want: "Confluence Cloud"},
		{name: "server space", err: firstContentStatesError(ListSpaceContentStates(context.Background(), server, "ENG")), want: "Confluence Cloud"},
		{name: "server settings", err: firstContentStateSettingsError(GetContentStateSettings(context.Background(), server, "ENG")), want: "Confluence Cloud"},
		{name: "server content", err: firstContentError(ListContentWithState(context.Background(), server, ContentWithStateOptions{SpaceKey: "ENG", StateID: int64Ptr(7)})), want: "Confluence Cloud"},
		{name: "missing content id", err: firstContentStateError(GetContentState(context.Background(), cloud, "", "")), want: "content id"},
		{name: "missing space key", err: firstContentStatesError(ListSpaceContentStates(context.Background(), cloud, "")), want: "space key"},
		{name: "missing state id", err: firstContentError(ListContentWithState(context.Background(), cloud, ContentWithStateOptions{SpaceKey: "ENG"})), want: "state id"},
	}
	for _, tc := range cases {
		if tc.err == nil || !strings.Contains(tc.err.Error(), tc.want) {
			t.Fatalf("%s error = %v, want %q", tc.name, tc.err, tc.want)
		}
	}
}

func firstContentStateError(_ *ContentStateResponse, err error) error {
	return err
}

func firstAvailableContentStateError(_ *AvailableContentStates, err error) error {
	return err
}

func firstContentStatesError(_ []ContentState, err error) error {
	return err
}

func firstContentStateSettingsError(_ *ContentStateSettings, err error) error {
	return err
}

func firstContentError(_ []Content, err error) error {
	return err
}

func int64Ptr(v int64) *int64 {
	return &v
}
