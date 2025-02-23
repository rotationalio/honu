package store

import (
	"github.com/rotationalio/honu/pkg/config"
	"github.com/rotationalio/honu/pkg/store/engine"
	"github.com/rotationalio/honu/pkg/store/engine/leveldb"
	"github.com/rotationalio/honu/pkg/store/lamport"
)

// Store is a Honu embedded database.
type Store struct {
	lamport.Clock
	db engine.Engine
}

func Open(conf config.Config) (s *Store, err error) {
	s = &Store{
		Clock: lamport.New(conf.PID),
	}

	if s.db, err = leveldb.Open(conf.Store); err != nil {
		return nil, err
	}

	return s, nil
}
