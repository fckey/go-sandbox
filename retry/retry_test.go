package retry

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	ctx := context.Background()
	// Without a context deadline, retry will run until the function
	// says not to retry any more.
	n := 0
	endRetry := errors.New("end retry")
	err := retry(ctx, DefaultBackoff(),
		func() (bool, error) {
			n++
			if n < 5 {
				return false, nil
			}
			return true, endRetry
		},
		func(context.Context, time.Duration) error { return nil })
	if got, want := err, endRetry; got != want {
		t.Errorf("got %v, want %v", err, endRetry)
	}
	if n != 5 {
		t.Errorf("n: got %d, want %d", n, 10)
	}

	// If the context has a deadline, sleep will return an error
	// and end the function.
	n = 0
	err = retry(ctx, DefaultBackoff(),
		func() (bool, error) { return false, nil },
		func(context.Context, time.Duration) error {
			n++
			if n < 10 {
				return nil
			}
			return context.DeadlineExceeded
		})
	if err == nil {
		t.Errorf("got nil, want error: %v", err)
	}
}

func TestRunWithRetry(t *testing.T) {
	ctx := context.Background()
	// Without a context deadline, retry will run until the function
	// says not to retry any more.
	n := 0
	retriableMsg := "retriable"
	middleRetry := errors.New(retriableMsg)
	endRetry := errors.New("end retry")

	retriableFunc := func() error {
		n++
		if n < 5 {
			return middleRetry
		}
		return endRetry
	}
	checkerFunc := func(err error) bool {
		if strings.Compare(err.Error(), retriableMsg) == 0 {
			return true
		}
		return false
	}
	//backoff := NewConfig(1*time.Second, 32*time.Second, 2, 10)
	err := RunWithRetry(ctx, DefaultBackoff(), retriableFunc, checkerFunc)
	if n != 5 {
		t.Errorf("n: got %d, want %d", n, 10)
	}

	if err == nil {
		t.Error("got nil, want error")
	}
}
