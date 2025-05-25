package main

import (
	"fmt"
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

func counterBait(r *http.Request) any {
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
	fmt.Print(c)
	return c
}
