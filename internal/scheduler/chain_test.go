package scheduler

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestJobChain_Empty(t *testing.T) {
	c := NewJobChain(0)
	if err := c.Run(context.Background()); err != nil {
		t.Fatalf("empty chain should not error, got: %v", err)
	}
}

func TestJobChain_AllSucceed(t *testing.T) {
	var order []string
	c := NewJobChain(0).
		Add("a", func(_ context.Context) error { order = append(order, "a"); return nil }).
		Add("b", func(_ context.Context) error { order = append(order, "b"); return nil })

	if err := c.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(order) != 2 || order[0] != "a" || order[1] != "b" {
		t.Errorf("unexpected execution order: %v", order)
	}
}

func TestJobChain_StopsOnFirstFailure(t *testing.T) {
	var ran []string
	errBoom := errors.New("boom")
	c := NewJobChain(0).
		Add("step1", func(_ context.Context) error { ran = append(ran, "step1"); return nil }).
		Add("step2", func(_ context.Context) error { ran = append(ran, "step2"); return errBoom }).
		Add("step3", func(_ context.Context) error { ran = append(ran, "step3"); return nil })

	err := c.Run(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "step2") {
		t.Errorf("error should mention failing step, got: %v", err)
	}
	if !errors.Is(err, errBoom) {
		t.Errorf("error should wrap original error")
	}
	if len(ran) != 2 {
		t.Errorf("expected 2 steps to run, got %d", len(ran))
	}
}

func TestJobChain_TimeoutEnforced(t *testing.T) {
	c := NewJobChain(20 * time.Millisecond).
		Add("slow", func(ctx context.Context) error {
			select {
			case <-time.After(500 * time.Millisecond):
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		})

	err := c.Run(context.Background())
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

func TestJobChain_NamesAndLen(t *testing.T) {
	c := NewJobChain(0).
		Add("x", func(_ context.Context) error { return nil }).
		Add("y", func(_ context.Context) error { return nil }).
		Add("z", func(_ context.Context) error { return nil })

	if c.Len() != 3 {
		t.Errorf("expected Len 3, got %d", c.Len())
	}
	names := c.Names()
	expected := []string{"x", "y", "z"}
	for i, n := range expected {
		if names[i] != n {
			t.Errorf("expected name %q at index %d, got %q", n, i, names[i])
		}
	}
}
