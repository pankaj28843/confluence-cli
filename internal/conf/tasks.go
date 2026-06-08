package conf

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

// Task is a Confluence Cloud task.
type Task struct {
	ID         string `json:"id,omitempty"`
	LocalID    string `json:"localId,omitempty"`
	SpaceID    string `json:"spaceId,omitempty"`
	PageID     string `json:"pageId,omitempty"`
	BlogPostID string `json:"blogPostId,omitempty"`
	Status     string `json:"status,omitempty"` // complete | incomplete
	Body       struct {
		Storage struct {
			Representation string `json:"representation,omitempty"`
			Value          string `json:"value,omitempty"`
		} `json:"storage,omitempty"`
		AtlasDocFormat struct {
			Representation string `json:"representation,omitempty"`
			Value          string `json:"value,omitempty"`
		} `json:"atlas_doc_format,omitempty"`
	} `json:"body,omitempty"`
	CreatedBy   string `json:"createdBy,omitempty"`
	AssignedTo  string `json:"assignedTo,omitempty"`
	CompletedBy string `json:"completedBy,omitempty"`
	CreatedAt   string `json:"createdAt,omitempty"`
	UpdatedAt   string `json:"updatedAt,omitempty"`
	DueAt       string `json:"dueAt,omitempty"`
	CompletedAt string `json:"completedAt,omitempty"`
}

// TaskFilter narrows Cloud task listing.
type TaskFilter struct {
	Status            string
	TaskIDs           []string
	SpaceIDs          []string
	PageID            string
	BlogPostID        string
	CreatedBy         []string
	AssignedTo        []string
	CompletedBy       []string
	IncludeBlankTasks bool
	BodyFormat        string
	Limit             int
}

// ListTasks lists Confluence Cloud content tasks.
func ListTasks(ctx context.Context, c *client.Client, f TaskFilter) ([]Task, error) {
	if !isCloud(c) {
		return nil, fmt.Errorf("content tasks are only supported on Confluence Cloud")
	}
	return listTasksCloudV2(ctx, c, f)
}

// GetTask returns one Confluence Cloud content task by id.
func GetTask(ctx context.Context, c *client.Client, id, bodyFormat string) (*Task, error) {
	if id == "" {
		return nil, fmt.Errorf("GetTask: id is required")
	}
	if !isCloud(c) {
		return nil, fmt.Errorf("content tasks are only supported on Confluence Cloud")
	}
	params := url.Values{}
	if bodyFormat != "" {
		params.Set("body-format", bodyFormat)
	}
	data, _, err := c.Get(ctx, "/api/v2/tasks/"+url.PathEscape(id), params)
	if err != nil {
		return nil, err
	}
	return parseTask(data, "task")
}

// UpdateTaskStatus updates one Confluence Cloud task status.
func UpdateTaskStatus(ctx context.Context, c *client.Client, id, status string) (*Task, error) {
	if id == "" {
		return nil, fmt.Errorf("UpdateTaskStatus: id is required")
	}
	if !isCloud(c) {
		return nil, fmt.Errorf("content tasks are only supported on Confluence Cloud")
	}
	if status != "complete" && status != "incomplete" {
		return nil, fmt.Errorf("UpdateTaskStatus: status must be complete or incomplete")
	}
	body := map[string]any{
		"id":     id,
		"status": status,
	}
	data, _, err := c.Put(ctx, "/api/v2/tasks/"+url.PathEscape(id), nil, body)
	if err != nil {
		return nil, err
	}
	return parseTask(data, "task")
}

func listTasksCloudV2(ctx context.Context, c *client.Client, f TaskFilter) ([]Task, error) {
	limit := normalizeTaskLimit(f.Limit)
	out := make([]Task, 0, limit)
	path := "/api/v2/tasks"
	params := taskQuery(f, limit)

	for len(out) < limit {
		data, hdrs, err := c.Get(ctx, path, params)
		if err != nil {
			return nil, err
		}
		var page struct {
			Results []Task `json:"results"`
			Links   struct {
				Next string `json:"next,omitempty"`
			} `json:"_links"`
		}
		if err := json.Unmarshal(data, &page); err != nil {
			return nil, fmt.Errorf("parse tasks: %w", err)
		}
		for _, task := range page.Results {
			out = append(out, task)
			if len(out) == limit {
				break
			}
		}

		next := nextPageURL(hdrs, page.Links.Next)
		if next == "" || len(page.Results) == 0 || len(out) == limit {
			break
		}
		path, params, err = nextPageRequest(c, next)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func taskQuery(f TaskFilter, limit int) url.Values {
	params := url.Values{"limit": {strconv.Itoa(limit)}}
	if f.Status != "" {
		params.Set("status", f.Status)
	}
	if f.PageID != "" {
		params.Set("page-id", f.PageID)
	}
	if f.BlogPostID != "" {
		params.Set("blogpost-id", f.BlogPostID)
	}
	if f.BodyFormat != "" {
		params.Set("body-format", f.BodyFormat)
	}
	if f.IncludeBlankTasks {
		params.Set("include-blank-tasks", "true")
	}
	addValues(params, "task-id", f.TaskIDs)
	addValues(params, "space-id", f.SpaceIDs)
	addValues(params, "created-by", f.CreatedBy)
	addValues(params, "assigned-to", f.AssignedTo)
	addValues(params, "completed-by", f.CompletedBy)
	return params
}

func addValues(params url.Values, key string, values []string) {
	for _, value := range values {
		if value != "" {
			params.Add(key, value)
		}
	}
}

func parseTask(data []byte, what string) (*Task, error) {
	var out Task
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parse %s: %w", what, err)
	}
	return &out, nil
}

// LongTask is a Server/Data Center long-running task.
type LongTask struct {
	ID   string `json:"id,omitempty"`
	Name struct {
		Key         string `json:"key,omitempty"`
		Args        []any  `json:"args,omitempty"`
		Translation string `json:"translation,omitempty"`
	} `json:"name,omitempty"`
	ElapsedTime        int64 `json:"elapsedTime,omitempty"`
	PercentageComplete int   `json:"percentageComplete,omitempty"`
	Successful         bool  `json:"successful,omitempty"`
	Messages           []struct {
		Key         string `json:"key,omitempty"`
		Args        []any  `json:"args,omitempty"`
		Translation string `json:"translation,omitempty"`
	} `json:"messages,omitempty"`
}

// LongTaskFilter narrows Server/Data Center long task listing.
type LongTaskFilter struct {
	Expand string
	Limit  int
}

// ListLongTasks lists Server/Data Center long-running tasks.
func ListLongTasks(ctx context.Context, c *client.Client, f LongTaskFilter) ([]LongTask, error) {
	if isCloud(c) {
		return nil, fmt.Errorf("long tasks are only supported on Confluence Server/Data Center")
	}
	limit := normalizeTaskLimit(f.Limit)
	out := make([]LongTask, 0, limit)
	start := 0

	for len(out) < limit {
		remaining := limit - len(out)
		if remaining > 200 {
			remaining = 200
		}
		params := url.Values{
			"limit": {strconv.Itoa(remaining)},
			"start": {strconv.Itoa(start)},
		}
		if f.Expand != "" {
			params.Set("expand", f.Expand)
		}
		data, _, err := c.Get(ctx, "/rest/api/longtask", params)
		if err != nil {
			return nil, err
		}
		var page struct {
			Results []LongTask `json:"results"`
			Size    int        `json:"size"`
			Links   struct {
				Next string `json:"next,omitempty"`
			} `json:"_links"`
		}
		if err := json.Unmarshal(data, &page); err != nil {
			return nil, fmt.Errorf("parse long tasks: %w", err)
		}
		for _, task := range page.Results {
			out = append(out, task)
			if len(out) == limit {
				break
			}
		}
		if page.Links.Next == "" || len(page.Results) == 0 || len(out) == limit {
			break
		}
		if page.Size > 0 {
			start += page.Size
		} else {
			start += len(page.Results)
		}
	}
	return out, nil
}

// GetLongTask returns one Server/Data Center long-running task by id.
func GetLongTask(ctx context.Context, c *client.Client, id, expand string) (*LongTask, error) {
	if id == "" {
		return nil, fmt.Errorf("GetLongTask: id is required")
	}
	if isCloud(c) {
		return nil, fmt.Errorf("long tasks are only supported on Confluence Server/Data Center")
	}
	params := url.Values{}
	if expand != "" {
		params.Set("expand", expand)
	}
	data, _, err := c.Get(ctx, "/rest/api/longtask/"+url.PathEscape(id), params)
	if err != nil {
		return nil, err
	}
	var out LongTask
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parse long task: %w", err)
	}
	return &out, nil
}

func normalizeTaskLimit(limit int) int {
	if limit <= 0 {
		return 25
	}
	if limit > 200 {
		return 200
	}
	return limit
}
