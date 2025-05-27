package element

import "fmt"

type HTMLInputNumber struct {
	ID          string
	Label       string
	Placeholder string
	Key         string
	Value       string
	Disabled    bool
	Pattern     string
	Min         int
	Max         int
}

func NewHTMLInputNumber(name string) HTMLInputNumber {
	return HTMLInputNumber{
		ID:          fmt.Sprintf("id-%s", name),
		Label:       name,
		Placeholder: "",
		Disabled:    false,
		Key:         fmt.Sprintf("%s", name),
		Value:       "",
		Pattern:     "",
		Min:         -100,
		Max:         100,
	}
}
