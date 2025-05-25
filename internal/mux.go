package internal

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"text/tabwriter"
)

func htmlxHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/javascript")
	w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", browserCacheDurationSeconds))
	w.Header().Add("Content-Length", strconv.Itoa(len(htmlx)))
	w.Write(htmlx)
}

// NewMux provides a mux to with patterns based on go templates in the specified directory
func NewMux(templateDirPath string, verbose bool) (*http.ServeMux, error) {
	mux := http.NewServeMux()

	mux.HandleFunc("/assets/htmlx.2.0.4.js", htmlxHandler)

	items := map[string][]Item{}
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	templateDir := filepath.Join(wd, templateDirPath)
	err = collect(&items, templateDir, templateDir, nil)
	if err != nil {
		return nil, err
	}

	// prevents duplicate pattern registration
	// expected since children share stylesheets
	pattensAdded := map[string]bool{}

	var tw *tabwriter.Writer

	if verbose {
		tw = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		tw.Write([]byte("kind\tpattern\n"))
		tw.Write([]byte("--\t--\n"))
	}

	for path, items := range items {
		if len(items) == 0 {
			fmt.Printf("no patterns for: %s\n", path)
			continue
		}

		for _, item := range items {
			if item.kind != htmlItemKindPage {
				continue
			}

			if _, exists := pattensAdded[item.Pattern]; exists {
				continue
			}
			if tw != nil {
				tw.Write(fmt.Appendf(nil, "page\t%s\n", item.Pattern))
			}

			mux.HandleFunc(item.Pattern, item.handler)
			pattensAdded[item.Pattern] = true

			for _, child := range item.children {
				if child.kind != htmlItemKindStyle {
					continue
				}
				if _, exists := pattensAdded[child.Pattern]; exists {
					continue
				}
				if tw != nil {
					tw.Write(fmt.Appendf(nil, "style\t%s\n", child.Pattern))
				}
				mux.HandleFunc(child.Pattern, child.handler)
				pattensAdded[child.Pattern] = true
			}
		}
	}
	if tw != nil {
		tw.Flush()
	}

	return mux, nil
}
