package internal

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"text/tabwriter"
)

func htmlxHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/javascript")
	w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", browserCacheDurationSeconds))
	w.Header().Add("Content-Length", strconv.Itoa(len(htmlx)))
	w.Write(htmlx)
}

// Pond is a collection of files from a dir with functions
// to get a server running
type Pond struct {
	items          map[string][]Item
	pathBase       string
	templateDir    string
	globalChildren []Item
}

// NewPond provides a new pond based on dir
func NewPond(templateDirPath string) (*Pond, error) {
	p := Pond{
		items: map[string][]Item{},
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	templateDir := filepath.Join(wd, templateDirPath)
	p.templateDir = templateDir

	err = p.collect(templateDir)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// collect will gather html and css from template dir
func (p *Pond) collect(pathBase string) error {
	if p.items == nil {
		p.items = map[string][]Item{}
	}
	entries, err := os.ReadDir(pathBase)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrNoTemplateDir
		}
		return err
	}

	isRoot := pathBase == p.templateDir

	children := []Item{}

	if p.globalChildren != nil {
		for _, e := range p.globalChildren {
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

		item, err := newItem(e, pathBase, p.templateDir)
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

		itemsDeref := p.items
		_, exists := itemsDeref[pathBase]
		if !exists {
			itemsDeref[pathBase] = []Item{}
		}
		itemsDeref[pathBase] = append(itemsDeref[pathBase], *pageItem)
	}

	if p.globalChildren == nil && isRoot {
		p.globalChildren = children
	}

	// now we can look at nested dirs
	for _, e := range dirs {
		p.collect(filepath.Join(pathBase, e.Name()))
	}

	return nil
}

// CastLines provides a mux to with patterns based on go templates in the specified directory
func (p *Pond) CastLines(verbose bool) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/assets/htmlx.2.0.4.js", htmlxHandler)

	// prevents duplicate pattern registration
	// expected since children share stylesheets
	pattensAdded := map[string]bool{}

	var tw *tabwriter.Writer

	if verbose {
		tw = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		tw.Write([]byte("kind\tpattern\n"))
		tw.Write([]byte("--\t--\n"))
	}

	for path, fish := range p.items {
		if len(fish) == 0 {
			fmt.Printf("no patterns for: %s\n", path)
			continue
		}

		for _, item := range fish {
			if _, exists := pattensAdded[item.Pattern]; exists {
				continue
			}
			if item.kind != htmlItemKindPage {
				continue
			}
			if tw != nil {
				tw.Write(fmt.Appendf(nil, "page\t%s\n", item.Pattern))
			}

			mux.HandleFunc(item.Pattern, item.handler)
			pattensAdded[item.Pattern] = true

			for _, child := range item.children {
				if _, exists := pattensAdded[child.Pattern]; exists {
					continue
				}

				if child.kind == htmlItemKindStyle {
					if tw != nil {
						tw.Write(fmt.Appendf(nil, "style\t%s\n", child.Pattern))
					}
					mux.HandleFunc(child.Pattern, child.handler)
					pattensAdded[child.Pattern] = true
				}

				if child.kind == htmlItemKindIsland {
					if tw != nil {
						tw.Write(fmt.Appendf(nil, "island\t%s\n", child.Pattern))
					}

					mux.HandleFunc(child.Pattern, child.handler)
					pattensAdded[child.Pattern] = true
				}

			}
		}
	}
	if tw != nil {
		tw.Flush()
	}

	return mux
}
