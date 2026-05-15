// Package scheduler provides scheduling primitives for cron-based job execution.
//
// # RateLimit
//
// RateLimit enforces a maximum number of executions for a named job within a
// sliding time window. It is safe for concurrent use.
//
// Basic usage:
//
//	rl := scheduler.NewRateLimit(5, time.Minute)
//
//	if rl.Allow("my-job", time.Now()) {
//		// proceed with execution
//	} else {
//		// skip or queue
//	}
//
// Middleware usage:
//
//	handler = scheduler.Chain(
//		baseHandler,
//		scheduler.WithRateLimit(rl, "my-job"),
//	)
//
// The sliding window means that the oldest recorded timestamps are pruned on
// each call to Allow, so the limit applies to the most recent window duration
// rather than a fixed calendar period.
package scheduler
