package ventures

import (
	"database/sql"
	"encoding/json"
	"io"
	"strconv"
	"strings"

	u "github.com/PaulioRandall/go-qlueless-assembly-api/internal/pkg/utils"
)

// NewVenture represents a new Venture.
type NewVenture struct {
	Description string `json:"description"`
	OrderIDs    string `json:"order_ids"`
	State       string `json:"state"`
	Extra       string `json:"extra"`
}

// DecodeNewVenture decodes a NewVenture from data obtained via a Reader
func DecodeNewVenture(r io.Reader) (NewVenture, error) {
	var v NewVenture
	d := json.NewDecoder(r)
	err := d.Decode(&v)
	return v, err
}

// Clean removes redundent whitespace from property values within a Venture
// except where whitespace is allowable.
func (nv *NewVenture) Clean() {
	nv.Description = strings.TrimSpace(nv.Description)
	nv.OrderIDs = u.StripWhitespace(nv.OrderIDs)
	nv.State = strings.TrimSpace(nv.State)
}

// Validate checks each field contains valid content returning a non-empty
// slice of human readable error messages detailing the violations found or an
// empty slice if all is well. These messages are suitable for returning to
// clients.
func (nv *NewVenture) Validate() []string {
	errMsgs := []string{}

	errMsgs = u.AppendIfEmpty(nv.Description, errMsgs,
		"Ventures must have a description.")

	if nv.OrderIDs != "" {
		errMsgs = u.AppendIfNotPositiveIntCSV(nv.OrderIDs, errMsgs,
			"Child OrderIDs within a Venture must all be positive integers.")
	}

	errMsgs = u.AppendIfEmpty(nv.State, errMsgs, "Ventures must have a state.")
	return errMsgs
}

// Insert inserts the NewVenture into the database
//
// @UNTESTED
func (nv *NewVenture) Insert(db *sql.DB) (*Venture, bool) {

	id, err := _findNextID(db)
	if u.LogIfErr(err) {
		return nil, false
	}

	stmt, err := db.Prepare(`INSERT INTO venture (
		id, description, order_ids, state, extra
	) VALUES (
		?, ?, ?, ?, ?
	);`)

	if stmt != nil {
		defer stmt.Close()
	}

	if u.LogIfErr(err) {
		return nil, false
	}

	_, err = nv._execInsert(id, stmt)
	if u.LogIfErr(err) {
		return nil, false
	}

	ven, err := QueryFor(db, id)
	if u.LogIfErr(err) {
		return nil, false
	}

	return ven, true
}

// _findNextID returns the next free Venture ID.
func _findNextID(db *sql.DB) (string, error) {
	stmt, err := db.Prepare(`SELECT COALESCE(MAX(id), 0) FROM venture;`)

	if stmt != nil {
		defer stmt.Close()
	}

	if err != nil {
		return "", err
	}

	var id int64
	err = stmt.QueryRow().Scan(&id)
	if err != nil {
		return "", err
	}

	id++
	r := strconv.FormatInt(id, 10)
	return r, nil
}

// _execInsert is a file private function that executes the supplied insert
// statement
func (nv *NewVenture) _execInsert(id string, stmt *sql.Stmt) (*Venture, error) {
	ven := Venture{
		ID:          id,
		Description: nv.Description,
		OrderIDs:    nv.OrderIDs,
		State:       nv.State,
		Extra:       nv.Extra,
	}

	_, err := stmt.Exec(ven.ID,
		ven.Description,
		ven.OrderIDs,
		ven.State,
		ven.Extra)

	if err != nil {
		return nil, err
	}

	return &ven, nil
}
