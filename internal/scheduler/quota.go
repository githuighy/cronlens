package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Quota tracks how many times a job may run within a rolling time window.
type Quota struct {
	mu      sync.Mutex
	limits  map[string]int
	counts  map[string][]time.Time
	window  time.Duration
}

// NewQuota creates a Quota with the given rolling window duration.
func NewQuota(window time.Duration) *Quota {
	return &Quota{
		limits: make(map[string]int),
		counts: make(map[string][]time.Time),
		window: window,
	}
}

// Set configures the maximum number of allowed runs for the named job.
func (q *Quota) Set(name string, max int) error {
	if name == "" {
		return fmt.Errorf("quota: name must not be empty")
	}
	if max < 1 {
		return fmt.Errorf("quota: max must be at least 1")
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	q.limits[name] = max
	return nil
}

// Delete removes the quota entry for the named job.
func (q *Quota) Delete(name string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	delete(q.limits, name)
	delete(q.counts, name)
}

// Allow reports whether the named job is permitted to run at now.
// It prunes stale timestamps and records the current run if allowed.
func (q *Quota) Allow(name string, now time.Time) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	max, ok := q.limits[name]
	if !ok {
		return true
	}

	cutoff := now.Add(-q.window)
	times := q.counts[name]
	active := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			active = append(active, t)
		}
	}

	if len(active) >= max {
		q.counts[name] = active
		return false
	}

	q.counts[name] = append(active, now)
	return true
}

// WithQuota wraps a job function, skipping execution when the quota is exceeded.
func WithQuota(q *Quota, name string, job func(ctx context.Context) error) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		if !q.Allow(name, time.Now()) {
			return nil
		}
		return job(ctx)
	}
}
