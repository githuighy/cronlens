package scheduler

import (
	"testing"
	"time"
)

var (
	t0 = time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC)
	t1 = time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	t2 = time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC)
)

func TestHistory_RecordAndAll(t *testing.T) {
	h := NewHistory(5)

	h.Record(t0, t0.Add(time.Millisecond), "job-a")
	h.Record(t1, t1.Add(time.Millisecond), "job-b")

	records := h.All()
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	if records[0].Label != "job-a" {
		t.Errorf("expected first label 'job-a', got %q", records[0].Label)
	}
	if records[1].Label != "job-b" {
		t.Errorf("expected second label 'job-b', got %q", records[1].Label)
	}
}

func TestHistory_CapEviction(t *testing.T) {
	h := NewHistory(2)

	h.Record(t0, t0, "first")
	h.Record(t1, t1, "second")
	h.Record(t2, t2, "third")

	if h.Len() != 2 {
		t.Fatalf("expected cap of 2, got %d", h.Len())
	}

	records := h.All()
	if records[0].Label != "second" {
		t.Errorf("expected 'second' after eviction, got %q", records[0].Label)
	}
	if records[1].Label != "third" {
		t.Errorf("expected 'third' as latest, got %q", records[1].Label)
	}
}

func TestHistory_Last(t *testing.T) {
	h := NewHistory(5)

	_, ok := h.Last()
	if ok {
		t.Error("expected no record on empty history")
	}

	h.Record(t0, t0, "only")
	rec, ok := h.Last()
	if !ok {
		t.Fatal("expected a record after insert")
	}
	if rec.Label != "only" {
		t.Errorf("expected label 'only', got %q", rec.Label)
	}
}

func TestHistory_DefaultCap(t *testing.T) {
	h := NewHistory(0)
	for i := 0; i < 15; i++ {
		h.Record(t0, t0, "x")
	}
	if h.Len() != 10 {
		t.Errorf("expected default cap 10, got %d", h.Len())
	}
}
