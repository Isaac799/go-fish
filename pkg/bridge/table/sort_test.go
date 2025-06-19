package table

import (
	"testing"

	"github.com/Isaac799/go-fish/pkg/bridge"
)

func TestSort_Define(t *testing.T) {
	tbl := mockTable()

	defineSort, _ := Sort(mockColName, mockColHabitat)

	err := tbl.Modify(defineSort)
	if err != nil {
		t.Fatal(err)
	}

	_ = tbl.SetSort(mockColName, SortAsc)
	_ = tbl.SetSort(mockColHabitat, SortDesc)

	c, err := tbl.Constraints()
	if err != nil {
		t.Fatal(err)
	}

	assert(t, c.Sort[mockColName], SortAsc)
	assert(t, c.Sort[mockColHabitat], SortDesc)
}

func TestSort_Modify(t *testing.T) {
	tbl := mockTable()

	defineSort, modifySort := Sort(mockColName, mockColHabitat)

	err := tbl.Modify(defineSort)
	if err != nil {
		t.Fatal(err)
	}

	_ = tbl.SetSort(mockColName, SortAsc)
	_ = tbl.SetSort(mockColHabitat, SortDesc)

	err = tbl.Modify(modifySort)
	if err != nil {
		t.Fatal(err)
	}

	sort := tbl.El.FindAll(bridge.LikeSubmitButton)
	assert(t, len(sort), 2)

	assert(t, sort[0].InnerText, "arrow_upward")
	assert(t, sort[1].InnerText, "arrow_downward")
}
