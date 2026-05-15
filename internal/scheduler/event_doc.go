// Package scheduler provides scheduling primitives built on top of parsed cron
// expressions.
//
// # EventBus
//
// EventBus is a lightweight publish/subscribe mechanism that lets consumers
// observe what is happening inside a running scheduler without coupling them to
// its internals.
//
// Typical usage:
//
//	bus := scheduler.NewEventBus()
//	bus.Subscribe(func(e scheduler.Event) {
//		fmt.Printf("[%s] %s – %s\n", e.Kind, e.Schedule, e.At.Format(time.RFC3339))
//	})
//
// Events are published synchronously in the order listeners were registered.
// If a listener panics the panic will propagate to the publisher; wrap
// sensitive listeners with recover if necessary.
package scheduler
