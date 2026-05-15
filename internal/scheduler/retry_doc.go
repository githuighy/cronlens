// Package scheduler provides scheduling utilities for cron expressions,
// including next-run prediction, last-run tracking, duration formatting,
// run history, tagging, and retry policies.
//
// # Retry
//
// The retry sub-feature allows callers to execute an arbitrary function
// according to a [RetryPolicy], automatically re-running it on failure
// up to a configurable maximum number of attempts with an optional
// exponential back-off delay between attempts.
//
// Basic usage:
//
//	policy := scheduler.RetryPolicy{
//		MaxAttempts:   5,
//		Delay:         2 * time.Second,
//		BackoffFactor: 2.0,
//	}
//
//	results, err := scheduler.RunWithRetry(policy, func() error {
//		return doWork()
//	})
//	if errors.Is(err, scheduler.ErrMaxAttemptsReached) {
//		log.Printf("all %d attempts failed", len(results))
//	}
//
// Use [DefaultRetryPolicy] for a sensible out-of-the-box configuration
// (3 attempts, 5 s initial delay, 2× back-off factor).
package scheduler
