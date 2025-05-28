package element

import "fmt"

type InputHidden struct {
	ID    string
	Key   string
	Value string
}

func NewInputHidden(name string, value string) InputHidden {
	return InputHidden{
		ID:    fmt.Sprintf("id-%s", name),
		Key:   name,
		Value: value,
	}
}
