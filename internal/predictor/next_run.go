package predictor

import (
	"time"

	"github.com/cronlens/cronlens/internal/parser"
)

// NextRun returns the next time the given cron expression will trigger
// after the provided reference time, evaluated in the given timezone.
func NextRun(expr *parser.Expression, after time.Time, loc *time.Location) (time.Time, error) {
	if loc == nil {
		loc = time.UTC
	}

	// Normalize to the given timezone and truncate sub-minute precision.
	t := after.In(loc).Truncate(time.Minute).Add(time.Minute)

	// Search up to 4 years ahead to avoid infinite loops on impossible expressions.
	limit := t.Add(4 * 365 * 24 * time.Hour)

	for t.Before(limit) {
		if !monthMatches(expr, t) {
			t = advanceToNextMonth(t)
			continue
		}
		if !expr.Matches(t) {
			t = t.Add(time.Minute)
			continue
		}
		return t, nil
	}

	return time.Time{}, ErrNoNextRun
}

// NextN returns the next n scheduled times after the given reference time.
func NextN(expr *parser.Expression, after time.Time, loc *time.Location, n int) ([]time.Time, error) {
	results := make([]time.Time, 0, n)
	current := after

	for len(results) < n {
		next, err := NextRun(expr, current, loc)
		if err != nil {
			return results, err
		}
		results = append(results, next)
		current = next
	}

	return results, nil
}

// monthMatches checks only the month field to allow fast month-level skipping.
func monthMatches(expr *parser.Expression, t time.Time) bool {
	month := int(t.Month())
	for _, v := range expr.Month.Values {
		if v == month {
			return true
		}
	}
	return false
}

// advanceToNextMonth moves t to the first minute of the next month.
func advanceToNextMonth(t time.Time) time.Time {
	first := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	return first.AddDate(0, 1, 0)
}
