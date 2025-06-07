package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"github.com/Isaac799/go-fish/pkg/bridge"
)

type user struct {
	ID        int
	FirstName string
	LastName  string
}

type tables struct {
	Basic    bridge.HTMLElement
	Stateful bridge.HTMLElement
}

type fishData struct {
	Season string
	User   *user
	Table  *tables
	Form   *bridge.HTMLElement
}

func queriedSeason(r *http.Request) *fishData {
	data := fishData{}
	season, ok := r.Context().Value(queryCtxKey).(string)
	if !ok {
		return nil
	}
	data.Season = season
	return &data
}

func userInfo(r *http.Request) *fishData {
	data := fishData{}
	user, ok := r.Context().Value(userCtxKey).(user)
	if !ok {
		return nil
	}
	data.User = &user
	return &data
}

type labeledValue struct {
	label string
	value uint64
}

func (s labeledValue) Print() string {
	return s.label
}

func tableInfo(r *http.Request) *fishData {
	fd := fishData{}

	fishCSV := `ID,Name,Habitat,Average Weight KG,Price USD,Stock
	1,Tuna,Marine,250.0,10.99,50
	2,Anchovies,Marine,0.02,2.99,300
	3,Sardines,Marine,0.15,3.49,220
	4,Clownfish,Marine,0.25,15.00,25
	5,Salmon,Freshwater/Marine,4.5,12.99,60
	6,Halibut,Marine,30.0,14.50,18
	7,Cod,Marine,12.0,11.75,35
	8,Trout,Freshwater,2.5,9.99,40
	9,Mackerel,Marine,1.0,6.99,80
	10,Herring,Marine,0.5,4.25,150`

	reader := bytes.NewReader([]byte(fishCSV))
	csvReader := csv.NewReader(reader)

	tableEl, err := bridge.NewTable(csvReader)
	if err != nil {
		fmt.Print(err)
		return &fd
	}

	reader2 := bytes.NewReader([]byte(fishCSV))
	csvReader2 := csv.NewReader(reader2)
	statefulTable := buildStatefulTable(csvReader2, r)

	tbls := tables{
		Basic:    *tableEl,
		Stateful: *statefulTable,
	}

	fd = fishData{Table: &tbls}
	return &fd
}

// buildStatefulTable is so cool. 3 main parts
//  1. Build element (default values)l
//  2. Populate element based on request form. name:name attributes align, ignoring mismatches
//  3. Modify element based on its values, strict adherence.
//     e.g. ignoring select not defined by me
func buildStatefulTable(csvReader2 *csv.Reader, r *http.Request) *bridge.HTMLElement {
	const (
		// the sardine for htmlx to post to
		templatePath = "table/_stateful_table"

		// for htmlx target to replace outer html
		formRootID = "fancy-table"

		// prefixes are used with the column index to make form keys
		formKeyPrefixFilterBy = "f"
		formKeyPrefixSortBy   = "s"

		// form keys for page and limit
		formKeyPaginationLimit = "limit"
		formKeyPaginationPage  = "page"
	)

	// Row dropdown. uint because _unsigned_ is what I want for any limits or pages.
	// Plus doing this ensures I don't have to do less than zero checks since
	// parsing for them will just err which is what I want
	var (
		rowLimitSm      = labeledValue{label: "10", value: uint64(10)}
		rowLimitMd      = labeledValue{label: "50", value: uint64(50)}
		rowLimitLg      = labeledValue{label: "250", value: uint64(250)}
		rowLimitXl      = labeledValue{label: "1000", value: uint64(1000)}
		rowLimitOptions = []labeledValue{rowLimitSm, rowLimitMd, rowLimitLg, rowLimitXl}
	)

	// CSV column identifiers
	const (
		ColID = iota
		ColName
		ColHabitat
		ColAverage
		ColPrice
	)

	// Sort directions
	const (
		SortNone = iota
		SortAsc
		SortDesc
	)

	var (
		// Columns I want to add sort/filter
		colsToSortAndFilter = []int{ColName, ColHabitat, ColPrice}

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

	// 1: Defining what the element is

	table, err := bridge.NewTable(csvReader2)

	form := bridge.Stateful(table, nil)
	if err != nil {
		fmt.Print(err)
		return nil
	}
	form.EnsureAttributes()
	form.Attributes["id"] = formRootID

	headers := form.FindAll(bridge.LikeTag("th"))

	for i := range headers {
		if !slices.Contains(colsToSortAndFilter, i) {
			continue
		}
		// So we keep the sort of items even if not clicked
		sortKey := fmt.Sprintf("%s%d", formKeyPrefixSortBy, i)
		preservedSport := bridge.NewInputHidden(sortKey, "0")

		// Gives a text input to filter by. We can define these now
		// since they are not modified later and not using the 'hidden' flow.
		// Also, opted to not use my new text element fn since that adds a label
		// that I don't want.
		filterKey := fmt.Sprintf("%s%d", formKeyPrefixFilterBy, i)
		filterByDiv := bridge.NewInputText(filterKey, bridge.InputKindText, 0, 20)
		filterByDiv.DeleteFirst(bridge.LikeTag("label"))

		textBox := filterByDiv.FindFirst(bridge.LikeInput)
		textBoxHTMLX := map[bridge.AttributeKey]string{
			"hx-trigger":       "keyup changed delay:500ms",
			"hx-post":          templatePath,
			"hx-target":        fmt.Sprintf("#%s", formRootID),
			"hx-swap":          "outerHTML",
			bridge.Placeholder: "filter...",
		}
		textBox.GiveAttributes(textBoxHTMLX)

		headers[i].Children = append(headers[i].Children, filterByDiv, preservedSport)
	}

	// Add pagination state to fill in from request
	pageHiddenEl := bridge.NewInputHidden(formKeyPaginationPage, "")
	form.Children = append(form.Children, pageHiddenEl)

	// Add rout count to fill in from request
	rowCountDiv := bridge.NewInputSelect(formKeyPaginationLimit, rowLimitOptions)
	rowCountDiv.DeleteFirst(bridge.LikeTag("label"))

	rowCountSel := rowCountDiv.FindFirst(bridge.LikeInput)
	rowCountHTMLX := map[bridge.AttributeKey]string{
		"hx-trigger":       "change",
		"hx-post":          templatePath,
		"hx-target":        fmt.Sprintf("#%s", formRootID),
		"hx-swap":          "outerHTML",
		bridge.Placeholder: "filter...",
	}
	rowCountSel.GiveAttributes(rowCountHTMLX)

	form.Children = append(form.Children, rowCountDiv)

	// 2: Populating the element form request
	//
	// Now we have the final table with all input elements
	// and can get its form to see what all we need to parse
	// from a request to populate the form values from the request
	form.FillForm(r)

	// 3: Modify the element based on its values
	//
	// e.g. since a sort icon/value is dependent on current sort we must
	// modify them after having determined the current values. Same
	// thing with the pagination being based on page and limit

	for i := range headers {
		colsToSortAndFilter := []int{ColName, ColHabitat, ColPrice}
		if !slices.Contains(colsToSortAndFilter, i) {
			continue
		}

		// Gives a button to change sort
		sortKey := fmt.Sprintf("%s%d", formKeyPrefixSortBy, i)
		hiddenSort := headers[i].FindFirst(
			bridge.LikeAttribute("type", string(bridge.InputKindHidden)),
			bridge.LikeAttribute("name", sortKey),
		)
		if hiddenSort == nil {
			continue
		}

		previousSort := SortNone
		desiredSort, err := hiddenSort.ValueInt()
		if err == nil {
			if slices.Contains(acceptableSort, desiredSort) {
				previousSort = desiredSort
			}
		}

		sortBtn := bridge.HTMLElement{
			Tag: "button",
			Attributes: bridge.Attributes{
				"class":     "material-icons",
				"type":      "submit",
				"hx-post":   templatePath,
				"name":      sortKey,
				"value":     strconv.Itoa(nextSort[previousSort]),
				"hx-target": fmt.Sprintf("#%s", formRootID),
				"hx-swap":   "outerHTML",
			},
			InnerText: icons[previousSort],
		}

		headers[i].Children = append(headers[i].Children, sortBtn)
	}

	var page uint64 = 1
	parsedCurrent, err := pageHiddenEl.ValueUint64()
	if err == nil {
		page = parsedCurrent
	}

	// v, err := bridge.ValueElementSelected(rowCountSel, rowLimitOptions)
	// fmt.Print(v)
	var limit uint64 = 10
	parsedLimit, err := bridge.ValueElementSelected(&rowCountDiv, rowLimitOptions)
	if err == nil && len(parsedLimit) == 1 {
		limit = parsedLimit[0].value
	}

	// build pagination navigation buttons
	var fakeTotalCount uint64 = 1000
	pagination := bridge.NewPagination(limit, page, fakeTotalCount)
	paginationEl := pagination.Element(
		fmt.Sprintf("#%s", formRootID),
		formKeyPaginationPage,
		templatePath,
	)
	form.Children = append(form.Children, paginationEl)

	return form
}
