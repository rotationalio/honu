package honu

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// DSN represents the parsed components of an embedded database service.
type DSN struct {
	Scheme string
	Path   string
}

// DSN Parsing and Handling
func ParseDSN(uri string) (_ *DSN, err error) {
	dsn, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("could not parse dsn: %s", err)
	}

	if dsn.Scheme == "" || dsn.Path == "" {
		return nil, errors.New("could not parse dsn, specify scheme:///relative/path/to/db")
	}

	return &DSN{
		Scheme: dsn.Scheme,
		Path:   strings.TrimPrefix(dsn.Path, "/"),
	}, nil
}
