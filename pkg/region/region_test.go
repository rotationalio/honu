package region_test

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/honu/pkg/region"
	"go.rtnl.ai/honu/pkg/store/lani"
)

func TestList(t *testing.T) {
	regions := region.List()
	require.Greater(t, len(regions), 0, "expected at least one region")
	require.NotContains(t, regions, region.UNKNOWN, "expected UNKNOWN region to be excluded from list")
}

func TestParse(t *testing.T) {
	regions := region.List()

	t.Run("Strings", func(t *testing.T) {
		for _, r := range regions {
			parsed, err := region.Parse(r.String())
			require.NoError(t, err, "failed to parse valid region string %q", r.String())
			require.Equal(t, r, parsed, "parsed region does not match original for string %q", r.String())
		}
	})

	t.Run("InvalidStrings", func(t *testing.T) {
		invalids := []string{"", "INVALID", "unknown", "us_east_1", "123"}
		for _, s := range invalids {
			parsed, err := region.Parse(s)
			require.Error(t, err, "expected error when parsing invalid region string %q", s)
			require.Equal(t, region.UNKNOWN, parsed, "expected UNKNOWN region when parsing invalid string %q", s)
		}
	})

	t.Run("CaseInsensitivity", func(t *testing.T) {
		for _, r := range regions {
			s := r.String()
			upper := strings.ToUpper(s)
			lower := strings.ToLower(s)

			parsed, err := region.Parse(upper)
			require.NoError(t, err, "failed to parse valid region string %q", upper)
			require.Equal(t, r, parsed, "parsed region does not match original for string %q", upper)

			parsed, err = region.Parse(lower)
			require.NoError(t, err, "failed to parse valid region string %q", lower)
			require.Equal(t, r, parsed, "parsed region does not match original for string %q", lower)
		}
	})

	t.Run("TrimSpace", func(t *testing.T) {
		for _, r := range regions {
			s := addSpaces(r.String())
			parsed, err := region.Parse(s)
			require.NoError(t, err, "failed to parse valid region string %q", s)
			require.Equal(t, r, parsed, "parsed region does not match original for string %q", s)
		}
	})

	t.Run("Hyphens", func(t *testing.T) {
		for _, r := range regions {
			s := strings.ReplaceAll(r.String(), "_", "-")
			parsed, err := region.Parse(s)
			require.NoError(t, err, "failed to parse valid region string %q", s)
			require.Equal(t, r, parsed, "parsed region does not match original for string %q", s)
		}
	})

	t.Run("Underscores", func(t *testing.T) {
		for _, r := range regions {
			s := strings.ReplaceAll(r.String(), "-", "_")
			parsed, err := region.Parse(s)
			require.NoError(t, err, "failed to parse valid region string %q", s)
			require.Equal(t, r, parsed, "parsed region does not match original for string %q", s)
		}
	})

	t.Run("Uint32", func(t *testing.T) {
		for _, r := range regions {
			parsed, err := region.Parse(uint32(r))
			require.NoError(t, err, "failed to parse valid region uint32 %d", uint32(r))
			require.Equal(t, r, parsed, "parsed region does not match original for uint32 %d", uint32(r))
		}
	})

	t.Run("UnknownUint32", func(t *testing.T) {
		maxVal := uint32(0)
		for _, r := range regions {
			if uint32(r) > maxVal {
				maxVal = uint32(r)
			}
		}

		parsed, _ := region.Parse(maxVal + 1)
		require.Equal(t, region.UNKNOWN.String(), parsed.String(), "expected UNKNOWN region when parsing unknown uint32 %d", maxVal+1)
	})

	t.Run("Region", func(t *testing.T) {
		for _, r := range regions {
			parsed, err := region.Parse(r)
			require.NoError(t, err, "failed to parse valid region %v", r)
			require.Equal(t, r, parsed, "parsed region does not match original for region %v", r)
		}
	})

	t.Run("InvalidType", func(t *testing.T) {
		parsed, err := region.Parse(3.14)
		require.Error(t, err, "expected error when parsing invalid type")
		require.Equal(t, region.UNKNOWN, parsed, "expected UNKNOWN region when parsing invalid type")
	})
}

func TestJSON(t *testing.T) {
	regions := region.List()

	t.Run("Standard", func(t *testing.T) {
		for _, r := range regions {
			data, err := r.MarshalJSON()
			require.NoError(t, err, "failed to marshal region %v to JSON", r)

			var unmarshaled region.Region
			err = unmarshaled.UnmarshalJSON(data)
			require.NoError(t, err, "failed to unmarshal region %v from JSON", r)
			require.Equal(t, r, unmarshaled, "unmarshaled region does not match original for region %v", r)
		}
	})

	t.Run("Unmarshal", func(t *testing.T) {
		t.Run("Strings", func(t *testing.T) {
			for _, r := range regions {
				data := []byte(`"` + r.String() + `"`)
				var unmarshaled region.Region
				err := unmarshaled.UnmarshalJSON(data)
				require.NoError(t, err, "failed to unmarshal region %v from JSON string", r)
				require.Equal(t, r, unmarshaled, "unmarshaled region does not match original for region %v from JSON string", r)
			}
		})

		t.Run("Numeric", func(t *testing.T) {
			for _, r := range regions {
				data := []byte(strconv.FormatUint(uint64(r), 10))
				var unmarshaled region.Region
				err := unmarshaled.UnmarshalJSON(data)
				require.NoError(t, err, "failed to unmarshal region %v from JSON numeric", r)
				require.Equal(t, r, unmarshaled, "unmarshaled region does not match original for region %v from JSON numeric", r)
			}
		})

		t.Run("Invalid", func(t *testing.T) {
			invalids := [][]byte{[]byte(`"`), []byte(`"INVALID"`), []byte(`123.45`), []byte(`{}`), []byte(`[]`), []byte(`true`)}
			for _, data := range invalids {
				var unmarshaled region.Region
				err := unmarshaled.UnmarshalJSON(data)
				require.Error(t, err, "expected error when unmarshaling invalid JSON %s", string(data))
				require.Equal(t, region.UNKNOWN, unmarshaled, "expected UNKNOWN region when unmarshaling invalid JSON %s", string(data))
			}
		})
	})
}

func TestLani(t *testing.T) {
	regions := region.List()

	for _, r := range regions {
		require.Equal(t, 5, r.Size())

		data, err := lani.Marshal(r)
		require.NoError(t, err, "failed to marshal region %v with lani", r)

		var unmarshaled region.Region
		err = lani.Unmarshal(data, &unmarshaled)
		require.NoError(t, err, "failed to unmarshal region %v with lani", r)
		require.Equal(t, r, unmarshaled, "unmarshaled region does not match original for region %v with lani", r)
	}

	require.Equal(t, 10+(5*len(regions)), regions.Size())

	data, err := lani.Marshal(regions)
	require.NoError(t, err, "failed to marshal regions list with lani")

	var unmarshaled region.Regions
	err = lani.Unmarshal(data, &unmarshaled)
	require.NoError(t, err, "failed to unmarshal regions list with lani")
	require.Equal(t, regions, unmarshaled, "unmarshaled regions list does not match original with lani")
}

func addSpaces(s string) string {
	fn := strings.Repeat(" ", rand.Intn(8))
	bn := strings.Repeat(" ", rand.Intn(4))
	return fn + s + bn
}
