package element

import "fmt"

type InputTextKind = string

const (
	InputTextKindText     InputTextKind = "text"
	InputTextKindPassword InputTextKind = "password"
	InputTextKindEmail    InputTextKind = "email"
	InputTextKindSearch   InputTextKind = "search"
	InputTextKindTel      InputTextKind = "tel"
	InputTextKindUrl      InputTextKind = "url"
)

type InputText struct {
	ID          string
	Label       string
	Kind        InputTextKind
	Placeholder string
	Key         string
	Value       string
	Disabled    bool
	Required    bool
	Readonly    bool
	Pattern     string
	MinLen      int
	MaxLen      int
}

func NewInputText(name string) InputText {
	return InputText{
		ID:          fmt.Sprintf("id-%s", name),
		Label:       name,
		Kind:        InputTextKindText,
		Placeholder: "",
		Key:         name,
		Value:       "",
		Disabled:    false,
		Required:    false,
		Readonly:    false,
		Pattern:     "",
		MinLen:      0,
		MaxLen:      50,
	}
}
