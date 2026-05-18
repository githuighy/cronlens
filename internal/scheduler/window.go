package scheduler

import (
	"fmt"
	"sync"
	"time"
)

// Window represents a time range during which a job is allowed to run.
type Window struct {
	Start time.Duration // offset from midnight
	End   time.Duration // offset from midnight
}

// WindowStore holds named execution windows.
type WindowStore struct {
	mu      sync.RWMutex
	windows map[string]Window
}

// NewWindowStore creates an empty WindowStore.
func NewWindowStore() *WindowStore {
	return &WindowStore{
		windows: make(map[string]Window),
	}
}

// Set registers a named window. Start must be before End.
func (ws *WindowStore) Set(name string, start, end time.Duration) error {
	if name == "" {
		return fmt.Errorf("window name must not be empty")
	}
	if start >= end {
		return fmt.Errorf("window start (%v) must be before end (%v)", start, end)
	}
	ws.mu.Lock()
	defer ws.mu.Unlock()
	ws.windows[name] = Window{Start: start, End: end}
	return nil
}

// Get retrieves a window by name.
func (ws *WindowStore) Get(name string) (Window, bool) {
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	w, ok := ws.windows[name]
	return w, ok
}

// Delete removes a window by name.
func (ws *WindowStore) Delete(name string) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	delete(ws.windows, name)
}

// InWindow reports whether t falls within the named window.
// Returns false if the window does not exist.
func (ws *WindowStore) InWindow(name string, t time.Time) bool {
	w, ok := ws.Get(name)
	if !ok {
		return false
	}
	midnight := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	offset := t.Sub(midnight)
	return offset >= w.Start && offset < w.End
}

// WithWindow wraps a JobFunc so it only runs when the current time falls
// within the named window. Skipped runs are not treated as errors.
func WithWindow(ws *WindowStore, name string, fn JobFunc) JobFunc {
	return func(ctx context.Context) error {
		if !ws.InWindow(name, time.Now()) {
			return nil
		}
		return fn(ctx)
	}
}
