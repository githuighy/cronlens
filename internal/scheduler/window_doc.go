// Package scheduler provides scheduling primitives for cronlens.
//
// # Window
//
// WindowStore allows restricting job execution to specific time-of-day
// windows. Each window is identified by a name and defined as a pair of
// durations measured from midnight (00:00:00) in the job's local timezone.
//
// Example — only run between 09:00 and 17:00:
//
//	ws := scheduler.NewWindowStore()
//	_ = ws.Set("business-hours", 9*time.Hour, 17*time.Hour)
//
//	wrapped := scheduler.WithWindow(ws, "business-hours", myJob)
//
// If the window does not exist or the current time is outside the window,
// WithWindow returns nil without calling the underlying job.
package scheduler
