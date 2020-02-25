package retry

import (
	"context"
	"fmt"
	"time"
)

// Retry calls the supplied function f repeatedly according to the provided
// backoff parameters. It returns when one of the following occurs:
// When f's first return value is true, Retry immediately returns with f's second return value.
// When the provided context is done, Retry returns with an error that
// includes both ctx.Error() and the last error returned by f.
func Retry(ctx context.Context, bo Backoff, f func() (stop bool, err error)) error {
	return retry(ctx, bo, f, Sleep)
}

func retry(ctx context.Context, bo Backoff, f func() (stop bool, err error),
	sleep func(context.Context, time.Duration) error) error {
	var lastErr error
	for {
		stop, err := f()
		if stop {
			return err
		}
		// Remember the last error from f.
		if err != nil && err != context.Canceled && err != context.DeadlineExceeded {
			lastErr = err
		}
		p := bo.Pause()
		if cerr := sleep(ctx, p); cerr != nil {
			if lastErr != nil {
				return fmt.Errorf("retry failed with %v; last error: %v", cerr, lastErr)
			}
			return cerr
		}
	}
}

// Sleep is similar to time.Sleep, but it can be interrupted by ctx.Done() closing.
// If interrupted, Sleep returns ctx.Err().
func Sleep(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	select {
	case <-ctx.Done():
		t.Stop()
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

// Is retryable is an interface to check if given error is retriable or not
type RetriableChecker interface {
	IsRetryableError(err error) bool
}

// Is retryable is an adopter of RetriableChecker
type IsRetryable func(err error) bool
// IsRetryableError calls r(err)
func (r IsRetryable) IsRetryableError(err error) bool {
	return r(err)
}

// RunWithRetry calls the function until it returns nil or a non-retryable error, or
// the context is done.
// See the similar function in ../storage/invoke.go. The main difference is the
// reason for retrying.
func RunWithRetry(ctx context.Context, backoff Backoff, call func() error, checker IsRetryable) error {
	return Retry(ctx, backoff, func() (stop bool, err error) {
		err = call()
		if err == nil {
			return true, nil
		}
		return !checker(err), err
	})
}
