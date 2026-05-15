package scheduler

import (
	"errors"
	"time"

	"github.com/cronlens/cronlens/internal/parser"
)

// ErrNoPreviousRun is returned when no previous run can be determined.
var ErrNoPreviousRun = errors.New("scheduler: no previous run found within search window")

// maxSearchBack is the maximum window we search backwards for a previous run.
const maxSearchBack = 366 * 24 * time.Hour

// lastRun finds the most recent time before `now` that the expression matched.
// It walks backwards minute by minute up to maxSearchBack.
func lastRun(expr *parser.Expression, now time.Time) (time.Time, error) {
	// Start from the previous minute
	t := now.Add(-time.Minute).Truncate(time.Minute)
	deadline := now.Add(-maxSearchBack)

	for t.After(deadline) {
		if expr.Matches(t) {
			return t, nil
		}
		t = t.Add(-time.Minute)
	}

	return time.Time{}, ErrNoPreviousRun
}
