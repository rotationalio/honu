package pebble

import "github.com/rotationalio/honu/config"

func Open(conf config.ReplicaConfig) (*PebbleEngine, error) {
	return &PebbleEngine{}, nil
}

type PebbleEngine struct{}

func (db *PebbleEngine) Engine() string {
	return "pebble"
}

func (db *PebbleEngine) Close() error {
	return nil
}
