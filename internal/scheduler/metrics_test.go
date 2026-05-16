package scheduler

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestMetrics_InitialState(t *testing.T) {
	m := NewMetrics("test")
	snap := m.Snapshot()

	if snap.TotalRuns != 0 || snap.SuccessRuns != 0 || snap.FailureRuns != 0 {
		t.Fatalf("expected zero counts, got %+v", snap)
	}
	if snap.SuccessRate() != 0 {
		t.Fatalf("expected 0 success rate on empty metrics")
	}
}

func TestMetrics_RecordSuccess(t *testing.T) {
	m := NewMetrics("test")
	m.Record(10*time.Millisecond, nil)

	snap := m.Snapshot()
	if snap.TotalRuns != 1 || snap.SuccessRuns != 1 || snap.FailureRuns != 0 {
		t.Fatalf("unexpected counts: %+v", snap)
	}
	if snap.AvgLatency != 10*time.Millisecond {
		t.Fatalf("expected 10ms avg latency, got %v", snap.AvgLatency)
	}
	if snap.LastRun.IsZero() || snap.LastSuccess.IsZero() {
		t.Fatal("expected non-zero LastRun and LastSuccess")
	}
	if !snap.LastFailure.IsZero() {
		t.Fatal("expected zero LastFailure")
	}
}

func TestMetrics_RecordFailure(t *testing.T) {
	m := NewMetrics("test")
	sentinel := errors.New("boom")
	m.Record(5*time.Millisecond, sentinel)

	snap := m.Snapshot()
	if snap.FailureRuns != 1 || snap.SuccessRuns != 0 {
		t.Fatalf("unexpected counts: %+v", snap)
	}
	if snap.LastError == nil || snap.LastError.Error() != "boom" {
		t.Fatalf("expected last error 'boom', got %v", snap.LastError)
	}
	if snap.SuccessRate() != 0 {
		t.Fatalf("expected 0 success rate")
	}
}

func TestMetrics_SuccessRate(t *testing.T) {
	m := NewMetrics("test")
	m.Record(1*time.Millisecond, nil)
	m.Record(1*time.Millisecond, nil)
	m.Record(1*time.Millisecond, errors.New("err"))

	snap := m.Snapshot()
	got := snap.SuccessRate()
	want := 2.0 / 3.0
	if got < want-0.001 || got > want+0.001 {
		t.Fatalf("expected success rate ~%.3f, got %.3f", want, got)
	}
}

func TestMetrics_Reset(t *testing.T) {
	m := NewMetrics("test")
	m.Record(1*time.Millisecond, nil)
	m.Reset()

	snap := m.Snapshot()
	if snap.TotalRuns != 0 || !snap.LastRun.IsZero() {
		t.Fatalf("expected reset state, got %+v", snap)
	}
}

func TestWithMetrics_RecordsExecution(t *testing.T) {
	m := NewMetrics("wrapped")
	called := false
	job := func(ctx context.Context) error {
		called = true
		return nil
	}

	wrapped := WithMetrics(m, job)
	if err := wrapped(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected job to be called")
	}
	if m.Snapshot().TotalRuns != 1 {
		t.Fatal("expected 1 recorded run")
	}
}

func TestWithMetrics_NilMetricsPassthrough(t *testing.T) {
	called := false
	job := func(ctx context.Context) error { called = true; return nil }
	wrapped := WithMetrics(nil, job)
	_ = wrapped(context.Background())
	if !called {
		t.Fatal("expected job to be called even with nil metrics")
	}
}
