// Package scheduler provides utilities for computing cron schedule summaries,
// including next N run times and elapsed time since last run.
package scheduler

import (
	"time"

	"github.com/cronlens/cronlens/internal/parser"
	"github.com/cronlens/cronlens/internal/predictor"
)

// Schedule holds a parsed cron expression and its associated timezone.
type Schedule struct {
	Expression *parser.Expression
	Location   *time.Location
}

// New creates a new Schedule from a cron expression string and timezone name.
// If tzName is empty, UTC is used.
func New(expr string, tzName string) (*Schedule, error) {
	loc, err := resolveLocation(tzName)
	if err != nil {
		return nil, err
	}

	parsed, err := parser.Parse(expr)
	if err != nil {
		return nil, err
	}

	return &Schedule{
		Expression: parsed,
		Location:   loc,
	}, nil
}

// NextRun returns the next time the schedule will fire after now.
func (s *Schedule) NextRun(now time.Time) (time.Time, error) {
	return predictor.NextRun(s.Expression, now.In(s.Location))
}

// NextN returns the next n run times after now.
func (s *Schedule) NextN(now time.Time, n int) ([]time.Time, error) {
	return predictor.NextN(s.Expression, now.In(s.Location), n)
}

// TimeSinceLast returns the duration elapsed since the most recent scheduled
// run before now. Returns ErrNoPreviousRun if no previous run can be found
// within a reasonable search window.
func (s *Schedule) TimeSinceLast(now time.Time) (time.Duration, time.Time, error) {
	now = now.In(s.Location)
	last, err := lastRun(s.Expression, now)
	if err != nil {
		return 0, time.Time{}, err
	}
	return now.Sub(last), last, nil
}

// resolveLocation parses a timezone name into a *time.Location.
func resolveLocation(tzName string) (*time.Location, error) {
	if tzName == "" {
		return time.UTC, nil
	}
	return time.LoadLocation(tzName)
}
