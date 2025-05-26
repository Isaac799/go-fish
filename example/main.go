// Package main is an example usage of the go-fish tool
package main

import (
	"fmt"
	"net/http"

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

	stock := map[string]gofish.Fish{
		"home": {
			Bait:     incrementQueryCount,
			Licenses: []gofish.License{requireSeason, springOnly},
		},
		"about page": {
			Bait: incrementQueryCount,
		},
		"user.id": {
			Bait: findUser,
		},
		"user.id.edit": {
			Bait: findUser,
		},
	}

	pond.Stock(stock)

	verbose := true
	mux := pond.CastLines(verbose)

	fmt.Println("gone fishing")
	http.ListenAndServe(":8080", mux)
}
