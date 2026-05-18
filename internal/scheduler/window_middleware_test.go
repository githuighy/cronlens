package scheduler

import (
	"context"
	"errors"
	"testing"
	"time"
)

// TestWithWindow_PropagatesError ensures errors from the inner job are
// returned when execution falls inside the allowed window.
func TestWithWindow_PropagatesError(t *testing.T) {
	ws := NewWindowStore()
	_ = ws.Set("open", 0, 24*time.Hour)
	want := errors.New("job failed")
	job := WithWindow(ws, "open", func(ctx context.Context) error {
		return want
	})
	if got := job(context.Background()); got != want {
		t.Errorf("expected %v, got %v", want, got)
	}
}

// TestWithWindow_MissingWindowSkips verifies that a missing window name
// causes the job to be silently skipped (nil returned).
func TestWithWindow_MissingWindowSkips(t *testing.T) {
	ws := NewWindowStore()
	called := false
	job := WithWindow(ws, "missing", func(ctx context.Context) error {
		called = true
		return nil
	})
	if err := job(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("job should not have been called for missing window")
	}
}

// TestWithWindow_ComposesWithChain verifies WithWindow composes correctly
// with other middleware such as Chain.
func TestWithWindow_ComposesWithChain(t *testing.T) {
	ws := NewWindowStore()
	_ = ws.Set("open", 0, 24*time.Hour)

	var order []string
	a := func(ctx context.Context) error { order = append(order, "a"); return nil }
	b := func(ctx context.Context) error { order = append(order, "b"); return nil }

	chained := Chain(a, b)
	windowed := WithWindow(ws, "open", chained)

	if err := windowed(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(order) != 2 || order[0] != "a" || order[1] != "b" {
		t.Errorf("unexpected execution order: %v", order)
	}
}

// TestWindowStore_OverwriteWindow ensures that setting a window with an
// existing name replaces the previous value.
func TestWindowStore_OverwriteWindow(t *testing.T) {
	ws := NewWindowStore()
	_ = ws.Set("slot", 8*time.Hour, 12*time.Hour)
	_ = ws.Set("slot", 13*time.Hour, 18*time.Hour)

	w, ok := ws.Get("slot")
	if !ok {
		t.Fatal("expected window to exist")
	}
	if w.Start != 13*time.Hour || w.End != 18*time.Hour {
		t.Errorf("expected overwritten window, got %+v", w)
	}

	// Old range should no longer match
	morning := time.Date(2024, 6, 1, 9, 0, 0, 0, time.UTC)
	if ws.InWindow("slot", morning) {
		t.Error("expected morning time to be outside overwritten window")
	}
}
