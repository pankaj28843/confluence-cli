package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	DefaultHTTPTimeout = 20 * time.Second
	maxErrorBodyLen    = 400
)

// Client speaks Confluence REST — Server/DC (Bearer) or Cloud (Basic email:token).
type Client struct {
	BaseURL      string
	Flavor       Flavor
	PAT          string
	Email        string
	APIToken     string
	DefaultSpace string
	HTTPClient   *http.Client
	UserAgent    string
	Debug        bool
}

type Option func(*Client)

func WithHTTPClient(hc *http.Client) Option { return func(c *Client) { c.HTTPClient = hc } }
func WithDebug(d bool) Option               { return func(c *Client) { c.Debug = d } }
func WithUserAgent(ua string) Option        { return func(c *Client) { c.UserAgent = ua } }

// New constructs a Client from Config.
func New(cfg Config, opts ...Option) (*Client, error) {
	tr := &http.Transport{
		TLSClientConfig:       &tls.Config{MinVersion: tls.VersionTLS12},
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 15 * time.Second,
		IdleConnTimeout:       30 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	c := &Client{
		BaseURL:      cfg.BaseURL,
		Flavor:       cfg.Flavor,
		PAT:          cfg.PAT,
		Email:        cfg.Email,
		APIToken:     cfg.APIToken,
		DefaultSpace: cfg.DefaultSpace,
		UserAgent:    DefaultUserAgent,
		Debug:        cfg.Debug,
	}
	c.HTTPClient = &http.Client{
		Timeout:   DefaultHTTPTimeout,
		Transport: wrapTransport(tr, &c.Debug),
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.HTTPClient != nil && c.HTTPClient.Transport != nil {
		if _, ok := c.HTTPClient.Transport.(*debugRoundTripper); !ok {
			c.HTTPClient.Transport = wrapTransport(c.HTTPClient.Transport, &c.Debug)
		}
	}
	return c, nil
}

// authHeader returns the correct Authorization value for this flavor.
func (c *Client) authHeader() string {
	if c.Flavor == FlavorCloud {
		return "Basic " + base64.StdEncoding.EncodeToString([]byte(c.Email+":"+c.APIToken))
	}
	return "Bearer " + c.PAT
}

// BuildURL joins the base URL with the supplied path. Path should start with '/'.
// Callers use /rest/api/... or /wiki/api/v2/... directly — the client does not
// inject prefixes.
func (c *Client) BuildURL(path string, params url.Values) string {
	u := strings.TrimRight(c.BaseURL, "/") + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}
	return u
}

// Get issues a GET request.
func (c *Client) Get(ctx context.Context, path string, params url.Values) ([]byte, *http.Header, error) {
	return c.do(ctx, http.MethodGet, path, params, nil, "")
}

// Post issues a POST with a JSON body.
func (c *Client) Post(ctx context.Context, path string, params url.Values, body any) ([]byte, *http.Header, error) {
	return c.do(ctx, http.MethodPost, path, params, body, "application/json")
}

// Put issues a PUT with a JSON body.
func (c *Client) Put(ctx context.Context, path string, params url.Values, body any) ([]byte, *http.Header, error) {
	return c.do(ctx, http.MethodPut, path, params, body, "application/json")
}

// Delete issues a DELETE request.
func (c *Client) Delete(ctx context.Context, path string, params url.Values) ([]byte, *http.Header, error) {
	return c.do(ctx, http.MethodDelete, path, params, nil, "")
}

// PostRawReader streams a caller-supplied body with a caller-specified
// Content-Type. Use it for multipart uploads where buffering large files would
// be wasteful.
func (c *Client) PostRawReader(ctx context.Context, path string, params url.Values, body io.Reader, contentType string, extra map[string]string) ([]byte, *http.Header, error) {
	return c.doRawReader(ctx, http.MethodPost, path, params, body, contentType, extra)
}

// PostRaw sends raw bytes with a caller-specified Content-Type (for multipart).
// Extra headers are merged in.
func (c *Client) PostRaw(ctx context.Context, path string, params url.Values, body []byte, contentType string, extra map[string]string) ([]byte, *http.Header, error) {
	return c.doRaw(ctx, http.MethodPost, path, params, body, contentType, extra)
}

func (c *Client) do(ctx context.Context, method, path string, params url.Values, body any, contentType string) ([]byte, *http.Header, error) {
	var raw []byte
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, nil, fmt.Errorf("marshal body: %w", err)
		}
		raw = b
	}
	return c.doRaw(ctx, method, path, params, raw, contentType, nil)
}

func (c *Client) doRaw(ctx context.Context, method, path string, params url.Values, body []byte, contentType string, extra map[string]string) ([]byte, *http.Header, error) {
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}
	return c.doRawReader(ctx, method, path, params, reader, contentType, extra)
}

func (c *Client) doRawReader(ctx context.Context, method, path string, params url.Values, body io.Reader, contentType string, extra map[string]string) ([]byte, *http.Header, error) {
	u := c.BuildURL(path, params)

	req, err := http.NewRequestWithContext(ctx, method, u, body)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Authorization", c.authHeader())
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	for k, v := range extra {
		req.Header.Set(k, v)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, c.networkError(err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		hdrs := resp.Header.Clone()
		return nil, &hdrs, &APIError{
			Status:      resp.StatusCode,
			Message:     fmt.Sprintf("authentication failed (HTTP %d)", resp.StatusCode),
			Body:        truncate(string(data), maxErrorBodyLen),
			UserFixable: true,
		}
	}
	if resp.StatusCode >= 400 {
		hdrs := resp.Header.Clone()
		return nil, &hdrs, &APIError{
			Status:  resp.StatusCode,
			Message: fmt.Sprintf("HTTP %d", resp.StatusCode),
			Body:    truncate(string(data), maxErrorBodyLen),
		}
	}
	hdrs := resp.Header.Clone()
	return data, &hdrs, nil
}

// APIError wraps non-2xx responses with a flag for user-fixable config issues.
type APIError struct {
	Status      int
	Message     string
	Body        string
	UserFixable bool
}

func (e *APIError) Error() string {
	if e.Body == "" {
		return e.Message
	}
	return fmt.Sprintf("%s: %s", e.Message, e.Body)
}

// IsUserFixable reports whether the error maps to exit-code 2.
func IsUserFixable(err error) bool {
	var ae *APIError
	if errors.As(err, &ae) {
		return ae.UserFixable
	}
	return errors.Is(err, ErrMissingURL) ||
		errors.Is(err, ErrMissingServerAuth) ||
		errors.Is(err, ErrMissingCloudAuth)
}

func (c *Client) networkError(err error) error {
	msg := err.Error()
	hint := ""
	switch {
	case strings.Contains(msg, "no such host"):
		hint = "Check CONFLUENCE_URL and DNS."
	case strings.Contains(msg, "certificate"):
		hint = "Set SSL_CERT_FILE to a PEM bundle if your host uses a private CA."
	case strings.Contains(msg, "connection refused"), strings.Contains(msg, "i/o timeout"), strings.Contains(msg, "deadline exceeded"):
		hint = "Host not reachable (network / VPN?)."
	}
	if hint == "" {
		return err
	}
	return fmt.Errorf("%w (hint: %s)", err, errors.New(hint))
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
