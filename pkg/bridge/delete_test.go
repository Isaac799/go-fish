package bridge

import "testing"

func TestDeleteFirst(t *testing.T) {
	textBox := NewInputText("first name", InputKindText, 3, 30)
	textBox.DeleteFirst(LikeTag("label"))
	assert(t, len(textBox.Children), 1)
	assert(t, textBox.Children[0].Tag, "input")
}

func TestDeleteNth(t *testing.T) {
	colorRadio := NewInputRadio("colors", mockColors)
	colorRadio.DeleteNth(2, LikeTag("input"))
	divs := colorRadio.Children

	assert(t, len(divs), 4)
	// legend := divs[0] // for note

	red := divs[1]
	assert(t, len(red.Children), 2)
	assert(t, red.Children[0].Tag, "label")
	assert(t, red.Children[1].Tag, "input")

	green := divs[2]
	assert(t, len(green.Children), 1)
	assert(t, green.Children[0].Tag, "label")

	blue := divs[3]
	assert(t, len(blue.Children), 2)
	assert(t, blue.Children[0].Tag, "label")
	assert(t, blue.Children[1].Tag, "input")
}

func TestDeleteAll(t *testing.T) {
	colorRadio := NewInputRadio("colors", mockColors)
	colorRadio.DeleteAll(LikeTag("input"))
	divs := colorRadio.Children

	assert(t, len(divs), 4)
	// legend := divs[0] // for note

	red := divs[1]
	assert(t, len(red.Children), 1)
	assert(t, red.Children[0].Tag, "label")

	green := divs[2]
	assert(t, len(green.Children), 1)
	assert(t, green.Children[0].Tag, "label")

	blue := divs[3]
	assert(t, len(blue.Children), 1)
	assert(t, blue.Children[0].Tag, "label")
}
