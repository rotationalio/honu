package db_test

import (
	"testing"

	. "github.com/rotationalio/honu/db"
	"github.com/stretchr/testify/require"
)

func TestDSNParsing(t *testing.T) {
	cases := []struct {
		uri string
		dsn *DSN
	}{
		{"leveldb:///fixtures/db", &DSN{"leveldb", "fixtures/db"}},
		{"sqlite3:///fixtures/db", &DSN{"sqlite3", "fixtures/db"}},
		{"leveldb:////data/db", &DSN{"leveldb", "/data/db"}},
		{"sqlite3:////data/db", &DSN{"sqlite3", "/data/db"}},
	}

	for _, tc := range cases {
		dsn, err := ParseDSN(tc.uri)
		require.NoError(t, err)
		require.Equal(t, tc.dsn, dsn)
	}

	// Test error cases
	_, err := ParseDSN("foo")
	require.Error(t, err)

	_, err = ParseDSN("foo://")
	require.Error(t, err)
}
