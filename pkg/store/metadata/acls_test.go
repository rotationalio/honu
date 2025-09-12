package metadata_test

import (
	"testing"

	"go.rtnl.ai/honu/pkg/store/metadata"
)

func TestAccessControl(t *testing.T) {
	staticSize := 17 // 16 bytes for ClientID and 1 byte for permissions (uint8)

	// Create a test generic case and execute the tests
	testCase := &TestCase{
		Name:        "AccessControl",
		Fixture:     "acl.json",
		StaticSize:  staticSize,
		FixtureSize: 17,
		New:         func() TestObject { return &metadata.AccessControl{} },
	}

	t.Run("StaticSize", testCase.TestStaticSize)
	t.Run("VariableSize", testCase.TestVariableSize)
	t.Run("Serialization", testCase.TestSerialization)
}
