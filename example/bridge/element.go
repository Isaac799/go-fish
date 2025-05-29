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

// InputChild gives the first child with the input related tag
func (el *HTMLElement) InputChild() *HTMLElement {
	inputTags := []string{"input", "select", "textarea"}

	if el.Children == nil {
		return nil
	}
	for _, c := range el.Children {
		for _, tag := range inputTags {
			if c.Tag == tag {
				return &c
			}
		}
	}
	return nil
}
