package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	. "github.com/rotationalio/honu/pkg/server/middleware"
	"github.com/stretchr/testify/require"
)

func MakeTestMiddleware(name string, abort bool, calls *Calls) Middleware {
	return func(next httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			calls.before(name)
			if !abort {
				next(w, r, p)
				calls.after(name)
			}
		}
	}
}

func MakeTestHandler(calls *Calls) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		calls.call("handler")
		fmt.Fprintln(w, "success")
	}
}

func TestChain(t *testing.T) {
	calls := &Calls{}
	h := Chain(
		MakeTestHandler(calls),
		MakeTestMiddleware("A", false, calls),
		MakeTestMiddleware("B", false, calls),
		MakeTestMiddleware("C", false, calls),
		MakeTestMiddleware("D", false, calls),
	)

	srv := testServer(t, h)
	_, err := srv.Client().Get(srv.URL + "/")
	require.NoError(t, err, "expected no error making request")

	expected := []string{
		"A-before",
		"B-before",
		"C-before",
		"D-before",
		"handler",
		"D-after",
		"C-after",
		"B-after",
		"A-after",
	}
	require.Equal(t, len(expected), calls.calls, "incorrect number of calls")
	require.Equal(t, expected, calls.callers, "middleware not chained correctly")
}

func TestAbort(t *testing.T) {
	calls := &Calls{}
	h := Chain(
		MakeTestHandler(calls),
		MakeTestMiddleware("A", false, calls),
		MakeTestMiddleware("B", true, calls),
		MakeTestMiddleware("C", false, calls),
		MakeTestMiddleware("D", false, calls),
	)

	srv := testServer(t, h)
	_, err := srv.Client().Get(srv.URL + "/")
	require.NoError(t, err, "expected no error making request")

	expected := []string{
		"A-before",
		"B-before",
		"A-after",
	}
	require.Equal(t, len(expected), calls.calls, "incorrect number of calls")
	require.Equal(t, expected, calls.callers, "middleware not chained correctly")
}

func TestChainWithNil(t *testing.T) {
	calls := &Calls{}
	h := Chain(
		MakeTestHandler(calls),
		nil,
		MakeTestMiddleware("A", false, calls),
		nil, nil, nil,
		MakeTestMiddleware("B", false, calls),
		nil, nil, nil,
	)

	srv := testServer(t, h)
	_, err := srv.Client().Get(srv.URL + "/")
	require.NoError(t, err, "expected no error making request")

	expected := []string{
		"A-before",
		"B-before",
		"handler",
		"B-after",
		"A-after",
	}
	require.Equal(t, len(expected), calls.calls, "incorrect number of calls")
	require.Equal(t, expected, calls.callers, "middleware not chained correctly")
}

func testServer(t *testing.T, h httprouter.Handle) *httptest.Server {
	// Setup the test server and router
	router := httprouter.New()
	router.GET("/", h)

	srv := httptest.NewServer(router)

	// Ensure the server is closed when the test is complete
	t.Cleanup(srv.Close)

	return srv
}

type Calls struct {
	calls   int
	callers []string
}

func (c *Calls) call(name string) {
	if c.callers == nil {
		c.callers = make([]string, 0, 16)
	}

	c.callers = append(c.callers, name)
	c.calls++
}

func (c *Calls) before(name string) {
	c.call(fmt.Sprintf("%s-before", name))
}

func (c *Calls) after(name string) {
	c.call(fmt.Sprintf("%s-after", name))
}
