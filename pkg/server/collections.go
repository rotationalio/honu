package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s *Server) ListCollections(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}

func (s *Server) CreateCollection(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}

func (s *Server) RetrieveCollection(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}

func (s *Server) UpdateCollection(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}

func (s *Server) DeleteCollection(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}

func (s *Server) ListIndexes(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}

func (s *Server) CreateIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}

func (s *Server) RetrieveIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}

func (s *Server) UpdateIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}

func (s *Server) DeleteIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}
