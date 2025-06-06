package bridge

import (
	"bytes"
	"encoding/csv"
	"testing"
)

func TestTable(t *testing.T) {
	fishCSV := `id,name,habitat,price_usd,stock
1,Tuna,Marine,10.99,50
2,Anchovies,Marine,2.99,200
3,Sardines,Marine,3.49,180
4,Clown Fish,Marine,15.00,25`

	reader := bytes.NewReader([]byte(fishCSV))
	csvReader := csv.NewReader(reader)

	tableEl, err := NewTable(csvReader)
	if err != nil {
		t.Fatal(err)
	}

	thead := tableEl.Children[0]
	eq(t, len(thead.Children), 1)

	headRow := thead.Children[0]
	headCol := headRow.Children[1]
	eq(t, headCol.InnerText, "name")

	tbody := tableEl.Children[1]
	eq(t, len(tbody.Children), 4)

	anchoviesRow := tbody.Children[1]
	anchoviesCol := anchoviesRow.Children[1]
	eq(t, anchoviesCol.InnerText, "Anchovies")
}
