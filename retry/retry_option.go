package retry

import (
	"time"

	"github.com/googleapis/gax-go/v2"
)

// Config is wrapper of gax.Config.
// This can encapsulate gax setting in this module
type Config struct {
	gax.Backoff
	MaxRetry int
	count    int
}

// NewConfig gives new backoff setting
func NewConfig(initDuration time.Duration, maxDuration time.Duration, multiplier float64, maxRetry int) Config {
	return Config{
		Backoff: gax.Backoff{
			Initial:    initDuration,
			Max:        maxDuration,
			Multiplier: multiplier,
		},
		MaxRetry: maxRetry,
	}
}

// DefaultBackoff return default setting of backoff value as follow.
// Initial:    100 * time.Millisecond,
// Max:        30000 * time.Millisecond,
// Multiplier: 1.3,
func DefaultBackoff() Config {
	return NewConfig(100*time.Millisecond, 30000*time.Millisecond, 1.3, 10)
}
