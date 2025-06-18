package table

import (
	"errors"
	"strconv"

	"github.com/Isaac799/go-fish/pkg/bridge"
)

// Sort indications direction for ordering a column
const (
	SortNone = iota
	SortAsc
	SortDesc
)

const (
	defaultPage  = 1
	defaultLimit = 10
)

// random as to not overlap with consumer field names
var (
	// form keys for page and limit
	formKeyPaginationLimit = bridge.RandomID()
	formKeyPaginationPage  = bridge.RandomID()
)

var (
	// Icons used after sort state is determined from request
	icons = map[int]string{
		SortNone: "unfold_more",
		SortAsc:  "arrow_upward",
		SortDesc: "arrow_downward",
	}

	// The next value to set the sort button to after sort
	// is determined from request
	nextSort = map[int]int{
		SortNone: SortAsc,
		SortAsc:  SortDesc,
		SortDesc: SortNone,
	}

	// only sort values of this nature are acceptable
	// mismatches are ignored
	acceptableSort = []int{SortNone, SortAsc, SortDesc}
)

// LabeledValue is a simple value and label that satisfies
// fmt.Stringer. Used in limit select input
type LabeledValue struct {
	Label string
	Value int
}

func (s LabeledValue) String() string {
	return s.Label
}

// DefaultPaginationLimitOptions is a reasonable
// limit options to be used in a select for a table
var DefaultPaginationLimitOptions = []LabeledValue{
	{Label: "10", Value: 10},
	{Label: "50", Value: 50},
	{Label: "250", Value: 250},
	{Label: "1000", Value: 1000},
}

var (
	// ErrInvalidLimit is given when making pagination
	// and a limit is unacceptable, such as a 0 value
	ErrInvalidLimit = errors.New("invalid page limit")
)

// paginationConfig has info for pagination of a table
type paginationConfig struct {
	CurrentPage  int
	TotalPages   int
	TotalCount   int
	Limit        int
	PreviousPage int
	NextPage     int
	Offset       int
	ShowFirst    bool
	ShowPrevious bool
	ShowNext     bool
	ShowLast     bool
}

// newPaginationConf provides a pagination metadata
func newPaginationConf(limit, page, totalCount int) (*paginationConfig, error) {
	if limit < 1 {
		return nil, ErrInvalidLimit
	}
	totalPages := (totalCount + limit - 1) / limit
	if page > totalPages {
		page = totalPages
	}
	previousPage := page - 1
	nextPage := page + 1

	pagination := paginationConfig{
		CurrentPage:  page,
		TotalPages:   totalPages,
		TotalCount:   totalCount,
		Limit:        limit,
		PreviousPage: page - 1,
		NextPage:     nextPage,
		Offset:       (page - 1) * limit,
		ShowFirst:    page > 2,
		ShowPrevious: previousPage > 0,
		ShowNext:     nextPage <= totalPages,
		ShowLast:     totalPages > 1 && page < totalPages-1,
	}
	return &pagination, nil
}

// Element provides a div element based around htmlx to be used
// in a form. Uses my classes and material icons
func (p *paginationConfig) Element(hxTarget, formPageKey, hxPost string) bridge.HTMLElement {
	paginationDiv := bridge.HTMLElement{
		Tag: "div",
		Attributes: bridge.Attributes{
			"class": "fr fr-center g2",
		},
		Children: make([]bridge.HTMLElement, 0, 7),
	}

	spreadIcon := bridge.HTMLElement{
		Tag: "span",
		Attributes: bridge.Attributes{
			"class": "material-icons",
		},
		InnerText: "more_horiz",
	}

	if p.ShowFirst {
		pageFirst := bridge.HTMLElement{
			Tag: "button",
			Attributes: bridge.Attributes{
				"type":      "submit",
				"hx-post":   hxPost,
				"name":      formPageKey,
				"value":     strconv.FormatInt(1, 10),
				"hx-target": hxTarget,
				"hx-swap":   "outerHTML",
			},
			InnerText: strconv.FormatInt(1, 10),
		}
		paginationDiv.Children = append(paginationDiv.Children, pageFirst, spreadIcon)
	}

	if p.ShowPrevious {
		pagePrev := bridge.HTMLElement{
			Tag: "button",
			Attributes: bridge.Attributes{
				"type":      "submit",
				"hx-post":   hxPost,
				"name":      formPageKey,
				"value":     strconv.Itoa(p.PreviousPage),
				"hx-target": hxTarget,
				"hx-swap":   "outerHTML",
			},
			InnerText: strconv.Itoa(p.PreviousPage),
		}
		paginationDiv.Children = append(paginationDiv.Children, pagePrev)
	}

	pageCurr := bridge.HTMLElement{
		Tag: "button",
		Attributes: bridge.Attributes{
			"class":           "is-current-page",
			"hx-disabled-elt": "this",
		},
		InnerText: strconv.Itoa(p.CurrentPage),
	}
	paginationDiv.Children = append(paginationDiv.Children, pageCurr)

	if p.ShowNext {
		pageNext := bridge.HTMLElement{
			Tag: "button",
			Attributes: bridge.Attributes{
				"type":      "submit",
				"hx-post":   hxPost,
				"name":      formPageKey,
				"value":     strconv.Itoa(p.NextPage),
				"hx-target": hxTarget,
				"hx-swap":   "outerHTML",
			},
			InnerText: strconv.Itoa(p.NextPage),
		}
		paginationDiv.Children = append(paginationDiv.Children, pageNext)
	}

	if p.ShowLast {
		pageLast := bridge.HTMLElement{
			Tag: "button",
			Attributes: bridge.Attributes{
				"type":      "submit",
				"hx-post":   hxPost,
				"name":      formPageKey,
				"value":     strconv.Itoa(p.TotalPages),
				"hx-target": hxTarget,
				"hx-swap":   "outerHTML",
			},
			InnerText: strconv.Itoa(p.TotalPages),
		}
		paginationDiv.Children = append(paginationDiv.Children, spreadIcon, pageLast)
	}

	return paginationDiv
}

// Pagination provides the Define and Modify functions needed
// enable pagination of a table
func Pagination() (Mod, Mod) {
	var defineFn Mod = func(table *HTMLTable) error {
		if table.conf.LimitOptions == nil {
			return ErrMissingTableLimitOptions
		}
		limitDiv := bridge.NewInputSelect(formKeyPaginationLimit, table.conf.LimitOptions)
		limitDiv.DeleteFirst(bridge.LikeTag("label"))

		limitSel := limitDiv.FindFirst(bridge.LikeInput)
		rowCountHTMLX := map[string]string{
			"hx-trigger": "change",
			"hx-post":    table.conf.HxPost,
			"hx-target":  table.conf.HxSwapTarget,
			"hx-swap":    "outerHTML",
		}
		limitSel.GiveAttributes(rowCountHTMLX)

		// store in config for easier access later
		pageHiddenEl := bridge.NewInputHidden(formKeyPaginationPage, "")
		table.pageHiddenEl = &pageHiddenEl
		table.limitHiddenEl = &limitDiv

		table.El.Children = append(table.El.Children, pageHiddenEl, limitDiv)

		return nil
	}

	var modifyFn Mod = func(table *HTMLTable) error {
		if table.conf.LimitOptions == nil {
			return ErrMissingTableLimitOptions
		}
		if table.pageHiddenEl == nil {
			return ErrMissingTablePage
		}
		if table.limitHiddenEl == nil {
			return ErrMissingTableLimit
		}

		var page int = defaultPage
		parsedCurrent, err := table.pageHiddenEl.ParseInt()
		if err == nil {
			page = parsedCurrent
		}

		var limit int = defaultLimit
		parsedLimit, err := bridge.InputSelectedValue(table.limitHiddenEl, table.conf.LimitOptions)
		if err == nil && len(parsedLimit) == 1 {
			limit = parsedLimit[0].Value
		}

		pagination, err := newPaginationConf(limit, page, table.RecordCount)
		if err != nil {
			return err
		}
		paginationEl := pagination.Element(table.conf.HxSwapTarget, formKeyPaginationPage, table.conf.HxPost)
		table.El.Children = append(table.El.Children, paginationEl)

		return nil
	}

	return defineFn, modifyFn
}
