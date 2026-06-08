package conf

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

func TestListCommentsServerUsesV1ChildCommentEndpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorServer, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/content/12345/child/comment" {
			t.Fatalf("path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("expand"); got != "body.view,version,extensions.inlineProperties,extensions.resolution" {
			t.Fatalf("expand: %q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"results": []map[string]any{{
				"id": "c1",
				"version": map[string]any{
					"when": "2026-06-08T10:00:00Z",
					"by":   map[string]any{"displayName": "Server User"},
				},
				"body": map[string]any{"view": map[string]any{"value": "<p>server body</p>"}},
				"extensions": map[string]any{
					"location": "footer",
				},
			}},
		})
	})

	got, err := ListComments(context.Background(), c, "12345", []string{"footer"}, 2)
	if err != nil {
		t.Fatalf("ListComments: %v", err)
	}
	if len(got) != 1 || got[0].Author != "Server User" || got[0].Body != "server body" {
		t.Fatalf("unexpected comments: %+v", got)
	}
}

func TestListCommentsCloudUsesV2FooterAndInlineEndpoints(t *testing.T) {
	requests := make([]string, 0, 2)
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.URL.RequestURI())
		switch r.URL.Path {
		case "/wiki/api/v2/pages/12345/footer-comments":
			if got := r.URL.Query().Get("body-format"); got != "STORAGE" {
				t.Fatalf("footer body-format: %q", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{
					"id":     "f1",
					"pageId": "12345",
					"version": map[string]any{
						"createdAt": "2026-06-08T10:00:00Z",
						"authorId":  "author-1",
					},
					"body":   map[string]any{"storage": map[string]any{"value": "<p>footer body</p>"}},
					"_links": map[string]any{"webui": "/pages/12345?focusedCommentId=f1"},
				}},
				"_links": map[string]any{},
			})
		case "/wiki/api/v2/pages/12345/inline-comments":
			if got := r.URL.Query().Get("resolution-status"); got != "open" {
				t.Fatalf("inline resolution-status: %q", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{
					"id":               "i1",
					"pageId":           "12345",
					"resolutionStatus": "open",
					"version": map[string]any{
						"createdAt": "2026-06-08T10:01:00Z",
						"authorId":  "author-2",
					},
					"body": map[string]any{"storage": map[string]any{"value": "<p>inline body</p>"}},
					"properties": map[string]any{
						"inlineOriginalSelection": "selected text",
					},
				}},
				"_links": map[string]any{},
			})
		default:
			t.Fatalf("path: %s", r.URL.Path)
		}
	})

	got, err := ListComments(context.Background(), c, "12345", []string{"footer", "inline"}, 2)
	if err != nil {
		t.Fatalf("ListComments: %v", err)
	}
	if len(requests) != 2 {
		t.Fatalf("requests = %d, want 2 (%v)", len(requests), requests)
	}
	if len(got) != 2 {
		t.Fatalf("len(comments) = %d, want 2 (%+v)", len(got), got)
	}
	if got[0].Location != "footer" || got[0].Author != "author-1" || got[0].Body != "footer body" {
		t.Fatalf("unexpected footer comment: %+v", got[0])
	}
	if got[1].Location != "inline" || got[1].Author != "author-2" || got[1].InlineOriginalSelection != "selected text" {
		t.Fatalf("unexpected inline comment: %+v", got[1])
	}
}

func TestListCommentsCloudResolvedOnlyUsesInlineResolvedFilter(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/wiki/api/v2/pages/12345/inline-comments" {
			t.Fatalf("path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("resolution-status"); got != "resolved" {
			t.Fatalf("resolution-status: %q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"results": []map[string]any{{
				"id":               "r1",
				"resolutionStatus": "resolved",
				"version":          map[string]any{"authorId": "author-3"},
				"body":             map[string]any{"storage": map[string]any{"value": "<p>resolved body</p>"}},
			}},
			"_links": map[string]any{},
		})
	})

	got, err := ListComments(context.Background(), c, "12345", []string{"resolved"}, 5)
	if err != nil {
		t.Fatalf("ListComments: %v", err)
	}
	if len(got) != 1 || !got[0].Resolved || got[0].Location != "resolved" {
		t.Fatalf("unexpected comments: %+v", got)
	}
}

func TestListCommentsCloudFollowsV2Pagination(t *testing.T) {
	requests := make([]string, 0, 2)
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.URL.RequestURI())
		switch len(requests) {
		case 1:
			if r.URL.RequestURI() != "/wiki/api/v2/pages/12345/footer-comments?body-format=STORAGE&limit=2" {
				t.Fatalf("first request: %s", r.URL.RequestURI())
			}
			w.Header().Set("Link", "<http://"+r.Host+"/wiki/api/v2/pages/12345/footer-comments?cursor=abc>; rel=\"next\"")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{
					"id":      "f1",
					"version": map[string]any{"authorId": "author-1"},
					"body":    map[string]any{"storage": map[string]any{"value": "<p>first</p>"}},
				}},
				"_links": map[string]any{},
			})
		case 2:
			if r.URL.RequestURI() != "/wiki/api/v2/pages/12345/footer-comments?cursor=abc" {
				t.Fatalf("second request: %s", r.URL.RequestURI())
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{
					"id":      "f2",
					"version": map[string]any{"authorId": "author-2"},
					"body":    map[string]any{"storage": map[string]any{"value": "<p>second</p>"}},
				}},
				"_links": map[string]any{},
			})
		default:
			t.Fatalf("unexpected extra request: %s", r.URL.RequestURI())
		}
	})

	got, err := ListComments(context.Background(), c, "12345", []string{"footer"}, 2)
	if err != nil {
		t.Fatalf("ListComments: %v", err)
	}
	if len(requests) != 2 {
		t.Fatalf("requests = %d, want 2 (%v)", len(requests), requests)
	}
	if len(got) != 2 || got[0].Body != "first" || got[1].Body != "second" {
		t.Fatalf("unexpected comments: %+v", got)
	}
}

func TestCreateFooterCommentCloudUsesV2Endpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/wiki/api/v2/footer-comments" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		body := readJSONMap(t, r.Body)
		if got := body["pageId"]; got != "12345" {
			t.Fatalf("pageId = %v, want 12345", got)
		}
		commentBody := mapValue(t, body, "body")
		if got := commentBody["representation"]; got != "storage" {
			t.Fatalf("body.representation = %v, want storage", got)
		}
		if got := commentBody["value"]; got != "<p>hello</p>" {
			t.Fatalf("body.value = %v, want <p>hello</p>", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":     "c1",
			"pageId": "12345",
			"version": map[string]any{
				"number":    1,
				"createdAt": "2026-06-08T10:00:00Z",
				"authorId":  "author-1",
			},
			"body": map[string]any{"storage": map[string]any{"value": "<p>hello</p>"}},
		})
	})

	got, err := CreateFooterComment(context.Background(), c, CommentInput{
		PageID:    "12345",
		BodyValue: "<p>hello</p>",
	})
	if err != nil {
		t.Fatalf("CreateFooterComment: %v", err)
	}
	if got.ID != "c1" || got.Body != "hello" || got.Location != "footer" || got.VersionNumber != 1 {
		t.Fatalf("unexpected comment: %+v", got)
	}
}

func TestCreateFooterCommentCloudSupportsBlogPostTarget(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/wiki/api/v2/footer-comments" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		body := readJSONMap(t, r.Body)
		if got := body["blogPostId"]; got != "2001" {
			t.Fatalf("blogPostId = %v, want 2001", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":         "c2",
			"blogPostId": "2001",
			"version":    map[string]any{"number": 1, "authorId": "author-2"},
			"body":       map[string]any{"storage": map[string]any{"value": "<p>post comment</p>"}},
		})
	})

	got, err := CreateFooterComment(context.Background(), c, CommentInput{
		BlogPostID: "2001",
		BodyValue:  "<p>post comment</p>",
	})
	if err != nil {
		t.Fatalf("CreateFooterComment: %v", err)
	}
	if got.ID != "c2" || got.Body != "post comment" {
		t.Fatalf("unexpected comment: %+v", got)
	}
}

func TestCreateFooterCommentCloudSupportsParentReply(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/wiki/api/v2/footer-comments" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		body := readJSONMap(t, r.Body)
		if got := body["parentCommentId"]; got != "parent-1" {
			t.Fatalf("parentCommentId = %v, want parent-1", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":              "reply-1",
			"parentCommentId": "parent-1",
			"version":         map[string]any{"number": 1, "authorId": "author-3"},
			"body":            map[string]any{"storage": map[string]any{"value": "<p>reply</p>"}},
		})
	})

	got, err := CreateFooterComment(context.Background(), c, CommentInput{
		ParentCommentID: "parent-1",
		BodyValue:       "<p>reply</p>",
	})
	if err != nil {
		t.Fatalf("CreateFooterComment: %v", err)
	}
	if got.ID != "reply-1" || got.Body != "reply" {
		t.Fatalf("unexpected reply: %+v", got)
	}
}

func TestCreateFooterCommentServerUsesContentResource(t *testing.T) {
	_, c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/rest/api/content" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		body := readJSONMap(t, r.Body)
		if got := body["type"]; got != "comment" {
			t.Fatalf("type = %v, want comment", got)
		}
		container := mapValue(t, body, "container")
		if got := container["id"]; got != "12345" {
			t.Fatalf("container.id = %v, want 12345", got)
		}
		if got := container["type"]; got != "page" {
			t.Fatalf("container.type = %v, want page", got)
		}
		commentBody := mapValue(t, body, "body")
		storage := mapValue(t, commentBody, "storage")
		if got := storage["representation"]; got != "storage" {
			t.Fatalf("body.storage.representation = %v, want storage", got)
		}
		if got := storage["value"]; got != "<p>server</p>" {
			t.Fatalf("body.storage.value = %v, want <p>server</p>", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":      "c1",
			"type":    "comment",
			"version": map[string]any{"number": 1, "when": "2026-06-08T10:00:00Z", "by": map[string]any{"displayName": "Server User"}},
			"body":    map[string]any{"storage": map[string]any{"value": "<p>server</p>"}},
		})
	})

	got, err := CreateFooterComment(context.Background(), c, CommentInput{
		PageID:    "12345",
		BodyValue: "<p>server</p>",
	})
	if err != nil {
		t.Fatalf("CreateFooterComment: %v", err)
	}
	if got.ID != "c1" || got.Author != "Server User" || got.Body != "server" || got.VersionNumber != 1 {
		t.Fatalf("unexpected comment: %+v", got)
	}
}

func TestGetFooterCommentCloudUsesV2Endpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/wiki/api/v2/footer-comments/c1" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		q := r.URL.Query()
		if got := q.Get("body-format"); got != "STORAGE" {
			t.Fatalf("body-format = %q, want STORAGE", got)
		}
		if got := q.Get("include-version"); got != "true" {
			t.Fatalf("include-version = %q, want true", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":      "c1",
			"version": map[string]any{"number": 4, "createdAt": "2026-06-08T10:00:00Z", "authorId": "author-1"},
			"body":    map[string]any{"storage": map[string]any{"value": "<p>current</p>"}},
		})
	})

	got, err := GetFooterComment(context.Background(), c, "c1")
	if err != nil {
		t.Fatalf("GetFooterComment: %v", err)
	}
	if got.ID != "c1" || got.Body != "current" || got.VersionNumber != 4 {
		t.Fatalf("unexpected comment: %+v", got)
	}
}

func TestGetFooterCommentServerUsesContentEndpoint(t *testing.T) {
	_, c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/rest/api/content/c1" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		if got := r.URL.Query().Get("expand"); got != "body.storage,version" {
			t.Fatalf("expand = %q, want body.storage,version", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":      "c1",
			"type":    "comment",
			"version": map[string]any{"number": 4, "when": "2026-06-08T10:00:00Z", "by": map[string]any{"displayName": "Server User"}},
			"body":    map[string]any{"storage": map[string]any{"value": "<p>current</p>"}},
		})
	})

	got, err := GetFooterComment(context.Background(), c, "c1")
	if err != nil {
		t.Fatalf("GetFooterComment: %v", err)
	}
	if got.ID != "c1" || got.Author != "Server User" || got.Body != "current" || got.VersionNumber != 4 {
		t.Fatalf("unexpected comment: %+v", got)
	}
}

func TestUpdateFooterCommentCloudUsesV2Endpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/wiki/api/v2/footer-comments/c1" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		body := readJSONMap(t, r.Body)
		version := mapValue(t, body, "version")
		if got := version["number"]; got != float64(5) {
			t.Fatalf("version.number = %v, want 5", got)
		}
		commentBody := mapValue(t, body, "body")
		if got := commentBody["representation"]; got != "storage" {
			t.Fatalf("body.representation = %v, want storage", got)
		}
		if got := commentBody["value"]; got != "<p>updated</p>" {
			t.Fatalf("body.value = %v, want <p>updated</p>", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":      "c1",
			"version": map[string]any{"number": 5, "createdAt": "2026-06-08T10:00:00Z", "authorId": "author-1"},
			"body":    map[string]any{"storage": map[string]any{"value": "<p>updated</p>"}},
		})
	})

	got, err := UpdateFooterComment(context.Background(), c, CommentInput{
		ID:            "c1",
		BodyValue:     "<p>updated</p>",
		VersionNumber: 4,
	})
	if err != nil {
		t.Fatalf("UpdateFooterComment: %v", err)
	}
	if got.ID != "c1" || got.Body != "updated" || got.VersionNumber != 5 {
		t.Fatalf("unexpected comment: %+v", got)
	}
}

func TestUpdateFooterCommentServerUsesContentResource(t *testing.T) {
	_, c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/rest/api/content/c1" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		body := readJSONMap(t, r.Body)
		if got := body["type"]; got != "comment" {
			t.Fatalf("type = %v, want comment", got)
		}
		version := mapValue(t, body, "version")
		if got := version["number"]; got != float64(5) {
			t.Fatalf("version.number = %v, want 5", got)
		}
		commentBody := mapValue(t, body, "body")
		storage := mapValue(t, commentBody, "storage")
		if got := storage["value"]; got != "<p>updated</p>" {
			t.Fatalf("body.storage.value = %v, want <p>updated</p>", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":      "c1",
			"type":    "comment",
			"version": map[string]any{"number": 5, "when": "2026-06-08T10:00:00Z", "by": map[string]any{"displayName": "Server User"}},
			"body":    map[string]any{"storage": map[string]any{"value": "<p>updated</p>"}},
		})
	})

	got, err := UpdateFooterComment(context.Background(), c, CommentInput{
		ID:            "c1",
		BodyValue:     "<p>updated</p>",
		VersionNumber: 4,
	})
	if err != nil {
		t.Fatalf("UpdateFooterComment: %v", err)
	}
	if got.ID != "c1" || got.Body != "updated" || got.VersionNumber != 5 {
		t.Fatalf("unexpected comment: %+v", got)
	}
}

func TestDeleteFooterCommentCloudUsesV2Endpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/wiki/api/v2/footer-comments/c1" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		w.WriteHeader(http.StatusNoContent)
	})

	if err := DeleteFooterComment(context.Background(), c, "c1"); err != nil {
		t.Fatalf("DeleteFooterComment: %v", err)
	}
}

func TestDeleteFooterCommentServerUsesContentEndpoint(t *testing.T) {
	_, c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/rest/api/content/c1" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		w.WriteHeader(http.StatusNoContent)
	})

	if err := DeleteFooterComment(context.Background(), c, "c1"); err != nil {
		t.Fatalf("DeleteFooterComment: %v", err)
	}
}
