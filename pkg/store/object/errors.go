package object

import "errors"

var (
	ErrBadVersion = errors.New("object is malformed: cannot decode specified version")
	ErrNoMetadata = errors.New("object is malformed: no associated metadata")
)
