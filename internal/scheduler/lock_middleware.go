package scheduler

import (
	"context"
	"time"
)

// WithLock wraps a JobFunc so that at most one execution of name can run at a
// time. If the lock cannot be acquired the invocation is silently skipped.
// The lock is automatically released when the wrapped function returns.
//
// ttl controls the maximum duration a lock entry is considered valid.
// Pass 0 for no automatic expiry (the lock is held until Release is called).
//
//	guard := scheduler.NewLock(30 * time.Second)
//	safe := scheduler.WithLock(guard, "nightly-report", myJob)
func WithLock(l *Lock, name string, next JobFunc) JobFunc {
	return func(ctx context.Context) error {
		if err := l.Acquire(name); err != nil {
			// Another instance is running – skip this tick.
			return nil
		}
		defer l.Release(name)
		return next(ctx)
	}
}

// WithLockOrError is like WithLock but returns ErrLockTimeout to the caller
// instead of silently skipping when the lock is already held.
func WithLockOrError(l *Lock, name string, next JobFunc) JobFunc {
	return func(ctx context.Context) error {
		if err := l.Acquire(name); err != nil {
			return err
		}
		defer l.Release(name)
		return next(ctx)
	}
}

// JobFunc is the signature expected by middleware helpers in this package.
// It matches the context-aware handler used by WithLogging and WithTimeout.
type JobFunc func(ctx context.Context) error

// ensure compile-time that time import is used (TTL constant below)
var _ = time.Second
