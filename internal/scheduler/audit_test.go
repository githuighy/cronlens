package scheduler

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestAuditLog_RecordAndAll(t *testing.T) {
	log := NewAuditLog(10)
	log.Record("job", AuditEvent{Name: "job", StartedAt: time.Now(), Duration: time.Millisecond})
	log.Record("job", AuditEvent{Name: "job", StartedAt: time.Now(), Duration: 2 * time.Millisecond})

	evs := log.All("job")
	if len(evs) != 2 {
		t.Fatalf("expected 2 events, got %d", len(evs))
	}
}

func TestAuditLog_CapEviction(t *testing.T) {
	log := NewAuditLog(3)
	for i := 0; i < 5; i++ {
		log.Record("job", AuditEvent{Name: "job", Duration: time.Duration(i) * time.Millisecond})
	}
	evs := log.All("job")
	if len(evs) != 3 {
		t.Fatalf("expected cap of 3, got %d", len(evs))
	}
	// oldest evicted: durations should be 2ms, 3ms, 4ms
	if evs[0].Duration != 2*time.Millisecond {
		t.Errorf("expected first remaining duration 2ms, got %v", evs[0].Duration)
	}
}

func TestAuditLog_Last(t *testing.T) {
	log := NewAuditLog(10)

	_, ok := log.Last("missing")
	if ok {
		t.Fatal("expected false for missing job")
	}

	log.Record("job", AuditEvent{Name: "job", Duration: time.Second})
	ev, ok := log.Last("job")
	if !ok {
		t.Fatal("expected event")
	}
	if ev.Duration != time.Second {
		t.Errorf("unexpected duration: %v", ev.Duration)
	}
}

func TestAuditLog_Clear(t *testing.T) {
	log := NewAuditLog(10)
	log.Record("job", AuditEvent{Name: "job"})
	log.Clear("job")
	if evs := log.All("job"); len(evs) != 0 {
		t.Fatalf("expected empty after clear, got %d", len(evs))
	}
}

func TestWithAudit_RecordsSuccess(t *testing.T) {
	log := NewAuditLog(10)
	job := WithAudit("j", log, func(_ context.Context) error { return nil })
	if err := job(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ev, ok := log.Last("j")
	if !ok {
		t.Fatal("expected audit event")
	}
	if ev.Err != nil {
		t.Errorf("expected nil err, got %v", ev.Err)
	}
	if ev.Duration <= 0 {
		t.Error("expected positive duration")
	}
}

func TestWithAudit_RecordsFailure(t *testing.T) {
	log := NewAuditLog(10)
	sentinel := errors.New("boom")
	job := WithAudit("j", log, func(_ context.Context) error { return sentinel })
	_ = job(context.Background())

	ev, ok := log.Last("j")
	if !ok {
		t.Fatal("expected audit event")
	}
	if !errors.Is(ev.Err, sentinel) {
		t.Errorf("expected sentinel error, got %v", ev.Err)
	}
}

func TestAuditLog_DefaultCap(t *testing.T) {
	log := NewAuditLog(0) // should default to 100
	for i := 0; i < 100; i++ {
		log.Record("j", AuditEvent{})
	}
	if len(log.All("j")) != 100 {
		t.Error("expected default cap of 100")
	}
}
