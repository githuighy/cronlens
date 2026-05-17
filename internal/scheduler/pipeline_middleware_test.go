package scheduler

import (
	"context"
	"errors"
	"testing"
)

func TestWithPipelineLogging_AllSucceed(t *testing.T) {
	executed := 0
	p := NewPipeline().
		Add("s1", func(ctx context.Context) error { executed++; return nil }).
		Add("s2", func(ctx context.Context) error { executed++; return nil })

	logged := WithPipelineLogging(p)
	if err := logged.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if executed != 2 {
		t.Errorf("expected 2 executions, got %d", executed)
	}
}

func TestWithPipelineLogging_StopsOnError(t *testing.T) {
	errFail := errors.New("fail")
	executed := 0
	p := NewPipeline().
		Add("ok", func(ctx context.Context) error { executed++; return nil }).
		Add("bad", func(ctx context.Context) error { executed++; return errFail }).
		Add("skip", func(ctx context.Context) error { executed++; return nil })

	logged := WithPipelineLogging(p)
	err := logged.Run(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, errFail) {
		t.Errorf("expected errFail, got %v", err)
	}
	if executed != 2 {
		t.Errorf("expected 2 executions before stop, got %d", executed)
	}
}

func TestWithPipelineLogging_PreservesNames(t *testing.T) {
	p := NewPipeline().
		Add("x", func(ctx context.Context) error { return nil }).
		Add("y", func(ctx context.Context) error { return nil })

	logged := WithPipelineLogging(p)
	names := logged.Names()
	if len(names) != 2 || names[0] != "x" || names[1] != "y" {
		t.Errorf("unexpected names after logging wrap: %v", names)
	}
}
