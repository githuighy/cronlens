package humanizer

import (
	"fmt"
	"strings"

	"github.com/cronlens/cronlens/internal/parser"
)

// Humanize converts a parsed cron expression into a human-readable description.
func Humanize(expr *parser.Expression) string {
	parts := []string{}

	minDesc := describeField(expr.Minute, "minute", 0, 59)
	hourDesc := describeField(expr.Hour, "hour", 0, 23)
	dayDesc := describeField(expr.DayOfMonth, "day of month", 1, 31)
	monthDesc := describeField(expr.Month, "month", 1, 12)
	weekDesc := describeField(expr.DayOfWeek, "day of week", 0, 6)

	if isEvery(expr.Minute) && isEvery(expr.Hour) && isEvery(expr.DayOfMonth) && isEvery(expr.Month) && isEvery(expr.DayOfWeek) {
		return "Every minute"
	}

	if isEvery(expr.Minute) {
		parts = append(parts, "every minute")
	} else {
		parts = append(parts, minDesc)
	}

	if !isEvery(expr.Hour) {
		parts = append(parts, hourDesc)
	}

	if !isEvery(expr.DayOfMonth) {
		parts = append(parts, dayDesc)
	}

	if !isEvery(expr.Month) {
		parts = append(parts, monthDesc)
	}

	if !isEvery(expr.DayOfWeek) {
		parts = append(parts, weekDesc)
	}

	result := strings.Join(parts, ", ")
	return capitalize(result)
}

func isEvery(values []int) bool {
	if len(values) == 0 {
		return true
	}
	expected := values[len(values)-1] - values[0] + 1
	return len(values) == expected
}

func describeField(values []int, unit string, min, max int) string {
	if len(values) == 0 {
		return fmt.Sprintf("every %s", unit)
	}
	if len(values) == 1 {
		return fmt.Sprintf("%s %d", unit, values[0])
	}
	strs := make([]string, len(values))
	for i, v := range values {
		strs[i] = fmt.Sprintf("%d", v)
	}
	return fmt.Sprintf("%s %s", unit, strings.Join(strs, "/"))
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
