package scheduler

import (
	"errors"
	"testing"
	"time"
)

func TestEventBus_NoListeners(t *testing.T) {
	bus := NewEventBus()
	// Should not panic with no listeners.
	bus.Publish(Event{Kind: EventFired, Schedule: "test"})
}

func TestEventBus_Subscribe_NilIgnored(t *testing.T) {
	bus := NewEventBus()
	bus.Subscribe(nil)
	if bus.Len() != 0 {
		t.Fatalf("expected 0 listeners, got %d", bus.Len())
	}
}

func TestEventBus_Publish_CallsAllListeners(t *testing.T) {
	bus := NewEventBus()

	var calls []EventKind
	bus.Subscribe(func(e Event) { calls = append(calls, e.Kind) })
	bus.Subscribe(func(e Event) { calls = append(calls, e.Kind) })

	bus.Publish(Event{Kind: EventFired, Schedule: "@hourly"})

	if len(calls) != 2 {
		t.Fatalf("expected 2 calls, got %d", len(calls))
	}
	for _, k := range calls {
		if k != EventFired {
			t.Errorf("unexpected kind %q", k)
		}
	}
}

func TestEventBus_Publish_SetsTimestamp(t *testing.T) {
	bus := NewEventBus()

	var received Event
	bus.Subscribe(func(e Event) { received = e })

	before := time.Now()
	bus.Publish(Event{Kind: EventSkipped, Schedule: "* * * * *"})
	after := time.Now()

	if received.At.Before(before) || received.At.After(after) {
		t.Errorf("timestamp %v out of expected range [%v, %v]", received.At, before, after)
	}
}

func TestEventBus_Publish_PreservesExplicitTimestamp(t *testing.T) {
	bus := NewEventBus()

	fixed := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	var received Event
	bus.Subscribe(func(e Event) { received = e })

	bus.Publish(Event{Kind: EventFired, Schedule: "0 9 * * 1", At: fixed})

	if !received.At.Equal(fixed) {
		t.Errorf("expected timestamp %v, got %v", fixed, received.At)
	}
}

func TestEventBus_Publish_ErrorField(t *testing.T) {
	bus := NewEventBus()

	sentinel := errors.New("job failed")
	var received Event
	bus.Subscribe(func(e Event) { received = e })

	bus.Publish(Event{Kind: EventFailed, Schedule: "*/5 * * * *", Err: sentinel})

	if !errors.Is(received.Err, sentinel) {
		t.Errorf("expected sentinel error, got %v", received.Err)
	}
}
