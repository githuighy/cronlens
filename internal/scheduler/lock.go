package scheduler

import (
	"errors"
	"sync"
	"time"
)

// ErrLockTimeout is returned when a lock cannot be acquired within the deadline.
var ErrLockTimeout = errors.New("scheduler: lock acquisition timed out")

// Lock provides a named, time-bounded mutual exclusion mechanism for
// preventing concurrent execution of scheduled jobs by name.
type Lock struct {
	mu    sync.Mutex
	held  map[string]time.Time
	tttl  time.Duration
}

// NewLock creates a Lock where each acquired entry expires after ttl.
// A zero ttl means locks never expire automatically.
func NewLock(ttl time.Duration) *Lock {
	return &Lock{
		held: make(map[string]time.Time),
		tttl: ttl,
	}
}

// Acquire attempts to lock name. Returns nil on success.
// If the lock is already held (and not expired), ErrLockTimeout is returned.
func (l *Lock) Acquire(name string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if exp, ok := l.held[name]; ok {
		if l.tttl == 0 || time.Now().Before(exp) {
			return ErrLockTimeout
		}
	}

	var expiry time.Time
	if l.tttl > 0 {
		expiry = time.Now().Add(l.tttl)
	}
	l.held[name] = expiry
	return nil
}

// Release releases the lock for name. It is a no-op if the lock is not held.
func (l *Lock) Release(name string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.held, name)
}

// IsHeld reports whether name is currently locked (and not expired).
func (l *Lock) IsHeld(name string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	exp, ok := l.held[name]
	if !ok {
		return false
	}
	if l.tttl > 0 && !time.Now().Before(exp) {
		delete(l.held, name)
		return false
	}
	return true
}
