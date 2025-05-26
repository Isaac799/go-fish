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

// Stock enables developer to provide what fish they think
// the pond should have. Keyed by file name (without extension)
type Stock map[string]Fish

// NewPondOptions gives the options available when creating a new pond
type NewPondOptions struct {
	// Licenses for a pond are applied to all fish in the pond
	// and are checked before a fish license in the order added.
	// To catch a fish all pond and fish licenses must be met.
	Licenses []License
}

// Pond is a collection of files from a dir with functions
// to get a server running
type Pond struct {
	pathBase       string
	templateDir    string
	globalChildren []Fish
	// fish are the items available for catch in a pond
	fish map[string][]Fish
	// licenses are required for any fish to be caught
	licenses []License
}

// Stock puts a stock into the pond. They will find their equal
// and be gobbled. So you can set fish bait and licenses, and
// feed then into the pond so the ponds fish inherit their stuff
func (p *Pond) Stock(stock Stock) {
	for stockPattern, stockFish := range stock {
		found := false
		for _, pondFish := range p.FishFinder() {
			if stockPattern != pondFish.templateName {
				continue
			}
			found = true
			pondFish.Gobble(stockFish)
		}
		if !found {
			fmt.Println("did not find matching fish for page")
			os.Exit(1)
		}
	}
}

// FishFinder provides a slice of all fish
func (p *Pond) FishFinder() []*Fish {
	all := []*Fish{}
	for _, fishes := range p.fish {
		for i := range fishes {
			all = append(all, &fishes[i])
		}
	}
	return all
}

// NewPond provides a new pond based on dir
func NewPond(templateDirPath string, options NewPondOptions) (Pond, error) {
	p := Pond{
		fish:     map[string][]Fish{},
		licenses: options.Licenses,
	}

	if p.licenses == nil {
		// do this to avoid nil deref in handle funcs
		// prefer this to checking for nil
		p.licenses = make([]License, 0)
	}

	wd, err := os.Getwd()
	if err != nil {
		return p, err
	}
	templateDir := filepath.Join(wd, templateDirPath)
	p.templateDir = templateDir

	err = p.collect(templateDir)
	if err != nil {
		return p, err
	}
	return p, nil
}

// collect will gather html and css from template dir
func (p *Pond) collect(pathBase string) error {
	if p.fish == nil {
		p.fish = map[string][]Fish{}
	}
	entries, err := os.ReadDir(pathBase)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrNoTemplateDir
		}
		return err
	}

	isRoot := pathBase == p.templateDir

	children := []Fish{}

	if p.globalChildren != nil {
		for _, e := range p.globalChildren {
			children = append(children, e)
		}
	}

	pageItems := []*Fish{}
	dirs := []os.DirEntry{}

	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, e)
			continue
		}

		item, err := newFish(e, pathBase, p)
		if errors.Is(err, ErrInvalidExtension) {
			continue
		}
		if err != nil {
			return err
		}

		if item.kind == FishKindTuna {
			pageItems = append(pageItems, item)
			continue
		}

		children = append(children, *item)
	}

	for _, pageItem := range pageItems {
		for _, c := range children {
			pageItem.children = append(pageItem.children, c)
		}

		itemsDeref := p.fish
		_, exists := itemsDeref[pathBase]
		if !exists {
			itemsDeref[pathBase] = []Fish{}
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

	for path, fish := range p.fish {
		if len(fish) == 0 {
			fmt.Printf("no patterns for: %s\n", path)
			continue
		}

		for _, item := range fish {
			if _, exists := pattensAdded[item.pattern]; exists {
				continue
			}
			if item.kind != FishKindTuna {
				continue
			}
			if tw != nil {
				tw.Write(fmt.Appendf(nil, "tuna\t%s\n", item.pattern))
			}

			mux.HandleFunc(item.pattern, item.handler)
			pattensAdded[item.pattern] = true

			for _, child := range item.children {
				if _, exists := pattensAdded[child.pattern]; exists {
					continue
				}
				if child.kind == FishKindTuna {
					// unreachable
					continue
				}
				if child.kind == FiskKindClown {
					if tw != nil {
						tw.Write(fmt.Appendf(nil, "clown\t%s\n", child.pattern))
					}
				}
				if child.kind == FishKindSardine {
					if tw != nil {
						tw.Write(fmt.Appendf(nil, "sardine\t%s\n", child.pattern))
					}
				}
				if child.kind == FiskKindAnchovy {
					if tw != nil {
						tw.Write(fmt.Appendf(nil, "anchovy\t%s\n", child.pattern))
					}
				}
				mux.HandleFunc(child.pattern, child.handler)
				pattensAdded[child.pattern] = true
			}
		}
	}
	if tw != nil {
		tw.Flush()
	}

	return mux
}
