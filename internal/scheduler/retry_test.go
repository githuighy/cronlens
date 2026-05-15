package scheduler

import (
	"errors"
	"testing"
	"time"
)

var errFake = errors.New("fake error")

func TestDefaultRetryPolicy(t *testing.T) {
	p := DefaultRetryPolicy()
	if p.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts=3, got %d", p.MaxAttempts)
	}
	if p.Delay != 5*time.Second {
		t.Errorf("expected Delay=5s, got %s", p.Delay)
	}
	if p.BackoffFactor != 2.0 {
		t.Errorf("expected BackoffFactor=2.0, got %f", p.BackoffFactor)
	}
}

func TestRunWithRetry_SuccessFirstAttempt(t *testing.T) {
	policy := RetryPolicy{MaxAttempts: 3, Delay: 0, BackoffFactor: 1.0}
	calls := 0
	results, err := RunWithRetry(policy, func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
	if len(results) != 1 || !results[0].Succeeded {
		t.Errorf("expected one successful result")
	}
}

func TestRunWithRetry_SuccessOnSecondAttempt(t *testing.T) {
	policy := RetryPolicy{MaxAttempts: 3, Delay: 0, BackoffFactor: 1.0}
	calls := 0
	results, err := RunWithRetry(policy, func() error {
		calls++
		if calls < 2 {
			return errFake
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 2 {
		t.Errorf("expected 2 calls, got %d", calls)
	}
	if !results[len(results)-1].Succeeded {
		t.Errorf("last result should be successful")
	}
}

func TestRunWithRetry_AllAttemptsFail(t *testing.T) {
	policy := RetryPolicy{MaxAttempts: 3, Delay: 0, BackoffFactor: 1.0}
	calls := 0
	results, err := RunWithRetry(policy, func() error {
		calls++
		return errFake
	})
	if !errors.Is(err, ErrMaxAttemptsReached) {
		t.Fatalf("expected ErrMaxAttemptsReached, got %v", err)
	}
	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
	if len(results) != 3 {
		t.Errorf("expected 3 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Succeeded {
			t.Errorf("result attempt %d should not be succeeded", r.Attempt)
		}
	}
}

func TestRunWithRetry_ZeroMaxAttemptsTreatedAsOne(t *testing.T) {
	policy := RetryPolicy{MaxAttempts: 0, Delay: 0, BackoffFactor: 1.0}
	calls := 0
	_, err := RunWithRetry(policy, func() error {
		calls++
		return errFake
	})
	if !errors.Is(err, ErrMaxAttemptsReached) {
		t.Fatalf("expected ErrMaxAttemptsReached, got %v", err)
	}
	if calls != 1 {
		t.Errorf("expected exactly 1 call, got %d", calls)
	}
}
