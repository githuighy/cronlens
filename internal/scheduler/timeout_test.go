package scheduler

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestTimeoutStore_DefaultApplied(t *testing.T) {
	store := NewTimeoutStore(100 * time.Millisecond)
	if got := store.Get("any"); got != 100*time.Millisecond {
		t.Fatalf("expected 100ms, got %v", got)
	}
}

func TestTimeoutStore_OverrideApplied(t *testing.T) {
	store := NewTimeoutStore(100 * time.Millisecond)
	if err := store.Set("job", 5*time.Second); err != nil {
		t.Fatal(err)
	}
	if got := store.Get("job"); got != 5*time.Second {
		t.Fatalf("expected 5s, got %v", got)
	}
}

func TestTimeoutStore_DeleteRevertsToDefault(t *testing.T) {
	store := NewTimeoutStore(50 * time.Millisecond)
	_ = store.Set("job", time.Minute)
	store.Delete("job")
	if got := store.Get("job"); got != 50*time.Millisecond {
		t.Fatalf("expected default 50ms after delete, got %v", got)
	}
}

func TestTimeoutStore_SetEmptyNameErrors(t *testing.T) {
	store := NewTimeoutStore(time.Second)
	if err := store.Set("", time.Second); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestTimeoutStore_ZeroDurationRemovesOverride(t *testing.T) {
	store := NewTimeoutStore(10 * time.Millisecond)
	_ = store.Set("job", time.Hour)
	_ = store.Set("job", 0) // should remove override
	if got := store.Get("job"); got != 10*time.Millisecond {
		t.Fatalf("expected default after zero set, got %v", got)
	}
}

func TestWithTimeout_CancelsSlowJob(t *testing.T) {
	store := NewTimeoutStore(50 * time.Millisecond)

	slow := func(ctx context.Context) error {
		select {
		case <-time.After(5 * time.Second):
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	wrapped := store.WithTimeout("slow", slow)
	err := wrapped(context.Background())
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
}

func TestWithTimeout_FastJobSucceeds(t *testing.T) {
	store := NewTimeoutStore(500 * time.Millisecond)

	fast := func(ctx context.Context) error { return nil }
	wrapped := store.WithTimeout("fast", fast)
	if err := wrapped(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWithTimeout_ZeroTimeoutNoDeadline(t *testing.T) {
	store := NewTimeoutStore(0) // no default timeout

	called := false
	job := func(ctx context.Context) error {
		called = true
		if _, ok := ctx.Deadline(); ok {
			t.Error("expected no deadline on context")
		}
		return nil
	}
	wrapped := store.WithTimeout("nolimit", job)
	_ = wrapped(context.Background())
	if !called {
		t.Fatal("job was not called")
	}
}
