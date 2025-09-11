package metadata_test

import (
	"encoding/binary"
	"testing"

	"go.rtnl.ai/honu/pkg/store/metadata"
)

func TestCompression(t *testing.T) {
	var staticSize int
	staticSize += 1                     // Algorithm (uint8)
	staticSize += binary.MaxVarintLen64 // Level (int64)

	// Create a test generic case and execute the tests
	testCase := &TestCase{
		Name:        "Compression",
		Fixture:     "compression.json",
		StaticSize:  staticSize,
		FixtureSize: 11,
		New:         func() TestObject { return &metadata.Compression{} },
	}

	t.Run("StaticSize", testCase.TestStaticSize)
	t.Run("VariableSize", testCase.TestVariableSize)
	t.Run("Serialization", testCase.TestSerialization)
}

func TestCompressionAlgorithm(t *testing.T) {
	testCase := &TestEnumCase{
		Name: "CompressionAlgorithm",
		Values: []TestEnum{
			metadata.None,
			metadata.GZIP,
			metadata.COMPRESS,
			metadata.DEFLATE,
			metadata.BROTLI,
		},
		Strings: []string{
			"NONE",
			"GZIP",
			"COMPRESS",
			"DEFLATE",
			"BROTLI",
		},
		Unknowns: "UNKNOWN",
		ICase:    true,
		ISpace:   true,
		Parse:    func(s string) (TestEnum, error) { return metadata.ParseCompressionAlgorithm(s) },
		New:      func(i uint8) Serializable { val := metadata.CompressionAlgorithm(i); return &val },
	}

	t.Run("String", testCase.TestString)
	t.Run("StringBounds", testCase.TestStringBounds)
	t.Run("Parse", testCase.TestParse)
	t.Run("JSON", testCase.TestJSON)
}
