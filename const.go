package main

import "errors"

var (
	ErrNoTemplateDir = errors.New("cannot find template directory. make sure it is at the same level as the executable.")
)
