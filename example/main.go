// Package main is an example usage of the go-fish tool
package main

import (
	"fmt"
	"net/http"
	"regexp"
	"text/template"

	"github.com/Isaac799/go-fish/example/bridge"
	gofish "github.com/Isaac799/go-fish/internal"
)

type anchors struct {
	Home       bridge.HTMLElement
	Water      bridge.HTMLElement
	Season     bridge.HTMLElement
	User       bridge.HTMLElement
	Form       bridge.HTMLElement
	Table      bridge.HTMLElement
	UserID     bridge.HTMLElement
	UserIDEdit bridge.HTMLElement
}

func newAnchors[T, K any](p *gofish.Pond[T, K]) anchors {
	return anchors{
		Home:       bridge.NewAnchor("Home", "/", p),
		Water:      bridge.NewAnchor("Water", "/water", p),
		Season:     bridge.NewAnchor("Seasons", "/season", p),
		Table:      bridge.NewAnchor("Table", "/table", p),
		Form:       bridge.NewAnchor("Form", "/form", p),
		User:       bridge.NewAnchor("User", "/user", p),
		UserID:     bridge.NewAnchor("View User", "/user/3", p),
		UserIDEdit: bridge.NewAnchor("Edit User", "/user/3/edit", p),
	}
}

type GlobalData struct {
	Anchors anchors
}

func setupPond[T, K any]() gofish.Pond[T, K] {
	uxPond, err := gofish.NewPond[T, K](
		"ux",
		gofish.NewPondOptions{
			Licenses: []gofish.License{
				visitorLog,
			},
		},
	)
	if err != nil {
		panic(err)
	}

	assetPond, err := gofish.NewPond[T, K](
		"asset",
		gofish.NewPondOptions{
			GlobalSmallFish: true,
		},
	)

	if err != nil {
		panic(err)
	}

	elementPond, err := gofish.NewPond[T, K](
		"bridge",
		gofish.NewPondOptions{
			GlobalSmallFish: true,
		},
	)
	if err != nil {
		panic(err)
	}

	gofish.FlowsInto(&assetPond, &uxPond)
	gofish.FlowsInto(&elementPond, &uxPond)

	return uxPond
}

func main() {
	pond := setupPond[GlobalData, *fishData]()

	stockFish := map[*regexp.Regexp]gofish.Fish[*fishData]{
		regexp.MustCompile("season"): {
			Licenses: []gofish.License{optionQuery},
			Bait:     queriedSeason,
		},
		regexp.MustCompile("user/.id"): {
			Bait:     userInfo,
			Licenses: []gofish.License{requireUser},
		},
		regexp.MustCompile("water"): {
			Bait: waterInfo,
			Tackle: template.FuncMap{
				"printTime": printTime,
			},
		},
		regexp.MustCompile("/form"): {
			Bait: exampleFormBait,
		},
		regexp.MustCompile("/table"): {
			Bait: tableInfo,
		},
	}

	gofish.StockPond(&pond, stockFish)

	gd := GlobalData{
		Anchors: newAnchors(&pond),
	}

	pond.GlobalBait = func(_ *http.Request) GlobalData {
		a := gd
		return a
	}

	verbose := true

	mux := gofish.CastLines(&pond, verbose)

	// mux.HandleFunc("/submit/test", func(w http.ResponseWriter, r *http.Request) {
	// 	form := exampleForm()
	// 	s := bridge.FormFromRequest(r, &form)
	// 	b, err := json.Marshal(s)
	// 	if err != nil {
	// 		fmt.Print(err)
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		return
	// 	}
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.Write(b)
	// })

	fmt.Println("gone fishing")
	http.ListenAndServe(":8080", mux)
}
