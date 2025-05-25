// Package internal provides the inner workings
package internal

import (
	"errors"
	"os"
	"path/filepath"
)

// collect will gather html and css from template dir
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

	pageItems := []*Item{}
	dirs := []os.DirEntry{}

	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, e)
			continue
		}

		item, err := newItem(e, pathBase, templateDir)
		if errors.Is(err, ErrInvalidExtension) {
			continue
		}
		if err != nil {
			return err
		}

		if item.kind == htmlItemKindPage {
			pageItems = append(pageItems, item)
			continue
		}

		children = append(children, *item)
	}

	for _, pageItem := range pageItems {
		for _, c := range children {
			pageItem.children = append(pageItem.children, c)
		}

		itemsDeref := *items
		_, exists := itemsDeref[pathBase]
		if !exists {
			itemsDeref[pathBase] = []Item{}
		}
		itemsDeref[pathBase] = append(itemsDeref[pathBase], *pageItem)
	}

	if globalChildren == nil && isRoot {
		globalChildren = children
	}

	// now we can look at nested dirs
	for _, e := range dirs {
		collect(items, filepath.Join(pathBase, e.Name()), templateDir, globalChildren)
	}

	return nil
}
