package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"
	"go.rtnl.ai/honu/pkg/api/v1/credentials"
	"go.rtnl.ai/ulid"
	"go.rtnl.ai/x/backoff"
	"go.rtnl.ai/x/query"
)

const (
	userAgent    = "HonuDB API Client/v1"
	accept       = "application/json"
	acceptLang   = "en-US,en"
	acceptEncode = "gzip, deflate, br"
	contentType  = "application/json; charset=utf-8"
)

// New creates a new APIv1 client that implements the Client interface.
func New(endpoint string, opts ...ClientOption) (_ Client, err error) {
	c := &APIv1{}
	if c.endpoint, err = url.Parse(endpoint); err != nil {
		return nil, fmt.Errorf("could not parse endpoint: %s", err)
	}

	// Apply our options
	for _, opt := range opts {
		if err = opt(c); err != nil {
			return nil, err
		}
	}

	// If an http client isn't specified, create a default client.
	if c.client == nil {
		c.client = &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Timeout:       30 * time.Second,
		}

		// Create cookie jar for CSRF
		if c.client.Jar, err = cookiejar.New(nil); err != nil {
			return nil, fmt.Errorf("could not create cookiejar: %w", err)
		}
	}

	return c, nil
}

// APIv1 implements the v1 Client interface for making requests to a honudb replica.
type APIv1 struct {
	endpoint *url.URL                // the base url for all requests
	client   *http.Client            // used to make http requests to the server
	creds    credentials.Credentials // used to authenticate requests with the server
}

// Ensure the APIv1 implements the Client interface
var _ Client = &APIv1{}

//===========================================================================
// Client Methods
//===========================================================================

const statusEP = "/v1/status"

func (s *APIv1) Status(ctx context.Context) (out *StatusReply, err error) {
	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, statusEP, nil, nil); err != nil {
		return nil, err
	}

	// NOTE: we cannot use s.Do because we want to parse 503 Unavailable errors
	var rep *http.Response
	if rep, err = s.client.Do(req); err != nil {
		return nil, err
	}
	defer rep.Body.Close()

	// Detect other errors
	if rep.StatusCode != http.StatusOK && rep.StatusCode != http.StatusServiceUnavailable {
		return nil, fmt.Errorf("%s", rep.Status)
	}

	// Deserialize the JSON data from the response
	out = &StatusReply{}
	if err = json.NewDecoder(rep.Body).Decode(out); err != nil {
		return nil, fmt.Errorf("could not deserialize status reply: %s", err)
	}
	return out, nil
}

//===========================================================================
// Client Utility Methods
//===========================================================================

// Wait for ready polls the node's status endpoint until it responds with an 200
// response, retrying with exponential backoff or until the context deadline is expired.
// If the user does not supply a context with a deadline, then a default deadline of
// 5 minutes is used so that this method does not block indefinitely. If the node API
// service is ready (e.g. responds to a status request) then no error is returned,
// otherwise an error is returned if the node never responds.
//
// NOTE: if the node returns a 503 Service Unavailable because it is in maintenance
// mode, this method will continue to wait until the deadline for the node to exit
// from maintenance mode and be ready again.
func (s *APIv1) WaitForReady(ctx context.Context) (err error) {
	// If context does not have a deadline, create a context with a default deadline.
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()
	}

	// Create the status request to send until ready
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, "/v1/status", nil, nil); err != nil {
		return err
	}

	// Create a closure to repeatedly call the status endpoint
	checkReady := func() (_ bool, err error) {
		var rep *http.Response
		if rep, err = s.client.Do(req); err != nil {
			return false, err
		}
		defer rep.Body.Close()

		if rep.StatusCode < 200 || rep.StatusCode >= 300 {
			return false, &StatusError{Code: rep.StatusCode, Reply: Reply{Success: false, Error: http.StatusText(rep.StatusCode)}}
		}
		return true, nil
	}

	// Create notification handler for backoff
	notify := func(err error, delay time.Duration) {
		log.Warn().Err(err).
			Str("endpoint", s.endpoint.String()).
			Dur("delay", delay).
			Msg("waiting to connect to HonuDB node")
	}

	if _, err = backoff.Retry(ctx, checkReady, backoff.WithNotify(notify)); err != nil {
		return err
	}
	return nil
}

//===========================================================================
// REST Resource Methods
//===========================================================================

func (s *APIv1) List(ctx context.Context, endpoint string, in *PageQuery, out interface{}) (err error) {
	var params url.Values
	if params, err = query.Values(in); err != nil {
		return fmt.Errorf("could not encode page query: %w", err)
	}

	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, endpoint, nil, &params); err != nil {
		return err
	}

	if _, err = s.Do(req, &out, true); err != nil {
		return err
	}

	return nil
}

func (s *APIv1) Create(ctx context.Context, endpoint string, in, out interface{}) (err error) {
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPost, endpoint, in, nil); err != nil {
		return err
	}

	if _, err = s.Do(req, &out, true); err != nil {
		return err
	}
	return nil
}

func (s *APIv1) Detail(ctx context.Context, endpoint string, out interface{}) (err error) {
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, endpoint, nil, nil); err != nil {
		return err
	}

	if _, err = s.Do(req, &out, true); err != nil {
		return err
	}
	return nil
}

func (s *APIv1) Update(ctx context.Context, endpoint string, in, out interface{}) (err error) {
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPut, endpoint, in, nil); err != nil {
		return err
	}

	if _, err = s.Do(req, &out, true); err != nil {
		return err
	}
	return nil
}

func (s *APIv1) Delete(ctx context.Context, endpoint string) (err error) {
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodDelete, endpoint, nil, nil); err != nil {
		return err
	}

	if _, err = s.Do(req, nil, true); err != nil {
		return err
	}
	return nil
}

//===========================================================================
// Helper Methods
//===========================================================================

func (s *APIv1) NewRequest(ctx context.Context, method, path string, data interface{}, params *url.Values) (req *http.Request, err error) {
	// Resolve the URL reference from the path
	url := s.endpoint.ResolveReference(&url.URL{Path: path})
	if params != nil && len(*params) > 0 {
		url.RawQuery = params.Encode()
	}

	var body io.ReadWriter
	switch {
	case data == nil:
		body = nil
	default:
		body = &bytes.Buffer{}
		if err = json.NewEncoder(body).Encode(data); err != nil {
			return nil, fmt.Errorf("could not serialize request data as json: %s", err)
		}
	}

	// Create the http request
	if req, err = http.NewRequestWithContext(ctx, method, url.String(), body); err != nil {
		return nil, fmt.Errorf("could not create request: %s", err)
	}

	// Set the headers on the request
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Accept", accept)
	req.Header.Add("Accept-Language", acceptLang)
	req.Header.Add("Accept-Encoding", acceptEncode)
	req.Header.Add("Content-Type", contentType)

	// If there is a request ID on the context, set it on the request, otherwise generate one
	var requestID string
	if requestID, _ = RequestIDFromContext(ctx); requestID == "" {
		requestID = ulid.Make().String()
	}
	req.Header.Add("X-Request-ID", requestID)

	// Add authentication and authorization header.
	if s.creds != nil {
		var token string
		if token, err = s.creds.AccessToken(); err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+token)
	}

	// Add CSRF protection if its available
	if s.client.Jar != nil {
		cookies := s.client.Jar.Cookies(url)
		for _, cookie := range cookies {
			if cookie.Name == "csrf_token" {
				req.Header.Add("X-CSRF-TOKEN", cookie.Value)
			}
		}
	}

	return req, nil
}

// Do executes an http request against the server, performs error checking, and
// deserializes the response data into the specified struct.
func (s *APIv1) Do(req *http.Request, data interface{}, checkStatus bool) (rep *http.Response, err error) {
	if rep, err = s.client.Do(req); err != nil {
		return rep, fmt.Errorf("could not execute request: %s", err)
	}
	defer rep.Body.Close()

	// Detect http status errors if they've occurred
	if checkStatus {
		if rep.StatusCode < 200 || rep.StatusCode >= 300 {
			// Attempt to read the error response from JSON, if available
			serr := &StatusError{
				Code: rep.StatusCode,
			}

			if err = json.NewDecoder(rep.Body).Decode(&serr.Reply); err == nil {
				return rep, serr
			}

			serr.Reply = Unsuccessful
			return rep, serr
		}
	}

	// Deserialize the JSON data from the body
	if data != nil && rep.StatusCode >= 200 && rep.StatusCode < 300 && rep.StatusCode != http.StatusNoContent {
		ct := rep.Header.Get("Content-Type")
		if ct != "" {
			mt, _, err := mime.ParseMediaType(ct)
			if err != nil {
				return nil, fmt.Errorf("malformed content-type header: %w", err)
			}

			if mt != accept {
				return nil, fmt.Errorf("unexpected content type: %q", mt)
			}
		}

		if err = json.NewDecoder(rep.Body).Decode(data); err != nil {
			return nil, fmt.Errorf("could not deserialize response data: %s", err)
		}
	}

	return rep, nil
}
