package conf

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

// TemplateListOptions controls Cloud template listing.
type TemplateListOptions struct {
	SpaceKey string
	Expand   []string
	Limit    int
}

// ContentTemplate is a Confluence content or blueprint template.
type ContentTemplate struct {
	TemplateID           string `json:"templateId,omitempty"`
	Name                 string `json:"name,omitempty"`
	Description          string `json:"description,omitempty"`
	TemplateType         string `json:"templateType,omitempty"`
	EditorVersion        string `json:"editorVersion,omitempty"`
	ReferencingBlueprint string `json:"referencingBlueprint,omitempty"`
	OriginalTemplate     struct {
		PluginKey string `json:"pluginKey,omitempty"`
		ModuleKey string `json:"moduleKey,omitempty"`
	} `json:"originalTemplate,omitempty"`
	Space  map[string]any `json:"space,omitempty"`
	Labels []Label        `json:"labels,omitempty"`
	Body   TemplateBody   `json:"body,omitempty"`
	Links  map[string]any `json:"_links,omitempty"`
}

// TemplateBody contains the documented content-template body representations.
type TemplateBody struct {
	View                TemplateBodyRepresentation `json:"view,omitempty"`
	ExportView          TemplateBodyRepresentation `json:"export_view,omitempty"`
	StyledView          TemplateBodyRepresentation `json:"styled_view,omitempty"`
	Storage             TemplateBodyRepresentation `json:"storage,omitempty"`
	Editor              TemplateBodyRepresentation `json:"editor,omitempty"`
	Editor2             TemplateBodyRepresentation `json:"editor2,omitempty"`
	Wiki                TemplateBodyRepresentation `json:"wiki,omitempty"`
	AtlasDocFormat      TemplateBodyRepresentation `json:"atlas_doc_format,omitempty"`
	AnonymousExportView TemplateBodyRepresentation `json:"anonymous_export_view,omitempty"`
}

// TemplateBodyRepresentation is one content-template body value.
type TemplateBodyRepresentation struct {
	Value          string         `json:"value,omitempty"`
	Representation string         `json:"representation,omitempty"`
	Links          map[string]any `json:"_links,omitempty"`
}

// ListContentTemplates lists Cloud content templates.
func ListContentTemplates(ctx context.Context, c *client.Client, opts TemplateListOptions) ([]ContentTemplate, error) {
	if !isCloud(c) {
		return nil, fmt.Errorf("content templates are only supported on Confluence Cloud")
	}
	return listTemplatesCloudV1(ctx, c, "/rest/api/template/page", opts, "content templates")
}

// ListBlueprintTemplates lists Cloud blueprint templates.
func ListBlueprintTemplates(ctx context.Context, c *client.Client, opts TemplateListOptions) ([]ContentTemplate, error) {
	if !isCloud(c) {
		return nil, fmt.Errorf("blueprint templates are only supported on Confluence Cloud")
	}
	return listTemplatesCloudV1(ctx, c, "/rest/api/template/blueprint", opts, "blueprint templates")
}

// GetContentTemplate returns one Cloud content template.
func GetContentTemplate(ctx context.Context, c *client.Client, id string, expand []string) (*ContentTemplate, error) {
	if id == "" {
		return nil, fmt.Errorf("GetContentTemplate: id is required")
	}
	if !isCloud(c) {
		return nil, fmt.Errorf("content templates are only supported on Confluence Cloud")
	}
	params := url.Values{}
	addValues(params, "expand", expand)
	data, _, err := c.Get(ctx, "/rest/api/template/"+url.PathEscape(id), params)
	if err != nil {
		return nil, err
	}
	return parseTemplate(data, "content template")
}

func listTemplatesCloudV1(ctx context.Context, c *client.Client, path string, opts TemplateListOptions, what string) ([]ContentTemplate, error) {
	limit := clampLimit(opts.Limit)
	out := make([]ContentTemplate, 0, limit)
	params := url.Values{
		"start": {"0"},
		"limit": {strconv.Itoa(limit)},
	}
	if opts.SpaceKey != "" {
		params.Set("spaceKey", opts.SpaceKey)
	}
	addValues(params, "expand", opts.Expand)

	for len(out) < limit {
		data, _, err := c.Get(ctx, path, params)
		if err != nil {
			return nil, err
		}
		var page struct {
			Results []ContentTemplate `json:"results"`
			Size    int               `json:"size"`
			Links   struct {
				Next string `json:"next,omitempty"`
			} `json:"_links"`
		}
		if err := json.Unmarshal(data, &page); err != nil {
			return nil, fmt.Errorf("parse %s: %w", what, err)
		}
		for _, template := range page.Results {
			out = append(out, template)
			if len(out) == limit {
				break
			}
		}
		if page.Links.Next == "" || len(page.Results) == 0 || len(out) == limit {
			break
		}
		nextPath, nextParams, err := nextPageRequest(c, page.Links.Next)
		if err != nil {
			return nil, err
		}
		path, params = nextPath, nextParams
	}
	return out, nil
}

func parseTemplate(data []byte, what string) (*ContentTemplate, error) {
	var out ContentTemplate
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parse %s: %w", what, err)
	}
	return &out, nil
}
