package element

import (
	"fmt"
	"slices"
)

type Printable interface {
	Print() string
}

func InputPickOne(a, b int) bool {
	return a == b
}

func InputPickMany(a []int, b int) bool {
	return slices.Contains(a, b)
}

type InputSelect[T Printable] struct {
	ID           string
	Label        string
	Placeholder  string
	PromptSelect bool
	Key          string
	Disabled     bool
	Multiple     bool
	Required     bool
	Size         int
	// Value is selected index(s) of Options. Typically just one, unless Multiple is true
	Value   []int
	Options []T
}

func NewInputSelect[T Printable](name string, Options []T) InputSelect[T] {
	return InputSelect[T]{
		ID:           fmt.Sprintf("id-%s", name),
		Label:        name,
		Placeholder:  "",
		PromptSelect: true,
		Disabled:     false,
		Required:     false,
		Multiple:     false,
		Key:          name,
		Size:         0,
		Value:        []int{},
		Options:      Options,
	}
}

func (el *InputSelect[T]) Listbox() {
	el.Size = len(el.Options)
}
