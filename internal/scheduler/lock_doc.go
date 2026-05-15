// Package scheduler provides scheduling primitives for cron-based job execution.
//
// # Lock
//
// Lock offers named mutual exclusion for job runs, preventing the same job
// from executing concurrently across goroutines. Each lock entry can carry
// an optional TTL so that stale locks (e.g. from a crashed process) are
// automatically released after the deadline expires.
//
// Basic usage:
//
//	l := scheduler.NewLock(30 * time.Second)
//
//	if err := l.Acquire("my-job"); err != nil {
//		// job is already running – skip this tick
//		return
//	}
//	defer l.Release("my-job")
//
//	// … run the job …
//
// A zero TTL creates locks that never expire automatically; callers must
// always call Release to avoid permanent blockage.
package scheduler
