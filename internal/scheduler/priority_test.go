package scheduler

import (
	"context"
	"errors"
	"testing"
)

func TestPriorityQueue_EnqueueDequeue(t *testing.T) {
	pq := NewPriorityQueue()

	order := []string{}
	makeJob := func(name string) func(ctx context.Context) error {
		return func(_ context.Context) error {
			order = append(order, name)
			return nil
		}
	}

	pq.Enqueue("low", PriorityLow, makeJob("low"))
	pq.Enqueue("high", PriorityHigh, makeJob("high"))
	pq.Enqueue("normal", PriorityNormal, makeJob("normal"))

	if err := DrainQueue(context.Background(), pq); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"high", "normal", "low"}
	for i, name := range want {
		if order[i] != name {
			t.Errorf("position %d: got %q, want %q", i, order[i], name)
		}
	}
}

func TestPriorityQueue_EmptyDequeue(t *testing.T) {
	pq := NewPriorityQueue()
	_, _, ok := pq.Dequeue()
	if ok {
		t.Error("expected false on empty dequeue")
	}
}

func TestPriorityQueue_Len(t *testing.T) {
	pq := NewPriorityQueue()
	if pq.Len() != 0 {
		t.Errorf("expected 0, got %d", pq.Len())
	}
	noop := func(_ context.Context) error { return nil }
	pq.Enqueue("a", PriorityNormal, noop)
	pq.Enqueue("b", PriorityNormal, noop)
	if pq.Len() != 2 {
		t.Errorf("expected 2, got %d", pq.Len())
	}
}

func TestDrainQueue_StopsOnError(t *testing.T) {
	pq := NewPriorityQueue()
	executed := 0
	errBoom := errors.New("boom")

	pq.Enqueue("fail", PriorityHigh, func(_ context.Context) error { executed++; return errBoom })
	pq.Enqueue("ok", PriorityLow, func(_ context.Context) error { executed++; return nil })

	err := DrainQueue(context.Background(), pq)
	if !errors.Is(err, errBoom) {
		t.Fatalf("expected errBoom, got %v", err)
	}
	if executed != 1 {
		t.Errorf("expected 1 execution before stop, got %d", executed)
	}
	if pq.Len() != 1 {
		t.Errorf("expected 1 remaining job, got %d", pq.Len())
	}
}

func TestWithPriority_EnqueuesInsteadOfRunning(t *testing.T) {
	pq := NewPriorityQueue()
	ran := false
	wrapped := WithPriority(pq, "deferred", PriorityNormal, func(_ context.Context) error {
		ran = true
		return nil
	})

	// Calling the wrapper should enqueue, not run.
	if err := wrapped(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ran {
		t.Error("job should not have run yet")
	}
	if pq.Len() != 1 {
		t.Errorf("expected 1 queued job, got %d", pq.Len())
	}

	// Draining should execute it.
	if err := DrainQueue(context.Background(), pq); err != nil {
		t.Fatalf("drain error: %v", err)
	}
	if !ran {
		t.Error("job should have run after drain")
	}
}
