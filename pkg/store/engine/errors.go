package engine

import "errors"

var (
	ErrNotFound      = errors.New("object not found")
	ErrReadOnlyDB    = errors.New("cannot execute operation in readonly mode")
	ErrReadOnlyTx    = errors.New("cannot execute operation: transaction is read only")
	ErrClosed        = errors.New("database engine has been closed")
	ErrAlreadyExists = errors.New("specified key already exists")
	ErrNotSupported  = errors.New("operation not supported")
)
