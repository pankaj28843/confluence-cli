package client

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func newTestServer(t *testing.T, flavor Flavor, handler http.HandlerFunc) (*httptest.Server, *Client) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	cfg := Config{BaseURL: srv.URL, Flavor: flavor, PAT: "server-pat", Email: "x@example.com", APIToken: "cloud-token"}
	c, err := New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	c.HTTPClient = &http.Client{
		Timeout:   5 * time.Second,
		Transport: wrapTransport(http.DefaultTransport, &c.Debug),
	}
	return srv, c
}

func TestAuthHeaderBearerOnServer(t *testing.T) {
	var got string
	srv, c := newTestServer(t, FlavorServer, func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Get("Authorization")
		_, _ = w.Write([]byte(`{"ok":true}`))
	})
	defer srv.Close()
	_, _, err := c.Get(context.Background(), "/rest/api/user/current", nil)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	want := "Bearer server-pat"
	if got != want {
		t.Fatalf("Server Authorization: got %q want %q", got, want)
	}
}

func TestAuthHeaderBasicOnCloud(t *testing.T) {
	var got string
	srv, c := newTestServer(t, FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Get("Authorization")
		_, _ = w.Write([]byte(`{"ok":true}`))
	})
	defer srv.Close()
	_, _, err := c.Get(context.Background(), "/rest/api/user/current", nil)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	want := "Basic " + base64.StdEncoding.EncodeToString([]byte("x@example.com:cloud-token"))
	if got != want {
		t.Fatalf("Cloud Authorization: got %q want %q", got, want)
	}
	if strings.HasPrefix(got, "Bearer ") {
		t.Fatalf("Bearer used on Cloud — Cloud must be Basic")
	}
}

func TestFlavorDetection(t *testing.T) {
	cases := []struct {
		in   string
		want Flavor
	}{
		{"https://wiki.example.com", FlavorServer},
		{"https://confluence.example.com/confluence", FlavorServer},
		{"https://example.atlassian.net/wiki", FlavorCloud},
		{"https://example.atlassian.net", FlavorCloud},
		{"https://wiki.internal.example/wiki", FlavorCloud},
	}
	for _, c := range cases {
		if got := DetectFlavor(c.in); got != c.want {
			t.Errorf("DetectFlavor(%q) = %v want %v", c.in, got, c.want)
		}
	}
}

func TestRetryOn5xx(t *testing.T) {
	var count int32
	srv, c := newTestServer(t, FlavorServer, func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&count, 1) <= 2 {
			http.Error(w, "boom", http.StatusServiceUnavailable)
			return
		}
		_, _ = w.Write([]byte(`{"ok":true}`))
	})
	defer srv.Close()

	if _, _, err := c.Get(context.Background(), "/rest/api/user/current", nil); err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got := atomic.LoadInt32(&count); got != 3 {
		t.Fatalf("want 3 attempts, got %d", got)
	}
}

func TestRetryAfterHonoured(t *testing.T) {
	var count int32
	var times []time.Time
	srv, c := newTestServer(t, FlavorServer, func(w http.ResponseWriter, r *http.Request) {
		times = append(times, time.Now())
		if atomic.AddInt32(&count, 1) == 1 {
			w.Header().Set("Retry-After", "1")
			http.Error(w, "slow down", http.StatusTooManyRequests)
			return
		}
		_, _ = w.Write([]byte(`{"ok":true}`))
	})
	defer srv.Close()
	if _, _, err := c.Get(context.Background(), "/rest/api/user/current", nil); err != nil {
		t.Fatalf("Get: %v", err)
	}
	if len(times) >= 2 && times[1].Sub(times[0]) < 900*time.Millisecond {
		t.Fatalf("Retry-After not honoured")
	}
}

func TestUnauthorizedIsUserFixable(t *testing.T) {
	srv, c := newTestServer(t, FlavorServer, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad pat", http.StatusUnauthorized)
	})
	defer srv.Close()
	_, _, err := c.Get(context.Background(), "/rest/api/user/current", nil)
	if err == nil || !IsUserFixable(err) {
		t.Fatalf("401 should be user-fixable, got %v", err)
	}
}

func TestBuildURLDoesNotDoublePrefixWiki(t *testing.T) {
	c := &Client{BaseURL: "https://example.atlassian.net/wiki", Flavor: FlavorCloud}
	got := c.BuildURL("/rest/api/user/current", nil)
	want := "https://example.atlassian.net/wiki/rest/api/user/current"
	if got != want {
		t.Fatalf("BuildURL: got %s want %s", got, want)
	}
}
