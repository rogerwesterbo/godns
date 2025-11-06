package v1allowedlans

import (
	"context"
	"encoding/json"
	"fmt"
	"net/netip"

	"github.com/rogerwesterbo/godns/pkg/interfaces/valkeyinterface"
	"github.com/vitistack/common/pkg/loggers/vlog"
)

const allowedLANsKey = "dns:config:allowedlans"

// AllowedLANsService manages allowed LAN networks for DNS recursion
type AllowedLANsService struct {
	valkeyClient valkeyinterface.ValkeyInterface
	prefixes     []netip.Prefix // cached in memory
}

// AllowedLANsConfig represents the stored configuration
type AllowedLANsConfig struct {
	Prefixes []string `json:"prefixes"`
}

// NewAllowedLANsService creates a new allowed LANs service
func NewAllowedLANsService(valkeyClient valkeyinterface.ValkeyInterface) *AllowedLANsService {
	return &AllowedLANsService{
		valkeyClient: valkeyClient,
		prefixes:     []netip.Prefix{},
	}
}

// SeedDefaults seeds default allowed LANs if none exist in Valkey
// This is safe for multiple pods - only seeds if key doesn't exist
func (s *AllowedLANsService) SeedDefaults(ctx context.Context, defaultPrefixes []netip.Prefix) error {
	// Check if configuration already exists
	_, err := s.valkeyClient.GetData(ctx, allowedLANsKey)
	if err == nil {
		// Configuration exists, load it
		vlog.Info("allowed LANs configuration already exists, loading from Valkey")
		return s.LoadFromValkey(ctx)
	}

	// Configuration doesn't exist, seed defaults
	vlog.Info("seeding default allowed LANs configuration")

	prefixStrings := make([]string, len(defaultPrefixes))
	for i, prefix := range defaultPrefixes {
		prefixStrings[i] = prefix.String()
	}

	config := AllowedLANsConfig{
		Prefixes: prefixStrings,
	}

	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal allowed LANs config: %w", err)
	}

	if err := s.valkeyClient.SetData(ctx, allowedLANsKey, string(data)); err != nil {
		return fmt.Errorf("failed to seed allowed LANs: %w", err)
	}

	s.prefixes = defaultPrefixes
	vlog.Infof("seeded %d default allowed LAN prefixes", len(defaultPrefixes))
	return nil
}

// LoadFromValkey loads the allowed LANs configuration from Valkey
func (s *AllowedLANsService) LoadFromValkey(ctx context.Context) error {
	data, err := s.valkeyClient.GetData(ctx, allowedLANsKey)
	if err != nil {
		return fmt.Errorf("failed to get allowed LANs from Valkey: %w", err)
	}

	var config AllowedLANsConfig
	if err := json.Unmarshal([]byte(data), &config); err != nil {
		return fmt.Errorf("failed to unmarshal allowed LANs config: %w", err)
	}

	prefixes := make([]netip.Prefix, 0, len(config.Prefixes))
	for _, prefixStr := range config.Prefixes {
		prefix, err := netip.ParsePrefix(prefixStr)
		if err != nil {
			vlog.Warnf("invalid prefix in configuration: %s: %v", prefixStr, err)
			continue
		}
		prefixes = append(prefixes, prefix)
	}

	s.prefixes = prefixes
	vlog.Infof("loaded %d allowed LAN prefixes from Valkey", len(prefixes))
	return nil
}

// IsAllowed checks if an IP address is in an allowed LAN
func (s *AllowedLANsService) IsAllowed(ip netip.Addr) bool {
	if !ip.IsValid() {
		return false
	}

	// Normalize IPv4-mapped IPv6 addresses (::ffff:192.0.2.1) to IPv4 (192.0.2.1)
	// This ensures Docker containers using IPv6 stack can match IPv4 prefixes
	if ip.Is4In6() {
		ip = ip.Unmap()
	}

	for _, prefix := range s.prefixes {
		if prefix.Contains(ip) {
			return true
		}
	}

	return false
}

// AddPrefix adds a new allowed network prefix and saves to Valkey
func (s *AllowedLANsService) AddPrefix(ctx context.Context, prefix netip.Prefix) error {
	s.prefixes = append(s.prefixes, prefix)

	prefixStrings := make([]string, len(s.prefixes))
	for i, p := range s.prefixes {
		prefixStrings[i] = p.String()
	}

	config := AllowedLANsConfig{
		Prefixes: prefixStrings,
	}

	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal allowed LANs config: %w", err)
	}

	if err := s.valkeyClient.SetData(ctx, allowedLANsKey, string(data)); err != nil {
		return fmt.Errorf("failed to save allowed LANs: %w", err)
	}

	return nil
}

// GetPrefixes returns all allowed network prefixes
func (s *AllowedLANsService) GetPrefixes() []netip.Prefix {
	return s.prefixes
}
