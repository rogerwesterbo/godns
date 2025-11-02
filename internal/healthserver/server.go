package healthserver

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/vitistack/common/pkg/loggers/vlog"
)

// Server manages liveness and readiness probe HTTP servers
type Server struct {
	livenessAddr  string
	readinessAddr string
	ready         atomic.Bool
	livenessHTTP  *http.Server
	readinessHTTP *http.Server
}

// New creates a new health check server instance
func New(livenessAddr, readinessAddr string) *Server {
	hs := &Server{
		livenessAddr:  livenessAddr,
		readinessAddr: readinessAddr,
	}

	// Liveness probe - always returns 200 OK if the server is running
	livenessMux := http.NewServeMux()
	livenessMux.HandleFunc("/health/live", hs.handleLiveness)
	hs.livenessHTTP = &http.Server{
		Addr:              livenessAddr,
		Handler:           livenessMux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Readiness probe - returns 200 OK only when ready
	readinessMux := http.NewServeMux()
	readinessMux.HandleFunc("/health/ready", hs.handleReadiness)
	hs.readinessHTTP = &http.Server{
		Addr:              readinessAddr,
		Handler:           readinessMux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return hs
}

// Start begins both HTTP health check servers
func (hs *Server) Start() error {
	errCh := make(chan error, 2)

	// Start liveness probe server
	go func() {
		vlog.Infof("Liveness probe listening on %s", hs.livenessAddr)
		if err := hs.livenessHTTP.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("liveness server error: %w", err)
		}
	}()

	// Start readiness probe server
	go func() {
		vlog.Infof("Readiness probe listening on %s", hs.readinessAddr)
		if err := hs.readinessHTTP.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("readiness server error: %w", err)
		}
	}()

	// Check for immediate startup errors
	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}

// SetReady marks the service as ready
func (hs *Server) SetReady(ready bool) {
	hs.ready.Store(ready)
	if ready {
		vlog.Info("Service marked as ready")
	} else {
		vlog.Info("Service marked as not ready")
	}
}

// Shutdown gracefully shuts down both health check servers
func (hs *Server) Shutdown() error {
	if err := hs.livenessHTTP.Shutdown(context.TODO()); err != nil {
		return err
	}
	return hs.readinessHTTP.Shutdown(context.TODO())
}

// handleLiveness responds to liveness probe requests
func (hs *Server) handleLiveness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

// handleReadiness responds to readiness probe requests
func (hs *Server) handleReadiness(w http.ResponseWriter, r *http.Request) {
	if hs.ready.Load() {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Ready"))
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("Not Ready"))
	}
}
