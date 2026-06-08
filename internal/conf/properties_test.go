package conf

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

func TestListContentPropertiesServerUsesV1Endpoint(t *testing.T) {
	_, c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/rest/api/content/12345/property" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		if got := r.URL.Query().Get("start"); got != "0" {
			t.Fatalf("start = %q, want 0", got)
		}
		if got := r.URL.Query().Get("limit"); got != "2" {
			t.Fatalf("limit = %q, want 2", got)
		}
		if got := r.URL.Query().Get("expand"); got != "version" {
			t.Fatalf("expand = %q, want version", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"results": []map[string]any{{
				"id":      "prop1",
				"key":     "foo",
				"value":   map[string]any{"enabled": true},
				"version": map[string]any{"number": 3},
			}},
			"size":   1,
			"_links": map[string]any{},
		})
	})

	got, err := ListContentProperties(context.Background(), c, "12345", "", 2)
	if err != nil {
		t.Fatalf("ListContentProperties: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len(properties) = %d, want 1", len(got))
	}
	if got[0].ID != "prop1" || got[0].Key != "foo" || got[0].Version.Number != 3 {
		t.Fatalf("unexpected property: %+v", got[0])
	}
	value, ok := got[0].Value.(map[string]any)
	if !ok || value["enabled"] != true {
		t.Fatalf("value = %#v, want enabled=true", got[0].Value)
	}
}

func TestListContentPropertiesCloudUsesV2PageEndpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/wiki/api/v2/pages/12345/properties" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		if got := r.URL.Query().Get("start"); got != "" {
			t.Fatalf("cloud v2 should not send start, got %q", got)
		}
		if got := r.URL.Query().Get("key"); got != "foo" {
			t.Fatalf("key = %q, want foo", got)
		}
		if got := r.URL.Query().Get("limit"); got != "2" {
			t.Fatalf("limit = %q, want 2", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"results": []map[string]any{{
				"id":      "prop1",
				"key":     "foo",
				"value":   map[string]any{"enabled": true},
				"version": map[string]any{"number": 3},
			}},
			"_links": map[string]any{},
		})
	})

	got, err := ListContentProperties(context.Background(), c, "12345", "foo", 2)
	if err != nil {
		t.Fatalf("ListContentProperties: %v", err)
	}
	if len(got) != 1 || got[0].ID != "prop1" || got[0].Key != "foo" {
		t.Fatalf("unexpected properties: %+v", got)
	}
}

func TestSetContentPropertyCloudResolvesKeyAndUpdatesByPropertyID(t *testing.T) {
	var requests []string
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.Method+" "+r.URL.RequestURI())
		w.Header().Set("Content-Type", "application/json")

		switch len(requests) {
		case 1:
			if r.Method != http.MethodGet || r.URL.Path != "/wiki/api/v2/pages/12345/properties" {
				t.Fatalf("lookup request: %s %s", r.Method, r.URL.RequestURI())
			}
			if got := r.URL.Query().Get("key"); got != "foo" {
				t.Fatalf("lookup key = %q, want foo", got)
			}
			if got := r.URL.Query().Get("limit"); got != "1" {
				t.Fatalf("lookup limit = %q, want 1", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{
					"id":      "prop1",
					"key":     "foo",
					"version": map[string]any{"number": 3},
				}},
				"_links": map[string]any{},
			})
		case 2:
			if r.Method != http.MethodPut || r.URL.Path != "/wiki/api/v2/pages/12345/properties/prop1" {
				t.Fatalf("update request: %s %s", r.Method, r.URL.RequestURI())
			}
			body := readJSONMap(t, r.Body)
			if got := body["key"]; got != "foo" {
				t.Fatalf("key = %v, want foo", got)
			}
			value := mapValue(t, body, "value")
			if got := value["enabled"]; got != true {
				t.Fatalf("value.enabled = %v, want true", got)
			}
			version := mapValue(t, body, "version")
			if got := version["number"]; got != float64(4) {
				t.Fatalf("version.number = %v, want 4", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":      "prop1",
				"key":     "foo",
				"value":   map[string]any{"enabled": true},
				"version": map[string]any{"number": 4},
			})
		default:
			t.Fatalf("unexpected extra request: %s", r.URL.RequestURI())
		}
	})

	got, err := SetContentProperty(context.Background(), c, "12345", "foo", map[string]any{"enabled": true})
	if err != nil {
		t.Fatalf("SetContentProperty: %v", err)
	}
	if got.ID != "prop1" || got.Version.Number != 4 {
		t.Fatalf("unexpected property: %+v", got)
	}
	if len(requests) != 2 {
		t.Fatalf("requests = %v, want lookup then update", requests)
	}
}

func TestSetSpacePropertyCloudResolvesSpaceIDAndPropertyID(t *testing.T) {
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
			if got := r.URL.Query().Get("limit"); got != "1" {
				t.Fatalf("space lookup limit = %q, want 1", got)
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
			if r.Method != http.MethodGet || r.URL.Path != "/wiki/api/v2/spaces/1001/properties" {
				t.Fatalf("property lookup request: %s %s", r.Method, r.URL.RequestURI())
			}
			if got := r.URL.Query().Get("key"); got != "foo" {
				t.Fatalf("property lookup key = %q, want foo", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{
					"id":      "spaceprop1",
					"key":     "foo",
					"version": map[string]any{"number": 7},
				}},
				"_links": map[string]any{},
			})
		case 3:
			if r.Method != http.MethodPut || r.URL.Path != "/wiki/api/v2/spaces/1001/properties/spaceprop1" {
				t.Fatalf("update request: %s %s", r.Method, r.URL.RequestURI())
			}
			body := readJSONMap(t, r.Body)
			value := mapValue(t, body, "value")
			if got := value["retention"]; got != "30d" {
				t.Fatalf("value.retention = %v, want 30d", got)
			}
			version := mapValue(t, body, "version")
			if got := version["number"]; got != float64(8) {
				t.Fatalf("version.number = %v, want 8", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":      "spaceprop1",
				"key":     "foo",
				"value":   map[string]any{"retention": "30d"},
				"version": map[string]any{"number": 8},
			})
		default:
			t.Fatalf("unexpected extra request: %s", r.URL.RequestURI())
		}
	})

	got, err := SetSpaceProperty(context.Background(), c, "ENG", "foo", map[string]any{"retention": "30d"})
	if err != nil {
		t.Fatalf("SetSpaceProperty: %v", err)
	}
	if got.ID != "spaceprop1" || got.Version.Number != 8 {
		t.Fatalf("unexpected property: %+v", got)
	}
	if len(requests) != 3 {
		t.Fatalf("requests = %v, want space lookup, property lookup, update", requests)
	}
}

func TestDeleteSpacePropertyServerUsesV1KeyEndpoint(t *testing.T) {
	_, c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/rest/api/space/ENG/property/foo" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		w.WriteHeader(http.StatusNoContent)
	})

	if err := DeleteSpaceProperty(context.Background(), c, "ENG", "foo"); err != nil {
		t.Fatalf("DeleteSpaceProperty: %v", err)
	}
}
