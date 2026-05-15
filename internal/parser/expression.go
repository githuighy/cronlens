package parser

import (
	"fmt"
	"strings"
)

// Expression holds all five parsed cron fields.
type Expression struct {
	Minute     *Field
	Hour       *Field
	DayOfMonth *Field
	Month      *Field
	DayOfWeek  *Field
	Raw        string
}

// Parse parses a standard 5-field cron expression string and returns an
// Expression or a descriptive error.
func Parse(expr string) (*Expression, error) {
	expr = strings.TrimSpace(expr)
	parts := strings.Fields(expr)
	if len(parts) != 5 {
		return nil, fmt.Errorf("expected 5 fields, got %d", len(parts))
	}

	orders := []struct {
		ft    FieldType
		token string
	}{
		{FieldMinute, parts[0]},
		{FieldHour, parts[1]},
		{FieldDayOfMonth, parts[2]},
		{FieldMonth, parts[3]},
		{FieldDayOfWeek, parts[4]},
	}

	fields := make([]*Field, 5)
	for i, o := range orders {
		f, err := ParseField(o.token, o.ft)
		if err != nil {
			return nil, fmt.Errorf("parse error at field %d (%s): %w", i+1, o.token, err)
		}
		fields[i] = f
	}

	return &Expression{
		Minute:     fields[0],
		Hour:       fields[1],
		DayOfMonth: fields[2],
		Month:      fields[3],
		DayOfWeek:  fields[4],
		Raw:        expr,
	}, nil
}

// Matches reports whether the given time components satisfy the expression.
func (e *Expression) Matches(minute, hour, dom, month, dow int) bool {
	return contains(e.Minute.Values, minute) &&
		contains(e.Hour.Values, hour) &&
		contains(e.DayOfMonth.Values, dom) &&
		contains(e.Month.Values, month) &&
		contains(e.DayOfWeek.Values, dow)
}

func contains(values []int, v int) bool {
	for _, x := range values {
		if x == v {
			return true
		}
	}
	return false
}
