package bridge

func (el *HTMLElement) deleteNth(count *uint, occurrence uint, filters ...ElementLike) {
	if el.Children == nil {
		return
	}

	kept := make([]HTMLElement, 0, len(el.Children))

	for _, child := range el.Children {
		isMatch := true
		for _, fn := range filters {
			pass := fn(&child)
			if !pass {
				isMatch = false
				break
			}
		}
		if !isMatch {
			kept = append(kept, child)
			continue
		}
		*count++
		if *count != occurrence {
			kept = append(kept, child)
			continue
		}
		continue
	}
	el.Children = kept

	if *count == occurrence {
		return
	}

	for i := range el.Children {
		el.Children[i].deleteNth(count, occurrence, filters...)
		if *count == occurrence {
			return
		}
	}
}

func (el *HTMLElement) deleteAll(filters ...ElementLike) {
	if el.Children == nil {
		return
	}

	kept := make([]HTMLElement, 0, len(el.Children))

	for _, child := range el.Children {
		isMatch := true
		for _, fn := range filters {
			pass := fn(&child)
			if !pass {
				isMatch = false
				break
			}
		}
		if !isMatch {
			kept = append(kept, child)
			continue
		}
		continue
	}
	el.Children = kept

	for i := range el.Children {
		el.Children[i].deleteAll(filters...)
	}
}

// DeleteFirst provides a consumer way to remove the first child
// element matching a criteria from the element tree
func (el *HTMLElement) DeleteFirst(filters ...ElementLike) {
	var c uint
	el.deleteNth(&c, 1, filters...)
}

// DeleteNth provides a consumer way to remove the nth child
// element matching a criteria from the element tree
func (el *HTMLElement) DeleteNth(occurrence uint, filters ...ElementLike) {
	var c uint
	el.deleteNth(&c, occurrence, filters...)
}

// DeleteAll provides a consumer way to remove the all children
// matching a criteria from the element tree
func (el *HTMLElement) DeleteAll(filters ...ElementLike) {
	el.deleteAll(filters...)
}
