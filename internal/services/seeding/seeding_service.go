package seeding

import (
	"context"
	"fmt"
	"net/netip"

	"github.com/rogerwesterbo/godns/internal/models"
	"github.com/rogerwesterbo/godns/internal/services/v1allowedlans"
	"github.com/rogerwesterbo/godns/internal/services/v1upstream"
	"github.com/rogerwesterbo/godns/internal/services/v1zoneservice"
	"github.com/rogerwesterbo/godns/pkg/consts"
	"github.com/spf13/viper"
	"github.com/vitistack/common/pkg/loggers/vlog"
)

// SeedingService handles initialization and seeding of configuration data
type SeedingService struct {
	allowedLANsService *v1allowedlans.AllowedLANsService
	upstreamService    *v1upstream.UpstreamService
	zoneService        *v1zoneservice.V1ZoneService
}

// NewSeedingService creates a new seeding service
func NewSeedingService(
	allowedLANsService *v1allowedlans.AllowedLANsService,
	upstreamService *v1upstream.UpstreamService,
	zoneService *v1zoneservice.V1ZoneService,
) *SeedingService {
	return &SeedingService{
		allowedLANsService: allowedLANsService,
		upstreamService:    upstreamService,
		zoneService:        zoneService,
	}
}

// SeedDefaults seeds all default configuration if not already present
// This is safe to call from multiple pods - it only seeds if keys don't exist
func (s *SeedingService) SeedDefaults(ctx context.Context, config SeedConfig) error {
	vlog.Info("starting configuration seeding...")

	// Seed test data if in development mode
	if viper.GetBool(consts.DEVELOPMENT) {
		// Seed allowed LANs
		if err := s.allowedLANsService.SeedDefaults(ctx, config.DefaultAllowedPrefixes); err != nil {
			return err
		}

		// Seed upstream DNS server
		if err := s.upstreamService.SeedDefault(ctx, config.DefaultUpstreamServer); err != nil {
			return err
		}

		vlog.Info("development mode detected - seeding test data...")
		if err := s.seedTestData(ctx); err != nil {
			vlog.Warnf("failed to seed test data: %v", err)
			// Don't fail the whole seeding process if test data fails
		} else {
			vlog.Info("test data seeded successfully")
		}
	}

	vlog.Info("configuration seeding completed successfully")
	return nil
}

// SeedConfig holds the default configuration values to seed
type SeedConfig struct {
	DefaultAllowedPrefixes []netip.Prefix
	DefaultUpstreamServer  string
}

// seedTestData seeds test zones and records for development/testing
func (s *SeedingService) seedTestData(ctx context.Context) error {
	// Check if test data already exists
	zones, err := s.zoneService.ListZones(ctx)
	if err != nil {
		return fmt.Errorf("failed to check existing zones: %w", err)
	}

	// Skip seeding if zones already exist
	if len(zones) > 0 {
		vlog.Infof("skipping test data seeding - %d zones already exist", len(zones))
		return nil
	}

	vlog.Info("seeding test zones...")

	// Seed Zone 1: home.lan (Home network)
	if err := s.zoneService.CreateZone(ctx, &models.DNSZone{
		Domain: "home.lan",
		Records: []models.DNSRecord{
			// Zone authority
			models.NewSOARecord("home.lan.", "ns1.home.lan.", "hostmaster.home.lan.", 2024110601, 3600, 1800, 604800, 300, 3600),
			models.NewNSRecord("home.lan.", "ns1.home.lan.", 3600),
			models.NewNSRecord("home.lan.", "ns2.home.lan.", 3600),
			// Name servers
			models.NewARecord("ns1.home.lan.", "192.168.1.1", 3600),
			models.NewARecord("ns2.home.lan.", "192.168.1.2", 3600),
			// Devices
			models.NewARecord("router.home.lan.", "192.168.1.1", 300),
			models.NewARecord("nas.home.lan.", "192.168.1.10", 300),
			models.NewAAAARecord("nas.home.lan.", "fd00::10", 300),
			models.NewARecord("server.home.lan.", "192.168.1.100", 300),
			models.NewARecord("printer.home.lan.", "192.168.1.50", 300),
			models.NewARecord("pi.home.lan.", "192.168.1.200", 300),
			// Aliases
			models.NewCNAMERecord("www.home.lan.", "server.home.lan.", 300),
			// Metadata
			models.NewTXTRecord("home.lan.", "v=spf1 ip4:192.168.1.0/24 -all", 300),
		},
	}); err != nil {
		return fmt.Errorf("failed to create home.lan zone: %w", err)
	}

	// Seed Zone 2: dev.local (Development environment)
	if err := s.zoneService.CreateZone(ctx, &models.DNSZone{
		Domain: "dev.local",
		Records: []models.DNSRecord{
			// Zone authority
			models.NewSOARecord("dev.local.", "ns.dev.local.", "hostmaster.dev.local.", 2024110601, 3600, 1800, 604800, 300, 3600),
			models.NewNSRecord("dev.local.", "ns.dev.local.", 3600),
			// Name server
			models.NewARecord("ns.dev.local.", "127.0.0.1", 3600),
			// Services
			models.NewARecord("api.dev.local.", "127.0.0.1", 300),
			models.NewARecord("db.dev.local.", "127.0.0.1", 300),
			models.NewARecord("cache.dev.local.", "127.0.0.1", 300),
			models.NewARecord("web.dev.local.", "127.0.0.1", 300),
			models.NewARecord("mail.dev.local.", "127.0.0.1", 300),
			// Mail configuration
			models.NewMXRecord("dev.local.", 10, "mail.dev.local.", 300),
			// Aliases
			models.NewCNAMERecord("www.dev.local.", "web.dev.local.", 300),
			models.NewCNAMERecord("admin.dev.local.", "web.dev.local.", 300),
			// Service discovery
			models.NewSRVRecord("_http._tcp.dev.local.", 10, 60, 80, "web.dev.local.", 300),
			// SPF record
			models.NewTXTRecord("dev.local.", "v=spf1 ip4:127.0.0.1 -all", 300),
		},
	}); err != nil {
		return fmt.Errorf("failed to create dev.local zone: %w", err)
	}

	// Seed Zone 3: k8s.local (Kubernetes cluster)
	if err := s.zoneService.CreateZone(ctx, &models.DNSZone{
		Domain: "k8s.local",
		Records: []models.DNSRecord{
			// Zone authority
			models.NewSOARecord("k8s.local.", "ns.k8s.local.", "hostmaster.k8s.local.", 2024110601, 3600, 1800, 604800, 300, 3600),
			models.NewNSRecord("k8s.local.", "ns.k8s.local.", 3600),
			// Name server
			models.NewARecord("ns.k8s.local.", "10.0.1.1", 3600),
			// Cluster nodes
			models.NewARecord("master.k8s.local.", "10.0.1.10", 300),
			models.NewARecord("worker1.k8s.local.", "10.0.1.20", 300),
			models.NewARecord("worker2.k8s.local.", "10.0.1.21", 300),
			models.NewARecord("worker3.k8s.local.", "10.0.1.22", 300),
			models.NewARecord("ingress.k8s.local.", "10.0.1.100", 300),
			// IPv6 support
			models.NewAAAARecord("master.k8s.local.", "fd00:10::10", 300),
			// Wildcard for apps
			models.NewARecord("*.apps.k8s.local.", "10.0.1.100", 300),
			// Alias
			models.NewCNAMERecord("api.k8s.local.", "master.k8s.local.", 300),
			// Service discovery for etcd
			models.NewSRVRecord("_etcd-server._tcp.k8s.local.", 10, 60, 2380, "master.k8s.local.", 300),
			// Cluster info
			models.NewTXTRecord("k8s.local.", "cluster=production version=1.28", 300),
		},
	}); err != nil {
		return fmt.Errorf("failed to create k8s.local zone: %w", err)
	}

	// Seed Zone 4: docker.local (Docker services)
	if err := s.zoneService.CreateZone(ctx, &models.DNSZone{
		Domain: "docker.local",
		Records: []models.DNSRecord{
			// Zone authority
			models.NewSOARecord("docker.local.", "ns.docker.local.", "hostmaster.docker.local.", 2024110601, 3600, 1800, 604800, 300, 3600),
			models.NewNSRecord("docker.local.", "ns.docker.local.", 3600),
			// Name server
			models.NewARecord("ns.docker.local.", "172.17.0.1", 3600),
			// Docker services
			models.NewARecord("portainer.docker.local.", "172.17.0.10", 300),
			models.NewARecord("registry.docker.local.", "172.17.0.20", 300),
			models.NewARecord("traefik.docker.local.", "172.17.0.30", 300),
			models.NewARecord("grafana.docker.local.", "172.17.0.40", 300),
			models.NewARecord("prometheus.docker.local.", "172.17.0.41", 300),
			// Aliases
			models.NewCNAMERecord("monitor.docker.local.", "grafana.docker.local.", 300),
			models.NewCNAMERecord("metrics.docker.local.", "prometheus.docker.local.", 300),
			// Service discovery
			models.NewSRVRecord("_metrics._tcp.docker.local.", 10, 100, 9090, "prometheus.docker.local.", 300),
			// Info
			models.NewTXTRecord("docker.local.", "network=bridge subnet=172.17.0.0/16", 300),
		},
	}); err != nil {
		return fmt.Errorf("failed to create docker.local zone: %w", err)
	}

	// Seed Zone 5: example.lan (Example/demo zone)
	if err := s.zoneService.CreateZone(ctx, &models.DNSZone{
		Domain: "example.lan",
		Records: []models.DNSRecord{
			// Zone authority
			models.NewSOARecord("example.lan.", "ns1.example.lan.", "hostmaster.example.lan.", 2024110601, 3600, 1800, 604800, 300, 3600),
			models.NewNSRecord("example.lan.", "ns1.example.lan.", 3600),
			models.NewNSRecord("example.lan.", "ns2.example.lan.", 3600),
			// Name servers
			models.NewARecord("ns1.example.lan.", "192.168.100.1", 3600),
			models.NewARecord("ns2.example.lan.", "192.168.100.2", 3600),
			// Web servers
			models.NewARecord("www.example.lan.", "192.168.100.10", 300),
			models.NewARecord("ftp.example.lan.", "192.168.100.30", 300),
			models.NewARecord("db.example.lan.", "192.168.100.40", 300),
			// Zone apex alias (ALIAS can be used at @ unlike CNAME)
			models.NewALIASRecord("example.lan.", "www.example.lan.", 300),
			// Mail servers
			models.NewARecord("mail.example.lan.", "192.168.100.20", 300),
			models.NewMXRecord("example.lan.", 10, "mail.example.lan.", 300),
			// Email security
			models.NewTXTRecord("example.lan.", "v=spf1 mx ip4:192.168.100.0/24 -all", 300),
			models.NewTXTRecord("_dmarc.example.lan.", "v=DMARC1; p=quarantine; rua=mailto:dmarc@example.lan", 300),
		},
	}); err != nil {
		return fmt.Errorf("failed to create example.lan zone: %w", err)
	}

	vlog.Info("successfully seeded 5 test zones with realistic DNS records")
	return nil
}
