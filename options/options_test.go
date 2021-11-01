package options_test

import (
	"errors"
	"testing"

	"github.com/rotationalio/honu/options"
	"github.com/stretchr/testify/require"
)

//Test to ensure that the correct error is thrown when an improperly
//formated error string is passed to options.parse()
func TestInvalidParseOptions(t *testing.T) {
	expectedError := errors.New("improperly formated option string")
	ldbOptions := options.LeveldbOptions{}

	options := "SingleElementOptionString"
	_, err := ldbOptions.Write(options)
	require.Equal(t, err, expectedError)

	options = "Invalid String Without Comma"
	_, err = ldbOptions.Read(options)
	require.Equal(t, err, expectedError)
}

func TestLevelDBReadOptions(t *testing.T) {
	ldbOptions := options.LeveldbOptions{}
	options := "DontFillCache true, Strict 1"
	readOptions, err := ldbOptions.Read(options)
	require.NoError(t, err)
	require.Equal(t, int(readOptions.Strict), 1)
	require.True(t, readOptions.DontFillCache)

	options = "foo true"
	_, err = ldbOptions.Write(options)
	require.Error(t, err, "foo is not a valid leveldb readoption")
}

func TestLevelDBWriteOptions(t *testing.T) {
	ldbOptions := options.LeveldbOptions{}
	options := "NoWriteMerge true, Sync false"
	WriteOptions, err := ldbOptions.Write(options)
	require.NoError(t, err)
	require.True(t, WriteOptions.NoWriteMerge)
	require.False(t, WriteOptions.Sync)

	options = "bar true"
	_, err = ldbOptions.Write(options)
	require.Error(t, err, "bar is not a valid leveldb writeoption")
}

func TestPebbleWriteOptions(t *testing.T) {
	pebbleOptions := options.PebbleOptions{}
	options := "Sync true"
	WriteOptions, err := pebbleOptions.Write(options)
	require.NoError(t, err)
	require.True(t, WriteOptions.Sync)

	options = "foobar true"
	_, err = pebbleOptions.Write(options)
	require.Error(t, err, "foobar is not a valid pebble writeoption")
}
