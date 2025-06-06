package bridge

import "testing"

func TestDeleteFirst(t *testing.T) {
	textBox := NewInputText("first name", InputKindText, 3, 30)
	textBox.DeleteFirst(LikeTag("label"))
	eq(t, len(textBox.Children), 1)
	eq(t, textBox.Children[0].Tag, "input")
}

func TestDeleteNth(t *testing.T) {
	colorRadio := NewInputRadio("colors", mockColors)
	colorRadio.DeleteNth(2, LikeTag("input"))
	divs := colorRadio.Children

	eq(t, len(divs), 4)
	// legend := divs[0] // for note

	red := divs[1]
	eq(t, len(red.Children), 2)
	eq(t, red.Children[0].Tag, "label")
	eq(t, red.Children[1].Tag, "input")

	green := divs[2]
	eq(t, len(green.Children), 1)
	eq(t, green.Children[0].Tag, "label")

	blue := divs[3]
	eq(t, len(blue.Children), 2)
	eq(t, blue.Children[0].Tag, "label")
	eq(t, blue.Children[1].Tag, "input")
}

func TestDeleteAll(t *testing.T) {
	colorRadio := NewInputRadio("colors", mockColors)
	colorRadio.DeleteAll(LikeTag("input"))
	divs := colorRadio.Children

	eq(t, len(divs), 4)
	// legend := divs[0] // for note

	red := divs[1]
	eq(t, len(red.Children), 1)
	eq(t, red.Children[0].Tag, "label")

	green := divs[2]
	eq(t, len(green.Children), 1)
	eq(t, green.Children[0].Tag, "label")

	blue := divs[3]
	eq(t, len(blue.Children), 1)
	eq(t, blue.Children[0].Tag, "label")
}
