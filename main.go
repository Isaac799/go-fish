package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func main() {
	mux := http.NewServeMux()

	items := map[string][]string{}
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	templateDir := filepath.Join(wd, "template")
	err = gatherTemplates(&items, templateDir)
	if err != nil {
		panic(err)
	}
	for path, files := range items {
		if len(files) == 0 {
			fmt.Printf("no patterns for: %s\n", path)
			continue
		}
		for _, fileName := range files {
			ext := filepath.Ext(fileName)
			name := strings.TrimSuffix(fileName, ext)

			filePath := filepath.Join(path, fileName)
			pattern := filepath.Join(path, name)
			pattern = strings.Replace(pattern, templateDir, "", 1)
			pattern = strings.ReplaceAll(pattern, "\\", "/")
			pattern = strings.ReplaceAll(pattern, "//", "/")
			fmt.Printf("pattern: %s\n", pattern)
			mux.HandleFunc(pattern, simpleServe(filePath, pattern))
		}
	}
	fmt.Println("gone fishing")
	http.ListenAndServe(":8080", mux)
}

func gatherTemplates(items *map[string][]string, base string) error {
	if items == nil {
		items = &map[string][]string{}
	}
	entries, err := os.ReadDir(base)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrNoTemplateDir
		}
		return err
	}

	for _, e := range entries {
		if e.IsDir() {
			gatherTemplates(items, filepath.Join(base, e.Name()))
			continue
		}
		itemsDeref := *items
		_, exists := itemsDeref[base]
		if !exists {
			itemsDeref[base] = []string{}
		}
		info, err := e.Info()
		if err != nil {
			return err
		}
		ext := filepath.Ext(info.Name())
		if ext != ".html" {
			continue
		}
		itemsDeref[base] = append(itemsDeref[base], info.Name())
	}

	return nil
}

func simpleServe(filePath string, templateName string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		t := template.New(templateName)
		parsed, err := t.ParseFiles(filePath)
		if err != nil {
			fmt.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = parsed.ExecuteTemplate(w, templateName, nil)
		if err != nil {
			fmt.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
