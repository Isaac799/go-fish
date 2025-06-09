package bridge

import "testing"

func TestStateful(t *testing.T) {
	state := map[string]string{
		"page":        "1",
		"limit":       "10",
		"sort_name":   "asc",
		"filter_name": "marine",
	}

	divEl := HTMLElement{
		Tag: "div",
	}
	form := ElementWithState(&divEl, state)
	assert(t, len(form.Children), 5)
	assert(t, form.Children[0].Tag, "div")
}
