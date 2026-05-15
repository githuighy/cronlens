package scheduler

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestOverlapGuard_Allow(t *testing.T) {
	guard := NewOverlapGuard(OverlapAllow)
	var count int32
	fn := guard.Wrap(func(ctx interface{ Done() <-chan struct{} }, _ time.Time) error {
		atomic.AddInt32(&count, 1)
		time.Sleep(20 * time.Millisecond)
		return nil
	})

	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = fn(context.Background(), time.Now())
		}()
	}
	wg.Wait()
	if atomic.LoadInt32(&count) != 3 {
		t.Errorf("expected 3 executions, got %d", count)
	}
}

func TestOverlapGuard_Skip(t *testing.T) {
	guard := NewOverlapGuard(OverlapSkip)
	var executed int32
	var skipped int32

	ready := make(chan struct{})
	fn := guard.Wrap(func(ctx interface{ Done() <-chan struct{} }, _ time.Time) error {
		close(ready)
		time.Sleep(40 * time.Millisecond)
		atomic.AddInt32(&executed, 1)
		return nil
	})

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = fn(context.Background(), time.Now())
	}()

	<-ready
	err := fn(context.Background(), time.Now())
	if err != ErrSkipped {
		t.Errorf("expected ErrSkipped, got %v", err)
	}
	atomic.AddInt32(&skipped, 1)

	wg.Wait()
	if atomic.LoadInt32(&executed) != 1 {
		t.Errorf("expected 1 execution, got %d", executed)
	}
	if atomic.LoadInt32(&skipped) != 1 {
		t.Errorf("expected 1 skip, got %d", skipped)
	}
}

func TestOverlapGuard_Queue(t *testing.T) {
	guard := NewOverlapGuard(OverlapQueue)
	var order []int
	var mu sync.Mutex

	fn := guard.Wrap(func(ctx interface{ Done() <-chan struct{} }, _ time.Time) error {
		time.Sleep(10 * time.Millisecond)
		return nil
	})

	for i := 0; i < 2; i++ {
		i := i
		go func() {
			_ = fn(context.Background(), time.Now())
			mu.Lock()
			order = append(order, i)
			mu.Unlock()
		}()
	}
	time.Sleep(60 * time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	if len(order) != 2 {
		t.Errorf("expected 2 completions, got %d", len(order))
	}
}
