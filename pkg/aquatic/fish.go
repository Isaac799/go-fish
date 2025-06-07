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
	// FishKindMackerel is essential to a healthy pond.
	// Is a "system" fish provided by me. Not discovered in file system.
	// Not meant to be caught. Handled different than most other fish.
	// Similar to a global sardine usage but instead of reading file
	// it uses pre-defined bytes.
	FishKindMackerel
)

const (
	// browserCacheDurationSeconds is used to cache documents
	// such as .css. To help prevent invalid cache we replace
	// the names with a hash of their content
	browserCacheDurationSeconds = 86400 // 1 day
)

// mackerelHTMLElement provides a system fish for
// the all powerful element.
func mackerelHTMLElement[K any]() Fish[K] {
	// A template element that works with bridge.HTMLelement.
	// Whitespace sensitive because text area value is inner text.
	elementTemplate := []byte(`{{define "_element"}}{{if .Tag}}{{if .SelfClosing}}<{{.Tag}}{{range $key, $value := .Attributes}} {{$key}}="{{$value}}" {{end}} />{{if .Children}}{{range $key, $value := .Children}}{{template "_element" $value}}{{end}}{{end}}{{else}}<{{.Tag}} {{range $key, $value := .Attributes}} {{$key}}="{{$value}}" {{end}}>{{.InnerText}}{{range $key, $value := .Children}} {{template "_element" $value}}{{end}}</{{.Tag}}>{{end}}{{end}}{{end}}`)
	randomStr := "3b5d5c3712955042212316173ccf37be"
	mackerel := Fish[K]{
		kind:      FishKindMackerel,
		isLanding: false,
		mime:      "text/html",
		bytes:     elementTemplate,
		// used by key seeing which fish eaten when handling tuna or sardine.
		// Value just needs to be unique.
		templateName: randomStr,
		// used by key in global fish.
		// Value just needs to be unique
		filePath: randomStr,
	}
	return mackerel
}

// License is a requirement to catch a fish.
// acts as a middleware. Return true if license is passed
type License func(next http.Handler) http.Handler

type masterBait[T, K any] struct {
	Local  K
	Global T
}

// Bait is to be gobbled up by a fish before catching it.
// A func that has access to the request and returns template data
type Bait[T any] func(r *http.Request) T

// Fish is an item found form the template dir.
type Fish[K any] struct {
	kind           int
	isLanding      bool
	mime           string
	hash           string
	templateName   string
	pattern        string
	scopedFilePath string
	filePath       string
	children       []Fish[K]

	// FishKindMackerel have it pre-defined. Other fish
	// cached read of template since it ony needs to be read once
	bytes []byte

	// Licenses is a collection of licenses a user must have
	// to catch a fish. Checked after pond licenses, in the
	// order added. To catch a fish all pond and fish licenses
	// must be met.
	Licenses []License

	// Bait fn is called and the result is passed into the
	// executed template, or eaten by the fish before caught
	Bait Bait[K]

	// Tackle helps catch a fish.
	// Given to a template to help transform the data.
	Tackle template.FuncMap
}

// Patten is the pattern of a fish used by mux
func Patten[K any](f *Fish[K]) string {
	return f.pattern
}

// Gobble has one fish gobble up another. Gaining its Licenses, Tackle, and Bait (if not already has some).
func Gobble[T any](f *Fish[T], f2 *Fish[T]) {
	if f.Bait == nil && f2.Bait != nil {
		f.Bait = f2.Bait
	}
	if f.Licenses == nil {
		f.Licenses = make([]License, 0, len(f2.Licenses))
	}
	for _, l := range f2.Licenses {
		f.Licenses = append(f.Licenses, l)
	}
	if f.Tackle == nil {
		f.Tackle = make(template.FuncMap, len(f2.Tackle))
	}
	maps.Copy(f.Tackle, f2.Tackle)
	for i := range f.children {
		if f.children[i].Bait == nil && f2.Bait != nil {
			f.children[i].Bait = f2.Bait
		}
		if f.children[i].kind != FishKindSardine {
			continue
		}
		if f.children[i].Licenses == nil {
			f.children[i].Licenses = make([]License, 0, len(f2.Licenses))
		}
		for _, l := range f2.Licenses {
			f.children[i].Licenses = append(f.children[i].Licenses, l)
		}
		if f.children[i].Tackle == nil {
			f.children[i].Tackle = make(template.FuncMap, len(f2.Tackle))
		}
		maps.Copy(f.children[i].Tackle, f2.Tackle)
	}
}

// Kind reads back the kind of a fish
func Kind[T, K any](f *Fish[K]) int {
	return f.kind
}

func newFish[T, K any](entry os.DirEntry, pathBase string, pond *Pond[T, K]) (*Fish[K], error) {
	info, err := entry.Info()
	if err != nil {
		return nil, err
	}

	ext := filepath.Ext(info.Name())
	mime := mime.TypeByExtension(ext)

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

	name := info.Name()
	if kind == FishKindTuna || kind == FishKindSardine {
		name = strings.TrimSuffix(info.Name(), ext)
	}
	templateName := name

	// since I want to cache styling while preventing
	// an invalid cache we make the name based on a hash
	// of its content
	file, err := os.Open(filepath.Join(pathBase, entry.Name()))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	b, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	hash := fmt.Sprintf("%x", md5.Sum(b))

	filePath := filepath.Join(pathBase, info.Name())
	scopedFilePath := strings.Replace(filePath, pond.templateDir, "", 1)

	pattern := filepath.Join(pathBase, name)
	pattern = strings.Replace(pattern, pond.templateDir, "", 1)
	pattern = strings.ReplaceAll(pattern, " ", "-")

	pattern = strings.ToLower(pattern)

	isLanding := false
	if kind == FishKindTuna {
		fileParts := strings.Split(pathBase, "/")
		if len(fileParts) > 0 {
			parentDir := fileParts[len(fileParts)-1]
			isLanding = parentDir == name
		}
	}

	if kind == FishKindTuna || kind == FishKindSardine {
		patternParts := strings.Split(pattern, ".")
		newPatternParts := []string{}
		for i, e := range patternParts {
			if (i+1)%2 == 0 {
				param := "{" + e + "}"
				newPatternParts = append(newPatternParts, param)
				continue
			}
			newPatternParts = append(newPatternParts, e)
		}
		pattern = strings.Join(newPatternParts, "/")
	}

	pattern = strings.ReplaceAll(pattern, "\\", "/")
	pattern = strings.ReplaceAll(pattern, "//", "/")

	if kind == FishKindTuna {
		if isLanding {
			pattern = strings.TrimSuffix(pattern, templateName)
		}
		if pattern != "/" && strings.HasSuffix(pattern, "/") {
			pattern = strings.TrimSuffix(pattern, "/")
		}
	}

	f := Fish[K]{
		kind:           kind,
		mime:           mime,
		hash:           hash,
		pattern:        pattern,
		isLanding:      isLanding,
		templateName:   templateName,
		filePath:       filePath,
		scopedFilePath: scopedFilePath,
		Licenses:       []License{},
	}

	return &f, nil
}

// TemplateBuffer will wrap a file content in the define syntax.
// Enforcing template name scheme and reducing template lines n - 2.
func templateBytes[K any](f *Fish[K]) ([]byte, error) {
	if f.bytes != nil {
		return f.bytes, nil
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
	f.bytes = buffer

	return buffer, nil
}
