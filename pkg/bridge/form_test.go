// Package test holds a few tests based on mock data
// to help me reproducible development issues.
package bridge

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"
)

var (
	errFoundErr         = errors.New("an err was found")
	errNotEqual         = errors.New("a and b where not equal")
	errUnexpectedValue  = errors.New("value was not as expected")
	errUnexpectedLength = errors.New("length was not as expected")
)

var mockColors = []mockChoose{
	{label: "red", value: "#FF0000"},
	{label: "green", value: "#00FF00"},
	{label: "blue", value: "#0000FF"},
}

type mockChoose struct {
	label string
	value string
}

func (s mockChoose) Print() string {
	return s.label
}

// mockFormSubmitReq is a mock request to submit a form
func mockFormSubmitReq() *http.Request {
	u := "http://localhost:8080/submit/test?name=Jane+Doe&bio=I+am+a+writer&favorite+number=27&color=%2300ff00&profile+picture=&shh=cat+and+mouse&birthday=1980-01-01&clock+in=10%3A15&vacation+start=1999-01-01T10%3A15&sel+color=0&radio+color=1&cb+color=0&cb+color=2"
	r, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		panic(err)
	}
	return r
}

func mockEmptyForm() HTMLElement {
	form := NewHTMLElement("form")
	form.Children = make([]HTMLElement, 12)

	form.Children[0] = NewInputText("name", InputKindText, 0, 30)
	form.Children[1] = NewInputTextarea("bio", 0, 30, 30, 5)
	form.Children[2] = NewInputNumber("favorite number", 0, 30)
	form.Children[3] = NewInputColor("color")
	form.Children[4] = NewInputFile("profile picture")
	form.Children[5] = NewInputHidden("shh", "cat and mouse")
	form.Children[6] = NewInputDate("birthday", nil, nil)
	form.Children[7] = NewInputTime("clock in", nil, nil)
	form.Children[8] = NewInputDateTime("vacation start", nil, nil)
	form.Children[9] = NewInputSel("sel color", mockColors)
	form.Children[10] = NewInputRadio("radio color", mockColors)
	form.Children[11] = NewInputCheckbox("cb color", mockColors)

	form.EnsureAttributes()
	form.Attributes["action"] = "/submit/test"

	submit := NewHTMLElement("button")
	submit.EnsureAttributes()
	submit.Attributes["type"] = "submit"
	submit.InnerText = "Submit Me!"

	form.Children = append(form.Children, submit)

	return form
}

func formValEq(t *testing.T, m map[string][]string, i int, key string, value string) {
	if len(m[key]) < i {
		fmt.Println(key)
		fmt.Println(value)
		t.Fatal(errUnexpectedLength)
	}
	if m[key][i] != value {
		fmt.Println(key)
		fmt.Printf("'%s' != '%s'", value, m[key][i])
		t.Fatal(errUnexpectedValue)
	}
}

func eq[T comparable](t *testing.T, a, b T) {
	if a == b {
		return
	}
	fmt.Println(a, b)
	t.Fatal(errNotEqual)
}

func eqIndexes(t *testing.T, a, b []int) {
	eq(t, len(a), len(b))
	for i := range a {
		eq(t, a[i], b[i])
	}
}

func noErr(t *testing.T, err error) {
	if err == nil {
		return
	}
	fmt.Printf(err.Error())
	t.Fatal(errFoundErr)
}

func TestParseRequest(t *testing.T) {
	el := mockEmptyForm()
	r := mockFormSubmitReq()
	el.FillForm(r)
	formValues := el.Form()
	formValEq(t, formValues, 0, "name", "Jane Doe")
	formValEq(t, formValues, 0, "bio", "I am a writer")
	formValEq(t, formValues, 0, "favorite number", "27")
	formValEq(t, formValues, 0, "color", "#00ff00")
	formValEq(t, formValues, 0, "shh", "cat and mouse")
	formValEq(t, formValues, 0, "birthday", "1980-01-01")
	formValEq(t, formValues, 0, "clock in", "10:15")
	formValEq(t, formValues, 0, "vacation start", "1999-01-01T10:15")
	formValEq(t, formValues, 0, "sel color", "0")
	formValEq(t, formValues, 0, "radio color", "1")
	formValEq(t, formValues, 0, "cb color", "0")
	formValEq(t, formValues, 1, "cb color", "2")
}

func TestParseString(t *testing.T) {
	form := mockEmptyForm()
	r := mockFormSubmitReq()
	form.FillForm(r)
	formValues := form.Form()

	s, err := formValues.String("name")
	noErr(t, err)
	eq(t, s, "Jane Doe")
}

func TestParseNumber(t *testing.T) {
	form := mockEmptyForm()
	r := mockFormSubmitReq()
	form.FillForm(r)
	formValues := form.Form()

	n, err := formValues.Int("favorite number")
	noErr(t, err)
	eq(t, n, 27)
}

func TestParseDate(t *testing.T) {
	form := mockEmptyForm()
	r := mockFormSubmitReq()
	form.FillForm(r)
	formValues := form.Form()

	d, err := formValues.Date("birthday")
	noErr(t, err)
	expectedTime, _ := time.Parse(TimeFormatHTMLDate, "1980-01-01")
	eq(t, *d, expectedTime)
}

func TestParseTime(t *testing.T) {
	form := mockEmptyForm()
	r := mockFormSubmitReq()
	form.FillForm(r)
	formValues := form.Form()

	d, err := formValues.Time("clock in")
	noErr(t, err)
	expectedTime, _ := time.Parse(TimeFormatHTMLTime, "10:15")
	eq(t, *d, expectedTime)
}

func TestParseDateTime(t *testing.T) {
	form := mockEmptyForm()
	r := mockFormSubmitReq()
	form.FillForm(r)
	formValues := form.Form()

	d, err := formValues.DateTime("vacation start")
	noErr(t, err)
	expectedTime, _ := time.Parse(TimeFormatHTMLDateTime, "1999-01-01T10:15")
	eq(t, *d, expectedTime)
}

func TestFormIndex(t *testing.T) {
	form := mockEmptyForm()
	r := mockFormSubmitReq()
	form.FillForm(r)
	formValues := form.Form()

	indexes, err := formValues.Indexes("cb color")
	noErr(t, err)
	expectedIndexes := []int{0, 2}
	eqIndexes(t, indexes, expectedIndexes)
}

func TestParseSelection(t *testing.T) {
	form := mockEmptyForm()
	r := mockFormSubmitReq()
	form.FillForm(r)
	formValues := form.Form()

	selectedColor, err := FormSelected(formValues, "sel color", mockColors)
	noErr(t, err)

	if len(selectedColor) != 1 {
		t.Fatal(errUnexpectedLength)
	}

	eq(t, selectedColor[0].Print(), mockColors[0].Print())
}
