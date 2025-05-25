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

// Fish is an item found form the template dir.
type Fish struct {
	kind     int
	Pattern  string
	filePath string
	children []Fish
}

func isSardine(s string) bool {
	return strings.HasPrefix(s, "_")
}

func newFish(entry os.DirEntry, pathBase, templateDir string) (*Fish, error) {
	info, err := entry.Info()
	if err != nil {
		return nil, err
	}

	kindMap := map[string]int{
		".html": fishKindTuna,
		".css":  fiskKindClown,
	}

	ext := filepath.Ext(info.Name())

	kind, exists := kindMap[ext]
	if !exists {
		return nil, ErrInvalidExtension
	}

	if isSardine(info.Name()) {
		kind = fishKindSardine
	}

	name := info.Name()
	if kind == fishKindTuna || kind == fishKindSardine {
		name = strings.TrimSuffix(info.Name(), ext)
	}

	// since I want to cache styling while preventing
	// an invalid cache we make the name based on a hash
	// of its content
	if kind == fiskKindClown {
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
	pattern = strings.Replace(pattern, templateDir, "", 1)
	pattern = strings.ReplaceAll(pattern, "\\", "/")
	pattern = strings.ReplaceAll(pattern, "//", "/")

	return &Fish{
		kind:     kind,
		Pattern:  pattern,
		filePath: filePath,
	}, nil

}

func (f *Fish) templateName() string {
	if len(f.Pattern) == 0 {
		return "unknown"
	}
	parts := strings.Split(f.Pattern, "/")
	return parts[len(parts)-1]
}

func (f *Fish) handlerIsland(w http.ResponseWriter) {
	templateName := f.templateName()

	t := template.New(templateName)
	parsed, err := t.ParseFiles(f.filePath)

	// buffer for html doc
	b := []byte{}
	buff := bytes.NewBuffer(b)
	err = parsed.ExecuteTemplate(buff, templateName, f)
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/html")
	w.Header().Add("Content-Length", strconv.Itoa(len(buff.Bytes())))
	w.Write(buff.Bytes())
}

func (f *Fish) handlerCSS(w http.ResponseWriter) {
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

func (f *Fish) handlerHTMLPage(w http.ResponseWriter) {
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
		if e.kind != fishKindSardine {
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
		if e.kind != fiskKindClown {
			continue
		}
		b := fmt.Appendf(nil, `<link rel="stylesheet" href="%s">`, e.Pattern)
		buff.Write(b)
	}
	buff.Write(htmlEndHeadStartBody)

	err = parsed.ExecuteTemplate(buff, templateName, f)
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

func (f *Fish) handler(w http.ResponseWriter, _ *http.Request) {
	if f.kind == fishKindSardine {
		f.handlerIsland(w)
		return
	}

	if f.kind == fiskKindClown {
		f.handlerCSS(w)
		return
	}

	f.handlerHTMLPage(w)
}
