package internal

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// NewMux provides a mux to with patterns based on go templates in the specified directory
func NewMux(templateDirPath string) (*http.ServeMux, error) {
	mux := http.NewServeMux()

	items := map[string][]HTMLItem{}
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	templateDir := filepath.Join(wd, templateDirPath)
	err = collect(&items, templateDir, templateDir, nil)
	if err != nil {
		return nil, err
	}
	for path, items := range items {
		if len(items) == 0 {
			fmt.Printf("no patterns for: %s\n", path)
			continue
		}

		for _, item := range items {
			if item.kind == htmlItemKindIsland {
				continue
			}
			fmt.Printf("pattern: %s\n", item.pattern)
			mux.HandleFunc(item.pattern, item.handler)
		}
	}
	return mux, nil
}
