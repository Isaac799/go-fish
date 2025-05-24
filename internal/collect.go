// Package internal provides the inner workings
package internal

import (
	"errors"
	"os"
	"path/filepath"
)

// collect will gather html and css from template dir. Needs optimization
func collect(items *map[string][]Item, pathBase, templateDir string, globalChildren []Item) error {
	if items == nil {
		items = &map[string][]Item{}
	}
	entries, err := os.ReadDir(pathBase)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrNoTemplateDir
		}
		return err
	}

	isRoot := pathBase == templateDir

	children := []Item{}

	if globalChildren != nil {
		for _, e := range globalChildren {
			children = append(children, e)
		}
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		htmlItem, err := newItem(e, pathBase, templateDir)
		if errors.Is(err, ErrInvalidExtension) {
			continue
		}
		if err != nil {
			return err
		}
		if htmlItem.kind != htmlItemKindStyle && htmlItem.kind != htmlItemKindIsland {
			continue
		}
		children = append(children, *htmlItem)
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		htmlItem, err := newItem(e, pathBase, templateDir)
		if errors.Is(err, ErrInvalidExtension) {
			continue
		}
		if err != nil {
			return err
		}
		if htmlItem.kind != htmlItemKindPage {
			continue
		}

		for _, c := range children {
			htmlItem.children = append(htmlItem.children, c)
		}

		itemsDeref := *items
		_, exists := itemsDeref[pathBase]
		if !exists {
			itemsDeref[pathBase] = []Item{}
		}
		itemsDeref[pathBase] = append(itemsDeref[pathBase], *htmlItem)
	}

	if globalChildren == nil && isRoot {
		globalChildren = children
	}

	// now we can look at nested dirs
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		collect(items, filepath.Join(pathBase, e.Name()), templateDir, globalChildren)
	}

	return nil
}
