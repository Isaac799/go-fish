package internal

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

// License is a requirement to catch a fish.
// acts as a middleware. Return true if license is passed
type License = func(w http.ResponseWriter, r *http.Request) bool

// Bait is to be gobbled up by a fish before catching it.
// A func that has access to the request and returns template data
type Bait = func(r *http.Request) any

// Fish is an item found form the template dir.
type Fish struct {
	pond     *Pond
	kind     int
	mime     string
	pattern  string
	filePath string
	children []Fish

	// Licenses is a collection of licenses a user must have
	// to catch a fish. Checked after pond licenses, in the
	// order added. To catch a fish all pond and fish licenses
	// must be met.
	Licenses []License
	Bait     Bait
}

// AddLicense appends a license required to catch a fish
func (f *Fish) AddLicense(l License) {
	f.Licenses = append(f.Licenses, l)
}

// Kind reads back the kind of a fish
func (f *Fish) Kind() int {
	return f.kind
}

// Pattern reads back the pattern a fish will bite for
func (f *Fish) Pattern() string {
	return f.pattern
}

func newFish(entry os.DirEntry, pathBase string, pond *Pond) (*Fish, error) {
	info, err := entry.Info()
	if err != nil {
		return nil, err
	}

	ext := filepath.Ext(info.Name())
	mime := mime.TypeByExtension(ext)

	kind := -1

	fmt.Println(info.Name(), mime)

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
		strings.HasPrefix(mime, "video") {
		kind = FiskKindAnchovy
	}

	if kind == -1 {
		return nil, ErrInvalidExtension
	}

	name := info.Name()
	if kind == FishKindTuna || kind == FishKindSardine {
		name = strings.TrimSuffix(info.Name(), ext)
	}

	// since I want to cache styling while preventing
	// an invalid cache we make the name based on a hash
	// of its content
	if kind == FiskKindClown {
		f, err := os.Open(filepath.Join(pathBase, entry.Name()))
		if err != nil {
			return nil, err
		}
		defer f.Close()
		b, err := io.ReadAll(f)
		if err != nil {
			return nil, err
		}
		hash := md5.Sum(b)
		name = fmt.Sprintf("%x%s", hash, ext)
	}

	filePath := filepath.Join(pathBase, info.Name())
	pattern := filepath.Join(pathBase, name)
	pattern = strings.Replace(pattern, pond.templateDir, "", 1)
	pattern = strings.ReplaceAll(pattern, "\\", "/")
	pattern = strings.ReplaceAll(pattern, "//", "/")

	return &Fish{
		kind:     kind,
		mime:     mime,
		pattern:  pattern,
		filePath: filePath,
		Licenses: []License{},
		pond:     pond,
	}, nil

}

func (f *Fish) templateName() string {
	if len(f.pattern) == 0 {
		return "unknown"
	}
	parts := strings.Split(f.pattern, "/")
	return parts[len(parts)-1]
}

func (f *Fish) handlerSardine(w http.ResponseWriter, r *http.Request) {
	templateName := f.templateName()

	t := template.New(templateName)
	parsed, err := t.ParseFiles(f.filePath)

	// buffer for html doc
	b := []byte{}
	buff := bytes.NewBuffer(b)

	var pageData any
	if f.Bait != nil {
		pageData = f.Bait(r)
	}
	err = parsed.ExecuteTemplate(buff, templateName, pageData)
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/html")
	w.Header().Add("Content-Length", strconv.Itoa(len(buff.Bytes())))
	w.Write(buff.Bytes())
}

func (f *Fish) handlerClownAnchovy(w http.ResponseWriter) {
	file, err := os.Open(f.filePath)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()
	b, err := io.ReadAll(file)
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", mime.TypeByExtension(filepath.Ext(file.Name())))
	w.Header().Add("Content-Length", strconv.Itoa(len(b)))
	w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", browserCacheDurationSeconds))
	w.Write(b)
}

func (f *Fish) handlerTuna(w http.ResponseWriter, r *http.Request) {
	// 3 main parts to the document I will add in between
	htmlStartHead := []byte(`<!DOCTYPE html><html lang="en">
<head>
    <meta charset="UTF-8">
	<script src="/assets/htmlx.2.0.4.js"></script>
    <meta name="viewport" content="width=device-width, initial-scale=1.0" >`)
	htmlEndHeadStartBody := []byte(`
</head>
<body>`)
	htmlEndBody := []byte(`</body></html>`)

	// gather islands
	collectedFilePaths := []string{}
	for _, e := range f.children {
		if e.kind != FishKindSardine {
			continue
		}
		collectedFilePaths = append(collectedFilePaths, e.filePath)
	}
	collectedFilePaths = append(collectedFilePaths, f.filePath)

	templateName := f.templateName()

	t := template.New(templateName)
	parsed, err := t.ParseFiles(collectedFilePaths...)

	// buffer for html doc
	b := []byte{}
	buff := bytes.NewBuffer(b)
	buff.Write(htmlStartHead)

	// styling
	for _, e := range f.children {
		if e.kind != FiskKindClown {
			continue
		}
		if strings.HasPrefix(e.mime, "text/css") {
			b := fmt.Appendf(nil, `<link rel="stylesheet" href="%s">`, e.pattern)
			buff.Write(b)
		}
		if strings.HasPrefix(e.mime, "text/javascript") {
			b := fmt.Appendf(nil, `<script src="%s"></script>`, e.pattern)
			buff.Write(b)
		}
	}
	buff.Write(htmlEndHeadStartBody)

	var pageData any
	if f.Bait != nil {
		pageData = f.Bait(r)
	}
	err = parsed.ExecuteTemplate(buff, templateName, pageData)
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	buff.Write(htmlEndBody)

	w.Header().Add("Content-Type", "text/html")
	w.Header().Add("Content-Length", strconv.Itoa(len(buff.Bytes())))
	w.Write(buff.Bytes())
}

func (f *Fish) handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	for _, license := range f.pond.licenses {
		allowed := license(w, r)
		if !allowed {
			return
		}
	}

	for _, license := range f.Licenses {
		allowed := license(w, r)
		if !allowed {
			return
		}
	}

	if f.kind == FishKindSardine {
		f.handlerSardine(w, r)
		return
	}

	if f.kind == FiskKindClown || f.kind == FiskKindAnchovy {
		f.handlerClownAnchovy(w)
		return
	}

	f.handlerTuna(w, r)
}
