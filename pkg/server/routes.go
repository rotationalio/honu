package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.rtnl.ai/honu/pkg/server/middleware"
)

// Sets up the server's middleware and routes.
func (s *Server) setupRoutes() (err error) {
	middleware := []middleware.Middleware{
		s.Maintenance(),
	}

	// Kubernetes liveness probes added before middleware.
	s.router.GET("/healthz", s.Healthz)
	s.router.GET("/livez", s.Healthz)
	s.router.GET("/readyz", s.Readyz)

	// API Routes
	// Status/Heartbeat endpoint
	s.addRoute(http.MethodGet, "/v1/status", s.Status, middleware...)

	// Collections resource
	s.addRoute(http.MethodGet, "/v1/collections", s.ListCollections, middleware...)
	s.addRoute(http.MethodPost, "/v1/collections", s.CreateCollection, middleware...)
	s.addRoute(http.MethodGet, "/v1/collections/:collectionID", s.RetrieveCollection, middleware...)
	s.addRoute(http.MethodPut, "/v1/collections/:collectionID", s.UpdateCollection, middleware...)
	s.addRoute(http.MethodDelete, "/v1/collections/:collectionID", s.DeleteCollection, middleware...)

	// Indexes resource
	s.addRoute(http.MethodGet, "/v1/collections/:collectionID/indexes", s.ListIndexes, middleware...)
	s.addRoute(http.MethodPost, "/v1/collections/:collectionID/indexes", s.CreateIndex, middleware...)
	s.addRoute(http.MethodGet, "/v1/collections/:collectionID/indexes/:indexID", s.RetrieveIndex, middleware...)
	s.addRoute(http.MethodPut, "/v1/collections/:collectionID/indexes/:indexID", s.UpdateIndex, middleware...)
	s.addRoute(http.MethodDelete, "/v1/collections/:collectionID/indexes/:indexID", s.DeleteIndex, middleware...)

	return nil
}

func (s *Server) addRoute(method, path string, h httprouter.Handle, m ...middleware.Middleware) {
	s.router.Handle(method, path, middleware.Chain(h, m...))
}
