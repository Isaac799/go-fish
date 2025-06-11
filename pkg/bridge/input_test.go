package bridge

import (
	"testing"
	"time"
)

func TestSetValueAttribute(t *testing.T) {
	elText := NewInputText("name", InputKindText, 0, 30)
	elText.SetFirstValue("Jane Doe")
	input := elText.FindFirst(LikeInput)
	s, err := input.ParseString()
	assertNoError(t, err)
	assert(t, "Jane Doe", s)
}

func TestSetValueTextArea(t *testing.T) {
	elTextarea := NewInputTextarea("bio", 0, 30, 30, 5)
	elTextarea.SetFirstValue("I am a writer.")
	input := elTextarea.FindFirst(LikeInput)

	s, err := input.ParseString()
	assertNoError(t, err)
	assert(t, "I am a writer.", s)
}

func TestSetCheckedRadio(t *testing.T) {
	elRadio := NewInputRadio("radio color", mockColors)
	elRadio.SetFirstValue("true")

	indexes, err := elRadio.ParseIndexes()
	assertNoError(t, err)
	assertIndexes(t, indexes, []int{0})
}

func TestSetCheckedCheckbox(t *testing.T) {
	elCheckbox := NewInputCheckbox("cb color", mockColors)
	elCheckbox.SetNthValue(1, "true")
	elCheckbox.SetNthValue(3, "true")
	indexes, err := elCheckbox.ParseIndexes()
	assertNoError(t, err)
	assertIndexes(t, indexes, []int{0, 2})
}

func TestSetValueSelectable(t *testing.T) {
	elSelect := NewInputSelect("sel color", mockColors)
	elSelect.SetSelectOption(1, true)
	indexes, err := elSelect.ParseIndexes()

	assertNoError(t, err)
	assertIndexes(t, indexes, []int{1})
}

// TODO consider overlap of these tests

func TestParse(t *testing.T) {
	testParseString(t)
	testParseBool(t)
	testParseFloat(t)
	testParseInt(t)
	testParseTime(t)
	testParseIndexes(t)
}

func testParseString(t *testing.T) {
	el := mockTextEl()
	el.Attributes["value"] = "Jane Doe"
	v, err := el.ParseString()
	assertNoError(t, err)
	assert(t, v, "Jane Doe")
}
func testParseBool(t *testing.T) {
	el := mockNumberEl()
	el.Attributes["value"] = "false"
	v, err := el.ParseBool()
	assertNoError(t, err)
	assert(t, v, false)

}
func testParseFloat(t *testing.T) {
	el := mockNumberEl()
	el.Attributes["value"] = "30.76"
	v, err := el.ParseFloat()
	assertNoError(t, err)
	assert(t, v, 30.76)

}
func testParseInt(t *testing.T) {
	el := mockNumberEl()
	el.Attributes["value"] = "30"
	v, err := el.ParseInt()
	assertNoError(t, err)
	assert(t, v, 30)
}
func testParseTime(t *testing.T) {
	expectedTime, _ := time.Parse(TimeFormatHTMLDateTime, "1998-01-01T10:15")
	el := NewInputDateTime("vacation start", nil, nil)
	el.SetFirstValue("1998-01-01T10:15")
	v, err := el.ParseTime()
	assertNoError(t, err)
	assert(t, v.Unix(), expectedTime.Unix())
}

func testParseIndexes(t *testing.T) {
	elCheckbox := NewInputCheckbox("cb color", mockColors)
	elCheckbox.SetNthValue(1, "true")
	elCheckbox.SetNthValue(3, "true")
	indexes, err := elCheckbox.ParseIndexes()
	assertNoError(t, err)
	assertIndexes(t, indexes, []int{0, 2})
}
