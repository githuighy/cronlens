package scheduler

import (
	"context"
	"errors"
	"testing"
	"time"
)

var errJob = errors.New("job error")

func TestCircuitBreaker_InitiallyClosed(t *testing.T) {
	cb := NewCircuitBreaker(3, time.Second)
	if cb.State() != CircuitClosed {
		t.Fatalf("expected Closed, got %v", cb.State())
	}
	if err := cb.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestCircuitBreaker_OpensAfterThreshold(t *testing.T) {
	cb := NewCircuitBreaker(3, time.Minute)
	for i := 0; i < 3; i++ {
		cb.RecordFailure()
	}
	if cb.State() != CircuitOpen {
		t.Fatalf("expected Open after 3 failures, got %v", cb.State())
	}
	if err := cb.Allow(); !errors.Is(err, ErrCircuitOpen) {
		t.Fatalf("expected ErrCircuitOpen, got %v", err)
	}
}

func TestCircuitBreaker_ClosesOnSuccess(t *testing.T) {
	cb := NewCircuitBreaker(2, time.Minute)
	cb.RecordFailure()
	cb.RecordFailure()
	cb.RecordSuccess()
	if cb.State() != CircuitClosed {
		t.Fatalf("expected Closed after success, got %v", cb.State())
	}
}

func TestCircuitBreaker_HalfOpenAfterTimeout(t *testing.T) {
	cb := NewCircuitBreaker(1, 10*time.Millisecond)
	cb.RecordFailure() // opens circuit
	time.Sleep(20 * time.Millisecond)
	if cb.State() != CircuitHalfOpen {
		t.Fatalf("expected HalfOpen after timeout, got %v", cb.State())
	}
	if err := cb.Allow(); err != nil {
		t.Fatalf("HalfOpen should allow probe, got %v", err)
	}
}

func TestWithCircuitBreaker_BlocksWhenOpen(t *testing.T) {
	cb := NewCircuitBreaker(1, time.Minute)
	cb.RecordFailure()

	called := false
	job := WithCircuitBreaker(cb, func(ctx context.Context) error {
		called = true
		return nil
	})

	err := job(context.Background())
	if !errors.Is(err, ErrCircuitOpen) {
		t.Fatalf("expected ErrCircuitOpen, got %v", err)
	}
	if called {
		t.Fatal("job should not have been called when circuit is open")
	}
}

func TestWithCircuitBreaker_RecordsFailure(t *testing.T) {
	cb := NewCircuitBreaker(3, time.Minute)
	job := WithCircuitBreaker(cb, func(ctx context.Context) error {
		return errJob
	})

	_ = job(context.Background())
	_ = job(context.Background())
	_ = job(context.Background())

	if cb.State() != CircuitOpen {
		t.Fatalf("expected Open after 3 job failures, got %v", cb.State())
	}
}

func TestNewCircuitBreaker_Defaults(t *testing.T) {
	cb := NewCircuitBreaker(0, 0)
	if cb.maxFailures != 3 {
		t.Fatalf("expected default maxFailures=3, got %d", cb.maxFailures)
	}
	if cb.resetTimeout != 30*time.Second {
		t.Fatalf("expected default resetTimeout=30s, got %v", cb.resetTimeout)
	}
}
