# go-fish

a ssr framework using go templating and htmlx

## Concepts

The main two things are:

- **page**: is available via mux
- **island**: a template to be used within other templates, prefixed with "_"

Scoping for islands is show in the example below.

## Styling

Styling scope is show in the example below too. 

I wanted to note here that we hash the content to make a new file name, once serving. This allows us to cache the document in the browser for faster loading times and less server overhead. Especially compared to loading css into the template every time. Also this lets me make changes without fear that a user will miss out.

## Example

Primary entry point

```go
package main

import (
	"fmt"
	"net/http"

	gofish "github.com/Isaac799/go-fish/internal"
)

func main() {
	pond, err := gofish.NewPond("template")
	if err != nil {
		panic(err)
	}

	verbose := true
	mux := pond.CastLines(verbose)

	fmt.Println("gone fishing")
	http.ListenAndServe(":8080", mux)
}

```

Template Directory

```txt
└── template
    ├── about.html              /about
    ├── blog
    │   ├── 2006-01-02.html     /blog/2006-01-02
    │   ├── 2006-01-03.html     /blog/2006-01-03
    │   ├── blog-style.css	    /blog/dd55ea4c2e29831e355a68015bc12d00.css
    │   └── _greeting.html      'greeting' can be used in a blogs
    ├── home.html               /home
    ├── _nav.html               'nav' can be used on all pages
    └── style.css               /7d084444a097620c49fe94852b215eb2.css
```

Template HTML define same as file name

```html
<!-- about.html -->
{{define "about"}}

<p>Hello World </p>

{{end}}
```