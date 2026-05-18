package scheduler

import (
	"context"
	"errors"
	"testing"
)

func TestDependencyStore_ReadyWithNoDeps(t *testing.T) {
	store := NewDependencyStore()
	if !store.Ready("standalone") {
		t.Fatal("job with no declared deps should always be ready")
	}
}

func TestDependencyStore_NotReadyUntilDepDone(t *testing.T) {
	store := NewDependencyStore()
	if err := store.Declare("jobB", "jobA"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if store.Ready("jobB") {
		t.Fatal("jobB should not be ready before jobA is done")
	}
	store.MarkDone("jobA")
	if !store.Ready("jobB") {
		t.Fatal("jobB should be ready after jobA is done")
	}
}

func TestDependencyStore_ResetClearsState(t *testing.T) {
	store := NewDependencyStore()
	store.MarkDone("jobA")
	store.Reset("jobA")
	if err := store.Declare("jobB", "jobA"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if store.Ready("jobB") {
		t.Fatal("jobB should not be ready after jobA was reset")
	}
}

func TestDependencyStore_CircularSelf(t *testing.T) {
	store := NewDependencyStore()
	err := store.Declare("jobA", "jobA")
	if !errors.Is(err, ErrCircularDependency) {
		t.Fatalf("expected ErrCircularDependency, got %v", err)
	}
}

func TestDependencyStore_CircularIndirect(t *testing.T) {
	store := NewDependencyStore()
	_ = store.Declare("jobB", "jobA")
	err := store.Declare("jobA", "jobB")
	if !errors.Is(err, ErrCircularDependency) {
		t.Fatalf("expected ErrCircularDependency, got %v", err)
	}
}

func TestWithDependency_BlocksWhenNotReady(t *testing.T) {
	store := NewDependencyStore()
	_ = store.Declare("jobB", "jobA")
	called := false
	job := WithDependency(store, "jobB", func(_ context.Context) error {
		called = true
		return nil
	})
	err := job(context.Background())
	if !errors.Is(err, ErrDependencyNotMet) {
		t.Fatalf("expected ErrDependencyNotMet, got %v", err)
	}
	if called {
		t.Fatal("job should not have been called")
	}
}

func TestWithDependency_RunsAndMarksDoneWhenReady(t *testing.T) {
	store := NewDependencyStore()
	_ = store.Declare("jobB", "jobA")
	store.MarkDone("jobA")
	called := false
	job := WithDependency(store, "jobB", func(_ context.Context) error {
		called = true
		return nil
	})
	if err := job(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("job should have been called")
	}
	if !store.Ready("jobC") {
		// jobC has no deps, always ready — sanity check
	}
	if !store.done["jobB"] {
		t.Fatal("jobB should be marked done after successful run")
	}
}
