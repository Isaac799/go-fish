package main

import (
	"net/http"
	"strconv"
	"time"

	bridge "github.com/Isaac799/go-fish/example/bridge"
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
	Text     bridge.HTMLElement
	Textarea bridge.HTMLElement
	Num      bridge.HTMLElement
	Sel      bridge.HTMLElement
	Radio    bridge.HTMLElement
	Checkbox bridge.HTMLElement
	Date     bridge.HTMLElement
	Time     bridge.HTMLElement
	DateTime bridge.HTMLElement
	Color    bridge.HTMLElement
	Hidden   bridge.HTMLElement
	File     bridge.HTMLElement
}

func inputs(_ *http.Request) any {
	options := []simpleSelect{
		{label: "red", value: "#FF0000"},
		{label: "blue", value: "#0000FF"},
		{label: "green", value: "#00FF00"},
	}

	b := variousInputs{
		Text:     bridge.NewInputText("name", bridge.InputKindText, 0, 30),
		Textarea: bridge.NewInputTextarea("bio", 0, 30, 30, 5),
		Num:      bridge.NewInputNumber("cell", 0, 30),
		Color:    bridge.NewInputColor("favorite color"),
		File:     bridge.NewInputFile("profile picture"),
		Hidden:   bridge.NewInputHidden("shh", "cat and mouse"),

		Date:     bridge.NewInputDate("birthday", nil, nil),
		Time:     bridge.NewInputTime("clock in", nil, nil),
		DateTime: bridge.NewInputDateTime("vacation start", nil, nil),

		Sel:      bridge.NewInputSel("favorite color", options),
		Radio:    bridge.NewInputRadio("favorite color", options),
		Checkbox: bridge.NewInputCheckbox("favorite colors", options),
	}
	return b
}
