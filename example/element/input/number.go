package element

import "fmt"

type HTMLInputNumber struct {
	ID          string
	Label       string
	Placeholder string
	Key         string
	Value       string
	Disabled    bool
	Required    bool
	Readonly    bool
	Min         int
	Max         int
}

func NewHTMLInputNumber(name string) HTMLInputNumber {
	return HTMLInputNumber{
		ID:          fmt.Sprintf("id-%s", name),
		Label:       name,
		Placeholder: "",
		Disabled:    false,
		Required:    false,
		Readonly:    false,
		Key:         name,
		Value:       "",
		Min:         -100,
		Max:         100,
	}
}
