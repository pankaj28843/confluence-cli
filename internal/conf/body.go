package conf

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

const defaultCloudBodyPollInterval = 500 * time.Millisecond

// BodyConversionInput describes a content-body conversion request.
type BodyConversionInput struct {
	From                  string
	To                    string
	Value                 string
	Expand                []string
	SpaceKeyContext       string
	ContentIDContext      string
	AllowCache            *bool
	EmbeddedContentRender string
	CloudPollAttempts     int
	CloudPollInterval     time.Duration
}

// BodyConversion is a normalized content-body conversion response.
type BodyConversion struct {
	AsyncID        string          `json:"asyncId,omitempty"`
	Representation string          `json:"representation,omitempty"`
	Value          string          `json:"value,omitempty"`
	RenderTaskID   string          `json:"renderTaskId,omitempty"`
	Status         string          `json:"status,omitempty"`
	Error          string          `json:"error,omitempty"`
	WebResource    WebResourceInfo `json:"webresource,omitempty"`
}

// WebResourceInfo contains web resource dependency details returned by body conversion.
type WebResourceInfo struct {
	Keys     []string       `json:"keys,omitempty"`
	Contexts []string       `json:"contexts,omitempty"`
	URIs     map[string]any `json:"uris,omitempty"`
	Tags     map[string]any `json:"tags,omitempty"`
}

// ConvertBody converts a Confluence content body between documented representations.
func ConvertBody(ctx context.Context, c *client.Client, in BodyConversionInput) (*BodyConversion, error) {
	if in.From == "" {
		return nil, fmt.Errorf("ConvertBody: from representation is required")
	}
	if in.To == "" {
		return nil, fmt.Errorf("ConvertBody: to representation is required")
	}
	if in.Value == "" {
		return nil, fmt.Errorf("ConvertBody: value is required")
	}
	if isCloud(c) {
		return convertBodyCloud(ctx, c, in)
	}
	return convertBodyServer(ctx, c, in)
}

func convertBodyServer(ctx context.Context, c *client.Client, in BodyConversionInput) (*BodyConversion, error) {
	params := url.Values{}
	addValues(params, "expand", in.Expand)
	body := map[string]any{
		"representation": in.From,
		"value":          in.Value,
	}
	data, _, err := c.Post(ctx, "/rest/api/contentbody/convert/"+url.PathEscape(in.To), params, body)
	if err != nil {
		return nil, err
	}
	return parseBodyConversion(data, "body conversion")
}

func convertBodyCloud(ctx context.Context, c *client.Client, in BodyConversionInput) (*BodyConversion, error) {
	params := url.Values{}
	addValues(params, "expand", in.Expand)
	if in.SpaceKeyContext != "" {
		params.Set("spaceKeyContext", in.SpaceKeyContext)
	}
	if in.ContentIDContext != "" {
		params.Set("contentIdContext", in.ContentIDContext)
	}
	if in.AllowCache != nil {
		params.Set("allowCache", strconv.FormatBool(*in.AllowCache))
	}
	if in.EmbeddedContentRender != "" {
		params.Set("embeddedContentRender", in.EmbeddedContentRender)
	}

	body := map[string]any{
		"representation": in.From,
		"value":          in.Value,
	}
	data, _, err := c.Post(ctx, "/rest/api/contentbody/convert/async/"+url.PathEscape(in.To), params, body)
	if err != nil {
		return nil, err
	}
	var queued struct {
		AsyncID string `json:"asyncId"`
	}
	if err := json.Unmarshal(data, &queued); err != nil {
		return nil, fmt.Errorf("parse body conversion queue: %w", err)
	}
	if queued.AsyncID == "" {
		return nil, fmt.Errorf("body conversion queue response missing asyncId")
	}

	attempts := in.CloudPollAttempts
	if attempts < 0 {
		attempts = 0
	}
	if attempts == 0 {
		return &BodyConversion{AsyncID: queued.AsyncID}, nil
	}
	interval := in.CloudPollInterval
	if interval == 0 {
		interval = defaultCloudBodyPollInterval
	}

	var last *BodyConversion
	for attempt := 0; attempt < attempts; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(interval):
			}
		}
		got, err := getCloudBodyConversion(ctx, c, queued.AsyncID)
		if err != nil {
			return nil, err
		}
		got.AsyncID = queued.AsyncID
		last = got
		if got.Value != "" || got.Error != "" {
			return got, nil
		}
	}
	if last != nil {
		return last, nil
	}
	return &BodyConversion{AsyncID: queued.AsyncID}, nil
}

func getCloudBodyConversion(ctx context.Context, c *client.Client, asyncID string) (*BodyConversion, error) {
	data, _, err := c.Get(ctx, "/rest/api/contentbody/convert/async/"+url.PathEscape(asyncID), nil)
	if err != nil {
		return nil, err
	}
	return parseBodyConversion(data, "body conversion result")
}

func parseBodyConversion(data []byte, what string) (*BodyConversion, error) {
	var out BodyConversion
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parse %s: %w", what, err)
	}
	return &out, nil
}
