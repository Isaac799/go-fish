package element

import (
	"fmt"
)

type HTMLInputRadio[T HTMLChoice] struct {
	ID       string
	Label    string
	Key      string
	Legend   string
	Disabled bool
	Required bool
	Value    *T
	Options  []T
}

func NewHTMLInputRadio[T HTMLChoice](name string, Options []T) HTMLInputRadio[T] {
	return HTMLInputRadio[T]{
		ID:       fmt.Sprintf("id-%s", name),
		Legend:   fmt.Sprintf("Choose a %s:", name),
		Label:    name,
		Disabled: false,
		Required: false,
		Key:      fmt.Sprintf("%s"),
		Value:    nil,
		Options:  Options,
	}
}
