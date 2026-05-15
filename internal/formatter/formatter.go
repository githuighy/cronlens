// Package formatter provides structured output formatting for cron expression analysis.
package formatter

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// OutputFormat represents the desired output format.
type OutputFormat string

const (
	FormatText OutputFormat = "text"
	FormatJSON OutputFormat = "json"
	FormatTable OutputFormat = "table"
)

// CronReport holds all analyzed data for a cron expression.
type CronReport struct {
	Expression  string    `json:"expression"`
	Human       string    `json:"human_readable"`
	Timezone    string    `json:"timezone"`
	NextRuns    []time.Time `json:"next_runs"`
	Valid       bool      `json:"valid"`
	Errors      []string  `json:"errors,omitempty"`
}

// Formatter renders a CronReport in the specified format.
type Formatter struct {
	Format    OutputFormat
	TimeLayout string
}

// New creates a Formatter with sensible defaults.
func New(format OutputFormat) *Formatter {
	return &Formatter{
		Format:     format,
		TimeLayout: "2006-01-02 15:04:05 MST",
	}
}

// Render converts a CronReport to its string representation.
func (f *Formatter) Render(r *CronReport) (string, error) {
	switch f.Format {
	case FormatJSON:
		return f.renderJSON(r)
	case FormatTable:
		return f.renderTable(r), nil
	default:
		return f.renderText(r), nil
	}
}

func (f *Formatter) renderJSON(r *CronReport) (string, error) {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", fmt.Errorf("formatter: json marshal: %w", err)
	}
	return string(b), nil
}

func (f *Formatter) renderText(r *CronReport) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Expression : %s\n", r.Expression))
	sb.WriteString(fmt.Sprintf("Meaning    : %s\n", r.Human))
	sb.WriteString(fmt.Sprintf("Timezone   : %s\n", r.Timezone))
	if len(r.NextRuns) > 0 {
		sb.WriteString("Next runs  :\n")
		for i, t := range r.NextRuns {
			sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, t.Format(f.TimeLayout)))
		}
	}
	if len(r.Errors) > 0 {
		sb.WriteString("Errors     :\n")
		for _, e := range r.Errors {
			sb.WriteString(fmt.Sprintf("  - %s\n", e))
		}
	}
	return sb.String()
}

func (f *Formatter) renderTable(r *CronReport) string {
	var sb strings.Builder
	line := strings.Repeat("-", 50)
	sb.WriteString(line + "\n")
	sb.WriteString(fmt.Sprintf("| %-47s|\n", "Cron Expression Analysis"))
	sb.WriteString(line + "\n")
	sb.WriteString(fmt.Sprintf("| %-12s | %-33s|\n", "Expression", r.Expression))
	sb.WriteString(fmt.Sprintf("| %-12s | %-33s|\n", "Meaning", truncate(r.Human, 33)))
	sb.WriteString(fmt.Sprintf("| %-12s | %-33s|\n", "Timezone", r.Timezone))
	sb.WriteString(line + "\n")
	for i, t := range r.NextRuns {
		sb.WriteString(fmt.Sprintf("| Next #%-6d| %-33s|\n", i+1, t.Format("2006-01-02 15:04 MST")))
	}
	if len(r.NextRuns) > 0 {
		sb.WriteString(line + "\n")
	}
	return sb.String()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
