package scheduler

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// ErrDependencyNotMet is returned when a job's dependency has not completed successfully.
var ErrDependencyNotMet = errors.New("dependency not met")

// ErrCircularDependency is returned when a circular dependency is detected.
var ErrCircularDependency = errors.New("circular dependency detected")

// DependencyStore tracks job completion state for dependency resolution.
type DependencyStore struct {
	mu       sync.RWMutex
	done     map[string]bool
	deps     map[string][]string
}

// NewDependencyStore creates a new DependencyStore.
func NewDependencyStore() *DependencyStore {
	return &DependencyStore{
		done: make(map[string]bool),
		deps: make(map[string][]string),
	}
}

// Declare registers that job `name` depends on all jobs in `deps`.
// Returns ErrCircularDependency if a cycle is detected.
func (d *DependencyStore) Declare(name string, deps ...string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, dep := range deps {
		if dep == name {
			return fmt.Errorf("%w: %q depends on itself", ErrCircularDependency, name)
		}
		if d.reachable(dep, name) {
			return fmt.Errorf("%w: %q -> %q", ErrCircularDependency, name, dep)
		}
	}
	d.deps[name] = deps
	return nil
}

// reachable reports whether `from` can reach `target` via existing deps (must hold mu).
func (d *DependencyStore) reachable(from, target string) bool {
	for _, dep := range d.deps[from] {
		if dep == target || d.reachable(dep, target) {
			return true
		}
	}
	return false
}

// MarkDone records that job `name` completed successfully.
func (d *DependencyStore) MarkDone(name string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.done[name] = true
}

// Reset clears the completion state for job `name`.
func (d *DependencyStore) Reset(name string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.done, name)
}

// Ready reports whether all declared dependencies for `name` are done.
func (d *DependencyStore) Ready(name string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	for _, dep := range d.deps[name] {
		if !d.done[dep] {
			return false
		}
	}
	return true
}

// WithDependency wraps a job so it only runs when all declared dependencies are met.
// If not ready, ErrDependencyNotMet is returned without executing the job.
func WithDependency(store *DependencyStore, name string, job func(context.Context) error) func(context.Context) error {
	return func(ctx context.Context) error {
		if !store.Ready(name) {
			return fmt.Errorf("%w: job %q", ErrDependencyNotMet, name)
		}
		if err := job(ctx); err != nil {
			return err
		}
		store.MarkDone(name)
		return nil
	}
}
