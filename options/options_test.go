package options_test

import (
	"testing"

	"github.com/rotationalio/honu/options"
	"github.com/stretchr/testify/require"
)

func TestInvalidParseOptions(t *testing.T) {
	ldbOptions := options.LeveldbOptions{}
	options := "SingleElementOptionString"
	_, err := ldbOptions.Write(&options)
	require.Error(t, err, "improperly formated option string")
}

func TestLevelDBReadOptions(t *testing.T) {
	ldbOptions := options.LeveldbOptions{}
	options := "DontFillCache true, Strict 1"
	readOptions, err := ldbOptions.Read(&options)
	require.NoError(t, err)
	require.Equal(t, int(readOptions.Strict), 1)
	require.True(t, readOptions.DontFillCache)

	options = "foo true"
	_, err = ldbOptions.Write(&options)
	require.Error(t, err, "foo is not a valid leveldb readoption")
}

func TestLevelDBWriteOptions(t *testing.T) {
	ldbOptions := options.LeveldbOptions{}
	options := "NoWriteMerge true, Sync false"
	WriteOptions, err := ldbOptions.Write(&options)
	require.NoError(t, err)
	require.True(t, WriteOptions.NoWriteMerge)
	require.False(t, WriteOptions.Sync)

	options = "bar true"
	_, err = ldbOptions.Write(&options)
	require.Error(t, err, "bar is not a valid leveldb writeoption")
}

func TestPebbleWriteOptions(t *testing.T) {
	pebbleOptions := options.PebbleOptions{}
	options := "Sync true"
	WriteOptions, err := pebbleOptions.Write(&options)
	require.NoError(t, err)
	require.True(t, WriteOptions.Sync)

	options = "foobar true"
	_, err = pebbleOptions.Write(&options)
	require.Error(t, err, "foobar is not a valid pebble writeoption")
}
