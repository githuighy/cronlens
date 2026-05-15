package predictor_test

import (
	"testing"
	"time"

	"github.com/cronlens/cronlens/internal/parser"
	"github.com/cronlens/cronlens/internal/predictor"
)

func mustParse(t *testing.T, expr string) *parser.Expression {
	t.Helper()
	e, err := parser.Parse(expr)
	if err != nil {
		t.Fatalf("failed to parse expression %q: %v", expr, err)
	}
	return e
}

func TestNextRun_EveryMinute(t *testing.T) {
	expr := mustParse(t, "* * * * *")
	after := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	want := time.Date(2024, 1, 15, 10, 31, 0, 0, time.UTC)

	got, err := predictor.NextRun(expr, after, time.UTC)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got.Equal(want) {
		t.Errorf("NextRun = %v, want %v", got, want)
	}
}

func TestNextRun_SpecificTime(t *testing.T) {
	// Every day at 09:00
	expr := mustParse(t, "0 9 * * *")
	after := time.Date(2024, 3, 10, 9, 1, 0, 0, time.UTC)
	want := time.Date(2024, 3, 11, 9, 0, 0, 0, time.UTC)

	got, err := predictor.NextRun(expr, after, time.UTC)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got.Equal(want) {
		t.Errorf("NextRun = %v, want %v", got, want)
	}
}

func TestNextRun_TimezoneAware(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Skip("timezone data not available")
	}

	// Every day at midnight New York time
	expr := mustParse(t, "0 0 * * *")
	// 23:30 UTC on Jan 15 = 18:30 EST — next midnight EST is Jan 16 00:00 EST = 05:00 UTC
	after := time.Date(2024, 1, 15, 23, 30, 0, 0, time.UTC)

	got, err := predictor.NextRun(expr, after, loc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Location() != loc {
		t.Errorf("expected result in %v, got %v", loc, got.Location())
	}
	if got.Hour() != 0 || got.Minute() != 0 {
		t.Errorf("expected midnight, got %02d:%02d", got.Hour(), got.Minute())
	}
}

func TestNextN(t *testing.T) {
	expr := mustParse(t, "0 * * * *") // top of every hour
	after := time.Date(2024, 6, 1, 0, 30, 0, 0, time.UTC)

	got, err := predictor.NextN(expr, after, time.UTC, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 results, got %d", len(got))
	}

	expected := []time.Time{
		time.Date(2024, 6, 1, 1, 0, 0, 0, time.UTC),
		time.Date(2024, 6, 1, 2, 0, 0, 0, time.UTC),
		time.Date(2024, 6, 1, 3, 0, 0, 0, time.UTC),
	}
	for i, want := range expected {
		if !got[i].Equal(want) {
			t.Errorf("NextN[%d] = %v, want %v", i, got[i], want)
		}
	}
}
