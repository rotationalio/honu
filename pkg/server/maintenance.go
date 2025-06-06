package server

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.rtnl.ai/honu/pkg"
	"go.rtnl.ai/honu/pkg/api/v1"
	"go.rtnl.ai/honu/pkg/server/middleware"
	"go.rtnl.ai/honu/pkg/server/render"
)

// If the server is in maintenance mode, aborts the current request and renders the
// maintenance mode page instead. Returns nil if not in maintenance mode.
func (s *Server) Maintenance() middleware.Middleware {
	if s.conf.Maintenance {
		return func(next httprouter.Handle) httprouter.Handle {
			return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
				render.Negotiate(r).Render(http.StatusServiceUnavailable, w, &api.StatusReply{
					Status:  "maintenance",
					Version: pkg.Version(),
					Uptime:  time.Since(s.started).String(),
				})
			}
		}
	}
	return nil
}
