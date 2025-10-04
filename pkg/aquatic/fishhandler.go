package aquatic

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

type pageData struct {
	Pond any
	Fish any
}

func handlerSardine(f *Fish, p *Pond) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := template.New(f.templateName)

		buff, err := f.cacheReef(p)
		if err != nil {
			fmt.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if f.Tackle != nil {
			t.Funcs(f.Tackle)
		}

		parsed, err := t.Parse(string(buff))
		if err != nil {
			fmt.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var fishData any
		if f.OnCatch != nil {
			fishData = f.OnCatch(r)
		}
		var pondData any
		if p.OnCatch != nil {
			pondData = p.OnCatch(r)
		}

		data := pageData{Pond: pondData, Fish: fishData}

		// want to exe template into this to get len for res
		resBytes := []byte{}
		resBuff := bytes.NewBuffer(resBytes)

		err = parsed.ExecuteTemplate(resBuff, f.templateName, data)
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

func handlerClownAnchovy(f *Fish) http.HandlerFunc {
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

func bobber(f *Fish, pond *Pond) []byte {
	if f.bobber != nil {
		return f.bobber
	}

	// unlikely more than 10 links in doc head
	// so realloc at least that many
	headLinks := make([][]byte, 0, 10)

	size := 0
	for _, e := range pond.shad {
		if e.kind != FiskKindClown {
			continue
		}
		if strings.HasPrefix(e.mime, "text/css") {
			b := fmt.Appendf(nil, `<link rel="stylesheet" href="%s?v=%s">`, e.pattern, e.hash)
			headLinks = append(headLinks, b)
			size += len(b)
		}
		if strings.HasPrefix(e.mime, "text/javascript") {
			b := fmt.Appendf(nil, `<script src="%s?v=%s"></script>`, e.pattern, e.hash)
			headLinks = append(headLinks, b)
			size += len(b)
		}
	}
	for _, e := range f.school {
		if e.kind != FiskKindClown {
			continue
		}
		if strings.HasPrefix(e.mime, "text/css") {
			b := fmt.Appendf(nil, `<link rel="stylesheet" href="%s?v=%s">`, e.pattern, e.hash)
			headLinks = append(headLinks, b)
			size += len(b)
		}
		if strings.HasPrefix(e.mime, "text/javascript") {
			b := fmt.Appendf(nil, `<script src="%s?v=%s"></script>`, e.pattern, e.hash)
			headLinks = append(headLinks, b)
			size += len(b)
		}
	}

	// Sort ensure links in lexicographical order (alphabetical)
	// Important for consistency in resolving css class conflicts and such
	sort.Slice(headLinks, func(i, j int) bool {
		return bytes.Compare(headLinks[i], headLinks[j]) < 0
	})

	b := make([]byte, size)
	last := 0
	for _, v := range headLinks {
		n := copy(b[last:last+len(v)], v)
		last += n
	}
	f.bobber = b
	return b
}

// handlerTuna wraps a fish reef in in html5 syntax and
// adds the bobber to the head of the document. It uses the
// bait for a fish and its pond for template data. It uses
// tackle for template funcs.
func handlerTuna(f *Fish, p *Pond) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			docStart  = []byte(`<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0" >`)
			bodyStart = []byte(`</head><body>`)
			docEnd    = []byte(`</body></html>`)
		)

		t := template.New(f.templateName)

		reef, err := f.cacheReef(p)
		if err != nil {
			fmt.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if f.Tackle != nil {
			t.Funcs(f.Tackle)
		}

		parsed, err := t.Parse(string(reef))
		if err != nil {
			fmt.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		headLinks := bobber(f, p)

		// this size is not perfect since the executed template size
		// cannot be know, but it helps some allocation before that
		size := len(docStart) + len(headLinks) + len(bodyStart) + len(docEnd)

		buff := bytes.NewBuffer(make([]byte, 0, size))
		buff.Write(docStart)
		buff.Write(headLinks)
		buff.Write(bodyStart)

		var fishData any
		if f.OnCatch != nil {
			fishData = f.OnCatch(r)
		}
		var pondData any
		if p.OnCatch != nil {
			pondData = p.OnCatch(r)
		}

		data := pageData{Pond: pondData, Fish: fishData}

		err = parsed.ExecuteTemplate(buff, f.templateName, data)
		if err != nil {
			fmt.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		buff.Write(docEnd)

		w.Header().Add("Cache-Control", "no-store")
		w.Header().Add("Content-Type", "text/html")
		w.Header().Add("Content-Length", strconv.Itoa(len(buff.Bytes())))
		_, err = w.Write(buff.Bytes())
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
func reel(f *Fish, pond *Pond) http.Handler {
	licenses := []License{}

	for _, license := range pond.licenses {
		licenses = append(licenses, license)
	}

	for _, license := range f.Licenses {
		licenses = append(licenses, license)
	}

	var (
		unaccountedFish = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("this fish was not accounted for"))
		})
	)

	handlerMap := map[int]http.HandlerFunc{
		FishKindSardine: handlerSardine(f, pond),
		FiskKindClown:   handlerClownAnchovy(f),
		FiskKindAnchovy: handlerClownAnchovy(f),
		FishKindTuna:    handlerTuna(f, pond),
	}

	finalHandler, exists := handlerMap[f.kind]
	if !exists {
		return unaccountedFish
	}

	return chainLicenses(finalHandler, licenses...)
}
