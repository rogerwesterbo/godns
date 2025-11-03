package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rogerwesterbo/godns/internal/httpserver/httproutes"
	"github.com/rogerwesterbo/godns/internal/services/v1zoneservice"
	"github.com/vitistack/common/pkg/loggers/vlog"
)

// HTTPServer represents the HTTP API server
type HTTPServer struct {
	address     string
	server      *http.Server
	zoneService *v1zoneservice.V1ZoneService
}

// New creates a new HTTP server instance
func New(address string, zoneService *v1zoneservice.V1ZoneService) *HTTPServer {
	return &HTTPServer{
		address:     address,
		zoneService: zoneService,
	}
}

// Start starts the HTTP server
func (s *HTTPServer) Start() error {
	// Create router with all routes
	router := httproutes.NewRouter(s.zoneService)

	// Configure HTTP server
	s.server = &http.Server{
		Addr:         s.address,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	vlog.Infof("Starting HTTP API server on %s", s.address)

	// Start server in a goroutine
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			vlog.Errorf("HTTP server error: %v", err)
		}
	}()

	return nil
}

// Stop gracefully stops the HTTP server
func (s *HTTPServer) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	vlog.Info("Stopping HTTP API server...")
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}

	vlog.Info("HTTP API server stopped")
	return nil
}
