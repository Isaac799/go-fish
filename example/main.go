// Package main is an example usage of the go-fish tool
package main

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/Isaac799/go-fish/pkg/aquatic"
)

type globalData struct{}

func setupPond() aquatic.Pond {
	config := aquatic.NewPondOptions{Licenses: []aquatic.License{visitorLog}}
	uxPond, err := aquatic.NewPond("ux", config)
	if err != nil {
		panic(err)
	}

	assetPond, err := aquatic.NewPond("asset", aquatic.NewPondOptions{GlobalSmallFish: true})
	if err != nil {
		panic(err)
	}

	assetPond.FlowsInto(&uxPond)
	return uxPond
}

func main() {
	pond := setupPond()
	rx := regexp.MustCompile

	stockFish := aquatic.Stock{
		rx("season"): {
			Licenses: []aquatic.License{optionQuery},
			OnCatch:  queriedSeason,
		},
		rx("user/.id"): {
			OnCatch:  userInfo,
			Licenses: []aquatic.License{requireUser},
		},
	}
	pond.Stock(stockFish)

	pond.OnCatch = func(_ *http.Request) any {
		fmt.Println("caught a fish")
		return nil
	}

	verbose := true
	mux := pond.CastLines(verbose)

	fmt.Println("gone fishing")
	http.ListenAndServe("localhost:8080", mux)
}
