package element

import (
	"fmt"
	"time"
)

func HTMLPrintDate(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02")
}

type HTMLInputDate struct {
	ID       string
	Label    string
	Key      string
	Value    *time.Time
	Disabled bool
	Required bool
	Readonly bool
	Min      *time.Time
	Max      *time.Time
}

func NewHTMLInputDate(name string) HTMLInputDate {
	return HTMLInputDate{
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
