// Package scheduler provides schedule management built on top of the cronlens
// parser and predictor packages.
//
// # History
//
// The History type offers a lightweight, thread-safe ring buffer for recording
// past schedule trigger events. It is useful for audit logging, dashboards, and
// debugging missed or late runs.
//
// Basic usage:
//
//	h := scheduler.NewHistory(50) // retain last 50 runs
//
//	// Record a run when it fires
//	h.Record(scheduledTime, time.Now(), "nightly-backup")
//
//	// Inspect history
//	for _, rec := range h.All() {
//		fmt.Printf("%s fired at %s (scheduled %s)\n",
//			rec.Label,
//			rec.TriggeredAt.Format(time.RFC3339),
//			rec.ScheduledAt.Format(time.RFC3339),
//		)
//	}
//
// The buffer is bounded: once the capacity is reached the oldest entry is
// discarded automatically. A capacity of zero or less defaults to 10.
package scheduler
