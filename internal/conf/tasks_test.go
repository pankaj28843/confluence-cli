package conf

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

func TestListTasksCloudUsesV2EndpointAndFilters(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/wiki/api/v2/tasks" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		q := r.URL.Query()
		if got := q.Get("limit"); got != "2" {
			t.Fatalf("limit = %q, want 2", got)
		}
		if got := q.Get("status"); got != "incomplete" {
			t.Fatalf("status = %q, want incomplete", got)
		}
		if got := q.Get("page-id"); got != "12345" {
			t.Fatalf("page-id = %q, want 12345", got)
		}
		if got := q.Get("assigned-to"); got != "acct-1" {
			t.Fatalf("assigned-to = %q, want acct-1", got)
		}
		if got := q.Get("include-blank-tasks"); got != "true" {
			t.Fatalf("include-blank-tasks = %q, want true", got)
		}
		if got := q.Get("body-format"); got != "storage" {
			t.Fatalf("body-format = %q, want storage", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"results": []map[string]any{{
				"id":         "42",
				"localId":    "1",
				"spaceId":    "1001",
				"pageId":     "12345",
				"status":     "incomplete",
				"assignedTo": "acct-1",
				"body": map[string]any{
					"storage": map[string]any{
						"representation": "storage",
						"value":          "<p>Ship it</p>",
					},
				},
			}},
			"_links": map[string]any{},
		})
	})

	got, err := ListTasks(context.Background(), c, TaskFilter{
		Status:            "incomplete",
		PageID:            "12345",
		AssignedTo:        []string{"acct-1"},
		IncludeBlankTasks: true,
		BodyFormat:        "storage",
		Limit:             2,
	})
	if err != nil {
		t.Fatalf("ListTasks: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len(tasks) = %d, want 1", len(got))
	}
	if got[0].ID != "42" || got[0].Status != "incomplete" || got[0].PageID != "12345" {
		t.Fatalf("unexpected task: %+v", got[0])
	}
	if got[0].Body.Storage.Value != "<p>Ship it</p>" {
		t.Fatalf("storage body = %q, want <p>Ship it</p>", got[0].Body.Storage.Value)
	}
}

func TestListTasksCloudFollowsV2Pagination(t *testing.T) {
	requests := make([]string, 0, 2)
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.URL.RequestURI())
		w.Header().Set("Content-Type", "application/json")
		switch len(requests) {
		case 1:
			if r.URL.Path != "/wiki/api/v2/tasks" {
				t.Fatalf("first path: %s", r.URL.Path)
			}
			w.Header().Add("Link", `</wiki/api/v2/tasks?cursor=abc>; rel="next"`)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{"id": "1", "status": "incomplete"}},
				"_links":  map[string]any{},
			})
		case 2:
			if r.URL.RequestURI() != "/wiki/api/v2/tasks?cursor=abc" {
				t.Fatalf("second request: %s", r.URL.RequestURI())
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{"id": "2", "status": "complete"}},
				"_links":  map[string]any{},
			})
		default:
			t.Fatalf("unexpected extra request: %s", r.URL.RequestURI())
		}
	})

	got, err := ListTasks(context.Background(), c, TaskFilter{Limit: 2})
	if err != nil {
		t.Fatalf("ListTasks: %v", err)
	}
	if len(requests) != 2 {
		t.Fatalf("requests = %v, want two pages", requests)
	}
	if len(got) != 2 || got[0].ID != "1" || got[1].ID != "2" {
		t.Fatalf("unexpected tasks: %+v", got)
	}
}

func TestGetTaskCloudUsesV2Endpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/wiki/api/v2/tasks/42" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		if got := r.URL.Query().Get("body-format"); got != "storage" {
			t.Fatalf("body-format = %q, want storage", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":     "42",
			"status": "complete",
			"pageId": "12345",
			"body": map[string]any{
				"storage": map[string]any{
					"representation": "storage",
					"value":          "<p>Done</p>",
				},
			},
		})
	})

	got, err := GetTask(context.Background(), c, "42", "storage")
	if err != nil {
		t.Fatalf("GetTask: %v", err)
	}
	if got.ID != "42" || got.Status != "complete" || got.Body.Storage.Value != "<p>Done</p>" {
		t.Fatalf("unexpected task: %+v", got)
	}
}

func TestUpdateTaskStatusCloudUsesV2Endpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/wiki/api/v2/tasks/42" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		body := readJSONMap(t, r.Body)
		if got := body["id"]; got != "42" {
			t.Fatalf("id = %v, want 42", got)
		}
		if got := body["status"]; got != "complete" {
			t.Fatalf("status = %v, want complete", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":          "42",
			"status":      "complete",
			"completedBy": "acct-2",
		})
	})

	got, err := UpdateTaskStatus(context.Background(), c, "42", "complete")
	if err != nil {
		t.Fatalf("UpdateTaskStatus: %v", err)
	}
	if got.ID != "42" || got.Status != "complete" || got.CompletedBy != "acct-2" {
		t.Fatalf("unexpected task: %+v", got)
	}
}

func TestListLongTasksServerUsesV1Endpoint(t *testing.T) {
	_, c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/rest/api/longtask" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		q := r.URL.Query()
		if got := q.Get("start"); got != "0" {
			t.Fatalf("start = %q, want 0", got)
		}
		if got := q.Get("limit"); got != "2" {
			t.Fatalf("limit = %q, want 2", got)
		}
		if got := q.Get("expand"); got != "messages" {
			t.Fatalf("expand = %q, want messages", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"results": []map[string]any{{
				"id":                 "lt1",
				"elapsedTime":        123,
				"percentageComplete": 50,
				"successful":         true,
				"name": map[string]any{
					"translation": "Reindex",
				},
			}},
			"size":   1,
			"_links": map[string]any{},
		})
	})

	got, err := ListLongTasks(context.Background(), c, LongTaskFilter{Expand: "messages", Limit: 2})
	if err != nil {
		t.Fatalf("ListLongTasks: %v", err)
	}
	if len(got) != 1 || got[0].ID != "lt1" || got[0].PercentageComplete != 50 {
		t.Fatalf("unexpected long tasks: %+v", got)
	}
	if got[0].Name.Translation != "Reindex" {
		t.Fatalf("name.translation = %q, want Reindex", got[0].Name.Translation)
	}
}

func TestGetLongTaskServerUsesV1Endpoint(t *testing.T) {
	_, c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/rest/api/longtask/lt1" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		if got := r.URL.Query().Get("expand"); got != "messages" {
			t.Fatalf("expand = %q, want messages", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":                 "lt1",
			"elapsedTime":        456,
			"percentageComplete": 100,
			"successful":         true,
			"name": map[string]any{
				"translation": "Import",
			},
		})
	})

	got, err := GetLongTask(context.Background(), c, "lt1", "messages")
	if err != nil {
		t.Fatalf("GetLongTask: %v", err)
	}
	if got.ID != "lt1" || got.Name.Translation != "Import" || got.PercentageComplete != 100 {
		t.Fatalf("unexpected long task: %+v", got)
	}
}
