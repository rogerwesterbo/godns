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
			{Name: "router.home.lan.", Type: "A", TTL: 300, Value: "192.168.1.1"},
			{Name: "nas.home.lan.", Type: "A", TTL: 300, Value: "192.168.1.10"},
			{Name: "server.home.lan.", Type: "A", TTL: 300, Value: "192.168.1.100"},
			{Name: "printer.home.lan.", Type: "A", TTL: 300, Value: "192.168.1.50"},
			{Name: "pi.home.lan.", Type: "A", TTL: 300, Value: "192.168.1.200"},
		},
	}); err != nil {
		return fmt.Errorf("failed to create home.lan zone: %w", err)
	}

	// Seed Zone 2: dev.local (Development environment)
	if err := s.zoneService.CreateZone(ctx, &models.DNSZone{
		Domain: "dev.local",
		Records: []models.DNSRecord{
			{Name: "api.dev.local.", Type: "A", TTL: 300, Value: "127.0.0.1"},
			{Name: "db.dev.local.", Type: "A", TTL: 300, Value: "127.0.0.1"},
			{Name: "cache.dev.local.", Type: "A", TTL: 300, Value: "127.0.0.1"},
			{Name: "web.dev.local.", Type: "A", TTL: 300, Value: "127.0.0.1"},
			{Name: "mail.dev.local.", Type: "A", TTL: 300, Value: "127.0.0.1"},
		},
	}); err != nil {
		return fmt.Errorf("failed to create dev.local zone: %w", err)
	}

	// Seed Zone 3: k8s.local (Kubernetes cluster)
	if err := s.zoneService.CreateZone(ctx, &models.DNSZone{
		Domain: "k8s.local",
		Records: []models.DNSRecord{
			{Name: "master.k8s.local.", Type: "A", TTL: 300, Value: "10.0.1.10"},
			{Name: "worker1.k8s.local.", Type: "A", TTL: 300, Value: "10.0.1.20"},
			{Name: "worker2.k8s.local.", Type: "A", TTL: 300, Value: "10.0.1.21"},
			{Name: "worker3.k8s.local.", Type: "A", TTL: 300, Value: "10.0.1.22"},
			{Name: "ingress.k8s.local.", Type: "A", TTL: 300, Value: "10.0.1.100"},
		},
	}); err != nil {
		return fmt.Errorf("failed to create k8s.local zone: %w", err)
	}

	// Seed Zone 4: docker.local (Docker services)
	if err := s.zoneService.CreateZone(ctx, &models.DNSZone{
		Domain: "docker.local",
		Records: []models.DNSRecord{
			{Name: "portainer.docker.local.", Type: "A", TTL: 300, Value: "172.17.0.10"},
			{Name: "registry.docker.local.", Type: "A", TTL: 300, Value: "172.17.0.20"},
			{Name: "traefik.docker.local.", Type: "A", TTL: 300, Value: "172.17.0.30"},
			{Name: "grafana.docker.local.", Type: "A", TTL: 300, Value: "172.17.0.40"},
			{Name: "prometheus.docker.local.", Type: "A", TTL: 300, Value: "172.17.0.41"},
		},
	}); err != nil {
		return fmt.Errorf("failed to create docker.local zone: %w", err)
	}

	// Seed Zone 5: example.lan (Example/demo zone)
	if err := s.zoneService.CreateZone(ctx, &models.DNSZone{
		Domain: "example.lan",
		Records: []models.DNSRecord{
			{Name: "www.example.lan.", Type: "A", TTL: 300, Value: "192.168.100.10"},
			{Name: "mail.example.lan.", Type: "A", TTL: 300, Value: "192.168.100.20"},
			{Name: "ftp.example.lan.", Type: "A", TTL: 300, Value: "192.168.100.30"},
			{Name: "db.example.lan.", Type: "A", TTL: 300, Value: "192.168.100.40"},
			{Name: "ns1.example.lan.", Type: "A", TTL: 3600, Value: "192.168.100.1"},
			{Name: "ns2.example.lan.", Type: "A", TTL: 3600, Value: "192.168.100.2"},
		},
	}); err != nil {
		return fmt.Errorf("failed to create example.lan zone: %w", err)
	}

	vlog.Info("successfully seeded 5 test zones with 26 DNS records")
	return nil
}
