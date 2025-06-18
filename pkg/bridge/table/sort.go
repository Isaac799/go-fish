package table

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/Isaac799/go-fish/pkg/bridge"
)

// random as to not overlap with consumer field names
var (
	// prefixes are used with the column index to make form keys
	formKeyPrefixSortBy = bridge.RandomID()
)

// Sort provides the Define and Modify functions needed
// enable sorting for columns of a table
func Sort(columnIndexes ...int) (Mod, Mod) {
	var defineFn Mod = func(table *HTMLTable) error {
		headers := table.El.FindAll(bridge.LikeTag("th"))
		if len(headers) == 0 {
			return ErrMissingExpectedElement
		}

		if table.sortInputs == nil {
			table.sortInputs = make(map[int]*bridge.HTMLInput, len(headers))
		}

		for i := range headers {
			if !slices.Contains(columnIndexes, i) {
				continue
			}
			// So we keep the sort of items even if not clicked
			sortKey := fmt.Sprintf("%s%d", formKeyPrefixSortBy, i)
			hiddenSortInput := bridge.NewInputHidden(sortKey, "0")
			headers[i].Children = append(headers[i].Children, hiddenSortInput)

			table.sortInputs[i] = &hiddenSortInput
		}

		return nil
	}

	var modifyFn Mod = func(table *HTMLTable) error {
		headers := table.El.FindAll(bridge.LikeTag("th"))
		if len(headers) == 0 {
			return ErrMissingExpectedElement
		}

		for index := range table.sortInputs {
			// parse the form of the header so we can extract the value of the sort
			desiredSort, err := table.sortInputs[index].ParseInt()
			if err != nil {
				continue
			}

			previousSort := SortNone
			if slices.Contains(acceptableSort, desiredSort) {
				previousSort = desiredSort
			}

			// So we keep the sort of items even if not clicked
			sortKey := fmt.Sprintf("%s%d", formKeyPrefixSortBy, index)

			sortBtn := bridge.HTMLElement{
				Tag: "button",
				Attributes: bridge.Attributes{
					"class":     "material-icons",
					"type":      "submit",
					"hx-post":   table.conf.HxPost,
					"name":      sortKey,
					"value":     strconv.Itoa(nextSort[previousSort]),
					"hx-target": table.conf.HxSwapTarget,
					"hx-swap":   "outerHTML",
				},
				InnerText: icons[previousSort],
			}

			headers[index].Children = append(headers[index].Children, sortBtn)
		}

		return nil
	}
	return defineFn, modifyFn
}
