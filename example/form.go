package main

import (
	"net/http"

	"github.com/Isaac799/go-fish/pkg/bridge"
)

type exampleColor struct {
	label string
	value string
}

func (s exampleColor) String() string {
	return s.label
}

func exampleFormBait(_ *http.Request) *fishData {
	data := fishData{}
	form := exampleForm()
	data.Form = &form
	return &data
}

func exampleForm() bridge.HTMLElement {
	form := bridge.NewHTMLElement("form")
	form.Children = make([]bridge.HTMLElement, 12)

	elText := bridge.NewInputText("name", bridge.InputKindText, 0, 30)
	elText.SetFirstValue("Jane Doe")
	form.Children[0] = elText

	elTextarea := bridge.NewInputTextarea("bio", 0, 30, 30, 5)
	elTextarea.SetFirstValue("I am a writer.")
	form.Children[1] = elTextarea

	elNumber := bridge.NewInputNumber("favorite number", 0, 30)
	elNumber.SetFirstValue("27")
	form.Children[2] = elNumber

	elColor := bridge.NewInputColor("color")
	elColor.SetFirstValue("#00FF00")
	form.Children[3] = elColor

	form.Children[4] = bridge.NewInputFile("profile picture")

	form.Children[5] = bridge.NewInputHidden("shh", "cat and mouse")

	elDate := bridge.NewInputDate("birthday", nil, nil)
	elDate.SetFirstValue("1980-01-01")
	form.Children[6] = elDate

	elTime := bridge.NewInputTime("clock in", nil, nil)
	elTime.SetFirstValue("10:15")
	form.Children[7] = elTime

	elDateTime := bridge.NewInputDateTime("vacation start", nil, nil)
	elDateTime.SetFirstValue("1999-01-01T10:15")
	form.Children[8] = elDateTime

	var exampleColors = []exampleColor{
		{label: "red", value: "#FF0000"},
		{label: "green", value: "#00FF00"},
		{label: "blue", value: "#0000FF"},
	}

	elSel := bridge.NewInputSelect("sel color", exampleColors)
	elSel.SetSelectOption(2, true)
	form.Children[9] = elSel

	elRadio := bridge.NewInputRadio("radio color", exampleColors)
	elRadio.SetFirstValue("true")
	form.Children[10] = elRadio

	elCheckbox := bridge.NewInputCheckbox("cb color", exampleColors)
	elCheckbox.SetFirstValue("true")
	elCheckbox.SetNthValue(3, "true")
	form.Children[11] = elCheckbox

	form.EnsureAttributes()
	form.Attributes["action"] = "/submit/test"

	submit := bridge.NewHTMLElement("button")
	submit.EnsureAttributes()
	submit.Attributes["type"] = "submit"
	submit.InnerText = "Submit Me!"

	form.Children = append(form.Children, submit)

	return form
}
