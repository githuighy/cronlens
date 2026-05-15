package scheduler

import (
	"math"
	"time"
)

// BackoffStrategy defines how retry delays are calculated.
type BackoffStrategy int

const (
	// BackoffFixed uses the same delay between every retry attempt.
	BackoffFixed BackoffStrategy = iota
	// BackoffLinear increases the delay linearly with each attempt.
	BackoffLinear
	// BackoffExponential doubles the delay with each attempt.
	BackoffExponential
)

// BackoffPolicy controls the delay between retry attempts.
type BackoffPolicy struct {
	// Strategy determines how the delay grows over attempts.
	Strategy BackoffStrategy
	// BaseDelay is the initial delay applied on the first retry.
	BaseDelay time.Duration
	// MaxDelay caps the computed delay so it never exceeds this value.
	MaxDelay time.Duration
	// Jitter adds a random fraction of BaseDelay when true (not implemented
	// here for determinism, but reserved for future use).
	Jitter bool
}

// DefaultBackoffPolicy returns a sensible exponential backoff starting at
// one second and capping at one minute.
func DefaultBackoffPolicy() BackoffPolicy {
	return BackoffPolicy{
		Strategy:  BackoffExponential,
		BaseDelay: time.Second,
		MaxDelay:  time.Minute,
	}
}

// Delay returns the wait duration for the given attempt number (1-based).
func (b BackoffPolicy) Delay(attempt int) time.Duration {
	if attempt < 1 {
		attempt = 1
	}

	var d time.Duration

	switch b.Strategy {
	case BackoffLinear:
		d = b.BaseDelay * time.Duration(attempt)
	case BackoffExponential:
		mult := math.Pow(2, float64(attempt-1))
		d = time.Duration(float64(b.BaseDelay) * mult)
	default: // BackoffFixed
		d = b.BaseDelay
	}

	if b.MaxDelay > 0 && d > b.MaxDelay {
		d = b.MaxDelay
	}

	return d
}
