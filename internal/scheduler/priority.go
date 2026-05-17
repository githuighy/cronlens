package scheduler

import (
	"context"
	"sort"
	"sync"
)

// Priority represents the execution priority of a job.
type Priority int

const (
	PriorityLow    Priority = 0
	PriorityNormal Priority = 50
	PriorityHigh   Priority = 100
)

// priorityJob wraps a job function with a name and priority level.
type priorityJob struct {
	name     string
	priority Priority
	job      func(ctx context.Context) error
}

// PriorityQueue holds jobs ordered by priority (highest first).
type PriorityQueue struct {
	mu   sync.Mutex
	jobs []priorityJob
}

// NewPriorityQueue creates an empty PriorityQueue.
func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{}
}

// Enqueue adds a job with the given priority.
func (pq *PriorityQueue) Enqueue(name string, priority Priority, job func(ctx context.Context) error) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	pq.jobs = append(pq.jobs, priorityJob{name: name, priority: priority, job: job})
	sort.SliceStable(pq.jobs, func(i, j int) bool {
		return pq.jobs[i].priority > pq.jobs[j].priority
	})
}

// Dequeue removes and returns the highest-priority job.
// Returns false if the queue is empty.
func (pq *PriorityQueue) Dequeue() (string, func(ctx context.Context) error, bool) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	if len(pq.jobs) == 0 {
		return "", nil, false
	}
	j := pq.jobs[0]
	pq.jobs = pq.jobs[1:]
	return j.name, j.job, true
}

// Len returns the number of queued jobs.
func (pq *PriorityQueue) Len() int {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	return len(pq.jobs)
}

// WithPriority wraps a job so it is enqueued into the given PriorityQueue
// rather than executed immediately. Call DrainQueue to run pending jobs.
func WithPriority(pq *PriorityQueue, name string, priority Priority, job func(ctx context.Context) error) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		pq.Enqueue(name, priority, job)
		return nil
	}
}

// DrainQueue executes all queued jobs in priority order using the provided context.
// Execution stops on the first error.
func DrainQueue(ctx context.Context, pq *PriorityQueue) error {
	for {
		_, job, ok := pq.Dequeue()
		if !ok {
			return nil
		}
		if err := job(ctx); err != nil {
			return err
		}
	}
}
