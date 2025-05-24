// Package internal provides the inner workings
package internal

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func isIsland(s string) bool {
	return strings.HasPrefix(s, "_")
}

func collect(items *map[string][]HTMLItem, pathBase, templateDir string, globalIslands []HTMLItem) error {
	if items == nil {
		items = &map[string][]HTMLItem{}
	}
	entries, err := os.ReadDir(pathBase)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrNoTemplateDir
		}
		return err
	}

	isRoot := pathBase == templateDir

	// ReadDir is sorted by filename and island templates start with "_" they will be collected first
	// but to ensure (and fix) for nested dir we sort too
	sort.Slice(entries, func(i, j int) bool {
		if !isRoot {
			return strings.Compare(entries[i].Name(), entries[j].Name()) > 0
		}
		return strings.Compare(entries[i].Name(), entries[j].Name()) < 0
	})

	islands := []HTMLItem{}

	if globalIslands != nil {
		for _, e := range globalIslands {
			islands = append(islands, e)
		}
	}

	// first just looking at the file entries so we can
	// snag global islands before going down nested dirs
	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		info, err := e.Info()
		if err != nil {
			return err
		}

		ext := filepath.Ext(info.Name())
		if ext != ".html" {
			continue
		}

		name := strings.TrimSuffix(info.Name(), ext)

		filePath := filepath.Join(pathBase, info.Name())
		pattern := filepath.Join(pathBase, name)
		pattern = strings.Replace(pattern, templateDir, "", 1)
		pattern = strings.ReplaceAll(pattern, "\\", "/")
		pattern = strings.ReplaceAll(pattern, "//", "/")

		if isIsland(info.Name()) {
			islands = append(islands, HTMLItem{
				kind:     htmlItemKindIsland,
				pattern:  pattern,
				filePath: filePath,
			})
			continue
		}

		// At this point we have gone through all islands as the arr is ordered
		if globalIslands == nil && isRoot {
			globalIslands = islands
		}

		itemsDeref := *items
		_, exists := itemsDeref[pathBase]
		if !exists {
			itemsDeref[pathBase] = []HTMLItem{}
		}

		itemsDeref[pathBase] = append(itemsDeref[pathBase], HTMLItem{
			kind:     htmlItemKindPage,
			pattern:  pattern,
			filePath: filePath,
			islands:  islands,
		})
	}

	// now we can look at nested dirs
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		collect(items, filepath.Join(pathBase, e.Name()), templateDir, globalIslands)
	}

	return nil
}
