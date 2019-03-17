package pkg

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogIfErr___1(t *testing.T) {
	act := LogIfErr(nil)
	assert.False(t, act)
	// Output:
	//
}

func TestLogIfErr___2(t *testing.T) {
	var err error = errors.New("Computer says no!")
	act := LogIfErr(err)
	assert.True(t, act)
	// Output:
	// Computer says no!
}

func dummyWorkItems() *[]WorkItem {
	return &[]WorkItem{
		WorkItem{
			Description: "# Outline the saga\nCreate a rough outline of the new saga.",
			WorkItemID:  "1",
			TagID:       "mid",
			StatusID:    "in_progress",
		},
		WorkItem{
			Description:      "# Name the saga\nThink of a name for the saga.",
			WorkItemID:       "2",
			ParentWorkItemID: "1",
			TagID:            "mid",
			StatusID:         "potential",
		},
		WorkItem{
			Description:      "# Outline the first chapter",
			WorkItemID:       "3",
			ParentWorkItemID: "1",
			TagID:            "mid",
			StatusID:         "delivered",
			Additional:       "archive_note:Done but not a compelling start",
		},
		WorkItem{
			Description:      "# Outline the second chapter",
			WorkItemID:       "4",
			ParentWorkItemID: "1",
			TagID:            "mid",
			StatusID:         "in_progress",
		},
	}
}

// When given an empty string, true is returned
func TestIsBlank___1(t *testing.T) {
	act := IsBlank("")
	assert.True(t, act)
}

// When given a string with whitespace, true is returned
func TestIsBlank___2(t *testing.T) {
	act := IsBlank("\r\n \t\f")
	assert.True(t, act)
}

// When given a string with no whitespaces, false is returned
func TestIsBlank___3(t *testing.T) {
	act := IsBlank("Captain Vimes")
	assert.False(t, act)
}

// When a value is present, it is returned
func TestValueOrEmpty___1(t *testing.T) {
	m := make(map[string]interface{})
	m["key"] = "value"
	act := ValueOrEmpty(m, "key")
	assert.Equal(t, "value", act)
}

// When a value is not present, empty string is returned
func TestValueOrEmpty___2(t *testing.T) {
	m := make(map[string]interface{})
	m["key"] = "value"
	act := ValueOrEmpty(m, "responsibilities")
	assert.Empty(t, act)
}

// When given an empty string, false is returned
func TestIsWrapperProp___1(t *testing.T) {
	act := isWrapperProp("")
	assert.False(t, act)
}

// When given an invalid wrapper property, false is returned
func TestIsWrapperProp___2(t *testing.T) {
	act := isWrapperProp("invalid")
	assert.False(t, act)
}

// When given a valid wrapper property, true is returned
func TestIsWrapperProp___3(t *testing.T) {
	act := isWrapperProp("message")
	assert.True(t, act)
}

// When given an empty string, an error is returned
func TestWrapWith___1(t *testing.T) {
	act, err := wrapWith("")
	assert.Nil(t, act)
	assert.NotNil(t, err)
}

// When given a dot separated list containing an empty value, an error is
// returned
func TestWrapWith___2(t *testing.T) {
	act, err := wrapWith("message..self")
	assert.Nil(t, act)
	assert.NotNil(t, err)
}

// When given a dot separated list containing an invalid value, an error is
// returned
func TestWrapWith___3(t *testing.T) {
	act, err := wrapWith("invalid")
	assert.Nil(t, act)
	assert.NotNil(t, err)
}

// When given a dot separated list containing an invalid value, an error is
// returned
func TestWrapWith___4(t *testing.T) {
	act, err := wrapWith("message.invalid")
	assert.Nil(t, act)
	assert.NotNil(t, err)
}

// When given a dot separated list containing a single valid value, the value
// is returned within an ordered array
func TestWrapWith___5(t *testing.T) {
	act, err := wrapWith("message")
	assert.NotNil(t, act)
	assert.Nil(t, err)
	assert.Len(t, act, 1)
	assert.Equal(t, "message", act[0])
}

// When given a dot separated list containing a multiple valid values, the
// values are returned within an ordered array
func TestWrapWith___6(t *testing.T) {
	act, err := wrapWith("message.self")
	assert.NotNil(t, act)
	assert.Nil(t, err)
	assert.Len(t, act, 2)
	assert.Equal(t, "message", act[0])
	assert.Equal(t, "self", act[1])
}

// When given a request without a 'wrap_with' query parameter, no values should
// be returned
func TestWrapData___1(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/", nil)
	if err != nil {
		panic(err)
	}
	r := Reply{
		Req: req,
	}
	act, err := wrapData(&r)
	assert.Nil(t, act)
	assert.Nil(t, err)
}

// When given a request with an invalid 'wrap_with' query parameter, the values
// should be returned
func TestWrapData___2(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/?wrap_with=abc.efg", nil)
	if err != nil {
		panic(err)
	}
	r := Reply{
		Req: req,
	}
	act, err := wrapData(&r)
	assert.Nil(t, act)
	assert.NotNil(t, err)
}

// When given a request with a valid 'wrap_with' query parameter, the values
// should be returned
func TestWrapData___3(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/?wrap_with=message.self", nil)
	if err != nil {
		panic(err)
	}
	r := Reply{
		Req: req,
	}
	act, err := wrapData(&r)
	assert.NotNil(t, act)
	assert.Nil(t, err)
	assert.Len(t, act, 2)
	assert.Equal(t, "message", act[0])
	assert.Equal(t, "self", act[1])
}
