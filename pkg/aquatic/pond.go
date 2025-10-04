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
type Stock map[*regexp.Regexp]Fish

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
type Pond struct {
	options     NewPondOptions
	pathBase    string
	templateDir string

	// OnCatch is a fn where you can provide the pond data for a template
	// to access.
	OnCatch func(r *http.Request) any

	// strictly for small fish to be used by tuna and sardines
	shad map[string]*Fish

	// fish are the items available for catch in a pond
	fish map[string][]Fish

	// licenses are required for any fish to be caught
	licenses []License
}

// FlowsInto can make global fish in one pond apply to another pond
// Note that only anchovy and clown are allowed to flow (assets)
// Useful to setup 2 ponds. one for assets, one for pages
func (p *Pond) FlowsInto(p2 *Pond) {
	for _, f := range p.shad {
		p2.shad[f.filePath] = f
	}
}

// Stock puts a stock into the pond. They will find their matches
// and be gobbled. So you can set fish bait and licenses, and
// feed then into the pond so the ponds fish inherit their stuff.
// Regex match done against relative file path to pond base dir
func (p *Pond) Stock(stock Stock) {
	for stockFishRegex, stockFish := range stock {
		found := false
		for _, pondFish := range p.FishFinder() {
			matched := stockFishRegex.Match([]byte(pondFish.scopedFilePath))
			if !matched {
				continue
			}

			found = true
			pondFish.Gobble(&stockFish)
		}
		if !found {
			fmt.Println("did not find matching fish for regex: " + stockFishRegex.String())
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

	p.options = options

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
	templateDir = filepath.Clean(templateDir)
	templateDir = filepath.ToSlash(templateDir)
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

	smallFishes := []*Fish{}
	bigFishes := []*Fish{}

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

	if p.shad == nil && isRoot {
		if p.shad == nil {
			p.shad = make(map[string]*Fish, len(smallFishes))
		}
		for _, f := range smallFishes {
			p.shad[f.filePath] = f
		}

	} else if p.options.GlobalSmallFish {
		if p.shad == nil {
			p.shad = make(map[string]*Fish, len(smallFishes))
		}
		for _, c := range smallFishes {
			p.shad[c.filePath] = c
		}
	}

	for _, pageItem := range bigFishes {
		for _, c := range smallFishes {
			pageItem.school = append(pageItem.school, *c)
		}

		itemsDeref := p.fish
		_, exists := itemsDeref[pathBase]
		if !exists {
			itemsDeref[pathBase] = []Fish{}
		}
		itemsDeref[pathBase] = append(itemsDeref[pathBase], *pageItem)
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

	var tw *tabwriter.Writer

	if verbose {
		tw = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		tw.Write([]byte("kind\tpattern\tfile\n"))
	}

	// allows us to collect fish before
	fishToRegister := make(map[string]*Fish)

	for _, child := range p.shad {
		if child.kind == FishKindTuna {
			// unreachable
			continue
		}
		fishToRegister[child.pattern] = child
	}

	// all dirs
	for path, fishes := range p.fish {
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

			for _, child := range fish.school {
				if child.kind == FishKindTuna {
					// unreachable
					continue
				}

				fishToRegister[child.pattern] = &child
			}
		}
	}

	sortedFish := make([]*Fish, 0, len(fishToRegister))
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
		mux.Handle(fish.pattern, reel(fish, p))
	}

	if tw != nil {
		tw.Flush()
	}

	return mux
}
