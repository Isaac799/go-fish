package main

import (
	"fmt"
	"os"

	gofish "github.com/Isaac799/go-fish/internal"
)

type page struct {
	pattern    string
	data       gofish.Bait
	middleware []gofish.License
}

func (p *page) useFish(f *gofish.Fish) {
	f.Bait = p.data
	f.Licenses = p.middleware
}

func (p *page) isFor(f *gofish.Fish) bool {
	return p.pattern == f.Pattern()
}

func setupPages(pond *gofish.Pond, pages []page) {
	for _, page := range pages {
		found := false
		for _, fish := range pond.FishFinder() {
			if !page.isFor(fish) {
				continue
			}
			found = true
			page.useFish(fish)
		}
		if !found {
			fmt.Println("did not find matching fish for page")
			os.Exit(1)
		}
	}
}
