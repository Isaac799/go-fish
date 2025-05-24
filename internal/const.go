package internal

import "errors"

const (
	htmlItemKindPage = iota
	htmlItemKindIsland
	htmlItemKindStyle
)

const (
	// browserCacheDurationSeconds is used to cache documents
	// such as .css. To help prevent invalid cache we replace
	// the names with a hash of their content
	browserCacheDurationSeconds = 86400 // 1 day
)

var (
	// ErrNoTemplateDir is given if I cannot find the desired dir to setup a new mux
	ErrNoTemplateDir = errors.New("cannot find template directory relative to working dir")
	// ErrInvalidExtension is given if a file is discoverd that I did not anticipate
	ErrInvalidExtension = errors.New("invalid file extension")
)
