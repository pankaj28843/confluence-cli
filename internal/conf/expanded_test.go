package conf

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

// One servemux exercises many read + write endpoints to pull coverage up
// above the 60% target without proliferating per-file tests.
func TestExpandedReadWrite(t *testing.T) {
	var lastMethod, lastPath string
	var lastBody []byte
	_, c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		lastMethod = r.Method
		lastPath = r.URL.Path
		lastBody, _ = io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")

		switch {
		case strings.HasSuffix(r.URL.Path, "/rest/api/user/current"):
			_ = json.NewEncoder(w).Encode(CurrentUser{Username: "alice", DisplayName: "Alice Example"})
		case strings.HasPrefix(r.URL.Path, "/rest/api/user"):
			_ = json.NewEncoder(w).Encode(User{Username: "bob", DisplayName: "Bob"})
		case strings.HasPrefix(r.URL.Path, "/rest/api/group/"):
			_ = json.NewEncoder(w).Encode(map[string]any{"results": []User{{Username: "alice"}}})
		case strings.HasPrefix(r.URL.Path, "/rest/api/group"):
			_ = json.NewEncoder(w).Encode(map[string]any{"results": []Group{{Name: "engineering", Type: "group"}}})
		case strings.HasPrefix(r.URL.Path, "/rest/api/search"):
			_ = json.NewEncoder(w).Encode(map[string]any{"results": []SearchHit{{Title: "hit", EntityType: "content"}}})
		case strings.HasSuffix(r.URL.Path, "/child/attachment"):
			if r.Method == http.MethodPost {
				_ = json.NewEncoder(w).Encode(map[string]any{"results": []Attachment{{ID: "att1", Title: "report.pdf"}}})
			} else {
				_ = json.NewEncoder(w).Encode(map[string]any{"results": []Attachment{{ID: "att1", Title: "logo.png"}}})
			}
		case strings.HasSuffix(r.URL.Path, "/child/comment"):
			_ = json.NewEncoder(w).Encode(map[string]any{"results": []any{
				map[string]any{"id": "c1", "body": map[string]any{"view": map[string]string{"value": "<p>hi</p>"}},
					"extensions": map[string]any{"location": "footer"}, "version": map[string]any{"when": "2026-01-01"}},
			}})
		case strings.HasSuffix(r.URL.Path, "/label"):
			if r.Method == http.MethodDelete {
				_, _ = w.Write([]byte(`{}`))
			} else if r.Method == http.MethodPost {
				_ = json.NewEncoder(w).Encode(map[string]any{"results": []Label{{Name: "shipped"}}})
			} else {
				_ = json.NewEncoder(w).Encode(map[string]any{"results": []Label{{Name: "existing"}}})
			}
		case strings.HasSuffix(r.URL.Path, "/history"), strings.HasSuffix(r.URL.Path, "/version"):
			_, _ = w.Write([]byte(`{"history":true}`))
		case strings.Contains(r.URL.Path, "/notification/child-created"):
			_, _ = w.Write([]byte(`{"watchers":[]}`))
		case strings.HasSuffix(r.URL.Path, "/watch"):
			_, _ = w.Write([]byte(`{"watchers":[]}`))
		case strings.HasSuffix(r.URL.Path, "/restriction"):
			_, _ = w.Write([]byte(`{"restrictions":[]}`))
		case strings.HasPrefix(r.URL.Path, "/rest/api/content/12345") && r.Method == http.MethodPut:
			cnt := Content{ID: "12345"}
			cnt.Version.Number = 5
			_ = json.NewEncoder(w).Encode(cnt)
		case strings.HasPrefix(r.URL.Path, "/rest/api/content/"):
			cnt := Content{ID: "12345", Title: "Demo", Type: "page"}
			cnt.Version.Number = 4
			cnt.Body.Storage.Value = "<h1>H</h1><p>p</p>"
			_ = json.NewEncoder(w).Encode(cnt)
		default:
			t.Fatalf("unhandled: %s %s", r.Method, r.URL.Path)
		}
	})

	ctx := context.Background()

	// User + group endpoints
	if u, err := GetCurrentUser(ctx, c); err != nil || u.Label() != "Alice Example" {
		t.Fatalf("GetCurrentUser: %v %+v", err, u)
	}
	if u, err := GetUser(ctx, c, "username", "bob"); err != nil || u.DisplayName != "Bob" {
		t.Fatalf("GetUser: %v %+v", err, u)
	}
	if gs, err := ListGroups(ctx, c, 5); err != nil || len(gs) != 1 {
		t.Fatalf("ListGroups: %v %+v", err, gs)
	}
	if ms, err := ListGroupMembers(ctx, c, "engineering", 5); err != nil || len(ms) != 1 {
		t.Fatalf("ListGroupMembers: %v %+v", err, ms)
	}
	if hs, err := SearchUsers(ctx, c, "alice", 5); err != nil || len(hs) != 1 {
		t.Fatalf("SearchUsers: %v %+v", err, hs)
	}

	// Content surface
	if p, err := GetContent(ctx, c, "12345", ""); err != nil || p.Version.Number != 4 {
		t.Fatalf("GetContent: %v %+v", err, p)
	}
	if _, err := GetHistory(ctx, c, "12345"); err != nil {
		t.Fatalf("GetHistory: %v", err)
	}
	if _, err := ListVersions(ctx, c, "12345"); err != nil {
		t.Fatalf("ListVersions: %v", err)
	}
	if hs, err := SearchGeneric(ctx, c, "type=space", 5); err != nil || len(hs) != 1 {
		t.Fatalf("SearchGeneric: %v %+v", err, hs)
	}

	// Attachments + labels + comments
	if as, err := ListAttachments(ctx, c, "12345", 10); err != nil || len(as) != 1 {
		t.Fatalf("ListAttachments: %v %+v", err, as)
	}
	if ls, err := ListLabels(ctx, c, "12345"); err != nil || len(ls) != 1 {
		t.Fatalf("ListLabels: %v %+v", err, ls)
	}
	if ls, err := AddLabels(ctx, c, "12345", []string{"shipped"}); err != nil || ls[0].Name != "shipped" {
		t.Fatalf("AddLabels: %v %+v", err, ls)
	}
	if err := RemoveLabel(ctx, c, "12345", "shipped"); err != nil {
		t.Fatalf("RemoveLabel: %v", err)
	}
	if cs, err := ListComments(ctx, c, "12345", nil, 5); err != nil || len(cs) != 1 {
		t.Fatalf("ListComments: %v %+v", err, cs)
	}

	// Watchers + restrictions
	if _, err := GetContentWatchers(ctx, c, "12345"); err != nil {
		t.Fatalf("GetContentWatchers: %v", err)
	}
	if _, err := GetSpaceWatchers(ctx, c, "ENG"); err != nil {
		t.Fatalf("GetSpaceWatchers: %v", err)
	}
	if _, err := GetContentRestrictions(ctx, c, "12345"); err != nil {
		t.Fatalf("GetContentRestrictions: %v", err)
	}

	// Writes: page update + attachment upload
	out, err := UpdatePage(ctx, c, UpdatePageInput{ID: "12345", Title: "New", BodyValue: "<p>x</p>", VersionNumber: 4})
	if err != nil || out.Version.Number != 5 {
		t.Fatalf("UpdatePage: %v %+v", err, out)
	}
	if lastMethod != http.MethodPut {
		t.Fatalf("UpdatePage should PUT, got %s %s", lastMethod, lastPath)
	}

	atts, err := UploadAttachment(ctx, c, "12345", "report.pdf", strings.NewReader("PDF-BYTES"), "v1")
	if err != nil || atts[0].Title != "report.pdf" {
		t.Fatalf("UploadAttachment: %v %+v", err, atts)
	}
	if !strings.Contains(string(lastBody), "PDF-BYTES") {
		t.Fatalf("UploadAttachment did not include the file bytes")
	}
}
