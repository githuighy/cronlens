package parser

import (
	"fmt"
	"strconv"
	"strings"
)

// FieldType represents which position in a cron expression a field occupies.
type FieldType int

const (
	FieldMinute FieldType = iota
	FieldHour
	FieldDayOfMonth
	FieldMonth
	FieldDayOfWeek
)

// fieldMeta holds the valid range and human-readable name for a cron field.
type fieldMeta struct {
	Name string
	Min  int
	Max  int
}

var fieldMetas = map[FieldType]fieldMeta{
	FieldMinute:     {"minute", 0, 59},
	FieldHour:       {"hour", 0, 23},
	FieldDayOfMonth: {"day-of-month", 1, 31},
	FieldMonth:      {"month", 1, 12},
	FieldDayOfWeek:  {"day-of-week", 0, 6},
}

// Field represents a parsed cron field with its resolved set of values.
type Field struct {
	Type   FieldType
	Values []int // sorted, deduplicated list of matching values
	Raw    string
}

// ParseField parses a single cron field token (e.g. "*/5", "1-5", "3,7") for
// the given field type and returns a Field or an error.
func ParseField(raw string, ft FieldType) (*Field, error) {
	meta := fieldMetas[ft]
	values, err := resolveValues(raw, meta.Min, meta.Max)
	if err != nil {
		return nil, fmt.Errorf("field %s: %w", meta.Name, err)
	}
	return &Field{Type: ft, Values: values, Raw: raw}, nil
}

func resolveValues(token string, min, max int) ([]int, error) {
	set := map[int]struct{}{}
	for _, part := range strings.Split(token, ",") {
		if err := expandPart(part, min, max, set); err != nil {
			return nil, err
		}
	}
	result := make([]int, 0, len(set))
	for v := min; v <= max; v++ {
		if _, ok := set[v]; ok {
			result = append(result, v)
		}
	}
	return result, nil
}

func expandPart(part string, min, max int, set map[int]struct{}) error {
	step := 1
	if idx := strings.Index(part, "/"); idx != -1 {
		var err error
		step, err = strconv.Atoi(part[idx+1:])
		if err != nil || step < 1 {
			return fmt.Errorf("invalid step %q", part[idx+1:])
		}
		part = part[:idx]
	}
	var lo, hi int
	if part == "*" {
		lo, hi = min, max
	} else if idx := strings.Index(part, "-"); idx != -1 {
		var err error
		lo, err = strconv.Atoi(part[:idx])
		if err != nil {
			return fmt.Errorf("invalid range start %q", part[:idx])
		}
		hi, err = strconv.Atoi(part[idx+1:])
		if err != nil {
			return fmt.Errorf("invalid range end %q", part[idx+1:])
		}
	} else {
		v, err := strconv.Atoi(part)
		if err != nil {
			return fmt.Errorf("invalid value %q", part)
		}
		lo, hi = v, v
	}
	if lo < min || hi > max || lo > hi {
		return fmt.Errorf("value %d-%d out of range [%d-%d]", lo, hi, min, max)
	}
	for v := lo; v <= hi; v += step {
		set[v] = struct{}{}
	}
	return nil
}
