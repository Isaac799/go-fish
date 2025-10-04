package main

import (
	"net/http"
)

type dragDropItem struct {
	ID,
	X,
	Y int
}

type user struct {
	ID        int
	FirstName string
	LastName  string
}

func queriedSeason(r *http.Request) any {
	season, ok := r.Context().Value(queryCtxKey).(string)
	if !ok {
		return nil
	}
	return &season
}

func userInfo(r *http.Request) any {
	user, ok := r.Context().Value(userCtxKey).(user)
	if !ok {
		return nil
	}
	return &user
}
