package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rotationalio/honu/pkg/server/middleware"
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

	return nil
}

func (s *Server) addRoute(method, path string, h httprouter.Handle, m ...middleware.Middleware) {
	s.router.Handle(method, path, middleware.Chain(h, m...))
}
