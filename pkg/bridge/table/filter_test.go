package table

import (
	"testing"
)

func TestFilter_Define(t *testing.T) {
	tbl := mockTable()

	defineFilter := Filter(mockColName, mockColHabitat)

	err := tbl.Modify(defineFilter)
	if err != nil {
		t.Fatal(err)
	}

	nameFilter := "dog"
	habitatFilter := "cat"

	_ = tbl.SetFilter(mockColName, nameFilter)
	_ = tbl.SetFilter(mockColHabitat, habitatFilter)

	c, err := tbl.Constraints()
	if err != nil {
		t.Fatal(err)
	}

	assert(t, c.Filter[mockColName], nameFilter)
	assert(t, c.Filter[mockColHabitat], habitatFilter)
}
