package scheduler

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// Group manages a named collection of Schedules, allowing batch operations
// such as retrieving the next run across all members.
type Group struct {
	mu        sync.RWMutex
	schedules map[string]*Schedule
}

// NewGroup returns an empty Group.
func NewGroup() *Group {
	return &Group{
		schedules: make(map[string]*Schedule),
	}
}

// Add registers a Schedule under the given name.
// Returns an error if the name is already taken.
func (g *Group) Add(name string, s *Schedule) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, exists := g.schedules[name]; exists {
		return fmt.Errorf("scheduler/group: name %q already registered", name)
	}
	g.schedules[name] = s
	return nil
}

// Remove deletes the Schedule registered under name.
func (g *Group) Remove(name string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.schedules, name)
}

// Get returns the Schedule for name and whether it was found.
func (g *Group) Get(name string) (*Schedule, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	s, ok := g.schedules[name]
	return s, ok
}

// Names returns a sorted list of registered schedule names.
func (g *Group) Names() []string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	names := make([]string, 0, len(g.schedules))
	for n := range g.schedules {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// NextEntry holds the name of a schedule and its upcoming run time.
type NextEntry struct {
	Name string
	Next time.Time
}

// NextAll returns the next run time for every schedule in the group,
// computed relative to from, sorted by ascending run time.
func (g *Group) NextAll(from time.Time) []NextEntry {
	g.mu.RLock()
	defer g.mu.RUnlock()
	entries := make([]NextEntry, 0, len(g.schedules))
	for name, s := range g.schedules {
		next := s.NextRun(from)
		entries = append(entries, NextEntry{Name: name, Next: next})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Next.Before(entries[j].Next)
	})
	return entries
}
