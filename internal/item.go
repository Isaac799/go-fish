package internal

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"
)

// HTMLItem is an item found form the template dir.
type HTMLItem struct {
	kind     int
	Pattern  string
	filePath string
	islands  []HTMLItem
}

func (hi *HTMLItem) templateName() string {
	if len(hi.Pattern) == 0 {
		return "unknown"
	}
	parts := strings.Split(hi.Pattern, "/")
	return parts[len(parts)-1]
}

func (hi *HTMLItem) handler(w http.ResponseWriter, _ *http.Request) {
	collectedFilePaths := []string{}
	for _, e := range hi.islands {
		if e.kind != htmlItemKindIsland {
			continue
		}
		collectedFilePaths = append(collectedFilePaths, e.filePath)
	}
	collectedFilePaths = append(collectedFilePaths, hi.filePath)

	key := hi.templateName()

	t := template.New(key)
	parsed, err := t.ParseFiles(collectedFilePaths...)

	if err != nil {
		fmt.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = parsed.ExecuteTemplate(w, key, hi)
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
