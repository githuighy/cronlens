package formatter_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/cronlens/internal/formatter"
)

func sampleReport() *formatter.CronReport {
	loc, _ := time.LoadLocation("UTC")
	t1 := time.Date(2024, 6, 1, 12, 0, 0, 0, loc)
	t2 := time.Date(2024, 6, 1, 13, 0, 0, 0, loc)
	return &formatter.CronReport{
		Expression: "0 12 * * *",
		Human:      "At 12:00 every day",
		Timezone:   "UTC",
		NextRuns:   []time.Time{t1, t2},
		Valid:      true,
	}
}

func TestRender_Text(t *testing.T) {
	f := formatter.New(formatter.FormatText)
	out, err := f.Render(sampleReport())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "0 12 * * *") {
		t.Error("text output missing expression")
	}
	if !strings.Contains(out, "At 12:00 every day") {
		t.Error("text output missing human description")
	}
	if !strings.Contains(out, "UTC") {
		t.Error("text output missing timezone")
	}
}

func TestRender_JSON(t *testing.T) {
	f := formatter.New(formatter.FormatJSON)
	out, err := f.Render(sampleReport())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded formatter.CronReport
	if err := json.Unmarshal([]byte(out), &decoded); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if decoded.Expression != "0 12 * * *" {
		t.Errorf("expected expression '0 12 * * *', got %q", decoded.Expression)
	}
	if !decoded.Valid {
		t.Error("expected valid=true in JSON output")
	}
}

func TestRender_Table(t *testing.T) {
	f := formatter.New(formatter.FormatTable)
	out, err := f.Render(sampleReport())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Cron Expression Analysis") {
		t.Error("table output missing header")
	}
	if !strings.Contains(out, "Next #1") {
		t.Error("table output missing next run entry")
	}
}

func TestRender_WithErrors(t *testing.T) {
	r := sampleReport()
	r.Valid = false
	r.Errors = []string{"field out of range", "invalid step"}
	f := formatter.New(formatter.FormatText)
	out, err := f.Render(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "field out of range") {
		t.Error("text output missing error detail")
	}
}
