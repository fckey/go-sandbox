package retry

import (
	"context"
	"fmt"
	"time"
)

type retryIter func() (stop bool, err error)

// Retry calls the supplied function retryIter repeatedly according to the provided
// backoff parameters in cfg. It returns when one of the following occurs:
// When retryIter's first return value is true, Retry immediately returns with retryIter's value in err.
// When the provided context is done, Retry returns with an error that
// includes both ctx.Error() and the last error returned by retryIter.
func Retry(ctx context.Context, cfg Config, f retryIter) error {
	return retry(ctx, cfg, f, sleep)
}

func retry(ctx context.Context, cfg Config, f retryIter, s Sleep) error {
	var lastErr error
	for {
		cfg.count++
		stop, err := f()
		if stop {
			return err
		}
		// Remember the last error from f.
		if err != nil && err != context.Canceled && err != context.DeadlineExceeded {
			lastErr = err
		}
		p := cfg.Backoff.Pause()
		if cerr := s(ctx, p); cerr != nil {
			if lastErr != nil {
				return fmt.Errorf("retry failed with %v; last error: %v", cerr, lastErr)
			}
			return cerr
		}
		if cfg.count >= cfg.MaxRetry {
			if lastErr != nil {
				return fmt.Errorf("maximum traial exceeded; last error: %v", lastErr)
			}
			return fmt.Errorf("operation was not succeded withnin %d trial", cfg.MaxRetry)
		}
	}
}

// Sleep is similar to time.Sleep, but it can be interrupted by ctx.Done() closing.
// Error is expected to be returned when it's interrupted.
type Sleep func(context.Context, time.Duration) error

// sleep is implementation of retry.Sleep
func sleep(ctx context.Context, d time.Duration) error {
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
func RunWithRetry(ctx context.Context, cfg Config, call func() error, checker IsRetryable) error {
	return Retry(ctx, cfg, func() (stop bool, err error) {
		err = call()
		if err == nil {
			return true, nil
		}
		return !checker(err), err
	})
}
