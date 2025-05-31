package bridge

import (
	"fmt"
	"strings"

	gofish "github.com/Isaac799/go-fish/internal"
)

// NewAnchor enables lazy but safe development of anchor tags
// for a page by looking at the tuna fish in a pond and their patterns
// taking into account path values. A pattern like "/user/3/edit" will match
// "/user/{id}/edit" by omitting the path value from comparison.
// Panics to prevent serving invalid anchors.
func NewAnchor[T, K any](innerText, patternLike string, pond *gofish.Pond[T, K]) HTMLElement {
	fish := gofish.FishFinder(pond)

	patternLikeParts := strings.Split(patternLike, "/")

	var anchorFish *gofish.Fish[K]
	for _, fish := range fish {
		if gofish.Kind[T](fish) != gofish.FishKindTuna {
			continue
		}

		if patternLike == gofish.Patten(fish) {
			anchorFish = fish
			break
		}

		p2 := gofish.Patten(fish)
		fishLikeParts := strings.Split(p2, "/")

		if len(fishLikeParts) != len(patternLikeParts) {
			continue
		}

		match := true
		for i := range fishLikeParts {
			desire := patternLikeParts[i]
			actual := fishLikeParts[i]
			isPathValue := strings.HasPrefix(actual, "{") && strings.HasSuffix(actual, "}")
			if isPathValue {
				continue
			}
			if desire == actual {
				continue
			}
			match = false
		}
		if !match {
			continue
		}

		anchorFish = fish
		break
	}

	if anchorFish == nil {
		s := fmt.Sprintf("cannot find fish with an anchor like %s", patternLike)
		panic(s)
	}

	pattern := gofish.Patten(anchorFish)
	anchor := HTMLElement{
		Tag:       "a",
		InnerText: innerText,
		Attributes: map[AttributeKey]string{
			HRef: pattern,
		},
	}

	return anchor
}
