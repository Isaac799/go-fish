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

	home := page{
		pattern:    "/home",
		data:       incrementQueryCount,
		middleware: []gofish.License{springOnly},
	}
	about := page{
		pattern:    "/about-page",
		middleware: []gofish.License{springOnly},
	}

	pages := []page{home, about}

	setupPages(&pond, pages)

	verbose := true
	mux := pond.CastLines(verbose)

	fmt.Println("gone fishing")
	http.ListenAndServe(":8080", mux)
}
