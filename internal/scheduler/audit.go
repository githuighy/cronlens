package scheduler

import (
	"context"
	"sync"
	"time"
)

// AuditEvent records a single job execution attempt.
type AuditEvent struct {
	Name      string
	StartedAt time.Time
	EndedAt   time.Time
	Duration  time.Duration
	Err       error
	Skipped   bool
}

// AuditLog stores a bounded list of audit events per job.
type AuditLog struct {
	mu     sync.Mutex
	events map[string][]AuditEvent
	cap    int
}

// NewAuditLog creates an AuditLog with the given per-job capacity.
func NewAuditLog(cap int) *AuditLog {
	if cap <= 0 {
		cap = 100
	}
	return &AuditLog{
		events: make(map[string][]AuditEvent),
		cap:    cap,
	}
}

// Record appends an AuditEvent for the named job, evicting the oldest if at capacity.
func (a *AuditLog) Record(name string, event AuditEvent) {
	a.mu.Lock()
	defer a.mu.Unlock()
	evs := a.events[name]
	if len(evs) >= a.cap {
		evs = evs[1:]
	}
	a.events[name] = append(evs, event)
}

// All returns a copy of all audit events recorded for the named job.
func (a *AuditLog) All(name string) []AuditEvent {
	a.mu.Lock()
	defer a.mu.Unlock()
	src := a.events[name]
	out := make([]AuditEvent, len(src))
	copy(out, src)
	return out
}

// Last returns the most recent audit event for the named job, or false if none.
func (a *AuditLog) Last(name string) (AuditEvent, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	evs := a.events[name]
	if len(evs) == 0 {
		return AuditEvent{}, false
	}
	return evs[len(evs)-1], true
}

// Clear removes all recorded events for the named job.
func (a *AuditLog) Clear(name string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.events, name)
}

// WithAudit wraps a job function so each execution is recorded in the AuditLog.
func WithAudit(name string, log *AuditLog, job func(context.Context) error) func(context.Context) error {
	return func(ctx context.Context) error {
		start := time.Now()
		err := job(ctx)
		end := time.Now()
		log.Record(name, AuditEvent{
			Name:      name,
			StartedAt: start,
			EndedAt:   end,
			Duration:  end.Sub(start),
			Err:       err,
		})
		return err
	}
}
