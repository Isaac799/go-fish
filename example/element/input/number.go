package element

import "fmt"

type InputNumber struct {
	ID          string
	Label       string
	Placeholder string
	Key         string
	Value       string
	Disabled    bool
	Required    bool
	Readonly    bool
	Step        float32
	Min         int
	Max         int
}

func NewInputNumber(name string) InputNumber {
	return InputNumber{
		ID:          fmt.Sprintf("id-%s", name),
		Label:       name,
		Placeholder: "",
		Disabled:    false,
		Required:    false,
		Readonly:    false,
		Key:         name,
		Step:        0,
		Value:       "",
		Min:         -100,
		Max:         100,
	}
}
