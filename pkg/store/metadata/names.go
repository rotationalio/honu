package metadata

import (
	"fmt"
	"regexp"
	"unicode"

	"go.rtnl.ai/honu/pkg/errors"
)

// Names in Honu are used to identify collections or to create indexable key/value pairs
// for faster lookups of objects. In order to ensure a consistent development
// experience, names cannot be any string, but must follow a specific set of rules.
// Names are case sensitive, must not contain spaces or punctuation except for
// underscores and dashes, and must not start with a number.
func ValidateName(s string) error {
	if s == "" {
		return errors.ErrEmptyName
	}

	for i, c := range s {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) && c != '_' && c != '-' {
			return fmt.Errorf("%w: invalid character '%c' at position %d", errors.ErrNameChar, c, i)
		}

		if i == 0 && (unicode.IsDigit(c) || c == '-') {
			return errors.ErrNameDigitStart
		}
	}

	return nil
}

var identifierRE = regexp.MustCompile(`^[\p{L}_][\p{L}\p{Nd}_-]*$`)

// Name identifier validation uses a regular expression to quickly check if the name
// is valid. This method is primarily used for benchmark comparisons and test coverage
// validation since the ValidateName method actually has better benchmark performance.
func ValidateNameRegex(s string) error {
	if !identifierRE.MatchString(s) {
		return errors.ErrInvalidName
	}
	return nil
}
