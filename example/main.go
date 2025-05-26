// Package main is an example usage of the go-fish tool
package main

import (
	"fmt"
	"net/http"
	"regexp"

	gofish "github.com/Isaac799/go-fish/internal"
)

func setupPond() gofish.Pond {
	options := gofish.NewPondOptions{
		Licenses: []gofish.License{
			visitorLog,
		},
	}

	pond, err := gofish.NewPond(
		"template",
		options,
	)

	if err != nil {
		panic(err)
	}

	return pond
}

func main() {
	pond := setupPond()

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

	pond.Stock(stockFish)

	verbose := false

	mux := pond.CastLines(verbose)

	fmt.Println("gone fishing")
	http.ListenAndServe(":8080", mux)
}
