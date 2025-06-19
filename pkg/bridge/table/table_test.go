package table

import (
	"bytes"
	"encoding/csv"
	"testing"

	"github.com/Isaac799/go-fish/pkg/bridge"
)

func TestNew(t *testing.T) {
	headers := []string{"ID", "Name"}
	config := NewConfig(bridge.RandomID(), headers, "")
	tbl, err := New(config)
	if err != nil {
		t.Fatal(err)
	}

	tableEl := tbl.El.FindFirst(bridge.LikeTag("table"))

	thead := tableEl.Children[0]
	assert(t, len(thead.Children), 1)

	headRow := thead.Children[0]
	headCol := headRow.Children[1]
	assert(t, headCol.InnerText, "Name")

	tbody := tableEl.Children[1]
	assert(t, len(tbody.Children), 0)
}

func TestSetPage(t *testing.T) {
	tbl := mockTable()

	definePagination, _ := Pagination()

	err := tbl.Modify(definePagination)
	if err != nil {
		t.Fatal(err)
	}

	_ = tbl.SetPage(3)

	v, err := tbl.pageHiddenEl.ParseInt()
	if err != nil {
		t.Fatal(err)
	}
	assert(t, v, 3)
}

func TestSetLimit(t *testing.T) {
	tbl := mockTable()

	definePagination, _ := Pagination()

	err := tbl.Modify(definePagination)
	if err != nil {
		t.Fatal(err)
	}

	_ = tbl.SetLimit(1)

	v, err := tbl.limitHiddenEl.ParseIndexes()
	if err != nil {
		t.Fatal(err)
	}
	assert(t, len(v), 1)
	assert(t, v[0], 1)
}

func TestSetFilter(t *testing.T) {
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

	s, _ := tbl.filterInputs[mockColName].ParseString()
	assert(t, s, nameFilter)

	s, _ = tbl.filterInputs[mockColHabitat].ParseString()
	assert(t, s, habitatFilter)
}

func TestSetSort(t *testing.T) {
	tbl := mockTable()

	defineSort, _ := Sort(mockColName, mockColHabitat)

	err := tbl.Modify(defineSort)
	if err != nil {
		t.Fatal(err)
	}

	_ = tbl.SetSort(mockColName, SortAsc)
	_ = tbl.SetSort(mockColHabitat, SortDesc)

	v, _ := tbl.sortInputs[mockColName].ParseInt()
	assert(t, v, SortAsc)

	v, _ = tbl.sortInputs[mockColHabitat].ParseInt()
	assert(t, v, SortDesc)

}
func TestSetData(t *testing.T) {
	tbl := mockTable()

	fishCSVData := `1,Tuna,Marine,250.0,10.99,50
2,Anchovies,Marine,0.02,2.99,300
3,Sardines,Marine,0.15,3.49,220
4,Clownfish,Marine,0.25,15.00,25
5,Salmon,Freshwater/Marine,4.5,12.99,60
6,Halibut,Marine,30.0,14.50,18
7,Cod,Marine,12.0,11.75,35
8,Trout,Freshwater,2.5,9.99,40
9,Mackerel,Marine,1.0,6.99,80
10,Herring,Marine,0.5,4.25,150`

	reader2 := bytes.NewReader([]byte(fishCSVData))
	csvReader2 := csv.NewReader(reader2)
	records, _ := csvReader2.ReadAll()

	err := tbl.SetData(records)
	if err != nil {
		t.Fatal(err)
	}

	body := tbl.El.FindFirst(bridge.LikeTag("tbody"))
	rows := body.FindAll(bridge.LikeTag("tr"))
	assert(t, len(rows), 10)

	cf := rows[3].FindNth(2, bridge.LikeTag("td"))
	assert(t, cf.InnerText, "Clownfish")
}
