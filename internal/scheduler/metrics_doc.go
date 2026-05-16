// Package scheduler provides scheduling primitives for cronlens.
//
// # Metrics
//
// Metrics tracks per-schedule execution statistics including total runs,
// success/failure counts, average latency, and timestamps of the last
// successful and failed executions.
//
// Basic usage:
//
//	m := scheduler.NewMetrics("my-job")
//	wrapped := scheduler.WithMetrics(m, myJobFunc)
//
//	// later…
//	snap := m.Snapshot()
//	fmt.Printf("success rate: %.1f%%\n", snap.SuccessRate()*100)
//
// WithMetrics is a middleware-style wrapper that is composable with the
// other middleware helpers (Chain, WithLogging, WithTimeout, etc.).
//
// Metrics is safe for concurrent use. Reset() clears all counters and
// timestamps, which is useful in tests or when rotating reporting windows.
package scheduler
