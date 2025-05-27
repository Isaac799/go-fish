package element

import "fmt"

type HTMLInputText struct {
	ID          string
	Label       string
	Placeholder string
	Key         string
	Value       string
	Disabled    bool
	Pattern     string
	MinLen      int
	MaxLen      int
}

func NewHTMLInputText(name string) HTMLInputText {
	return HTMLInputText{
		ID:          fmt.Sprintf("id-%s", name),
		Label:       name,
		Placeholder: "",
		Key:         fmt.Sprintf("%s", name),
		Value:       "",
		Disabled:    false,
		Pattern:     "",
		MinLen:      0,
		MaxLen:      50,
	}
}
