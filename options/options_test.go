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
	require.Equal(t, options.NamespaceDefault, opts.Namespace)
	require.False(t, opts.Force)
	require.False(t, opts.Tombstones)

	// Test setting multiple options
	opts, err = options.New(options.WithLevelDBRead(&ldb.ReadOptions{Strict: ldb.StrictJournal}), options.WithNamespace("foo"))
	require.NoError(t, err, "could not create options")
	require.Equal(t, "foo", opts.Namespace)
	require.NotNil(t, opts.LevelDBRead)
	require.Equal(t, ldb.StrictJournal, opts.LevelDBRead.Strict)

	// Ensuring setting empty string namespace still ends up as the default namespace
	opts, err = options.New(options.WithNamespace(""))
	require.NoError(t, err, "could not create options with empty string namespace")
	require.Equal(t, options.NamespaceDefault, opts.Namespace)

	// Test boolean options
	opts, err = options.New(options.WithForce(), options.WithTombstones())
	require.NoError(t, err, "boolean options returned an error")
	require.True(t, opts.Force)
	require.True(t, opts.Tombstones)
}

func TestLevelDBReadOptions(t *testing.T) {
	readOptions := &ldb.ReadOptions{
		DontFillCache: true,
		Strict:        1,
	}
	cfg := &options.Options{}
	ldbReadFunc := options.WithLevelDBRead(readOptions)
	ldbReadFunc(cfg)
	require.Equal(t, cfg.LevelDBRead, readOptions)
}

func TestLevelDBWriteOptions(t *testing.T) {
	writeOptions := &ldb.WriteOptions{
		NoWriteMerge: true,
		Sync:         true,
	}
	cfg := &options.Options{}
	ldbWriteFunc := options.WithLevelDBWrite(writeOptions)
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
