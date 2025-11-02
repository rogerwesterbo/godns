package v1upstream

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/miekg/dns"
	"github.com/rogerwesterbo/godns/pkg/interfaces/valkeyinterface"
	"github.com/vitistack/common/pkg/loggers/vlog"
)

const upstreamConfigKey = "dns:config:upstream"

// UpstreamService handles forwarding DNS queries to upstream servers
type UpstreamService struct {
	valkeyClient valkeyinterface.ValkeyInterface
	upstreamAddr string
	timeout      time.Duration
}

// UpstreamConfig represents the stored upstream configuration
type UpstreamConfig struct {
	Address string `json:"address"`
}

// NewUpstreamService creates a new upstream service
func NewUpstreamService(valkeyClient valkeyinterface.ValkeyInterface, timeout time.Duration) *UpstreamService {
	return &UpstreamService{
		valkeyClient: valkeyClient,
		timeout:      timeout,
	}
}

// SeedDefault seeds default upstream server if none exists in Valkey
// This is safe for multiple pods - only seeds if key doesn't exist
func (s *UpstreamService) SeedDefault(ctx context.Context, defaultUpstream string) error {
	// Check if configuration already exists
	_, err := s.valkeyClient.GetData(ctx, upstreamConfigKey)
	if err == nil {
		// Configuration exists, load it
		vlog.Info("upstream DNS configuration already exists, loading from Valkey")
		return s.LoadFromValkey(ctx)
	}

	// Configuration doesn't exist, seed default
	vlog.Infof("seeding default upstream DNS server: %s", defaultUpstream)

	config := UpstreamConfig{
		Address: defaultUpstream,
	}

	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal upstream config: %w", err)
	}

	if err := s.valkeyClient.SetData(ctx, upstreamConfigKey, string(data)); err != nil {
		return fmt.Errorf("failed to seed upstream config: %w", err)
	}

	s.upstreamAddr = defaultUpstream
	vlog.Infof("seeded upstream DNS server: %s", defaultUpstream)
	return nil
}

// LoadFromValkey loads the upstream configuration from Valkey
func (s *UpstreamService) LoadFromValkey(ctx context.Context) error {
	data, err := s.valkeyClient.GetData(ctx, upstreamConfigKey)
	if err != nil {
		return fmt.Errorf("failed to get upstream config from Valkey: %w", err)
	}

	var config UpstreamConfig
	if err := json.Unmarshal([]byte(data), &config); err != nil {
		return fmt.Errorf("failed to unmarshal upstream config: %w", err)
	}

	s.upstreamAddr = config.Address
	vlog.Infof("loaded upstream DNS server from Valkey: %s", s.upstreamAddr)
	return nil
}

// Forward forwards a DNS query to the upstream server
func (s *UpstreamService) Forward(ctx context.Context, query *dns.Msg) (*dns.Msg, error) {
	c := &dns.Client{
		Net:            "udp",
		Timeout:        s.timeout,
		Dialer:         &net.Dialer{Timeout: s.timeout},
		SingleInflight: true,
	}

	type result struct {
		msg *dns.Msg
		err error
	}

	ch := make(chan result, 1)
	go func() {
		in, _, err := c.ExchangeContext(ctx, query, s.upstreamAddr)
		// Fallback to TCP if truncated
		if err == nil && in != nil && in.Truncated {
			cTCP := &dns.Client{Net: "tcp", Timeout: s.timeout}
			in, _, err = cTCP.ExchangeContext(ctx, query, s.upstreamAddr)
		}
		ch <- result{in, err}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case r := <-ch:
		return r.msg, r.err
	}
}

// SetUpstream changes the upstream DNS server address and saves to Valkey
func (s *UpstreamService) SetUpstream(ctx context.Context, addr string) error {
	config := UpstreamConfig{
		Address: addr,
	}

	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal upstream config: %w", err)
	}

	if err := s.valkeyClient.SetData(ctx, upstreamConfigKey, string(data)); err != nil {
		return fmt.Errorf("failed to save upstream config: %w", err)
	}

	s.upstreamAddr = addr
	vlog.Infof("updated upstream DNS server: %s", addr)
	return nil
}

// GetUpstream returns the current upstream DNS server address
func (s *UpstreamService) GetUpstream() string {
	return s.upstreamAddr
}
