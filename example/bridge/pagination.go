package bridge

import "strconv"

// Pagination has info for pagination of a table
type Pagination struct {
	CurrentPage  uint64
	TotalPages   uint64
	TotalCount   uint64
	Limit        uint64
	PreviousPage uint64
	NextPage     uint64
	Offset       uint64
	ShowFirst    bool
	ShowPrevious bool
	ShowNext     bool
	ShowLast     bool
}

// NewPagination provides a pagination metadata given total records
func NewPagination(limit, page, totalCount uint64) Pagination {
	totalPages := (totalCount + limit - 1) / limit
	if page > totalPages {
		page = totalPages
	}
	previousPage := page - 1
	nextPage := page + 1
	return Pagination{
		CurrentPage:  page,
		TotalPages:   totalPages,
		TotalCount:   totalCount,
		Limit:        limit,
		PreviousPage: page - 1,
		NextPage:     nextPage,
		Offset:       (page - 1) * limit,
		ShowFirst:    page > 2,
		ShowPrevious: previousPage > 0,
		ShowNext:     nextPage <= totalPages,
		ShowLast:     totalPages > 1 && page < totalPages-1,
	}
}

// Element provides a div element based around HTMLX to be used
// in a form. Uses my classes and material icons
func (p *Pagination) Element(hxTarget, formPageKey, hxPost string) HTMLElement {
	paginationDiv := HTMLElement{
		Tag: "div",
		Attributes: Attributes{
			"class": "fr fr-center g2",
		},
		Children: make([]HTMLElement, 0, 7),
	}

	spreadIcon := HTMLElement{
		Tag: "span",
		Attributes: Attributes{
			"class": "material-icons",
		},
		InnerText: "more_horiz",
	}

	if p.ShowFirst {
		pageFirst := HTMLElement{
			Tag: "button",
			Attributes: Attributes{
				"type":      "submit",
				"hx-post":   hxPost,
				"name":      formPageKey,
				"value":     strconv.FormatUint(1, 10),
				"hx-target": hxTarget,
				"hx-swap":   "outerHTML",
			},
			InnerText: strconv.FormatUint(1, 10),
		}
		paginationDiv.Children = append(paginationDiv.Children, pageFirst, spreadIcon)
	}

	if p.ShowPrevious {
		pagePrev := HTMLElement{
			Tag: "button",
			Attributes: Attributes{
				"type":      "submit",
				"hx-post":   hxPost,
				"name":      formPageKey,
				"value":     strconv.FormatUint(p.PreviousPage, 10),
				"hx-target": hxTarget,
				"hx-swap":   "outerHTML",
			},
			InnerText: strconv.FormatUint(p.PreviousPage, 10),
		}
		paginationDiv.Children = append(paginationDiv.Children, pagePrev)
	}

	pageCurr := HTMLElement{
		Tag: "button",
		Attributes: Attributes{
			"class":           "is-current-page",
			"hx-disabled-elt": "this",
		},
		InnerText: strconv.FormatUint(p.CurrentPage, 10),
	}
	paginationDiv.Children = append(paginationDiv.Children, pageCurr)

	if p.ShowNext {
		pageNext := HTMLElement{
			Tag: "button",
			Attributes: Attributes{
				"type":      "submit",
				"hx-post":   hxPost,
				"name":      formPageKey,
				"value":     strconv.FormatUint(p.NextPage, 10),
				"hx-target": hxTarget,
				"hx-swap":   "outerHTML",
			},
			InnerText: strconv.FormatUint(p.NextPage, 10),
		}
		paginationDiv.Children = append(paginationDiv.Children, pageNext)
	}

	if p.ShowLast {
		pageLast := HTMLElement{
			Tag: "button",
			Attributes: Attributes{
				"type":      "submit",
				"hx-post":   hxPost,
				"name":      formPageKey,
				"value":     strconv.FormatUint(p.TotalPages, 10),
				"hx-target": hxTarget,
				"hx-swap":   "outerHTML",
			},
			InnerText: strconv.FormatUint(p.TotalPages, 10),
		}
		paginationDiv.Children = append(paginationDiv.Children, spreadIcon, pageLast)
	}

	// currentPage := NewInputHidden(formPageKey, strconv.FormatUint(p, 10CurrentPage)))
	// currentLimit := NewInputHidden(formLimitKey, strconv.FormatUint(p, 10Limit)))

	// paginationDiv.Children[7] = currentPage
	// paginationDiv.Children[8] = currentLimit

	return paginationDiv
}
