package scheduler

import (
	"testing"
	"time"
)

func TestSummarize(t *testing.T) {
	s, err := New("0 9 * * 1-5", "UTC")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Wednesday 2024-03-13 at 10:00 UTC
	now := time.Date(2024, 3, 13, 10, 0, 0, 0, time.UTC)
	summary := s.Summarize(now)

	if summary.Expression != "0 9 * * 1-5" {
		t.Errorf("unexpected expression: %q", summary.Expression)
	}
	if summary.Timezone != "UTC" {
		t.Errorf("unexpected timezone: %q", summary.Timezone)
	}
	if summary.Description == "" {
		t.Error("expected non-empty description")
	}
	if summary.LastRun.IsZero() {
		t.Error("expected non-zero last run")
	}
	if summary.NextRun.IsZero() {
		t.Error("expected non-zero next run")
	}
	if summary.TimeUntil <= 0 {
		t.Errorf("expected positive TimeUntil, got %v", summary.TimeUntil)
	}
	if summary.TimeSince <= 0 {
		t.Errorf("expected positive TimeSince, got %v", summary.TimeSince)
	}
}

func TestSummarize_Timezone(t *testing.T) {
	s, err := New("0 12 * * *", "America/New_York")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	now := time.Date(2024, 3, 15, 15, 0, 0, 0, time.UTC)
	summary := s.Summarize(now)

	if summary.Timezone != "America/New_York" {
		t.Errorf("unexpected timezone: %q", summary.Timezone)
	}
	if summary.NextRun.IsZero() {
		t.Error("expected non-zero next run")
	}
}
