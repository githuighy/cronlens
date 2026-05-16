package scheduler

import (
	"sync"
	"time"
)

// Metrics tracks execution statistics for a named schedule.
type Metrics struct {
	mu           sync.RWMutex
	name         string
	totalRuns    int64
	successRuns  int64
	failureRuns  int64
	totalLatency time.Duration
	lastRun      time.Time
	lastSuccess  time.Time
	lastFailure  time.Time
	lastError    error
}

// NewMetrics creates a new Metrics instance for the given schedule name.
func NewMetrics(name string) *Metrics {
	return &Metrics{name: name}
}

// Record records the result of a single job execution.
func (m *Metrics) Record(latency time.Duration, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	m.totalRuns++
	m.totalLatency += latency
	m.lastRun = now

	if err != nil {
		m.failureRuns++
		m.lastFailure = now
		m.lastError = err
	} else {
		m.successRuns++
		m.lastSuccess = now
	}
}

// Snapshot returns a point-in-time copy of the current metrics.
func (m *Metrics) Snapshot() MetricsSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var avgLatency time.Duration
	if m.totalRuns > 0 {
		avgLatency = m.totalLatency / time.Duration(m.totalRuns)
	}

	return MetricsSnapshot{
		Name:         m.name,
		TotalRuns:    m.totalRuns,
		SuccessRuns:  m.successRuns,
		FailureRuns:  m.failureRuns,
		AvgLatency:   avgLatency,
		LastRun:      m.lastRun,
		LastSuccess:  m.lastSuccess,
		LastFailure:  m.lastFailure,
		LastError:    m.lastError,
	}
}

// Reset clears all recorded metrics.
func (m *Metrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.totalRuns = 0
	m.successRuns = 0
	m.failureRuns = 0
	m.totalLatency = 0
	m.lastRun = time.Time{}
	m.lastSuccess = time.Time{}
	m.lastFailure = time.Time{}
	m.lastError = nil
}

// MetricsSnapshot is an immutable point-in-time view of Metrics.
type MetricsSnapshot struct {
	Name        string
	TotalRuns   int64
	SuccessRuns int64
	FailureRuns int64
	AvgLatency  time.Duration
	LastRun     time.Time
	LastSuccess time.Time
	LastFailure time.Time
	LastError   error
}

// SuccessRate returns the ratio of successful runs to total runs (0–1).
func (s MetricsSnapshot) SuccessRate() float64 {
	if s.TotalRuns == 0 {
		return 0
	}
	return float64(s.SuccessRuns) / float64(s.TotalRuns)
}
