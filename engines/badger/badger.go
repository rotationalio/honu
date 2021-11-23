package badger

import (
	"errors"

	"github.com/rotationalio/honu/config"
	engine "github.com/rotationalio/honu/engines"
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

func (db *BadgerEngine) Begin(readonly bool) (engine.Transaction, error) {
	return nil, errors.New("not implemented yet")
}
