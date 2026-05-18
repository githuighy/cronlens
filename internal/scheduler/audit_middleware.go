package scheduler

import (
	"context"
	"log"
	"time"
)

// WithAuditLogging wraps a job so that each run is recorded in the AuditLog
// and a summary line is printed to the standard logger.
func WithAuditLogging(name string, al *AuditLog, job func(context.Context) error) func(context.Context) error {
	return func(ctx context.Context) error {
		start := time.Now()
		err := job(ctx)
		end := time.Now()

		ev := AuditEvent{
			Name:      name,
			StartedAt: start,
			EndedAt:   end,
			Duration:  end.Sub(start),
			Err:       err,
		}
		al.Record(name, ev)

		if err != nil {
			log.Printf("[audit] job=%s status=error duration=%s err=%v", name, ev.Duration, err)
		} else {
			log.Printf("[audit] job=%s status=ok duration=%s", name, ev.Duration)
		}
		return err
	}
}

// WithAuditSkip wraps a job so that skipped executions (when the job returns
// ErrSkipped) are flagged in the audit record rather than treated as errors.
func WithAuditSkip(name string, al *AuditLog, job func(context.Context) error) func(context.Context) error {
	return func(ctx context.Context) error {
		start := time.Now()
		err := job(ctx)
		end := time.Now()

		ev := AuditEvent{
			Name:      name,
			StartedAt: start,
			EndedAt:   end,
			Duration:  end.Sub(start),
		}

		if err == ErrSkipped {
			ev.Skipped = true
			ev.Err = nil
		} else {
			ev.Err = err
		}

		al.Record(name, ev)
		return err
	}
}
