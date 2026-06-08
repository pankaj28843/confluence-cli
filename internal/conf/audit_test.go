package conf

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

func TestListAuditRecordsCloudPaginatesAndFilters(t *testing.T) {
	var seen []string
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		seen = append(seen, r.URL.String())
		if r.URL.Path != "/wiki/rest/api/audit" {
			t.Fatalf("path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("startDate") != "1700000000000" || q.Get("endDate") != "1700100000000" || q.Get("searchString") != "space" {
			t.Fatalf("filters: %s", r.URL.RawQuery)
		}
		if q.Get("limit") != "2" {
			t.Fatalf("limit: %q", q.Get("limit"))
		}
		switch q.Get("start") {
		case "0":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []AuditRecord{{CreationDate: 1700000000001, Summary: "space exported", Category: "space", Author: User{AccountID: "acct-1", DisplayName: "Ada"}}},
				"_links":  map[string]any{"next": "/wiki/rest/api/audit?start=1&limit=2&startDate=1700000000000&endDate=1700100000000&searchString=space"},
			})
		case "1":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []AuditRecord{{CreationDate: 1700000000002, Summary: "space permission changed", Category: "permissions"}},
				"_links":  map[string]any{},
			})
		default:
			t.Fatalf("start: %q", q.Get("start"))
		}
	})

	got, err := ListAuditRecords(context.Background(), c, AuditListOptions{
		StartDate:    "1700000000000",
		EndDate:      "1700100000000",
		SearchString: "space",
		Limit:        2,
	})
	if err != nil {
		t.Fatalf("ListAuditRecords: %v", err)
	}
	if len(got) != 2 || got[0].Summary != "space exported" || got[0].Author.AccountID != "acct-1" {
		t.Fatalf("records = %+v", got)
	}
	if len(seen) != 2 {
		t.Fatalf("requests = %v", seen)
	}
}

func TestListAuditRecordsSinceCloudUsesSinceEndpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/wiki/rest/api/audit/since" {
			t.Fatalf("path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("number") != "7" || q.Get("units") != "DAYS" || q.Get("searchString") != "group" || q.Get("start") != "0" || q.Get("limit") != "3" {
			t.Fatalf("query: %s", r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"results": []AuditRecord{{CreationDate: 1700000000003, Summary: "group membership changed"}},
			"_links":  map[string]any{},
		})
	})

	got, err := ListAuditRecords(context.Background(), c, AuditListOptions{
		SinceNumber:  7,
		SinceUnit:    "DAYS",
		SearchString: "group",
		Limit:        3,
	})
	if err != nil {
		t.Fatalf("ListAuditRecords since: %v", err)
	}
	if len(got) != 1 || got[0].Summary != "group membership changed" {
		t.Fatalf("records = %+v", got)
	}
}

func TestListAuditRecordsServerUsesDeprecatedReadEndpoint(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorServer, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/audit" {
			t.Fatalf("path: %s", r.URL.Path)
		}
		if r.URL.RawQuery != "" {
			t.Fatalf("server query should be empty: %s", r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"results": []AuditRecord{
				{CreationDate: 1700000000004, Summary: "first"},
				{CreationDate: 1700000000005, Summary: "second"},
			},
		})
	})

	got, err := ListAuditRecords(context.Background(), c, AuditListOptions{Limit: 1})
	if err != nil {
		t.Fatalf("ListAuditRecords server: %v", err)
	}
	if len(got) != 1 || got[0].Summary != "first" {
		t.Fatalf("records = %+v", got)
	}
}

func TestGetAuditRetentionCloud(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/wiki/rest/api/audit/retention" {
			t.Fatalf("path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(AuditRetention{Number: 6, Units: "MONTHS"})
	})

	got, err := GetAuditRetention(context.Background(), c)
	if err != nil {
		t.Fatalf("GetAuditRetention: %v", err)
	}
	if got.Number != 6 || got.Units != "MONTHS" {
		t.Fatalf("retention = %+v", got)
	}
}

func TestAuditHelpersRejectUnsupportedInputs(t *testing.T) {
	_, server := testClientWithFlavor(t, client.FlavorServer, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected server request: %s", r.URL.String())
	})
	_, cloud := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected cloud request: %s", r.URL.String())
	})

	cases := []struct {
		name string
		err  error
		want string
	}{
		{name: "server filters", err: firstAuditError(ListAuditRecords(context.Background(), server, AuditListOptions{SearchString: "space"})), want: "unsupported"},
		{name: "server since", err: firstAuditError(ListAuditRecords(context.Background(), server, AuditListOptions{SinceNumber: 3, SinceUnit: "MONTHS"})), want: "unsupported"},
		{name: "cloud since missing number", err: firstAuditError(ListAuditRecords(context.Background(), cloud, AuditListOptions{SinceUnit: "DAYS"})), want: "number"},
		{name: "server retention", err: firstAuditRetentionError(GetAuditRetention(context.Background(), server)), want: "Cloud"},
	}
	for _, tc := range cases {
		if tc.err == nil || !strings.Contains(tc.err.Error(), tc.want) {
			t.Fatalf("%s error = %v, want %q", tc.name, tc.err, tc.want)
		}
	}
}

func firstAuditError(_ []AuditRecord, err error) error {
	return err
}

func firstAuditRetentionError(_ *AuditRetention, err error) error {
	return err
}
