package scheduler

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestThrottle_Allow_FirstCallAlwaysPasses(t *testing.T) {
	th := NewThrottle(ThrottlePolicy{MinInterval: time.Hour})
	if !th.Allow("job1") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestThrottle_Allow_SecondCallBlocked(t *testing.T) {
	th := NewThrottle(ThrottlePolicy{MinInterval: time.Hour})
	th.Allow("job1")
	if th.Allow("job1") {
		t.Fatal("expected second call within interval to be blocked")
	}
}

func TestThrottle_Allow_PassesAfterInterval(t *testing.T) {
	th := NewThrottle(ThrottlePolicy{MinInterval: 10 * time.Millisecond})
	th.Allow("job1")
	time.Sleep(20 * time.Millisecond)
	if !th.Allow("job1") {
		t.Fatal("expected call after interval to be allowed")
	}
}

func TestThrottle_Allow_IndependentNames(t *testing.T) {
	th := NewThrottle(ThrottlePolicy{MinInterval: time.Hour})
	th.Allow("job1")
	if !th.Allow("job2") {
		t.Fatal("expected different job name to be allowed independently")
	}
}

func TestThrottle_Reset_ClearsState(t *testing.T) {
	th := NewThrottle(ThrottlePolicy{MinInterval: time.Hour})
	th.Allow("job1")
	th.Reset("job1")
	if !th.Allow("job1") {
		t.Fatal("expected call after reset to be allowed")
	}
}

func TestWithThrottle_AllowsFirstRun(t *testing.T) {
	th := NewThrottle(ThrottlePolicy{MinInterval: time.Hour})
	called := false
	job := func(ctx context.Context) error {
		called = true
		return nil
	}
	mw := WithThrottle(th, "myjob")
	err := mw(job)(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected job to be called")
	}
}

func TestWithThrottle_BlocksSecondRun(t *testing.T) {
	th := NewThrottle(ThrottlePolicy{MinInterval: time.Hour})
	job := func(ctx context.Context) error { return nil }
	mw := WithThrottle(th, "myjob")
	_ = mw(job)(context.Background())
	err := mw(job)(context.Background())
	if err == nil {
		t.Fatal("expected throttle error on second call")
	}
	if !strings.Contains(err.Error(), "throttled") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestWithThrottle_PropagatesJobError(t *testing.T) {
	th := NewThrottle(ThrottlePolicy{MinInterval: 0})
	expected := errors.New("job failed")
	job := func(ctx context.Context) error { return expected }
	mw := WithThrottle(th, "myjob")
	err := mw(job)(context.Background())
	if !errors.Is(err, expected) {
		t.Fatalf("expected job error, got: %v", err)
	}
}
