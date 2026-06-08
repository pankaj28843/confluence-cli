package conf

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

func TestListBlogPostsCloudUsesV2SpaceBlogPostsEndpoint(t *testing.T) {
	var requests []string
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.Method+" "+r.URL.RequestURI())
		w.Header().Set("Content-Type", "application/json")

		switch len(requests) {
		case 1:
			if r.Method != http.MethodGet || r.URL.Path != "/wiki/api/v2/spaces" {
				t.Fatalf("space lookup request: %s %s", r.Method, r.URL.RequestURI())
			}
			if got := r.URL.Query().Get("keys"); got != "ENG" {
				t.Fatalf("space lookup keys = %q, want ENG", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{
					"id":   "1001",
					"key":  "ENG",
					"name": "Engineering",
				}},
			})
		case 2:
			if r.Method != http.MethodGet || r.URL.Path != "/wiki/api/v2/spaces/1001/blogposts" {
				t.Fatalf("list request: %s %s", r.Method, r.URL.RequestURI())
			}
			q := r.URL.Query()
			if got := q.Get("limit"); got != "25" {
				t.Fatalf("limit = %q, want 25", got)
			}
			if got := q.Get("title"); got != "Weekly" {
				t.Fatalf("title = %q, want Weekly", got)
			}
			if got := q.Get("status"); got != "current" {
				t.Fatalf("status = %q, want current", got)
			}
			if got := q.Get("body-format"); got != "storage" {
				t.Fatalf("body-format = %q, want storage", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{
					"id":      "2001",
					"title":   "Weekly",
					"spaceId": "1001",
					"status":  "current",
					"version": map[string]any{"number": 3},
				}},
			})
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
	})

	got, err := ListBlogPosts(context.Background(), c, BlogPostListOptions{
		SpaceKey: "ENG",
		Title:    "Weekly",
		Status:   "current",
		Limit:    25,
	})
	if err != nil {
		t.Fatalf("ListBlogPosts: %v", err)
	}
	if len(got) != 1 || got[0].ID != "2001" || got[0].Type != "blogpost" || got[0].Space.Key != "ENG" {
		t.Fatalf("unexpected blogposts: %+v", got)
	}
	if len(requests) != 2 {
		t.Fatalf("requests = %v, want space lookup then blogpost list", requests)
	}
}

func TestListBlogPostsServerUsesContentTypeBlogpost(t *testing.T) {
	_, c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/rest/api/content" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		q := r.URL.Query()
		if got := q.Get("type"); got != "blogpost" {
			t.Fatalf("type = %q, want blogpost", got)
		}
		if got := q.Get("spaceKey"); got != "ENG" {
			t.Fatalf("spaceKey = %q, want ENG", got)
		}
		if got := q.Get("postingDay"); got != "2026-06-08" {
			t.Fatalf("postingDay = %q, want 2026-06-08", got)
		}
		if got := q.Get("title"); got != "Weekly" {
			t.Fatalf("title = %q, want Weekly", got)
		}
		if got := q.Get("limit"); got != "10" {
			t.Fatalf("limit = %q, want 10", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"results": []map[string]any{{
				"id":     "3001",
				"type":   "blogpost",
				"title":  "Weekly",
				"status": "current",
				"space":  map[string]any{"key": "ENG"},
			}},
		})
	})

	got, err := ListBlogPosts(context.Background(), c, BlogPostListOptions{
		SpaceKey:   "ENG",
		PostingDay: "2026-06-08",
		Title:      "Weekly",
		Limit:      10,
	})
	if err != nil {
		t.Fatalf("ListBlogPosts: %v", err)
	}
	if len(got) != 1 || got[0].ID != "3001" || got[0].Type != "blogpost" {
		t.Fatalf("unexpected blogposts: %+v", got)
	}
}

func TestGetBlogPostCloudUsesV2Endpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/wiki/api/v2/blogposts/2001" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		q := r.URL.Query()
		if got := q.Get("body-format"); got != "storage" {
			t.Fatalf("body-format = %q, want storage", got)
		}
		if got := q.Get("include-labels"); got != "true" {
			t.Fatalf("include-labels = %q, want true", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":      "2001",
			"title":   "Weekly",
			"spaceId": "1001",
			"status":  "current",
			"body":    map[string]any{"storage": map[string]any{"value": "<p>x</p>", "representation": "storage"}},
			"version": map[string]any{"number": 4},
		})
	})

	got, err := GetBlogPost(context.Background(), c, "2001")
	if err != nil {
		t.Fatalf("GetBlogPost: %v", err)
	}
	if got.ID != "2001" || got.Type != "blogpost" || got.Body.Storage.Value != "<p>x</p>" {
		t.Fatalf("unexpected blogpost: %+v", got)
	}
}

func TestCreateBlogPostCloudUsesV2EndpointAndSpaceID(t *testing.T) {
	var requests []string
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.Method+" "+r.URL.RequestURI())
		w.Header().Set("Content-Type", "application/json")

		switch len(requests) {
		case 1:
			if r.Method != http.MethodGet || r.URL.Path != "/wiki/api/v2/spaces" {
				t.Fatalf("space lookup request: %s %s", r.Method, r.URL.RequestURI())
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{"id": "1001", "key": "ENG"}},
			})
		case 2:
			if r.Method != http.MethodPost || r.URL.Path != "/wiki/api/v2/blogposts" {
				t.Fatalf("create request: %s %s", r.Method, r.URL.RequestURI())
			}
			body := readJSONMap(t, r.Body)
			if got := body["spaceId"]; got != "1001" {
				t.Fatalf("spaceId = %v, want 1001", got)
			}
			if got := body["status"]; got != "draft" {
				t.Fatalf("status = %v, want draft", got)
			}
			if got := body["title"]; got != "Draft Post" {
				t.Fatalf("title = %v, want Draft Post", got)
			}
			postBody := mapValue(t, body, "body")
			if got := postBody["representation"]; got != "storage" {
				t.Fatalf("body.representation = %v, want storage", got)
			}
			if got := postBody["value"]; got != "<p>draft</p>" {
				t.Fatalf("body.value = %v, want <p>draft</p>", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":      "2002",
				"title":   "Draft Post",
				"status":  "draft",
				"spaceId": "1001",
				"version": map[string]any{"number": 1},
			})
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
	})

	got, err := CreateBlogPost(context.Background(), c, BlogPostInput{
		SpaceKey:  "ENG",
		Title:     "Draft Post",
		BodyValue: "<p>draft</p>",
		Status:    "draft",
	})
	if err != nil {
		t.Fatalf("CreateBlogPost: %v", err)
	}
	if got.ID != "2002" || got.Type != "blogpost" || got.Space.Key != "ENG" {
		t.Fatalf("unexpected blogpost: %+v", got)
	}
}

func TestCreateBlogPostServerRejectsDraftStatus(t *testing.T) {
	_, c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
	})

	_, err := CreateBlogPost(context.Background(), c, BlogPostInput{
		SpaceKey:  "ENG",
		Title:     "Draft Post",
		BodyValue: "<p>draft</p>",
		Status:    "draft",
	})
	if err == nil {
		t.Fatalf("CreateBlogPost: got nil error, want unsupported draft error")
	}
}

func TestUpdateBlogPostCloudUsesV2Endpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/wiki/api/v2/blogposts/2001" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		body := readJSONMap(t, r.Body)
		if got := body["id"]; got != "2001" {
			t.Fatalf("id = %v, want 2001", got)
		}
		if got := body["status"]; got != "current" {
			t.Fatalf("status = %v, want current", got)
		}
		if got := body["title"]; got != "Updated" {
			t.Fatalf("title = %v, want Updated", got)
		}
		version := mapValue(t, body, "version")
		if got := version["number"]; got != float64(6) {
			t.Fatalf("version.number = %v, want 6", got)
		}
		postBody := mapValue(t, body, "body")
		if got := postBody["representation"]; got != "storage" {
			t.Fatalf("body.representation = %v, want storage", got)
		}
		if got := postBody["value"]; got != "<p>updated</p>" {
			t.Fatalf("body.value = %v, want <p>updated</p>", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":      "2001",
			"title":   "Updated",
			"status":  "current",
			"version": map[string]any{"number": 6},
		})
	})

	got, err := UpdateBlogPost(context.Background(), c, BlogPostInput{
		ID:            "2001",
		Title:         "Updated",
		BodyValue:     "<p>updated</p>",
		VersionNumber: 5,
	})
	if err != nil {
		t.Fatalf("UpdateBlogPost: %v", err)
	}
	if got.ID != "2001" || got.Title != "Updated" || got.Version.Number != 6 || got.Type != "blogpost" {
		t.Fatalf("unexpected blogpost: %+v", got)
	}
}

func TestDeleteBlogPostCloudSupportsPurgeAndDraftFlags(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/wiki/api/v2/blogposts/2001" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		if got := r.URL.Query().Get("purge"); got != "true" {
			t.Fatalf("purge = %q, want true", got)
		}
		if got := r.URL.Query().Get("draft"); got != "true" {
			t.Fatalf("draft = %q, want true", got)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	if err := DeleteBlogPost(context.Background(), c, "2001", BlogPostDeleteOptions{Purge: true, Draft: true}); err != nil {
		t.Fatalf("DeleteBlogPost: %v", err)
	}
}

func TestDeleteBlogPostServerUsesContentDeleteEndpoint(t *testing.T) {
	_, c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/rest/api/content/3001" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		if got := r.URL.Query().Get("status"); got != "trashed" {
			t.Fatalf("status = %q, want trashed", got)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	if err := DeleteBlogPost(context.Background(), c, "3001", BlogPostDeleteOptions{Purge: true}); err != nil {
		t.Fatalf("DeleteBlogPost: %v", err)
	}
}
