package bridge

import (
	"errors"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

var (
	// ErrKeyDoesNotExist is given if extracting a form value with a key,
	// but that key is not part of the form.
	ErrKeyDoesNotExist = errors.New("key not found in form")
	// ErrMultipleValues is given if extracting a form value with a key,
	// but that key's []string length is not as expected
	ErrMultipleValues = errors.New("form value has multiple values")
	// ErrInvalidSelection is given when discovering selected items from
	// a list of indexes
	ErrInvalidSelection = errors.New("form value has an invalid selection")
	// ErrDuplicateSelection is given when discovering selected items from
	// a list of indexes and an index is given more than onces
	ErrDuplicateSelection = errors.New("form value has the same selection more than once")
)

// ParsedForm is the result of comparing a request against a predefined form
// with helpful methods for parsing values.
type ParsedForm map[string][]string

// FormKeys looks at all the children (recursively) and grabs the name-value
// pairs. These attributes are expected for a form.
func (el *HTMLElement) FormKeys() []string {
	m := []string{}
	if el.Children == nil {
		return nil
	}
	flat := flatten(el.Children)
	for _, el := range flat {
		if el.Attributes == nil {
			continue
		}
		key := el.Attributes["name"]
		m = append(m, key)
	}
	return m
}

// FormFromRequest will compare a request form to that of
// an HTML element with inputs in its tree. The form element
// is only used to gather what keys are relevant to a form being
// submitted. The values are only from the form.
func FormFromRequest(r *http.Request, form HTMLElement) ParsedForm {
	m := form.FormKeys()
	m2 := make(map[string][]string)
	r.ParseForm()
	for _, key := range m {
		if !r.Form.Has(key) {
			continue
		}
		m2[key] = r.Form[key]
	}

	return m2
}

// String parses a string from a form
func (form ParsedForm) String(key string) (string, error) {
	s, exists := form[key]
	if !exists {
		return "", ErrKeyDoesNotExist
	}
	if len(s) != 1 {
		return "", ErrMultipleValues
	}
	return s[0], nil
}

// Time parses a time from a form
func (form ParsedForm) Time(key string) (*time.Time, error) {
	s, err := form.String(key)
	if err != nil {
		return nil, err
	}
	t, err := time.Parse(TimeFormatHTMLTime, s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// Date parses a date from a form
func (form ParsedForm) Date(key string) (*time.Time, error) {
	s, err := form.String(key)
	if err != nil {
		return nil, err
	}
	t, err := time.Parse(TimeFormatHTMLDate, s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// DateTime parses a date from a form
func (form ParsedForm) DateTime(key string) (*time.Time, error) {
	s, err := form.String(key)
	if err != nil {
		return nil, err
	}
	t, err := time.Parse(TimeFormatHTMLDateTime, s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// Number parses a number a form
func (form ParsedForm) Number(key string) (int, error) {
	s, err := form.String(key)
	if err != nil {
		return 0, err
	}
	parsed, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return parsed, nil
}

// Indexes parses any select, checkbox, or radio values into the indexes chosen
func (form ParsedForm) Indexes(key string) ([]int, error) {
	values, exists := form[key]
	if !exists {
		return nil, ErrKeyDoesNotExist
	}
	indexes := make([]int, len(values))
	for i, s := range values {
		parsed, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		indexes[i] = parsed
	}
	return indexes, nil
}

// FormSelected gives a slice of what values where selected
// similar to how we created a select, radio, or checkbox it
// can tell which of the items was chosen.
func FormSelected[T Printable](form ParsedForm, key string, pool []T) ([]T, error) {
	indexes, err := form.Indexes(key)
	if err != nil {
		return nil, ErrKeyDoesNotExist
	}
	consumed := make(map[int]bool, len(indexes))
	items := make([]T, len(indexes))
	for i, index := range indexes {
		if index < 0 || index > len(indexes) {
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

// FormFilesFromRequest will compare a request form files to that of
// an HTML element with inputs in its tree
func FormFilesFromRequest(r *http.Request, form HTMLElement) (map[string]multipart.File, error) {
	m := form.FormKeys()
	m2 := make(map[string]multipart.File)
	for _, key := range m {
		// FormFile calls required parsing
		file, _, err := r.FormFile(key)
		if err != nil {
			return nil, err
		}
		m2[key] = file
	}

	return m2, nil
}
