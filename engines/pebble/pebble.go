package pebble

import (
	"errors"

	"github.com/cockroachdb/pebble"
	"github.com/rotationalio/honu/config"
	engine "github.com/rotationalio/honu/engines"
	"github.com/rotationalio/honu/iterator"
	opts "github.com/rotationalio/honu/options"
)

type PebbleEngine struct {
	pebble *pebble.DB
}

// TODO: Allow Passing Pebble Options
func Open(path string, conf config.ReplicaConfig) (_ *PebbleEngine, err error) {
	engine := &PebbleEngine{}
	if engine.pebble, err = pebble.Open(path, nil); err != nil {
		return nil, err
	}
	return engine, nil
}

//Returns a string giving the engine type.
func (db *PebbleEngine) Engine() string {
	return "pebble"
}

//Close the database.
func (db *PebbleEngine) Close() error {
	return db.pebble.Close()
}

// Get the latest version of the object stored by the key.
func (db *PebbleEngine) Get(key []byte, options ...opts.SetOptions) (value []byte, err error) {
	if len(options) > 0 {
		return nil, errors.New("pebble does not take read options")
	}
	value, closer, err := db.pebble.Get(key)
	if err != nil && errors.Is(err, pebble.ErrNotFound) {
		return value, engine.ErrNotFound
	}
	if err := closer.Close(); err != nil {
		return nil, err
	}
	return value, nil
}

// Put a new value to the specified key and update the version.
func (db *PebbleEngine) Put(key, value []byte, options ...opts.SetOptions) error {
	var cfg *opts.Options
	cfg.PebbleWrite = nil
	for _, setOption := range options {
		if err := setOption(cfg); err != nil {
			return err
		}
	}
	return db.pebble.Set(key, value, cfg.PebbleWrite)
}

// Delete the object represented by the key, creating a tombstone object.
func (db *PebbleEngine) Delete(key []byte, options ...opts.SetOptions) error {
	var cfg *opts.Options
	cfg.PebbleWrite = nil
	for _, setOption := range options {
		if err := setOption(cfg); err != nil {
			return err
		}
	}
	return db.pebble.Delete(key, cfg.PebbleWrite)
}

//TODO: Implement pebble iteration (engines/pebble/iter.go)
func (db *PebbleEngine) Iter(prefix []byte) (i iterator.Iterator, err error) {
	return nil, nil
}
