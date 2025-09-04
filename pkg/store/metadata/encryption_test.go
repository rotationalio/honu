package metadata_test

import (
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/honu/pkg/store/lani"
	"go.rtnl.ai/honu/pkg/store/metadata"
)

func TestEncryptionSize(t *testing.T) {
	// Compute the static size of the Encryption struct
	var staticSize int
	staticSize += binary.MaxVarintLen64 // Length of PublicKeyID
	staticSize += binary.MaxVarintLen64 // Length of EncryptionKey
	staticSize += binary.MaxVarintLen64 // Length of HMACSecret
	staticSize += binary.MaxVarintLen64 // Length of Signature
	staticSize += 3                     // SealingAlgorithm, EncryptionAlgorithm, SignatureAlgorithm (all uint8)

	t.Run("StaticSize", func(t *testing.T) {
		encryption := &metadata.Encryption{}
		require.Equal(t, staticSize, encryption.Size(), "expected zero valued encryption to have a static size of %d bytes", staticSize)
	})

	t.Run("VariableSize", func(t *testing.T) {
		var encryption metadata.Encryption
		loadFixture(t, "encryption.json", &encryption)
		require.Equal(t, 197, encryption.Size(), "expected encryption to have a size of 197 bytes as computed from fixture")
	})
}

func TestEncryptionSerialization(t *testing.T) {
	var obj *metadata.Encryption
	loadFixture(t, "encryption.json", &obj)

	data, err := lani.Marshal(obj)
	require.NoError(t, err, "could not marshal encryption")

	cmp := &metadata.Encryption{}
	err = lani.Unmarshal(data, cmp)
	require.NoError(t, err, "could not unmarshal encryption")
	require.Equal(t, obj, cmp, "deserialized encryption does not match original")
}

func TestParseEncryptionAlgorithm(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		tests := []struct {
			input    string
			expected metadata.EncryptionAlgorithm
		}{
			{"plaintext", metadata.Plaintext},
			{"aes256_gcm", metadata.AES256_GCM},
			{"aes192_gcm", metadata.AES192_GCM},
			{"aes128_gcm", metadata.AES128_GCM},
			{"hmac_sha256", metadata.HMAC_SHA256},
			{"rsa_oeap_sha512", metadata.RSA_OEAP_SHA512},
			{"AES256_GCM", metadata.AES256_GCM},
			{"AES192_GCM", metadata.AES192_GCM},
			{"AES128_GCM", metadata.AES128_GCM},
			{"HMAC_SHA256", metadata.HMAC_SHA256},
			{"RSA_OEAP_SHA512", metadata.RSA_OEAP_SHA512},
		}

		for _, tc := range tests {
			algo, err := metadata.ParseEncryptionAlgorithm(tc.input)
			require.NoError(t, err, "expected no error for valid encryption algorithm: %s", tc.input)
			require.Equal(t, tc.expected, algo, "expected %s to parse to %d but got %d", tc.input, tc.expected, algo)
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		testCases := []string{
			"",
			"foo",
			"AES 182 GCM",
		}

		for _, tc := range testCases {
			_, err := metadata.ParseEncryptionAlgorithm(tc)
			require.Error(t, err, "expected error for invalid encryption algorithm: %s", tc)
		}
	})
}

func TestEncryptionJSON(t *testing.T) {
	tests := []metadata.EncryptionAlgorithm{
		metadata.Plaintext,
		metadata.AES256_GCM,
		metadata.AES192_GCM,
		metadata.AES128_GCM,
		metadata.HMAC_SHA256,
		metadata.RSA_OEAP_SHA512,
	}

	for _, orig := range tests {
		origs := orig.String()
		data, err := orig.MarshalJSON()
		require.NoError(t, err, "failed to marshal encryption algorithm %s", origs)
		require.Equal(t, fmt.Sprintf(`"%s"`, origs), string(data), "expected JSON to match for %s", origs)

		var clone metadata.EncryptionAlgorithm
		err = clone.UnmarshalJSON(data)
		require.NoError(t, err, "failed to unmarshal encryption algorithm %s", origs)
		require.Equal(t, orig, clone, "expected unmarshaled encryption algorithm to match original %s", origs)
	}

}
