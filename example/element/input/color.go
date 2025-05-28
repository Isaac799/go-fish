package element

import "fmt"

type InputColor struct {
	ID       string
	Label    string
	Key      string
	Value    string
	Disabled bool
	Required bool
}

func NewInputColor(name string) InputColor {
	return InputColor{
		ID:       fmt.Sprintf("id-%s", name),
		Label:    name,
		Key:      name,
		Value:    "",
		Disabled: false,
		Required: false,
	}
}
