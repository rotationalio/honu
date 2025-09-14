package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.rtnl.ai/honu/pkg/api/v1"
	"go.rtnl.ai/honu/pkg/server/render"
	"go.rtnl.ai/honu/pkg/store/metadata"
)

func (s *Server) ListCollections(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		err         error
		collections []*metadata.Collection
	)

	renderer := render.Negotiate(r)
	if collections, err = s.db.Collections(); err != nil {
		renderer.Render(http.StatusInternalServerError, w, api.Error(err))
		return
	}
	renderer.Render(http.StatusOK, w, collections)
}

func (s *Server) CreateCollection(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		err        error
		collection *metadata.Collection
	)

	if err = s.db.New(collection); err != nil {
		return
	}
	render.Negotiate(r).Render(http.StatusCreated, w, collection)
}

func (s *Server) RetrieveCollection(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}

func (s *Server) UpdateCollection(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}

func (s *Server) DeleteCollection(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}

func (s *Server) ListIndexes(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}

func (s *Server) CreateIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}

func (s *Server) RetrieveIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}

func (s *Server) UpdateIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}

func (s *Server) DeleteIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}
