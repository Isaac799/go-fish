package bridge

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
)

var (
	// ErrNotInputKindTime is returned when setting time for an input that is not time related
	ErrNotInputKindTime = errors.New("can only set time on a date time or datetime input kind")
)

const (
	// InputKindText is for input tag with attribute type of text
	InputKindText = "text"
	// InputKindPassword is for input tag with attribute type of password
	InputKindPassword = "password"
	// InputKindEmail is for input tag with attribute type of email
	InputKindEmail = "email"
	// InputKindSearch is for input tag with attribute type of search
	InputKindSearch = "search"
	// InputKindTel is for input tag with attribute type of tel
	InputKindTel = "tel"
	// InputKindURL is for input tag with attribute type of url
	InputKindURL = "url"

	// InputKindTextarea is for textarea tag with col and row attributes
	InputKindTextarea = "textarea"
	// InputKindNumber is for input tag with attribute type of num
	InputKindNumber = "number"
	// InputKindColor is for input tag with attribute type of color
	InputKindColor = "color"
	// InputKindHidden is for input tag with attribute type of hidden
	InputKindHidden = "hidden"
	// InputKindFile is for input tag with attribute type of file
	InputKindFile = "file"

	// InputKindDate is for input tag with attribute type of date
	InputKindDate = "date"
	// InputKindTime is for input tag with attribute type of time
	InputKindTime = "time"
	// InputKindDateTime is for input tag with attribute type of datetime
	InputKindDateTime = "datetime-local"

	// InputKindSelect is for select tag with options
	InputKindSelect = "select"
	// InputKindRadio is for input tags with attribute type of radio
	InputKindRadio = "radio"
	// InputKindCheckbox is for input tags with attribute type of checkbox
	InputKindCheckbox = "checkbox"

	// InputKindSubmit is for input tags with attribute type of submit
	InputKindSubmit = "submit"
)

// HTMLInput is a specific kind of html element with a focus on
// user input or storing state
type HTMLInput = HTMLElement

func newInput(kind string, name string) HTMLInput {
	id := fmt.Sprintf("id-%s", name)

	if kind == InputKindHidden {
		return HTMLElement{
			Tag:         "input",
			SelfClosing: true,
			Attributes: map[string]string{
				"id":    id,
				"type":  kind,
				"name":  name,
				"value": "",
			},
		}
	}

	label := HTMLElement{
		Tag:       "label",
		InnerText: name,
		Attributes: map[string]string{
			"for": id,
		},
	}

	input := HTMLElement{
		Tag:         "input",
		SelfClosing: true,
		Attributes: map[string]string{
			"id":    id,
			"type":  kind,
			"name":  name,
			"value": "",
		},
	}

	children := make([]HTMLElement, 0, 2)
	children = append(children, label, input)

	div := HTMLElement{
		Tag:      "div",
		Children: children,
	}

	return div
}

func newTextArea(name string, col, row uint) HTMLInput {
	id := fmt.Sprintf("id-%s", name)

	label := HTMLElement{
		Tag:       "label",
		InnerText: name,
		Attributes: map[string]string{
			"for": id,
		},
	}

	input := HTMLElement{
		Tag:       "textarea",
		InnerText: "",
		Attributes: map[string]string{
			"id":   id,
			"type": string(InputKindTextarea),
			"name": name,
			"col":  fmt.Sprintf("%d", col),
			"row":  fmt.Sprintf("%d", row),
		},
	}

	children := make([]HTMLElement, 0, 2)
	children = append(children, label, input)

	div := HTMLElement{
		Tag:      "div",
		Children: children,
	}

	return div
}

func newSelect[T fmt.Stringer](name string, options []T) HTMLInput {
	id := fmt.Sprintf("id-%s", name)

	label := HTMLElement{
		Tag:       "label",
		InnerText: fmt.Sprintf("Choose %s:", name),
		Attributes: map[string]string{
			"for": id,
		},
	}

	input := HTMLElement{
		Tag: "select",
		Attributes: map[string]string{
			"id":   id,
			"name": name,
		},
	}

	for i, option := range options {
		id := fmt.Sprintf("id-%s", name)
		el := HTMLElement{
			Tag:       "option",
			InnerText: option.String(),
			Attributes: map[string]string{
				"id":    id,
				"name":  name,
				"value": fmt.Sprintf("%d", i),
			},
		}
		input.Children = append(input.Children, el)
	}

	children := make([]HTMLElement, 0, 2)
	children = append(children, label, input)

	div := HTMLElement{
		Tag:      "div",
		Children: children,
	}

	return div
}

func newRadioCheckbox[T fmt.Stringer](kind string, name string, options []T) HTMLInput {
	legend := HTMLElement{
		Tag:       "legend",
		InnerText: fmt.Sprintf("Choose %s:", name),
	}

	children := make([]HTMLElement, 0, len(options)+1)
	children = append(children, legend)

	for i, option := range options {
		label := HTMLElement{
			Tag:       "label",
			InnerText: option.String(),
			Attributes: map[string]string{
				"id": fmt.Sprintf("id-%s", name),
			},
		}
		input := HTMLElement{
			Tag:         "input",
			SelfClosing: true,
			InnerText:   option.String(),
			Attributes: map[string]string{
				"id":    fmt.Sprintf("id-%s", name),
				"type":  kind,
				"name":  name,
				"value": fmt.Sprintf("%d", i),
			},
		}

		optionChildren := make([]HTMLElement, 0, 2)
		optionChildren = append(optionChildren, label, input)

		div := HTMLElement{
			Tag:      "div",
			Children: optionChildren,
		}
		children = append(children, div)
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
func NewInputText(name string, kind string, minLen, maxLen uint) HTMLInput {
	el := newInput(kind, name)
	input := el.FindFirst(LikeInput)
	input.Attributes["minlength"] = fmt.Sprintf("%d", minLen)
	input.Attributes["maxlength"] = fmt.Sprintf("%d", maxLen)
	return el
}

// NewInputTextarea is a div element with labeled textarea child
func NewInputTextarea(name string, minLen, maxLen, col, row uint) HTMLInput {
	el := newTextArea(name, col, row)
	input := el.FindFirst(LikeInput)
	input.Attributes["minlength"] = fmt.Sprintf("%d", minLen)
	input.Attributes["maxlength"] = fmt.Sprintf("%d", maxLen)
	return el
}

// NewInputNumber is a div element with labeled number child
func NewInputNumber(name string, min, max float32) HTMLInput {
	el := newInput(InputKindNumber, name)
	input := el.FindFirst(LikeInput)
	input.Attributes["min"] = fmt.Sprintf("%f", min)
	input.Attributes["max"] = fmt.Sprintf("%f", max)
	return el
}

// NewInputColor is a div element with labeled color child
func NewInputColor(name string) HTMLInput {
	el := newInput(InputKindColor, name)
	return el
}

// NewInputHidden gives a hidden input element
func NewInputHidden(name string, value string) HTMLInput {
	el := newInput(InputKindHidden, name)
	el.Attributes["value"] = value
	return el
}

// NewInputFile is a div element with labeled file child
func NewInputFile(name string) HTMLInput {
	el := newInput(InputKindFile, name)
	return el
}

// NewInputDate is a div element with labeled date child
func NewInputDate(name string, min, max *time.Time) HTMLInput {
	el := newInput(InputKindDate, name)
	input := el.FindFirst(LikeInput)
	input.Attributes["min"] = PrintDate(min)
	input.Attributes["max"] = PrintDate(max)
	return el
}

// NewInputTime is a div element with labeled time child
func NewInputTime(name string, min, max *time.Time) HTMLInput {
	el := newInput(InputKindTime, name)
	input := el.FindFirst(LikeInput)
	input.Attributes["min"] = PrintTime(min)
	input.Attributes["max"] = PrintTime(max)
	return el
}

// NewInputDateTime is a div element with labeled datetime-local child
func NewInputDateTime(name string, min, max *time.Time) HTMLInput {
	el := newInput(InputKindDateTime, name)
	input := el.FindFirst(LikeInput)
	input.Attributes["min"] = PrintDateTime(min)
	input.Attributes["max"] = PrintDateTime(max)
	return el
}

// NewInputSelect is a div element with labeled select child
// One to many selections allowed.
func NewInputSelect[T fmt.Stringer](name string, options []T) HTMLInput {
	el := newSelect(name, options)
	return el
}

// NewInputRadio is a div element with labeled radio input children
// One selection allowed.
func NewInputRadio[T fmt.Stringer](name string, options []T) HTMLInput {
	el := newRadioCheckbox(InputKindRadio, name, options)
	return el
}

// NewInputCheckbox is a div element with labeled checkbox input children.
// Many selections allowed.
func NewInputCheckbox[T fmt.Stringer](name string, options []T) HTMLInput {
	el := newRadioCheckbox(InputKindCheckbox, name, options)
	return el
}

// SetFirstValue finds the first input element and sets its
// value. Recommended use is with elements that utilize the
// value attribute. Usage with a select, radio, or checkbox
// will just set the first element value
func (el *HTMLInput) SetFirstValue(s string) error {
	return el.SetNthValue(1, s)
}

// SetNthValue finds the nth occurrence on an input element.
// Depending on that element it will set the value accordingly.
// Not for use with select as the option children are treated
// differently, see SetSelectOption.
func (el *HTMLInput) SetNthValue(occurrence uint, s string) error {
	var c uint
	input := el.findNth(&c, occurrence, LikeInput)
	if input == nil {
		return ErrNoInputElement
	}
	if input.Tag == InputKindTextarea {
		input.InnerText = s
		return nil
	}
	if input.Attributes["type"] == InputKindCheckbox ||
		input.Attributes["type"] == InputKindRadio {
		input.EnsureAttributes()
		input.Attributes["checked"] = strconv.FormatBool(s == "t" || s == "true")
		return nil
	}
	input.EnsureAttributes()
	input.Attributes["value"] = s
	return nil
}

// SetSelectOption modifies a select option's selected attribute
func (el *HTMLInput) SetSelectOption(index uint, b bool) error {
	var c uint
	input := el.findNth(&c, 1, LikeTag("select"))
	if input == nil {
		return ErrNoInputElement
	}
	if len(input.Children) < int(index) {
		return ErrNoInputElement
	}
	option := input.Children[index]
	option.EnsureAttributes()
	option.Attributes["selected"] = strconv.FormatBool(b)
	return nil
}

// InputSelectedValue parses the chosen items of a select, checkbox, or radio
// from inputs names the same as the first input element
func InputSelectedValue[T fmt.Stringer](el *HTMLInput, pool []T) ([]T, error) {
	indexes, err := el.ParseIndexes()
	if err != nil {
		return nil, err
	}
	consumed := make(map[int]bool, len(indexes))
	items := make([]T, len(indexes))
	for i, index := range indexes {
		if index < 0 || index > len(pool) {
			return nil, ErrInvalidSelection
		}
		if _, exists := consumed[i]; exists {
			return nil, ErrDuplicateSelection
		}
		consumed[i] = true
		items[i] = pool[index]
	}
	return items, nil
}

func (el *HTMLInput) formKey() (ParsedForm, string, error) {
	firstInput := el.FindFirst(LikeInput)
	if firstInput.Attributes == nil {
		return nil, "", ErrAttributesNil
	}
	name, exists := firstInput.Attributes["name"]
	if !exists {
		return nil, "", ErrAttrNotExist
	}

	form := el.Form()
	return form, name, nil
}

// ParseString provides the string value of an input
func (el *HTMLInput) ParseString() (string, error) {
	form, key, err := el.formKey()
	if err != nil {
		return "", err
	}
	s := form[key]
	return s, nil
}

// ParseBool provides the boolean value of an input
func (el *HTMLInput) ParseBool() (bool, error) {
	form, key, err := el.formKey()
	if err != nil {
		return false, err
	}
	s := form[key]
	return strconv.ParseBool(s)
}

// ParseFloat provides the numerical value of an input
func (el *HTMLInput) ParseFloat() (float64, error) {
	form, key, err := el.formKey()
	if err != nil {
		return 0.0, err
	}
	s := form[key]
	return strconv.ParseFloat(s, 64)
}

// ParseInt provides the numerical value of an input
func (el *HTMLInput) ParseInt() (int, error) {
	form, key, err := el.formKey()
	if err != nil {
		return 0, err
	}
	s := form[key]
	v, err := strconv.Atoi(s)
	return int(v), err
}

// ParseTime provides the date, time, or datetime value of an input
func (el *HTMLInput) ParseTime() (*time.Time, error) {
	form, key, err := el.formKey()
	if err != nil {
		return nil, err
	}
	s := form[key]

	// there are 3 ays I can parse out time for HTML,
	// so we look over the 3 html date, time, and datetime formats
	// until one works and we can set the value, otherwise fail
	var parsed time.Time
	layouts := []string{TimeFormatHTMLDate, TimeFormatHTMLTime, TimeFormatHTMLDateTime}
	for _, layout := range layouts {
		parsed, err = time.Parse(layout, s)
		if err != nil {
			// we expect err since it may not be this format
			continue
		}
		break
	}
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

// ParseIndexes provides the indexes of selected items
func (el *HTMLInput) ParseIndexes() ([]int, error) {
	form, key, err := el.formKey()
	if err != nil {
		return nil, err
	}
	s := form[key]
	parts := strings.Split(s, ",")
	result := make([]int, 0, len(parts))
	for _, s := range parts {
		s := strings.TrimSpace(s)
		v, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return nil, err
		}
		result = append(result, int(v))
	}
	return result, nil
}
