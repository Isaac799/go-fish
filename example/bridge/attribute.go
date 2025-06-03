// Package bridge is an simple abstract on HTML concepts to close the gap
// between Go and HTML. It can generate the HTML and provide functionality
// to grantee the user experience matches what the server expects.
// Imagine a well defined form with seamless front-back validation, or a
// powerful table with features (e.g. pagination) implemented directly. Or
// even navigation links that match our available routes.
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

// AttributeKey is a is the key in an html element kye-value pair
type AttributeKey string

// Attributes is the attributes of an html element
type Attributes map[AttributeKey]string

const (
	// ID is the html attribute id
	ID AttributeKey = "id"
	// Name is the html attribute name. Key in a forms key-value pair
	Name AttributeKey = "name"
	// For is the html attribute for. Used by labels
	For AttributeKey = "for"
	// Type is the html attribute type. Used by inputs.
	Type AttributeKey = "type"
	// Value is the html attribute value. Value in a forms key-value pair
	Value AttributeKey = "value"
	// Required is the html attribute required
	Required AttributeKey = "required"
	// Readonly is the html attribute readonly
	Readonly AttributeKey = "readonly"
	// Col is the html attribute cols for textarea
	Col AttributeKey = "cols"
	// Row is the html attribute rows for textarea
	Row AttributeKey = "rows"
	// Size is the html attribute size for select listbox behavior
	Size AttributeKey = "size"
	// Step is the html attribute step for number input type
	Step AttributeKey = "step"
	// Min is the html attribute min validation
	Min AttributeKey = "min"
	// Max is the html attribute max validation
	Max AttributeKey = "max"
	// MinLength is the html attribute minlength validation
	MinLength AttributeKey = "minlength"
	// MaxLength is the html attribute maxlength validation
	MaxLength AttributeKey = "maxlength"
	// Disabled is the html attribute disabled
	Disabled AttributeKey = "disabled"
	// Multiple is the html attribute multiple for select enabling many
	Multiple AttributeKey = "multiple"
	// Accept is the html attribute accept for file input type
	Accept AttributeKey = "accept"
	// Pattern is the html attribute pattern for validation
	Pattern AttributeKey = "pattern"
	// Placeholder is the html attribute placeholder
	Placeholder AttributeKey = "placeholder"
	// HRef is the html attribute href
	HRef AttributeKey = "href"
)

// Ensure helps safe access of attributes.
func (a Attributes) Ensure() Attributes {
	if a != nil {
		return a
	}
	return make(map[AttributeKey]string, 0)
}

// SetChecked sets the checked attribute of an element.
// If its attributes are nil it will create them.
func (a Attributes) SetChecked(b bool) error {
	if a == nil {
		return ErrAttributesNil
	}
	if b {
		a["checked"] = "true"
	} else {
		a["checked"] = "false"
	}
	return nil
}

// SetTime will set the value attribute of an element to a specific time
// depending on the type of input. Err only if the input is not a time.
// If its attributes are nil it will create them.
// Nil times are empty string, which HTML treats as empty.
func (a Attributes) SetTime(kind InputKind, t *time.Time) error {
	if t == nil {
		a["value"] = ""
		return nil
	}
	m := make(map[InputKind]HTMLPrintTime, 3)
	m[InputKindDate] = PrintDate
	m[InputKindTime] = PrintTime
	m[InputKindDateTime] = PrintDateTime

	fn, exists := m[kind]
	if !exists {
		return ErrNotInputKindTime
	}

	a["value"] = fn(t)
	return nil
}
