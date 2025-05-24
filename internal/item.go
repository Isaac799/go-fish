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

// Item is an item found form the template dir.
type Item struct {
	kind     int
	Pattern  string
	filePath string
	children []Item
}

func isIsland(s string) bool {
	return strings.HasPrefix(s, "_")
}

func newItem(entry os.DirEntry, pathBase, templateDir string) (*Item, error) {
	info, err := entry.Info()
	if err != nil {
		return nil, err
	}

	kindMap := map[string]int{
		".html": htmlItemKindPage,
		".css":  htmlItemKindStyle,
	}

	ext := filepath.Ext(info.Name())

	kind, exists := kindMap[ext]
	if !exists {
		return nil, ErrInvalidExtension
	}

	if isIsland(info.Name()) {
		kind = htmlItemKindIsland
	}

	name := info.Name()
	if kind == htmlItemKindPage {
		name = strings.TrimSuffix(info.Name(), ext)
	}

	// since I want to cache styling while preventing
	// an invalid cache we make the name based on a hash
	// of its content
	if kind == htmlItemKindStyle {
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

	return &Item{
		kind:     kind,
		Pattern:  pattern,
		filePath: filePath,
	}, nil

}

func (hi *Item) templateName() string {
	if len(hi.Pattern) == 0 {
		return "unknown"
	}
	parts := strings.Split(hi.Pattern, "/")
	return parts[len(parts)-1]
}

func (hi *Item) handler(w http.ResponseWriter, _ *http.Request) {

	if hi.kind == htmlItemKindIsland {
		return
	}

	if hi.kind == htmlItemKindStyle {
		f, err := os.Open(hi.filePath)
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
		defer f.Close()
		b, err := io.ReadAll(f)
		if err != nil {
			fmt.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", mime.TypeByExtension(filepath.Ext(f.Name())))
		w.Header().Add("Content-Length", strconv.Itoa(len(b)))
		w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", browserCacheDurationSeconds))
		w.Write(b)
		return
	}

	// 3 main parts to the document I will add in between
	htmlStartHead := []byte(`<!DOCTYPE html><html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0" >`)
	htmlEndHeadStartBody := []byte(`
</head>
<body>`)
	htmlEndBody := []byte(`</body></html>`)

	// gather islands
	collectedFilePaths := []string{}
	for _, e := range hi.children {
		if e.kind != htmlItemKindIsland {
			continue
		}
		collectedFilePaths = append(collectedFilePaths, e.filePath)
	}
	collectedFilePaths = append(collectedFilePaths, hi.filePath)

	templateName := hi.templateName()

	t := template.New(templateName)
	parsed, err := t.ParseFiles(collectedFilePaths...)

	// buffer for html doc
	b := []byte{}
	buff := bytes.NewBuffer(b)
	buff.Write(htmlStartHead)

	// styling
	for _, e := range hi.children {
		if e.kind != htmlItemKindStyle {
			continue
		}
		b := fmt.Appendf(nil, `<link rel="stylesheet" href="%s">`, e.Pattern)
		buff.Write(b)
	}
	buff.Write(htmlEndHeadStartBody)

	err = parsed.ExecuteTemplate(buff, templateName, hi)
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
