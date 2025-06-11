package bridge

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	reflectionStructTag = "gofish"
)

var (
	// ErrCannotSet is given if I cannot access a field
	// most likely because it was not exported, the field is not pointed to
	// or I forgot to use Elem.
	ErrCannotSet = errors.New("cannot set unexported field")
	// ErrNotSupported is given if that type is not supported
	// since I do not have full support for all types
	ErrNotSupported = errors.New("conversion is not supported")
	// ErrMalformedTime is given if a date, time, or datetime cannot
	// be parsed given one of the 3 formats
	ErrMalformedTime = errors.New("date time represented as string did not follow an html format")
)

// HTMLValueType are expected types I can parse from an attribute.
//
// Seems like a wild one, but I wanted to be explicit about
// the limitations of my reflection.
//
// Align with possible HTML values represented as strings.
type HTMLValueType interface {
	~string | ~bool | ~[]string | ~[]bool |
		~int | ~uint | ~[]int | ~[]uint |
		~int8 | ~uint8 | ~[]int8 | ~[]uint8 |
		~int16 | ~uint16 | ~[]int16 | ~[]uint16 |
		~int32 | ~uint32 | ~[]int32 | ~[]uint32 |
		~int64 | ~uint64 | ~[]int64 | ~[]uint64 |
		~float32 | ~float64 | ~[]float32 | ~[]float64 |
		time.Time | []time.Time
}

// AttributesToStruct will parse out a string map into something that
// is easier to work with. Max struct depth is 1 (so no nested structs).
// Thanks to this blog for getting me started https://go.dev/blog/laws-of-reflection
func AttributesToStruct(m Attributes, v any) error {
	// value
	to := reflect.ValueOf(v).Elem()
	// concrete type
	toType := to.Type()

	// while we are iterating over the number of fields in the _value_,
	// its order and len are same as the _type_ metadata. So we can access
	// them both at the same time
	for i := range to.NumField() {
		// value
		toField := to.Field(i)
		// concrete type
		toFieldType := toType.Field(i)

		toFieldTag := toFieldType.Tag.Get(reflectionStructTag)
		if toFieldTag == "" {
			continue
		}

		for key, str := range m {
			// align map keys with struct tag
			if key != toFieldTag {
				continue
			}

			// put the map string into the value
			err := htmlStrIntoReflectVal(toField, str)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// htmlStrIntoReflectVal allows parsing a string into a field.
// Designed to be used with ParsedForm, mainly for string to be
// formatted as expected (mainly slices to be comma separated).
//
// Also mainly used to parse HTML attributes, not meant for slices
// with that use case.
//
// Also having this broken out enables me to re-use it for reflecting
// into custom struct (e.g. validation)
func htmlStrIntoReflectVal(toField reflect.Value, str string) error {
	// ensures we can set a element
	canSetToField := toField.CanSet()
	if !canSetToField {
		return ErrCannotSet
	}

	// Based on what we are reflecting _to_
	switch toField.Kind() {
	case reflect.String:
		// strings no change
		toField.SetString(str)

	case reflect.Bool:
		// bool simple enough, strconv to the rescue
		v, err := strconv.ParseBool(str)
		if err != nil {
			return err
		}
		toField.SetBool(v)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// int simple too thanks to strconv
		v, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return err
		}
		toField.SetInt(v)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// uint just strconv
		v, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return err
		}
		toField.SetUint(v)

	case reflect.Float32, reflect.Float64:
		// uint just strconv
		v, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return err
		}
		toField.SetFloat(v)

	case reflect.Struct:
		// here is where things get a little different, not bad tho

		// I only support reflecting to time, not custom structs
		if toField.Type() != reflect.TypeOf(time.Time{}) {
			return ErrNotSupported
		}

		// there are 3 ays I can parse out time for HTML,
		// so we look over the 3 html date, time, and datetime formats
		// until one works and we can set the value, otherwise fail
		var parsed time.Time
		var err error

		layouts := []string{TimeFormatHTMLDate, TimeFormatHTMLTime, TimeFormatHTMLDateTime}
		for _, layout := range layouts {
			parsed, err = time.Parse(layout, str)
			if err != nil {
				// we expect err since it may not be this format
				continue
			}
			break
		}
		if err != nil {
			return ErrMalformedTime
		}
		toField.Set(reflect.ValueOf(parsed))

	case reflect.Slice:
		// slice is the most complex, in terms of reflection
		// but the simplest in terms of parsing, since to parse
		// I can just call what I already made over and over

		// gets the underlying type, so []int is int
		sliceType := toField.Type().Elem()

		// the usage of this function is meant to be in conjunction with
		// my ParsedForm struct, which will format arrays in a comma separated
		strs := strings.Split(str, ",")

		// make a newSlice of the desired type with cap to hold app parsed values
		// and I append to it
		newSlice := reflect.MakeSlice(toField.Type(), 0, len(strs))

		for _, s := range strs {
			// must remove white space because "true, false" -> []string{"true", " false"}
			// and the space breaks value parsing
			s := strings.TrimSpace(s)

			// new _value_ with same type as slice
			elemVal := reflect.New(sliceType).Elem()

			// just do what I already did
			err := htmlStrIntoReflectVal(elemVal, s)
			if err != nil {
				return err
			}

			// once the _value_ has been set we append it to parsed slice
			newSlice = reflect.Append(newSlice, elemVal)
		}
		toField.Set(newSlice)
	default:
		return ErrNotSupported
	}

	// tada!
	return nil
}
