package bridge

import (
	"errors"
	"net/http"
	"slices"
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

	// used to help parse and populate elements and their values
	// for elements that do not store their value in a value tag
	nonValueTags = []string{"select", "textarea"}
	// used to help parse and populate elements and their values
	// for elements that do not store their value in a single value tag
	nonValueKinds = []string{"checkbox", "radio"}
)

// ParsedForm is the result of comparing a request against a predefined form
// with helpful methods for parsing values.
type ParsedForm map[string][]string

// String parses a string from the form given a key
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

// Time parses a time from the form given a key
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

// Date parses a date from the form given a key
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

// DateTime parses a date and time from the form given a key
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

// Int parses a int from the form given a key
func (form ParsedForm) Int(key string) (int, error) {
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

// Uint64 parses a uint64 from the form given a key
// useful for things like a table page where the number must be unsigned
func (el *HTMLElement) Uint64() (uint64, error) {
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

// Indexes parses the chosen indexes of a select, checkbox, or radio
// from the form given a key
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

// FormSelected parses the chosen items of a select, checkbox, or radio
// from the form given a key
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

// Form provides a parsed value of all the input related elements.
// Useful when defining what something should look like, then
// getting the form that will actually be provided to a client.
func (el *HTMLElement) Form() ParsedForm {
	m := make(map[string][]string)

	inputs := el.FindAll(LikeInput)

	// For selects we must look at their selected options
	for _, input := range inputs {
		if input.Attributes == nil {
			continue
		}
		key, exists := input.Attributes["name"]
		if !exists {
			continue
		}

		if input.Children == nil {
			continue
		}
		for selectOptionIndex, child := range input.Children {
			if child.Tag != "option" {
				continue
			}
			if child.Attributes == nil {
				continue
			}
			v, exists := child.Attributes["selected"]
			if !exists {
				continue
			}
			if v != "true" {
				continue
			}

			indexStr := strconv.Itoa(selectOptionIndex)

			if m[key] == nil {
				m[key] = make([]string, 0, 1)
				m[key] = append(m[key], indexStr)
			} else {
				m[key] = append(m[key], indexStr)
			}
		}
	}

	// For radios and checkbox we must look at their elements
	// and compare to the values given and add the checked attribute
	radioGroups := map[string][]*HTMLElement{}
	checkboxGroups := map[string][]*HTMLElement{}

	for _, input := range inputs {
		if input.Attributes == nil {
			continue
		}
		kind := input.Attributes["type"]
		if kind != "radio" && kind != "checkbox" {
			continue
		}
		if input.Attributes == nil {
			continue
		}
		key, exists := input.Attributes["name"]
		if !exists {
			continue
		}
		if kind == "radio" {
			if radioGroups[key] == nil {
				radioGroups[key] = make([]*HTMLElement, 1)
				radioGroups[key][0] = input
			} else {
				radioGroups[key] = append(radioGroups[key], input)
			}
		}
		if kind == "checkbox" {
			if checkboxGroups[key] == nil {
				checkboxGroups[key] = make([]*HTMLElement, 1)
				checkboxGroups[key][0] = input
			} else {
				checkboxGroups[key] = append(checkboxGroups[key], input)
			}
		}
	}

	for key, inputs := range radioGroups {
		for index, input := range inputs {
			if input.Attributes == nil {
				continue
			}
			v, exists := input.Attributes["checked"]
			if !exists {
				continue
			}
			if v != "true" {
				continue
			}
			indexStr := strconv.Itoa(index)
			if m[key] == nil {
				m[key] = make([]string, 0, 1)
				m[key] = append(m[key], indexStr)
			} else {
				m[key] = append(m[key], indexStr)
			}
		}
	}

	for key, inputs := range checkboxGroups {
		for index, input := range inputs {
			if input.Attributes == nil {
				continue
			}
			v, exists := input.Attributes["checked"]
			if !exists {
				continue
			}
			if v != "true" {
				continue
			}
			indexStr := strconv.Itoa(index)
			if m[key] == nil {
				m[key] = make([]string, 0, 1)
				m[key] = append(m[key], indexStr)
			} else {
				m[key] = append(m[key], indexStr)
			}
		}
	}

	// For textarea we must look at inner html
	for _, input := range inputs {
		if input.Tag != "textarea" {
			continue
		}
		key, exists := input.Attributes["name"]
		if !exists {
			continue
		}
		m[key] = make([]string, 0, 1)
		m[key] = append(m[key], input.InnerText)
	}

	// for most other inputs we can look at the value attribute
	for _, input := range inputs {
		if input.Attributes == nil {
			continue
		}
		if slices.Contains(nonValueTags, input.Tag) {
			continue
		}
		kind, exists := input.Attributes["type"]
		if !exists {
			continue
		}
		if slices.Contains(nonValueKinds, kind) {
			continue
		}
		key, exists := input.Attributes["name"]
		if !exists {
			continue
		}
		value, exists := input.Attributes["value"]
		if !exists {
			continue
		}
		m[key] = make([]string, 0, 1)
		m[key] = append(m[key], value)
	}

	return m
}

// FillForm will look at the form in a request and compare it
// to all the inputs in an element to set their attributes accordingly.
// Attributes: value, checked, and selected. Or inner text if needed.
// Useful when you know what a element is, and want to preserve state
// from a users form submission.
func (el *HTMLElement) FillForm(r *http.Request) {
	r.ParseForm()
	inputs := el.FindAll(LikeInput)

	// For selects we must look at their selected options
	// by comparing the options to the values given
	// and add the selected attribute to the options
	for rKey := range r.Form {
		for _, input := range inputs {
			if input.Tag != "select" {
				continue
			}
			if input.Attributes == nil {
				continue
			}
			key, exists := input.Attributes["name"]
			if !exists {
				continue
			}
			if rKey != key {
				continue
			}
			if input.Children == nil {
				continue
			}
			for _, selectedIndexStr := range r.Form[rKey] {
				selectedIndex, err := strconv.Atoi(selectedIndexStr)
				if err != nil {
					continue
				}
				for selectOptionIndex, child := range input.Children {
					if child.Tag != "option" {
						continue
					}
					child.EnsureAttributes()
					if selectedIndex != selectOptionIndex {
						continue
					}
					child.Attributes["selected"] = "true"
				}
			}
		}
	}

	// For radios and checkbox we must look at their elements
	// and compare to the values given and add the checked attribute
	radioGroups := map[string][]*HTMLElement{}
	checkboxGroups := map[string][]*HTMLElement{}

	for _, input := range inputs {
		if input.Attributes == nil {
			continue
		}
		kind := input.Attributes["type"]
		if kind != "radio" && kind != "checkbox" {
			continue
		}
		if input.Attributes == nil {
			continue
		}
		key, exists := input.Attributes["name"]
		if !exists {
			continue
		}
		if kind == "radio" {
			if radioGroups[key] == nil {
				radioGroups[key] = make([]*HTMLElement, 0, 1)
				radioGroups[key] = append(radioGroups[key], input)
			} else {
				radioGroups[key] = append(radioGroups[key], input)
			}
			continue
		}
		if radioGroups[key] == nil {
			radioGroups[key] = make([]*HTMLElement, 0, 1)
			radioGroups[key] = append(radioGroups[key], input)
		} else {
			radioGroups[key] = append(radioGroups[key], input)
		}
	}

	for rKey := range r.Form {
		for _, selectedIndexStr := range r.Form[rKey] {
			selectedIndex, err := strconv.Atoi(selectedIndexStr)
			if err != nil {
				continue
			}
			for key, inputs := range radioGroups {
				if rKey != key {
					continue
				}
				for checkedOptionIndex, input := range inputs {
					if input.Attributes == nil {
						continue
					}
					input.EnsureAttributes()
					if selectedIndex != checkedOptionIndex {
						continue
					}
					input.Attributes["checked"] = "true"
				}
			}
			for key, inputs := range checkboxGroups {
				if rKey != key {
					continue
				}
				for checkedOptionIndex, input := range inputs {
					if input.Attributes == nil {
						continue
					}
					input.EnsureAttributes()
					if selectedIndex != checkedOptionIndex {
						continue
					}
					input.Attributes["checked"] = "true"
				}
			}
		}
	}

	// For textarea we must look at inner html
	for rKey := range r.Form {
		for _, input := range inputs {
			if input.Tag != "textarea" {
				continue
			}
			key, exists := input.Attributes["name"]
			if !exists {
				continue
			}
			if rKey != key {
				continue
			}
			v := r.Form[rKey]
			if len(v) == 0 {
				continue
			}
			input.InnerText = v[len(v)-1]
		}
	}

	// Next we can look at non select elements
	for rKey := range r.Form {
		for _, input := range inputs {
			if input.Attributes == nil {
				continue
			}
			if slices.Contains(nonValueTags, input.Tag) {
				continue
			}
			kind, exists := input.Attributes["type"]
			if !exists {
				continue
			}
			if slices.Contains(nonValueKinds, kind) {
				continue
			}

			key, exists := input.Attributes["name"]
			if !exists {
				continue
			}
			if rKey != key {
				continue
			}
			_, exists = input.Attributes["value"]
			if !exists {
				continue
			}
			// I elected to take the last of a form value since the query
			// params are parsed after the post body. this ensures if I
			// make a post with a 'overwritten' value it will appear later
			// in the slice.
			// e.g. sort on a table has a hidden input to preserve
			// sort state for columns that are sorted. There is also a button that
			// when clicked will submit the 'overwritten' value in the query params
			// but just for that column
			newVal := r.Form[rKey][len(r.Form[rKey])-1]
			input.Attributes["value"] = newVal
		}
	}
}
