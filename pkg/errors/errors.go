package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// TODO: remove this error when we release Honu v1.0.0!
var ErrNotImplemented = Status(http.StatusNotImplemented, "not feature has not been implemented yet")

// Storage and query errors directly related to database operations.
var (
	ErrNotFound             = Status(http.StatusNotFound, "object not found")
	ErrReadOnlyDB           = Status(http.StatusUnprocessableEntity, "cannot execute operation in readonly mode")
	ErrReadOnlyTx           = Status(http.StatusUnprocessableEntity, "cannot execute operation: transaction is read only")
	ErrClosed               = Status(http.StatusGone, "database engine has been closed")
	ErrTxClosed             = Status(http.StatusGone, "transaction has already been committed or rolled back")
	ErrAlreadyExists        = Status(http.StatusConflict, "specified key already exists")
	ErrNoCollection         = Status(http.StatusNotFound, "collection with specified ID or name does not exist")
	ErrCollectionExists     = Status(http.StatusConflict, "collection with specified name already exists")
	ErrCollectionIdentifier = Status(http.StatusBadRequest, "collection identifier must be a name or ULID")
	ErrRepairCollection     = Status(http.StatusInternalServerError, "collection is malformed or incorrectly initialized requiring repair")
	ErrNotSupported         = Status(http.StatusNotImplemented, "operation not supported")
	ErrCreateID             = Status(http.StatusBadRequest, "cannot specify ID when creating new object")
	ErrIDMismatch           = Status(http.StatusBadRequest, "specified ID does not match resource ID")
	ErrNameMismatch         = Status(http.StatusBadRequest, "specified name does not match resource name")
	ErrNotInitialized       = Status(http.StatusInternalServerError, "store has not been properly initialized with system state")
)

// Access control errors
var (
	ErrAccessDenied = Status(http.StatusForbidden, "permission denied")
)

// Iteration Errors
var (
	ErrIterReleased = Status(http.StatusGone, "iterator has been released")
)

// Name validation errors (for collections and other restricted objects)
var (
	ErrInvalidName    = Status(http.StatusBadRequest, "identifier names must be alphanumeric or contain underscores and dashes")
	ErrEmptyName      = Status(http.StatusBadRequest, "identifier names cannot be empty")
	ErrNameChar       = Status(http.StatusBadRequest, "identifier names cannot contain spaces or punctuation")
	ErrNameDigitStart = Status(http.StatusBadRequest, "identifier names cannot start with a digit")
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
