package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Deadline represents a named deadline with an expiry time.
type Deadline struct {
	Name      string
	ExpiresAt time.Time
}

// DeadlineStore tracks named deadlines and whether they have been met.
type DeadlineStore struct {
	mu        sync.Mutex
	deadlines map[string]time.Time
}

// NewDeadlineStore creates an empty DeadlineStore.
func NewDeadlineStore() *DeadlineStore {
	return &DeadlineStore{
		deadlines: make(map[string]time.Time),
	}
}

// Set registers or updates a named deadline.
func (d *DeadlineStore) Set(name string, expiresAt time.Time) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.deadlines[name] = expiresAt
}

// Get returns the deadline for the given name and whether it exists.
func (d *DeadlineStore) Get(name string) (time.Time, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	t, ok := d.deadlines[name]
	return t, ok
}

// Expired reports whether the named deadline exists and has passed.
func (d *DeadlineStore) Expired(name string, now time.Time) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	t, ok := d.deadlines[name]
	if !ok {
		return false
	}
	return now.After(t)
}

// Delete removes a named deadline.
func (d *DeadlineStore) Delete(name string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.deadlines, name)
}

// WithDeadline wraps a job so it fails if the named deadline has expired.
// The deadline is checked at invocation time using the provided store.
func WithDeadline(store *DeadlineStore, name string, job func(context.Context) error) func(context.Context) error {
	return func(ctx context.Context) error {
		if store.Expired(name, time.Now()) {
			return fmt.Errorf("deadline %q has expired", name)
		}
		return job(ctx)
	}
}
