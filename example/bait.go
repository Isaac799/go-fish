package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/Isaac799/go-fish/pkg/bridge"
	"github.com/Isaac799/go-fish/pkg/bridge/table"
)

type dragDropItem struct {
	ID,
	X,
	Y int
}

type user struct {
	ID        int
	FirstName string
	LastName  string
}

type fishData struct {
	Season   string
	User     *user
	Table    *bridge.HTMLElement
	Form     *bridge.HTMLElement
	DragDrop *bridge.HTMLElement
}

var (
	statefulTableID = bridge.RandomID()
)

func dragDrop(r *http.Request) *fishData {
	// default positions
	items := []dragDropItem{
		{ID: 1, X: 50, Y: 50},
		{ID: 2, X: 200, Y: 150},
	}

	container := bridge.ElementWithState(&bridge.HTMLElement{
		Tag:      "div",
		Children: make([]bridge.HTMLElement, 0, len(items)+1),
	}, nil)
	container.EnsureAttributes()
	container.Attributes["id"] = "container"

	updateEl := bridge.HTMLElement{
		Tag: "span",
		Attributes: bridge.Attributes{
			"hx-post": "/drag-drop/_container",
			// positionUpdated emitted by JS on drag end
			"hx-trigger": "positionUpdated from:body",
			"hx-target":  "#container",
			"hx-swap":    "outerHTML",
		},
	}
	container.Children = append(container.Children, updateEl)

	var nodeIdentifiers = func(i int) (string, string, string) {
		nodeID := fmt.Sprintf("n%d", i)
		nameX := fmt.Sprintf("%sx", nodeID)
		nameY := fmt.Sprintf("%sy", nodeID)
		return nodeID, nameX, nameY
	}

	xInputs := make([]*bridge.HTMLInput, 0, len(items))
	yInputs := make([]*bridge.HTMLInput, 0, len(items))

	for i, item := range items {
		nodeID, nameX, nameY := nodeIdentifiers(i)
		draggableEl := bridge.HTMLElement{
			Tag: "div",
			Attributes: bridge.Attributes{
				"class": "node",
				"id":    nodeID,
			},
			Children: make([]bridge.HTMLElement, 0, 2),
		}

		draggableX := bridge.NewInputHidden(nameX, strconv.Itoa(item.X))
		xInputs = append(xInputs, &draggableX)
		draggableY := bridge.NewInputHidden(nameY, strconv.Itoa(item.Y))
		yInputs = append(yInputs, &draggableY)

		draggableEl.Children = append(draggableEl.Children, draggableX, draggableY)
		container.Children = append(container.Children, draggableEl)
	}

	container.FormFill(r)

	nodes := container.FindAll(bridge.LikeAttribute("class", "node"))

	// example of enforcing data constraint server side
	var roundToNearest = func(val, nearest float64) int {
		return int(math.Round(val/nearest) * nearest)
	}

	for i, el := range nodes {
		el.EnsureAttributes()

		x, err := xInputs[i].ParseFloat()
		if err != nil {
			print(err.Error())
			continue
		}

		y, err := yInputs[i].ParseFloat()
		if err != nil {
			print(err.Error())
			continue
		}

		roundBy := 25.0
		roundedX := roundToNearest(x, roundBy)
		roundedY := roundToNearest(y, roundBy)

		el.Attributes["style"] = fmt.Sprintf("transform: translate(%dpx, %dpx)", roundedX, roundedY)
	}

	data := fishData{DragDrop: container}
	return &data
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

func tableInfo(r *http.Request) *fishData {
	fd := fishData{}

	fd = fishData{Table: buildStatefulTable(r)}
	return &fd
}

// buildStatefulTable is so cool. 3 main parts
//  1. Define the element
//  2. Populate it (based on request form -  name:name attributes align, ignoring mismatches)
//  3. Modify it based on its values
func buildStatefulTable(r *http.Request) *bridge.HTMLElement {
	// CSV column identifiers
	const (
		ColID = iota
		ColName
		ColHabitat
		ColAverage
		ColPrice
		ColStock
	)

	// 1: Defining what the element is

	headers := []string{"ID", "Name", "Habitat", "Average Weight KG", "Price USD", "Stock"}
	hxPost := "table/_stateful_table"
	config := table.NewConfig(statefulTableID, headers, hxPost)
	tbl, _ := table.New(config)

	definePagination, modifyPagination := table.Pagination()
	defineSort, modifySort := table.Sort(ColAverage, ColPrice)
	defineFilter := table.Filter(ColName, ColHabitat)

	err := tbl.Modify(defineSort, definePagination, defineFilter)
	if err != nil {
		fmt.Print(err)
		return nil
	}

	// 2: Populating the element form request
	//
	// Now we have the final table with all input elements
	// and can get its form to see what all we need to parse
	// from a request to populate the form values from the request
	tbl.El.FormFill(r)

	constraints, err := tbl.Constraints()
	if err != nil {
		fmt.Print(err)
		return nil
	}

	// allows compile
	_ = constraints

	// simulate results from constraints
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

	// setting arbitrary limit for testing, in reality
	// derive from a query
	tbl.RecordCount = 1000

	// 3: Modify the element based on its values
	//
	// e.g. since a sort icon/value is dependent on current sort we must
	// modify them after having determined the current values. Same
	// thing with the pagination being based on page and limit

	err = tbl.Modify(modifySort, modifyPagination)
	if err != nil {
		fmt.Print(err)
		return nil
	}

	err = tbl.SetData(records)
	if err != nil {
		fmt.Print(err)
		return nil
	}

	return tbl.El
}
