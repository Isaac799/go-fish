package main

import (
	"net/http"
	"strconv"
	"time"

	element "github.com/Isaac799/go-fish/example/element"
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
	Text     element.HTMLElement
	Textarea element.HTMLElement
	Num      element.HTMLElement
	Sel      element.HTMLElement
	Radio    element.HTMLElement
	Checkbox element.HTMLElement
	Date     element.HTMLElement
	Time     element.HTMLElement
	DateTime element.HTMLElement
	Color    element.HTMLElement
	Hidden   element.HTMLElement
	File     element.HTMLElement
}

func inputs(_ *http.Request) any {
	options := []simpleSelect{
		{label: "red", value: "#FF0000"},
		{label: "blue", value: "#0000FF"},
		{label: "green", value: "#00FF00"},
	}

	b := variousInputs{
		Text:     element.NewInputText("name", element.InputKindText, 0, 30),
		Textarea: element.NewInputTextarea("bio", 0, 30, 30, 5),
		Num:      element.NewInputNumber("cell", 0, 30),
		Color:    element.NewInputColor("favorite color"),
		File:     element.NewInputFile("profile picture"),
		Hidden:   element.NewInputHidden("shh", "cat and mouse"),

		Date:     element.NewInputDate("birthday", nil, nil),
		Time:     element.NewInputTime("clock in", nil, nil),
		DateTime: element.NewInputDateTime("vacation start", nil, nil),

		Sel:      element.NewInputSel("favorite color", options),
		Radio:    element.NewInputRadio("favorite color", options),
		Checkbox: element.NewInputCheckbox("favorite colors", options),
	}
	return b
}
