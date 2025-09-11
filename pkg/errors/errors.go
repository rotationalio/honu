package errors

import (
	"errors"
	"fmt"
)

// Common errors returned by the store package.
var (
	ErrNotFound      = errors.New("object not found")
	ErrReadOnlyDB    = errors.New("cannot execute operation in readonly mode")
	ErrReadOnlyTx    = errors.New("cannot execute operation: transaction is read only")
	ErrClosed        = errors.New("database engine has been closed")
	ErrAlreadyExists = errors.New("specified key already exists")
	ErrNotSupported  = errors.New("operation not supported")
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
