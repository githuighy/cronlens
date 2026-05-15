package scheduler

import (
	"bytes"
	"context"
	"errors"
	"log"
	"testing"
	"time"
)

func TestChain_ExecutesInOrder(t *testing.T) {
	var order []int
	m1 := func(next JobFunc) JobFunc {
		return func(ctx context.Context) error {
			order = append(order, 1)
			err := next(ctx)
			order = append(order, 4)
			return err
		}
	}
	m2 := func(next JobFunc) JobFunc {
		return func(ctx context.Context) error {
			order = append(order, 2)
			err := next(ctx)
			order = append(order, 3)
			return err
		}
	}
	job := Chain(func(ctx context.Context) error { return nil }, m1, m2)
	if err := job(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []int{1, 2, 3, 4}
	for i, v := range expected {
		if order[i] != v {
			t.Errorf("order[%d] = %d, want %d", i, order[i], v)
		}
	}
}

func TestWithLogging_LogsCompletion(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	job := Chain(func(ctx context.Context) error { return nil }, WithLogging(logger))
	if err := job(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("completed")) {
		t.Errorf("expected log to contain 'completed', got: %s", buf.String())
	}
}

func TestWithLogging_LogsError(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	wantErr := errors.New("boom")
	job := Chain(func(ctx context.Context) error { return wantErr }, WithLogging(logger))
	if err := job(context.Background()); !errors.Is(err, wantErr) {
		t.Fatalf("expected wantErr, got: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("failed")) {
		t.Errorf("expected log to contain 'failed', got: %s", buf.String())
	}
}

func TestWithTimeout_CancelsContext(t *testing.T) {
	job := Chain(
		func(ctx context.Context) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(500 * time.Millisecond):
				return nil
			}
		},
		WithTimeout(10*time.Millisecond),
	)
	err := job(context.Background())
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected DeadlineExceeded, got: %v", err)
	}
}

func TestWithRecover_CatchesPanic(t *testing.T) {
	job := Chain(
		func(ctx context.Context) error { panic("unexpected panic") },
		WithRecover(),
	)
	err := job(context.Background())
	if err == nil {
		t.Fatal("expected error from panic recovery, got nil")
	}
}
