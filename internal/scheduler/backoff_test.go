package scheduler

import (
	"testing"
	"time"
)

func TestBackoffPolicy_Fixed(t *testing.T) {
	p := BackoffPolicy{
		Strategy:  BackoffFixed,
		BaseDelay: 5 * time.Second,
		MaxDelay:  time.Minute,
	}

	for _, attempt := range []int{1, 2, 5, 10} {
		if got := p.Delay(attempt); got != 5*time.Second {
			t.Errorf("attempt %d: expected 5s, got %v", attempt, got)
		}
	}
}

func TestBackoffPolicy_Linear(t *testing.T) {
	p := BackoffPolicy{
		Strategy:  BackoffLinear,
		BaseDelay: 2 * time.Second,
		MaxDelay:  time.Minute,
	}

	cases := []struct {
		attempt int
		want    time.Duration
	}{
		{1, 2 * time.Second},
		{2, 4 * time.Second},
		{3, 6 * time.Second},
		{5, 10 * time.Second},
	}

	for _, tc := range cases {
		if got := p.Delay(tc.attempt); got != tc.want {
			t.Errorf("attempt %d: expected %v, got %v", tc.attempt, tc.want, got)
		}
	}
}

func TestBackoffPolicy_Exponential(t *testing.T) {
	p := BackoffPolicy{
		Strategy:  BackoffExponential,
		BaseDelay: time.Second,
		MaxDelay:  time.Minute,
	}

	cases := []struct {
		attempt int
		want    time.Duration
	}{
		{1, 1 * time.Second},
		{2, 2 * time.Second},
		{3, 4 * time.Second},
		{4, 8 * time.Second},
		{5, 16 * time.Second},
	}

	for _, tc := range cases {
		if got := p.Delay(tc.attempt); got != tc.want {
			t.Errorf("attempt %d: expected %v, got %v", tc.attempt, tc.want, got)
		}
	}
}

func TestBackoffPolicy_MaxDelayCap(t *testing.T) {
	p := BackoffPolicy{
		Strategy:  BackoffExponential,
		BaseDelay: time.Second,
		MaxDelay:  10 * time.Second,
	}

	// attempt 5 would be 16s without cap
	if got := p.Delay(5); got != 10*time.Second {
		t.Errorf("expected delay capped at 10s, got %v", got)
	}
}

func TestDefaultBackoffPolicy(t *testing.T) {
	p := DefaultBackoffPolicy()

	if p.Strategy != BackoffExponential {
		t.Errorf("expected exponential strategy")
	}
	if p.BaseDelay != time.Second {
		t.Errorf("expected 1s base delay, got %v", p.BaseDelay)
	}
	if p.MaxDelay != time.Minute {
		t.Errorf("expected 1m max delay, got %v", p.MaxDelay)
	}
}

func TestBackoffPolicy_ZeroAttemptClamped(t *testing.T) {
	p := DefaultBackoffPolicy()
	// attempt 0 should behave like attempt 1
	if got := p.Delay(0); got != time.Second {
		t.Errorf("expected 1s for attempt 0, got %v", got)
	}
}
