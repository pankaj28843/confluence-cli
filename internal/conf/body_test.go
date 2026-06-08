package conf

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

func TestConvertBodyServerUsesContentBodyConvertEndpoint(t *testing.T) {
	_, c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/rest/api/contentbody/convert/view" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		if got := r.URL.Query().Get("expand"); got != "webresource.uris.css" {
			t.Fatalf("expand = %q, want webresource.uris.css", got)
		}
		body := readJSONMap(t, r.Body)
		if got := body["representation"]; got != "storage" {
			t.Fatalf("representation = %v, want storage", got)
		}
		if got := body["value"]; got != "<p>Hello</p>" {
			t.Fatalf("value = %v, want <p>Hello</p>", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"representation": "view",
			"value":          "<p>Hello</p>",
			"webresource": map[string]any{
				"keys": []string{"confluence.web.resources"},
			},
		})
	})

	got, err := ConvertBody(context.Background(), c, BodyConversionInput{
		From:   "storage",
		To:     "view",
		Value:  "<p>Hello</p>",
		Expand: []string{"webresource.uris.css"},
	})
	if err != nil {
		t.Fatalf("ConvertBody: %v", err)
	}
	if got.Representation != "view" || got.Value != "<p>Hello</p>" {
		t.Fatalf("unexpected conversion: %+v", got)
	}
	if len(got.WebResource.Keys) != 1 || got.WebResource.Keys[0] != "confluence.web.resources" {
		t.Fatalf("webresource keys = %+v", got.WebResource.Keys)
	}
}

func TestConvertBodyCloudUsesAsyncEndpointAndPollsResult(t *testing.T) {
	requests := make([]string, 0, 2)
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.Method+" "+r.URL.RequestURI())
		w.Header().Set("Content-Type", "application/json")

		switch len(requests) {
		case 1:
			if r.Method != http.MethodPost || r.URL.Path != "/wiki/rest/api/contentbody/convert/async/view" {
				t.Fatalf("queue request: %s %s", r.Method, r.URL.RequestURI())
			}
			q := r.URL.Query()
			if got := q.Get("spaceKeyContext"); got != "ENG" {
				t.Fatalf("spaceKeyContext = %q, want ENG", got)
			}
			if got := q.Get("contentIdContext"); got != "12345" {
				t.Fatalf("contentIdContext = %q, want 12345", got)
			}
			if got := q.Get("allowCache"); got != "false" {
				t.Fatalf("allowCache = %q, want false", got)
			}
			if got := q.Get("embeddedContentRender"); got != "current" {
				t.Fatalf("embeddedContentRender = %q, want current", got)
			}
			if got := q["expand"]; len(got) != 2 || got[0] != "webresource.uris.css" || got[1] != "webresource.uris.js" {
				t.Fatalf("expand = %v, want css/js values", got)
			}
			body := readJSONMap(t, r.Body)
			if got := body["representation"]; got != "storage" {
				t.Fatalf("representation = %v, want storage", got)
			}
			if got := body["value"]; got != "<p>Hello</p>" {
				t.Fatalf("value = %v, want <p>Hello</p>", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"asyncId": "async-1"})
		case 2:
			if r.Method != http.MethodGet || r.URL.Path != "/wiki/rest/api/contentbody/convert/async/async-1" {
				t.Fatalf("poll request: %s %s", r.Method, r.URL.RequestURI())
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"representation": "view",
				"value":          "<p>Hello</p>",
				"renderTaskId":   "async-1",
				"status":         "COMPLETE",
			})
		default:
			t.Fatalf("unexpected extra request: %s", r.URL.RequestURI())
		}
	})

	got, err := ConvertBody(context.Background(), c, BodyConversionInput{
		From:                  "storage",
		To:                    "view",
		Value:                 "<p>Hello</p>",
		Expand:                []string{"webresource.uris.css", "webresource.uris.js"},
		SpaceKeyContext:       "ENG",
		ContentIDContext:      "12345",
		AllowCache:            boolPtr(false),
		EmbeddedContentRender: "current",
		CloudPollAttempts:     1,
	})
	if err != nil {
		t.Fatalf("ConvertBody: %v", err)
	}
	if len(requests) != 2 {
		t.Fatalf("requests = %v, want queue and poll", requests)
	}
	if got.AsyncID != "async-1" || got.Representation != "view" || got.Value != "<p>Hello</p>" {
		t.Fatalf("unexpected conversion: %+v", got)
	}
}

func TestConvertBodyCloudReturnsAsyncIDWhenPollingDisabled(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/wiki/rest/api/contentbody/convert/async/view" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"asyncId": "async-2"})
	})

	got, err := ConvertBody(context.Background(), c, BodyConversionInput{
		From:              "storage",
		To:                "view",
		Value:             "<p>Hello</p>",
		CloudPollAttempts: 0,
	})
	if err != nil {
		t.Fatalf("ConvertBody: %v", err)
	}
	if got.AsyncID != "async-2" || got.Value != "" {
		t.Fatalf("unexpected async-only conversion: %+v", got)
	}
}

func boolPtr(v bool) *bool {
	return &v
}
