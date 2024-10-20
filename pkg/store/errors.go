package store

import "errors"

var (
	ErrNoLength     = errors.New("field not written with length value")
	ErrParseBoolean = errors.New("could not parse boolean value")
)
