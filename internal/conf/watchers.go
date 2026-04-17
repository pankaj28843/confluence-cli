package conf

import (
	"context"
	"net/url"

	"github.com/pankaj28843/confluence-cli/internal/client"
)

// GetContentWatchers returns raw bytes from /rest/api/content/{id}/notification/child-created
// equivalent endpoints — the exact shape varies between versions, so we pass through.
// In practice callers use /rest/api/user/watch/content/{id} for "am I watching?".
func GetContentWatchers(ctx context.Context, c *client.Client, contentID string) ([]byte, error) {
	// Server/DC exposes /rest/experimental/content/{id}/notification/ for watchers;
	// a stable endpoint is /rest/api/content/{id}/notification/child-created. We
	// surface notifications under a read-only passthrough:
	data, _, err := c.Get(ctx, "/rest/api/content/"+url.PathEscape(contentID)+"/notification/child-created", nil)
	return data, err
}

// GetSpaceWatchers returns users watching a space (Server/DC v11 endpoint).
func GetSpaceWatchers(ctx context.Context, c *client.Client, spaceKey string) ([]byte, error) {
	data, _, err := c.Get(ctx, "/rest/api/space/"+url.PathEscape(spaceKey)+"/watch", nil)
	return data, err
}

// GetContentRestrictions returns ACEs for a content id.
func GetContentRestrictions(ctx context.Context, c *client.Client, contentID string) ([]byte, error) {
	data, _, err := c.Get(ctx, "/rest/api/content/"+url.PathEscape(contentID)+"/restriction", url.Values{"expand": {"restrictions.user,restrictions.group"}})
	return data, err
}
