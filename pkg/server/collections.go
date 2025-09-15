package server

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/julienschmidt/httprouter"
	"go.rtnl.ai/honu/pkg/errors"
	"go.rtnl.ai/honu/pkg/mime"
	"go.rtnl.ai/honu/pkg/server/render"
	"go.rtnl.ai/honu/pkg/store/metadata"
	"go.rtnl.ai/ulid"
)

func (s *Server) ListCollections(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		err         error
		collections []*metadata.Collection
	)

	if collections, err = s.db.Collections(); err != nil {
		render.Error(w, r, err)
		return
	}
	render.Negotiate(r).Render(http.StatusOK, w, collections)
}

func (s *Server) CreateCollection(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		err        error
		collection *metadata.Collection
	)

	collection = &metadata.Collection{}
	if err = mime.Bind(w, r, &collection); err != nil {
		render.Error(w, r, err)
		return
	}

	// TODO: validate collection name and fields

	if err = s.db.New(collection); err != nil {
		render.Error(w, r, err)
		return
	}

	render.Negotiate(r).Render(http.StatusCreated, w, collection)
}

func (s *Server) RetrieveCollection(w http.ResponseWriter, r *http.Request, q httprouter.Params) {
	var (
		err        error
		identifier any
		collection *metadata.Collection
	)

	// Attempt to parse the identifier as a ULID first, otherwise use it as a name string.
	identifier = parseIdentifier(q[0])
	if collection, err = s.db.Collection(identifier); err != nil {
		render.Error(w, r, err)
		return
	}

	render.Negotiate(r).Render(http.StatusOK, w, collection)
}

func (s *Server) UpdateCollection(w http.ResponseWriter, r *http.Request, q httprouter.Params) {
	var (
		err        error
		identifier any
		collection *metadata.Collection
	)

	// Attempt to parse the identifier as a ULID first, otherwise use it as a name string.
	identifier = parseIdentifier(q[0])

	// Parse the collection metadata from the request body.
	collection = &metadata.Collection{}
	if err = mime.Bind(w, r, &collection); err != nil {
		render.Error(w, r, err)
		return
	}

	// Ensure that the URL param is on the collection object if missing or that it matches.
	switch id := identifier.(type) {
	case ulid.ULID:
		if !collection.ID.IsZero() && !collection.ID.Equals(id) {
			render.Error(w, r, errors.ErrIDMismatch)
			return
		} else {
			collection.ID = id
		}
	case string:
		if collection.Name != "" && collection.Name != id {
			render.Error(w, r, errors.ErrNameMismatch)
			return
		} else {
			collection.Name = id
		}
	}

	// TODO: validate collection name and fields

	if err = s.db.Modify(collection); err != nil {
		render.Error(w, r, err)
		return
	}

	render.Negotiate(r).Render(http.StatusOK, w, collection)
}

func (s *Server) DeleteCollection(w http.ResponseWriter, r *http.Request, q httprouter.Params) {
	var (
		err        error
		params     url.Values
		identifier any
		operation  string
	)

	// Parse the operation from the query parameters
	if params, err = url.ParseQuery(r.URL.RawQuery); err != nil {
		render.Error(w, r, errors.Status(http.StatusBadRequest, "invalid query parameters"))
		return
	}

	// Attempt to parse the identifier as a ULID first, otherwise use it as a name string.
	identifier = parseIdentifier(q[0])

	operation = strings.ToLower(strings.TrimSpace(params.Get("operation")))
	switch operation {
	case "drop":
		if err = s.db.Drop(identifier); err != nil {
			render.Error(w, r, err)
			return
		}
	case "truncate":
		if err = s.db.Truncate(identifier); err != nil {
			render.Error(w, r, err)
			return
		}
	case "":
		render.Error(w, r, errors.Status(http.StatusBadRequest, "missing operation parameter"))
		return
	default:
		render.Error(w, r, errors.Status(http.StatusBadRequest, "invalid operation parameter"))
		return
	}
}

func (s *Server) ListIndexes(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}

func (s *Server) CreateIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}

func (s *Server) RetrieveIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}

func (s *Server) UpdateIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}

func (s *Server) DeleteIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}

// Returns a ULID from the parameter otherwise simply returns the parameter string value.
func parseIdentifier(param httprouter.Param) any {
	if id, err := ulid.Parse(param.Value); err == nil {
		return id
	}
	return param.Value
}
