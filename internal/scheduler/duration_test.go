package scheduler

import (
	"testing"
	"time"
)

func TestTimeUntilNext(t *testing.T) {
	s, err := New("0 * * * *", "UTC")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	now := time.Date(2024, 3, 15, 10, 30, 0, 0, time.UTC)
	dur := s.TimeUntilNext(now)

	expected := 30 * time.Minute
	if dur != expected {
		t.Errorf("expected %v, got %v", expected, dur)
	}
}

func TestTimeSinceLast(t *testing.T) {
	s, err := New("0 * * * *", "UTC")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	now := time.Date(2024, 3, 15, 10, 30, 0, 0, time.UTC)
	dur := s.TimeSinceLast(now)

	expected := 30 * time.Minute
	if dur != expected {
		t.Errorf("expected %v, got %v", expected, dur)
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Duration
		expected string
	}{
		{"seconds", 45 * time.Second, "45s"},
		{"minutes", 5*time.Minute + 30*time.Second, "5m 30s"},
		{"hours", 2*time.Hour + 15*time.Minute, "2h 15m"},
		{"days", 26*time.Hour + 10*time.Minute, "1d 2h 10m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatDuration(tt.input)
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}
