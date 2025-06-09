// Package bridge is an simple abstract on HTML concepts to close the gap
// between Go and HTML. It can generate the HTML and provide functionality
// to grantee the user experience matches what the server expects.
// Imagine a well defined form with seamless front-back validation, or a
// powerful table with features (e.g. pagination) implemented directly. Or
// even navigation links that match our available routes.
// Also bridges are water related thus matching this fish themed repo :-)
package bridge

import (
	"errors"
	"time"
)

var (
	// ErrAttributesNil is given when trying to preform an action on nil attributes
	ErrAttributesNil = errors.New("attributes are nil")
	// ErrValueDoesNotExist is given when trying to access an the value entry
	ErrValueDoesNotExist = errors.New("value does not exist")
)

// Attributes is the attributes of an html element
type Attributes map[string]string

// SetTime will set the value attribute of an element to a specific time
// depending on the type of input. Err only if the input is not a time.
// If its attributes are nil it will create them.
// Nil times are empty string, which HTML treats as empty.
func (attrs Attributes) SetTime(kind string, t *time.Time) error {
	if t == nil {
		attrs["value"] = ""
		return nil
	}
	m := make(map[string]HTMLPrintTime, 3)
	m[InputKindDate] = PrintDate
	m[InputKindTime] = PrintTime
	m[InputKindDateTime] = PrintDateTime

	fn, exists := m[kind]
	if !exists {
		return ErrNotInputKindTime
	}

	attrs["value"] = fn(t)
	return nil
}
