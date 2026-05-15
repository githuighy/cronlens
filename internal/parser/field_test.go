package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseField_Wildcard(t *testing.T) {
	f, err := ParseField("*", FieldMinute)
	require.NoError(t, err)
	assert.Equal(t, 60, len(f.Values)) // 0-59
	assert.Equal(t, 0, f.Values[0])
	assert.Equal(t, 59, f.Values[59])
}

func TestParseField_Step(t *testing.T) {
	f, err := ParseField("*/15", FieldMinute)
	require.NoError(t, err)
	assert.Equal(t, []int{0, 15, 30, 45}, f.Values)
}

func TestParseField_Range(t *testing.T) {
	f, err := ParseField("9-17", FieldHour)
	require.NoError(t, err)
	assert.Equal(t, []int{9, 10, 11, 12, 13, 14, 15, 16, 17}, f.Values)
}

func TestParseField_List(t *testing.T) {
	f, err := ParseField("1,3,5", FieldDayOfWeek)
	require.NoError(t, err)
	assert.Equal(t, []int{1, 3, 5}, f.Values)
}

func TestParseField_RangeWithStep(t *testing.T) {
	f, err := ParseField("0-30/10", FieldMinute)
	require.NoError(t, err)
	assert.Equal(t, []int{0, 10, 20, 30}, f.Values)
}

func TestParseField_SingleValue(t *testing.T) {
	f, err := ParseField("7", FieldMonth)
	require.NoError(t, err)
	assert.Equal(t, []int{7}, f.Values)
}

func TestParseField_OutOfRange(t *testing.T) {
	_, err := ParseField("99", FieldMinute)
	assert.Error(t, err)
}

func TestParseField_InvalidStep(t *testing.T) {
	_, err := ParseField("*/abc", FieldHour)
	assert.Error(t, err)
}

func TestParseField_InvalidRange(t *testing.T) {
	_, err := ParseField("10-5", FieldMinute) // lo > hi
	assert.Error(t, err)
}
