package metadata_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/honu/pkg/store/lani"
)

//===========================================================================
// Generic Tests for all Metadata types
//===========================================================================

type TestObject interface {
	lani.Encodable
	lani.Decodable
}

type TestCase struct {
	Name        string
	Fixture     string
	StaticSize  int
	FixtureSize int
	New         func() TestObject
}

func (tc *TestCase) TestStaticSize(t *testing.T) {
	t.Helper()
	obj := tc.New()
	require.Equal(t, tc.StaticSize, obj.Size(), "expected zero valued %s to have a static size of %d bytes", tc.Name, tc.StaticSize)
}

func (tc *TestCase) TestVariableSize(t *testing.T) {
	t.Helper()
	obj := tc.New()
	loadFixture(t, tc.Fixture, obj)
	require.Equal(t, tc.FixtureSize, obj.Size(), "expected %s to have a size of %d bytes as computed from fixture", tc.Name, tc.FixtureSize)
}

func (tc *TestCase) TestSerialization(t *testing.T) {
	t.Helper()
	obj := tc.New()
	loadFixture(t, tc.Fixture, obj)

	data, err := lani.Marshal(obj)
	require.NoError(t, err, "could not marshal %s", tc.Name)

	cmp := tc.New()
	err = lani.Unmarshal(data, cmp)
	require.NoError(t, err, "could not unmarshal %s", tc.Name)
	require.Equal(t, obj, cmp, "deserialized %s does not match original", tc.Name)
}

func loadFixture(t *testing.T, name string, v interface{}) {
	path := filepath.Join("testdata", name)
	f, err := os.Open(path)
	require.NoError(t, err, "could not open %s", path)
	defer f.Close()

	err = json.NewDecoder(f).Decode(v)
	require.NoError(t, err, "could not decode %s", path)
}

//===========================================================================
// Generic Tests for all Metadata Enums
//===========================================================================

type TestEnum interface {
	fmt.Stringer
	Value() uint8
}

type Serializable interface {
	TestEnum
	json.Marshaler
	json.Unmarshaler
}

type TestEnumCase struct {
	Name     string     // name of the enum
	Values   []TestEnum // the ordered enum values
	Strings  []string   // the expected string representations (same order as values)
	Invalid  []string   // invalid strings that should error on unmarshaling/parse
	Unknowns string     // the expected string representation of unknown or zero value
	ICase    bool       // whether the enum is case insensitive (adds lower case reprs to parse)
	ISpace   bool       // whether the enum is space insensitive (adds space reprs to parse)

	// Parsing function for the enum
	Parse func(string) (TestEnum, error)
	New   func(uint8) Serializable
}

func (tc *TestEnumCase) TestString(t *testing.T) {
	t.Helper()
	for i, val := range tc.Values {
		require.Equal(t, tc.Strings[i], val.String(), "expected %s to have string representation %q", tc.Name, tc.Strings[i])
	}
}

func (tc *TestEnumCase) TestStringBounds(t *testing.T) {
	t.Helper()
	max := uint8(0)
	min := uint8(255)

	for i := range tc.Values {
		if uint8(i) > max {
			max = uint8(i)
		}
		if uint8(i) < min {
			min = uint8(i)
		}
	}

	// Test one above the max
	above := tc.New(max + 1)
	require.Equal(t, tc.Unknowns, above.String(), "expected %s to have string representation %q for unknown value", tc.Name, tc.Unknowns)

	// Test zero value
	if min > 0 {
		zero := tc.New(0)
		require.Equal(t, tc.Unknowns, zero.String(), "expected %s to have string representation %q for zero value", tc.Name, tc.Unknowns)
	}
}

func (tc *TestEnumCase) TestParse(t *testing.T) {
	t.Helper()
	t.Run("Valid", func(t *testing.T) {
		type testCase struct {
			input    string
			expected TestEnum
		}

		tests := make([]testCase, 0, len(tc.Values)*3)

		// Add all the string representations
		for i, str := range tc.Strings {
			expected := tc.Values[i]
			tests = append(tests, testCase{input: str, expected: expected})
			if tc.ICase {
				tests = append(tests, testCase{input: strings.ToLower(str), expected: expected})
				tests = append(tests, testCase{input: strings.ToUpper(str), expected: expected})
			}
			if tc.ISpace {
				tests = append(tests, testCase{input: addSpaces(str), expected: expected})
			}
		}

		for _, test := range tests {
			actual, err := tc.Parse(test.input)
			require.NoError(t, err, "expected no error parsing input %q", test.input)
			require.Equal(t, test.expected, actual, "expected %q to parse to %v", test.input, test.expected)
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		invalid := []string{"", "foo", "bar", "123", "INVALID"}
		invalid = append(invalid, tc.Invalid...)
		for _, input := range invalid {
			_, err := tc.Parse(input)
			require.Error(t, err, "expected error parsing invalid input %q", input)
		}
	})
}

func (tc *TestEnumCase) TestJSON(t *testing.T) {
	t.Helper()
	// Test marshaling
	for _, val := range tc.Values {
		orig := tc.New(val.Value())
		data, err := json.Marshal(orig)
		require.NoError(t, err, "could not marshal %s value %q", tc.Name, orig.String())

		cmp := tc.New(0)
		err = json.Unmarshal(data, &cmp)
		require.NoError(t, err, "could not unmarshal %s value %q", tc.Name, orig.String())
		require.Equal(t, orig, cmp, "unmarshaled %s does not match original", tc.Name)
	}
}

func addSpaces(s string) string {
	fn := strings.Repeat(" ", rand.Intn(8))
	bn := strings.Repeat(" ", rand.Intn(4))
	return fn + s + bn
}
