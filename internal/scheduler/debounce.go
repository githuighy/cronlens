package scheduler

import (
	"context"
	"sync"
	"time"
)

// Debounce holds state for a debounced job execution.
// It ensures that rapid successive calls only result in a single
// execution after the specified wait duration has elapsed.
type Debounce struct {
	mu      sync.Mutex
	timers  map[string]*time.Timer
	wait    time.Duration
}

// NewDebounce creates a new Debounce with the given wait duration.
func NewDebounce(wait time.Duration) *Debounce {
	return &Debounce{
		timers: make(map[string]*time.Timer),
		wait:   wait,
	}
}

// Trigger schedules fn to run after the debounce wait period for the given name.
// If Trigger is called again for the same name before the timer fires, the
// previous timer is cancelled and a new one is started.
func (d *Debounce) Trigger(name string, fn func(ctx context.Context) error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if t, ok := d.timers[name]; ok {
		t.Stop()
	}

	d.timers[name] = time.AfterFunc(d.wait, func() {
		d.mu.Lock()
		delete(d.timers, name)
		d.mu.Unlock()

		_ = fn(context.Background())
	})
}

// Cancel stops a pending debounced call for the given name, if any.
func (d *Debounce) Cancel(name string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if t, ok := d.timers[name]; ok {
		t.Stop()
		delete(d.timers, name)
	}
}

// Pending returns true if a debounced call is waiting to fire for the given name.
func (d *Debounce) Pending(name string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	_, ok := d.timers[name]
	return ok
}

// WithDebounce wraps a job function so that concurrent rapid invocations
// are collapsed into a single execution after the wait period.
func WithDebounce(d *Debounce, name string, fn func(ctx context.Context) error) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		d.Trigger(name, fn)
		return nil
	}
}
