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
type License = func(next http.Handler) http.Handler

// Bait is to be gobbled up by a fish before catching it.
// A func that has access to the request and returns template data
type Bait = func(r *http.Request) any

// Fish is an item found form the template dir.
type Fish struct {
	pond         *Pond
	kind         int
	mime         string
	templateName string
	pattern      string
	filePath     string
	children     []Fish

	// Licenses is a collection of licenses a user must have
	// to catch a fish. Checked after pond licenses, in the
	// order added. To catch a fish all pond and fish licenses
	// must be met.
	Licenses []License

	// Bait fn is called and the result is passed into the
	// executed template, or eaten by the fish before caught
	Bait Bait
}

// AddLicense appends a license required to catch a fish
func (f *Fish) AddLicense(l License) {
	f.Licenses = append(f.Licenses, l)
}

// Gobble has one fish gobble up another. Gaining its Licenses and Bait.
func (f *Fish) Gobble(f2 Fish) {
	f.Licenses = f2.Licenses
	f.Bait = f2.Bait
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
	templateName := name

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
	pattern = strings.ReplaceAll(pattern, " ", "-")

	pattern = strings.ToLower(pattern)

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

	fmt.Println(pattern)

	return &Fish{
		kind:         kind,
		mime:         mime,
		pattern:      pattern,
		templateName: templateName,
		filePath:     filePath,
		Licenses:     []License{},
		pond:         pond,
	}, nil
}

// TemplateBuffer will wrap a file content in the define syntax.
// Enforcing template name scheme and reducing template lines n - 2.
func (f *Fish) templateBuffer() (*bytes.Buffer, error) {
	file, err := os.Open(f.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// buffer for wrap file in define
	fileBuffer := bytes.NewBuffer([]byte{})

	start := fmt.Appendf(nil, "{{define \"%s\"}}", f.templateName)
	fileBuffer.Grow(len(start))
	fileBuffer.Write(start)

	fileBuffer.Grow(len(fileBytes))
	fileBuffer.Write(fileBytes)

	end := []byte("{{end}}")
	fileBuffer.Grow(len(end))
	fileBuffer.Write(end)

	return fileBuffer, nil
}

func (f *Fish) handlerSardine(w http.ResponseWriter, r *http.Request) {
	t := template.New(f.templateName)

	buff, err := f.templateBuffer()
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	parsed, err := t.Parse(buff.String())

	var pageData any
	if f.Bait != nil {
		pageData = f.Bait(r)
	}

	// want to exe template into this to get len for res
	resBytes := []byte{}
	resBuff := bytes.NewBuffer(resBytes)

	err = parsed.ExecuteTemplate(resBuff, f.templateName, pageData)
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/html")
	w.Header().Add("Content-Length", strconv.Itoa(len(resBuff.Bytes())))
	w.Write(resBuff.Bytes())
}

func (f *Fish) handlerClownAnchovy(w http.ResponseWriter, _ *http.Request) {
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
	w.Header().Add("Content-Type", f.mime)
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

	t := template.New(f.templateName)

	allFishBuff := bytes.NewBuffer([]byte{})

	// Define the sardines first
	for _, e := range f.children {
		if e.kind != FishKindSardine {
			continue
		}
		buff, err := e.templateBuffer()
		if err != nil {
			fmt.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		allFishBuff.Grow(buff.Len())
		allFishBuff.Write(buff.Bytes())
	}

	// Define the tuna last
	buff, err := f.templateBuffer()
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	allFishBuff.Grow(buff.Len())
	allFishBuff.Write(buff.Bytes())

	parsed, err := t.Parse(allFishBuff.String())

	// buffer for html doc we will wrap in html5 syntax
	// want to exec template into this to get len for res
	resBytes := []byte{}
	resBuff := bytes.NewBuffer(resBytes)

	resBuff.Grow(len(htmlStartHead))
	resBuff.Write(htmlStartHead)

	// styling
	for _, e := range f.children {
		if e.kind != FiskKindClown {
			continue
		}
		if strings.HasPrefix(e.mime, "text/css") {
			b := fmt.Appendf(nil, `<link rel="stylesheet" href="%s">`, e.pattern)
			resBuff.Grow(len(b))
			resBuff.Write(b)
		}
		if strings.HasPrefix(e.mime, "text/javascript") {
			b := fmt.Appendf(nil, `<script src="%s"></script>`, e.pattern)
			resBuff.Grow(len(b))
			resBuff.Write(b)
		}
	}
	resBuff.Grow(len(htmlEndHeadStartBody))
	resBuff.Write(htmlEndHeadStartBody)

	var pageData any
	if f.Bait != nil {
		pageData = f.Bait(r)
	}

	err = parsed.ExecuteTemplate(resBuff, f.templateName, pageData)
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resBuff.Grow(len(htmlEndBody))
	resBuff.Write(htmlEndBody)

	w.Header().Add("Content-Type", "text/html")
	w.Header().Add("Content-Length", strconv.Itoa(len(resBuff.Bytes())))
	w.Write(resBuff.Bytes())
}

// chainLicense essentially is a Russian nesting doll like so
// (fin, A, B) is ran as A(B(fin))
// get it? fish have fin. do you get it?
func chainLicenses(fin http.Handler, licenses ...License) http.Handler {
	for i := len(licenses) - 1; i >= 0; i-- {
		license := licenses[i]
		fin = license(fin)
	}
	return fin
}

// reel enables catching a fish. It will chain license
// together to ensure you are allowed to catch
func (f *Fish) reel() http.Handler {
	licenses := []License{}

	for _, license := range f.pond.licenses {
		licenses = append(licenses, license)
	}

	for _, license := range f.Licenses {
		licenses = append(licenses, license)
	}

	handlerMap := map[int]http.HandlerFunc{
		FishKindSardine: f.handlerSardine,
		FiskKindClown:   f.handlerClownAnchovy,
		FiskKindAnchovy: f.handlerClownAnchovy,
		FishKindTuna:    f.handlerTuna,
	}

	finalHandler, exists := handlerMap[f.kind]
	if !exists {
		panic("no valid handler for fish kind")
	}

	return chainLicenses(finalHandler, licenses...)
}
