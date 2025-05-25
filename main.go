package main

import (
	"fmt"
	"net/http"
	"time"

	gofish "github.com/Isaac799/go-fish/internal"
)

func springOnly(w http.ResponseWriter, r *http.Request) bool {
	if !r.URL.Query().Has("season") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("a season is required"))
		return false
	}

	season := r.URL.Query().Get("season")
	if season != "spring" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("cannot fish outside of spring"))
		return false
	}

	return true
}

func visitorLog(_ http.ResponseWriter, _ *http.Request) bool {
	fmt.Printf("visited %s\n", time.Now().Format(time.RFC1123))
	return true
}

func main() {
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

	for _, fish := range pond.FishFinder() {
		if fish.Pattern() == "/home" {
			fish.AddLicense(springOnly)
		}
	}

	verbose := true
	mux := pond.CastLines(verbose)

	fmt.Println("gone fishing")
	http.ListenAndServe(":8080", mux)
}
