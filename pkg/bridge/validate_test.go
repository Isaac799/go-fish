package bridge

import (
	"strconv"
	"testing"
	"time"
)

func TestValidateString(t *testing.T) {
	testValidString(t)
	testStringNotRequired(t)
	testStringPatternMismatch(t)
	testTooShortString(t)
	testTooLongString(t)
}

func TestNumberValidate(t *testing.T) {
	testNumberValid(t)
	testNumberNotRequired(t)
	testNumberBelowMin(t)
	testNumberAboveMax(t)
	testNumberInvalid(t)
}

func TestDateValidate(t *testing.T) {
	testDateValid(t)
	testDateBeforeMin(t)
	testDateAfterMax(t)
	testDateInvalidFormat(t)
}

func TestTimeValidate(t *testing.T) {
	testTimeValid(t)
	testTimeBeforeMin(t)
	testTimeAfterMax(t)
	testInvalidTimeFormat(t)
}

func TestDateTimeValidate(t *testing.T) {
	testDateTimeValid(t)
	testDateTimeBeforeMin(t)
	testDateTimeAfterMax(t)
	testDateTimeFormatInvalid(t)
}

func mockTextEl() *HTMLElement {
	elText := NewInputText("name", InputKindText, 5, 10)
	elInput := elText.FindFirst(LikeInput)
	return elInput
}

func testValidString(t *testing.T) {
	el := mockTextEl()
	el.Attributes["value"] = "Jane Doe"
	okay, err := el.Validate()
	assertNoError(t, err)
	assert(t, okay, true)
}

func testStringNotRequired(t *testing.T) {
	el := mockTextEl()
	el.Attributes["required"] = strconv.FormatBool(false)
	el.Attributes["value"] = ""
	okay, err := el.Validate()
	assertNoError(t, err)
	assert(t, okay, true)
}

func testStringPatternMismatch(t *testing.T) {
	el := mockTextEl()
	el.Attributes["pattern"] = "John"
	el.Attributes["value"] = "Jane"
	okay, err := el.Validate()
	assertNoError(t, err)
	assert(t, okay, false)
}

func testTooShortString(t *testing.T) {
	el := mockTextEl()
	el.Attributes["value"] = "Jane"
	okay, err := el.Validate()
	assertNoError(t, err)
	assert(t, okay, false)
}

func testTooLongString(t *testing.T) {
	el := mockTextEl()
	el.Attributes["value"] = "Jane Doe Jane Doe Jane Doe"
	okay, err := el.Validate()
	assertNoError(t, err)
	assert(t, okay, false)
}

// numbers

func mockNumberEl() *HTMLElement {
	elNumber := NewInputNumber("age", 18, 99)
	return elNumber.FindFirst(LikeInput)
}

func testNumberValid(t *testing.T) {
	el := mockNumberEl()
	el.Attributes["value"] = "30"
	okay, err := el.Validate()
	assertNoError(t, err)
	assert(t, okay, true)
}

func testNumberNotRequired(t *testing.T) {
	el := mockNumberEl()
	el.Attributes["required"] = strconv.FormatBool(false)
	el.Attributes["value"] = ""
	_, err := el.Validate()
	if err == nil {
		t.Fatal()
	}
}

func testNumberBelowMin(t *testing.T) {
	el := mockNumberEl()
	el.Attributes["value"] = "17"
	okay, err := el.Validate()
	assertNoError(t, err)
	assert(t, okay, false)
}

func testNumberAboveMax(t *testing.T) {
	el := mockNumberEl()
	el.Attributes["value"] = "100"
	okay, err := el.Validate()
	assertNoError(t, err)
	assert(t, okay, false)
}

func testNumberInvalid(t *testing.T) {
	el := mockNumberEl()
	el.Attributes["value"] = "notanumber"
	_, err := el.Validate()
	if err == nil {
		t.Fatal()
	}
}

// date

func mockDateEl() *HTMLElement {
	min, _ := time.Parse(TimeFormatHTMLDate, "1970-01-01")
	max, _ := time.Parse(TimeFormatHTMLDate, "1990-01-01")
	input := NewInputDate("birthday", &min, &max)
	return input.FindFirst(LikeInput)
}

func testDateValid(t *testing.T) {
	el := mockDateEl()
	el.Attributes["value"] = "1980-06-15"
	okay, err := el.Validate()
	assertNoError(t, err)
	assert(t, okay, true)
}

func testDateBeforeMin(t *testing.T) {
	el := mockDateEl()
	el.Attributes["value"] = "1969-12-31"
	okay, err := el.Validate()
	assertNoError(t, err)
	assert(t, okay, false)
}

func testDateAfterMax(t *testing.T) {
	el := mockDateEl()
	el.Attributes["value"] = "1991-01-01"
	okay, err := el.Validate()
	assertNoError(t, err)
	assert(t, okay, false)
}

func testDateInvalidFormat(t *testing.T) {
	el := mockDateEl()
	el.Attributes["value"] = "15-06-1980"
	_, err := el.Validate()
	if err == nil {
		t.Fatal()
	}
}

// time

func mockTimeEl() *HTMLElement {
	min, _ := time.Parse(TimeFormatHTMLTime, "10:00")
	max, _ := time.Parse(TimeFormatHTMLTime, "10:30")
	input := NewInputTime("clock in", &min, &max)
	return input.FindFirst(LikeInput)
}

func testTimeValid(t *testing.T) {
	el := mockTimeEl()
	el.Attributes["value"] = "10:15"
	okay, err := el.Validate()
	assertNoError(t, err)
	assert(t, okay, true)
}

func testTimeBeforeMin(t *testing.T) {
	el := mockTimeEl()
	el.Attributes["value"] = "09:59"
	okay, err := el.Validate()
	assertNoError(t, err)
	assert(t, okay, false)
}

func testTimeAfterMax(t *testing.T) {
	el := mockTimeEl()
	el.Attributes["value"] = "10:31"
	okay, err := el.Validate()
	assertNoError(t, err)
	assert(t, okay, false)
}

func testInvalidTimeFormat(t *testing.T) {
	el := mockTimeEl()
	el.Attributes["value"] = "25:61"
	_, err := el.Validate()
	if err == nil {
		t.Fatal()
	}
}

// datetime

func mockDateTimeEl() *HTMLElement {
	min, _ := time.Parse(TimeFormatHTMLDateTime, "1998-01-01T10:15")
	max, _ := time.Parse(TimeFormatHTMLDateTime, "2000-01-01T10:15")
	input := NewInputDateTime("vacation start", &min, &max)
	return input.FindFirst(LikeInput)
}

func testDateTimeValid(t *testing.T) {
	el := mockDateTimeEl()
	el.Attributes["value"] = "1999-06-01T12:00"
	okay, err := el.Validate()
	assertNoError(t, err)
	assert(t, okay, true)
}

func testDateTimeBeforeMin(t *testing.T) {
	el := mockDateTimeEl()
	el.Attributes["value"] = "1997-12-31T10:14"
	okay, err := el.Validate()
	assertNoError(t, err)
	assert(t, okay, false)
}

func testDateTimeAfterMax(t *testing.T) {
	el := mockDateTimeEl()
	el.Attributes["value"] = "2000-01-01T10:16"
	okay, err := el.Validate()
	assertNoError(t, err)
	assert(t, okay, false)
}

func testDateTimeFormatInvalid(t *testing.T) {
	el := mockDateTimeEl()
	el.Attributes["value"] = "01/01/1999 12:00"
	_, err := el.Validate()
	if err == nil {
		t.Fatal()
	}
}
