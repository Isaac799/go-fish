package main

import "github.com/Isaac799/go-fish/example/bridge"

type exampleColor struct {
	label string
	value string
}

func (s exampleColor) Print() string {
	return s.label
}

func exampleForm() bridge.HTMLElement {
	form := bridge.NewHTMLElement("form")
	form.Children = make([]bridge.HTMLElement, 12)

	elText := bridge.NewInputText("name", bridge.InputKindText, 0, 30)
	elText.SetValue(1, "Jane Doe")
	form.Children[0] = elText

	elTextarea := bridge.NewInputTextarea("bio", 0, 30, 30, 5)
	elTextarea.SetValue(1, "I am a writer.")
	form.Children[1] = elTextarea

	elNumber := bridge.NewInputNumber("favorite number", 0, 30)
	elNumber.SetValue(1, "27")
	form.Children[2] = elNumber

	elColor := bridge.NewInputColor("color")
	elColor.SetValue(1, "#00FF00")
	form.Children[3] = elColor

	form.Children[4] = bridge.NewInputFile("profile picture")

	form.Children[5] = bridge.NewInputHidden("shh", "cat and mouse")

	elDate := bridge.NewInputDate("birthday", nil, nil)
	elDate.SetValue(1, "1980-01-01")
	form.Children[6] = elDate

	elTime := bridge.NewInputTime("clock in", nil, nil)
	elTime.SetValue(1, "10:15")
	form.Children[7] = elTime

	elDateTime := bridge.NewInputDateTime("vacation start", nil, nil)
	elDateTime.SetValue(1, "1999-01-01T10:15")
	form.Children[8] = elDateTime

	var exampleColors = []exampleColor{
		{label: "red", value: "#FF0000"},
		{label: "green", value: "#00FF00"},
		{label: "blue", value: "#0000FF"},
	}

	elSel := bridge.NewInputSel("sel color", exampleColors)
	elSel.SetValue(1, "2")
	form.Children[9] = elSel

	elRadio := bridge.NewInputRadio("radio color", exampleColors)
	elRadio.SetChecked(1, true)
	form.Children[10] = elRadio

	elCheckbox := bridge.NewInputCheckbox("cb color", exampleColors)
	elCheckbox.SetChecked(1, true)
	elCheckbox.SetChecked(3, true)
	form.Children[11] = elCheckbox

	form.Attributes = form.Attributes.Ensure()
	form.Attributes["action"] = "/submit/test"

	submit := bridge.NewHTMLElement("button")
	submit.Attributes = submit.Attributes.Ensure()
	submit.Attributes["type"] = "submit"
	submit.InnerText = "Submit Me!"

	form.Children = append(form.Children, submit)

	return form
}
