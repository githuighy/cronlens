package scheduler

import (
	"context"
	"fmt"
	"time"
)

// ChainEntry holds a named job within a chain.
type ChainEntry struct {
	Name string
	Job  func(ctx context.Context) error
}

// JobChain executes a sequence of jobs in order, stopping on the first failure.
type JobChain struct {
	entries []ChainEntry
	timeout time.Duration
}

// NewJobChain creates a new JobChain with an optional per-step timeout.
func NewJobChain(timeout time.Duration) *JobChain {
	return &JobChain{timeout: timeout}
}

// Add appends a named job to the chain.
func (c *JobChain) Add(name string, job func(ctx context.Context) error) *JobChain {
	c.entries = append(c.entries, ChainEntry{Name: name, Job: job})
	return c
}

// Run executes all jobs in sequence. If any job fails, execution stops and the
// error is returned wrapped with the failing step name.
func (c *JobChain) Run(ctx context.Context) error {
	for _, entry := range c.entries {
		stepCtx := ctx
		var cancel context.CancelFunc
		if c.timeout > 0 {
			stepCtx, cancel = context.WithTimeout(ctx, c.timeout)
		}
		err := entry.Job(stepCtx)
		if cancel != nil {
			cancel()
		}
		if err != nil {
			return fmt.Errorf("chain step %q failed: %w", entry.Name, err)
		}
	}
	return nil
}

// Len returns the number of steps in the chain.
func (c *JobChain) Len() int {
	return len(c.entries)
}

// Names returns the ordered list of step names.
func (c *JobChain) Names() []string {
	names := make([]string, len(c.entries))
	for i, e := range c.entries {
		names[i] = e.Name
	}
	return names
}
