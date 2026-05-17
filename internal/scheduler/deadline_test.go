package scheduler

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDeadlineStore_SetAndGet(t *testing.T) {
	store := NewDeadlineStore()
	expiry := time.Now().Add(time.Hour)
	store.Set("job1", expiry)

	got, ok := store.Get("job1")
	if !ok {
		t.Fatal("expected deadline to exist")
	}
	if !got.Equal(expiry) {
		t.Errorf("expected %v, got %v", expiry, got)
	}
}

func TestDeadlineStore_GetMissing(t *testing.T) {
	store := NewDeadlineStore()
	_, ok := store.Get("missing")
	if ok {
		t.Error("expected missing deadline to return false")
	}
}

func TestDeadlineStore_Expired_NotYet(t *testing.T) {
	store := NewDeadlineStore()
	store.Set("future", time.Now().Add(time.Hour))
	if store.Expired("future", time.Now()) {
		t.Error("deadline should not be expired yet")
	}
}

func TestDeadlineStore_Expired_Past(t *testing.T) {
	store := NewDeadlineStore()
	store.Set("past", time.Now().Add(-time.Second))
	if !store.Expired("past", time.Now()) {
		t.Error("deadline should be expired")
	}
}

func TestDeadlineStore_Expired_Missing(t *testing.T) {
	store := NewDeadlineStore()
	if store.Expired("nonexistent", time.Now()) {
		t.Error("missing deadline should not be considered expired")
	}
}

func TestDeadlineStore_Delete(t *testing.T) {
	store := NewDeadlineStore()
	store.Set("temp", time.Now().Add(time.Hour))
	store.Delete("temp")
	_, ok := store.Get("temp")
	if ok {
		t.Error("expected deadline to be deleted")
	}
}

func TestWithDeadline_NotExpired(t *testing.T) {
	store := NewDeadlineStore()
	store.Set("job", time.Now().Add(time.Hour))

	called := false
	job := WithDeadline(store, "job", func(ctx context.Context) error {
		called = true
		return nil
	})

	if err := job(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected job to be called")
	}
}

func TestWithDeadline_Expired(t *testing.T) {
	store := NewDeadlineStore()
	store.Set("job", time.Now().Add(-time.Second))

	job := WithDeadline(store, "job", func(ctx context.Context) error {
		return nil
	})

	err := job(context.Background())
	if err == nil {
		t.Fatal("expected error for expired deadline")
	}
	if !errors.Is(err, err) {
		t.Errorf("unexpected error type: %v", err)
	}
}
