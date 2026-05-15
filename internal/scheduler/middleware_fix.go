package scheduler

// JobFunc is the function signature for a schedulable job.
type JobFunc func(ctx interface{ Done() <-chan struct{} }) error
