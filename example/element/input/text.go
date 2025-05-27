package element

import "fmt"

type HTMLInputText struct {
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
}

func NewHTMLInputText(name string) HTMLInputText {
	return HTMLInputText{
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
		MaxLen:      50,
	}
}
