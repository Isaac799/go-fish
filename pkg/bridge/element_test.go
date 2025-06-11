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

func TestValueElementSelected(t *testing.T) {
	elSelect := NewInputSelect("sel color", mockColors)
	elSelect.SetSelectOption(1, true)

	chosen, err := InputSelectedValue(&elSelect, mockColors)
	assertNoError(t, err)
	assert(t, len(chosen), 1)
	assert(t, chosen[0].String(), "green")
}

func TestClass(t *testing.T) {
	el := HTMLElement{}
	el.Attributes = map[string]string{"class": "red blue"}
	classes := el.Class()
	assertIndexes(t, classes, []string{"red", "blue"})
}

func TestAppendClass(t *testing.T) {
	el := HTMLElement{}
	el.Attributes = map[string]string{"class": "red"}
	el.AppendClass("blue")
	s := el.Attributes["class"]
	assert(t, s, "red blue")
}

func TestStyle(t *testing.T) {
	el := HTMLElement{}
	el.Attributes = map[string]string{"style": "color: red"}
	style := el.Style()
	s := style["color"]
	assert(t, s, "red")
}

func TestAppendStyle(t *testing.T) {
	el := HTMLElement{}
	el.Attributes = map[string]string{"style": "color: red"}
	el.AppendStyle("background-color", "blue")
	s := el.Attributes["style"]
	assert(t, s, "color:red;background-color:blue")
}
