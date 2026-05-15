package scheduler

import (
	"context"
	"errors"
	"sync"
	"time"
)

// CircuitState represents the state of a circuit breaker.
type CircuitState int

const (
	CircuitClosed   CircuitState = iota // normal operation
	CircuitOpen                         // failing, requests blocked
	CircuitHalfOpen                     // probing for recovery
)

// ErrCircuitOpen is returned when the circuit breaker is open.
var ErrCircuitOpen = errors.New("circuit breaker is open")

// CircuitBreaker tracks consecutive failures and opens the circuit
// after a configurable threshold, allowing recovery after a timeout.
type CircuitBreaker struct {
	mu           sync.Mutex
	state        CircuitState
	failures     int
	maxFailures  int
	resetTimeout time.Duration
	openedAt     time.Time
}

// NewCircuitBreaker creates a CircuitBreaker that opens after maxFailures
// consecutive failures and attempts recovery after resetTimeout.
func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
	if maxFailures <= 0 {
		maxFailures = 3
	}
	if resetTimeout <= 0 {
		resetTimeout = 30 * time.Second
	}
	return &CircuitBreaker{
		state:        CircuitClosed,
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
	}
}

// State returns the current circuit state.
func (cb *CircuitBreaker) State() CircuitState {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.currentState()
}

// currentState computes state, transitioning Open→HalfOpen after timeout.
// Caller must hold cb.mu.
func (cb *CircuitBreaker) currentState() CircuitState {
	if cb.state == CircuitOpen && time.Since(cb.openedAt) >= cb.resetTimeout {
		cb.state = CircuitHalfOpen
	}
	return cb.state
}

// Allow returns nil if the request may proceed, or ErrCircuitOpen.
func (cb *CircuitBreaker) Allow() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	if cb.currentState() == CircuitOpen {
		return ErrCircuitOpen
	}
	return nil
}

// RecordSuccess resets failure count and closes the circuit.
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures = 0
	cb.state = CircuitClosed
}

// RecordFailure increments failure count and may open the circuit.
func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures++
	if cb.failures >= cb.maxFailures {
		cb.state = CircuitOpen
		cb.openedAt = time.Now()
	}
}

// WithCircuitBreaker wraps a JobFunc with circuit-breaker protection.
func WithCircuitBreaker(cb *CircuitBreaker, next JobFunc) JobFunc {
	return func(ctx context.Context) error {
		if err := cb.Allow(); err != nil {
			return err
		}
		err := next(ctx)
		if err != nil {
			cb.RecordFailure()
		} else {
			cb.RecordSuccess()
		}
		return err
	}
}
