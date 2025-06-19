package table

import (
	"testing"

	"github.com/Isaac799/go-fish/pkg/bridge"
)

func TestPagination_Define(t *testing.T) {
	tbl := mockTable()

	definePagination, _ := Pagination()

	err := tbl.Modify(definePagination)
	if err != nil {
		t.Fatal(err)
	}

	_ = tbl.SetPage(3)
	// DefaultPaginationLimitOptions second option is 50
	_ = tbl.SetLimit(1)

	c, err := tbl.Constraints()
	if err != nil {
		t.Fatal(err)
	}

	assert(t, c.Page, 3)
	assert(t, c.Limit, 50)
	assert(t, c.Offset, 100)
}

func TestPagination_Modify(t *testing.T) {
	tbl := mockTable()

	definePagination, modifyPagination := Pagination()

	err := tbl.Modify(definePagination)
	if err != nil {
		t.Fatal(err)
	}

	tbl.RecordCount = 1000

	_ = tbl.SetPage(7)
	// DefaultPaginationLimitOptions second option is 50
	_ = tbl.SetLimit(1)

	err = tbl.Modify(modifyPagination)
	if err != nil {
		t.Fatal(err)
	}

	submitBtns := tbl.El.FindAll(bridge.LikeSubmitButton)
	assert(t, len(submitBtns), 4)

	// first
	assert(t, submitBtns[0].InnerText, "1")
	// previous
	assert(t, submitBtns[1].InnerText, "6")
	// next
	assert(t, submitBtns[2].InnerText, "8")
	// last
	assert(t, submitBtns[3].InnerText, "20")
}
