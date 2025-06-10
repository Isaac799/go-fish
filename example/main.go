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
	config := aquatic.NewPondOptions{Licenses: []aquatic.License{visitorLog}}
	uxPond, err := aquatic.NewPond[T, K]("ux", config)
	if err != nil {
		panic(err)
	}

	assetPond, err := aquatic.NewPond[T, K]("asset", aquatic.NewPondOptions{GlobalSmallFish: true})
	if err != nil {
		panic(err)
	}

	aquatic.FlowsInto(&assetPond, &uxPond)
	return uxPond
}

func main() {
	pond := setupPond[globalData, *fishData]()
	rx := regexp.MustCompile

	stockFish := aquatic.Stock[globalData, *fishData]{
		rx("season"): {
			Licenses: []aquatic.License{optionQuery},
			Bait:     queriedSeason,
		},
		rx("user/.id"): {
			Bait:     userInfo,
			Licenses: []aquatic.License{requireUser},
		},
		rx("/form"): {
			Bait: exampleFormBait,
		},
		rx("/table"): {
			Bait: tableInfo,
		},
		rx("drag-drop"): {
			Bait: dragDrop,
		},
	}
	aquatic.StockPond(&pond, stockFish)

	gd := globalData{}
	pond.Chum = func(_ *http.Request) globalData {
		return gd
	}

	verbose := true
	mux := aquatic.CastLines(&pond, verbose)

	fmt.Println("gone fishing")
	http.ListenAndServe("localhost:8080", mux)
}
