package main

import (
	"context"
	"net/netip"
	"time"

	"github.com/rogerwesterbo/godns/internal/clients"
	"github.com/rogerwesterbo/godns/internal/dnsserver"
	"github.com/rogerwesterbo/godns/internal/dnsserver/handlers"
	"github.com/rogerwesterbo/godns/internal/httpserver"
	"github.com/rogerwesterbo/godns/internal/services/seeding"
	"github.com/rogerwesterbo/godns/internal/services/v1allowedlans"
	"github.com/rogerwesterbo/godns/internal/services/v1dnsservice"
	"github.com/rogerwesterbo/godns/internal/services/v1upstream"
	"github.com/rogerwesterbo/godns/internal/services/v1zoneservice"
	"github.com/rogerwesterbo/godns/internal/settings"
	"github.com/rogerwesterbo/godns/pkg/consts"
	"github.com/rogerwesterbo/godns/pkg/validation"
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

	vlog.Info("GoDNS starting...")

	clients.Init()

	// Create context for initialization operations
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Initialize DNS service
	dnsService := v1dnsservice.NewDNSService(clients.V1ValkeyClient)

	// Initialize allowed LANs service with Valkey backend
	allowedLANsService := v1allowedlans.NewAllowedLANsService(clients.V1ValkeyClient)

	// Initialize upstream DNS service with Valkey backend
	upstreamService := v1upstream.NewUpstreamService(clients.V1ValkeyClient, 3*time.Second)

	// Initialize zone service for HTTP API and seeding
	zoneService := v1zoneservice.NewV1ZoneService(clients.V1ValkeyClient)

	// Initialize and run seeding service
	seedingService := seeding.NewSeedingService(allowedLANsService, upstreamService, zoneService)

	// Get default upstream server from config or use Cloudflare
	defaultUpstream := viper.GetString(consts.DNS_UPSTREAM_SERVER)
	if defaultUpstream == "" {
		defaultUpstream = "1.1.1.1:53" // Cloudflare DNS
	}

	// Seed default configuration (safe for multiple pods - only seeds if not exists)
	seedConfig := seeding.SeedConfig{
		DefaultAllowedPrefixes: []netip.Prefix{
			netip.MustParsePrefix("192.168.0.0/16"), // Private IPv4
			netip.MustParsePrefix("10.0.0.0/8"),     // Private IPv4
			netip.MustParsePrefix("172.16.0.0/12"),  // Private IPv4
			netip.MustParsePrefix("fd00::/8"),       // IPv6 ULA
		},
		DefaultUpstreamServer: defaultUpstream,
	}
	if err := seedingService.SeedDefaults(ctx, seedConfig); err != nil {
		vlog.Fatalf("failed to seed configuration: %v", err)
	}

	// Create DNS handler with all services
	dnsHandler := handlers.NewDNSHandler(dnsService, allowedLANsService, upstreamService)

	createHttpServer := viper.GetBool(consts.DNS_ENABLE_HTTP_API)
	if createHttpServer {
		// Create and start the HTTP API server
		httpAPIAddress := viper.GetString(consts.HTTP_API_PORT)
		httpServer := httpserver.New(httpAPIAddress, zoneService)
		if err := httpServer.Start(); err != nil {
			vlog.Fatalf("failed to start HTTP API server: %v", err)
		}
	}

	// Create and start the DNS server
	dnsAddress := viper.GetString(consts.DNS_SERVER_PORT)
	livenessProbePort := viper.GetString(consts.DNS_SERVER_LIVENESS_PROBE_PORT)
	readinessProbePort := viper.GetString(consts.DNS_SERVER_READYNESS_PROBE_PORT)
	if err := validation.ValidateDNSAddress(livenessProbePort); err != nil {
		vlog.Fatalf("Invalid DNS server address: %v", err)
	}
	if err := validation.ValidateDNSAddress(readinessProbePort); err != nil {
		vlog.Fatalf("Invalid DNS server address: %v", err)
	}
	if err := validation.ValidateDNSAddress(dnsAddress); err != nil {
		vlog.Fatalf("Invalid DNS server address: %v", err)
	}

	server := dnsserver.New(dnsAddress, livenessProbePort, readinessProbePort, dnsHandler)
	if err := server.Start(); err != nil {
		vlog.Fatalf("server error: %v", err)
	}

	vlog.Info("GoDNS stopped.")
}
