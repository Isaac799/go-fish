package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Isaac799/go-fish/example/bridge"
)

type user struct {
	ID        int
	FirstName string
	LastName  string
}

type water struct {
	ServerTime time.Time
	RotateDeg  int
}

type fishData struct {
	Season string
	User   *user
	Water  *water
	Form   *bridge.HTMLElement
}

func queriedSeason(r *http.Request) *fishData {
	data := fishData{}
	season, ok := r.Context().Value(queryCtxKey).(string)
	if !ok {
		return nil
	}
	data.Season = season
	return &data
}

func userInfo(r *http.Request) *fishData {
	data := fishData{}
	user, ok := r.Context().Value(userCtxKey).(user)
	if !ok {
		return nil
	}
	data.User = &user
	return &data
}

func waterInfo(r *http.Request) *fishData {
	data := fishData{}
	posStr := r.URL.Query().Get("pos")
	offsetStr := r.URL.Query().Get("off")

	w := water{
		RotateDeg:  0,
		ServerTime: time.Now(),
	}

	off, err := strconv.Atoi(offsetStr)
	if err != nil {
		data.Water = &w
		return &data
	}
	pos, err := strconv.Atoi(posStr)
	if err != nil {
		data.Water = &w
		return &data
	}

	w.RotateDeg = pos + off

	data.Water = &w
	return &data
}
