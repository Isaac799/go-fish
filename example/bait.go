package main

import (
	"net/http"
	"strconv"
	"time"

	element "github.com/Isaac799/go-fish/example/element/input"
)

type user struct {
	ID        int
	FirstName string
	LastName  string
}

type water struct {
	ServerTime time.Time
	RotateDeg  int
}

func queriedSeason(r *http.Request) any {
	season, ok := r.Context().Value(queryCtxKey).(string)
	if !ok {
		return nil
	}
	return season
}

func userInfo(r *http.Request) any {
	user, ok := r.Context().Value(userCtxKey).(user)
	if !ok {
		return nil
	}
	return user
}

func waterInfo(r *http.Request) any {
	posStr := r.URL.Query().Get("pos")
	offsetStr := r.URL.Query().Get("off")

	w := water{
		RotateDeg:  0,
		ServerTime: time.Now(),
	}

	off, err := strconv.Atoi(offsetStr)
	if err != nil {
		return w
	}
	pos, err := strconv.Atoi(posStr)
	if err != nil {
		return w
	}

	w.RotateDeg = pos + off

	return w
}

type simpleSelect struct {
	label string
	value string
}

func (s simpleSelect) Print() string {
	return s.label
}
func (s simpleSelect) Value() any {
	return s.value
}

type variousInputs struct {
	Text     element.InputText
	Textarea element.InputTextArea
	Num      element.InputNumber
	Sel      element.InputSelect[simpleSelect]
	Radio    element.InputRadio[simpleSelect]
	Checkbox element.InputCheckbox[simpleSelect]
	Date     element.InputDate
	Time     element.InputTime
	DateTime element.InputDateTime
	Color    element.InputColor
}

func inputs(_ *http.Request) any {
	textEl := element.NewInputText("name")
	numberEl := element.NewInputNumber("age")
	textareaEl := element.NewInputTextArea("bio")
	selectEl := element.NewInputSelect("favorite color", []simpleSelect{
		{label: "red", value: "#FF0000"},
		{label: "blue", value: "#0000FF"},
		{label: "green", value: "#00FF00"},
	})
	radioEl := element.NewInputRadio("second favorite color", []simpleSelect{
		{label: "red", value: "#FF0000"},
		{label: "blue", value: "#0000FF"},
		{label: "green", value: "#00FF00"},
	})
	cbEl := element.NewInputCheckbox("third favorite color", []simpleSelect{
		{label: "red", value: "#FF0000"},
		{label: "blue", value: "#0000FF"},
		{label: "green", value: "#00FF00"},
	})

	dateEl := element.NewInputDate("birthday")
	timeEl := element.NewInputTime("clock in")
	datetimeEl := element.NewInputDateTime("vacation start")
	colorEl := element.NewInputColor("backup color")

	return variousInputs{
		Text:     textEl,
		Textarea: textareaEl,
		Num:      numberEl,
		Sel:      selectEl,
		Radio:    radioEl,
		Checkbox: cbEl,
		Date:     dateEl,
		Time:     timeEl,
		DateTime: datetimeEl,
		Color:    colorEl,
	}
}
