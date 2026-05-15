package scheduler

import (
	"time"

	"github.com/nabeelvalley/cronlens/internal/humanizer"
)

// Summary holds a human-readable overview of a schedule at a point in time.
type Summary struct {
	Expression  string
	Timezone    string
	Description string
	LastRun     time.Time
	NextRun     time.Time
	TimeUntil   time.Duration
	TimeSince   time.Duration
}

// Summarize returns a Summary of the schedule relative to the given time.
func (s *Schedule) Summarize(from time.Time) Summary {
	description, _ := humanizer.Humanize(s.expr)

	return Summary{
		Expression:  s.raw,
		Timezone:    s.loc.String(),
		Description: description,
		LastRun:     s.LastRun(from),
		NextRun:     s.NextRun(from),
		TimeUntil:   s.TimeUntilNext(from),
		TimeSince:   s.TimeSinceLast(from),
	}
}
