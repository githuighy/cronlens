package scheduler

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestCooldown_FirstCallAlwaysPasses(t *testing.T) {
	c := NewCooldown(5 * time.Second)
	if !c.Allow("job1") {
		t.Error("expected first call to always pass")
	}
}

func TestCooldown_SecondCallBlocked(t *testing.T) {
	c := NewCooldown(5 * time.Second)
	c.Allow("job1")
	if c.Allow("job1") {
		t.Error("expected second call to be blocked within cooldown period")
	}
}

func TestCooldown_PassesAfterPeriod(t *testing.T) {
	c := NewCooldown(50 * time.Millisecond)
	c.Allow("job1")
	time.Sleep(60 * time.Millisecond)
	if !c.Allow("job1") {
		t.Error("expected call to pass after cooldown period")
	}
}

func TestCooldown_IndependentNames(t *testing.T) {
	c := NewCooldown(5 * time.Second)
	c.Allow("job1")
	if !c.Allow("job2") {
		t.Error("expected independent job names to have separate cooldowns")
	}
}

func TestCooldown_Reset_ClearsState(t *testing.T) {
	c := NewCooldown(5 * time.Second)
	c.Allow("job1")
	c.Reset("job1")
	if !c.Allow("job1") {
		t.Error("expected Allow to pass after Reset")
	}
}

func TestCooldown_Remaining_ZeroWhenReady(t *testing.T) {
	c := NewCooldown(5 * time.Second)
	if r := c.Remaining("job1"); r != 0 {
		t.Errorf("expected 0 remaining for unseen job, got %v", r)
	}
}

func TestCooldown_Remaining_PositiveWhenActive(t *testing.T) {
	c := NewCooldown(5 * time.Second)
	c.Allow("job1")
	if r := c.Remaining("job1"); r <= 0 {
		t.Errorf("expected positive remaining duration, got %v", r)
	}
}

func TestWithCooldown_SkipsWhenActive(t *testing.T) {
	c := NewCooldown(5 * time.Second)
	called := 0
	job := func(ctx context.Context) error {
		called++
		return nil
	}
	wrapped := WithCooldown("job1", c, job)

	// First call should succeed
	if err := wrapped(context.Background()); err != nil {
		t.Errorf("unexpected error on first call: %v", err)
	}
	// Second call should be blocked
	if err := wrapped(context.Background()); err == nil {
		t.Error("expected error on second call within cooldown")
	}
	if called != 1 {
		t.Errorf("expected job called once, got %d", called)
	}
}

func TestWithCooldown_PropagatesJobError(t *testing.T) {
	c := NewCooldown(5 * time.Second)
	expected := errors.New("job failed")
	job := func(ctx context.Context) error { return expected }
	wrapped := WithCooldown("job1", c, job)

	if err := wrapped(context.Background()); !errors.Is(err, expected) {
		t.Errorf("expected job error to propagate, got %v", err)
	}
}
