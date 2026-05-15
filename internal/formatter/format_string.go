package formatter

import "strings"

// ParseFormat converts a string flag value into an OutputFormat.
// It is case-insensitive and returns FormatText for unrecognised values.
func ParseFormat(s string) OutputFormat {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "json":
		return FormatJSON
	case "table":
		return FormatTable
	default:
		return FormatText
	}
}

// ValidFormats returns all recognised format names for use in help text.
func ValidFormats() []string {
	return []string{string(FormatText), string(FormatJSON), string(FormatTable)}
}

// String implements the fmt.Stringer interface.
func (f OutputFormat) String() string {
	return string(f)
}
