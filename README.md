# Go Fish

a fish themed ssr framework using go templating and htmlx

## Concepts

The main 3 things are:

- **Tuna**: is the big fish of the app, a page that people visit. Available via mux
  - consumes sardines
- **Sardine**: is the small fry, prefixed with "_"
  - a template to be used within other templates 
  - can be fetched with htmlx
- **Clown**fish  is the styling or css of the app. 
  - hashed for their name to enable browser cache
  - served independently, not embedded in template

Sardine and Clown scope is show in the example below.

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
    │   ├── blog-style.css	    /blog/dd55ea4c2e29.css
    │   └── _greeting.html      'greeting' can be used in a blogs
    ├── home.html               /home
    ├── _nav.html               'nav' can be used on all pages
    └── style.css               /7d084444a097620c4.css
```

Template HTML define same as file name

```html
<!-- about.html -->
{{define "about"}}

<p>Hello World </p>

{{end}}
```

Using htmlx to load a sardine

```html
<div id="_nav"></div>
<div>
    <button
        hx-get="/_nav"
        hx-target="#_nav"
    >
        Load Nav
    </button>
</div>
```