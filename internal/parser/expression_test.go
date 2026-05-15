package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_ValidExpressions(t *testing.T) {
	tests := []struct {
		expr        string
		wantMinutes []int
		wantHours   []int
	}{
		{"* * * * *", rangeTo(0, 59), rangeTo(0, 23)},
		{"0 * * * *", []int{0}, rangeTo(0, 23)},
		{"*/15 9-17 * * 1-5", []int{0, 15, 30, 45}, rangeTo(9, 17)},
		{"5,10,15 0 * * *", []int{5, 10, 15}, []int{0}},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			expr, err := Parse(tt.expr)
			require.NoError(t, err)
			assert.Equal(t, tt.wantMinutes, expr.Minute.Values)
			assert.Equal(t, tt.wantHours, expr.Hour.Values)
		})
	}
}

func TestParse_InvalidExpressions(t *testing.T) {
	tests := []string{
		"* * * *",         // too few fields
		"* * * * * *",     // too many fields
		"60 * * * *",      // minute out of range
		"* 25 * * *",      // hour out of range
		"* * 0 * *",       // dom out of range (min is 1)
		"* * * 13 *",      // month out of range
		"abc * * * *",     // non-numeric
		"*/0 * * * *",     // step of 0
	}
	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			_, err := Parse(tt)
			assert.Error(t, err)
		})
	}
}

func TestExpression_Matches(t *testing.T) {
	expr, err := Parse("30 9 * * 1-5")
	require.NoError(t, err)

	// Monday 09:30
	assert.True(t, expr.Matches(30, 9, 15, 6, 1))
	// Saturday 09:30 — weekend, should not match
	assert.False(t, expr.Matches(30, 9, 15, 6, 6))
	// Monday 09:00 — wrong minute
	assert.False(t, expr.Matches(0, 9, 15, 6, 1))
}

func rangeTo(lo, hi int) []int {
	result := make([]int, hi-lo+1)
	for i := range result {
		result[i] = lo + i
	}
	return result
}
