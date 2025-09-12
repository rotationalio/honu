package metadata_test

import (
	"encoding/binary"
	"testing"

	"go.rtnl.ai/honu/pkg/store/metadata"
)

func TestPublisher(t *testing.T) {
	// Compute the static size of the Encryption struct
	var staticSize int
	staticSize += 16 + 16               // PublisherID and ClientID (both ULID)
	staticSize += binary.MaxVarintLen64 // Length of IPAddress
	staticSize += binary.MaxVarintLen64 // Length of UserAgent

	// Create a test generic case and execute the tests
	testCase := &TestCase{
		Name:        "Publisher",
		Fixture:     "publisher.json",
		StaticSize:  staticSize,
		FixtureSize: 77,
		New:         func() TestObject { return &metadata.Publisher{} },
	}

	t.Run("StaticSize", testCase.TestStaticSize)
	t.Run("VariableSize", testCase.TestVariableSize)
	t.Run("Serialization", testCase.TestSerialization)
}
