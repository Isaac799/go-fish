package element

import "fmt"

type InputTextArea struct {
	ID          string
	Label       string
	Placeholder string
	Key         string
	Value       string
	Disabled    bool
	Required    bool
	Readonly    bool
	Pattern     string
	MinLen      int
	MaxLen      int
	Col         uint
	Row         uint
}

func NewInputTextArea(name string) InputTextArea {
	return InputTextArea{
		ID:          fmt.Sprintf("id-%s", name),
		Label:       name,
		Placeholder: "",
		Key:         name,
		Value:       "",
		Disabled:    false,
		Required:    false,
		Readonly:    false,
		Pattern:     "",
		MinLen:      0,
		MaxLen:      200,
		Col:         30,
		Row:         10,
	}
}
