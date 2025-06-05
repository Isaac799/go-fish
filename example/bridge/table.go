package bridge

import (
	"encoding/csv"
)

// NewTable takes in a csv reader to build an HTML table.
// Giving th consumer of this function all the power on how to
// transform data before its rendered.
func NewTable(csvReader *csv.Reader) (*HTMLElement, error) {
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	table := HTMLElement{
		Tag:      "table",
		Children: make([]HTMLElement, 2),
	}
	tHead := HTMLElement{
		Tag:      "thead",
		Children: make([]HTMLElement, 1),
	}
	tBody := HTMLElement{
		Tag:      "tbody",
		Children: make([]HTMLElement, len(records)-1),
	}
	for y, row := range records {
		tr := HTMLElement{
			Tag:      "tr",
			Children: make([]HTMLElement, len(row)),
		}

		for x, col := range row {
			if y == 0 {
				// header
				th := HTMLElement{
					Tag:       "th",
					InnerText: col,
				}
				tr.Children[x] = th
				continue
			}
			// row
			td := HTMLElement{
				Tag:       "td",
				InnerText: col,
			}
			tr.Children[x] = td
		}
		if y == 0 {
			// header
			tHead.Children[0] = tr
			continue
		}
		// row
		tBody.Children[y-1] = tr
	}

	table.Children[0] = tHead
	table.Children[1] = tBody

	return &table, nil
}
