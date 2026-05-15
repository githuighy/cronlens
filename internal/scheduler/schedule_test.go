package scheduler_test

import (
	"testing"
	"time"

	"github.com/cronlens/cronlens/internal/scheduler"
)

func TestNew_ValidExpression(t *testing.T) {
	s, err := scheduler.New("*/5 * * * *", "UTC")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil schedule")
	}
}

func TestNew_InvalidTimezone(t *testing.T) {
	_, err := scheduler.New("*/5 * * * *", "Not/AZone")
	if err == nil {
		t.Fatal("expected error for invalid timezone")
	}
}

func TestNew_InvalidExpression(t *testing.T) {
	_, err := scheduler.New("bad expr", "UTC")
	if err == nil {
		t.Fatal("expected error for invalid expression")
	}
}

func TestSchedule_NextRun(t *testing.T) {
	s, err := scheduler.New("0 9 * * 1", "UTC") // Every Monday at 09:00
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Use a known Monday
	now := time.Date(2024, 1, 8, 8, 0, 0, 0, time.UTC) // Monday 08:00
	next, err := s.NextRun(now)
	if err != nil {
		t.Fatalf("NextRun error: %v", err)
	}
	if next.Weekday() != time.Monday {
		t.Errorf("expected Monday, got %v", next.Weekday())
	}
	if next.Hour() != 9 || next.Minute() != 0 {
		t.Errorf("expected 09:00, got %02d:%02d", next.Hour(), next.Minute())
	}
}

func TestSchedule_NextN(t *testing.T) {
	s, err := scheduler.New("*/15 * * * *", "UTC")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	now := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	runs, err := s.NextN(now, 4)
	if err != nil {
		t.Fatalf("NextN error: %v", err)
	}
	if len(runs) != 4 {
		t.Fatalf("expected 4 runs, got %d", len(runs))
	}
	for i := 1; i < len(runs); i++ {
		if !runs[i].After(runs[i-1]) {
			t.Errorf("runs not in ascending order at index %d", i)
		}
	}
}

func TestSchedule_TimeSinceLast(t *testing.T) {
	s, err := scheduler.New("* * * * *", "UTC") // every minute
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	now := time.Date(2024, 3, 15, 10, 30, 45, 0, time.UTC)
	dur, last, err := s.TimeSinceLast(now)
	if err != nil {
		t.Fatalf("TimeSinceLast error: %v", err)
	}
	if dur <= 0 {
		t.Errorf("expected positive duration, got %v", dur)
	}
	if last.IsZero() {
		t.Error("expected non-zero last run time")
	}
}
