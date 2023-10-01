package engine

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrReadOnlyTx    = errors.New("cannot execute a write operation in a read only transaction")
	ErrAlreadyExists = errors.New("specified key already exists in the database")
)
