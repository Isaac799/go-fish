package bridge

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"
)

// Printable ensures that a list of items has a way to be
// displayed in the template. Is used for select, radio, and checkbox
type Printable interface {
	Print() string
}

var (
	// ErrNotInputKindTime is returned when setting time for an input that is not time related
	ErrNotInputKindTime = errors.New("can only set time on a date time or datetime input kind")
)

// InputKind is a pre defined list of html inputs that I support
type InputKind string

const (
	// InputKindText is for input tag with attribute type of text
	InputKindText InputKind = "text"
	// InputKindPassword is for input tag with attribute type of password
	InputKindPassword InputKind = "password"
	// InputKindEmail is for input tag with attribute type of email
	InputKindEmail InputKind = "email"
	// InputKindSearch is for input tag with attribute type of search
	InputKindSearch InputKind = "search"
	// InputKindTel is for input tag with attribute type of tel
	InputKindTel InputKind = "tel"
	// InputKindURL is for input tag with attribute type of url
	InputKindURL InputKind = "url"

	// InputKindTextarea is for textarea tag with col and row attributes
	InputKindTextarea InputKind = "textarea"
	// InputKindNum is for input tag with attribute type of num
	InputKindNum InputKind = "number"
	// InputKindColor is for input tag with attribute type of color
	InputKindColor InputKind = "color"
	// InputKindHidden is for input tag with attribute type of hidden
	InputKindHidden InputKind = "hidden"
	// InputKindFile is for input tag with attribute type of file
	InputKindFile InputKind = "file"

	// InputKindDate is for input tag with attribute type of date
	InputKindDate InputKind = "date"
	// InputKindTime is for input tag with attribute type of time
	InputKindTime InputKind = "time"
	// InputKindDateTime is for input tag with attribute type of datetime
	InputKindDateTime InputKind = "datetime-local"

	// InputKindSel is for select tag with options
	InputKindSel InputKind = "select"
	// InputKindRadio is for input tags with attribute type of radio
	InputKindRadio InputKind = "radio"
	// InputKindCheckbox is for input tags with attribute type of checkbox
	InputKindCheckbox InputKind = "checkbox"
)

func newInput(kind InputKind, name string) HTMLElement {
	id := fmt.Sprintf("id-%s", name)

	if kind == InputKindHidden {
		return HTMLElement{
			Tag:         "input",
			SelfClosing: true,
			Attributes: map[AttributeKey]string{
				ID:    id,
				Type:  string(kind),
				Name:  name,
				Value: "",
			},
		}
	}

	label := HTMLElement{
		Tag:       "label",
		InnerText: name,
		Attributes: map[AttributeKey]string{
			For: id,
		},
	}

	input := HTMLElement{
		Tag:         "input",
		SelfClosing: true,
		Attributes: map[AttributeKey]string{
			ID:    id,
			Type:  string(kind),
			Name:  name,
			Value: "",
		},
	}

	children := make([]HTMLElement, 2)
	children[0] = label
	children[1] = input

	div := HTMLElement{
		Tag:      "div",
		Children: children,
	}

	return div
}

func newTextArea(name string, col, row uint) HTMLElement {
	id := fmt.Sprintf("id-%s", name)

	label := HTMLElement{
		Tag:       "label",
		InnerText: name,
		Attributes: map[AttributeKey]string{
			For: id,
		},
	}

	input := HTMLElement{
		Tag:       "textarea",
		InnerText: "",
		Attributes: map[AttributeKey]string{
			ID:   id,
			Type: string(InputKindTextarea),
			Name: name,
			Col:  fmt.Sprintf("%d", col),
			Row:  fmt.Sprintf("%d", row),
		},
	}

	children := make([]HTMLElement, 2)
	children[0] = label
	children[1] = input

	div := HTMLElement{
		Tag:      "div",
		Children: children,
	}

	return div
}

func newSelect[T Printable](name string, options []T) HTMLElement {
	id := fmt.Sprintf("id-%s", name)

	label := HTMLElement{
		Tag:       "label",
		InnerText: fmt.Sprintf("Choose %s:", name),
		Attributes: map[AttributeKey]string{
			For: id,
		},
	}

	input := HTMLElement{
		Tag: "select",
		Attributes: map[AttributeKey]string{
			ID:    id,
			Name:  name,
			Value: "",
		},
	}

	for i, option := range options {
		id := fmt.Sprintf("id-%s", name)
		el := HTMLElement{
			Tag:       "option",
			InnerText: option.Print(),
			Attributes: map[AttributeKey]string{
				ID:    id,
				Name:  name,
				Value: fmt.Sprintf("%d", i),
			},
		}
		input.Children = append(input.Children, el)
	}

	children := make([]HTMLElement, 2)
	children[0] = label
	children[1] = input

	div := HTMLElement{
		Tag:      "div",
		Children: children,
	}

	return div
}

func newRadioCheckbox[T Printable](kind InputKind, name string, options []T) HTMLElement {
	legend := HTMLElement{
		Tag:       "legend",
		InnerText: fmt.Sprintf("Choose %s:", name),
	}

	children := make([]HTMLElement, len(options)+1)
	children[0] = legend

	for i, option := range options {
		label := HTMLElement{
			Tag:       "label",
			InnerText: option.Print(),
			Attributes: map[AttributeKey]string{
				ID: fmt.Sprintf("id-%s", name),
			},
		}
		input := HTMLElement{
			Tag:         "input",
			SelfClosing: true,
			InnerText:   option.Print(),
			Attributes: map[AttributeKey]string{
				ID:    fmt.Sprintf("id-%s", name),
				Type:  string(kind),
				Name:  name,
				Value: fmt.Sprintf("%d", i),
			},
		}

		optionChildren := make([]HTMLElement, 2)
		optionChildren[0] = label
		optionChildren[1] = input

		div := HTMLElement{
			Tag:      "div",
			Children: optionChildren,
		}
		children[i+1] = div
	}

	div := HTMLElement{
		Tag:      "div",
		Children: children,
	}

	return div
}

// InputPickOne allows comparison in templates of a a single selection
// say from select or radio
func InputPickOne(a, b int) bool {
	return a == b
}

// InputPickMany allows comparison in templates of a multiple selection
// say from select multiple or checkbox
func InputPickMany(a []int, b int) bool {
	return slices.Contains(a, b)
}

// InputJoinComma will join an input with commas. Useful when
// setting the value of a select with multiple attribute
func InputJoinComma(s []string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.Join(s, ",")
}

// NewInputText is a div element with labeled text child
// To be called with [ text | password | email | search | tel | url ]
func NewInputText(name string, kind InputKind, minLen, maxLen uint) HTMLElement {
	el := newInput(kind, name)
	input := el.FindFirst(LikeInput)
	input.Attributes["minLength"] = fmt.Sprintf("%d", minLen)
	input.Attributes["maxLength"] = fmt.Sprintf("%d", maxLen)
	return el
}

// NewInputTextarea is a div element with labeled textarea child
func NewInputTextarea(name string, minLen, maxLen, col, row uint) HTMLElement {
	el := newTextArea(name, col, row)
	input := el.FindFirst(LikeInput)
	input.Attributes["minLength"] = fmt.Sprintf("%d", minLen)
	input.Attributes["maxLength"] = fmt.Sprintf("%d", maxLen)
	return el
}

// NewInputNumber is a div element with labeled number child
func NewInputNumber(name string, min, max float32) HTMLElement {
	el := newInput(InputKindNum, name)
	input := el.FindFirst(LikeInput)
	input.Attributes["min"] = fmt.Sprintf("%f", min)
	input.Attributes["max"] = fmt.Sprintf("%f", max)
	return el
}

// NewInputColor is a div element with labeled color child
func NewInputColor(name string) HTMLElement {
	el := newInput(InputKindColor, name)
	return el
}

// NewInputHidden gives a hidden input element
func NewInputHidden(name string, value string) HTMLElement {
	el := newInput(InputKindHidden, name)
	el.Attributes["value"] = value
	return el
}

// NewInputFile is a div element with labeled file child
func NewInputFile(name string) HTMLElement {
	el := newInput(InputKindFile, name)
	return el
}

// NewInputDate is a div element with labeled date child
func NewInputDate(name string, min, max *time.Time) HTMLElement {
	el := newInput(InputKindDate, name)
	input := el.FindFirst(LikeInput)
	input.Attributes["min"] = PrintDate(min)
	input.Attributes["max"] = PrintDate(max)
	return el
}

// NewInputTime is a div element with labeled time child
func NewInputTime(name string, min, max *time.Time) HTMLElement {
	el := newInput(InputKindTime, name)
	input := el.FindFirst(LikeInput)
	input.Attributes["min"] = PrintTime(min)
	input.Attributes["max"] = PrintTime(max)
	return el
}

// NewInputDateTime is a div element with labeled datetime-local child
func NewInputDateTime(name string, min, max *time.Time) HTMLElement {
	el := newInput(InputKindDateTime, name)
	input := el.FindFirst(LikeInput)
	input.Attributes["min"] = PrintDateTime(min)
	input.Attributes["max"] = PrintDateTime(max)
	return el
}

// NewInputSel is a div element with labeled select child
// One to many selections allowed.
func NewInputSel[T Printable](name string, options []T) HTMLElement {
	el := newSelect(name, options)
	return el
}

// NewInputRadio is a div element with labeled radio input children
// One selection allowed.
func NewInputRadio[T Printable](name string, options []T) HTMLElement {
	el := newRadioCheckbox(InputKindRadio, name, options)
	return el
}

// NewInputCheckbox is a div element with labeled checkbox input children.
// Many selections allowed.
func NewInputCheckbox[T Printable](name string, options []T) HTMLElement {
	el := newRadioCheckbox(InputKindCheckbox, name, options)
	return el
}
