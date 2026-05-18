package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// TimeoutStore manages per-job timeout durations.
type TimeoutStore struct {
	mu       sync.RWMutex
	default_ time.Duration
	overrides map[string]time.Duration
}

// NewTimeoutStore creates a TimeoutStore with the given default timeout.
// A zero or negative default disables timeout enforcement by default.
func NewTimeoutStore(defaultTimeout time.Duration) *TimeoutStore {
	return &TimeoutStore{
		default_:  defaultTimeout,
		overrides: make(map[string]time.Duration),
	}
}

// Set registers a timeout duration for a specific job name.
// A zero or negative value removes any override, reverting to the default.
func (t *TimeoutStore) Set(name string, d time.Duration) error {
	if name == "" {
		return fmt.Errorf("timeout: job name must not be empty")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if d <= 0 {
		delete(t.overrides, name)
		return nil
	}
	t.overrides[name] = d
	return nil
}

// Get returns the effective timeout for the given job name.
// Returns the per-name override if set, otherwise the store default.
// A zero duration means no timeout.
func (t *TimeoutStore) Get(name string) time.Duration {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if d, ok := t.overrides[name]; ok {
		return d
	}
	return t.default_
}

// Delete removes a per-name override, reverting to the default.
func (t *TimeoutStore) Delete(name string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.overrides, name)
}

// WithTimeout wraps a job function, cancelling the context after the
// duration returned by the store for the given name. If the duration is
// zero the job runs without an additional deadline.
func (t *TimeoutStore) WithTimeout(name string, job func(context.Context) error) func(context.Context) error {
	return func(ctx context.Context) error {
		d := t.Get(name)
		if d <= 0 {
			return job(ctx)
		}
		ctx, cancel := context.WithTimeout(ctx, d)
		defer cancel()
		return job(ctx)
	}
}
