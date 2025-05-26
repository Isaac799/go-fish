// Package main is an example usage of the go-fish tool
package main

import (
	"fmt"
	"net/http"
	"regexp"

	gofish "github.com/Isaac799/go-fish/internal"
)

func setupPonds() (gofish.Pond, gofish.Pond) {
	appPond, err := gofish.NewPond(
		"template",
		gofish.NewPondOptions{
			Licenses: []gofish.License{
				visitorLog,
			},
		},
	)

	assetPond, err := gofish.NewPond(
		"asset",
		gofish.NewPondOptions{
			GlobalAnchovyAndClown: true,
		},
	)

	if err != nil {
		panic(err)
	}

	return appPond, assetPond
}

func main() {
	appPond, assetPond := setupPonds()

	stockFish := map[*regexp.Regexp]gofish.Fish{
		regexp.MustCompile("blog"): {
			Bait: incrementQueryCount,
		},
		regexp.MustCompile("home"): {
			Bait:     incrementQueryCount,
			Licenses: []gofish.License{requireSeason, springOnly},
		},
		regexp.MustCompile("about page"): {
			Bait: incrementQueryCount,
		},
		regexp.MustCompile("user"): {
			Bait: findUser,
		},
	}

	appPond.Stock(stockFish)
	assetPond.FlowsInto(&appPond)

	verbose := true

	mux := appPond.CastLines(verbose)

	fmt.Println("gone fishing")
	http.ListenAndServe(":8080", mux)
}
