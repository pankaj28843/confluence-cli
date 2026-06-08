package conf

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

func TestBulkGetUsersCloudUsesV2UsersBulk(t *testing.T) {
	_, c := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method: %s", r.Method)
		}
		if r.URL.Path != "/wiki/api/v2/users-bulk" {
			t.Fatalf("path: %s", r.URL.Path)
		}
		var body struct {
			AccountIDs []string `json:"accountIds"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if strings.Join(body.AccountIDs, ",") != "acct-1,acct-2" {
			t.Fatalf("accountIds = %#v", body.AccountIDs)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"results": []User{{
				AccountID:              "acct-1",
				DisplayName:            "Ada Lovelace",
				PublicName:             "Ada",
				Email:                  "ada@example.com",
				TimeZone:               "Europe/London",
				PersonalSpaceID:        "123",
				IsExternalCollaborator: true,
				AccountStatus:          "active",
				AccountType:            "atlassian",
			}},
		})
	})

	got, err := BulkGetUsers(context.Background(), c, []string{" acct-1 ", "", "acct-2"})
	if err != nil {
		t.Fatalf("BulkGetUsers: %v", err)
	}
	if len(got) != 1 || got[0].AccountID != "acct-1" || got[0].TimeZone != "Europe/London" || !got[0].IsExternalCollaborator {
		t.Fatalf("users = %+v", got)
	}
}

func TestBulkGetUsersRejectsUnsupportedInputs(t *testing.T) {
	_, cloud := testClientWithFlavor(t, client.FlavorCloud, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected cloud request: %s", r.URL.String())
	})
	_, server := testClientWithFlavor(t, client.FlavorServer, func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected server request: %s", r.URL.String())
	})

	if _, err := BulkGetUsers(context.Background(), server, []string{"acct-1"}); err == nil || !strings.Contains(err.Error(), "unsupported") {
		t.Fatalf("server error = %v", err)
	}
	if _, err := BulkGetUsers(context.Background(), cloud, nil); err == nil || !strings.Contains(err.Error(), "account id") {
		t.Fatalf("empty error = %v", err)
	}
}
