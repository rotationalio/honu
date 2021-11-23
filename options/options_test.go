package options_test

import (
	"testing"

	"github.com/cockroachdb/pebble"
	"github.com/rotationalio/honu/options"
	"github.com/stretchr/testify/require"
	ldb "github.com/syndtr/goleveldb/leveldb/opt"
)

func TestLevelDBReadOptions(t *testing.T) {
	readOptions := &ldb.ReadOptions{
		DontFillCache: true,
		Strict:        1,
	}
	cfg := &options.Options{}
	ldbReadFunc := options.WithLeveldbRead(readOptions)
	ldbReadFunc(cfg)
	require.Equal(t, cfg.LeveldbRead, readOptions)
}

func TestLevelDBWriteOptions(t *testing.T) {
	writeOptions := &ldb.WriteOptions{
		NoWriteMerge: true,
		Sync:         true,
	}
	cfg := &options.Options{}
	ldbWriteFunc := options.WithLeveldbWrite(writeOptions)
	ldbWriteFunc(cfg)
	require.Equal(t, cfg.LeveldbWrite, writeOptions)

}

func TestPebbleWriteOptions(t *testing.T) {
	writeOptions := &pebble.WriteOptions{
		Sync: true,
	}
	cfg := &options.Options{}
	pebbleWriteFunc := options.WithPebbleWrite(writeOptions)
	pebbleWriteFunc(cfg)
	require.Equal(t, cfg.PebbleWrite, writeOptions)
}
