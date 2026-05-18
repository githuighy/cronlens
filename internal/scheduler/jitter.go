package scheduler

import (
	"context"
	"math/rand"
	"sync"
	"time"
)

// Jitter adds randomized delay before executing a job, preventing thundering-herd
// problems when many schedules fire simultaneously.
type Jitter struct {
	mu      sync.Mutex
	max     map[string]time.Duration
	default_ time.Duration
}

// NewJitter creates a Jitter with a default maximum delay applied to all jobs
// unless overridden per-name via Set.
func NewJitter(defaultMax time.Duration) *Jitter {
	return &Jitter{
		max:      make(map[string]time.Duration),
		default_: defaultMax,
	}
}

// Set overrides the maximum jitter duration for a specific job name.
func (j *Jitter) Set(name string, max time.Duration) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.max[name] = max
}

// Delete removes the per-name override, reverting to the default.
func (j *Jitter) Delete(name string) {
	j.mu.Lock()
	defer j.mu.Unlock()
	delete(j.max, name)
}

// Delay returns a random duration in [0, max) for the given name.
func (j *Jitter) Delay(name string) time.Duration {
	j.mu.Lock()
	max, ok := j.max[name]
	if !ok {
		max = j.default_
	}
	j.mu.Unlock()

	if max <= 0 {
		return 0
	}
	//nolint:gosec // non-cryptographic randomness is acceptable for jitter
	return time.Duration(rand.Int63n(int64(max)))
}

// WithJitter wraps a job function, sleeping a random delay before execution.
// The name parameter is used to look up the configured maximum jitter.
func WithJitter(j *Jitter, name string, job func(ctx context.Context) error) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		delay := j.Delay(name)
		if delay > 0 {
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		return job(ctx)
	}
}
