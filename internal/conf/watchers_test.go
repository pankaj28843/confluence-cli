package conf

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

func TestListContentWatchersServerUsesWatchersEndpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorServer, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/rest/api/content/12345/watchers" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		if got := r.URL.Query().Get("start"); got != "0" {
			t.Fatalf("start = %q", got)
		}
		if got := r.URL.Query().Get("limit"); got != "5" {
			t.Fatalf("limit = %q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"results": []map[string]any{{"username": "alice", "displayName": "Alice Example"}},
			"start":   0,
			"limit":   5,
			"size":    1,
		})
	})

	got, err := ListContentWatchers(context.Background(), c, "12345", 5)
	if err != nil {
		t.Fatalf("ListContentWatchers: %v", err)
	}
	if len(got.Results) != 1 || got.Results[0].Watcher.DisplayName != "Alice Example" {
		t.Fatalf("watchers = %+v", got.Results)
	}
}

func TestListContentWatchersCloudUsesNotificationEndpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/wiki/rest/api/content/12345/notification/child-created" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		if got := r.URL.Query().Get("limit"); got != "3" {
			t.Fatalf("limit = %q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"results": []map[string]any{{
				"type":      "page",
				"contentId": 12345,
				"watcher":   map[string]any{"accountId": "abc", "displayName": "Cloud User"},
			}},
			"start": 0,
			"limit": 3,
			"size":  1,
		})
	})

	got, err := ListContentWatchers(context.Background(), c, "12345", 3)
	if err != nil {
		t.Fatalf("ListContentWatchers: %v", err)
	}
	if len(got.Results) != 1 || got.Results[0].ContentID != 12345 || got.Results[0].Watcher.AccountID != "abc" {
		t.Fatalf("watchers = %+v", got.Results)
	}
}

func TestListSpaceWatchersUsesFlavorSpecificEndpoints(t *testing.T) {
	tests := []struct {
		name   string
		flavor client.Flavor
		path   string
		body   map[string]any
	}{
		{
			name:   "server",
			flavor: client.FlavorServer,
			path:   "/rest/api/space/ENG/watchers",
			body: map[string]any{
				"results": []map[string]any{{"username": "alice", "displayName": "Alice Example"}},
				"size":    1,
			},
		},
		{
			name:   "cloud",
			flavor: client.FlavorCloud,
			path:   "/wiki/rest/api/space/ENG/watch",
			body: map[string]any{
				"results": []map[string]any{{
					"type":     "space",
					"spaceKey": "ENG",
					"watcher":  map[string]any{"accountId": "abc", "displayName": "Cloud User"},
				}},
				"size": 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, c := testClientWithFlavor(t, tt.flavor, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet || r.URL.Path != tt.path {
					t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
				}
				if got := r.URL.Query().Get("limit"); got != "7" {
					t.Fatalf("limit = %q", got)
				}
				_ = json.NewEncoder(w).Encode(tt.body)
			})

			got, err := ListSpaceWatchers(context.Background(), c, "ENG", 7)
			if err != nil {
				t.Fatalf("ListSpaceWatchers: %v", err)
			}
			if len(got.Results) != 1 {
				t.Fatalf("watchers = %+v", got.Results)
			}
		})
	}
}

func TestGetWatchStatusParsesCloudObjectAndServerBoolean(t *testing.T) {
	_, cloud := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/wiki/rest/api/user/watch/content/12345" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		if got := r.URL.Query().Get("accountId"); got != "abc" {
			t.Fatalf("accountId = %q", got)
		}
		_, _ = w.Write([]byte(`{"watching":true}`))
	})
	contentStatus, err := GetContentWatchStatus(context.Background(), cloud, "12345", WatchStatusOptions{AccountID: "abc"})
	if err != nil {
		t.Fatalf("GetContentWatchStatus: %v", err)
	}
	if !contentStatus.Watching {
		t.Fatalf("content status = %+v", contentStatus)
	}

	_, server := testClientWithFlavor(t, client.FlavorServer, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/rest/api/user/watch/space/ENG" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		if got := r.URL.Query().Get("username"); got != "alice" {
			t.Fatalf("username = %q", got)
		}
		if got := r.URL.Query().Get("contentType"); got != "blogpost" {
			t.Fatalf("contentType = %q", got)
		}
		_, _ = w.Write([]byte(`true`))
	})
	spaceStatus, err := GetSpaceWatchStatus(context.Background(), server, "ENG", WatchStatusOptions{Username: "alice", ContentType: "blogpost"})
	if err != nil {
		t.Fatalf("GetSpaceWatchStatus: %v", err)
	}
	if !spaceStatus.Watching {
		t.Fatalf("space status = %+v", spaceStatus)
	}
}

func TestWatcherHelpersValidateRequiredInputs(t *testing.T) {
	if _, err := ListContentWatchers(context.Background(), nil, "", 25); err == nil {
		t.Fatal("ListContentWatchers should require content id")
	}
	if _, err := ListSpaceWatchers(context.Background(), nil, "", 25); err == nil {
		t.Fatal("ListSpaceWatchers should require space key")
	}
	if _, err := GetContentWatchStatus(context.Background(), nil, "", WatchStatusOptions{}); err == nil {
		t.Fatal("GetContentWatchStatus should require content id")
	}
	if _, err := GetSpaceWatchStatus(context.Background(), nil, "", WatchStatusOptions{}); err == nil {
		t.Fatal("GetSpaceWatchStatus should require space key")
	}
}
