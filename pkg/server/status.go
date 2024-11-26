package server

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rotationalio/honu/pkg"
	"github.com/rotationalio/honu/pkg/api/v1"
	"github.com/rotationalio/honu/pkg/server/render"
)

const (
	serverStatusOK          = "ok"
	serverStatusNotReady    = "not ready"
	serverStatusUnhealthy   = "unhealthy"
	serverStatusMaintenance = "maintenance"
)

// Status reports the version and uptime of the server
func (s *Server) Status(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var state string
	s.RLock()
	switch {
	case s.healthy && s.ready:
		state = serverStatusOK
	case s.healthy && !s.ready:
		state = serverStatusNotReady
	case !s.healthy:
		state = serverStatusUnhealthy
	}
	s.RUnlock()

	render.Negotiate(r).Render(http.StatusOK, w, &api.StatusReply{
		Status:  state,
		Version: pkg.Version(),
		Uptime:  time.Since(s.started).String(),
	})
}

// Healthz is used to alert k8s to the health/liveness status of the server.
func (s *Server) Healthz(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s.RLock()
	healthy := s.healthy
	s.RUnlock()

	if !healthy {
		render.Text(http.StatusServiceUnavailable, w, serverStatusUnhealthy)
		return
	}

	render.Text(http.StatusOK, w, serverStatusOK)
}

// Readyz is used to alert k8s to the readiness status of the server.
func (s *Server) Readyz(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s.RLock()
	ready := s.ready
	s.RUnlock()

	if !ready {
		render.Text(http.StatusServiceUnavailable, w, serverStatusNotReady)
		return
	}

	render.Text(http.StatusOK, w, serverStatusOK)
}
