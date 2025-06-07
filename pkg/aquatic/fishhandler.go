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
			b, err := templateBytes(&e)
			if err != nil {
				fmt.Print(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			allFishBuff.Grow(len(b))
			allFishBuff.Write(b)
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
			b, err := templateBytes(e)
			if err != nil {
				fmt.Print(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			allFishBuff.Grow(len(b))
			allFishBuff.Write(b)
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
			b, err := templateBytes(&e)
			if err != nil {
				fmt.Print(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			allFishBuff.Grow(len(b))
			allFishBuff.Write(b)
			sardinesEaten[e.templateName] = true
		}

		// Also just the mackerel, aka the system fish
		for _, e := range pond.globalSmallFish {
			if e.kind != FishKindMackerel {
				continue
			}
			if _, exists := sardinesEaten[e.templateName]; exists {
				continue
			}
			allFishBuff.Grow(len(e.bytes))
			allFishBuff.Write(e.bytes)
			sardinesEaten[e.templateName] = true
		}

		if _, exists := sardinesEaten[f.templateName]; !exists {
			// Define the tuna 'sardine' last
			b, err := templateBytes(f)
			if err != nil {
				fmt.Print(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			allFishBuff.Grow(len(b))
			allFishBuff.Write(b)
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
			b, err := templateBytes(&e)
			if err != nil {
				fmt.Print(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			allFishBuff.Grow(len(b))
			allFishBuff.Write(b)
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
			b, err := templateBytes(e)
			if err != nil {
				fmt.Print(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			allFishBuff.Grow(len(b))
			allFishBuff.Write(b)
			sardinesEaten[e.templateName] = true
		}

		// Also just the mackerel, aka the system fish
		for _, e := range pond.globalSmallFish {
			if e.kind != FishKindMackerel {
				continue
			}
			if _, exists := sardinesEaten[e.templateName]; exists {
				continue
			}
			allFishBuff.Grow(len(e.bytes))
			allFishBuff.Write(e.bytes)
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
			b, err := templateBytes(&e)
			if err != nil {
				fmt.Print(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			allFishBuff.Grow(len(b))
			allFishBuff.Write(b)
			sardinesEaten[e.templateName] = true
		}

		if _, exists := sardinesEaten[f.templateName]; exists {
			fmt.Println("sardine name conflicts with tuna name: ", f.templateName)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Define the tuna last
		b, err := templateBytes(f)
		if err != nil {
			fmt.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		allFishBuff.Grow(len(b))
		allFishBuff.Write(b)

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

	var (
		cannotCatch = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("this fish is for catching"))
		})
		unaccountedFish = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("this fish was not accounted for"))
		})
	)

	handlerMap := map[int]http.HandlerFunc{
		FishKindSardine:  handlerSardine(f, pond),
		FiskKindClown:    handlerClownAnchovy[T](f),
		FiskKindAnchovy:  handlerClownAnchovy[T](f),
		FishKindTuna:     handlerTuna(f, pond),
		FishKindMackerel: cannotCatch,
	}

	finalHandler, exists := handlerMap[f.kind]
	if !exists {
		return unaccountedFish
	}

	return chainLicenses(finalHandler, licenses...)
}
