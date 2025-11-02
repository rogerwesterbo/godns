package seeding

import (
	"context"
	"net/netip"

	"github.com/rogerwesterbo/godns/internal/services/v1allowedlans"
	"github.com/rogerwesterbo/godns/internal/services/v1upstream"
	"github.com/vitistack/common/pkg/loggers/vlog"
)

// SeedingService handles initialization and seeding of configuration data
type SeedingService struct {
	allowedLANsService *v1allowedlans.AllowedLANsService
	upstreamService    *v1upstream.UpstreamService
}

// NewSeedingService creates a new seeding service
func NewSeedingService(
	allowedLANsService *v1allowedlans.AllowedLANsService,
	upstreamService *v1upstream.UpstreamService,
) *SeedingService {
	return &SeedingService{
		allowedLANsService: allowedLANsService,
		upstreamService:    upstreamService,
	}
}

// SeedDefaults seeds all default configuration if not already present
// This is safe to call from multiple pods - it only seeds if keys don't exist
func (s *SeedingService) SeedDefaults(ctx context.Context, config SeedConfig) error {
	vlog.Info("starting configuration seeding...")

	// Seed allowed LANs
	if err := s.allowedLANsService.SeedDefaults(ctx, config.DefaultAllowedPrefixes); err != nil {
		return err
	}

	// Seed upstream DNS server
	if err := s.upstreamService.SeedDefault(ctx, config.DefaultUpstreamServer); err != nil {
		return err
	}

	vlog.Info("configuration seeding completed successfully")
	return nil
}

// SeedConfig holds the default configuration values to seed
type SeedConfig struct {
	DefaultAllowedPrefixes []netip.Prefix
	DefaultUpstreamServer  string
}
