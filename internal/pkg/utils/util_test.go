package utils

import (
	"errors"
	"testing"

	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
)

// ****************************************************************************
// LogIfErr()
// ****************************************************************************

func TestLogIfErr_1(t *testing.T) {
	act := LogIfErr(nil)
	assert.False(t, act)
	// Output:
	//
}

func TestLogIfErr_2(t *testing.T) {
	err := errors.New("Computer says no!")
	act := LogIfErr(err)
	assert.True(t, act)
	// Output:
	// [ERROR] Computer says no!
}

// ****************************************************************************
// WarnIfErr()
// ****************************************************************************

func TestWarnIfErr_1(t *testing.T) {
	act := WarnIfErr(nil)
	assert.False(t, act)
	// Output:
	//
}

func TestWarnIfErr_2(t *testing.T) {
	err := errors.New("Computer says no!")
	act := WarnIfErr(err)
	assert.True(t, act)
	// Output:
	// [warning] Computer says no!
}

// ****************************************************************************
// StripWhitespace()
// ****************************************************************************

func TestStripWhitespace_1(t *testing.T) {
	act := StripWhitespace("Rince \n\t\f\r wind")
	assert.Equal(t, "Rincewind", act)
}

func TestStripWhitespace_2(t *testing.T) {
	act := StripWhitespace("\t \n\t \r\n\n\fRincewind")
	assert.Equal(t, "Rincewind", act)
}

func TestStripWhitespace_3(t *testing.T) {
	act := StripWhitespace("Rincewind\r\n \t\t\t\t \f \r  \v\v")
	assert.Equal(t, "Rincewind", act)
}

func TestStripWhitespace_4(t *testing.T) {
	act := StripWhitespace("\r\nRi \tn\tc\t\t ew\f \r  in\vd\v")
	assert.Equal(t, "Rincewind", act)
}

func TestStripWhitespace_5(t *testing.T) {
	act := StripWhitespace("Rincewind")
	assert.Equal(t, "Rincewind", act)
}

func TestStripWhitespace_6(t *testing.T) {
	act := StripWhitespace("")
	assert.Equal(t, "", act)
}

func TestStripWhitespace_7(t *testing.T) {
	act := StripWhitespace("\r\n \t\t \t\t \f \r  \v\v  ")
	assert.Equal(t, "", act)
}

// ****************************************************************************
// AppendIfEmpty()
// ****************************************************************************

func TestAppendIfEmpty_1(t *testing.T) {
	act := AppendIfEmpty("", []string{}, "abc")
	require.Len(t, act, 1)
	assert.Contains(t, act, "abc")
}

func TestAppendIfEmpty_2(t *testing.T) {
	act := AppendIfEmpty("", []string{"xyz"}, "abc")
	require.Len(t, act, 2)
	assert.Contains(t, act, "xyz")
	assert.Contains(t, act, "abc")
}

func TestAppendIfEmpty_3(t *testing.T) {
	act := AppendIfEmpty("NOT-EMPTY", []string{}, "abc")
	assert.Len(t, act, 0)
}

// ****************************************************************************
// AppendIfNotPositiveInt()
// ****************************************************************************

func TestAppendIfNotPositiveInt___1(t *testing.T) {
	act := AppendIfNotPositiveInt("5", []string{}, "abc")
	assert.Len(t, act, 0)
}

func TestAppendIfNotPositiveInt___2(t *testing.T) {
	act := AppendIfNotPositiveInt("0", []string{}, "abc")
	require.Len(t, act, 1)
	assert.Contains(t, act, "abc")
}

func TestAppendIfNotPositiveInt___3(t *testing.T) {
	act := AppendIfNotPositiveInt("-5", []string{}, "abc")
	require.Len(t, act, 1)
	assert.Contains(t, act, "abc")
}

func TestAppendIfNotPositiveInt___4(t *testing.T) {
	act := []string{}
	act = AppendIfNotPositiveInt("-1", act, "abc")
	act = AppendIfNotPositiveInt("-1", act, "efg")
	act = AppendIfNotPositiveInt("-1", act, "hij")
	require.Len(t, act, 3)
	assert.Contains(t, act, "abc")
	assert.Contains(t, act, "efg")
	assert.Contains(t, act, "hij")
}

// ****************************************************************************
// IsPositiveInt()
// ****************************************************************************

func TestIsPositiveIntCSV_1(t *testing.T) {
	act := IsPositiveIntCSV("5")
	assert.True(t, act)
}

func TestIsPositiveIntCSV_2(t *testing.T) {
	act := IsPositiveIntCSV("1,2,3,4")
	assert.True(t, act)
}

func TestIsPositiveIntCSV_3(t *testing.T) {
	act := IsPositiveIntCSV("")
	assert.False(t, act)
}

func TestIsPositiveIntCSV_4(t *testing.T) {
	act := IsPositiveIntCSV("abc")
	assert.False(t, act)
}

func TestIsPositiveIntCSV_5(t *testing.T) {
	act := IsPositiveIntCSV("abc,efg,xyz")
	assert.False(t, act)
}

// ****************************************************************************
// AppendIfNotPositiveIntCSV()
// ****************************************************************************

func TestAppendIfNotPositiveIntCSV_1(t *testing.T) {
	act := AppendIfNotPositiveIntCSV("1,2,99,4,3", []string{}, "abc")
	assert.Len(t, act, 0)
}

func TestAppendIfNotPositiveIntCSV_2(t *testing.T) {
	act := AppendIfNotPositiveIntCSV("4", []string{}, "abc")
	assert.Len(t, act, 0)
}

func TestAppendIfNotPositiveIntCSV_3(t *testing.T) {
	act := AppendIfNotPositiveIntCSV("0", []string{}, "abc")
	assert.Len(t, act, 1)
	assert.Equal(t, "abc", act[0])
}

func TestAppendIfNotPositiveIntCSV_4(t *testing.T) {
	act := AppendIfNotPositiveIntCSV("-99", []string{}, "abc")
	assert.Len(t, act, 1)
	assert.Equal(t, "abc", act[0])
}

func TestAppendIfNotPositiveIntCSV_5(t *testing.T) {
	act := AppendIfNotPositiveIntCSV("3,2,1,0", []string{}, "abc")
	assert.Len(t, act, 1)
	assert.Equal(t, "abc", act[0])
}

func TestAppendIfNotPositiveIntCSV_6(t *testing.T) {
	act := AppendIfNotPositiveIntCSV(",1,2", []string{}, "abc")
	assert.Len(t, act, 1)
	assert.Equal(t, "abc", act[0])
}

func TestAppendIfNotPositiveIntCSV_7(t *testing.T) {
	act := AppendIfNotPositiveIntCSV("", []string{}, "abc")
	assert.Len(t, act, 1)
	assert.Equal(t, "abc", act[0])
}
