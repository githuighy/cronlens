// Package scheduler provides scheduling primitives for cron-based job execution.
//
// # TimeoutStore
//
// TimeoutStore manages per-job execution timeouts. A global default can be
// set at construction time, and individual jobs can override that default
// with their own duration.
//
// Usage:
//
//	store := scheduler.NewTimeoutStore(30 * time.Second)
//
//	// Override timeout for a specific job.
//	store.Set("slow-report", 2 * time.Minute)
//
//	// Wrap a job to enforce its timeout.
//	protected := store.WithTimeout("slow-report", myJob)
//
// When the context passed to the wrapped job is cancelled before the
// timeout fires, the original cancellation takes precedence — the
// TimeoutStore simply adds an upper bound.
package scheduler
