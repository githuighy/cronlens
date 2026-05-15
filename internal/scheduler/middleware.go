package scheduler

import (
	"context"
	"log"
	"time"
)

// JobFunc is the function signature for a scheduled job.
type JobFunc func(ctx context.Context) error

// Middleware wraps a JobFunc with additional behavior.
type Middleware func(next JobFunc) JobFunc

// Chain applies a series of middlewares to a JobFunc, in order.
// The first middleware in the slice is the outermost wrapper.
func Chain(job JobFunc, middlewares ...Middleware) JobFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		job = middlewares[i](job)
	}
	return job
}

// WithLogging returns a Middleware that logs job start, end, and any error.
func WithLogging(logger *log.Logger) Middleware {
	return func(next JobFunc) JobFunc {
		return func(ctx context.Context) error {
			start := time.Now()
			logger.Printf("[cronlens] job starting")
			err := next(ctx)
			dur := time.Since(start)
			if err != nil {
				logger.Printf("[cronlens] job failed after %s: %v", dur, err)
			} else {
				logger.Printf("[cronlens] job completed in %s", dur)
			}
			return err
		}
	}
}

// WithTimeout returns a Middleware that cancels the job context after d.
func WithTimeout(d time.Duration) Middleware {
	return func(next JobFunc) JobFunc {
		return func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, d)
			defer cancel()
			return next(ctx)
		}
	}
}

// WithRecover returns a Middleware that recovers from panics and returns them as errors.
func WithRecover() Middleware {
	return func(next JobFunc) JobFunc {
		return func(ctx context.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("job panicked: %v", r)
				}
			}()
			return next(ctx)
		}
	}
}
