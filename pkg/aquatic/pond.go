package aquatic

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/tabwriter"
)

var (
	// ErrNoTemplateDir is given if I cannot find the desired dir to setup a new mux
	ErrNoTemplateDir = errors.New("cannot find template directory relative to working dir")
	// ErrInvalidExtension is given if a file is discoverd that I did not anticipate
	ErrInvalidExtension = errors.New("invalid file extension")
)

var fishKindStr = map[int]string{
	FishKindTuna:    "Tuna",
	FishKindSardine: "Sardine",
	FiskKindClown:   "Clown",
	FiskKindAnchovy: "Anchovy",
}

// Stock enables developer to provide what fish they think
// the pond should have. Keyed by file name (without extension)
type Stock[T, K any] map[*regexp.Regexp]Fish[K]

// NewPondOptions gives the options available when creating a new pond
type NewPondOptions struct {
	// Licenses for a pond are applied to all fish in the pond
	// and are checked before a fish license in the order added.
	// To catch a fish all pond and fish licenses must be met.
	Licenses []License
	// GlobalSmallFish makes all Anchovy, Sardine, and Clown fish
	// global scoped no matter where they are. Useful for an assets pond
	// that flows into another pond.
	GlobalSmallFish bool
}

// Pond is a collection of files from a dir with functions
// to get a server running
type Pond[T, K any] struct {
	options     NewPondOptions
	pathBase    string
	templateDir string
	GlobalBait  Bait[T]
	// strictly for small fish
	globalSmallFish map[string]*Fish[K]
	// fish are the items available for catch in a pond
	fish map[string][]Fish[K]
	// licenses are required for any fish to be caught
	licenses []License
}

// FlowsInto can make global fish in one pond apply to another pond
// Note that only anchovy and clown are allowed to flow (assets)
// Useful to setup 2 ponds. one for assets, one for pages
func FlowsInto[T, K any](p *Pond[T, K], p2 *Pond[T, K]) {
	for _, f := range p.globalSmallFish {
		p2.globalSmallFish[f.filePath] = f
	}
}

// StockPond puts a stock into the pond. They will find their matches
// and be gobbled. So you can set fish bait and licenses, and
// feed then into the pond so the ponds fish inherit their stuff.
// Regex match done against relative file path to pond base dir
func StockPond[T, K any](p *Pond[T, K], stock Stock[T, K]) {
	for stockFishRegex, stockFish := range stock {
		found := false
		for _, pondFish := range FishFinder(p) {
			matched := stockFishRegex.Match([]byte(pondFish.scopedFilePath))
			if !matched {
				continue
			}

			found = true
			Gobble(pondFish, &stockFish)
		}
		if !found {
			fmt.Println("did not find matching fish for regex: " + stockFishRegex.String())
		}
	}
}

// FishFinder provides a slice of all fish
func FishFinder[T, K any](p *Pond[T, K]) []*Fish[K] {
	all := []*Fish[K]{}
	for _, fishes := range p.fish {
		for i := range fishes {
			all = append(all, &fishes[i])
		}
	}
	return all
}

// NewPond provides a new pond based on dir
func NewPond[T, K any](templateDirPath string, options NewPondOptions) (Pond[T, K], error) {

	p := Pond[T, K]{
		fish:     map[string][]Fish[K]{},
		licenses: options.Licenses,
	}

	p.options = options

	if p.licenses == nil {
		// do this to avoid nil deref in handle funcs
		// prefer this to checking for nil
		p.licenses = make([]License, 0, 0)
	}

	wd, err := os.Getwd()
	if err != nil {
		return p, err
	}
	templateDir := filepath.Join(wd, templateDirPath)
	p.templateDir = templateDir

	err = collect(&p, templateDir)
	if err != nil {
		return p, err
	}
	return p, nil
}

// collect will gather html and css from template dir
func collect[T, K any](p *Pond[T, K], pathBase string) error {
	if p.fish == nil {
		p.fish = map[string][]Fish[K]{}
	}
	entries, err := os.ReadDir(pathBase)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrNoTemplateDir
		}
		return err
	}

	isRoot := pathBase == p.templateDir

	_elementFish := mackerelHTMLElement[K]()
	smallFishes := []*Fish[K]{
		&_elementFish,
	}
	bigFishes := []*Fish[K]{}

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
			bigFishes = append(bigFishes, item)
			continue
		}

		smallFishes = append(smallFishes, item)
	}

	if p.globalSmallFish == nil && isRoot {
		if p.globalSmallFish == nil {
			p.globalSmallFish = make(map[string]*Fish[K], len(smallFishes))
		}
		for _, f := range smallFishes {
			p.globalSmallFish[f.filePath] = f
		}

	} else if p.options.GlobalSmallFish {
		if p.globalSmallFish == nil {
			p.globalSmallFish = make(map[string]*Fish[K], len(smallFishes))
		}
		for _, c := range smallFishes {
			p.globalSmallFish[c.filePath] = c
		}
	}

	for _, pageItem := range bigFishes {
		for _, c := range smallFishes {
			pageItem.children = append(pageItem.children, *c)
		}

		itemsDeref := p.fish
		_, exists := itemsDeref[pathBase]
		if !exists {
			itemsDeref[pathBase] = []Fish[K]{}
		}
		itemsDeref[pathBase] = append(itemsDeref[pathBase], *pageItem)
	}

	// now we can look at nested dirs
	for _, e := range dirs {
		collect(p, filepath.Join(pathBase, e.Name()))
	}

	return nil
}

// CastLines provides a mux to with patterns based on go templates in the specified directory
func CastLines[T, K any](pond *Pond[T, K], verbose bool) *http.ServeMux {
	mux := http.NewServeMux()

	var tw *tabwriter.Writer

	if verbose {
		tw = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		tw.Write([]byte("kind\tpattern\tfile\n"))
	}

	// allows us to collect fish before
	fishToRegister := make(map[string]*Fish[K])

	for _, child := range pond.globalSmallFish {
		if child.kind == FishKindMackerel {
			// not to be served
			continue
		}
		if child.kind == FishKindTuna {
			// unreachable
			continue
		}
		fishToRegister[child.pattern] = child
	}

	// all dirs
	for path, fishes := range pond.fish {
		if len(fishes) == 0 {
			fmt.Printf("no patterns for: %s\n", path)
			continue
		}

		// all fish in dir
		for _, fish := range fishes {
			if fish.kind != FishKindTuna {
				continue
			}

			fishToRegister[fish.pattern] = &fish

			for _, child := range fish.children {
				if child.kind == FishKindTuna {
					// unreachable
					continue
				}

				fishToRegister[child.pattern] = &child
			}
		}
	}

	sortedFish := make([]*Fish[K], 0, len(fishToRegister))
	for pattern, fish := range fishToRegister {
		if len(pattern) == 0 {
			continue
		}
		sortedFish = append(sortedFish, fish)
	}

	// ensure more explicit routes matched first
	sort.Slice(sortedFish, func(i, j int) bool {
		return strings.Compare(sortedFish[i].pattern, sortedFish[j].pattern) > 0
	})

	for _, fish := range sortedFish {
		if tw != nil {
			tw.Write(fmt.Appendf(nil, "%s\t%s\t%s\n", fishKindStr[fish.kind], fish.pattern, fish.scopedFilePath))
		}
		mux.Handle(fish.pattern, reel(fish, pond))
	}

	if tw != nil {
		tw.Flush()
	}

	return mux
}
