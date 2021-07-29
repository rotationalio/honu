package honu_test

import (
	"testing"

	"github.com/rotationalio/honu"
	"github.com/stretchr/testify/require"
)

func TestDSNParsing(t *testing.T) {
	cases := []struct {
		uri string
		dsn *honu.DSN
	}{
		{"leveldb:///fixtures/db", &honu.DSN{"leveldb", "fixtures/db"}},
		{"sqlite3:///fixtures/db", &honu.DSN{"sqlite3", "fixtures/db"}},
		{"leveldb:////data/db", &honu.DSN{"leveldb", "/data/db"}},
		{"sqlite3:////data/db", &honu.DSN{"sqlite3", "/data/db"}},
	}

	for _, tc := range cases {
		dsn, err := honu.ParseDSN(tc.uri)
		require.NoError(t, err)
		require.Equal(t, tc.dsn, dsn)
	}

	// Test error cases
	_, err := honu.ParseDSN("foo")
	require.Error(t, err)

	_, err = honu.ParseDSN("foo://")
	require.Error(t, err)
}
