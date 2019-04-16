package uhttp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	a "github.com/PaulioRandall/go-qlueless-assembly-api/internal/pkg/asserts"
	w "github.com/PaulioRandall/go-qlueless-assembly-api/internal/pkg/wrapped"

	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
)

// ****************************************************************************
// RelURL()
// ****************************************************************************

func TestRelURL_1(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/character?q=Nobby", nil)
	require.Nil(t, err)

	act := RelURL(req)
	assert.Equal(t, "/character?q=Nobby", act)
}

// ****************************************************************************
// CheckReplyMessage()
// ****************************************************************************

func TestCheckNotEmpty_1(t *testing.T) {
	req, res, _ := SetupRequest("/")
	act := CheckNotEmpty(res, req, "message", "message")
	assert.True(t, act)
}

func TestCheckNotEmpty_2(t *testing.T) {
	req, res, rec := SetupRequest("/")
	act := CheckNotEmpty(res, req, "message", "")
	require.False(t, act)
	assert.Equal(t, 500, rec.Code)
}

// ****************************************************************************
// PrepResponseData()
// ****************************************************************************

func TestPrepResponseData___1(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/", nil)
	require.Nil(t, err)

	act := PrepResponseData(req, nil, "ignored")
	assert.Nil(t, act)
}

func TestPrepResponseData___2(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/", nil)
	require.Nil(t, err)
	data := make(map[string]interface{})
	data["album"] = "As Daylight Dies"

	act := PrepResponseData(req, data, "ignored")
	require.NotNil(t, act)
	assert.Equal(t, data, act)
}

func TestPrepResponseData___3(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/?wrap", nil)
	require.Nil(t, err)

	data := make(map[string]interface{})
	data["album"] = "As Daylight Dies"

	exp := w.WrappedReply{
		Message: "Cheese",
		Self:    req.URL.String(),
		Data:    data,
	}

	act := PrepResponseData(req, data, "Cheese")
	require.NotNil(t, act)
	assert.Equal(t, exp, act)
}

// ****************************************************************************
// AppendCORSHeaders()
// ****************************************************************************

func TestAppendCORSHeaders___1(t *testing.T) {
	rec := httptest.NewRecorder()
	var res http.ResponseWriter = rec
	AppendCORSHeaders(&res, "*")
	a.AssertHeadersEquals(t, (*rec).Header(), map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "*",
		"Access-Control-Allow-Methods": "*",
	})
}

// ****************************************************************************
// AppendJSONHeader()
// ****************************************************************************

func TestAppendJSONHeaders___1(t *testing.T) {
	rec := httptest.NewRecorder()
	var res http.ResponseWriter = rec
	AppendJSONHeader(&res, "")
	a.AssertHeadersEquals(t, (*rec).Header(), map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	})
}
