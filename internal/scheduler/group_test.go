package scheduler

import (
	"testing"
	"time"
)

func mustSchedule(t *testing.T, expr, tz string) *Schedule {
	t.Helper()
	s, err := New(expr, tz)
	if err != nil {
		t.Fatalf("New(%q, %q): %v", expr, tz, err)
	}
	return s
}

func TestGroup_AddAndGet(t *testing.T) {
	g := NewGroup()
	s := mustSchedule(t, "0 9 * * 1-5", "UTC")

	if err := g.Add("weekday-morning", s); err != nil {
		t.Fatalf("Add: %v", err)
	}
	got, ok := g.Get("weekday-morning")
	if !ok || got != s {
		t.Errorf("Get returned wrong schedule")
	}
}

func TestGroup_Add_DuplicateName(t *testing.T) {
	g := NewGroup()
	s := mustSchedule(t, "* * * * *", "UTC")
	_ = g.Add("job", s)
	if err := g.Add("job", s); err == nil {
		t.Error("expected error for duplicate name, got nil")
	}
}

func TestGroup_Remove(t *testing.T) {
	g := NewGroup()
	s := mustSchedule(t, "* * * * *", "UTC")
	_ = g.Add("job", s)
	g.Remove("job")
	if _, ok := g.Get("job"); ok {
		t.Error("expected schedule to be removed")
	}
}

func TestGroup_Names(t *testing.T) {
	g := NewGroup()
	for _, name := range []string{"c", "a", "b"} {
		_ = g.Add(name, mustSchedule(t, "* * * * *", "UTC"))
	}
	names := g.Names()
	want := []string{"a", "b", "c"}
	for i, n := range names {
		if n != want[i] {
			t.Errorf("Names()[%d] = %q, want %q", i, n, want[i])
		}
	}
}

func TestGroup_NextAll_Order(t *testing.T) {
	g := NewGroup()
	// every hour at minute 0
	_ = g.Add("hourly", mustSchedule(t, "0 * * * *", "UTC"))
	// every minute
	_ = g.Add("minutely", mustSchedule(t, "* * * * *", "UTC"))

	from := time.Date(2024, 6, 1, 12, 0, 30, 0, time.UTC)
	entries := g.NextAll(from)

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	// minutely fires sooner
	if entries[0].Name != "minutely" {
		t.Errorf("expected minutely first, got %q", entries[0].Name)
	}
	if !entries[0].Next.Before(entries[1].Next) {
		t.Error("expected entries sorted by ascending next-run time")
	}
}
