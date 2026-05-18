package scheduler

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestWindowStore_SetAndGet(t *testing.T) {
	ws := NewWindowStore()
	if err := ws.Set("work", 9*time.Hour, 17*time.Hour); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	w, ok := ws.Get("work")
	if !ok {
		t.Fatal("expected window to exist")
	}
	if w.Start != 9*time.Hour || w.End != 17*time.Hour {
		t.Errorf("unexpected window values: %+v", w)
	}
}

func TestWindowStore_SetEmptyName(t *testing.T) {
	ws := NewWindowStore()
	if err := ws.Set("", time.Hour, 2*time.Hour); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestWindowStore_SetInvalidRange(t *testing.T) {
	ws := NewWindowStore()
	if err := ws.Set("bad", 10*time.Hour, 9*time.Hour); err == nil {
		t.Fatal("expected error when start >= end")
	}
}

func TestWindowStore_Delete(t *testing.T) {
	ws := NewWindowStore()
	_ = ws.Set("temp", time.Hour, 2*time.Hour)
	ws.Delete("temp")
	if _, ok := ws.Get("temp"); ok {
		t.Fatal("expected window to be deleted")
	}
}

func TestWindowStore_InWindow_Inside(t *testing.T) {
	ws := NewWindowStore()
	_ = ws.Set("work", 9*time.Hour, 17*time.Hour)
	// 10:30 on any day
	now := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	if !ws.InWindow("work", now) {
		t.Error("expected time to be inside window")
	}
}

func TestWindowStore_InWindow_Outside(t *testing.T) {
	ws := NewWindowStore()
	_ = ws.Set("work", 9*time.Hour, 17*time.Hour)
	now := time.Date(2024, 1, 15, 20, 0, 0, 0, time.UTC)
	if ws.InWindow("work", now) {
		t.Error("expected time to be outside window")
	}
}

func TestWindowStore_InWindow_Missing(t *testing.T) {
	ws := NewWindowStore()
	now := time.Now()
	if ws.InWindow("nonexistent", now) {
		t.Error("expected false for missing window")
	}
}

func TestWithWindow_RunsInsideWindow(t *testing.T) {
	ws := NewWindowStore()
	_ = ws.Set("always", 0, 24*time.Hour)
	called := false
	job := WithWindow(ws, "always", func(ctx context.Context) error {
		called = true
		return nil
	})
	if err := job(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected job to be called inside window")
	}
}

func TestWithWindow_SkipsOutsideWindow(t *testing.T) {
	ws := NewWindowStore()
	// A window that never covers current time (far future offset — won't match)
	_ = ws.Set("never", 0, time.Nanosecond)
	sentinel := errors.New("should not run")
	job := WithWindow(ws, "never", func(ctx context.Context) error {
		return sentinel
	})
	if err := job(context.Background()); err != nil {
		t.Errorf("expected nil, got: %v", err)
	}
}
