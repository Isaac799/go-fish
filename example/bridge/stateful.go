package bridge

// Stateful wraps an element in a form and adds hidden html inputs
// to store key-value pairs. Enabling form submissions within an element
// to know the state it was rendered with. e.g. what page a table was on.
func Stateful(el *HTMLElement, state map[string]string) *HTMLElement {
	form := HTMLElement{
		Tag:      "form",
		Children: make([]HTMLElement, 0, len(state)+1),
	}
	form.Children = append(form.Children, *el)

	if state != nil {
		i := 1
		for k, v := range state {
			form.Children = append(form.Children, NewInputHidden(k, v))
			i++
		}
	}

	return &form
}
