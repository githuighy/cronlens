package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Cooldown enforces a minimum wait period between successive runs of a named job.
// If a job is triggered before the cooldown period has elapsed, it is skipped.
type Cooldown struct {
	mu       sync.Mutex
	lastRuns map[string]time.Time
	period   time.Duration
}

// NewCooldown creates a Cooldown with the given minimum period between runs.
func NewCooldown(period time.Duration) *Cooldown {
	return &Cooldown{
		lastRuns: make(map[string]time.Time),
		period:   period,
	}
}

// Allow returns true if the named job is allowed to run (i.e. the cooldown
// period has elapsed since the last run). It records the current time as the
// last run time when it returns true.
func (c *Cooldown) Allow(name string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	if last, ok := c.lastRuns[name]; ok {
		if now.Sub(last) < c.period {
			return false
		}
	}
	c.lastRuns[name] = now
	return true
}

// Reset clears the cooldown state for the named job, allowing it to run
// immediately on the next call to Allow.
func (c *Cooldown) Reset(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.lastRuns, name)
}

// Remaining returns the duration remaining in the cooldown period for the
// named job. Returns 0 if the job is ready to run.
func (c *Cooldown) Remaining(name string) time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()

	last, ok := c.lastRuns[name]
	if !ok {
		return 0
	}
	remaining := c.period - time.Since(last)
	if remaining < 0 {
		return 0
	}
	return remaining
}

// WithCooldown wraps a job function, skipping execution if the cooldown period
// has not elapsed since the last successful run.
func WithCooldown(name string, c *Cooldown, job func(ctx context.Context) error) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		if !c.Allow(name) {
			return fmt.Errorf("cooldown active for %q: %.0fs remaining", name, c.Remaining(name).Seconds())
		}
		return job(ctx)
	}
}
