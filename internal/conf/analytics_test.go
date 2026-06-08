package conf

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

func TestGetContentAnalyticsViewsCloudUsesV1Endpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/wiki/rest/api/analytics/content/123/views" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		if got := r.URL.Query().Get("fromDate"); got != "2026-01-02T00:00:00.000Z" {
			t.Fatalf("fromDate = %q, want 2026-01-02T00:00:00.000Z", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"id": 123, "count": 42})
	})

	got, err := GetContentViewCount(context.Background(), c, "123", AnalyticsOptions{
		FromDate: "2026-01-02T00:00:00.000Z",
	})
	if err != nil {
		t.Fatalf("GetContentViewCount: %v", err)
	}
	if got.ID != 123 || got.Count != 42 {
		t.Fatalf("view count = %+v", got)
	}
}

func TestGetContentAnalyticsViewersCloudUsesV1Endpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/wiki/rest/api/analytics/content/123/viewers" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		if got := r.URL.Query().Get("fromDate"); got != "" {
			t.Fatalf("fromDate = %q, want empty", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"id": 123, "count": 9})
	})

	got, err := GetContentViewerCount(context.Background(), c, "123", AnalyticsOptions{})
	if err != nil {
		t.Fatalf("GetContentViewerCount: %v", err)
	}
	if got.ID != 123 || got.Count != 9 {
		t.Fatalf("viewer count = %+v", got)
	}
}

func TestAnalyticsRejectServerFlavorAndMissingInputs(t *testing.T) {
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
		{name: "server views", err: firstAnalyticsCountError(GetContentViewCount(context.Background(), server, "123", AnalyticsOptions{})), want: "Confluence Cloud"},
		{name: "server viewers", err: firstAnalyticsCountError(GetContentViewerCount(context.Background(), server, "123", AnalyticsOptions{})), want: "Confluence Cloud"},
		{name: "missing views content id", err: firstAnalyticsCountError(GetContentViewCount(context.Background(), cloud, "", AnalyticsOptions{})), want: "content id"},
		{name: "missing viewers content id", err: firstAnalyticsCountError(GetContentViewerCount(context.Background(), cloud, "", AnalyticsOptions{})), want: "content id"},
	}
	for _, tc := range cases {
		if tc.err == nil || !strings.Contains(tc.err.Error(), tc.want) {
			t.Fatalf("%s error = %v, want %q", tc.name, tc.err, tc.want)
		}
	}
}

func firstAnalyticsCountError(_ *AnalyticsCount, err error) error {
	return err
}
