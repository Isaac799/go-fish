// Package main is an example usage of the go-fish tool
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"text/template"

	"github.com/Isaac799/go-fish/example/bridge"
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
		"bridge",
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
			Bait: func(_ *http.Request) any {
				form := exampleForm()
				return form
			},
		},
	}

	uxPond.Stock(stockFish)

	return uxPond
}

func main() {
	pond := setupPond()

	verbose := true

	mux := pond.CastLines(verbose)

	mux.HandleFunc("/submit/test", func(w http.ResponseWriter, r *http.Request) {
		form := exampleForm()
		s := bridge.FormFromRequest(r, form)
		b, err := json.Marshal(s)
		if err != nil {
			fmt.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	})

	fmt.Println("gone fishing")
	http.ListenAndServe(":8080", mux)
}
