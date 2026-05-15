package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ThrottlePolicy defines how a job should be rate-limited.
type ThrottlePolicy struct {
	// MinInterval is the minimum duration that must elapse between two runs.
	MinInterval time.Duration
}

// Throttle tracks last-run timestamps per job name and enforces a minimum
// interval between consecutive executions.
type Throttle struct {
	mu       sync.Mutex
	policy   ThrottlePolicy
	lastRuns map[string]time.Time
}

// NewThrottle creates a Throttle with the given policy.
func NewThrottle(policy ThrottlePolicy) *Throttle {
	return &Throttle{
		policy:   policy,
		lastRuns: make(map[string]time.Time),
	}
}

// Allow reports whether the job identified by name is allowed to run now.
// It updates the last-run timestamp when it returns true.
func (t *Throttle) Allow(name string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	if last, ok := t.lastRuns[name]; ok {
		if now.Sub(last) < t.policy.MinInterval {
			return false
		}
	}
	t.lastRuns[name] = now
	return true
}

// Reset clears the last-run record for the given job name.
func (t *Throttle) Reset(name string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.lastRuns, name)
}

// WithThrottle returns a middleware that skips execution if the job has run
// too recently according to the throttle policy. The job name is used as the
// throttle key.
func WithThrottle(th *Throttle, name string) func(JobFunc) JobFunc {
	return func(next JobFunc) JobFunc {
		return func(ctx context.Context) error {
			if !th.Allow(name) {
				return fmt.Errorf("throttled: job %q ran too recently", name)
			}
			return next(ctx)
		}
	}
}
