package table

import "github.com/Isaac799/go-fish/pkg/bridge"

// Constraints is the constraints placed on data
// such as pagination, filtering, and sorting. Specifically
// designed to aid in SQL queries. Maps are keyed by column index.
type Constraints struct {
	Page   int
	Limit  int
	Offset int

	Sort   map[int]int
	Filter map[int]string
}

// Constraints provides the constraints placed on data.
// To be used in fetching rows before modify table.
// Consider define step of Pagination, Sort, and Filter
// and filling the table form via request before use.
func (table *HTMLTable) Constraints() (*Constraints, error) {
	if table.conf.LimitOptions == nil {
		return nil, ErrMissingTableLimitOptions
	}
	if table.pageHiddenEl == nil {
		return nil, ErrMissingTablePage
	}
	if table.limitHiddenEl == nil {
		return nil, ErrMissingTableLimit
	}

	page := defaultPage
	limit := defaultLimit

	parsedPage, err := table.pageHiddenEl.ParseInt()
	if err == nil {
		page = parsedPage
	}

	selection, err := bridge.InputSelectedValue(table.limitHiddenEl, table.conf.LimitOptions)
	if err == nil && len(selection) == 1 {
		limit = selection[0].Value
	}

	dc := Constraints{
		Limit:  limit,
		Page:   page,
		Offset: (page - 1) * limit,
		Sort:   make(map[int]int, len(table.sortInputs)),
		Filter: make(map[int]string, len(table.filterInputs)),
	}

	for i := range table.filterInputs {
		v, err := table.filterInputs[i].ParseString()
		if err != nil {
			continue
		}
		dc.Filter[i] = v
	}
	for i := range table.sortInputs {
		v, err := table.sortInputs[i].ParseInt()
		if err != nil {
			continue
		}
		dc.Sort[i] = v
	}

	return &dc, nil
}
