package conf

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

// AuditRecord is one row from the Confluence audit log.
type AuditRecord struct {
	Author            User                `json:"author,omitempty"`
	RemoteAddress     string              `json:"remoteAddress,omitempty"`
	CreationDate      int64               `json:"creationDate,omitempty"`
	Summary           string              `json:"summary,omitempty"`
	Description       string              `json:"description,omitempty"`
	Category          string              `json:"category,omitempty"`
	SysAdmin          bool                `json:"sysAdmin,omitempty"`
	SuperAdmin        bool                `json:"superAdmin,omitempty"`
	AffectedObject    AuditObject         `json:"affectedObject,omitempty"`
	ChangedValues     []AuditChangedValue `json:"changedValues,omitempty"`
	AssociatedObjects []AuditObject       `json:"associatedObjects,omitempty"`
	Links             map[string]any      `json:"_links,omitempty"`
}

// AuditObject identifies an object associated with an audit record.
type AuditObject struct {
	Name       string `json:"name,omitempty"`
	ObjectType string `json:"objectType,omitempty"`
}

// AuditChangedValue describes one field change in an audit record.
type AuditChangedValue struct {
	Name           string `json:"name,omitempty"`
	OldValue       string `json:"oldValue,omitempty"`
	HiddenOldValue string `json:"hiddenOldValue,omitempty"`
	NewValue       string `json:"newValue,omitempty"`
	HiddenNewValue string `json:"hiddenNewValue,omitempty"`
}

// AuditRetention is the Cloud audit retention period.
type AuditRetention struct {
	Number int    `json:"number"`
	Units  string `json:"units"`
}

// AuditListOptions controls audit-log listing.
type AuditListOptions struct {
	StartDate    string
	EndDate      string
	SearchString string
	SinceNumber  int64
	SinceUnit    string
	Limit        int
}

// ListAuditRecords lists audit records. Cloud supports documented filters and
// time-period reads; Server/Data Center exposes only a deprecated read endpoint
// in the current official OpenAPI, so filters are rejected there.
func ListAuditRecords(ctx context.Context, c *client.Client, opts AuditListOptions) ([]AuditRecord, error) {
	if c == nil {
		return nil, fmt.Errorf("ListAuditRecords: client is required")
	}
	opts.Limit = clampLimit(opts.Limit)
	if isCloud(c) {
		return listAuditRecordsCloud(ctx, c, opts)
	}
	if opts.hasCloudOnlyFilters() {
		return nil, fmt.Errorf("unsupported Server/Data Center audit filters")
	}
	raw, _, err := c.Get(ctx, "/rest/api/audit", nil)
	if err != nil {
		return nil, err
	}
	records, err := parseAuditRecords(raw)
	if err != nil {
		return nil, err
	}
	return capAuditRecords(records, opts.Limit), nil
}

// GetAuditRetention returns the Cloud audit retention period.
func GetAuditRetention(ctx context.Context, c *client.Client) (*AuditRetention, error) {
	if c == nil {
		return nil, fmt.Errorf("GetAuditRetention: client is required")
	}
	if !isCloud(c) {
		return nil, fmt.Errorf("GetAuditRetention: Confluence Cloud only")
	}
	raw, _, err := c.Get(ctx, "/rest/api/audit/retention", nil)
	if err != nil {
		return nil, err
	}
	var retention AuditRetention
	if err := json.Unmarshal(raw, &retention); err != nil {
		return nil, fmt.Errorf("parse audit retention: %w", err)
	}
	return &retention, nil
}

func listAuditRecordsCloud(ctx context.Context, c *client.Client, opts AuditListOptions) ([]AuditRecord, error) {
	path := "/rest/api/audit"
	params := url.Values{
		"start": {"0"},
		"limit": {strconv.Itoa(opts.Limit)},
	}
	if opts.SearchString != "" {
		params.Set("searchString", opts.SearchString)
	}
	if opts.SinceNumber > 0 || opts.SinceUnit != "" {
		if opts.SinceNumber <= 0 {
			return nil, fmt.Errorf("audit since number is required")
		}
		if opts.SinceUnit == "" {
			opts.SinceUnit = "MONTHS"
		}
		if opts.StartDate != "" || opts.EndDate != "" {
			return nil, fmt.Errorf("audit since does not support start-date or end-date filters")
		}
		path = "/rest/api/audit/since"
		params.Set("number", strconv.FormatInt(opts.SinceNumber, 10))
		params.Set("units", strings.ToUpper(opts.SinceUnit))
	} else {
		if opts.StartDate != "" {
			params.Set("startDate", opts.StartDate)
		}
		if opts.EndDate != "" {
			params.Set("endDate", opts.EndDate)
		}
	}
	return collectAuditPages(ctx, c, path, params, opts.Limit)
}

func collectAuditPages(ctx context.Context, c *client.Client, path string, params url.Values, limit int) ([]AuditRecord, error) {
	out := make([]AuditRecord, 0, limit)
	for len(out) < limit {
		raw, headers, err := c.Get(ctx, path, params)
		if err != nil {
			return nil, err
		}
		var page struct {
			Results []AuditRecord `json:"results"`
			Links   struct {
				Next string `json:"next,omitempty"`
			} `json:"_links"`
		}
		if err := json.Unmarshal(raw, &page); err != nil {
			return nil, fmt.Errorf("parse audit records: %w", err)
		}
		for _, record := range page.Results {
			out = append(out, record)
			if len(out) == limit {
				break
			}
		}
		next := nextPageURL(headers, page.Links.Next)
		if next == "" || len(page.Results) == 0 || len(out) == limit {
			break
		}
		var nextErr error
		path, params, nextErr = nextPageRequest(c, next)
		if nextErr != nil {
			return nil, nextErr
		}
	}
	return out, nil
}

func parseAuditRecords(raw []byte) ([]AuditRecord, error) {
	var page struct {
		Results []AuditRecord `json:"results"`
	}
	pageErr := json.Unmarshal(raw, &page)
	if pageErr == nil && page.Results != nil {
		return page.Results, nil
	}
	var records []AuditRecord
	recordsErr := json.Unmarshal(raw, &records)
	if recordsErr == nil {
		return records, nil
	}
	return nil, fmt.Errorf("parse audit records: page=%v array=%v", pageErr, recordsErr)
}

func capAuditRecords(records []AuditRecord, limit int) []AuditRecord {
	if limit <= 0 {
		limit = 25
	}
	if len(records) > limit {
		return records[:limit]
	}
	return records
}

func (o AuditListOptions) hasCloudOnlyFilters() bool {
	return o.StartDate != "" ||
		o.EndDate != "" ||
		o.SearchString != "" ||
		o.SinceNumber > 0 ||
		o.SinceUnit != ""
}
