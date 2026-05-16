package scheduler

import (
	"context"
	"time"
)

// WithMetrics wraps a JobFunc so that every execution is recorded in m.
// It measures wall-clock latency and captures any returned error.
func WithMetrics(m *Metrics, job func(ctx context.Context) error) func(ctx context.Context) error {
	if m == nil {
		return job
	}
	return func(ctx context.Context) error {
		start := time.Now()
		err := job(ctx)
		m.Record(time.Since(start), err)
		return err
	}
}
