package validator_test

import (
	"testing"

	"github.com/yourorg/cronlens/internal/validator"
)

func TestValidate_ValidExpressions(t *testing.T) {
	cases := []struct {
		name string
		expr string
	}{
		{"every minute", "* * * * *"},
		{"specific time", "30 9 * * 1"},
		{"range", "0-30 8-18 * * *"},
		{"step", "*/5 * * * *"},
		{"list", "0 9,12,17 * * *"},
		{"range with step", "0 */2 1-15 * *"},
		{"full specific", "0 0 1 1 0"},
		{"day of week max", "0 0 * * 6"},
		{"month boundaries", "0 0 1 1 *"},
		{"month max", "0 0 31 12 *"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := validator.Validate(tc.expr)
			if !result.Valid {
				t.Errorf("expected valid, got errors: %s", result.Error())
			}
		})
	}
}

func TestValidate_InvalidExpressions(t *testing.T) {
	cases := []struct {
		name        string
		expr        string
		wantErrCount int
	}{
		{"too few fields", "* * * *", 1},
		{"too many fields", "* * * * * *", 1},
		{"minute out of range", "60 * * * *", 1},
		{"hour out of range", "* 24 * * *", 1},
		{"dom out of range", "* * 32 * *", 1},
		{"month out of range", "* * * 13 *", 1},
		{"dow out of range", "* * * * 7", 1},
		{"bad range inverted", "30-10 * * * *", 1},
		{"bad step zero", "*/0 * * * *", 1},
		{"non-integer", "abc * * * *", 1},
		{"multiple errors", "60 25 * * *", 2},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := validator.Validate(tc.expr)
			if result.Valid {
				t.Errorf("expected invalid expression %q to fail", tc.expr)
			}
			if len(result.Errors) != tc.wantErrCount {
				t.Errorf("expected %d error(s), got %d: %s",
					tc.wantErrCount, len(result.Errors), result.Error())
			}
		})
	}
}

func TestValidationError_Error(t *testing.T) {
	e := &validator.ValidationError{
		Field:  "minute",
		Value:  "99",
		Reason: "value 99 out of bounds [0,59]",
	}
	got := e.Error()
	if got == "" {
		t.Error("expected non-empty error string")
	}
}
