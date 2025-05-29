package bridge

import (
	"mime/multipart"
	"net/http"
)

// flatten recursively grabs elements
func flatten(elements []HTMLElement) []*HTMLElement {
	collected := []*HTMLElement{}
	for _, el := range elements {
		if el.Children != nil {
			childKeys := flatten(el.Children)
			for _, c := range childKeys {
				collected = append(collected, c)
			}
		}
		collected = append(collected, &el)
	}
	return collected
}

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
func FormFromRequest(r *http.Request, form HTMLElement) map[string][]string {
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
