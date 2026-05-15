// Package main is the entry point for the cronlens CLI tool.
// It provides a command-line interface for parsing cron expressions,
// generating human-readable descriptions, and predicting next run times
// with full timezone awareness.
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/yourusername/cronlens/internal/humanizer"
	"github.com/yourusername/cronlens/internal/parser"
	"github.com/yourusername/cronlens/internal/predictor"
)

const usage = `cronlens — Human-readable cron expression parser and next-run predictor

Usage:
  cronlens [flags] <cron expression>

Examples:
  cronlens "*/5 * * * *"
  cronlens -n 5 "0 9 * * MON-FRI"
  cronlens -tz America/New_York "30 8 * * *"

Flags:
`

func main() {
	var (
		tzName = flag.String("tz", "UTC", "Timezone for next-run predictions (e.g. America/New_York)")
		next   = flag.Int("n", 1, "Number of upcoming run times to display")
		quiet  = flag.Bool("q", false, "Suppress human-readable description, only show next run times")
	)

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	expr := flag.Arg(0)

	parsed, err := parser.Parse(expr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing expression: %v\n", err)
		os.Exit(1)
	}

	loc, err := time.LoadLocation(*tzName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid timezone %q: %v\n", *tzName, err)
		os.Exit(1)
	}

	if !*quiet {
		desc := humanizer.Humanize(parsed)
		fmt.Printf("Expression : %s\n", expr)
		fmt.Printf("Description: %s\n", desc)
		fmt.Printf("Timezone   : %s\n", loc)
		fmt.Println()
	}

	n := *next
	if n < 1 {
		n = 1
	}

	now := time.Now().In(loc)
	runs, err := predictor.NextN(parsed, now, n)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error computing next runs: %v\n", err)
		os.Exit(1)
	}

	if n == 1 {
		fmt.Printf("Next run: %s\n", formatTime(runs[0]))
	} else {
		fmt.Printf("Next %s runs:\n", strconv.Itoa(n))
		for i, t := range runs {
			fmt.Printf("  %2d. %s\n", i+1, formatTime(t))
		}
	}
}

// formatTime returns a human-friendly representation of a time value,
// including the timezone abbreviation for clarity.
func formatTime(t time.Time) string {
	return t.Format("Mon, 02 Jan 2006 15:04:05 MST")
}
