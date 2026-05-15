// Package validator provides cron expression validation with detailed error reporting.
package validator

import (
	"fmt"
	"strings"
)

// FieldConstraint defines the allowed range for a cron field.
type FieldConstraint struct {
	Name string
	Min  int
	Max  int
}

// Standard cron field constraints.
var fieldConstraints = []FieldConstraint{
	{Name: "minute", Min: 0, Max: 59},
	{Name: "hour", Min: 0, Max: 23},
	{Name: "day-of-month", Min: 1, Max: 31},
	{Name: "month", Min: 1, Max: 12},
	{Name: "day-of-week", Min: 0, Max: 6},
}

// ValidationError holds a field-level validation failure.
type ValidationError struct {
	Field   string
	Value   string
	Reason  string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid %s field %q: %s", e.Field, e.Value, e.Reason)
}

// Result holds the outcome of a validation run.
type Result struct {
	Valid  bool
	Errors []*ValidationError
}

func (r *Result) Error() string {
	msgs := make([]string, len(r.Errors))
	for i, e := range r.Errors {
		msgs[i] = e.Error()
	}
	return strings.Join(msgs, "; ")
}

// Validate checks a raw cron expression string and returns a Result.
func Validate(expr string) *Result {
	result := &Result{Valid: true}

	fields := strings.Fields(expr)
	if len(fields) != 5 {
		result.Valid = false
		result.Errors = append(result.Errors, &ValidationError{
			Field:  "expression",
			Value:  expr,
			Reason: fmt.Sprintf("expected 5 fields, got %d", len(fields)),
		})
		return result
	}

	for i, raw := range fields {
		constraint := fieldConstraints[i]
		if errs := validateField(raw, constraint); len(errs) > 0 {
			result.Valid = false
			result.Errors = append(result.Errors, errs...)
		}
	}

	return result
}

func validateField(raw string, c FieldConstraint) []*ValidationError {
	if raw == "*" {
		return nil
	}

	var errs []*ValidationError
	parts := strings.Split(raw, ",")
	for _, part := range parts {
		if e := validatePart(part, c); e != nil {
			errs = append(errs, e)
		}
	}
	return errs
}

func validatePart(part string, c FieldConstraint) *ValidationError {
	// Handle step: */n or range/n
	if strings.Contains(part, "/") {
		segments := strings.SplitN(part, "/", 2)
		step, err := parseInt(segments[1])
		if err != nil || step < 1 {
			return &ValidationError{Field: c.Name, Value: part, Reason: "step must be a positive integer"}
		}
		if segments[0] == "*" {
			return nil
		}
		part = segments[0]
	}

	// Handle range: a-b
	if strings.Contains(part, "-") {
		segments := strings.SplitN(part, "-", 2)
		lo, err1 := parseInt(segments[0])
		hi, err2 := parseInt(segments[1])
		if err1 != nil || err2 != nil {
			return &ValidationError{Field: c.Name, Value: part, Reason: "range bounds must be integers"}
		}
		if lo < c.Min || hi > c.Max {
			return &ValidationError{Field: c.Name, Value: part,
				Reason: fmt.Sprintf("range %d-%d out of bounds [%d,%d]", lo, hi, c.Min, c.Max)}
		}
		if lo > hi {
			return &ValidationError{Field: c.Name, Value: part, Reason: "range start must not exceed end"}
		}
		return nil
	}

	// Plain integer
	v, err := parseInt(part)
	if err != nil {
		return &ValidationError{Field: c.Name, Value: part, Reason: "value must be an integer"}
	}
	if v < c.Min || v > c.Max {
		return &ValidationError{Field: c.Name, Value: part,
			Reason: fmt.Sprintf("value %d out of bounds [%d,%d]", v, c.Min, c.Max)}
	}
	return nil
}

func parseInt(s string) (int, error) {
	var v int
	_, err := fmt.Sscanf(s, "%d", &v)
	return v, err
}
