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
