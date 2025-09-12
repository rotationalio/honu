package errors

import (
	"errors"
	"fmt"
)

// TODO: remove this error when we release Honu v1.0.0!
var ErrNotImplemented = errors.New("not feature has not been implemented yet")

// Storage errors that indicate a bad request.
var (
	ErrNotFound             = errors.New("object not found")
	ErrReadOnlyDB           = errors.New("cannot execute operation in readonly mode")
	ErrReadOnlyTx           = errors.New("cannot execute operation: transaction is read only")
	ErrClosed               = errors.New("database engine has been closed")
	ErrAlreadyExists        = errors.New("specified key already exists")
	ErrNoCollection         = errors.New("collection with specified ID or name does not exist")
	ErrCollectionExists     = errors.New("collection with specified name already exists")
	ErrCollectionIdentifier = errors.New("collection identifier must be a name or ULID")
	ErrNotSupported         = errors.New("operation not supported")
)

// Name validation errors (for collections and other restricted objects)
var (
	ErrInvalidName    = errors.New("identifier names must be alphanumeric or contain underscores and dashes")
	ErrEmptyName      = errors.New("identifier names cannot be empty")
	ErrNameChar       = errors.New("identifier names cannot contain spaces or punctuation")
	ErrNameDigitStart = errors.New("identifier names cannot start with a digit")
)

// In is a helper function to check if an error is in a list of errors.
func In(err error, targets ...error) bool {
	for _, target := range targets {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

// Reduce namespacing conflicts by adding error functions from the errors package.
var (
	New  = errors.New
	Fmt  = fmt.Errorf
	Is   = errors.Is
	As   = errors.As
	Join = errors.Join
)
