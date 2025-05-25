package main

import (
	"fmt"
	"net/http"
	"time"
)

func springOnly(w http.ResponseWriter, r *http.Request) bool {
	if !r.URL.Query().Has("season") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("a season is required"))
		return false
	}

	season := r.URL.Query().Get("season")
	if season != "spring" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("cannot fish outside of spring"))
		return false
	}

	return true
}

func visitorLog(_ http.ResponseWriter, _ *http.Request) bool {
	fmt.Printf("visited %s\n", time.Now().Format(time.RFC1123))
	return true
}
