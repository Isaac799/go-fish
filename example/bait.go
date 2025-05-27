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
	Text     element.HTMLInputText
	Textarea element.HTMLInputTextArea
	Num      element.HTMLInputNumber
	Sel      element.HTMLInputSelect[simpleSelect]
	Radio    element.HTMLInputRadio[simpleSelect]
}

func inputs(_ *http.Request) any {
	text := element.NewHTMLInputText("name")
	num := element.NewHTMLInputNumber("age")
	textarea := element.NewHTMLInputTextArea("bio")
	sel := element.NewHTMLInputSelect("favorite color", []simpleSelect{
		{label: "red", value: "#FF0000"},
		{label: "blue", value: "#0000FF"},
		{label: "green", value: "#00FF00"},
	})
	radio := element.NewHTMLInputRadio("second favorite color", []simpleSelect{
		{label: "red", value: "#FF0000"},
		{label: "blue", value: "#0000FF"},
		{label: "green", value: "#00FF00"},
	})

	return variousInputs{
		Text:     text,
		Textarea: textarea,
		Num:      num,
		Sel:      sel,
		Radio:    radio,
	}
}
