package bridge

import (
	"errors"
	"maps"
	"strconv"
	"time"
)

var (
	// ErrNoInputElement is given when using a fn on an HTML element
	// such as parsing a number from an input, but it cannot find a field
	// like an input within that element
	ErrNoInputElement = errors.New("no input element found within this element")
	// ErrNoValueOnInputElement is given when using a fn on an HTML element
	// such as parsing a number from an input, but it cannot find a the
	// value attribute on the first input element
	ErrNoValueOnInputElement = errors.New("no value attribute found on input element")
)

// HTMLElement is key-value pairs that make up an element
type HTMLElement struct {
	// Tag is the html tag
	Tag string
	// SelfClosing adds a closing tag
	SelfClosing bool
	// InnerText is the text inside an element, only rendered if SelfClosing is false
	InnerText string
	// Attributes make up an element
	Attributes Attributes
	// Children are rendered either inside (!SelfClosing) or below (SelfClosing) the parent
	Children []HTMLElement
}

// NewHTMLElement provides a new html element
func NewHTMLElement(tag string) HTMLElement {
	return HTMLElement{
		Tag:        tag,
		Attributes: make(map[AttributeKey]string, 0),
		Children:   make([]HTMLElement, 0),
	}
}

// GiveAttributes allows safely dumping several attributes at once.
// Useful when wanting to alias a group of attributes (say htmlx attributes)
// making the code a little more clear about why we are modifying the element.
func (el *HTMLElement) GiveAttributes(attrs map[AttributeKey]string) {
	el.EnsureAttributes()
	maps.Copy(el.Attributes, attrs)
}

// EnsureAttributes ensures attributes are not nil before usage
func (el *HTMLElement) EnsureAttributes() {
	if el.Attributes != nil {
		return
	}
	el.Attributes = make(Attributes)
}

// SetFirstValue finds the first input element and sets its
// value. Recommended use is with elements that utilize the
// value attribute. Usage with a select, radio, or checkbox
// will just set the first element value
func (el *HTMLElement) SetFirstValue(s string) error {
	return el.SetNthValue(1, s)
}

// SetNthValue finds the nth occurrence on an input element.
// Depending on that element it will set the value accordingly.
// Not for use with select as the option children are treated
// differently, see SetSelectOption.
func (el *HTMLElement) SetNthValue(occurrence uint, s string) error {
	var c uint
	input := el.findNth(&c, occurrence, LikeInput)
	if input == nil {
		return ErrNoInputElement
	}
	if input.Tag == string(InputKindTextarea) {
		input.InnerText = s
		return nil
	}
	if input.Attributes["type"] == string(InputKindCheckbox) ||
		input.Attributes["type"] == string(InputKindRadio) {
		input.EnsureAttributes()
		input.Attributes["checked"] = strconv.FormatBool(s == "t" || s == "true")
		return nil
	}
	input.EnsureAttributes()
	input.Attributes["value"] = s
	return nil
}

// SetSelectOption modifies a select option's selected attribute
func (el *HTMLElement) SetSelectOption(index uint, b bool) error {
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

// ValueString parses a string from the first input element found
func (el *HTMLElement) ValueString() (string, error) {
	input := el.FindFirst(LikeInput)
	if input.Attributes == nil {
		return "", ErrNoInputElement
	}
	if input.Tag == string(InputKindTextarea) {
		return input.InnerText, nil
	}
	s, exists := input.Attributes["value"]
	if !exists {
		return "", ErrNoValueOnInputElement
	}
	return s, nil
}

// ValueTime parses a time from the first input element found
func (el *HTMLElement) ValueTime() (*time.Time, error) {
	s, err := el.ValueString()
	if err != nil {
		return nil, ErrNoValueOnInputElement
	}
	t, err := time.Parse(TimeFormatHTMLTime, s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// ValueDate parses a date from the first input element found
func (el *HTMLElement) ValueDate() (*time.Time, error) {
	s, err := el.ValueString()
	if err != nil {
		return nil, ErrNoValueOnInputElement
	}
	t, err := time.Parse(TimeFormatHTMLDate, s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// ValueDateTime parses a date from the first input element found
func (el *HTMLElement) ValueDateTime() (*time.Time, error) {
	s, err := el.ValueString()
	if err != nil {
		return nil, ErrNoValueOnInputElement
	}
	t, err := time.Parse(TimeFormatHTMLDateTime, s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// ValueInt parses a number from the first input element found
func (el *HTMLElement) ValueInt() (int, error) {
	s, err := el.ValueString()
	if err != nil {
		return 0, ErrNoValueOnInputElement
	}
	parsed, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return parsed, nil
}

// ValueUint64 parses a number from the first input element found
// useful for things like a table page where the number must be unsigned
func (el *HTMLElement) ValueUint64() (uint64, error) {
	s, err := el.ValueString()
	if err != nil {
		return 0, ErrNoValueOnInputElement
	}
	parsed, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		return 0, err
	}
	return parsed, nil
}

// ValueIndexes parses the chosen indexes of a select, checkbox, or radio
// from inputs names the same as the first input element
func (el *HTMLElement) ValueIndexes() ([]int, error) {
	firstEl := el.FindFirst(LikeInput)
	if firstEl.Attributes == nil {
		return nil, ErrAttributesNil
	}
	name, exists := firstEl.Attributes["name"]
	if !exists {
		return nil, ErrKeyDoesNotExist
	}
	parsed := el.Form()
	indexes, err := parsed.Indexes(name)
	return indexes, err
}

// ValueElementSelected parses the chosen items of a select, checkbox, or radio
// from inputs names the same as the first input element
func ValueElementSelected[T Printable](el *HTMLElement, pool []T) ([]T, error) {
	indexes, err := el.ValueIndexes()
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
