package element

import (
	"fmt"
	"time"
)

func HTMLPrintTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("15:04")
}

type HTMLInputTime struct {
	ID       string
	Label    string
	Key      string
	Value    *time.Time
	Disabled bool
	Readonly bool
	Required bool
	Min      *time.Time
	Max      *time.Time
}

func NewHTMLInputTime(name string) HTMLInputTime {
	return HTMLInputTime{
		ID:       fmt.Sprintf("id-%s", name),
		Label:    name,
		Disabled: false,
		Required: false,
		Readonly: false,
		Key:      name,
		Value:    nil,
		Min:      nil,
		Max:      nil,
	}
}
