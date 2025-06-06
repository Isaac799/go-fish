package bridge

import "slices"

// ElementLike helps narrow down filtering for a specific element
type ElementLike = func(el *HTMLElement) bool

func (el *HTMLElement) findNth(count *uint, occurrence uint, filters ...ElementLike) *HTMLElement {
	if el.Children != nil {
		for i := range el.Children {
			d := el.Children[i].findNth(count, occurrence, filters...)
			if d == nil {
				continue
			}
			return d
		}
	}

	for _, fn := range filters {
		pass := fn(el)
		if pass {
			continue
		}
		return nil
	}

	if *count != occurrence {
		*count++
		return nil
	}

	return el
}

func (el *HTMLElement) findAll(into *[]*HTMLElement, filters ...ElementLike) {
	if el.Children != nil {
		for i := range el.Children {
			el.Children[i].findAll(into, filters...)
		}
	}

	for _, fn := range filters {
		pass := fn(el)
		if pass {
			continue
		}
		return
	}

	*into = append(*into, el)
}

// LikeInput is a function to see if an elements tag is similar to an
// input that may be used in a form
var LikeInput = ElementLike(func(element *HTMLElement) bool {
	inputTags := []string{"input", "select", "textarea"}
	if !slices.Contains(inputTags, element.Tag) {
		return false
	}
	return true
})

// LikeTag is a small alias to see elements like a certain tag
var LikeTag = func(tag string) ElementLike {
	return ElementLike(func(element *HTMLElement) bool {
		return element.Tag == tag
	})
}

// HasAttribute is a small alias to see elements where an attribute is
// non empty value
var HasAttribute = func(key AttributeKey, value string) ElementLike {
	return ElementLike(func(element *HTMLElement) bool {
		if element.Attributes == nil {
			return false
		}
		v, exists := element.Attributes[key]
		if !exists {
			return false
		}
		return v == value
	})
}

// HasName is a small alias to find an element with a name
// non empty value
var HasName = func(name string) ElementLike {
	return ElementLike(func(element *HTMLElement) bool {
		if element.Attributes == nil {
			return false
		}
		v, exists := element.Attributes["name"]
		if !exists {
			return false
		}
		return v == name
	})
}

// FindFirst provides a consumer way to search for the first element matching
// a criteria
func (el *HTMLElement) FindFirst(filters ...ElementLike) *HTMLElement {
	var c uint = 1
	return el.findNth(&c, 1, filters...)
}

// FindNth provides a consumer way to search for the nth element matching
// a criteria
func (el *HTMLElement) FindNth(occurrence uint, filters ...ElementLike) *HTMLElement {
	var c uint = 1
	return el.findNth(&c, occurrence, filters...)
}

// FindAll provides a consumer way to search for the all element matching
// a criteria
func (el *HTMLElement) FindAll(filters ...ElementLike) []*HTMLElement {
	items := make([]*HTMLElement, 0)
	el.findAll(&items, filters...)
	return items
}
