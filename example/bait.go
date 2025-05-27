package main

import (
	"net/http"
	"strconv"
	"time"
)

type user struct {
	ID        int
	FirstName string
	LastName  string
}

type water struct {
	ServerTimeStr string
	RotateDeg     int
}

func queriedSeason(r *http.Request) any {
	season, ok := r.Context().Value(queryCtxKey).(string)
	if !ok {
		return nil
	}
	return season
}

func userInfo(r *http.Request) any {
	user, ok := r.Context().Value(userCtxKey).(user)
	if !ok {
		return nil
	}
	return user
}

func waterInfo(r *http.Request) any {
	posStr := r.URL.Query().Get("pos")
	offsetStr := r.URL.Query().Get("off")

	w := water{
		RotateDeg:     0,
		ServerTimeStr: time.Now().Format(time.RFC1123),
	}

	off, err := strconv.Atoi(offsetStr)
	if err != nil {
		return w
	}
	pos, err := strconv.Atoi(posStr)
	if err != nil {
		return w
	}

	w.RotateDeg = pos + off

	return w
}
