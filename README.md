# go-fish

a simple ssr framework made with go templating and ws for regional reloading

## Concepts

- **page**: is available via mux
- **island**: a template to be used within other templates, prefixed with "_"

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
	mux, err := gofish.NewMux("template")
	if err != nil {
		panic(err)
	}
	fmt.Println("gone fishing")
	http.ListenAndServe(":8080", mux)
}
```

Template Directory

```txt
└── template
    ├── about.html             /about
    ├── blog
    │   ├── 2006-01-02.html    /blog/2006-01-03
    │   ├── 2006-01-03.html    /blog/2006-01-03
    │   └── _greeting.html     'greeting' is scoped to blogs
    ├── home.html              /home
    └── _nav.html              'nav' is globally scoped 
```

Template HTML define same as file name

```html
<!-- about.html -->
{{define "about"}}

<p>Hello World </p>

{{end}}
```