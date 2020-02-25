package retry

import (
	"github.com/googleapis/gax-go/v2"
	"time"
)

// CallSettings allow fine-grained control over how calls are made.
type CallSettings struct {
	// Retry returns a Retryer to be used to control retry logic of a method call.
	// If Retry is nil or the returned Retryer is nil, the call will not be retried.
	Retry func() Retryer
}

// Retryer is used by Invoke to determine retry behavior.
type Retryer interface {
	// Retry reports whether a request should be retriedand how long to pause before retrying
	// if the previous attempt returned with err. Invoke never calls Retry with nil error.
	Retry(err error) (pause time.Duration, shouldRetry bool)
}

// CallOption is an option used by Invoke to control behaviors of RPC calls.
// CallOption works by modifying relevant fields of CallSettings.
type CallOption interface {
	// Resolve applies the option by modifying cs.
	Resolve(cs *CallSettings)
}

type retryerOption func() Retryer

func (o retryerOption) Resolve(s *CallSettings) {
	s.Retry = o
}

// Backoff is wrapper of gax.Backoff.
// This can encapsulate gax setting in this module
type Backoff struct {
	gax.Backoff
}

// NewBackoff gives new backoff setting
func NewBackoff(Initial time.Duration, Max time.Duration, Multiplier float64) Backoff {
	return Backoff{
		gax.Backoff{
			Initial:    100 * time.Millisecond,
			Max:        60000 * time.Millisecond,
			Multiplier: 1.3,
		},
	}
}

// DefaultBackoff return default setting of backoff value as follow.
// Initial:    100 * time.Millisecond,
// Max:        30000 * time.Millisecond,
// Multiplier: 1.3,
func DefaultBackoff() Backoff {
	return NewBackoff(100*time.Millisecond, 60000*time.Millisecond, 1.3)
}

func OnStatusCode(statuses [] int, bo gax.Backoff) Retryer {
	return &boRetryer{
		backoff:      bo,
		httpStatuses: statuses,
	}
}

type boRetryer struct {
	backoff      gax.Backoff
	httpStatuses [] int
}

func (r *boRetryer) Retry(err error) (time.Duration, bool) {
	//
	//if !ok {
	//	return 0, false
	//}
	//
	//for _, rc := range r.httpStatuses {
	//	if c == rc {
	//		return r.backoff.Pause(), true
	//	}
	//}
	return 0, false
}
