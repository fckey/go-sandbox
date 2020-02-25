package retry

import (
	"context"
	"net/http"
	"time"
)

// RetryableHTTPRequest is http call which expect to be retried.
// This function should become a closure having *http.Client and *http.Request
type RetryableHTTPRequest func() (resp *http.Response, err error)

// Is retryable is an interface to check if given error is retriable or not
type RetriableHTTPResponseChecker interface {
	IsRetryableStatus(resp *http.Response, err error) bool
}

// Is retryable is an adopter of RetriableChecker
type IsRetryableHTTPResponse func(resp *http.Response, err error) bool

// IsRetryableError calls r(err)
func (i IsRetryableHTTPResponse) IsRetryableStatus(resp *http.Response, err error) bool {
	return i(resp, err)
}


func RunWithHTTPRetry(ctx context.Context, backoff Backoff,
	call RetryableHTTPRequest, checker IsRetryableHTTPResponse) (resp *http.Response, err error) {
	return resp, Retry(ctx, backoff, func() (stop bool, err error) {
		resp, err = call()
		if err == nil {
			return true, nil
		}
		return !checker(resp, err), err
	})
}

func WithRtriableHTTPResponse(statuses []int, err error) IsRetryableHTTPResponse {
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

		// Never retry client error response code like StatusBadRequest or StatusUnauthorized
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
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

// APICall is a user defined call stub.
type APICall func(context.Context, CallSettings) (*http.Response, error)

type sleeper func(ctx context.Context, d time.Duration) error

// invoke implements Invoke, taking an additional sleeper argument for testing.
func invoke(ctx context.Context, call APICall, settings CallSettings, sp sleeper) (resp *http.Response, err error) {
	var retryer Retryer
	for {
		resp, err := call(ctx, settings)
		if err == nil {
			return resp, nil
		}
		if settings.Retry == nil {
			return nil, err
		}
		// Never retry client error response code like StatusBadRequest or StatusUnauthorized
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			return nil, err
		}
		if retryer == nil {
			if r := settings.Retry(); r != nil {
				retryer = r
			} else {
				return nil, err
			}
		}
		if d, ok := retryer.Retry(err); !ok {
			return resp, err
		} else if err = sp(ctx, d); err != nil {
			return nil, err
		}
	}
}
