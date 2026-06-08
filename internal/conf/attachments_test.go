package conf

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

func attachmentTestClientWithFlavor(t *testing.T, flavor client.Flavor, handler http.HandlerFunc) (*httptest.Server, *client.Client) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	cfg := client.Config{BaseURL: srv.URL, Flavor: flavor, PAT: "x"}
	if flavor == client.FlavorCloud {
		cfg.BaseURL = srv.URL + "/wiki"
		cfg.Email = "x@example.com"
		cfg.APIToken = "cloud-token"
	}
	c, err := client.New(cfg)
	if err != nil {
		t.Fatalf("client.New: %v", err)
	}
	c.HTTPClient = &http.Client{Timeout: 5 * time.Second}
	return srv, c
}

func TestListAttachmentsServerUsesV1AttachmentEndpoint(t *testing.T) {
	_, c := attachmentTestClientWithFlavor(t, client.FlavorServer, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/content/12345/child/attachment" {
			t.Fatalf("path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("limit"); got != "2" {
			t.Fatalf("limit: %q", got)
		}
		if got := r.URL.Query().Get("expand"); got != "version,metadata,extensions" {
			t.Fatalf("expand: %q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"results": []Attachment{{
				ID:    "att1",
				Title: "report.pdf",
				Type:  "attachment",
			}},
		})
	})

	got, err := ListAttachments(context.Background(), c, "12345", 2)
	if err != nil {
		t.Fatalf("ListAttachments: %v", err)
	}
	if len(got) != 1 || got[0].Title != "report.pdf" {
		t.Fatalf("unexpected attachments: %+v", got)
	}
}

func TestListAttachmentsCloudUsesV2PageAttachmentEndpoint(t *testing.T) {
	_, c := attachmentTestClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/wiki/api/v2/pages/12345/attachments" {
			t.Fatalf("path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("limit"); got != "2" {
			t.Fatalf("limit: %q", got)
		}
		if got := r.URL.Query().Get("expand"); got != "" {
			t.Fatalf("cloud v2 should not send expand, got %q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"results": []map[string]any{{
				"id":           "att1",
				"title":        "report.pdf",
				"status":       "current",
				"mediaType":    "application/pdf",
				"comment":      "release report",
				"fileSize":     42,
				"downloadLink": "/download/attachments/12345/report.pdf?api=v2",
				"webuiLink":    "/pages/viewpageattachments.action?pageId=12345",
				"version": map[string]any{
					"number":    3,
					"createdAt": "2026-06-08T10:00:00Z",
				},
			}},
			"_links": map[string]any{},
		})
	})

	got, err := ListAttachments(context.Background(), c, "12345", 2)
	if err != nil {
		t.Fatalf("ListAttachments: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len(attachments) = %d, want 1", len(got))
	}
	if got[0].Title != "report.pdf" || got[0].Extensions.MediaType != "application/pdf" || got[0].Extensions.FileSize != 42 {
		t.Fatalf("unexpected attachment: %+v", got[0])
	}
	if got[0].Links.Download != "/download/attachments/12345/report.pdf?api=v2" {
		t.Fatalf("download link: %q", got[0].Links.Download)
	}
}

func TestListAttachmentsCloudFollowsLinkHeaderPagination(t *testing.T) {
	requests := make([]string, 0, 2)
	_, c := attachmentTestClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.URL.RequestURI())
		switch len(requests) {
		case 1:
			if r.URL.RequestURI() != "/wiki/api/v2/pages/12345/attachments?limit=2" {
				t.Fatalf("first request: %s", r.URL.RequestURI())
			}
			w.Header().Set("Link", "<http://"+r.Host+"/wiki/api/v2/pages/12345/attachments?cursor=abc>; rel=\"next\"")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{
					"id":        "att1",
					"title":     "first.pdf",
					"mediaType": "application/pdf",
					"fileSize":  10,
				}},
				"_links": map[string]any{},
			})
		case 2:
			if r.URL.RequestURI() != "/wiki/api/v2/pages/12345/attachments?cursor=abc" {
				t.Fatalf("second request: %s", r.URL.RequestURI())
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{
					"id":        "att2",
					"title":     "second.pdf",
					"mediaType": "application/pdf",
					"fileSize":  20,
				}},
				"_links": map[string]any{},
			})
		default:
			t.Fatalf("unexpected extra request: %s", r.URL.RequestURI())
		}
	})

	got, err := ListAttachments(context.Background(), c, "12345", 2)
	if err != nil {
		t.Fatalf("ListAttachments: %v", err)
	}
	if len(requests) != 2 {
		t.Fatalf("requests = %d, want 2 (%v)", len(requests), requests)
	}
	if len(got) != 2 || got[0].Title != "first.pdf" || got[1].Title != "second.pdf" {
		t.Fatalf("unexpected attachments: %+v", got)
	}
}

func TestAttachmentDownloadURLUsesLinksBaseWhenPresent(t *testing.T) {
	a := Attachment{}
	a.Links.Base = "https://example.atlassian.net/wiki"
	a.Links.Download = "/download/attachments/12345/report.pdf"
	if got := a.DownloadURL(); got != "https://example.atlassian.net/wiki/download/attachments/12345/report.pdf" {
		t.Fatalf("DownloadURL = %q", got)
	}
}

func TestAttachmentCloudFallbackKeepsDownloadLinkForDownloadCommand(t *testing.T) {
	_, c := attachmentTestClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/wiki/api/v2/pages/12345/attachments") {
			t.Fatalf("path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"results": []map[string]any{{
				"id":           "att1",
				"title":        "report.pdf",
				"downloadLink": "/download/attachments/12345/report.pdf?api=v2",
			}},
			"_links": map[string]any{},
		})
	})

	got, err := ListAttachments(context.Background(), c, "12345", 1)
	if err != nil {
		t.Fatalf("ListAttachments: %v", err)
	}
	if got[0].Links.Download == "" {
		t.Fatalf("expected download link fallback, got %+v", got[0])
	}
}

func TestDownloadAttachmentServerUsesAttachmentDownloadLink(t *testing.T) {
	_, c := attachmentTestClientWithFlavor(t, client.FlavorServer, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/download/attachments/12345/report.pdf" {
			t.Fatalf("path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte("server-bytes"))
	})

	got, err := DownloadAttachment(context.Background(), c, "12345", Attachment{
		ID: "att1",
		Links: struct {
			Download string `json:"download,omitempty"`
			WebUI    string `json:"webui,omitempty"`
			Base     string `json:"base,omitempty"`
		}{Download: "/download/attachments/12345/report.pdf"},
	})
	if err != nil {
		t.Fatalf("DownloadAttachment: %v", err)
	}
	if string(got) != "server-bytes" {
		t.Fatalf("DownloadAttachment = %q, want server-bytes", string(got))
	}
}

func TestDownloadAttachmentCloudUsesV1RawDownloadEndpoint(t *testing.T) {
	_, c := attachmentTestClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/wiki/download/attachments/12345/report.pdf":
			http.Error(w, "v2 metadata download link requires browser session", http.StatusUnauthorized)
		case "/wiki/rest/api/content/12345/child/attachment/att1/download":
			_, _ = w.Write([]byte("cloud-bytes"))
		default:
			t.Fatalf("path: %s", r.URL.Path)
		}
	})

	got, err := DownloadAttachment(context.Background(), c, "12345", Attachment{
		ID: "att1",
		Links: struct {
			Download string `json:"download,omitempty"`
			WebUI    string `json:"webui,omitempty"`
			Base     string `json:"base,omitempty"`
		}{Download: "/download/attachments/12345/report.pdf"},
	})
	if err != nil {
		t.Fatalf("DownloadAttachment: %v", err)
	}
	if string(got) != "cloud-bytes" {
		t.Fatalf("DownloadAttachment = %q, want cloud-bytes", string(got))
	}
}
