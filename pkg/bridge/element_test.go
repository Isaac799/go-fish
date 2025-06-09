package bridge

import (
	"testing"
)

func TestGiveAttributes(t *testing.T) {
	style := "color: red"
	placeholder := "enter something..."
	attrs := map[string]string{
		"style":       style,
		"placeholder": placeholder,
	}
	el := NewHTMLElement("input")
	el.GiveAttributes(attrs)
	assert(t, len(el.Attributes), 2)
	assert(t, el.Attributes["style"], style)
	assert(t, el.Attributes["placeholder"], placeholder)
}

func TestEnsureAttributes(t *testing.T) {
	el := HTMLElement{}
	if el.Attributes != nil {
		t.Fatal(errNotEqual)
	}
	el.EnsureAttributes()
	if el.Attributes == nil {
		t.Fatal(errNotEqual)
	}
}

func TestSetValueAttribute(t *testing.T) {
	elText := NewInputText("name", InputKindText, 0, 30)
	elText.SetFirstValue("Jane Doe")
	input := elText.FindFirst(LikeInput)
	s, err := ElementValue[string](input)
	assertNoError(t, err)
	assert(t, "Jane Doe", s)
}

func TestSetValueTextArea(t *testing.T) {
	elTextarea := NewInputTextarea("bio", 0, 30, 30, 5)
	elTextarea.SetFirstValue("I am a writer.")
	input := elTextarea.FindFirst(LikeInput)

	s, err := ElementValue[string](input)
	assertNoError(t, err)
	assert(t, "I am a writer.", s)
}

func TestSetCheckedRadio(t *testing.T) {
	elRadio := NewInputRadio("radio color", mockColors)
	elRadio.SetFirstValue("true")

	indexes, err := ElementValue[[]int](&elRadio)
	assertNoError(t, err)
	assertIndexes(t, indexes, []int{0})
}

func TestSetCheckedCheckbox(t *testing.T) {
	elCheckbox := NewInputCheckbox("cb color", mockColors)
	elCheckbox.SetNthValue(1, "true")
	elCheckbox.SetNthValue(3, "true")
	indexes, err := ElementValue[[]int](&elCheckbox)
	assertNoError(t, err)
	assertIndexes(t, indexes, []int{0, 2})
}

func TestSetValueSelectable(t *testing.T) {
	elSelect := NewInputSelect("sel color", mockColors)
	elSelect.SetSelectOption(1, true)
	indexes, err := ElementValue[[]int](&elSelect)

	assertNoError(t, err)
	assertIndexes(t, indexes, []int{1})
}

func TestValueElementSelected(t *testing.T) {
	elSelect := NewInputSelect("sel color", mockColors)
	elSelect.SetSelectOption(1, true)

	chosen, err := ElementSelectedValue(&elSelect, mockColors)
	assertNoError(t, err)
	assert(t, len(chosen), 1)
	assert(t, chosen[0].Print(), "green")
}
