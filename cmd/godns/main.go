package main

import (
	"context"
	"net/http"
	"net/netip"
	"time"

	"github.com/rogerwesterbo/godns/internal/clients"
	"github.com/rogerwesterbo/godns/internal/dnsserver"
	"github.com/rogerwesterbo/godns/internal/dnsserver/handlers"
	"github.com/rogerwesterbo/godns/internal/httpserver"
	"github.com/rogerwesterbo/godns/internal/services/seeding"
	"github.com/rogerwesterbo/godns/internal/services/v1allowedlans"
	"github.com/rogerwesterbo/godns/internal/services/v1cacheservice"
	"github.com/rogerwesterbo/godns/internal/services/v1dnsservice"
	"github.com/rogerwesterbo/godns/internal/services/v1healthcheckservice"
	"github.com/rogerwesterbo/godns/internal/services/v1loadbalancerservice"
	"github.com/rogerwesterbo/godns/internal/services/v1metricsservice"
	"github.com/rogerwesterbo/godns/internal/services/v1querylogservice"
	"github.com/rogerwesterbo/godns/internal/services/v1ratelimitservice"
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

	// Initialize DNS cache service
	var cacheService *v1cacheservice.DNSCache
	if viper.GetBool(consts.DNS_CACHE_ENABLED) {
		cacheSize := viper.GetInt(consts.DNS_CACHE_SIZE)
		cacheTTL := time.Duration(viper.GetInt(consts.DNS_CACHE_TTL_SECONDS)) * time.Second
		cacheService = v1cacheservice.NewDNSCache(cacheSize, cacheTTL)
		vlog.Infof("DNS cache enabled (size: %d, TTL: %s)", cacheSize, cacheTTL)
	}

	// Initialize rate limiter
	var rateLimiter *v1ratelimitservice.RateLimiter
	if viper.GetBool(consts.DNS_RATE_LIMIT_ENABLED) {
		ratePerSec := viper.GetInt(consts.DNS_RATE_LIMIT_QUERIES_PER_SEC)
		burst := viper.GetInt(consts.DNS_RATE_LIMIT_BURST)
		rateLimiter = v1ratelimitservice.NewRateLimiter(ratePerSec, burst)
		vlog.Infof("Rate limiting enabled (%d queries/sec, burst: %d)", ratePerSec, burst)
	}

	// Initialize load balancer
	var loadBalancer *v1loadbalancerservice.LoadBalancer
	if viper.GetBool(consts.DNS_LOAD_BALANCER_ENABLED) {
		strategy := v1loadbalancerservice.RoundRobin // Default strategy
		strategyStr := viper.GetString(consts.DNS_LOAD_BALANCER_STRATEGY)
		switch strategyStr {
		case "weighted-round-robin":
			strategy = v1loadbalancerservice.WeightedRoundRobin
		case "least-connections":
			strategy = v1loadbalancerservice.LeastConnections
		case "random":
			strategy = v1loadbalancerservice.Random
		}
		loadBalancer = v1loadbalancerservice.NewLoadBalancer(strategy)
		vlog.Infof("Load balancer enabled (strategy: %s)", strategyStr)
	}

	// Initialize health check service
	var healthCheckService *v1healthcheckservice.HealthCheckService
	if viper.GetBool(consts.DNS_HEALTH_CHECK_ENABLED) {
		healthCheckService = v1healthcheckservice.NewHealthCheckService()
		vlog.Info("Health check service enabled")
	}

	// Initialize query logging
	var queryLogService *v1querylogservice.QueryLogService
	if viper.GetBool(consts.DNS_QUERY_LOG_ENABLED) {
		bufferSize := viper.GetInt(consts.DNS_QUERY_LOG_BUFFER_SIZE)
		flushInterval, _ := time.ParseDuration(viper.GetString(consts.DNS_QUERY_LOG_FLUSH_INTERVAL))
		if flushInterval == 0 {
			flushInterval = 1 * time.Minute
		}
		queryLogService = v1querylogservice.NewQueryLogService(bufferSize, flushInterval, clients.V1ValkeyClient)
		vlog.Infof("Query logging enabled (buffer: %d, flush: %s)", bufferSize, flushInterval)
	}

	// Initialize metrics service
	var metricsService *v1metricsservice.MetricsService
	if viper.GetBool(consts.METRICS_ENABLED) {
		metricsService = v1metricsservice.NewMetricsService()

		// Start metrics HTTP server
		metricsPort := viper.GetString(consts.METRICS_PORT)
		go func() {
			mux := http.NewServeMux()
			mux.Handle("/metrics", metricsService.Handler())

			// Create server with timeouts for security
			metricsServer := &http.Server{
				Addr:         metricsPort,
				Handler:      mux,
				ReadTimeout:  15 * time.Second,
				WriteTimeout: 15 * time.Second,
				IdleTimeout:  60 * time.Second,
			}

			vlog.Infof("Starting metrics server on %s", metricsPort)
			if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				vlog.Errorf("Metrics server error: %v", err)
			}
		}()
	}

	// Initialize and run seeding service
	seedingService := seeding.NewSeedingService(allowedLANsService, upstreamService, zoneService)

	// Get default upstream server from config, system, or fallback to Cloudflare
	defaultUpstream := viper.GetString(consts.DNS_UPSTREAM_SERVER)
	if defaultUpstream == "" {
		defaultUpstream = "1.1.1.1:53"
		vlog.Infof("Using Cloudflare DNS as upstream fallback: %s", defaultUpstream)
	} else {
		vlog.Infof("Using configured upstream DNS: %s", defaultUpstream)
	}

	// Seed default configuration (safe for multiple pods - only seeds if not exists)
	seedConfig := seeding.SeedConfig{
		DefaultAllowedPrefixes: []netip.Prefix{
			netip.MustParsePrefix("127.0.0.0/8"),    // Localhost IPv4
			netip.MustParsePrefix("::1/128"),        // Localhost IPv6
			netip.MustParsePrefix("192.168.0.0/16"), // Private IPv4
			netip.MustParsePrefix("10.0.0.0/8"),     // Private IPv4
			netip.MustParsePrefix("172.16.0.0/12"),  // Private IPv4 (Docker default)
			netip.MustParsePrefix("fd00::/8"),       // IPv6 ULA
		},
		DefaultUpstreamServer: defaultUpstream,
	}
	if err := seedingService.SeedDefaults(ctx, seedConfig); err != nil {
		vlog.Fatalf("failed to seed configuration: %v", err)
	}

	// Create DNS handler with all services
	dnsHandler := handlers.NewDNSHandler(
		dnsService,
		allowedLANsService,
		upstreamService,
		cacheService,
		rateLimiter,
		loadBalancer,
		healthCheckService,
		queryLogService,
		metricsService,
	)

	createHttpServer := viper.GetBool(consts.DNS_ENABLE_HTTP_API)
	if createHttpServer {
		// Create and start the HTTP API server
		httpAPIAddress := viper.GetString(consts.HTTP_API_PORT)
		httpServer, err := httpserver.New(
			httpAPIAddress,
			zoneService,
			cacheService,
			rateLimiter,
			loadBalancer,
			healthCheckService,
			queryLogService,
		)
		if err != nil {
			vlog.Fatalf("failed to create HTTP API server: %v", err)
		}
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
