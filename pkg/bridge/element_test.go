package bridge

import (
	"testing"
)

func TestGiveAttributes(t *testing.T) {
	style := "color: red"
	placeholder := "enter something..."
	attrs := map[AttributeKey]string{
		"style":       style,
		"placeholder": placeholder,
	}
	el := NewHTMLElement("input")
	el.GiveAttributes(attrs)
	eq(t, len(el.Attributes), 2)
	eq(t, el.Attributes["style"], style)
	eq(t, el.Attributes["placeholder"], placeholder)
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

	s, err := elText.FindFirst(LikeInput).ValueString()
	noErr(t, err)
	eq(t, "Jane Doe", s)
}

func TestSetValueTextArea(t *testing.T) {
	elTextarea := NewInputTextarea("bio", 0, 30, 30, 5)
	elTextarea.SetFirstValue("I am a writer.")

	s, err := elTextarea.ValueString()
	noErr(t, err)
	eq(t, "I am a writer.", s)
}

func TestSetCheckedRadio(t *testing.T) {
	elRadio := NewInputRadio("radio color", mockColors)
	elRadio.SetFirstValue("true")
	indexes, err := elRadio.ValueIndexes()
	noErr(t, err)
	eqIndexes(t, indexes, []int{0})
}

func TestSetCheckedCheckbox(t *testing.T) {
	elCheckbox := NewInputCheckbox("cb color", mockColors)
	elCheckbox.SetNthValue(1, "true")
	elCheckbox.SetNthValue(3, "true")
	indexes, err := elCheckbox.ValueIndexes()
	noErr(t, err)
	eqIndexes(t, indexes, []int{0, 2})
}

func TestSetValueSelectable(t *testing.T) {
	elSelect := NewInputSelect("sel color", mockColors)
	elSelect.SetSelectOption(1, true)

	indexes, err := elSelect.ValueIndexes()
	noErr(t, err)
	eqIndexes(t, indexes, []int{1})
}

func TestValueElementSelected(t *testing.T) {
	elSelect := NewInputSelect("sel color", mockColors)
	elSelect.SetSelectOption(1, true)

	chosen, err := ValueElementSelected(&elSelect, mockColors)
	noErr(t, err)
	eq(t, len(chosen), 1)
	eq(t, chosen[0].Print(), "green")
}
