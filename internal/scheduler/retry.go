package scheduler

import (
	"errors"
	"time"
)

// RetryPolicy defines how missed or failed runs should be retried.
type RetryPolicy struct {
	// MaxAttempts is the maximum number of retry attempts (0 = no retries).
	MaxAttempts int
	// Delay is the duration to wait between retry attempts.
	Delay time.Duration
	// BackoffFactor multiplies Delay on each successive attempt (1.0 = no backoff).
	BackoffFactor float64
}

// DefaultRetryPolicy returns a sensible default retry policy.
func DefaultRetryPolicy() RetryPolicy {
	return RetryPolicy{
		MaxAttempts:   3,
		Delay:         5 * time.Second,
		BackoffFactor: 2.0,
	}
}

// ErrMaxAttemptsReached is returned when all retry attempts have been exhausted.
var ErrMaxAttemptsReached = errors.New("max retry attempts reached")

// RetryResult records the outcome of a single attempt.
type RetryResult struct {
	Attempt   int
	At        time.Time
	Succeeded bool
	Err       error
}

// RunWithRetry executes fn according to the given RetryPolicy, returning
// all attempt results and a final error if every attempt failed.
func RunWithRetry(policy RetryPolicy, fn func() error) ([]RetryResult, error) {
	max := policy.MaxAttempts
	if max < 1 {
		max = 1
	}

	results := make([]RetryResult, 0, max)
	delay := policy.Delay

	for attempt := 1; attempt <= max; attempt++ {
		err := fn()
		result := RetryResult{
			Attempt:   attempt,
			At:        time.Now(),
			Succeeded: err == nil,
			Err:       err,
		}
		results = append(results, result)

		if err == nil {
			return results, nil
		}

		if attempt < max {
			time.Sleep(delay)
			if policy.BackoffFactor > 0 {
				delay = time.Duration(float64(delay) * policy.BackoffFactor)
			}
		}
	}

	return results, ErrMaxAttemptsReached
}
