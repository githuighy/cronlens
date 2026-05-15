package scheduler

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestRateLimit_FirstCallAlwaysPasses(t *testing.T) {
	rl := NewRateLimit(3, time.Minute)
	now := time.Now()
	if !rl.Allow("job", now) {
		t.Fatal("expected first call to be allowed")
	}
}

func TestRateLimit_BlocksAfterMax(t *testing.T) {
	rl := NewRateLimit(2, time.Minute)
	now := time.Now()
	rl.Allow("job", now)
	rl.Allow("job", now.Add(time.Second))
	if rl.Allow("job", now.Add(2*time.Second)) {
		t.Fatal("expected third call to be blocked")
	}
}

func TestRateLimit_AllowsAfterWindowExpires(t *testing.T) {
	rl := NewRateLimit(1, 10*time.Millisecond)
	now := time.Now()
	rl.Allow("job", now)
	// Advance beyond the window.
	future := now.Add(20 * time.Millisecond)
	if !rl.Allow("job", future) {
		t.Fatal("expected call to be allowed after window expiry")
	}
}

func TestRateLimit_IndependentNames(t *testing.T) {
	rl := NewRateLimit(1, time.Minute)
	now := time.Now()
	rl.Allow("a", now)
	if !rl.Allow("b", now) {
		t.Fatal("expected independent name to be allowed")
	}
}

func TestRateLimit_Reset(t *testing.T) {
	rl := NewRateLimit(1, time.Minute)
	now := time.Now()
	rl.Allow("job", now)
	rl.Reset("job")
	if !rl.Allow("job", now) {
		t.Fatal("expected call to be allowed after reset")
	}
}

func TestRateLimit_Remaining(t *testing.T) {
	rl := NewRateLimit(3, time.Minute)
	now := time.Now()
	if rem := rl.Remaining("job", now); rem != 3 {
		t.Fatalf("expected 3 remaining, got %d", rem)
	}
	rl.Allow("job", now)
	if rem := rl.Remaining("job", now); rem != 2 {
		t.Fatalf("expected 2 remaining, got %d", rem)
	}
}

func TestWithRateLimit_Middleware(t *testing.T) {
	rl := NewRateLimit(1, time.Minute)
	calls := 0
	base := JobFunc(func(_ context.Context) error {
		calls++
		return nil
	})
	guarded := Chain(base, WithRateLimit(rl, "job"))

	if err := guarded(context.Background()); err != nil {
		t.Fatalf("unexpected error on first call: %v", err)
	}
	err := guarded(context.Background())
	if err == nil {
		t.Fatal("expected rate limit error on second call")
	}
	expected := fmt.Sprintf("rate limit exceeded for %q", "job")
	if err.Error() != expected {
		t.Fatalf("unexpected error message: %q", err.Error())
	}
	if calls != 1 {
		t.Fatalf("expected 1 execution, got %d", calls)
	}
}
