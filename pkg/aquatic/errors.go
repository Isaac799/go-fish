package aquatic

import "errors"

var (
	// ErrNewFish is provided if anything goes wrong generating a fish from a pond
	ErrNewFish = errors.New("making new fish")
	// ErrCoral is provided when failing to read the bytes for a fish
	ErrCoral = errors.New("reading coral bytes")
	// ErrReef is provided when failing to read the bytes for a fish
	// and its prerequisite fish (e.g. sardine)
	ErrReef = errors.New("reading reef bytes")
	// ErrHandle is provided when a http handler fails to 'catch' a fish
	ErrHandle = errors.New("http handler")
	// ErrNewPond is provided is something goes wrong making a pond
	ErrNewPond = errors.New("making new pond")
	// ErrCollectFish is provided if we have a problem collecting the fish from a new pond
	ErrCollectFish = errors.New("collecting fish")
)

var (
	// ErrNoTemplateDir is given if I cannot find the desired dir to setup a new mux
	ErrNoTemplateDir = errors.New("cannot find template directory relative to working dir")
	// ErrInvalidExtension is given if a file is discoverd that I did not anticipate
	ErrInvalidExtension = errors.New("invalid file extension")
)
