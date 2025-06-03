// Package test holds a few tests based on mock data
// to help me reproducible development issues.
package bridge

import (
	"testing"
)

func TestFindNth(t *testing.T) {
	form := mockEmptyForm()

	thirdInput := form.FindNth(3, LikeInput)
	eq(t, thirdInput.Attributes["name"], "favorite number")
}

func TestFindAll(t *testing.T) {
	form := mockEmptyForm()

	allInputs := form.FindAll(LikeInput)
	eq(t, len(allInputs), 16)
}

func TestFindFirst(t *testing.T) {
	form := mockEmptyForm()

	first := form.FindFirst(LikeInput)
	eq(t, first.Attributes["name"], "name")
}
