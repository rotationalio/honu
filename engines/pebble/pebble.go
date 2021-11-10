package pebble

import (
	"errors"

	"github.com/cockroachdb/pebble"
	"github.com/rotationalio/honu/config"
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

func (db *PebbleEngine) Engine() string {
	return "pebble"
}

func (db *PebbleEngine) Close() error {
	return db.pebble.Close()
}

func (db *PebbleEngine) Get(key []byte, options ...opts.SetOptions) (value []byte, err error) {
	if len(options) > 0 {
		return nil, errors.New("pebble does not take read options")
	}
	value, closer, err := db.pebble.Get(key)
	if err != nil {
		return nil, err
	}
	if err := closer.Close(); err != nil {
		return nil, err
	}
	return value, nil
}

func (db *PebbleEngine) Put(key, value []byte, options ...opts.SetOptions) error {
	var cfg *opts.Options
	cfg.PebbleWrite = nil
	for _, option := range options {
		if err := option(cfg); err != nil {
			return err
		}
	}
	return db.pebble.Set(key, value, cfg.PebbleWrite)
}

func (db *PebbleEngine) Delete(key []byte, options ...opts.SetOptions) error {
	var cfg *opts.Options
	cfg.PebbleWrite = nil
	for _, option := range options {
		if err := option(cfg); err != nil {
			return err
		}
	}
	return db.pebble.Delete(key, cfg.PebbleWrite)
}
