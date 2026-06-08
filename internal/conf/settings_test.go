package conf

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

func TestListThemesCloudUsesV1EndpointAndPagination(t *testing.T) {
	requests := make([]string, 0, 2)
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.URL.RequestURI())
		if r.Method != http.MethodGet || r.URL.Path != "/wiki/rest/api/settings/theme" {
			t.Fatalf("request: %s %s", r.Method, r.URL.RequestURI())
		}
		q := r.URL.Query()
		if got := q.Get("limit"); got != "2" {
			t.Fatalf("limit = %q, want 2", got)
		}
		switch q.Get("start") {
		case "0":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{"themeKey": "default", "name": "Default"}},
				"start":   0,
				"limit":   2,
				"size":    1,
				"_links":  map[string]any{"next": "/wiki/rest/api/settings/theme?start=1&limit=2"},
			})
		case "1":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{{"themeKey": "custom", "name": "Custom"}},
				"start":   1,
				"limit":   2,
				"size":    1,
				"_links":  map[string]any{},
			})
		default:
			t.Fatalf("start = %q", q.Get("start"))
		}
	})

	got, err := ListThemes(context.Background(), c, ThemeListOptions{Limit: 2})
	if err != nil {
		t.Fatalf("ListThemes: %v", err)
	}
	if len(requests) != 2 {
		t.Fatalf("requests = %v, want two paginated requests", requests)
	}
	if len(got) != 2 || got[0].ThemeKey != "default" || got[1].ThemeKey != "custom" {
		t.Fatalf("themes = %+v", got)
	}
}

func TestThemeReadsCloudUseV1Endpoints(t *testing.T) {
	requests := make([]string, 0, 3)
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.URL.RequestURI())
		switch r.URL.Path {
		case "/wiki/rest/api/settings/theme/selected":
			_ = json.NewEncoder(w).Encode(map[string]any{"themeKey": "global", "name": "Global"})
		case "/wiki/rest/api/settings/theme/global":
			_ = json.NewEncoder(w).Encode(map[string]any{"themeKey": "global", "name": "Global"})
		case "/wiki/rest/api/space/ENG/theme":
			_ = json.NewEncoder(w).Encode(map[string]any{"themeKey": "space-theme", "name": "Space Theme"})
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
	})

	global, err := GetGlobalTheme(context.Background(), c)
	if err != nil {
		t.Fatalf("GetGlobalTheme: %v", err)
	}
	view, err := GetTheme(context.Background(), c, "global")
	if err != nil {
		t.Fatalf("GetTheme: %v", err)
	}
	space, err := GetSpaceTheme(context.Background(), c, "ENG")
	if err != nil {
		t.Fatalf("GetSpaceTheme: %v", err)
	}
	if len(requests) != 3 {
		t.Fatalf("requests = %v, want three requests", requests)
	}
	if global.ThemeKey != "global" || view.ThemeKey != "global" || space.ThemeKey != "space-theme" {
		t.Fatalf("themes = global:%+v view:%+v space:%+v", global, view, space)
	}
}

func TestSettingsReadsCloudUseV1Endpoints(t *testing.T) {
	requests := make([]string, 0, 3)
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.URL.RequestURI())
		switch r.URL.Path {
		case "/wiki/rest/api/settings/systemInfo":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"cloudId":         "cloud-1",
				"commitHash":      "abcdef",
				"siteTitle":       "Example Wiki",
				"defaultLocale":   "en-US",
				"defaultTimeZone": "Europe/Copenhagen",
			})
		case "/wiki/rest/api/settings/lookandfeel":
			if got := r.URL.Query().Get("spaceKey"); got != "ENG" {
				t.Fatalf("spaceKey = %q, want ENG", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"headings":           map[string]any{"color": "#172B4D"},
				"links":              map[string]any{"color": "#0052CC"},
				"bordersAndDividers": map[string]any{"color": "#DFE1E6"},
			})
		case "/wiki/rest/api/space/ENG/settings":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"routeOverrideEnabled": true,
				"spaceKey":             "ENG",
				"editor":               map[string]any{"page": "v2", "blogpost": "v2", "default": "v2"},
				"_links":               map[string]any{"self": "/settings"},
			})
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
	})

	systemInfo, err := GetSystemInfo(context.Background(), c)
	if err != nil {
		t.Fatalf("GetSystemInfo: %v", err)
	}
	lookAndFeel, err := GetLookAndFeelSettings(context.Background(), c, "ENG")
	if err != nil {
		t.Fatalf("GetLookAndFeelSettings: %v", err)
	}
	spaceSettings, err := GetSpaceSettings(context.Background(), c, "ENG")
	if err != nil {
		t.Fatalf("GetSpaceSettings: %v", err)
	}
	if len(requests) != 3 {
		t.Fatalf("requests = %v, want three requests", requests)
	}
	if systemInfo.CloudID != "cloud-1" || systemInfo.SiteTitle != "Example Wiki" {
		t.Fatalf("system info = %+v", systemInfo)
	}
	if lookAndFeel.Headings.Color != "#172B4D" || lookAndFeel.Links.Color != "#0052CC" {
		t.Fatalf("look and feel = %+v", lookAndFeel)
	}
	if !spaceSettings.RouteOverrideEnabled || spaceSettings.SpaceKey != "ENG" || spaceSettings.Editor.Default != "v2" {
		t.Fatalf("space settings = %+v", spaceSettings)
	}
}

func TestSettingsAndThemesRejectServerFlavorAndMissingInputs(t *testing.T) {
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
		{name: "server list themes", err: firstThemesError(ListThemes(context.Background(), server, ThemeListOptions{})), want: "Confluence Cloud"},
		{name: "server global theme", err: firstThemeError(GetGlobalTheme(context.Background(), server)), want: "Confluence Cloud"},
		{name: "server theme", err: firstThemeError(GetTheme(context.Background(), server, "global")), want: "Confluence Cloud"},
		{name: "server space theme", err: firstThemeError(GetSpaceTheme(context.Background(), server, "ENG")), want: "Confluence Cloud"},
		{name: "server system info", err: firstSystemInfoError(GetSystemInfo(context.Background(), server)), want: "Confluence Cloud"},
		{name: "server look and feel", err: firstLookAndFeelError(GetLookAndFeelSettings(context.Background(), server, "")), want: "Confluence Cloud"},
		{name: "server space settings", err: firstSpaceSettingsError(GetSpaceSettings(context.Background(), server, "ENG")), want: "Confluence Cloud"},
		{name: "missing theme key", err: firstThemeError(GetTheme(context.Background(), cloud, "")), want: "theme key"},
		{name: "missing theme space", err: firstThemeError(GetSpaceTheme(context.Background(), cloud, "")), want: "space key"},
		{name: "missing settings space", err: firstSpaceSettingsError(GetSpaceSettings(context.Background(), cloud, "")), want: "space key"},
	}
	for _, tc := range cases {
		if tc.err == nil || !strings.Contains(tc.err.Error(), tc.want) {
			t.Fatalf("%s error = %v, want %q", tc.name, tc.err, tc.want)
		}
	}
}

func firstThemesError(_ []Theme, err error) error {
	return err
}

func firstThemeError(_ *Theme, err error) error {
	return err
}

func firstSystemInfoError(_ *SystemInfo, err error) error {
	return err
}

func firstLookAndFeelError(_ *LookAndFeelSettings, err error) error {
	return err
}

func firstSpaceSettingsError(_ *SpaceSettings, err error) error {
	return err
}
