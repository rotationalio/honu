package options

import (
	"github.com/cockroachdb/pebble"
	ldb "github.com/syndtr/goleveldb/leveldb/opt"
)

type Options struct {
	LeveldbRead  *ldb.ReadOptions
	LeveldbWrite *ldb.WriteOptions
	PebbleWrite  *pebble.WriteOptions
}

type SetOptions func(cfg *Options) error

func WithLeveldbRead(leveldbRead *ldb.ReadOptions) SetOptions {
	return func(cfg *Options) error {
		cfg.LeveldbRead = leveldbRead
		return nil
	}
}

func WithLeveldbWrite(leveldbWrite *ldb.WriteOptions) SetOptions {
	return func(cfg *Options) error {
		cfg.LeveldbWrite = leveldbWrite
		return nil
	}
}

func WithPebbleWrite(pebbleWrite *pebble.WriteOptions) SetOptions {
	return func(cfg *Options) error {
		cfg.PebbleWrite = pebbleWrite
		return nil
	}
}
