package v1zoneservice

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rogerwesterbo/godns/internal/models"
	"github.com/rogerwesterbo/godns/pkg/interfaces/valkeyinterface"
)

const (
	zoneKeyPrefix   = "zone:"
	zoneListKey     = "zones:list"
	recordKeyPrefix = "record:"
)

// V1ZoneService handles DNS zone and record operations
type V1ZoneService struct {
	client valkeyinterface.ValkeyInterface
}

// NewV1ZoneService creates a new zone service
func NewV1ZoneService(client valkeyinterface.ValkeyInterface) *V1ZoneService {
	return &V1ZoneService{
		client: client,
	}
}

// CreateZone creates a new DNS zone
func (s *V1ZoneService) CreateZone(ctx context.Context, zone *models.DNSZone) error {
	if zone.Domain == "" {
		return fmt.Errorf("zone domain cannot be empty")
	}

	// Normalize domain to end with a dot
	if !strings.HasSuffix(zone.Domain, ".") {
		zone.Domain += "."
	}

	// Check if zone already exists
	zoneKey := zoneKeyPrefix + zone.Domain
	_, err := s.client.GetData(ctx, zoneKey)
	if err == nil {
		return fmt.Errorf("zone %s already exists", zone.Domain)
	}

	// Validate records
	for _, record := range zone.Records {
		if err := s.validateRecord(&record); err != nil {
			return fmt.Errorf("invalid record: %w", err)
		}
	}

	// Save zone metadata
	zoneData, err := json.Marshal(zone)
	if err != nil {
		return fmt.Errorf("failed to marshal zone: %w", err)
	}

	if err := s.client.SetData(ctx, zoneKey, string(zoneData)); err != nil {
		return fmt.Errorf("failed to save zone: %w", err)
	}

	// Add zone to the list of zones
	zones, err := s.listZoneDomains(ctx)
	if err != nil && !strings.Contains(err.Error(), "key not found") {
		return fmt.Errorf("failed to get zone list: %w", err)
	}

	// Add domain if not already in list
	found := false
	for _, z := range zones {
		if z == zone.Domain {
			found = true
			break
		}
	}
	if !found {
		zones = append(zones, zone.Domain)
		zonesData, err := json.Marshal(zones)
		if err != nil {
			return fmt.Errorf("failed to marshal zone list: %w", err)
		}
		if err := s.client.SetData(ctx, zoneListKey, string(zonesData)); err != nil {
			return fmt.Errorf("failed to save zone list: %w", err)
		}
	}

	// Save individual records
	for _, record := range zone.Records {
		if err := s.saveRecord(ctx, zone.Domain, &record); err != nil {
			return fmt.Errorf("failed to save record: %w", err)
		}
	}

	return nil
}

// GetZone retrieves a DNS zone by domain
func (s *V1ZoneService) GetZone(ctx context.Context, domain string) (*models.DNSZone, error) {
	if !strings.HasSuffix(domain, ".") {
		domain += "."
	}

	zoneKey := zoneKeyPrefix + domain
	data, err := s.client.GetData(ctx, zoneKey)
	if err != nil {
		return nil, fmt.Errorf("zone not found: %w", err)
	}

	var zone models.DNSZone
	if err := json.Unmarshal([]byte(data), &zone); err != nil {
		return nil, fmt.Errorf("failed to unmarshal zone: %w", err)
	}

	return &zone, nil
}

// ListZones returns all DNS zones
func (s *V1ZoneService) ListZones(ctx context.Context) ([]models.DNSZone, error) {
	domains, err := s.listZoneDomains(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "key not found") {
			return []models.DNSZone{}, nil
		}
		return nil, err
	}

	zones := make([]models.DNSZone, 0, len(domains))
	for _, domain := range domains {
		zone, err := s.GetZone(ctx, domain)
		if err != nil {
			// Skip zones that can't be retrieved
			continue
		}
		zones = append(zones, *zone)
	}

	return zones, nil
}

// UpdateZone updates an existing DNS zone
func (s *V1ZoneService) UpdateZone(ctx context.Context, domain string, zone *models.DNSZone) error {
	if !strings.HasSuffix(domain, ".") {
		domain += "."
	}

	// Check if zone exists
	_, err := s.GetZone(ctx, domain)
	if err != nil {
		return fmt.Errorf("zone not found: %w", err)
	}

	// Update domain to match the key
	zone.Domain = domain

	// Validate records
	for _, record := range zone.Records {
		if err := s.validateRecord(&record); err != nil {
			return fmt.Errorf("invalid record: %w", err)
		}
	}

	// Delete old records for this zone
	if err := s.deleteZoneRecords(ctx, domain); err != nil {
		return fmt.Errorf("failed to delete old records: %w", err)
	}

	// Save updated zone
	zoneKey := zoneKeyPrefix + domain
	zoneData, err := json.Marshal(zone)
	if err != nil {
		return fmt.Errorf("failed to marshal zone: %w", err)
	}

	if err := s.client.SetData(ctx, zoneKey, string(zoneData)); err != nil {
		return fmt.Errorf("failed to update zone: %w", err)
	}

	// Save new records
	for _, record := range zone.Records {
		if err := s.saveRecord(ctx, domain, &record); err != nil {
			return fmt.Errorf("failed to save record: %w", err)
		}
	}

	return nil
}

// DeleteZone deletes a DNS zone and all its records
func (s *V1ZoneService) DeleteZone(ctx context.Context, domain string) error {
	if !strings.HasSuffix(domain, ".") {
		domain += "."
	}

	// Check if zone exists
	_, err := s.GetZone(ctx, domain)
	if err != nil {
		return fmt.Errorf("zone not found: %w", err)
	}

	// Delete all records for this zone
	if err := s.deleteZoneRecords(ctx, domain); err != nil {
		return fmt.Errorf("failed to delete zone records: %w", err)
	}

	// Delete zone metadata
	zoneKey := zoneKeyPrefix + domain
	if err := s.client.DeleteData(ctx, zoneKey); err != nil {
		return fmt.Errorf("failed to delete zone: %w", err)
	}

	// Remove from zone list
	zones, err := s.listZoneDomains(ctx)
	if err != nil {
		return fmt.Errorf("failed to get zone list: %w", err)
	}

	newZones := make([]string, 0, len(zones))
	for _, z := range zones {
		if z != domain {
			newZones = append(newZones, z)
		}
	}

	zonesData, err := json.Marshal(newZones)
	if err != nil {
		return fmt.Errorf("failed to marshal zone list: %w", err)
	}

	if err := s.client.SetData(ctx, zoneListKey, string(zonesData)); err != nil {
		return fmt.Errorf("failed to update zone list: %w", err)
	}

	return nil
}

// CreateRecord adds a new record to a zone
func (s *V1ZoneService) CreateRecord(ctx context.Context, domain string, record *models.DNSRecord) error {
	if !strings.HasSuffix(domain, ".") {
		domain += "."
	}

	// Check if zone exists
	zone, err := s.GetZone(ctx, domain)
	if err != nil {
		return fmt.Errorf("zone not found: %w", err)
	}

	// Validate record
	if err := s.validateRecord(record); err != nil {
		return err
	}

	// Check if record already exists
	for _, r := range zone.Records {
		if r.Name == record.Name && r.Type == record.Type {
			return fmt.Errorf("record %s of type %s already exists in zone", record.Name, record.Type)
		}
	}

	// Add record to zone
	zone.Records = append(zone.Records, *record)

	// Update zone
	zoneData, err := json.Marshal(zone)
	if err != nil {
		return fmt.Errorf("failed to marshal zone: %w", err)
	}

	zoneKey := zoneKeyPrefix + domain
	if err := s.client.SetData(ctx, zoneKey, string(zoneData)); err != nil {
		return fmt.Errorf("failed to update zone: %w", err)
	}

	// Save individual record
	if err := s.saveRecord(ctx, domain, record); err != nil {
		return fmt.Errorf("failed to save record: %w", err)
	}

	return nil
}

// GetRecord retrieves a specific record from a zone
func (s *V1ZoneService) GetRecord(ctx context.Context, domain, name, recordType string) (*models.DNSRecord, error) {
	if !strings.HasSuffix(domain, ".") {
		domain += "."
	}

	zone, err := s.GetZone(ctx, domain)
	if err != nil {
		return nil, fmt.Errorf("zone not found: %w", err)
	}

	for _, record := range zone.Records {
		if record.Name == name && record.Type == recordType {
			return &record, nil
		}
	}

	return nil, fmt.Errorf("record not found")
}

// UpdateRecord updates an existing record in a zone
func (s *V1ZoneService) UpdateRecord(ctx context.Context, domain, name, recordType string, record *models.DNSRecord) error {
	if !strings.HasSuffix(domain, ".") {
		domain += "."
	}

	// Get zone
	zone, err := s.GetZone(ctx, domain)
	if err != nil {
		return fmt.Errorf("zone not found: %w", err)
	}

	// Validate record
	if err := s.validateRecord(record); err != nil {
		return err
	}

	// Find and update the record
	found := false
	for i, r := range zone.Records {
		if r.Name == name && r.Type == recordType {
			zone.Records[i] = *record
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("record not found")
	}

	// Update zone
	zoneData, err := json.Marshal(zone)
	if err != nil {
		return fmt.Errorf("failed to marshal zone: %w", err)
	}

	zoneKey := zoneKeyPrefix + domain
	if err := s.client.SetData(ctx, zoneKey, string(zoneData)); err != nil {
		return fmt.Errorf("failed to update zone: %w", err)
	}

	// Delete old record
	oldRecordKey := recordKeyPrefix + domain + ":" + name + ":" + recordType
	_ = s.client.DeleteData(ctx, oldRecordKey)

	// Save updated record
	if err := s.saveRecord(ctx, domain, record); err != nil {
		return fmt.Errorf("failed to save record: %w", err)
	}

	return nil
}

// DeleteRecord deletes a record from a zone
func (s *V1ZoneService) DeleteRecord(ctx context.Context, domain, name, recordType string) error {
	if !strings.HasSuffix(domain, ".") {
		domain += "."
	}

	// Get zone
	zone, err := s.GetZone(ctx, domain)
	if err != nil {
		return fmt.Errorf("zone not found: %w", err)
	}

	// Find and remove the record
	found := false
	newRecords := make([]models.DNSRecord, 0, len(zone.Records))
	for _, r := range zone.Records {
		if r.Name == name && r.Type == recordType {
			found = true
			continue
		}
		newRecords = append(newRecords, r)
	}

	if !found {
		return fmt.Errorf("record not found")
	}

	zone.Records = newRecords

	// Update zone
	zoneData, err := json.Marshal(zone)
	if err != nil {
		return fmt.Errorf("failed to marshal zone: %w", err)
	}

	zoneKey := zoneKeyPrefix + domain
	if err := s.client.SetData(ctx, zoneKey, string(zoneData)); err != nil {
		return fmt.Errorf("failed to update zone: %w", err)
	}

	// Delete individual record
	recordKey := recordKeyPrefix + domain + ":" + name + ":" + recordType
	if err := s.client.DeleteData(ctx, recordKey); err != nil {
		return fmt.Errorf("failed to delete record: %w", err)
	}

	return nil
}

// Helper functions

func (s *V1ZoneService) listZoneDomains(ctx context.Context) ([]string, error) {
	data, err := s.client.GetData(ctx, zoneListKey)
	if err != nil {
		return nil, err
	}

	var zones []string
	if err := json.Unmarshal([]byte(data), &zones); err != nil {
		return nil, fmt.Errorf("failed to unmarshal zone list: %w", err)
	}

	return zones, nil
}

func (s *V1ZoneService) saveRecord(ctx context.Context, domain string, record *models.DNSRecord) error {
	recordKey := recordKeyPrefix + domain + ":" + record.Name + ":" + record.Type
	recordData, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal record: %w", err)
	}

	return s.client.SetData(ctx, recordKey, string(recordData))
}

func (s *V1ZoneService) deleteZoneRecords(ctx context.Context, domain string) error {
	zone, err := s.GetZone(ctx, domain)
	if err != nil {
		return err
	}

	for _, record := range zone.Records {
		recordKey := recordKeyPrefix + domain + ":" + record.Name + ":" + record.Type
		_ = s.client.DeleteData(ctx, recordKey) // Ignore errors for individual records
	}

	return nil
}

func (s *V1ZoneService) validateRecord(record *models.DNSRecord) error {
	if record.Name == "" {
		return fmt.Errorf("record name cannot be empty")
	}
	if record.Type == "" {
		return fmt.Errorf("record type cannot be empty")
	}
	if record.Value == "" {
		return fmt.Errorf("record value cannot be empty")
	}

	// Validate record type
	validTypes := map[string]bool{
		"A": true, "AAAA": true, "CNAME": true, "MX": true,
		"NS": true, "TXT": true, "PTR": true, "SRV": true,
		"SOA": true, "CAA": true,
	}
	if !validTypes[strings.ToUpper(record.Type)] {
		return fmt.Errorf("invalid record type: %s", record.Type)
	}

	// Normalize type to uppercase
	record.Type = strings.ToUpper(record.Type)

	// Set default TTL if not specified
	if record.TTL == 0 {
		record.TTL = 300 // 5 minutes default
	}

	return nil
}
