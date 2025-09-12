package metadata_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/honu/pkg/errors"
	. "go.rtnl.ai/honu/pkg/store/metadata"
)

func TestNameIdentifiers(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		testCases := []string{
			"myVariable",
			"età",
			"你好",
			"_variable1",
			"_42",
			"a೬",
			"t0",
			"my_awesome_collection",
			"_my_protected_collection",
			"my-cool-collection",
			"ends-with-dash-",
			"ends_with_underscore_",
			"_starts_with_underscore",
			"élève",
			"Straße",
			"пример",
			"שלום",
			"こんにちは",
			"안녕하세요",
			"résumé",
			"naïve",
			"español",
			"über",
		}

		t.Run("Regex", func(t *testing.T) {
			for i, identifier := range testCases {
				err := ValidateNameRegex(identifier)
				require.NoError(t, err, "test case %d: expected no error for identifier %q", i, identifier)
			}
		})

		t.Run("CharLoop", func(t *testing.T) {
			for i, identifier := range testCases {
				err := ValidateName(identifier)
				require.NoError(t, err, "test case %d: expected no error for identifier %q", i, identifier)
			}
		})
	})

	t.Run("Invalid", func(t *testing.T) {
		testCases := []struct {
			identifier string
			expected   error
		}{
			{"", errors.ErrEmptyName},
			{"1variable", errors.ErrNameDigitStart},
			{"೬variable", errors.ErrNameDigitStart},
			{"variable name", errors.ErrNameChar},
			{"variable-name!", errors.ErrNameChar},
			{"-starts-with-dash", errors.ErrNameDigitStart},
			{"abc※def", errors.ErrNameChar},
			{"_🤩_", errors.ErrNameChar},
		}

		t.Run("Regex", func(t *testing.T) {
			for i, tc := range testCases {
				err := ValidateNameRegex(tc.identifier)
				require.Error(t, err, "test case %d: expected error for identifier %q", i, tc.identifier)
			}
		})

		t.Run("CharLoop", func(t *testing.T) {
			for i, tc := range testCases {
				err := ValidateName(tc.identifier)
				require.ErrorIs(t, err, tc.expected, "test case %d: expected error for identifier %q", i, tc.identifier)
			}
		})
	})
}

func BenchmarkNameIdentifiers(b *testing.B) {
	b.Run("Valid", func(b *testing.B) {
		s := "über_೬12-你好"

		b.Run("Regex", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ValidateNameRegex(s)
			}
		})

		b.Run("CharLoop", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ValidateName(s)
			}
		})
	})

	b.Run("Invalid", func(b *testing.B) {
		s := "_abcdefgeh-你好-123-※def"

		b.Run("Regex", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ValidateNameRegex(s)
			}
		})

		b.Run("CharLoop", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ValidateName(s)
			}
		})
	})
}
