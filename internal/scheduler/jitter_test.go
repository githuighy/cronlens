package scheduler

import (
	"context"
	"testing"
	"time"
)

func TestJitter_DefaultDelay(t *testing.T) {
	j := NewJitter(100 * time.Millisecond)
	for i := 0; i < 20; i++ {
		d := j.Delay("any")
		if d < 0 || d >= 100*time.Millisecond {
			t.Fatalf("delay %v out of range [0, 100ms)", d)
		}
	}
}

func TestJitter_PerNameOverride(t *testing.T) {
	j := NewJitter(1 * time.Second)
	j.Set("fast", 10*time.Millisecond)

	for i := 0; i < 20; i++ {
		d := j.Delay("fast")
		if d >= 10*time.Millisecond {
			t.Fatalf("expected delay < 10ms, got %v", d)
		}
	}
}

func TestJitter_DeleteReverts(t *testing.T) {
	j := NewJitter(500 * time.Millisecond)
	j.Set("job", 1*time.Millisecond)
	j.Delete("job")

	// After deletion the default (500ms) applies; just check it is within range.
	d := j.Delay("job")
	if d < 0 || d >= 500*time.Millisecond {
		t.Fatalf("delay %v out of range after delete", d)
	}
}

func TestJitter_ZeroMaxReturnsZero(t *testing.T) {
	j := NewJitter(0)
	if d := j.Delay("x"); d != 0 {
		t.Fatalf("expected 0 delay, got %v", d)
	}
}

func TestWithJitter_RunsJob(t *testing.T) {
	j := NewJitter(0) // zero jitter so test is fast
	ran := false
	wrapped := WithJitter(j, "test", func(ctx context.Context) error {
		ran = true
		return nil
	})
	if err := wrapped(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ran {
		t.Fatal("job was not executed")
	}
}

func TestWithJitter_RespectsContextCancellation(t *testing.T) {
	j := NewJitter(10 * time.Second) // large delay
	j.Set("slow", 10*time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	ran := false
	wrapped := WithJitter(j, "slow", func(ctx context.Context) error {
		ran = true
		return nil
	})

	err := wrapped(ctx)
	if err == nil {
		t.Fatal("expected context error, got nil")
	}
	if ran {
		t.Fatal("job should not have run after context cancellation")
	}
}
