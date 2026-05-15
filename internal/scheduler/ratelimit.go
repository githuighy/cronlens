package scheduler

import (
	"fmt"
	"sync"
	"time"
)

// RateLimit enforces a maximum number of job executions within a sliding window.
type RateLimit struct {
	mu       sync.Mutex
	window   time.Duration
	maxCalls int
	records  map[string][]time.Time
}

// NewRateLimit creates a RateLimit that allows at most maxCalls executions
// per name within the given window duration.
func NewRateLimit(maxCalls int, window time.Duration) *RateLimit {
	return &RateLimit{
		window:   window,
		maxCalls: maxCalls,
		records:  make(map[string][]time.Time),
	}
}

// Allow reports whether a call for the given name is permitted at now.
// It prunes stale timestamps and records the call if allowed.
func (r *RateLimit) Allow(name string, now time.Time) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	cutoff := now.Add(-r.window)
	times := r.records[name]

	// Prune entries outside the window.
	active := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			active = append(active, t)
		}
	}

	if len(active) >= r.maxCalls {
		r.records[name] = active
		return false
	}

	r.records[name] = append(active, now)
	return true
}

// Reset clears all recorded calls for the given name.
func (r *RateLimit) Reset(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.records, name)
}

// Remaining returns the number of calls still permitted for name within the
// current window relative to now.
func (r *RateLimit) Remaining(name string, now time.Time) int {
	r.mu.Lock()
	defer r.mu.Unlock()

	cutoff := now.Add(-r.window)
	count := 0
	for _, t := range r.records[name] {
		if t.After(cutoff) {
			count++
		}
	}
	rem := r.maxCalls - count
	if rem < 0 {
		return 0
	}
	return rem
}

// WithRateLimit returns a middleware that skips execution when the rate limit
// for the given name has been exceeded.
func WithRateLimit(rl *RateLimit, name string) func(JobFunc) JobFunc {
	return func(next JobFunc) JobFunc {
		return func(ctx interface{ Done() <-chan struct{} }) error {
			if !rl.Allow(name, time.Now()) {
				return fmt.Errorf("rate limit exceeded for %q", name)
			}
			return next(ctx)
		}
	}
}
