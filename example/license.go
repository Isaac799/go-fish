package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type contextKey string

const seasonCtxKey contextKey = "season"
const userCtxKey contextKey = "user"
const queryCtxKey contextKey = "q"
const rotateCtxKey contextKey = "rotate"

const (
	right = "right"
	left  = "left"
)

var userDB = map[int]user{
	1: {
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
	},
	2: {
		ID:        2,
		FirstName: "Jane",
		LastName:  "Doe",
	},
	3: {
		ID:        3,
		FirstName: "Sally",
		LastName:  "Sue",
	},
}

func requireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if len(id) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("a user id is required"))
			return
		}
		i, err := strconv.Atoi(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("user id must be an integer"))
			return
		}

		u, exists := userDB[i]
		if !exists {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("user not found"))
			return
		}

		ctx := context.WithValue(r.Context(), userCtxKey, u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func optionQuery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		ctx := context.WithValue(r.Context(), queryCtxKey, q)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func visitorLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("visited %s\n", time.Now().Format(time.RFC1123))
		next.ServeHTTP(w, r)
	})
}
