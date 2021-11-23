package options

import (
	"github.com/cockroachdb/pebble"
	ldb "github.com/syndtr/goleveldb/leveldb/opt"
)

//Contains all available read/write options
//for the supported databases. Fields are set
//by functions using the SetOptions signature.
type Options struct {
	LeveldbRead  *ldb.ReadOptions
	LeveldbWrite *ldb.WriteOptions
	PebbleWrite  *pebble.WriteOptions
}

//Defining the signature of functions accepted
//as parameters by engine.Store functions.
type SetOptions func(cfg *Options) error

//Closure returning a function that adds the leveldbRead
//parameter to an Options struct's LeveldbRead field.
func WithLeveldbRead(leveldbRead *ldb.ReadOptions) SetOptions {
	return func(cfg *Options) error {
		cfg.LeveldbRead = leveldbRead
		return nil
	}
}

//Closure returning a function that adds the leveldbWrite
//parameter to an Options struct's LeveldbWrite field.
func WithLeveldbWrite(leveldbWrite *ldb.WriteOptions) SetOptions {
	return func(cfg *Options) error {
		cfg.LeveldbWrite = leveldbWrite
		return nil
	}
}

//Closure returning a function that adds the pebbleWrite
//parameter to an Options struct's PebbleWrite field.
func WithPebbleWrite(pebbleWrite *pebble.WriteOptions) SetOptions {
	return func(cfg *Options) error {
		cfg.PebbleWrite = pebbleWrite
		return nil
	}
}
