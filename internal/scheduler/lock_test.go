package scheduler

import (
	"testing"
	"time"
)

func TestLock_AcquireAndRelease(t *testing.T) {
	l := NewLock(0)

	if err := l.Acquire("job-a"); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if !l.IsHeld("job-a") {
		t.Fatal("expected lock to be held after Acquire")
	}

	l.Release("job-a")
	if l.IsHeld("job-a") {
		t.Fatal("expected lock to be released")
	}
}

func TestLock_DoubleAcquireFails(t *testing.T) {
	l := NewLock(0)

	if err := l.Acquire("job-b"); err != nil {
		t.Fatalf("first acquire failed: %v", err)
	}
	if err := l.Acquire("job-b"); err != ErrLockTimeout {
		t.Fatalf("expected ErrLockTimeout, got %v", err)
	}
	l.Release("job-b")
}

func TestLock_TTLExpiry(t *testing.T) {
	l := NewLock(20 * time.Millisecond)

	if err := l.Acquire("job-c"); err != nil {
		t.Fatalf("first acquire failed: %v", err)
	}

	// Lock should still be held immediately.
	if !l.IsHeld("job-c") {
		t.Fatal("expected lock to be held before TTL")
	}

	time.Sleep(40 * time.Millisecond)

	// After TTL the lock should have expired.
	if l.IsHeld("job-c") {
		t.Fatal("expected lock to have expired")
	}

	// Re-acquire should succeed after expiry.
	if err := l.Acquire("job-c"); err != nil {
		t.Fatalf("re-acquire after TTL failed: %v", err)
	}
	l.Release("job-c")
}

func TestLock_IndependentNames(t *testing.T) {
	l := NewLock(0)

	if err := l.Acquire("x"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := l.Acquire("y"); err != nil {
		t.Fatalf("different name should not be blocked: %v", err)
	}

	l.Release("x")
	l.Release("y")
}

func TestLock_ReleaseUnheld(t *testing.T) {
	l := NewLock(0)
	// Should not panic.
	l.Release("nonexistent")
}
