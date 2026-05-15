package scheduler

import (
	"testing"
	"time"
)

func TestLastRun_EveryMinute(t *testing.T) {
	s, err := New("* * * * *", "UTC")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	now := time.Date(2024, 3, 15, 10, 30, 45, 0, time.UTC)
	last := s.LastRun(now)

	expected := time.Date(2024, 3, 15, 10, 30, 0, 0, time.UTC)
	if !last.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, last)
	}
}

func TestLastRun_SpecificTime(t *testing.T) {
	s, err := New("30 9 * * *", "UTC")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	now := time.Date(2024, 3, 15, 12, 0, 0, 0, time.UTC)
	last := s.LastRun(now)

	expected := time.Date(2024, 3, 15, 9, 30, 0, 0, time.UTC)
	if !last.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, last)
	}
}

func TestLastRun_BeforeFirstRun(t *testing.T) {
	s, err := New("0 12 * * *", "UTC")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	now := time.Date(2024, 3, 15, 8, 0, 0, 0, time.UTC)
	last := s.LastRun(now)

	expected := time.Date(2024, 3, 14, 12, 0, 0, 0, time.UTC)
	if !last.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, last)
	}
}
