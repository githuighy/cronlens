package scheduler

import (
	"context"
	"fmt"
	"log"
)

// WithPriorityLogging wraps DrainQueue with structured log output, printing
// each job's name and outcome as it is dequeued and executed.
//
//	if err := WithPriorityLogging(ctx, pq, log.Default()); err != nil {
//	    // handle first failure
//	}
func WithPriorityLogging(ctx context.Context, pq *PriorityQueue, logger *log.Logger) error {
	for {
		name, job, ok := pq.Dequeue()
		if !ok {
			return nil
		}
		logger.Printf("[priority] running job %q", name)
		if err := job(ctx); err != nil {
			logger.Printf("[priority] job %q failed: %v", name, err)
			return fmt.Errorf("priority job %q: %w", name, err)
		}
		logger.Printf("[priority] job %q completed", name)
	}
}
