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
