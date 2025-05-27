package element

import (
	"fmt"
)

type HTMLInputRadio[T HTMLChoice] struct {
	ID          string
	Label       string
	Placeholder string
	Key         string
	Legend      string
	Disabled    bool
	Value       *T
	Options     []T
}

func NewHTMLInputRadio[T HTMLChoice](name string, Options []T) HTMLInputRadio[T] {
	return HTMLInputRadio[T]{
		ID:          fmt.Sprintf("id-%s", name),
		Legend:      fmt.Sprintf("Choose a %s:", name),
		Label:       name,
		Placeholder: "",
		Disabled:    false,
		Key:         fmt.Sprintf("%s"),
		Value:       nil,
		Options:     Options,
	}
}
