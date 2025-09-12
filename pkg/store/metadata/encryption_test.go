package metadata_test

import (
	"encoding/binary"
	"testing"

	"go.rtnl.ai/honu/pkg/store/metadata"
)

func TestEncryption(t *testing.T) {
	// Compute the static size of the Encryption struct
	var staticSize int
	staticSize += binary.MaxVarintLen64 // Length of PublicKeyID
	staticSize += binary.MaxVarintLen64 // Length of EncryptionKey
	staticSize += binary.MaxVarintLen64 // Length of HMACSecret
	staticSize += binary.MaxVarintLen64 // Length of Signature
	staticSize += 3                     // SealingAlgorithm, EncryptionAlgorithm, SignatureAlgorithm (all uint8)

	// Create a test generic case and execute the tests
	testCase := &TestCase{
		Name:        "Encryption",
		Fixture:     "encryption.json",
		StaticSize:  staticSize,
		FixtureSize: 197,
		New:         func() TestObject { return &metadata.Encryption{} },
	}

	t.Run("StaticSize", testCase.TestStaticSize)
	t.Run("VariableSize", testCase.TestVariableSize)
	t.Run("Serialization", testCase.TestSerialization)
}

func TestEncryptionAlgorithm(t *testing.T) {
	testCase := &TestEnumCase{
		Name: "EncryptionAlgorithm",
		Values: []TestEnum{
			metadata.Plaintext,
			metadata.AES256_GCM,
			metadata.AES192_GCM,
			metadata.AES128_GCM,
			metadata.HMAC_SHA256,
			metadata.RSA_OEAP_SHA512,
		},
		Strings: []string{
			"PLAINTEXT",
			"AES256_GCM",
			"AES192_GCM",
			"AES128_GCM",
			"HMAC_SHA256",
			"RSA_OEAP_SHA512",
		},
		Unknowns: "UNKNOWN",
		ICase:    true,
		ISpace:   true,
		Parse:    func(s string) (TestEnum, error) { return metadata.ParseEncryptionAlgorithm(s) },
		New:      func(i uint8) Serializable { val := metadata.EncryptionAlgorithm(i); return &val },
	}

	t.Run("String", testCase.TestString)
	t.Run("StringBounds", testCase.TestStringBounds)
	t.Run("Parse", testCase.TestParse)
	t.Run("JSON", testCase.TestJSON)
}
