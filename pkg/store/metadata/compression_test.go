package metadata_test

import (
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/honu/pkg/store/lani"
	"go.rtnl.ai/honu/pkg/store/metadata"
)

func TestCompressionSize(t *testing.T) {
	var staticSize int
	staticSize += 1                     // Algorithm (uint8)
	staticSize += binary.MaxVarintLen64 // Level (int64)

	t.Run("StaticSize", func(t *testing.T) {
		compression := &metadata.Compression{}
		require.Equal(t, staticSize, compression.Size(), "expected zero valued compression to have a static size of %d bytes", staticSize)
	})

	t.Run("VariableSize", func(t *testing.T) {
		// Technically, the size of a Compression object is static because it does not contain any variable-length fields.
		var compression metadata.Compression
		loadFixture(t, "compression.json", &compression)
		require.Equal(t, 11, compression.Size(), "expected compression to have a size of 10 bytes as computed from fixture")
	})
}

func TestCompressionSerialization(t *testing.T) {
	var obj *metadata.Compression
	loadFixture(t, "compression.json", &obj)

	data, err := lani.Marshal(obj)
	require.NoError(t, err, "could not marshal compression")

	cmp := &metadata.Compression{}
	err = lani.Unmarshal(data, cmp)
	require.NoError(t, err, "could not unmarshal compression")
	require.Equal(t, obj, cmp, "deserialized compression does not match original")
}

func TestParseCompressionAlgorithm(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		tests := []struct {
			input    string
			expected metadata.CompressionAlgorithm
		}{
			{"none", metadata.None},
			{"gzip", metadata.GZIP},
			{"compress", metadata.COMPRESS},
			{"deflate", metadata.DEFLATE},
			{"brotli", metadata.BROTLI},
			{"NONE", metadata.None},
			{"GZIP", metadata.GZIP},
			{"COMPRESS", metadata.COMPRESS},
			{"DEFLATE", metadata.DEFLATE},
			{"BROTLI", metadata.BROTLI},
		}

		for _, test := range tests {
			t.Run(test.input, func(t *testing.T) {
				result, err := metadata.ParseCompressionAlgorithm(test.input)
				require.NoError(t, err, "expected no error for input %s", test.input)
				require.Equal(t, test.expected, result, "expected %s to parse to %d", test.input, test.expected)
			})
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		_, err := metadata.ParseCompressionAlgorithm("invalid")
		require.Error(t, err, "expected error for invalid compression algorithm")
	})
}

func TestCompressionJSON(t *testing.T) {
	tests := []metadata.CompressionAlgorithm{
		metadata.None,
		metadata.GZIP,
		metadata.COMPRESS,
		metadata.DEFLATE,
		metadata.BROTLI,
	}

	for _, orig := range tests {
		origs := orig.String()
		data, err := orig.MarshalJSON()
		require.NoError(t, err, "failed to marshal compression algorithm %s", origs)
		require.Equal(t, fmt.Sprintf(`"%s"`, origs), string(data), "expected JSON to match for %s", origs)

		var clone metadata.CompressionAlgorithm
		err = clone.UnmarshalJSON(data)
		require.NoError(t, err, "failed to unmarshal compression algorithm %s", origs)
		require.Equal(t, orig, clone, "expected unmarshaled compression algorithm to match original %s", origs)
	}
}
