package store

import "errors"

var (
	ErrGrowNegative = errors.New("cannot grow encoder by negative value")
	ErrTooLarge     = errors.New("cannot allocate more space for encoding")
	ErrNoLength     = errors.New("field not written with length value")
	ErrParseBoolean = errors.New("could not parse boolean value")
	ErrParseVarInt  = errors.New("could not parse varint")
)
