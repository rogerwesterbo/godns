package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rogerwesterbo/godns/internal/clients"
	"github.com/rogerwesterbo/godns/internal/healthserver"
	"github.com/rogerwesterbo/godns/internal/httpserver"
	"github.com/rogerwesterbo/godns/internal/services/v1zoneservice"
	"github.com/rogerwesterbo/godns/internal/settings"
	"github.com/rogerwesterbo/godns/pkg/consts"
	"github.com/spf13/viper"
	"github.com/vitistack/common/pkg/loggers/vlog"
)

func main() {
	settings.Init()

	vlogOpts := vlog.Options{
		Level:             viper.GetString(consts.LOG_LEVEL),    // debug|info|warn|error|dpanic|panic|fatal
		JSON:              viper.GetBool(consts.LOG_JSON),       // default: structured JSON (fastest to parse)
		AddCaller:         viper.GetBool(consts.LOG_ADD_CALLER), // include caller file:line
		DisableStacktrace: viper.GetBool(consts.LOG_DISABLE_STACKTRACE),
		ColorizeLine:      viper.GetBool(consts.LOG_COLORIZE_LINE),      // set true only for human console viewing
		UnescapeMultiline: viper.GetBool(consts.LOG_UNESCAPE_MULTILINE), // set true only if you need pretty multi-line msg rendering in text mode
	}
	_ = vlog.Setup(vlogOpts)
	defer func() {
		_ = vlog.Sync()
	}()

	vlog.Info("GoDNS API Server starting...")

	// Initialize clients
	clients.Init()

	// Initialize zone service for HTTP API
	zoneService := v1zoneservice.NewV1ZoneService(clients.V1ValkeyClient)

	// Create and start health check servers
	livenessProbePort := viper.GetString(consts.HTTP_API_LIVENESS_PROBE_PORT)
	readinessProbePort := viper.GetString(consts.HTTP_API_READINESS_PROBE_PORT)
	healthServer := healthserver.New(livenessProbePort, readinessProbePort)

	// Register Valkey health check
	healthServer.AddHealthCheck(clients.V1ValkeyClient.Ping)

	if err := healthServer.Start(); err != nil {
		vlog.Fatalf("failed to start health check servers: %v", err)
	}

	// Mark service as ready after initialization
	healthServer.SetReady(true)

	// Create and start the HTTP API server
	httpAPIAddress := viper.GetString(consts.HTTP_API_PORT)
	httpServer := httpserver.New(httpAPIAddress, zoneService)
	if err := httpServer.Start(); err != nil {
		vlog.Fatalf("failed to start HTTP API server: %v", err)
	}

	// Wait for interrupt signal for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	vlog.Info("GoDNS API Server is running. Press Ctrl+C to stop.")
	<-sigChan

	// Graceful shutdown
	vlog.Info("Shutting down HTTP API server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Mark service as not ready during shutdown
	healthServer.SetReady(false)

	if err := httpServer.Stop(shutdownCtx); err != nil {
		vlog.Errorf("Error during HTTP server shutdown: %v", err)
	}

	if err := healthServer.Shutdown(); err != nil {
		vlog.Errorf("Error during health server shutdown: %v", err)
	}

	vlog.Info("GoDNS API Server stopped.")
}
