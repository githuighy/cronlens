package scheduler

import "context"

// JobFunc is the canonical job signature used throughout the scheduler package.
type JobFunc func(ctx context.Context) error
