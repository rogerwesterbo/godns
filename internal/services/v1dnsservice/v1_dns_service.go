package v1dnsservice

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/miekg/dns"
	"github.com/rogerwesterbo/godns/internal/models"
	"github.com/rogerwesterbo/godns/pkg/interfaces/valkeyinterface"
)

const (
	// Valkey key prefix for DNS zones
	zoneKeyPrefix = "dns:zone:"
	// Valkey key prefix for DNS records
	recordKeyPrefix = "dns:record:"
)

// DNSService handles DNS record management and storage
type DNSService struct {
	valkeyClient valkeyinterface.ValkeyInterface
}

// NewDNSService creates a new DNS service
func NewDNSService(valkeyClient valkeyinterface.ValkeyInterface) *DNSService {
	return &DNSService{
		valkeyClient: valkeyClient,
	}
}

// GetRecord retrieves a DNS record from Valkey
func (s *DNSService) GetRecord(ctx context.Context, name string, recordType string) (*models.DNSRecord, error) {
	key := s.buildRecordKey(name, recordType)

	data, err := s.valkeyClient.GetData(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get DNS record: %w", err)
	}

	var record models.DNSRecord
	if err := json.Unmarshal([]byte(data), &record); err != nil {
		return nil, fmt.Errorf("failed to unmarshal DNS record: %w", err)
	}

	return &record, nil
}

// SetRecord stores a DNS record in Valkey
func (s *DNSService) SetRecord(ctx context.Context, record *models.DNSRecord) error {
	key := s.buildRecordKey(record.Name, record.Type)

	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal DNS record: %w", err)
	}

	if err := s.valkeyClient.SetData(ctx, key, string(data)); err != nil {
		return fmt.Errorf("failed to set DNS record: %w", err)
	}

	return nil
}

// DeleteRecord removes a DNS record from Valkey
func (s *DNSService) DeleteRecord(ctx context.Context, name string, recordType string) error {
	key := s.buildRecordKey(name, recordType)

	if err := s.valkeyClient.DeleteData(ctx, key); err != nil {
		return fmt.Errorf("failed to delete DNS record: %w", err)
	}

	return nil
}

// GetZone retrieves all records for a zone
func (s *DNSService) GetZone(ctx context.Context, domain string) (*models.DNSZone, error) {
	key := zoneKeyPrefix + domain

	data, err := s.valkeyClient.GetData(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get DNS zone: %w", err)
	}

	var zone models.DNSZone
	if err := json.Unmarshal([]byte(data), &zone); err != nil {
		return nil, fmt.Errorf("failed to unmarshal DNS zone: %w", err)
	}

	return &zone, nil
}

// SetZone stores a complete zone configuration
func (s *DNSService) SetZone(ctx context.Context, zone *models.DNSZone) error {
	key := zoneKeyPrefix + zone.Domain

	data, err := json.Marshal(zone)
	if err != nil {
		return fmt.Errorf("failed to marshal DNS zone: %w", err)
	}

	if err := s.valkeyClient.SetData(ctx, key, string(data)); err != nil {
		return fmt.Errorf("failed to set DNS zone: %w", err)
	}

	// Also store individual records for faster lookup
	for _, record := range zone.Records {
		if err := s.SetRecord(ctx, &record); err != nil {
			return fmt.Errorf("failed to set record in zone: %w", err)
		}
	}

	return nil
}

// LookupRecord performs a DNS lookup and returns DNS resource records
func (s *DNSService) LookupRecord(ctx context.Context, name string, qtype uint16) ([]dns.RR, error) {
	recordType := dns.TypeToString[qtype]

	record, err := s.GetRecord(ctx, name, recordType)
	if err != nil {
		return nil, err
	}

	// Convert the record to DNS RR format
	rr, err := s.convertToRR(record)
	if err != nil {
		return nil, fmt.Errorf("failed to convert record to RR: %w", err)
	}

	return []dns.RR{rr}, nil
}

// HasZone checks if a domain is managed by this DNS server
func (s *DNSService) HasZone(ctx context.Context, name string) (string, bool) {
	// Try to find if the name belongs to any configured zone
	// We'll check common zone patterns
	parts := strings.Split(strings.TrimSuffix(name, "."), ".")

	// Try progressively larger domain parts
	for i := 0; i < len(parts); i++ {
		domain := strings.Join(parts[i:], ".") + "."

		_, err := s.GetZone(ctx, domain)
		if err == nil {
			return domain, true
		}
	}

	return "", false
}

// buildRecordKey creates a Valkey key for a DNS record
func (s *DNSService) buildRecordKey(name string, recordType string) string {
	return fmt.Sprintf("%s%s:%s", recordKeyPrefix, name, recordType)
}

// convertToRR converts a DNSRecord model to a dns.RR
func (s *DNSService) convertToRR(record *models.DNSRecord) (dns.RR, error) {
	// Build the RR string format: "name TTL class type value"
	rrString := fmt.Sprintf("%s %d IN %s %s", record.Name, record.TTL, record.Type, record.Value)

	rr, err := dns.NewRR(rrString)
	if err != nil {
		return nil, fmt.Errorf("failed to create RR from string: %w", err)
	}

	return rr, nil
}
