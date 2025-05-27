package element

import "fmt"

type HTMLChoice interface {
	Print() string
	Value() any
}

type HTMLInputSelect[T HTMLChoice] struct {
	ID           string
	Label        string
	Placeholder  string
	PromptSelect bool
	Key          string
	Disabled     bool
	Required     bool
	Value        *T
	Options      []T
}

func NewHTMLInputSelect[T HTMLChoice](name string, Options []T) HTMLInputSelect[T] {
	return HTMLInputSelect[T]{
		ID:           fmt.Sprintf("id-%s", name),
		Label:        name,
		Placeholder:  "",
		PromptSelect: true,
		Disabled:     false,
		Required:     false,
		Key:          name,
		Value:        nil,
		Options:      Options,
	}
}
