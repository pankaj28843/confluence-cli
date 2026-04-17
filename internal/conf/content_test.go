package conf

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestGetContent(t *testing.T) {
	_, c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/rest/api/content/12345") {
			t.Fatalf("path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("expand"); got != DefaultExpand {
			t.Fatalf("expand: %s", got)
		}
		cnt := Content{ID: "12345", Title: "Demo", Type: "page"}
		cnt.Body.Storage.Value = "<h1>Hello</h1><p>world</p>"
		_ = json.NewEncoder(w).Encode(cnt)
	})
	got, err := GetContent(context.Background(), c, "12345", "")
	if err != nil || got.ID != "12345" {
		t.Fatalf("GetContent: %v %+v", err, got)
	}
	if md := got.RenderMarkdown(false); !strings.Contains(md, "# Hello") {
		t.Fatalf("markdown conversion: %q", md)
	}
}

func TestSearchCQL(t *testing.T) {
	_, c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/content/search" {
			t.Fatalf("path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("cql"); got != "type=page" {
			t.Fatalf("cql: %s", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"results": []Content{{ID: "1", Title: "Hit"}}})
	})
	got, err := SearchCQL(context.Background(), c, "type=page", 3, "")
	if err != nil || len(got) != 1 {
		t.Fatalf("SearchCQL: %v %+v", err, got)
	}
}

func TestHTMLToMarkdown(t *testing.T) {
	// Block-level elements are preserved; inline emphasis/links inside a <p>
	// are flattened to plain text — matches the ported converter semantics.
	in := `<h1>Title</h1><p>Body text.</p><ul><li>one</li><li>two</li></ul><pre><code>fn()</code></pre>`
	md := HTMLToMarkdown(in)
	for _, want := range []string{"# Title", "Body text.", "- one", "- two", "```", "fn()"} {
		if !strings.Contains(md, want) {
			t.Errorf("missing %q in %q", want, md)
		}
	}
}

func TestGetChildren(t *testing.T) {
	_, c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/rest/api/content/12345/child/page") {
			t.Fatalf("path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"results": []Content{{ID: "c1", Title: "child"}}})
	})
	got, err := GetChildren(context.Background(), c, "12345", "", 10)
	if err != nil || len(got) != 1 {
		t.Fatalf("GetChildren: %v %+v", err, got)
	}
}
