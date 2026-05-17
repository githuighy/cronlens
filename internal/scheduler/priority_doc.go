// Package scheduler provides scheduling primitives for cron-based job execution.
//
// # Priority Queue
//
// PriorityQueue allows jobs to be enqueued with an explicit priority level and
// dequeued in highest-priority-first order. Three convenience constants are
// provided:
//
//	- PriorityLow    (0)   — background / best-effort work
//	- PriorityNormal (50)  — default priority for most jobs
//	- PriorityHigh   (100) — time-sensitive or critical jobs
//
// Custom integer values outside this range are also accepted.
//
// # Usage
//
//	pq := scheduler.NewPriorityQueue()
//	pq.Enqueue("cleanup", scheduler.PriorityLow, cleanupFn)
//	pq.Enqueue("alert",   scheduler.PriorityHigh, alertFn)
//
//	// Runs alertFn first, then cleanupFn.
//	if err := scheduler.DrainQueue(ctx, pq); err != nil {
//	    log.Println("job failed:", err)
//	}
package scheduler
