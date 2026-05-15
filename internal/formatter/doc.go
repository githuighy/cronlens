// Package formatter provides output rendering for CronReport values.
//
// It supports three output formats:
//
//   - text  — human-friendly plain text (default)
//   - json  — machine-readable JSON
//   - table — ASCII table for terminal display
//
// Usage:
//
//	f := formatter.New(formatter.FormatText)
//	out, err := f.Render(report)
//
The TimeLayout field on Formatter can be overridden to change how
timestamps are displayed in text and table modes.
package formatter
