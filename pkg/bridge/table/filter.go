package table

import (
	"fmt"
	"slices"

	"github.com/Isaac799/go-fish/pkg/bridge"
)

// random as to not overlap with consumer field names
var (
	// prefixes are used with the column i to make form keys
	formKeyPrefixFilterBy = bridge.RandomID()
)

// Filter provides the Define fn to add filters on a table.
// Filters are text boxes. It is up to the server to parse its value.
// No need to modify after form filled as it is a static input.
func Filter(columnIndexes ...int) Mod {
	var defineFn Mod = func(table *HTMLTable) error {
		headers := table.El.FindAll(bridge.LikeTag("th"))
		if len(headers) == 0 {
			return ErrMissingExpectedElement
		}

		textBoxHTMLX := map[string]string{
			"hx-trigger":  "keyup changed delay:500ms",
			"hx-post":     table.conf.HxPost,
			"hx-target":   table.conf.HxSwapTarget,
			"hx-swap":     "outerHTML",
			"placeholder": "filter...",
		}

		if table.filterInputs == nil {
			table.filterInputs = make(map[int]*bridge.HTMLInput, len(headers))
		}

		for i := range headers {
			if !slices.Contains(columnIndexes, i) {
				continue
			}
			// Gives a text input to filter by. We can define these now
			// since they are not modified later and not using the 'hidden' flow.
			filterKey := fmt.Sprintf("%s%d", formKeyPrefixFilterBy, i)
			filterByDiv := bridge.NewInputText(filterKey, bridge.InputKindText, 0, 20)
			filterByDiv.DeleteFirst(bridge.LikeTag("label"))

			textBox := filterByDiv.FindFirst(bridge.LikeInput)
			textBox.GiveAttributes(textBoxHTMLX)

			table.filterInputs[i] = &filterByDiv

			// find sort slot made on New table
			filterSlot := headers[i].FindFirst(bridge.LikeAttribute("id", table.conf.filterSlotID))

			// fallback to just placing in th
			if filterSlot == nil {
				headers[i].Children = append(headers[i].Children, filterByDiv)
				continue
			}

			filterSlot.Children = append(filterSlot.Children, filterByDiv)
		}

		return nil
	}
	return defineFn
}
