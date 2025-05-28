// Package main is an example usage of the go-fish tool
package main

import (
	"fmt"
	"net/http"
	"regexp"
	"text/template"

	gofish "github.com/Isaac799/go-fish/internal"
)

func setupPond() gofish.Pond {
	uxPond, err := gofish.NewPond(
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

	assetPond, err := gofish.NewPond(
		"asset",
		gofish.NewPondOptions{
			GlobalSmallFish: true,
		},
	)

	if err != nil {
		panic(err)
	}

	elementPond, err := gofish.NewPond(
		"element",
		gofish.NewPondOptions{
			GlobalSmallFish: true,
		},
	)
	if err != nil {
		panic(err)
	}

	assetPond.FlowsInto(&uxPond)
	elementPond.FlowsInto(&uxPond)

	stockFish := map[*regexp.Regexp]gofish.Fish{
		regexp.MustCompile("season"): {
			Licenses: []gofish.License{optionQuery},
			Bait:     queriedSeason,
		},
		regexp.MustCompile("user.id"): {
			Bait:     userInfo,
			Licenses: []gofish.License{requireUser},
		},
		regexp.MustCompile("water"): {
			Bait: waterInfo,
			Tackle: template.FuncMap{
				"printTime": printTime,
			},
		},
		regexp.MustCompile("[input]"): {
			Bait: inputs,
		},
	}

	uxPond.Stock(stockFish)

	return uxPond
}

func main() {
	pond := setupPond()

	verbose := true

	mux := pond.CastLines(verbose)

	fmt.Println("gone fishing")
	http.ListenAndServe(":8080", mux)
}
