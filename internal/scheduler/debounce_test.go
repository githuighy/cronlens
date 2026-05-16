package scheduler

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestDebounce_SingleTrigger(t *testing.T) {
	d := NewDebounce(50 * time.Millisecond)

	var count int32
	d.Trigger("job", func(_ context.Context) error {
		atomic.AddInt32(&count, 1)
		return nil
	})

	time.Sleep(100 * time.Millisecond)
	if got := atomic.LoadInt32(&count); got != 1 {
		t.Errorf("expected 1 execution, got %d", got)
	}
}

func TestDebounce_RapidTriggerCollapses(t *testing.T) {
	d := NewDebounce(80 * time.Millisecond)

	var count int32
	fn := func(_ context.Context) error {
		atomic.AddInt32(&count, 1)
		return nil
	}

	for i := 0; i < 5; i++ {
		d.Trigger("job", fn)
		time.Sleep(20 * time.Millisecond)
	}

	time.Sleep(150 * time.Millisecond)
	if got := atomic.LoadInt32(&count); got != 1 {
		t.Errorf("expected 1 execution after debounce, got %d", got)
	}
}

func TestDebounce_Cancel(t *testing.T) {
	d := NewDebounce(100 * time.Millisecond)

	var count int32
	d.Trigger("job", func(_ context.Context) error {
		atomic.AddInt32(&count, 1)
		return nil
	})

	if !d.Pending("job") {
		t.Error("expected job to be pending")
	}

	d.Cancel("job")

	if d.Pending("job") {
		t.Error("expected job to no longer be pending after cancel")
	}

	time.Sleep(150 * time.Millisecond)
	if got := atomic.LoadInt32(&count); got != 0 {
		t.Errorf("expected 0 executions after cancel, got %d", got)
	}
}

func TestDebounce_IndependentNames(t *testing.T) {
	d := NewDebounce(50 * time.Millisecond)

	var countA, countB int32
	d.Trigger("a", func(_ context.Context) error { atomic.AddInt32(&countA, 1); return nil })
	d.Trigger("b", func(_ context.Context) error { atomic.AddInt32(&countB, 1); return nil })

	time.Sleep(120 * time.Millisecond)

	if got := atomic.LoadInt32(&countA); got != 1 {
		t.Errorf("expected countA=1, got %d", got)
	}
	if got := atomic.LoadInt32(&countB); got != 1 {
		t.Errorf("expected countB=1, got %d", got)
	}
}

func TestWithDebounce_Wrapper(t *testing.T) {
	d := NewDebounce(50 * time.Millisecond)

	var count int32
	fn := func(_ context.Context) error {
		atomic.AddInt32(&count, 1)
		return nil
	}

	wrapped := WithDebounce(d, "wrapped-job", fn)
	_ = wrapped(context.Background())
	_ = wrapped(context.Background())

	time.Sleep(120 * time.Millisecond)
	if got := atomic.LoadInt32(&count); got != 1 {
		t.Errorf("expected 1 execution via WithDebounce wrapper, got %d", got)
	}
}
