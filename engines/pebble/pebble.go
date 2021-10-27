package pebble

import (
	"errors"

	"github.com/rotationalio/honu/config"
)

func Open(conf config.ReplicaConfig) (*PebbleEngine, error) {
	return &PebbleEngine{}, errors.New("not implemented yet")
}

type PebbleEngine struct{}

func (db *PebbleEngine) Engine() string {
	return "pebble"
}

func (db *PebbleEngine) Close() error {
	return errors.New("not implemented yet")
}
