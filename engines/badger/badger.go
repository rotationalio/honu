package badger

import "github.com/rotationalio/honu/config"

func Open(conf config.ReplicaConfig) (*BadgerEngine, error) {
	return &BadgerEngine{}, nil
}

type BadgerEngine struct{}

func (db *BadgerEngine) Engine() string {
	return "badger"
}

func (db *BadgerEngine) Close() error {
	return nil
}
