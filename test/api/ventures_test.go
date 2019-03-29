package api

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	p "github.com/PaulioRandall/go-qlueless-assembly-api/internal/pkg"
	. "github.com/PaulioRandall/go-qlueless-assembly-api/internal/pkg/asserts"
	. "github.com/PaulioRandall/go-qlueless-assembly-api/internal/pkg/ventures"
)

var ventureDefaultMethods = []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}

// TODO: Craft some test data and pre-inject it into a SQLite database
func TestGET_Ventures(t *testing.T) {
	t.Log(`Given some Ventures already exist on the server
		When all Ventures are requested
		AND the HTTP method is 'GET'
		Then ensure the response code is 200
		And the 'Content-Type' header contains 'application/json'
		And 'Access-Control-Allow-Origin' is '*'
		And 'Access-Control-Allow-Headers' is '*'
		And 'Access-Control-Allow-Methods' only contains GET, POST, PUT, DELETE, HEAD, and OPTIONS
		And the body is a JSON array of valid Ventures
		...`)

	req := APICall{
		URL:    "http://localhost:8080/ventures",
		Method: "GET",
	}
	res := req.fire()
	defer res.Body.Close()
	defer PrintResponse(t, res.Body)

	require.Equal(t, 200, res.StatusCode)
	assertDefaultHeaders(t, res, "application/json", ventureDefaultMethods)

	var ven []Venture
	err := json.NewDecoder(res.Body).Decode(&ven)
	require.Nil(t, err)
	AssertGenericVentures(t, ven)
}

// TODO: Craft some test data and pre-inject it into a SQLite database
func TestGET_Venture_1(t *testing.T) {
	t.Log(`Given some Ventures already exist on the server
		When a specific existing Venture is requested
		AND the HTTP method is 'GET'
		Then ensure the response code is 200
		And the 'Content-Type' header contains 'application/json'
		And 'Access-Control-Allow-Origin' is '*'
		And 'Access-Control-Allow-Headers' is '*'
		And 'Access-Control-Allow-Methods' only contains GET, POST, PUT, DELETE, HEAD, and OPTIONS
		And the body is a JSON object representing a valid Venture
		...`)

	req := APICall{
		URL:    "http://localhost:8080/ventures?id=1",
		Method: "GET",
	}
	res := req.fire()
	defer res.Body.Close()
	defer PrintResponse(t, res.Body)

	require.Equal(t, 200, res.StatusCode)
	assertDefaultHeaders(t, res, "application/json", ventureDefaultMethods)

	var ven Venture
	err := json.NewDecoder(res.Body).Decode(&ven)
	require.Nil(t, err)
	AssertGenericVenture(t, ven)
}

func TestGET_Venture_2(t *testing.T) {
	t.Log(`Given some Ventures already exist on the server
		When a specific non-existent Venture is requested
		AND the HTTP method is 'GET'
		Then ensure the response code is 404
		And the 'Content-Type' header contains 'application/json'
		And 'Access-Control-Allow-Origin' is '*'
		And 'Access-Control-Allow-Headers' is '*'
		And 'Access-Control-Allow-Methods' only contains GET, POST, PUT, DELETE, HEAD, and OPTIONS
		And the body is a JSON object representing an error response
		...`)

	req := APICall{
		URL:    "http://localhost:8080/ventures?id=999999",
		Method: "GET",
	}
	res := req.fire()
	defer res.Body.Close()
	defer PrintResponse(t, res.Body)

	require.Equal(t, 404, res.StatusCode)
	assertDefaultHeaders(t, res, "application/json", ventureDefaultMethods)

	var reply p.WrappedReply
	err := json.NewDecoder(res.Body).Decode(&reply)
	require.Nil(t, err)
	AssertGenericError(t, reply)
}

func TestPOST_Venture_1(t *testing.T) {
	t.Log(`Given some Ventures already exist on the server
		When a new valid Venture is posted
		AND the HTTP method is 'POST'
		Then ensure the response code is 201
		And the 'Content-Type' header contains 'application/json'
		And 'Access-Control-Allow-Origin' is '*'
		And 'Access-Control-Allow-Headers' is '*'
		And 'Access-Control-Allow-Methods' only contains GET, POST, PUT, DELETE, HEAD, and OPTIONS
		And the body is a JSON object representing the living input Venture with a new assigned ID
		...`)

	input := Venture{
		Description: "A new Venture",
		State:       "Not started",
		OrderIDs:    "1,2,3",
	}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(&input)

	req := APICall{
		URL:    "http://localhost:8080/ventures",
		Method: "POST",
		Body:   buf,
	}
	res := req.fire()
	defer res.Body.Close()
	defer PrintResponse(t, res.Body)

	require.Equal(t, 201, res.StatusCode)
	assertDefaultHeaders(t, res, "application/json", ventureDefaultMethods)

	var output Venture
	err := json.NewDecoder(res.Body).Decode(&output)
	require.Nil(t, err)
	AssertGenericVenture(t, output)

	input.ID = output.ID
	input.IsAlive = true
	assert.Equal(t, input, output)
}

func TestPOST_Venture_2(t *testing.T) {
	t.Log(`Given some Ventures already exist on the server
		When a new invalid Venture is posted
		AND the HTTP method is 'POST'
		Then ensure the response code is 400
		And the 'Content-Type' header contains 'application/json'
		And 'Access-Control-Allow-Origin' is '*'
		And 'Access-Control-Allow-Headers' is '*'
		And 'Access-Control-Allow-Methods' only contains GET, POST, PUT, DELETE, HEAD, and OPTIONS
		And the body is a JSON object representing an error response
		...`)

	input := Venture{
		Description: "",
		State:       "",
		OrderIDs:    "invalid",
	}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(&input)

	req := APICall{
		URL:    "http://localhost:8080/ventures",
		Method: "POST",
		Body:   buf,
	}
	res := req.fire()
	defer res.Body.Close()
	defer PrintResponse(t, res.Body)

	require.Equal(t, 400, res.StatusCode)
	assertDefaultHeaders(t, res, "application/json", ventureDefaultMethods)

	var reply p.WrappedReply
	err := json.NewDecoder(res.Body).Decode(&reply)
	require.Nil(t, err)
	AssertGenericError(t, reply)
}

// TODO: Write more POST tests
// TODO: Write HEAD and OPTIONS tests
