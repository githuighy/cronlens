// Package scheduler provides scheduling primitives for cronlens.
//
// # AuditLog
//
// AuditLog maintains a per-job ring buffer of execution events. Each event
// captures the job name, start/end timestamps, elapsed duration, and any
// error returned.
//
// Usage:
//
//	log := scheduler.NewAuditLog(50)
//
//	wrapped := scheduler.WithAudit("my-job", log, func(ctx context.Context) error {
//		// job logic
//		return nil
//	})
//
//	// After execution:
//	// last, ok := log.Last("my-job")
//	// all  := log.All("my-job")
//
// The capacity passed to NewAuditLog controls how many events are retained
// per job name. Once the cap is reached the oldest event is evicted.
package scheduler
