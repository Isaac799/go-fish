// Package table provides a htmlx based table
// with optional multi column sorting, filtering, and pagination.
// Following the flow: define->populate->modify.
//
// 1: define the  html element and all the input fields it has
// (particularly hidden ones for state)
//
// 2: populate is where the request's form is compared
// against the form of the table, to retain previous state.
//
// 3: modify enables us to alter the element as needed based on
// state change e.g.: sort icon arrow direction
//
// Using this flow, I have pre made functions provide the define
// and modify fns that can modify a table.
package table

import (
	"errors"

	"github.com/Isaac799/go-fish/pkg/bridge"
)

var (
	// ErrMissingTableLimitOptions is provided if trying
	// to get table constraints and there was no limit options to
	// choose form in the table config.
	ErrMissingTableLimitOptions = errors.New("missing table limit options")
	// ErrMissingTablePage is given if doing modification or
	// analysis of a table that did not have pagination defined
	// properly
	ErrMissingTablePage = errors.New("missing table page hidden input el")
	// ErrMissingTableLimit is given if doing modification or
	// analysis of a table that did not have pagination defined
	// properly
	ErrMissingTableLimit = errors.New("missing table limit hidden input el")
	// ErrMissingTableHeaders is provided when making a table but no
	// headers (columns) are provided.
	ErrMissingTableHeaders = errors.New("missing table headers")
	// ErrMismatchRecordHeaderLength is provided when setting data
	// and the number of columns does not match the number of
	// headers the table was made with
	ErrMismatchRecordHeaderLength = errors.New("a row length did not align with expected header")
	// ErrMissingExpectedElement is a problem with define step and should
	// never be encountered unless an element has been screwed up by
	// the developer
	ErrMissingExpectedElement = errors.New("missing an expected element")
)

// HTMLTable just an html element with configuration
// for defining static behavior (e.g. hx post target), and state
// to define dynamic behavior (e.g. record count).
type HTMLTable struct {
	El *bridge.HTMLElement

	// configuration of a table. Not meant to be changed after
	// table is initialized.
	conf Config

	// non static config below, since state != config

	// RecordCount is to be set if you want pagination
	// to behave most optimal.
	// Otherwise you can set an arbitrary limit.
	RecordCount int

	// below are qol to prevent re-find and are
	// used to store state to fill in from request
	pageHiddenEl  *bridge.HTMLInput
	limitHiddenEl *bridge.HTMLInput

	// below are keyed by col index
	sortInputs   map[int]*bridge.HTMLInput
	filterInputs map[int]*bridge.HTMLInput
}

// New will build am empty table. I recommend using NewHTMLTableConf.
// Configuration cannot change after set - however state can.
func New(conf Config) (*HTMLTable, error) {
	if conf.Headers == nil || len(conf.Headers) == 0 {
		return nil, ErrMissingTableHeaders
	}

	tr := bridge.HTMLElement{
		Tag:      "tr",
		Children: make([]bridge.HTMLElement, 0, len(conf.Headers)),
	}
	for _, col := range conf.Headers {
		// header
		th := bridge.HTMLElement{
			Tag:       "th",
			InnerText: col,
		}
		tr.Children = append(tr.Children, th)
	}

	tHead := bridge.HTMLElement{
		Tag:      "thead",
		Children: []bridge.HTMLElement{tr},
	}
	tBody := bridge.HTMLElement{
		Tag: "tbody",
	}

	table := bridge.HTMLElement{
		Tag:      "table",
		Children: []bridge.HTMLElement{tHead, tBody},
	}

	form := bridge.ElementWithState(&table, nil)
	form.EnsureAttributes()
	form.Attributes["id"] = conf.ID
	answer := HTMLTable{
		El:           form,
		sortInputs:   make(map[int]*bridge.HTMLInput, len(conf.Headers)),
		filterInputs: make(map[int]*bridge.HTMLInput, len(conf.Headers)),
		conf:         conf,
	}
	return &answer, nil
}

// Mod is a function that will modify a table.
// Modifications come in 2 pairs: Define and Modify
type Mod func(table *HTMLTable) error

// Modify will alter a table, often adding a feature
func (table *HTMLTable) Modify(mods ...Mod) error {
	for _, mod := range mods {
		err := mod(table)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetData populates the tbody. It should align with the
// headers the table was created with.
func (table *HTMLTable) SetData(records [][]string) error {
	body := table.El.FindFirst(bridge.LikeTag("tbody"))
	if body == nil {
		return ErrMissingExpectedElement
	}

	body.Children = make([]bridge.HTMLElement, 0, len(records)-1)
	for _, row := range records {
		if len(row) != len(table.conf.Headers) {
			return ErrMismatchRecordHeaderLength
		}

		tr := bridge.HTMLElement{
			Tag:      "tr",
			Children: make([]bridge.HTMLElement, 0, len(row)),
		}

		for _, col := range row {
			td := bridge.HTMLElement{
				Tag:       "td",
				InnerText: col,
			}
			tr.Children = append(tr.Children, td)
		}
		body.Children = append(body.Children, tr)
	}

	return nil
}
