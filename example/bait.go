package main

import (
	"bytes"
	"encoding/csv"
	"net/http"
	"strconv"
	"time"

	"github.com/Isaac799/go-fish/example/bridge"
)

type user struct {
	ID        int
	FirstName string
	LastName  string
}

type water struct {
	ServerTime time.Time
	RotateDeg  int
	Table      bridge.HTMLElement
}

type fishData struct {
	Season string
	User   *user
	Water  *water
	Table  *bridge.HTMLElement
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

func tableInfo(_ *http.Request) *fishData {
	fishCSV := `id,name,habitat,average_weight_kg,price_usd,stock
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
		return nil
	}

	state := map[string]string{
		"page":        "0",
		"limit":       "0",
		"sort_name":   "0",
		"filter_name": "0",
	}

	statefulTable := bridge.Stateful(*tableEl, state)

	fd := fishData{
		Table: &statefulTable,
	}

	return &fd
}

func waterInfo(r *http.Request) *fishData {
	data := fishData{}
	posStr := r.URL.Query().Get("pos")
	offsetStr := r.URL.Query().Get("off")

	w := water{
		RotateDeg:  0,
		ServerTime: time.Now(),
	}

	off, err := strconv.Atoi(offsetStr)
	if err != nil {
		data.Water = &w
		return &data
	}
	pos, err := strconv.Atoi(posStr)
	if err != nil {
		data.Water = &w
		return &data
	}

	w.RotateDeg = pos + off

	data.Water = &w
	return &data
}
