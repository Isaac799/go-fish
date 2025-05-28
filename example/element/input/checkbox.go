package element

import (
	"fmt"
)

type InputCheckbox[T Printable] struct {
	ID       string
	Label    string
	Key      string
	Legend   string
	Disabled bool
	Required bool
	// Value is checked index(s) of Options.
	Value   []int
	Options []T
}

func NewInputCheckbox[T Printable](name string, Options []T) InputCheckbox[T] {
	return InputCheckbox[T]{
		ID:       fmt.Sprintf("id-%s", name),
		Legend:   fmt.Sprintf("Choose %s:", name),
		Label:    name,
		Disabled: false,
		Required: false,
		Key:      name,
		Value:    []int{},
		Options:  Options,
	}
}
