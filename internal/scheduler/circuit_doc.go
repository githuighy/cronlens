// Package scheduler provides scheduling primitives for cron-based job execution.
//
// # Circuit Breaker
//
// CircuitBreaker protects jobs from repeated execution when downstream
// dependencies are consistently failing. It tracks consecutive failures
// and transitions through three states:
//
//   - Closed  – normal operation; all calls are allowed through.
//   - Open    – failure threshold exceeded; calls are rejected immediately
//     with ErrCircuitOpen, preventing further load on a failing system.
//   - HalfOpen – after resetTimeout elapses, one probe call is allowed;
//     success closes the circuit, failure reopens it.
//
// # Usage
//
//	cb := scheduler.NewCircuitBreaker(5, 30*time.Second)
//
//	protected := scheduler.WithCircuitBreaker(cb, func(ctx context.Context) error {
//		return callExternalService(ctx)
//	})
//
// WithCircuitBreaker can be composed with other middleware such as
// WithLogging, WithTimeout, and WithRecover via Chain.
package scheduler
