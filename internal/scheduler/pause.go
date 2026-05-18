package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// PauseState represents the current pause state of a named schedule.
type PauseState struct {
	Paused    bool
	PausedAt  time.Time
	ResumedAt time.Time
	Reason    string
}

// PauseStore tracks pause/resume state for named schedules.
type PauseStore struct {
	mu     sync.RWMutex
	states map[string]*PauseState
}

// NewPauseStore returns an empty PauseStore.
func NewPauseStore() *PauseStore {
	return &PauseStore{states: make(map[string]*PauseState)}
}

// Pause marks a named schedule as paused with an optional reason.
// Returns an error if the schedule is already paused.
func (p *PauseStore) Pause(name, reason string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if s, ok := p.states[name]; ok && s.Paused {
		return fmt.Errorf("schedule %q is already paused", name)
	}
	p.states[name] = &PauseState{
		Paused:   true,
		PausedAt: time.Now(),
		Reason:   reason,
	}
	return nil
}

// Resume marks a named schedule as resumed.
// Returns an error if the schedule is not currently paused.
func (p *PauseStore) Resume(name string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	s, ok := p.states[name]
	if !ok || !s.Paused {
		return fmt.Errorf("schedule %q is not paused", name)
	}
	s.Paused = false
	s.ResumedAt = time.Now()
	return nil
}

// IsPaused reports whether the named schedule is currently paused.
func (p *PauseStore) IsPaused(name string) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	s, ok := p.states[name]
	return ok && s.Paused
}

// State returns the PauseState for the named schedule, or nil if unknown.
func (p *PauseStore) State(name string) *PauseState {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if s, ok := p.states[name]; ok {
		copy := *s
		return &copy
	}
	return nil
}

// WithPause wraps a job so it skips execution when the named schedule is paused.
func WithPause(store *PauseStore, name string, job func(context.Context) error) func(context.Context) error {
	return func(ctx context.Context) error {
		if store.IsPaused(name) {
			return nil
		}
		return job(ctx)
	}
}
