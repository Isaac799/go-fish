package internal

import "errors"

const (
	htmlItemKindPage = iota
	htmlItemKindIsland
)

var (
	// ErrNoTemplateDir is given if I cannot find the desired dir to setup a new mux
	ErrNoTemplateDir = errors.New("cannot find template directory relative to working dir")
)
