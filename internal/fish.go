package internal

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"maps"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

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

// AddLicense appends a license required to catch a fish
func addLicense[K any](f *Fish[K], l License) {
	f.Licenses = append(f.Licenses, l)
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
func templateBuffer[K any](f *Fish[K]) (*bytes.Buffer, error) {
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
	_, err = fileBuffer.Write(start)
	if err != nil {
		return fileBuffer, err
	}

	fileBuffer.Grow(len(fileBytes))
	_, err = fileBuffer.Write(fileBytes)
	if err != nil {
		return fileBuffer, err
	}

	end := []byte("{{end}}")
	fileBuffer.Grow(len(end))
	_, err = fileBuffer.Write(end)
	if err != nil {
		return fileBuffer, err
	}

	return fileBuffer, nil
}

func handlerSardine[T, K any](f *Fish[K], pond *Pond[T, K]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := template.New(f.templateName)

		allFishBuff := bytes.NewBuffer([]byte{})
		sardinesEaten := map[string]bool{}

		// Define the local sardines first
		for _, e := range f.children {
			if e.kind != FishKindSardine {
				continue
			}
			if _, exists := sardinesEaten[e.templateName]; exists {
				continue
			}
			buff, err := templateBuffer(&e)
			if err != nil {
				fmt.Print(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			allFishBuff.Grow(buff.Len())
			allFishBuff.Write(buff.Bytes())
			sardinesEaten[e.templateName] = true
		}

		// Define the global sardines second, not to overwrite the local ones
		for _, e := range pond.globalSmallFish {
			if e.kind != FishKindSardine {
				continue
			}
			if _, exists := sardinesEaten[e.templateName]; exists {
				continue
			}
			buff, err := templateBuffer(e)
			if err != nil {
				fmt.Print(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			allFishBuff.Grow(buff.Len())
			allFishBuff.Write(buff.Bytes())
			sardinesEaten[e.templateName] = true
		}

		// Define the children last
		for _, e := range f.children {
			if e.kind != FishKindSardine {
				continue
			}
			if _, exists := sardinesEaten[e.templateName]; exists {
				continue
			}
			buff, err := templateBuffer(&e)
			if err != nil {
				fmt.Print(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			allFishBuff.Grow(buff.Len())
			allFishBuff.Write(buff.Bytes())
			sardinesEaten[e.templateName] = true
		}

		if _, exists := sardinesEaten[f.templateName]; !exists {
			// Define the tuna 'sardine' last
			buff, err := templateBuffer(f)
			if err != nil {
				fmt.Print(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			allFishBuff.Grow(buff.Len())
			allFishBuff.Write(buff.Bytes())
		}

		if f.Tackle != nil {
			t.Funcs(f.Tackle)
		}

		parsed, err := t.Parse(allFishBuff.String())
		if err != nil {
			fmt.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var globalBait T
		var localBait K

		if pond.GlobalBait != nil {
			globalBait = pond.GlobalBait(r)
		}
		if f.Bait != nil {
			localBait = f.Bait(r)
		}

		pageData := masterBait[T, K]{
			Local:  localBait,
			Global: globalBait,
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
		_, err = w.Write(resBuff.Bytes())
		if err != nil {
			fmt.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func handlerClownAnchovy[T, K any](f *Fish[K]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if f.kind == FiskKindClown {
			ver := r.URL.Query().Get("v")
			if ver != f.hash {
				fmt.Print("mismatched file ver request")
				w.WriteHeader(http.StatusNotFound)
				return
			}
		}

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
		_, err = w.Write(b)
		if err != nil {
			fmt.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func handlerTuna[T, K any](f *Fish[K], pond *Pond[T, K]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 3 main parts to the document I will add in between
		htmlStartHead := []byte(`<!DOCTYPE html><html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0" >`)
		htmlEndHeadStartBody := []byte(`
</head>
<body>`)
		htmlEndBody := []byte(`</body></html>`)

		t := template.New(f.templateName)

		allFishBuff := bytes.NewBuffer([]byte{})

		sardinesEaten := map[string]bool{}

		// Define the local sardines first
		for _, e := range f.children {
			if e.kind != FishKindSardine {
				continue
			}
			if _, exists := sardinesEaten[e.templateName]; exists {
				continue
			}
			buff, err := templateBuffer(&e)
			if err != nil {
				fmt.Print(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			allFishBuff.Grow(buff.Len())
			allFishBuff.Write(buff.Bytes())
			sardinesEaten[e.templateName] = true
		}

		// Define the global sardines first
		for _, e := range pond.globalSmallFish {
			if e.kind != FishKindSardine {
				continue
			}
			if _, exists := sardinesEaten[e.templateName]; exists {
				continue
			}
			buff, err := templateBuffer(e)
			if err != nil {
				fmt.Print(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			allFishBuff.Grow(buff.Len())
			allFishBuff.Write(buff.Bytes())
			sardinesEaten[e.templateName] = true
		}

		// Define the global sardines second, not to overwrite the local ones
		for _, e := range f.children {
			if e.kind != FishKindSardine {
				continue
			}
			if _, exists := sardinesEaten[e.templateName]; exists {
				continue
			}
			buff, err := templateBuffer(&e)
			if err != nil {
				fmt.Print(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			allFishBuff.Grow(buff.Len())
			allFishBuff.Write(buff.Bytes())
			sardinesEaten[e.templateName] = true
		}

		if _, exists := sardinesEaten[f.templateName]; exists {
			fmt.Println("sardine name conflicts with tuna name: ", f.templateName)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Define the tuna last
		buff, err := templateBuffer(f)
		if err != nil {
			fmt.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		allFishBuff.Grow(buff.Len())
		allFishBuff.Write(buff.Bytes())

		if f.Tackle != nil {
			t.Funcs(f.Tackle)
		}

		parsed, err := t.Parse(allFishBuff.String())
		if err != nil {
			fmt.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// buffer for html doc we will wrap in html5 syntax
		// want to exec template into this to get len for res
		resBytes := []byte{}
		resBuff := bytes.NewBuffer(resBytes)

		resBuff.Grow(len(htmlStartHead))
		_, err = resBuff.Write(htmlStartHead)
		if err != nil {
			fmt.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Ensure the links are always ordered the same
		// so the page loads them the same. Important for
		// css class conflicts and such
		orderedSmallFish := make([]*Fish[K], len(pond.globalSmallFish))

		// global styling
		k := 0
		for _, e := range pond.globalSmallFish {
			orderedSmallFish[k] = e
			k++
		}

		sort.Slice(orderedSmallFish, func(i, j int) bool {
			return strings.Compare(orderedSmallFish[i].pattern, orderedSmallFish[j].pattern) < 0
		})

		for _, e := range orderedSmallFish {
			if e.kind != FiskKindClown {
				continue
			}
			if strings.HasPrefix(e.mime, "text/css") {
				b := fmt.Appendf(nil, `<link rel="stylesheet" href="%s?v=%s">`, e.pattern, e.hash)
				resBuff.Grow(len(b))
				_, err = resBuff.Write(b)
				if err != nil {
					fmt.Print(err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
			if strings.HasPrefix(e.mime, "text/javascript") {
				b := fmt.Appendf(nil, `<script src="%s?v=%s"></script>`, e.pattern, e.hash)
				resBuff.Grow(len(b))
				_, err = resBuff.Write(b)
				if err != nil {
					fmt.Print(err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
		}

		// styling
		for _, e := range f.children {
			if e.kind != FiskKindClown {
				continue
			}
			if strings.HasPrefix(e.mime, "text/css") {
				b := fmt.Appendf(nil, `<link rel="stylesheet" href="%s">`, e.pattern)
				resBuff.Grow(len(b))
				_, err = resBuff.Write(b)
				if err != nil {
					fmt.Print(err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
			if strings.HasPrefix(e.mime, "text/javascript") {
				b := fmt.Appendf(nil, `<script src="%s"></script>`, e.pattern)
				resBuff.Grow(len(b))
				_, err = resBuff.Write(b)
				if err != nil {
					fmt.Print(err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
		}
		resBuff.Grow(len(htmlEndHeadStartBody))
		_, err = resBuff.Write(htmlEndHeadStartBody)
		if err != nil {
			fmt.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var globalBait T
		var localBait K
		if f.Bait != nil {
			localBait = f.Bait(r)
		}
		if pond.GlobalBait != nil {
			globalBait = pond.GlobalBait(r)
		}

		pageData := masterBait[T, K]{
			Local:  localBait,
			Global: globalBait,
		}

		err = parsed.ExecuteTemplate(resBuff, f.templateName, pageData)
		if err != nil {
			fmt.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resBuff.Grow(len(htmlEndBody))
		resBuff.Write(htmlEndBody)

		w.Header().Add("Cache-Control", "no-store")
		w.Header().Add("Content-Type", "text/html")
		w.Header().Add("Content-Length", strconv.Itoa(len(resBuff.Bytes())))
		_, err = w.Write(resBuff.Bytes())
		if err != nil {
			fmt.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
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
func reel[T, K any](f *Fish[K], pond *Pond[T, K]) http.Handler {
	licenses := []License{}

	for _, license := range pond.licenses {
		licenses = append(licenses, license)
	}

	for _, license := range f.Licenses {
		licenses = append(licenses, license)
	}

	handlerMap := map[int]http.HandlerFunc{
		FishKindSardine: handlerSardine(f, pond),
		FiskKindClown:   handlerClownAnchovy[T](f),
		FiskKindAnchovy: handlerClownAnchovy[T](f),
		FishKindTuna:    handlerTuna(f, pond),
	}

	finalHandler, exists := handlerMap[f.kind]
	if !exists {
		panic("no valid handler for fish kind")
	}

	return chainLicenses(finalHandler, licenses...)
}
