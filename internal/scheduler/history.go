package scheduler

import (
	"sync"
	"time"
)

// RunRecord captures a single scheduled run event.
type RunRecord struct {
	ScheduledAt time.Time
	TriggeredAt time.Time
	Label       string
}

// History tracks recent run records for a schedule in a bounded ring buffer.
type History struct {
	mu      sync.RWMutex
	records []RunRecord
	cap     int
}

// NewHistory creates a History that retains at most maxEntries records.
func NewHistory(maxEntries int) *History {
	if maxEntries <= 0 {
		maxEntries = 10
	}
	return &History{
		records: make([]RunRecord, 0, maxEntries),
		cap:     maxEntries,
	}
}

// Record appends a new run record, evicting the oldest if at capacity.
func (h *History) Record(scheduled, triggered time.Time, label string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if len(h.records) >= h.cap {
		h.records = h.records[1:]
	}
	h.records = append(h.records, RunRecord{
		ScheduledAt: scheduled,
		TriggeredAt: triggered,
		Label:       label,
	})
}

// All returns a copy of all stored records in chronological order.
func (h *History) All() []RunRecord {
	h.mu.RLock()
	defer h.mu.RUnlock()

	out := make([]RunRecord, len(h.records))
	copy(out, h.records)
	return out
}

// Last returns the most recent RunRecord and whether one exists.
func (h *History) Last() (RunRecord, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.records) == 0 {
		return RunRecord{}, false
	}
	return h.records[len(h.records)-1], true
}

// Len returns the number of records currently stored.
func (h *History) Len() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.records)
}
