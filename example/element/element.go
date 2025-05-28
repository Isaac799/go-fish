// Package element is an simple abstract on HTML element
package element

// AttributeKey is a is the key in an html element kye-value pair
type AttributeKey string

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
)

const (
	// True is a string alias for TRUE
	True = "true"
	// False is a string alias for FALSE
	False = "false"
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
	Attributes map[AttributeKey]string
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
