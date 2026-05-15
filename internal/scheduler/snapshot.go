package scheduler

import (
	"sync"
	"time"
)

// Snapshot captures the state of a Schedule at a point in time.
type Snapshot struct {
	Name        string            `json:"name"`
	Expression  string            `json:"expression"`
	Timezone    string            `json:"timezone"`
	Tags        map[string]string `json:"tags,omitempty"`
	NextRun     time.Time         `json:"next_run"`
	LastRun     *time.Time        `json:"last_run,omitempty"`
	CapturedAt  time.Time         `json:"captured_at"`
}

// SnapshotStore holds a bounded set of named snapshots.
type SnapshotStore struct {
	mu        sync.RWMutex
	snapshots map[string]Snapshot
}

// NewSnapshotStore creates an empty SnapshotStore.
func NewSnapshotStore() *SnapshotStore {
	return &SnapshotStore{
		snapshots: make(map[string]Snapshot),
	}
}

// Save stores a snapshot under its Name, overwriting any previous entry.
func (s *SnapshotStore) Save(snap Snapshot) {
	s.mu.Lock()
	defer s.mu.Unlock()
	snap.CapturedAt = time.Now().UTC()
	s.snapshots[snap.Name] = snap
}

// Get retrieves a snapshot by name. Returns false if not found.
func (s *SnapshotStore) Get(name string) (Snapshot, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	snap, ok := s.snapshots[name]
	return snap, ok
}

// Delete removes a snapshot by name.
func (s *SnapshotStore) Delete(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.snapshots, name)
}

// All returns a copy of all stored snapshots.
func (s *SnapshotStore) All() []Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Snapshot, 0, len(s.snapshots))
	for _, snap := range s.snapshots {
		out = append(out, snap)
	}
	return out
}

// Count returns the number of stored snapshots.
func (s *SnapshotStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.snapshots)
}
