package conf

import (
	"context"
	"net/http"
	"testing"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

func TestDeletePageCloudUsesV2PageEndpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/wiki/api/v2/pages/12345" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		if got := r.URL.Query().Get("purge"); got != "" {
			t.Fatalf("purge = %q, want empty", got)
		}
		if got := r.URL.Query().Get("draft"); got != "" {
			t.Fatalf("draft = %q, want empty", got)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	if err := DeletePage(context.Background(), c, "12345", PageDeleteOptions{}); err != nil {
		t.Fatalf("DeletePage: %v", err)
	}
}

func TestDeletePageCloudSupportsPurgeAndDraftFlags(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/wiki/api/v2/pages/12345" {
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

	err := DeletePage(context.Background(), c, "12345", PageDeleteOptions{Purge: true, Draft: true})
	if err != nil {
		t.Fatalf("DeletePage: %v", err)
	}
}

func TestDeletePageServerUsesContentDeleteEndpoint(t *testing.T) {
	_, c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/rest/api/content/12345" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		if got := r.URL.Query().Get("status"); got != "" {
			t.Fatalf("status = %q, want empty", got)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	if err := DeletePage(context.Background(), c, "12345", PageDeleteOptions{}); err != nil {
		t.Fatalf("DeletePage: %v", err)
	}
}

func TestDeletePageServerPurgeUsesTrashedStatusParameter(t *testing.T) {
	_, c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/rest/api/content/12345" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		if got := r.URL.Query().Get("status"); got != "trashed" {
			t.Fatalf("status = %q, want trashed", got)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := DeletePage(context.Background(), c, "12345", PageDeleteOptions{Purge: true})
	if err != nil {
		t.Fatalf("DeletePage: %v", err)
	}
}

func TestDeletePageServerRejectsDraftFlag(t *testing.T) {
	_, c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
	})

	err := DeletePage(context.Background(), c, "12345", PageDeleteOptions{Draft: true})
	if err == nil {
		t.Fatalf("DeletePage: got nil error, want unsupported draft error")
	}
}
