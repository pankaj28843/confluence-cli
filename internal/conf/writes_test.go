package conf

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

func TestCreatePageCloudUsesV2PageEndpointAndSpaceID(t *testing.T) {
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
				"_links": map[string]any{},
			})
		case 2:
			if r.Method != http.MethodPost || r.URL.Path != "/wiki/api/v2/pages" {
				t.Fatalf("create request: %s %s", r.Method, r.URL.RequestURI())
			}
			body := readJSONMap(t, r.Body)
			if got := body["spaceId"]; got != "1001" {
				t.Fatalf("spaceId = %v, want 1001; body=%v", got, body)
			}
			if got := body["title"]; got != "Created" {
				t.Fatalf("title = %v, want Created", got)
			}
			if got := body["parentId"]; got != "12345" {
				t.Fatalf("parentId = %v, want 12345", got)
			}
			if got := body["status"]; got != "current" {
				t.Fatalf("status = %v, want current", got)
			}
			if got := body["subtype"]; got != "live" {
				t.Fatalf("subtype = %v, want live", got)
			}
			pageBody := mapValue(t, body, "body")
			if got := pageBody["representation"]; got != "storage" {
				t.Fatalf("body.representation = %v, want storage", got)
			}
			if got := pageBody["value"]; got != "<p>x</p>" {
				t.Fatalf("body.value = %v, want <p>x</p>", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":       "67890",
				"status":   "current",
				"title":    "Created",
				"spaceId":  "1001",
				"parentId": "12345",
				"version":  map[string]any{"number": 1},
				"_links":   map[string]any{"webui": "/spaces/ENG/pages/67890"},
			})
		default:
			t.Fatalf("unexpected extra request: %s %s", r.Method, r.URL.RequestURI())
		}
	})

	got, err := CreatePage(context.Background(), c, CreatePageInput{
		SpaceKey: "ENG", Title: "Created", BodyValue: "<p>x</p>", ParentID: "12345",
	})
	if err != nil {
		t.Fatalf("CreatePage: %v", err)
	}
	if got.ID != "67890" || got.Title != "Created" || got.Version.Number != 1 {
		t.Fatalf("unexpected page: %+v", got)
	}
	if len(requests) != 2 {
		t.Fatalf("requests = %v, want space lookup then page create", requests)
	}
}

func TestUpdatePageCloudUsesV2PageEndpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/wiki/api/v2/pages/12345" {
			t.Fatalf("update request: %s %s", r.Method, r.URL.RequestURI())
		}
		body := readJSONMap(t, r.Body)
		if got := body["id"]; got != "12345" {
			t.Fatalf("id = %v, want 12345", got)
		}
		if got := body["status"]; got != "current" {
			t.Fatalf("status = %v, want current", got)
		}
		if got := body["title"]; got != "Updated" {
			t.Fatalf("title = %v, want Updated", got)
		}
		version := mapValue(t, body, "version")
		if got := version["number"]; got != float64(5) {
			t.Fatalf("version.number = %v, want 5", got)
		}
		pageBody := mapValue(t, body, "body")
		if got := pageBody["representation"]; got != "storage" {
			t.Fatalf("body.representation = %v, want storage", got)
		}
		if got := pageBody["value"]; got != "<p>updated</p>" {
			t.Fatalf("body.value = %v, want <p>updated</p>", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":      "12345",
			"status":  "current",
			"title":   "Updated",
			"version": map[string]any{"number": 5},
		})
	})

	got, err := UpdatePage(context.Background(), c, UpdatePageInput{
		ID: "12345", Title: "Updated", BodyValue: "<p>updated</p>", VersionNumber: 4,
	})
	if err != nil {
		t.Fatalf("UpdatePage: %v", err)
	}
	if got.ID != "12345" || got.Title != "Updated" || got.Version.Number != 5 {
		t.Fatalf("unexpected page: %+v", got)
	}
}

func TestPutAttachmentCloudKeepsV1MultipartUploadSurface(t *testing.T) {
	var requests []string
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.Method+" "+r.URL.RequestURI())
		w.Header().Set("Content-Type", "application/json")

		switch len(requests) {
		case 1:
			if r.Method != http.MethodGet || r.URL.Path != "/wiki/api/v2/pages/12345/attachments" {
				t.Fatalf("attachment lookup request: %s %s", r.Method, r.URL.RequestURI())
			}
			if got := r.URL.Query().Get("limit"); got != "200" {
				t.Fatalf("lookup limit = %q, want 200", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{
					"id":    "att1",
					"title": "report.pdf",
				}},
				"_links": map[string]any{},
			})
		case 2:
			if r.Method != http.MethodPost || r.URL.Path != "/wiki/rest/api/content/12345/child/attachment/att1/data" {
				t.Fatalf("attachment upload request: %s %s", r.Method, r.URL.RequestURI())
			}
			if got := r.Header.Get("X-Atlassian-Token"); got != "no-check" {
				t.Fatalf("X-Atlassian-Token = %q, want no-check", got)
			}
			if err := r.ParseMultipartForm(1 << 20); err != nil {
				t.Fatalf("ParseMultipartForm: %v", err)
			}
			file, _, err := r.FormFile("file")
			if err != nil {
				t.Fatalf("multipart file field: %v", err)
			}
			data, err := io.ReadAll(file)
			if err != nil {
				t.Fatalf("read multipart file: %v", err)
			}
			if string(data) != "PDF-BYTES" {
				t.Fatalf("uploaded bytes = %q, want PDF-BYTES", string(data))
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []Attachment{{ID: "att1", Title: "report.pdf"}},
			})
		default:
			t.Fatalf("unexpected extra request: %s %s", r.Method, r.URL.RequestURI())
		}
	})

	got, err := PutAttachment(context.Background(), c, "12345", "report.pdf", strings.NewReader("PDF-BYTES"), "v2")
	if err != nil {
		t.Fatalf("PutAttachment: %v", err)
	}
	if len(got) != 1 || got[0].ID != "att1" {
		t.Fatalf("unexpected attachments: %+v", got)
	}
	if len(requests) != 2 {
		t.Fatalf("requests = %v, want lookup then v1 multipart update", requests)
	}
}

func readJSONMap(t *testing.T, r io.Reader) map[string]any {
	t.Helper()
	var body map[string]any
	if err := json.NewDecoder(r).Decode(&body); err != nil {
		t.Fatalf("decode JSON body: %v", err)
	}
	return body
}

func mapValue(t *testing.T, body map[string]any, key string) map[string]any {
	t.Helper()
	value, ok := body[key].(map[string]any)
	if !ok {
		t.Fatalf("%s = %T %v, want object", key, body[key], body[key])
	}
	return value
}
