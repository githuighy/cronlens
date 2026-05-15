package humanizer_test

import (
	"testing"

	"github.com/cronlens/cronlens/internal/humanizer"
	"github.com/cronlens/cronlens/internal/parser"
)

func mustParse(t *testing.T, expr string) *parser.Expression {
	t.Helper()
	e, err := parser.Parse(expr)
	if err != nil {
		t.Fatalf("failed to parse %q: %v", expr, err)
	}
	return e
}

func TestHumanize_EveryMinute(t *testing.T) {
	expr := mustParse(t, "* * * * *")
	got := humanizer.Humanize(expr)
	if got != "Every minute" {
		t.Errorf("expected 'Every minute', got %q", got)
	}
}

func TestHumanize_SpecificMinute(t *testing.T) {
	expr := mustParse(t, "30 * * * *")
	got := humanizer.Humanize(expr)
	if got == "" {
		t.Error("expected non-empty description")
	}
	t.Logf("30 * * * * => %s", got)
}

func TestHumanize_SpecificTime(t *testing.T) {
	expr := mustParse(t, "0 9 * * *")
	got := humanizer.Humanize(expr)
	if got == "" {
		t.Error("expected non-empty description")
	}
	t.Logf("0 9 * * * => %s", got)
}

func TestHumanize_WeekdayAndMonth(t *testing.T) {
	expr := mustParse(t, "0 8 * 6 1")
	got := humanizer.Humanize(expr)
	if got == "" {
		t.Error("expected non-empty description")
	}
	t.Logf("0 8 * 6 1 => %s", got)
}

func TestMonthName(t *testing.T) {
	cases := []struct {
		input    int
		expected string
	}{
		{1, "January"},
		{6, "June"},
		{12, "December"},
		{13, ""},
	}
	for _, tc := range cases {
		got := humanizer.MonthName(tc.input)
		if got != tc.expected {
			t.Errorf("MonthName(%d): expected %q, got %q", tc.input, tc.expected, got)
		}
	}
}

func TestWeekdayName(t *testing.T) {
	cases := []struct {
		input    int
		expected string
	}{
		{0, "Sunday"},
		{5, "Friday"},
		{6, "Saturday"},
		{7, ""},
	}
	for _, tc := range cases {
		got := humanizer.WeekdayName(tc.input)
		if got != tc.expected {
			t.Errorf("WeekdayName(%d): expected %q, got %q", tc.input, tc.expected, got)
		}
	}
}
