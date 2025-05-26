package main

import (
	"net/http"
	"strconv"
)

type counter struct {
	Value int
}

func (c *counter) increment() {
	c.Value++
	if c.Value > 3 {
		c.Value = 0
	}
}

func incrementQueryCount(r *http.Request) any {
	count := r.URL.Query().Get("count")
	if len(count) == 0 {
		return counter{}
	}
	i, err := strconv.Atoi(count)
	if err != nil {
		return counter{}
	}
	c := counter{
		Value: i,
	}
	c.increment()
	return c
}

type user struct {
	ID        int
	FirstName string
	LastName  string
}

func findUser(r *http.Request) any {
	count := r.PathValue("id")
	if len(count) == 0 {
		return nil
	}
	i, err := strconv.Atoi(count)
	if err != nil {
		return nil
	}

	users := map[int]user{
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

	if u, exists := users[i]; exists {
		return u
	}
	return nil
}
