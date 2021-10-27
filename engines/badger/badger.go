package badger

import (
	"errors"

	"github.com/rotationalio/honu/config"
)

func Open(conf config.ReplicaConfig) (*BadgerEngine, error) {
	return &BadgerEngine{}, errors.New("not implemented yet")
}

type BadgerEngine struct{}

func (db *BadgerEngine) Engine() string {
	return "badger"
}

func (db *BadgerEngine) Close() error {
	return errors.New("not implemented yet")
}
