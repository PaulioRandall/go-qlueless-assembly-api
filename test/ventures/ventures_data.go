package ventures

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	a "github.com/PaulioRandall/go-qlueless-assembly-api/internal/pkg/asserts"
	q "github.com/PaulioRandall/go-qlueless-assembly-api/internal/pkg/qserver"
	v "github.com/PaulioRandall/go-qlueless-assembly-api/internal/pkg/ventures"
	w "github.com/PaulioRandall/go-qlueless-assembly-api/internal/pkg/wrapped"
	test "github.com/PaulioRandall/go-qlueless-assembly-api/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var VenHttpMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
var _httpMethods = "GET, POST, PUT, DELETE, OPTIONS"
var dbPath string = ""
var venDB *sql.DB = nil

// BeginEmptyTest is run at the start of a test to setup the server but does
// not inject any test data.
func BeginEmptyTest(relServerPath string) {
	dbPath = relServerPath + "/qlueless.db"
	DBReset()
	test.StartServer(relServerPath)
}

// BeginTest is run at the start of every test to setup the server and
// inject the test data.
func BeginTest(relServerPath string) {
	dbPath = relServerPath + "/qlueless.db"
	DBReset()
	DBInjectLiving()
	DBInjectDead()
	test.StartServer(relServerPath)
}

// EndTest should be deferred straight after BeginTest() is run to
// close resources at the end of every test.
func EndTest() {
	test.StopServer()
	DBClose()
}

// _deleteIfExists deletes the file at the path specified if it exist.
func _deleteIfExists(path string) {
	err := os.Remove(path)
	switch {
	case err == nil, os.IsNotExist(err):
	default:
		log.Fatal(err)
	}
}

// DBReset will reset the database by closing and deleting it then
// creating a new one.
func DBReset() {
	DBClose()
	_deleteIfExists(dbPath)

	var err error
	venDB, err = q.OpenSQLiteDatabase(dbPath)
	if err != nil {
		panic(err)
	}

	err = v.CreateTables(venDB)
	if err != nil {
		panic(err)
	}
}

// DBClose closes the test database.
func DBClose() {
	if venDB != nil {
		err := venDB.Close()
		if err != nil {
			panic(err)
		}
	}
	venDB = nil
}

// DBInject injects a Venture into the database.
func DBInject(new v.NewVenture) *v.Venture {
	ven, ok := new.Insert(venDB)
	if !ok {
		panic("Already printed above!")
	}
	return ven
}

// DBInjectLiving injects a default set of living Ventures into the database
func DBInjectLiving() {
	DBInject(v.NewVenture{
		Description: "White wizard",
		State:       "Not started",
		Extra:       "colour: white; power: 9000",
	})
	DBInject(v.NewVenture{
		Description: "Green lizard",
		State:       "In progress",
	})
	DBInject(v.NewVenture{
		Description: "Pink gizzard",
		State:       "Finished",
	})
	DBInject(v.NewVenture{
		Description: "Eddie Izzard",
		State:       "In Progress",
	})
	DBInject(v.NewVenture{
		Description: "The Count of Tuscany",
		State:       "In Progress",
	})
}

// DBInjectDead injects a default set of dead Ventures into the database
func DBInjectDead() {
	s := []v.Venture{
		*DBInject(v.NewVenture{
			Description: "Rose",
			State:       "Finised",
		}),
		*DBInject(v.NewVenture{
			Description: "Lily",
			State:       "Closed",
		}),
	}

	mod := v.ModVenture{
		Props: "is_dead",
		Values: v.Venture{
			Dead: true,
		},
	}

	for _, ven := range s {
		mod.ApplyMod(&ven)
		err := ven.Update(venDB)
		if err != nil {
			panic(err)
		}
	}
}

// DBQueryAll queries the database for all living ventures
func DBQueryAll() []v.Venture {
	rows, err := venDB.Query(`
		SELECT id, last_modified, description, order_ids, state, extra
		FROM ql_venture
	`)

	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		panic(err)
	}

	return _mapRows(rows)
}

// DBQueryMany queries the database for Ventures with the specified IDs
func DBQueryMany(ids string) []v.Venture {
	rows, err := venDB.Query(fmt.Sprintf(`
		SELECT id, last_modified, description, order_ids, state, extra
		FROM ql_venture
		WHERE id IN (%s)
	`, ids))

	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		panic(err)
	}

	return _mapRows(rows)
}

// DBQueryOne queries the database for a specific Venture
func DBQueryOne(id string) v.Venture {
	vens := DBQueryMany(id)
	if len(vens) != 1 {
		panic("Expected a single venture from query")
	}
	return vens[0]
}

// DBQueryFirst queries the database for the first Venture encountered
func DBQueryFirst() *v.Venture {
	vens := DBQueryAll()
	if len(vens) > 0 {
		return &vens[0]
	}
	return nil
}

// _mapRows is a file private function that maps rows from a database query into
// a slice of Ventures.
func _mapRows(rows *sql.Rows) []v.Venture {
	vens := []v.Venture{}

	for rows.Next() {
		vens = append(vens, *_mapRow(rows))
	}

	return vens
}

// _mapRow is a file private function that maps a single row from a database
// query into a Venture.
func _mapRow(rows *sql.Rows) *v.Venture {
	ven := v.Venture{}
	err := rows.Scan(&ven.ID,
		&ven.LastModified,
		&ven.Description,
		&ven.Orders,
		&ven.State,
		&ven.Extra)

	if err != nil {
		panic(err)
	}
	return &ven
}

// AssertHeaders asserts that the expected headers have been supplied.
func AssertHeaders(t *testing.T, h http.Header) {
	a.AssertHeadersEquals(t, h, map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "*",
		"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
		"Content-Type":                 "application/json; charset=utf-8",
	})
}

// AssertGenericReply asserts that reading from an io.Reader produces a generic
// reply with the expected values present.
func AssertGenericReply(t *testing.T, r io.Reader) {
	gr, err := w.DecodeFromReader(r)
	require.Nil(t, err)
	assert.NotEmpty(t, gr.Message)
	assert.NotEmpty(t, gr.Self)
	assert.Empty(t, gr.Data)
}
