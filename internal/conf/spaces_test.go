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

func testClient(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *client.Client) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	c, err := client.New(client.Config{BaseURL: srv.URL, Flavor: client.FlavorServer, PAT: "x"})
	if err != nil {
		t.Fatalf("client.New: %v", err)
	}
	c.HTTPClient = &http.Client{Timeout: 5 * time.Second}
	return srv, c
}

func TestListSpaces(t *testing.T) {
	_, c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/rest/api/space") {
			t.Fatalf("path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"results": []Space{{Key: "ENG", Name: "Engineering"}, {Key: "OPS", Name: "Operations"}},
			"size":    2,
		})
	})
	got, err := ListSpaces(context.Background(), c, SpaceFilter{Limit: 5})
	if err != nil || len(got) != 2 {
		t.Fatalf("ListSpaces: %v %+v", err, got)
	}
}

func TestGetSpace(t *testing.T) {
	_, c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/space/ENG" {
			t.Fatalf("path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(Space{Key: "ENG", Name: "Engineering"})
	})
	s, err := GetSpace(context.Background(), c, "ENG")
	if err != nil || s.Name != "Engineering" {
		t.Fatalf("GetSpace: %v %+v", err, s)
	}
}
