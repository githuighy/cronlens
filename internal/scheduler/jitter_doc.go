// Package scheduler provides scheduling primitives for cronlens.
//
// # Jitter
//
// Jitter introduces a randomized delay before a job runs. This is useful when
// multiple schedules share the same cron expression and would otherwise all
// fire at exactly the same second, causing resource contention.
//
// Basic usage:
//
//	j := scheduler.NewJitter(5 * time.Second)
//
//	// Optional: override for a specific job
//	j.Set("heavy-report", 10 * time.Second)
//
//	wrapped := scheduler.WithJitter(j, "heavy-report", myJob)
//
// WithJitter respects context cancellation: if the context is cancelled during
// the sleep, the wrapper returns ctx.Err() immediately without running the job.
//
// A zero or negative maximum disables the delay entirely for that name.
package scheduler
