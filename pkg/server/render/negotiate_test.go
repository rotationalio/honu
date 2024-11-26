package render_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rotationalio/honu/pkg/server/render"
	"github.com/stretchr/testify/require"
)

func TestParseAccept(t *testing.T) {
	testCases := []struct {
		header   string
		expected []string
	}{
		{
			"", []string{},
		},
		{
			"application/json",
			[]string{"application/json"},
		},
		{
			"application/json; charset=utf-8",
			[]string{"application/json"},
		},
		{
			"application/json,     ",
			[]string{"application/json"},
		},
		{
			"application/json,application/msgpack",
			[]string{"application/json", "application/msgpack"},
		},
	}

	for i, tc := range testCases {
		accepted := render.ParseAccept(tc.header)
		require.Equal(t, tc.expected, accepted, "test case %d failed", i)
	}
}

func TestAccepted(t *testing.T) {
	testCases := []struct {
		header   []string
		expected []string
	}{
		{
			nil, []string{},
		},
		{
			[]string{"application/json"},
			[]string{"application/json"},
		},
		{
			[]string{"application/json; charset=utf-8"},
			[]string{"application/json"},
		},
		{
			[]string{"application/json,     "},
			[]string{"application/json"},
		},
		{
			[]string{"application/json,application/msgpack"},
			[]string{"application/json", "application/msgpack"},
		},
	}

	for i, tc := range testCases {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header[render.Accept] = tc.header

		accepted := render.Accepted(r)
		require.Equal(t, tc.expected, accepted, "test case %d failed", i)
	}
}

func TestAcceptedNoHeader(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	require.Nil(t, render.Accepted(r), "expected no Accept header to return nil")
}
