// Package scheduler provides scheduling primitives for cronlens.
//
// # Quota
//
// Quota enforces a maximum number of job executions within a rolling time
// window. Unlike RateLimit, which counts calls per fixed window boundary,
// Quota uses a sliding window so bursts are smoothed over time.
//
// Example usage:
//
//	q := scheduler.NewQuota(24 * time.Hour)
//	_ = q.Set("daily-report", 3) // allow at most 3 runs per day
//
//	wrapped := scheduler.WithQuota(q, "daily-report", func(ctx context.Context) error {
//		// job logic
//		return nil
//	})
//
// When the quota is exceeded the wrapped job returns nil immediately and
// the run is silently skipped. Pair with WithLogging or WithMetrics to
// observe skipped executions.
package scheduler
