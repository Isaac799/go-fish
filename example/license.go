package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type contextKey string

const seasonCtxKey contextKey = "season"

func springOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		season, ok := r.Context().Value(seasonCtxKey).(string)

		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("season must be valid string"))
			return
		}

		if season != "spring" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("cannot fish outside of spring"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func requireSeason(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !r.URL.Query().Has("season") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("a season is required"))
			return
		}

		season := r.URL.Query().Get("season")
		if len(season) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("a season is required"))
			return
		}

		ctx := context.WithValue(r.Context(), seasonCtxKey, season)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func visitorLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("visited %s\n", time.Now().Format(time.RFC1123))
		next.ServeHTTP(w, r)
	})
}
