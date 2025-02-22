package store

import "github.com/rotationalio/honu/pkg/store/lamport"

// Store is a Honu embedded database.
type Store struct {
	lamport.Clock
}

func Open(dsn string) (*Store, error) {
	s := &Store{
		Clock: lamport.New(1),
	}
	return s, nil
}
