// Package aquatic has fish and pond.
// Fish are html documents either generated from templates or served.
// Pond hold all the fish and can inherit each others fish.
package aquatic

import (
	"crypto/md5"
	"fmt"
	"io"
	"maps"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	// FishKindTuna is a big fish. Served as a page. Consumes Sardines.
	// Identified by mime [ text/html ].
	// Not cached.
	FishKindTuna = iota
	// FishKindSardine is a small fish. Used by tuna. Smaller templates, served standalone too.
	// Identified by mime [ text/html ] & underscore prefix.
	// Not cached.
	FishKindSardine
	// FiskKindClown is a decorative fish. Used in head of document.
	// Identified by mime [ text/css | text/javascript ].
	// Is cached & name from hash.
	FiskKindClown
	// FiskKindAnchovy is supportive of the tuna.
	// Identified by mime [ image | audio | video ].
	// Is cached.
	FiskKindAnchovy
)

const (
	// browserCacheDurationSeconds is used to cache documents
	// such as .css. To help prevent invalid cache we replace
	// the names with a hash of their content
	browserCacheDurationSeconds = 86400 // 1 day
)

// BeforeCatchFn is a requirement to catch a fish, a middleware.
type BeforeCatchFn func(next http.Handler) http.Handler

// Fish is an item found form the template dir.
type Fish struct {
	kind           int
	isLanding      bool
	mime           string
	hash           string
	templateName   string
	pattern        string
	scopedFilePath string
	filePath       string

	// fish found in same dir
	school []Fish

	// coral is bytes of template since it ony needs to be read once.
	// Populated on first time parsing
	coral []byte

	// reef is combination of all coral (templates bytes) combined into
	// one beautiful string for the parse template use. It includes all
	// coral for other global, children... fish. Saves me from having to
	// rebuild known template combinations for sardine and tuna.
	reef []byte

	// bobber stays above a tuna. Is the head of the html document.
	// Only relevant for tuna. Saved for reuse after first determined.
	bobber []byte

	// BeforeCatch is middleware before a fish is caught.
	// Checked after pond middleware in order.
	BeforeCatch []BeforeCatchFn

	// OnCatch provides data that a fish as access to when it is caught.
	OnCatch func(r *http.Request) any

	// Tackle helps catch a fish.
	// Given to a template to help transform the data.
	Tackle template.FuncMap
}

// Gobble has one fish gobble up another. Gaining its Before and On catch fns, and Tackle (if not already has some).
func (f *Fish) Gobble(f2 *Fish) {
	if f.OnCatch == nil && f2.OnCatch != nil {
		f.OnCatch = f2.OnCatch
	}
	if f.BeforeCatch == nil {
		f.BeforeCatch = make([]BeforeCatchFn, 0, len(f2.BeforeCatch))
	}
	for _, l := range f2.BeforeCatch {
		f.BeforeCatch = append(f.BeforeCatch, l)
	}
	if f.Tackle == nil {
		f.Tackle = make(template.FuncMap, len(f2.Tackle))
	}
	maps.Copy(f.Tackle, f2.Tackle)
	for i := range f.school {
		if f.school[i].OnCatch == nil && f2.OnCatch != nil {
			f.school[i].OnCatch = f2.OnCatch
		}
		if f.school[i].kind != FishKindSardine {
			continue
		}
		if f.school[i].BeforeCatch == nil {
			f.school[i].BeforeCatch = make([]BeforeCatchFn, 0, len(f2.BeforeCatch))
		}
		for _, l := range f2.BeforeCatch {
			f.school[i].BeforeCatch = append(f.school[i].BeforeCatch, l)
		}
		if f.school[i].Tackle == nil {
			f.school[i].Tackle = make(template.FuncMap, len(f2.Tackle))
		}
		maps.Copy(f.school[i].Tackle, f2.Tackle)
	}
}

func patternFromRelativePath(relative, ext string, isHTML bool) string {
	pattern := filepath.ToSlash(relative)
	pattern = strings.ReplaceAll(pattern, " ", "-")
	pattern = strings.ToLower(pattern)

	if !isHTML {
		return pattern
	}

	pattern = strings.TrimSuffix(pattern, ext)

	arr := strings.Split(pattern, "/")
	arr2 := []string{}

	for i, pathItem := range arr {
		if len(pathItem) == 0 {
			continue
		}

		// exclude file prefix eq dir
		if i <= 2 && len(arr2) > 0 && arr2[len(arr2)-1] == pathItem {
			continue
		}

		if !strings.Contains(pathItem, ".") {
			arr2 = append(arr2, pathItem)
			continue
		}

		arr := strings.Split(pathItem, ".")
		for k, s := range arr {
			if len(s) == 0 {
				continue
			}
			// exclude file prefix eq dir
			if k == 0 && len(arr2) > 0 && arr2[len(arr2)-1] == s {
				continue
			}
			if k%2 == 0 {
				arr2 = append(arr2, s)
				continue
			}
			arr2 = append(arr2, fmt.Sprintf("{%s}", s))
		}
	}

	pattern = fmt.Sprintf("/%s", strings.Join(arr2, "/"))
	return pattern
}

func newFish(entry os.DirEntry, rootPath string, pond *Pond) (*Fish, error) {
	rootPath = filepath.ToSlash(rootPath)

	info, err := entry.Info()
	if err != nil {
		return nil, err
	}
	ext := filepath.Ext(info.Name())

	// since I want to cache styling while preventing
	// an invalid cache we make the name based on a hash
	// of its content
	file, err := os.Open(filepath.Join(rootPath, entry.Name()))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	b, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	hash := fmt.Sprintf("%x", md5.Sum(b))

	// so windows mime types suck and using mime package not always work
	// e.g. windows not knowing what a woff2 file was and causing
	// browser to "rejected by sanitizer" due to incorrect mime type
	// so we need the fallback
	mime := mime.TypeByExtension(ext)
	if len(mime) == 0 {
		// fallback
		mime = http.DetectContentType(b)
	}

	if len(mime) == 0 {
		fmt.Println("gofish warn: cannot determine mime type of: ", ext)
		return nil, ErrInvalidExtension
	}

	kind := -1
	if strings.HasPrefix(mime, "text/html") {
		if strings.HasPrefix(info.Name(), "_") {
			kind = FishKindSardine
		} else {
			kind = FishKindTuna
		}
	} else if strings.HasPrefix(mime, "text/css") ||
		strings.HasPrefix(mime, "text/javascript") {
		kind = FiskKindClown
	} else if strings.HasPrefix(mime, "image") ||
		strings.HasPrefix(mime, "audio") ||
		strings.HasPrefix(mime, "video") ||
		strings.HasPrefix(mime, "font") {
		kind = FiskKindAnchovy
	}

	if kind == -1 {
		return nil, ErrInvalidExtension
	}

	absolute := filepath.Join(rootPath, info.Name())
	absolute = filepath.ToSlash(absolute)

	relative := strings.Replace(absolute, pond.templateDir, "", 1)
	relativeNoSuffix := strings.TrimSuffix(relative, ext)

	isHTML := kind == FishKindTuna || kind == FishKindSardine

	rootArr := strings.Split(rootPath, "/")
	relNow := strings.TrimPrefix(relativeNoSuffix, "/")
	isLanding := len(rootArr) > 0 && rootArr[len(rootArr)-1] == relNow

	var pattern string
	if isLanding {
		pattern = "/"
	} else {
		pattern = patternFromRelativePath(relative, ext, isHTML)
	}

	f := Fish{
		kind:           kind,
		mime:           mime,
		hash:           hash,
		pattern:        pattern,
		isLanding:      false,
		templateName:   strings.TrimSuffix(info.Name(), ext),
		filePath:       absolute,
		scopedFilePath: relative,
		BeforeCatch:    []BeforeCatchFn{},
	}

	return &f, nil
}

// cacheCoral will wrap a file content in the define syntax.
// Enforcing template name scheme and reducing template lines n - 2.
// Once coral is discovered for the first time it is saved in the fish for re use.
func (f *Fish) cacheCoral() ([]byte, error) {
	if f.coral != nil {
		return f.coral, nil
	}

	file, err := os.Open(f.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	prefix := fmt.Appendf(nil, "{{define \"%s\"}}", f.templateName)
	suffix := []byte("{{end}}")

	size := len(prefix) + int(info.Size()) + len(suffix)
	buffer := make([]byte, size)

	copy(buffer, prefix)
	n, err := file.Read(buffer[len(prefix) : len(prefix)+int(info.Size())])
	if err != nil {
		return nil, err
	}

	copy(buffer[len(prefix)+n:], suffix)
	f.coral = buffer

	return buffer, nil
}

// cacheReef combines the coral of dependent fish and itself.
// Once a reef is discovered for the first time it is saved in the fish for re use.
func (f *Fish) cacheReef(pond *Pond) ([]byte, error) {
	if f.reef != nil {
		return f.reef, nil
	}

	// a map to store the various fish needed to be eaten
	// by this fish to give it access to all templates available
	// to it. Populated in a significant way to enable scoping
	eaten := map[string][]byte{}
	size := 0

	// local sardines first to give the consumer (tuna or sardine)
	// access to its local dependent templates
	for _, e := range f.school {
		if e.kind != FishKindSardine {
			continue
		}
		if _, exists := eaten[e.templateName]; exists {
			continue
		}
		b, err := e.cacheCoral()
		if err != nil {
			return nil, err
		}
		size += len(b)
		eaten[e.templateName] = b
	}

	// global sardines come after local ones so they do not
	// overwrite local ones. So if _nav in global scope and
	// _nav in this fish dir we already consumed the local
	// one, and it cannot be re defined
	for _, e := range pond.shad {
		if e.kind != FishKindSardine {
			continue
		}
		if _, exists := eaten[e.templateName]; exists {
			continue
		}
		b, err := e.cacheCoral()
		if err != nil {
			return nil, err
		}
		size += len(b)
		eaten[e.templateName] = b
	}

	// finally we can consume the 'main' fish (tuna or sardine)
	// this is to ensure not re define if is sardine
	if _, exists := eaten[f.templateName]; !exists {
		b, err := f.cacheCoral()
		if err != nil {
			return nil, err
		}
		size += len(b)
		eaten[f.templateName] = b
	}

	// now the cool part, a sliding copy into a single pre
	// alloc buff using references to the fish bytes since
	buff := make([]byte, size)
	last := 0
	for _, v := range eaten {
		n := copy(buff[last:last+len(v)], v)
		last += n
	}

	f.reef = buff
	return buff, nil
}
