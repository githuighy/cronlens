// Package scheduler provides scheduling primitives for cron-based job execution.
//
// # Webhook Notifier
//
// WebhookNotifier sends an HTTP POST request to a configured URL whenever a
// scheduled job fires. The request body is a JSON-encoded WebhookPayload that
// includes the job name, cron expression, fired-at timestamp, and outcome.
//
// Basic usage:
//
//	notifier := scheduler.NewWebhookNotifier("https://example.com/hooks/cron")
//
//	wrapped := scheduler.Chain(
//		myJob,
//		scheduler.WithWebhook(notifier, "backup", "0 2 * * *"),
//	)
//
// The notifier performs a best-effort delivery: if the HTTP call fails, the
// error is silently discarded so that the underlying job result is preserved.
//
// Timeout for each HTTP request defaults to 5 seconds and can be adjusted by
// setting WebhookNotifier.Timeout directly.
package scheduler
