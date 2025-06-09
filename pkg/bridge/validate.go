package bridge

import (
	"errors"
	"regexp"
	"time"
)

var (
	// ErrCannotValidate is given when a input is not available to
	// validate due to a certain tag or tag.
	ErrCannotValidate = errors.New("cannot validate this kind of input")
	// ErrMinMaxRequired is given validating an input
	// that does not have at lease min/max
	ErrMinMaxRequired = errors.New("cannot validate without a min/max")
)

// Validate will determine look at the first input in an element
// and based on its type apply the appropriate validation
func (el *HTMLElement) Validate() (bool, error) {
	kind, kindExists := el.Attributes["type"]
	if !kindExists {
		return false, ErrNoInputElement
	}

	switch kind {
	case InputKindFile, InputKindColor:
		// todo
	case InputKindCheckbox, InputKindSelect, InputKindRadio:
		// todo
	case InputKindText, InputKindPassword, InputKindEmail, InputKindSearch, InputKindTel, InputKindURL, InputKindTextarea, InputKindHidden:
		v, err := ValidationForElement[string](el)
		if err != nil {
			return false, err
		}
		if v.Required && len(v.Value) == 0 {
			return false, nil
		} else if !v.Required && len(v.Value) == 0 {
			return true, nil
		}
		if len(v.Pattern) > 0 {
			matched, err := regexp.Match(v.Pattern, []byte(v.Value))
			if err != nil {
				return false, err
			}
			if !matched {
				return false, nil
			}
		}
		withinMinMax := len(v.Value) >= int(v.MinLen) && len(v.Value) < int(v.MaxLen)
		return withinMinMax, nil
	case InputKindNumber:
		v, err := ValidationForElement[float64](el)
		if err != nil {
			return false, err
		}
		if v.Required && v.Value == 0 {
			return false, nil
		} else if !v.Required && v.Value == 0 {
			return true, nil
		}
		withinMinMax := v.Value >= v.Min && v.Value < v.Max
		return withinMinMax, nil
	case InputKindDate, InputKindTime, InputKindDateTime:
		v, err := ValidationForElement[time.Time](el)
		if err != nil {
			return false, err
		}
		if v.Required && v.Value.IsZero() {
			return false, nil
		} else if !v.Required && v.Value.IsZero() {
			return true, nil
		}
		withinMinMax := v.Min.Before(v.Value) && v.Max.After(v.Value)
		return withinMinMax, nil
	}
	return false, ErrCannotValidate
}

// Validation is derived from attributes to determine if an
// input is valid. All fields are optional
type Validation[T HTMLValueType] struct {
	Kind     string `gofish:"type"`
	Value    T      `gofish:"value"`
	Min      T      `gofish:"min"`
	Max      T      `gofish:"max"`
	MinLen   int    `gofish:"minlength"`
	MaxLen   int    `gofish:"maxlength"`
	Required bool   `gofish:"required"`
	Pattern  string `gofish:"pattern"`
}

// attributeValue helps me just unmarshal the value attr
type attributeValue[T HTMLValueType] struct {
	Value T `HTMLAttr:"value"`
}

// ValidationForElement provides fields to validate against derived
// from the attribute map
func ValidationForElement[T HTMLValueType](el *HTMLElement) (*Validation[T], error) {
	v := Validation[T]{}
	if el.Attributes == nil {
		return &v, ErrAttributesNil
	}
	err := AttributesToStruct(el.Attributes, &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}
