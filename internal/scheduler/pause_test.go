package scheduler

import (
	"context"
	"errors"
	"testing"
)

func TestPauseStore_PauseAndResume(t *testing.T) {
	store := NewPauseStore()

	if err := store.Pause("job1", "maintenance"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !store.IsPaused("job1") {
		t.Fatal("expected job1 to be paused")
	}

	if err := store.Resume("job1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if store.IsPaused("job1") {
		t.Fatal("expected job1 to be resumed")
	}
}

func TestPauseStore_DoublePauseFails(t *testing.T) {
	store := NewPauseStore()
	_ = store.Pause("job1", "")

	err := store.Pause("job1", "again")
	if err == nil {
		t.Fatal("expected error on double pause")
	}
}

func TestPauseStore_ResumeNotPausedFails(t *testing.T) {
	store := NewPauseStore()

	err := store.Resume("unknown")
	if err == nil {
		t.Fatal("expected error resuming unknown schedule")
	}
}

func TestPauseStore_State(t *testing.T) {
	store := NewPauseStore()

	if s := store.State("missing"); s != nil {
		t.Fatal("expected nil state for unknown schedule")
	}

	_ = store.Pause("job2", "testing")
	s := store.State("job2")
	if s == nil {
		t.Fatal("expected non-nil state")
	}
	if !s.Paused {
		t.Error("expected state to be paused")
	}
	if s.Reason != "testing" {
		t.Errorf("expected reason %q, got %q", "testing", s.Reason)
	}
}

func TestWithPause_SkipsWhenPaused(t *testing.T) {
	store := NewPauseStore()
	_ = store.Pause("job3", "")

	called := false
	job := WithPause(store, "job3", func(ctx context.Context) error {
		called = true
		return nil
	})

	if err := job(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Fatal("expected job to be skipped while paused")
	}
}

func TestWithPause_RunsWhenNotPaused(t *testing.T) {
	store := NewPauseStore()
	sentinel := errors.New("ran")

	job := WithPause(store, "job4", func(ctx context.Context) error {
		return sentinel
	})

	err := job(context.Background())
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}
