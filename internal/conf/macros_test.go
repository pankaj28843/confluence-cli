package conf

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

func TestGetMacroBodyCloudUsesV1MacroIDEndpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/wiki/rest/api/content/12345/history/7/macro/id/macro-1" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"name": "panel",
			"body": "<p>hello</p>",
			"parameters": map[string]any{
				"title": "Heads up",
			},
			"_links": map[string]any{
				"base": "https://example.atlassian.net/wiki",
			},
		})
	})

	got, err := GetMacroBody(context.Background(), c, MacroLookup{
		ContentID: "12345",
		Version:   7,
		MacroID:   "macro-1",
	})
	if err != nil {
		t.Fatalf("GetMacroBody: %v", err)
	}
	if got.Name != "panel" || got.Body != "<p>hello</p>" || got.Parameters["title"] != "Heads up" {
		t.Fatalf("unexpected macro instance: %+v", got)
	}
}

func TestGetMacroBodyServerUsesV1MacroIDEndpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorServer, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/rest/api/content/12345/history/7/macro/id/macro-1" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"name":       "toc",
			"body":       "",
			"parameters": map[string]any{"maxLevel": "3"},
		})
	})

	got, err := GetMacroBody(context.Background(), c, MacroLookup{
		ContentID: "12345",
		Version:   7,
		MacroID:   "macro-1",
	})
	if err != nil {
		t.Fatalf("GetMacroBody: %v", err)
	}
	if got.Name != "toc" || got.Parameters["maxLevel"] != "3" {
		t.Fatalf("unexpected macro instance: %+v", got)
	}
}

func TestGetMacroBodyServerSupportsDeprecatedHashEndpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorServer, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/rest/api/content/12345/history/7/macro/hash/hash-1" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"name": "info",
			"body": "<p>legacy hash lookup</p>",
		})
	})

	got, err := GetMacroBody(context.Background(), c, MacroLookup{
		ContentID: "12345",
		Version:   7,
		Hash:      "hash-1",
	})
	if err != nil {
		t.Fatalf("GetMacroBody: %v", err)
	}
	if got.Name != "info" || got.Body != "<p>legacy hash lookup</p>" {
		t.Fatalf("unexpected macro instance: %+v", got)
	}
}

func TestGetMacroBodyRejectsInvalidLookup(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
	})

	cases := []struct {
		name string
		in   MacroLookup
		want string
	}{
		{
			name: "missing content",
			in:   MacroLookup{Version: 1, MacroID: "macro-1"},
			want: "content id is required",
		},
		{
			name: "missing version",
			in:   MacroLookup{ContentID: "12345", MacroID: "macro-1"},
			want: "version must be greater than zero",
		},
		{
			name: "missing selector",
			in:   MacroLookup{ContentID: "12345", Version: 1},
			want: "exactly one",
		},
		{
			name: "both selectors",
			in:   MacroLookup{ContentID: "12345", Version: 1, MacroID: "macro-1", Hash: "hash-1"},
			want: "exactly one",
		},
		{
			name: "cloud hash",
			in:   MacroLookup{ContentID: "12345", Version: 1, Hash: "hash-1"},
			want: "only documented on Confluence Server/Data Center",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := GetMacroBody(context.Background(), c, tc.in)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("GetMacroBody error = %v, want contains %q", err, tc.want)
			}
		})
	}
}
