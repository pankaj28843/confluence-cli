package conf

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

func isCloud(c *client.Client) bool {
	return c != nil && c.Flavor == client.FlavorCloud
}

func nextPageRequest(c *client.Client, next string) (string, url.Values, error) {
	u, err := url.Parse(next)
	if err != nil {
		return "", nil, fmt.Errorf("parse next link %q: %w", next, err)
	}

	path := u.Path
	if path == "" {
		path = next
	}
	if c != nil && c.BaseURL != "" {
		base, err := url.Parse(c.BaseURL)
		if err != nil {
			return "", nil, fmt.Errorf("parse base url %q: %w", c.BaseURL, err)
		}
		basePath := strings.TrimRight(base.Path, "/")
		if basePath != "" && strings.HasPrefix(path, basePath+"/") {
			path = strings.TrimPrefix(path, basePath)
		}
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return path, u.Query(), nil
}

func nextPageURL(headers *http.Header, bodyNext string) string {
	if headers != nil {
		for _, raw := range headers.Values("Link") {
			for _, part := range strings.Split(raw, ",") {
				part = strings.TrimSpace(part)
				if !strings.Contains(part, `rel="next"`) {
					continue
				}
				semi := strings.Index(part, ";")
				if semi < 0 {
					continue
				}
				target := strings.TrimSpace(part[:semi])
				if strings.HasPrefix(target, "<") && strings.HasSuffix(target, ">") {
					return strings.TrimSuffix(strings.TrimPrefix(target, "<"), ">")
				}
			}
		}
	}
	return bodyNext
}
