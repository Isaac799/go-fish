package element

import (
	"fmt"
	"time"
)

func HTMLPrintDateTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02T15:04")
}

type HTMLInputDateTime struct {
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

func NewHTMLInputDateTime(name string) HTMLInputDateTime {
	return HTMLInputDateTime{
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
