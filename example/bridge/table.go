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
		Children: make([]HTMLElement, 0, 2),
	}
	tHead := HTMLElement{
		Tag:      "thead",
		Children: make([]HTMLElement, 0, 1),
	}
	tBody := HTMLElement{
		Tag:      "tbody",
		Children: make([]HTMLElement, 0, len(records)-1),
	}
	for y, row := range records {
		tr := HTMLElement{
			Tag:      "tr",
			Children: make([]HTMLElement, 0, len(row)),
		}

		for _, col := range row {
			if y == 0 {
				// header
				th := HTMLElement{
					Tag:       "th",
					InnerText: col,
				}
				tr.Children = append(tr.Children, th)
				continue
			}
			// row
			td := HTMLElement{
				Tag:       "td",
				InnerText: col,
			}
			tr.Children = append(tr.Children, td)
		}
		if y == 0 {
			// header
			tHead.Children = append(tHead.Children, tr)
			continue
		}
		// row
		tBody.Children = append(tBody.Children, tr)
	}

	table.Children = append(table.Children, tHead, tBody)

	return &table, nil
}
