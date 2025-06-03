package bridge

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

// EnsureAttributes ensures attributes are not nil before usage
func (el *HTMLElement) EnsureAttributes() {
	if el.Attributes != nil {
		return
	}
	el.Attributes = make(Attributes)
}

// SetValue finds the nth occurrence of an input searching nested
// and modifies the value attribute.
func (el *HTMLElement) SetValue(occurrence uint, s string) {
	var c uint = 1
	input := el.findNth(&c, occurrence)
	if input == nil {
		return
	}
	if input.Tag == "textarea" {
		input.InnerText = s
		return
	}
	input.EnsureAttributes()
	input.Attributes["value"] = s
	return
}

// SetChecked finds the nth occurrence of an input searching nested
// and modifies the checked attribute
func (el *HTMLElement) SetChecked(occurrence uint, b bool) {
	var c uint = 1
	input := el.findNth(&c, occurrence)
	if input == nil {
		return
	}
	input.EnsureAttributes()
	input.Attributes.SetChecked(b)
	return
}
