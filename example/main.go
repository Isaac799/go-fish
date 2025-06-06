// Package main is an example usage of the go-fish tool
package main

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/Isaac799/go-fish/pkg/aquatic"
)

type globalData struct{}

func setupPond[T, K any]() aquatic.Pond[T, K] {
	uxPond, err := aquatic.NewPond[T, K](
		"ux",
		aquatic.NewPondOptions{
			Licenses: []aquatic.License{
				visitorLog,
			},
		},
	)
	if err != nil {
		panic(err)
	}

	assetPond, err := aquatic.NewPond[T, K](
		"asset",
		aquatic.NewPondOptions{
			GlobalSmallFish: true,
		},
	)

	if err != nil {
		panic(err)
	}

	aquatic.FlowsInto(&assetPond, &uxPond)

	return uxPond
}

func main() {
	pond := setupPond[globalData, *fishData]()

	stockFish := map[*regexp.Regexp]aquatic.Fish[*fishData]{
		regexp.MustCompile("season"): {
			Licenses: []aquatic.License{optionQuery},
			Bait:     queriedSeason,
		},
		regexp.MustCompile("user/.id"): {
			Bait:     userInfo,
			Licenses: []aquatic.License{requireUser},
		},
		regexp.MustCompile("/form"): {
			Bait: exampleFormBait,
		},
		regexp.MustCompile("/table"): {
			Bait: tableInfo,
		},
	}

	aquatic.StockPond(&pond, stockFish)

	gd := globalData{}

	pond.GlobalBait = func(_ *http.Request) globalData {
		a := gd
		return a
	}

	verbose := true

	mux := aquatic.CastLines(&pond, verbose)

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
