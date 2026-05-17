package scheduler

import (
	"context"
	"errors"
	"testing"
)

func TestPipeline_Empty(t *testing.T) {
	p := NewPipeline()
	if err := p.Run(context.Background()); err != nil {
		t.Fatalf("expected nil error for empty pipeline, got %v", err)
	}
}

func TestPipeline_AllSucceed(t *testing.T) {
	order := []string{}
	p := NewPipeline().
		Add("a", func(ctx context.Context) error { order = append(order, "a"); return nil }).
		Add("b", func(ctx context.Context) error { order = append(order, "b"); return nil }).
		Add("c", func(ctx context.Context) error { order = append(order, "c"); return nil })

	if err := p.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(order) != 3 || order[0] != "a" || order[1] != "b" || order[2] != "c" {
		t.Errorf("unexpected execution order: %v", order)
	}
}

func TestPipeline_StopsOnFirstFailure(t *testing.T) {
	ran := []string{}
	errBoom := errors.New("boom")
	p := NewPipeline().
		Add("ok", func(ctx context.Context) error { ran = append(ran, "ok"); return nil }).
		Add("fail", func(ctx context.Context) error { ran = append(ran, "fail"); return errBoom }).
		Add("skip", func(ctx context.Context) error { ran = append(ran, "skip"); return nil })

	err := p.Run(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, errBoom) {
		t.Errorf("expected errBoom in chain, got %v", err)
	}
	if len(ran) != 2 || ran[1] != "fail" {
		t.Errorf("unexpected run order: %v", ran)
	}
}

func TestPipeline_CancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	p := NewPipeline().
		Add("should-not-run", func(ctx context.Context) error { return nil })

	if err := p.Run(ctx); err == nil {
		t.Fatal("expected error from cancelled context")
	}
}

func TestPipeline_NamesAndLen(t *testing.T) {
	p := NewPipeline().
		Add("alpha", func(ctx context.Context) error { return nil }).
		Add("beta", func(ctx context.Context) error { return nil })

	if p.Len() != 2 {
		t.Errorf("expected len 2, got %d", p.Len())
	}
	names := p.Names()
	if names[0] != "alpha" || names[1] != "beta" {
		t.Errorf("unexpected names: %v", names)
	}
}

func TestAsPipelineJob_WrapsError(t *testing.T) {
	errInner := errors.New("inner")
	p := NewPipeline().Add("fail", func(ctx context.Context) error { return errInner })
	job := AsPipelineJob(p, "my-pipeline")
	err := job(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, errInner) {
		t.Errorf("expected inner error in chain, got %v", err)
	}
}
