package object

import "errors"

var (
	ErrBadVersion = errors.New("object is malformed: cannot decode specified version")
	ErrMalformed  = errors.New("object is malformed: cannot parse data or metadata")
)
