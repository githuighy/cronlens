package scheduler

import "time"

// EventKind describes the type of scheduler event.
type EventKind string

const (
	EventFired   EventKind = "fired"
	EventSkipped EventKind = "skipped"
	EventFailed  EventKind = "failed"
	EventRetried EventKind = "retried"
)

// Event represents a single occurrence in a schedule's lifecycle.
type Event struct {
	Kind      EventKind
	Schedule  string
	At        time.Time
	Message   string
	Err       error
}

// EventBus collects and broadcasts scheduler events to registered listeners.
type EventBus struct {
	listeners []func(Event)
}

// NewEventBus returns an initialised EventBus.
func NewEventBus() *EventBus {
	return &EventBus{}
}

// Subscribe registers a listener that will be called for every published event.
func (b *EventBus) Subscribe(fn func(Event)) {
	if fn != nil {
		b.listeners = append(b.listeners, fn)
	}
}

// Publish sends an event to all registered listeners.
func (b *EventBus) Publish(e Event) {
	if e.At.IsZero() {
		e.At = time.Now()
	}
	for _, fn := range b.listeners {
		fn(e)
	}
}

// Len returns the number of registered listeners.
func (b *EventBus) Len() int {
	return len(b.listeners)
}
