package scheduler

import (
	"errors"
	"sync"
	"time"
)

// OverlapPolicy defines how concurrent executions of the same job are handled.
type OverlapPolicy int

const (
	// OverlapAllow permits multiple concurrent executions.
	OverlapAllow OverlapPolicy = iota
	// OverlapSkip skips a new execution if one is already running.
	OverlapSkip
	// OverlapQueue queues the next execution until the current one finishes.
	OverlapQueue
)

// ErrSkipped is returned when an execution is skipped due to overlap policy.
var ErrSkipped = errors.New("execution skipped: job already running")

// OverlapGuard wraps a JobFunc and enforces the given OverlapPolicy.
type OverlapGuard struct {
	policy  OverlapPolicy
	running bool
	mu      sync.Mutex
	queue   chan struct{}
}

// NewOverlapGuard creates an OverlapGuard with the specified policy.
func NewOverlapGuard(policy OverlapPolicy) *OverlapGuard {
	g := &OverlapGuard{policy: policy}
	if policy == OverlapQueue {
		g.queue = make(chan struct{}, 1)
	}
	return g
}

// Wrap returns a JobFunc that enforces the overlap policy around fn.
func (g *OverlapGuard) Wrap(fn JobFunc) JobFunc {
	return func(ctx interface{ Done() <-chan struct{} }, scheduledAt time.Time) error {
		switch g.policy {
		case OverlapSkip:
			g.mu.Lock()
			if g.running {
				g.mu.Unlock()
				return ErrSkipped
			}
			g.running = true
			g.mu.Unlock()
			defer func() {
				g.mu.Lock()
				g.running = false
				g.mu.Unlock()
			}()
		case OverlapQueue:
			g.queue <- struct{}{}
			defer func() { <-g.queue }()
		}
		return fn(ctx, scheduledAt)
	}
}
