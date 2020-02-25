package retry

import (
	"context"
	"net/http"
)

// RetryableHTTPRequest is http call which expect to be retried.
// This function should become a closure having *http.Client and *http.Request
type RetryableHTTPRequest func() (resp *http.Response, err error)

// Is retryable is an interface to check if given error is retriable or not
type RetriableHTTPResponseChecker interface {
	IsRetryableStatus(resp *http.Response, err error) bool
}

// Is retryable is an adopter of RetriableChecker
type IsHTTPRequestRetryable func(resp *http.Response, err error) bool

// IsRetryableError calls r(err)
func (i IsHTTPRequestRetryable) IsRetryableStatus(resp *http.Response, err error) bool {
	return i(resp, err)
}

func RunWithHTTPRetry(ctx context.Context, config Config,
	call RetryableHTTPRequest, checker IsHTTPRequestRetryable) (resp *http.Response, err error) {

	return resp, Retry(ctx, config, func() (stop bool, err error) {
		resp, err = call()
		if err == nil {
			return true, nil
		}
		return !checker(resp, err), err
	})
}

// DoWithRetry executes http.Client.Do(http.Request) of given http.Client and http.Request as retryable manner
// If simple Do is required, this function should be used.
// statuses is expected as list of StatusCode in http
func DoWithRetry(ctx context.Context, cfg Config,
	c *http.Client, r *http.Request, statuses ...int) (resp *http.Response, err error) {
	checker := WithRtriableHTTPResponse(statuses...)
	return resp, Retry(ctx, cfg, func() (stop bool, err error) {
		resp, err = c.Do(r)
		if err == nil && resp.StatusCode < 400 {
			return true, nil
		}
		return !checker(resp, err), err
	})
}

// WithRtriableHTTPResponse judges if response is retriable or not
// Reatriable response can be set by statuses
func WithRtriableHTTPResponse(statuses ...int) IsHTTPRequestRetryable {
	retriableStatuses := []int{
		http.StatusRequestTimeout,
		http.StatusGatewayTimeout,
	}
	retriableStatuses = append(retriableStatuses, statuses...)
	return func(resp *http.Response, err error) bool {
		// Err for HTTP request is not retryable by default
		if err != nil {
			return false
		}

		// Never retry StatusUnauthorized
		if resp.StatusCode == 401 {
			return false
		}

		for _, status := range retriableStatuses {
			if status == resp.StatusCode {
				return true
			}
		}
		return false
	}
}
