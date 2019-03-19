package orders

import (
	"testing"

	. "github.com/PaulioRandall/go-qlueless-assembly-api/internal/pkg"
	"github.com/stretchr/testify/assert"
)

// When a valid map is provided, a Thing is returned
func TestMapToOrder___1(t *testing.T) {
	m := make(map[string]interface{})
	m["description"] = "description"

	act := mapToOrder(m)
	assert.Equal(t, "description", act.Description)
	assert.Empty(t, act.ID)
	assert.Empty(t, act.ParentID)
	assert.Empty(t, act.StatusID)
	assert.Empty(t, act.Additional)
}

// When a valid map is provided, a Thing is returned
func TestMapToOrder___2(t *testing.T) {
	m := make(map[string]interface{})
	m["description"] = "description"
	m["id"] = "id"
	m["parent_id"] = "parent_id"
	m["status_id"] = "status_id"
	m["additional"] = "abc: xyz; colour: black"

	act := mapToOrder(m)
	assert.Equal(t, "description", act.Description)
	assert.Equal(t, "id", act.ID)
	assert.Equal(t, "parent_id", act.ParentID)
	assert.Equal(t, "status_id", act.StatusID)
	assert.Equal(t, "abc: xyz; colour: black", act.Additional)
}

func createDummyOrder() Thing {
	return Thing{
		Description: "# Outline the saga\nCreate a rough outline of the new saga.",
		StatusID:    "in_progress",
	}
}

// When given an order, returns an order ID
func TestAddOrder___1(t *testing.T) {
	o := createDummyOrder()
	act, err := addOrder(o)
	assert.Nil(t, err)
	assert.NotEmpty(t, act)
	o.ID = act

	stored, ok := orders[o.ID]
	assert.True(t, ok)
	assert.Equal(t, o, stored)
}

// When invoked twice with the same order, returns differnt IDs each time
func TestAddOrder___2(t *testing.T) {
	a := createDummyOrder()
	ID_1, err := addOrder(a)
	assert.Nil(t, err)
	a.ID = ID_1

	b := createDummyOrder()
	ID_2, err := addOrder(b)
	assert.Nil(t, err)
	b.ID = ID_2

	assert.NotEqual(t, ID_1, ID_2)

	stored, ok := orders[a.ID]
	assert.True(t, ok)
	assert.Equal(t, a, stored)

	stored, ok = orders[b.ID]
	assert.True(t, ok)
	assert.Equal(t, b, stored)
}
