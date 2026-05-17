package scheduler

import (
	"bytes"
	"context"
	"errors"
	"log"
	"strings"
	"testing"
)

func TestWithPriorityLogging_AllSucceed(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	pq := NewPriorityQueue()
	noop := func(_ context.Context) error { return nil }
	pq.Enqueue("alpha", PriorityHigh, noop)
	pq.Enqueue("beta", PriorityLow, noop)

	if err := WithPriorityLogging(context.Background(), pq, logger); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "alpha") {
		t.Error("expected log output for 'alpha'")
	}
	if !strings.Contains(out, "beta") {
		t.Error("expected log output for 'beta'")
	}
	if !strings.Contains(out, "completed") {
		t.Error("expected 'completed' in log output")
	}
}

func TestWithPriorityLogging_StopsOnError(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	errFail := errors.New("fail")

	pq := NewPriorityQueue()
	pq.Enqueue("bad", PriorityHigh, func(_ context.Context) error { return errFail })
	pq.Enqueue("good", PriorityLow, func(_ context.Context) error { return nil })

	err := WithPriorityLogging(context.Background(), pq, logger)
	if err == nil {
		t.Fatal("expected an error")
	}
	if !errors.Is(err, errFail) {
		t.Errorf("expected errFail wrapped, got %v", err)
	}
	if !strings.Contains(err.Error(), "bad") {
		t.Errorf("error should mention job name 'bad', got: %v", err)
	}
	if pq.Len() != 1 {
		t.Errorf("expected 1 unexecuted job remaining, got %d", pq.Len())
	}
}

func TestWithPriorityLogging_EmptyQueue(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	pq := NewPriorityQueue()

	if err := WithPriorityLogging(context.Background(), pq, logger); err != nil {
		t.Fatalf("unexpected error on empty queue: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no log output for empty queue, got: %s", buf.String())
	}
}
