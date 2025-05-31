// Package internal is the inner workings of the gofish package
package internal

import "errors"

const (
	// FishKindTuna is a big fish. Served as a page. Consumes Sardines.
	// Identified by mime [ text/html ].
	// Not cahced.
	FishKindTuna = iota
	// FishKindSardine is a small fish. Used by tuna. Smaller templates, served standalone too.
	// Identified by mime [ text/html ] & underscore prefix.
	// Not cahced.
	FishKindSardine
	// FiskKindClown is a decorative fish. Used in head of document.
	// Identified by mime [ text/css | text/javascript ].
	// Is cached & name from hash.
	FiskKindClown
	// FiskKindAnchovy is supportive of the tuna.
	// Identified by mime [ image | audio | video ].
	// Is cached.
	FiskKindAnchovy
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

var fishKindStr = map[int]string{
	FishKindTuna:    "Tuna",
	FishKindSardine: "Sardine",
	FiskKindClown:   "Clown",
	FiskKindAnchovy: "Anchovy",
}
