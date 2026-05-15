package scheduler

import (
	"fmt"
	"time"
)

// TimeUntilNext returns the duration until the next scheduled run from now.
func (s *Schedule) TimeUntilNext(from time.Time) time.Duration {
	next := s.NextRun(from)
	if next.IsZero() {
		return 0
	}
	return next.Sub(from)
}

// TimeSinceLast returns the duration since the last scheduled run from now.
func (s *Schedule) TimeSinceLast(from time.Time) time.Duration {
	last := s.LastRun(from)
	if last.IsZero() {
		return 0
	}
	return from.Sub(last)
}

// FormatDuration formats a duration into a human-readable string.
func FormatDuration(d time.Duration) string {
	if d < 0 {
		d = -d
	}

	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	switch {
	case days > 0:
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	case hours > 0:
		return fmt.Sprintf("%dh %dm", hours, minutes)
	case minutes > 0:
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	default:
		return fmt.Sprintf("%ds", seconds)
	}
}
