package api_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/rotationalio/honu/pkg/api/v1"
	"github.com/stretchr/testify/require"
)

var (
	ctx = context.Background()
)

func TestStatus(t *testing.T) {
	fixture := &api.StatusReply{}
	err := loadFixture("testdata/statusOK.json", fixture)
	require.NoError(t, err, "could not load status ok fixture")

	_, client := testServer(t, &testServerConfig{
		expectedMethod: http.MethodGet,
		expectedPath:   "/v1/status",
		fixture:        fixture,
		statusCode:     http.StatusOK,
	})

	rep, err := client.Status(ctx)
	require.NoError(t, err, "could not execute status request")
	require.Equal(t, fixture, rep, "expected reply to be equal to fixture")
}

type testServerConfig struct {
	expectedMethod string
	expectedPath   string
	statusCode     int
	fixture        interface{}
}

func testServer(t *testing.T, conf *testServerConfig) (*httptest.Server, api.Client) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != conf.expectedMethod {
			http.Error(w, fmt.Sprintf("expected method %s got %s", conf.expectedMethod, r.Method), http.StatusExpectationFailed)
			return
		}

		if r.URL.Path != conf.expectedPath {
			http.Error(w, fmt.Sprintf("expected path %s got %s", conf.expectedPath, r.URL.Path), http.StatusExpectationFailed)
			return
		}

		if conf.statusCode == 0 {
			conf.statusCode = http.StatusOK
		}

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(conf.statusCode)
		json.NewEncoder(w).Encode(conf.fixture)
	}))

	// Ensure the server is closed when the test is complete
	t.Cleanup(ts.Close)

	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")
	return ts, client
}

func loadFixture(path string, v interface{}) (err error) {
	var f *os.File
	if f, err = os.Open(path); err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(v)
}
