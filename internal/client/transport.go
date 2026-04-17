package client

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

// debugRoundTripper logs method/URL/status/elapsed to stderr when debug=true.
// Authorization header is always redacted. Retries 429 and 5xx with exp backoff.
type debugRoundTripper struct {
	base  http.RoundTripper
	debug *bool
}

const (
	maxRetries     = 3
	backoffBase    = 500 * time.Millisecond
	backoffCeiling = 4 * time.Second
)

func wrapTransport(base http.RoundTripper, debug *bool) *debugRoundTripper {
	return &debugRoundTripper{base: base, debug: debug}
}

func (rt *debugRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	var lastErr error
	var resp *http.Response

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			wait := backoffFor(resp, attempt)
			timer := time.NewTimer(wait)
			select {
			case <-req.Context().Done():
				timer.Stop()
				return nil, req.Context().Err()
			case <-timer.C:
			}
		}

		if attempt > 0 && req.Body != nil {
			if req.GetBody == nil {
				return resp, lastErr
			}
			nb, err := req.GetBody()
			if err != nil {
				return resp, err
			}
			req.Body = nb
		}

		resp, lastErr = rt.base.RoundTrip(req)
		rt.logDebug(req, resp, lastErr, time.Since(start))

		if lastErr != nil {
			return resp, lastErr
		}
		if !shouldRetry(resp) {
			return resp, nil
		}
		if resp.Body != nil {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
		}
	}
	return resp, lastErr
}

func shouldRetry(resp *http.Response) bool {
	if resp == nil {
		return false
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		return true
	}
	if resp.StatusCode >= 500 && resp.StatusCode <= 599 {
		return true
	}
	return false
}

func backoffFor(resp *http.Response, attempt int) time.Duration {
	if resp != nil {
		if ra := resp.Header.Get("Retry-After"); ra != "" {
			if secs, err := strconv.Atoi(ra); err == nil && secs >= 0 {
				return time.Duration(secs) * time.Second
			}
		}
		// Cloud returns X-RateLimit-Reset (epoch) on throttling.
		if rs := resp.Header.Get("X-RateLimit-Reset"); rs != "" {
			if epoch, err := strconv.ParseInt(rs, 10, 64); err == nil {
				if d := time.Until(time.Unix(epoch, 0)); d > 0 && d < backoffCeiling {
					return d
				}
			}
		}
	}
	d := backoffBase * (1 << (attempt - 1))
	if d > backoffCeiling {
		d = backoffCeiling
	}
	return d
}

func (rt *debugRoundTripper) logDebug(req *http.Request, resp *http.Response, err error, elapsed time.Duration) {
	if rt.debug == nil || !*rt.debug {
		return
	}
	status := "-"
	if resp != nil {
		status = resp.Status
	}
	errStr := ""
	if err != nil {
		errStr = " err=" + err.Error()
	}
	fmt.Fprintf(os.Stderr, "[confluence] %s %s -> %s (%.1fms)%s\n",
		req.Method, req.URL.Redacted(), status, float64(elapsed.Microseconds())/1000, errStr)
}
