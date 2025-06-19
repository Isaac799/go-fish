package table

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Isaac799/go-fish/pkg/bridge"
)

var (
	errNotEqual = errors.New("a and b where not equal")
)

const (
	mockColID = iota
	mockColName
	mockColHabitat
	mockColAverage
	mockColPrice
	mockColStock
)

func mockTable() *HTMLTable {
	headers := []string{"ID", "Name", "Habitat", "Average Weight KG", "Price USD", "Stock"}
	hxPost := "table/_stateful_table"
	config := NewConfig(bridge.RandomID(), headers, hxPost)
	tbl, _ := New(config)
	return tbl
}

func assert[T comparable](t *testing.T, a, b T) {
	if a == b {
		return
	}
	fmt.Println(a, b)
	t.Fatal(errNotEqual)
}

func TestConstraint_Accuracy(t *testing.T) {
	tbl := mockTable()

	definePagination, _ := Pagination()
	defineSort, _ := Sort(mockColName, mockColHabitat)
	defineFilter := Filter(mockColName, mockColHabitat)

	err := tbl.Modify(defineSort, definePagination, defineFilter)
	if err != nil {
		t.Fatal(err)
	}

	nameFilter := "dog"
	habitatFilter := "cat"

	_ = tbl.SetPage(3)
	// DefaultPaginationLimitOptions second option is 50
	_ = tbl.SetLimit(1)
	_ = tbl.SetFilter(mockColName, nameFilter)
	_ = tbl.SetFilter(mockColHabitat, habitatFilter)
	_ = tbl.SetSort(mockColName, SortAsc)
	_ = tbl.SetSort(mockColHabitat, SortDesc)

	c, err := tbl.Constraints()
	if err != nil {
		t.Fatal(err)
	}

	assert(t, c.Page, 3)
	assert(t, c.Limit, 50)
	assert(t, c.Offset, 100)
	assert(t, c.Sort[mockColName], SortAsc)
	assert(t, c.Sort[mockColHabitat], SortDesc)
	assert(t, c.Filter[mockColName], nameFilter)
	assert(t, c.Filter[mockColHabitat], habitatFilter)
}
