package scheduler

import (
	"context"
	"log"
	"time"
)

// WithChainLogging wraps each step in the chain with entry/exit logging.
// It replaces the chain's entries with instrumented versions and returns
// the same chain for fluent chaining.
func WithChainLogging(c *JobChain) *JobChain {
	instrumented := NewJobChain(0)
	for _, entry := range c.entries {
		e := entry // capture
		instrumented.Add(e.Name, func(ctx context.Context) error {
			start := time.Now()
			log.Printf("[chain] step %q starting", e.Name)
			err := e.Job(ctx)
			dur := time.Since(start)
			if err != nil {
				log.Printf("[chain] step %q failed after %s: %v", e.Name, dur, err)
			} else {
				log.Printf("[chain] step %q completed in %s", e.Name, dur)
			}
			return err
		})
	}
	// Copy timeout from original chain.
	instrumented.timeout = c.timeout
	return instrumented
}

// AsJob converts a JobChain into a single func(ctx context.Context) error
// suitable for use with middleware, groups, or retry policies.
func AsJob(c *JobChain) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		return c.Run(ctx)
	}
}
