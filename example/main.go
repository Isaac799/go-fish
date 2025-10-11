// Package main is an example usage of the go-fish tool
package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/Isaac799/go-fish/pkg/aquatic"
)

type globalData struct{}

func setupPond() (*aquatic.Pond, error) {
	config := aquatic.NewPondOptions{BeforeCatchFns: []aquatic.BeforeCatchFn{visitorLog}}
	uxPond, err := aquatic.NewPond("ux", config)
	if err != nil {
		return nil, err
	}

	assetPond, err := aquatic.NewPond("asset", aquatic.NewPondOptions{GlobalSmallFish: true})
	if err != nil {
		return nil, err
	}

	assetPond.FlowsInto(&uxPond)
	return &uxPond, nil
}

func main() {
	pond, err := setupPond()
	if err != nil {
		log.Fatal(err.Error())
	}

	rx := regexp.MustCompile

	stockFish := aquatic.Stock{
		rx("/season"): {
			BeforeCatch: []aquatic.BeforeCatchFn{optionQuery},
			OnCatch:     queriedSeason,
		},
		rx("/user/user.id"): {
			BeforeCatch: []aquatic.BeforeCatchFn{requireUser},
			OnCatch:     userInfo,
		},
	}
	pond.Stock(stockFish)

	pond.OnCatch = func(_ *http.Request) any {
		fmt.Println("caught a fish")
		return nil
	}

	go func() {
		for err := range pond.OnErr {
			fmt.Println(err.Error())
		}
	}()

	verbose := true
	mux := pond.CastLines(verbose)

	fmt.Println("gone fishing")
	http.ListenAndServe("localhost:8080", mux)
}
