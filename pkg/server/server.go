package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.rtnl.ai/honu/pkg/config"
	"go.rtnl.ai/honu/pkg/logger"
	"go.rtnl.ai/honu/pkg/store"
)

func init() {
	// Initialize zerolog with GCP logging requirements
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.TimestampFieldName = logger.GCPFieldKeyTime
	zerolog.MessageFieldName = logger.GCPFieldKeyMsg

	// Add the severity hook for GCP logging
	var gcpHook logger.SeverityHook
	log.Logger = zerolog.New(os.Stdout).Hook(gcpHook).With().Timestamp().Logger()
}

// Create a new Honu database server/replica instance using the specified configuration.
// This function is the main entry point to initializing a honudb instance and should
// be called rather than constructing a server directly. This method ensures that the
// configuration is correctly loaded from the environment, that the logging defaults
// are set correctly, and that any observability tools are correctly configured.
func New(conf config.Config) (s *Server, err error) {
	// Load the default configuration from the environment
	if conf.IsZero() {
		if conf, err = config.New(); err != nil {
			return nil, err
		}
	}

	// Set the global level
	zerolog.SetGlobalLevel(conf.GetLogLevel())

	// Set human readable logging if specified.
	if conf.ConsoleLog {
		console := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		log.Logger = zerolog.New(console).With().Timestamp().Logger()
	}

	// Create the server and prepare to serve.
	s = &Server{
		conf: conf,
		errc: make(chan error, 1),
	}

	// Initialize the underlying store
	if s.db, err = store.Open(conf); err != nil {
		return nil, err
	}

	// Create the httprouter
	s.router = httprouter.New()
	s.router.RedirectFixedPath = true
	s.router.HandleMethodNotAllowed = true
	s.router.RedirectTrailingSlash = true
	if err = s.setupRoutes(); err != nil {
		return nil, err
	}

	// Create the http server
	s.srv = &http.Server{
		Addr:              s.conf.BindAddr,
		Handler:           s.router,
		ErrorLog:          nil,
		ReadHeaderTimeout: s.conf.ReadTimeout,
		WriteTimeout:      s.conf.WriteTimeout,
		IdleTimeout:       s.conf.IdleTimeout,
	}

	return s, nil
}

// A Honu Database server implements several services as enabled for interaction with
// the Honu replica network including:
//
// 1. A database api client for users to interact with the database
// 2. An administrative client for managing peers and the replica
// 3. A replication service with auto-adapting anti-entropy replication
// 4. A metrics server for prometheus to scrape data
//
// The server may also implement background services as required.
type Server struct {
	sync.RWMutex
	db      *store.Store
	conf    config.Config
	srv     *http.Server
	router  *httprouter.Router
	url     *url.URL
	started time.Time
	errc    chan error
	healthy bool
	ready   bool
}

func (s *Server) Serve() (err error) {
	// Create a socket to listen on and infer the final URL.
	// NOTE: if the bindaddr is 127.0.0.1:0 for testing, a random port will be assigned,
	// manually creating the listener will allow us to determine which port.
	// When we start listening all incoming requests will be buffered until the server
	// actually starts up in its own go routine below.
	var sock net.Listener
	if sock, err = net.Listen("tcp", s.srv.Addr); err != nil {
		return fmt.Errorf("could not listen on bind addr %s: %s", s.srv.Addr, err)
	}

	s.setURL(sock.Addr())
	s.SetStatus(true, true)
	s.started = time.Now()

	// Listen for HTTP requests and handle them.
	go func() {
		// Make sure we don't use the external err to avoid data races.
		if serr := s.serve(sock); !errors.Is(serr, http.ErrServerClosed) {
			s.errc <- serr
		}
	}()

	log.Info().Str("url", s.URL()).Msg("honu database server started")
	return <-s.errc
}

// ServeTLS if a tls configuration is provided, otherwise Serve.
func (s *Server) serve(sock net.Listener) error {
	if s.srv.TLSConfig != nil {
		return s.srv.ServeTLS(sock, "", "")
	}
	return s.srv.Serve(sock)
}

func (s *Server) Shutdown() (err error) {
	log.Info().Msg("gracefully shutting down honu database server")
	s.SetStatus(false, false)

	ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
	defer cancel()

	s.srv.SetKeepAlivesEnabled(false)
	if err = s.srv.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

// SetStatus sets the health and ready status on the server, modifying the behavior of
// the kubernetes probe responses.
func (s *Server) SetStatus(health, ready bool) {
	s.Lock()
	s.healthy = health
	s.ready = ready
	s.Unlock()
	log.Debug().Bool("health", health).Bool("ready", ready).Msg("server status set")
}

// URL returns the endpoint of the server as determined by the configuration and the
// socket address and port (if specified).
func (s *Server) URL() string {
	s.RLock()
	defer s.RUnlock()
	return s.url.String()
}

func (s *Server) setURL(addr net.Addr) {
	s.Lock()
	defer s.Unlock()

	s.url = &url.URL{
		Scheme: "http",
		Host:   addr.String(),
	}

	if s.srv.TLSConfig != nil {
		s.url.Scheme = "https"
	}

	if tcp, ok := addr.(*net.TCPAddr); ok && tcp.IP.IsUnspecified() {
		s.url.Host = fmt.Sprintf("127.0.0.1:%d", tcp.Port)
	}
}
