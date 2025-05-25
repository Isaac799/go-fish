package main

import (
	gofish "github.com/Isaac799/go-fish/internal"
)

func prepHomeFish(fish *gofish.Fish) {
	fish.AddLicense(springOnly)
	fish.Bait = counterBait
}
