package scheduler

import (
	"context"
	"testing"
	"time"
)

func TestQuota_FirstCallAlwaysPasses(t *testing.T) {
	q := NewQuota(time.Hour)
	_ = q.Set("job", 2)
	if !q.Allow("job", time.Now()) {
		t.Fatal("expected first call to be allowed")
	}
}

func TestQuota_BlocksAfterMax(t *testing.T) {
	q := NewQuota(time.Hour)
	_ = q.Set("job", 2)
	now := time.Now()
	q.Allow("job", now)
	q.Allow("job", now.Add(time.Second))
	if q.Allow("job", now.Add(2*time.Second)) {
		t.Fatal("expected third call to be blocked")
	}
}

func TestQuota_AllowsAfterWindowExpires(t *testing.T) {
	q := NewQuota(time.Minute)
	_ = q.Set("job", 1)
	base := time.Now()
	q.Allow("job", base)
	// advance past the window
	if !q.Allow("job", base.Add(2*time.Minute)) {
		t.Fatal("expected call to be allowed after window expiry")
	}
}

func TestQuota_NoLimitSetAlwaysPasses(t *testing.T) {
	q := NewQuota(time.Hour)
	for i := 0; i < 100; i++ {
		if !q.Allow("unrestricted", time.Now()) {
			t.Fatal("expected unrestricted job to always pass")
		}
	}
}

func TestQuota_IndependentNames(t *testing.T) {
	q := NewQuota(time.Hour)
	_ = q.Set("a", 1)
	_ = q.Set("b", 1)
	now := time.Now()
	q.Allow("a", now)
	if !q.Allow("b", now) {
		t.Fatal("quota for 'b' should be independent of 'a'")
	}
}

func TestQuota_DeleteResetsState(t *testing.T) {
	q := NewQuota(time.Hour)
	_ = q.Set("job", 1)
	now := time.Now()
	q.Allow("job", now)
	q.Delete("job")
	if !q.Allow("job", now.Add(time.Second)) {
		t.Fatal("expected allow after delete")
	}
}

func TestQuota_SetInvalidMax(t *testing.T) {
	q := NewQuota(time.Hour)
	if err := q.Set("job", 0); err == nil {
		t.Fatal("expected error for max=0")
	}
}

func TestWithQuota_SkipsWhenExceeded(t *testing.T) {
	q := NewQuota(time.Hour)
	_ = q.Set("job", 1)
	// consume the quota
	q.Allow("job", time.Now())

	calls := 0
	wrapped := WithQuota(q, "job", func(ctx context.Context) error {
		calls++
		return nil
	})

	_ = wrapped(context.Background())
	if calls != 0 {
		t.Fatalf("expected 0 calls, got %d", calls)
	}
}
