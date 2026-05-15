package formatter_test

import (
	"testing"

	"github.com/user/cronlens/internal/formatter"
)

func TestParseFormat(t *testing.T) {
	cases := []struct {
		input    string
		want     formatter.OutputFormat
	}{
		{"json", formatter.FormatJSON},
		{"JSON", formatter.FormatJSON},
		{"Json", formatter.FormatJSON},
		{"table", formatter.FormatTable},
		{"TABLE", formatter.FormatTable},
		{"text", formatter.FormatText},
		{"TEXT", formatter.FormatText},
		{"", formatter.FormatText},
		{"unknown", formatter.FormatText},
		{" json ", formatter.FormatJSON},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got := formatter.ParseFormat(tc.input)
			if got != tc.want {
				t.Errorf("ParseFormat(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestValidFormats(t *testing.T) {
	formats := formatter.ValidFormats()
	if len(formats) != 3 {
		t.Errorf("expected 3 valid formats, got %d", len(formats))
	}
	seen := map[string]bool{}
	for _, f := range formats {
		seen[f] = true
	}
	for _, expected := range []string{"text", "json", "table"} {
		if !seen[expected] {
			t.Errorf("ValidFormats missing %q", expected)
		}
	}
}

func TestOutputFormat_String(t *testing.T) {
	if formatter.FormatJSON.String() != "json" {
		t.Errorf("expected \"json\", got %q", formatter.FormatJSON.String())
	}
}
