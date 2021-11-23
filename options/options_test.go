package options_test

import (
	"testing"

	"github.com/cockroachdb/pebble"
	"github.com/rotationalio/honu/options"
	"github.com/stretchr/testify/require"
	ldb "github.com/syndtr/goleveldb/leveldb/opt"
)

func TestHonuOptions(t *testing.T) {
	// Test default options
	opts, err := options.New()
	require.NoError(t, err, "could not create options")
	require.Equal(t, "default", opts.Namespace)
	require.False(t, opts.Destroy)

	// Test setting multiple options
	opts, err = options.New(options.WithDestroy(), options.WithNamespace("foo"))
	require.NoError(t, err, "could not create options")
	require.Equal(t, "foo", opts.Namespace)
	require.True(t, opts.Destroy)
}

func TestLevelDBReadOptions(t *testing.T) {
	readOptions := &ldb.ReadOptions{
		DontFillCache: true,
		Strict:        1,
	}
	cfg := &options.Options{}
	ldbReadFunc := options.WithLeveldbRead(readOptions)
	ldbReadFunc(cfg)
	require.Equal(t, cfg.LevelDBRead, readOptions)
}

func TestLevelDBWriteOptions(t *testing.T) {
	writeOptions := &ldb.WriteOptions{
		NoWriteMerge: true,
		Sync:         true,
	}
	cfg := &options.Options{}
	ldbWriteFunc := options.WithLeveldbWrite(writeOptions)
	ldbWriteFunc(cfg)
	require.Equal(t, cfg.LevelDBWrite, writeOptions)

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
