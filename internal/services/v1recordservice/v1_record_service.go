package v1recordservice

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
	recordKeyPrefix = "record:"
)

// V1RecordService handles DNS record operations
type V1RecordService struct {
	client valkeyinterface.ValkeyInterface
}

// NewV1RecordService creates a new record service
func NewV1RecordService(client valkeyinterface.ValkeyInterface) *V1RecordService {
	return &V1RecordService{
		client: client,
	}
}

// CreateRecord adds a new record to a zone
func (s *V1RecordService) CreateRecord(ctx context.Context, domain string, record *models.DNSRecord) error {
	if !strings.HasSuffix(domain, ".") {
		domain += "."
	}

	// Check if zone exists
	zone, err := s.getZone(ctx, domain)
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
func (s *V1RecordService) GetRecord(ctx context.Context, domain, name, recordType string) (*models.DNSRecord, error) {
	if !strings.HasSuffix(domain, ".") {
		domain += "."
	}

	zone, err := s.getZone(ctx, domain)
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
func (s *V1RecordService) UpdateRecord(ctx context.Context, domain, name, recordType string, record *models.DNSRecord) error {
	if !strings.HasSuffix(domain, ".") {
		domain += "."
	}

	// Get zone
	zone, err := s.getZone(ctx, domain)
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
func (s *V1RecordService) DeleteRecord(ctx context.Context, domain, name, recordType string) error {
	if !strings.HasSuffix(domain, ".") {
		domain += "."
	}

	// Get zone
	zone, err := s.getZone(ctx, domain)
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

// SetRecordEnabled sets the enabled status of a DNS record
func (s *V1RecordService) SetRecordEnabled(ctx context.Context, domain, name, recordType string, enabled bool) error {
	if !strings.HasSuffix(domain, ".") {
		domain += "."
	}

	// Load zone so the in-zone record listing stays in sync
	zone, err := s.getZone(ctx, domain)
	if err != nil {
		return fmt.Errorf("zone not found: %w", err)
	}

	var updatedRecord *models.DNSRecord
	for i := range zone.Records {
		if zone.Records[i].Name == name && zone.Records[i].Type == recordType {
			zone.Records[i].Disabled = !enabled
			updatedRecord = &zone.Records[i]
			break
		}
	}

	if updatedRecord == nil {
		return fmt.Errorf("record not found")
	}

	// Persist updated zone metadata
	zoneData, err := json.Marshal(zone)
	if err != nil {
		return fmt.Errorf("failed to marshal zone: %w", err)
	}

	zoneKey := zoneKeyPrefix + domain
	if err := s.client.SetData(ctx, zoneKey, string(zoneData)); err != nil {
		return fmt.Errorf("failed to update zone: %w", err)
	}

	// Persist individual record entry
	if err := s.saveRecord(ctx, domain, updatedRecord); err != nil {
		return fmt.Errorf("failed to update record status: %w", err)
	}

	return nil
}

// Helper functions

// getZone retrieves a zone from storage
func (s *V1RecordService) getZone(ctx context.Context, domain string) (*models.DNSZone, error) {
	zoneKey := zoneKeyPrefix + domain
	data, err := s.client.GetData(ctx, zoneKey)
	if err != nil {
		return nil, fmt.Errorf("zone not found: %w", err)
	}

	var zone models.DNSZone
	if err := json.Unmarshal([]byte(data), &zone); err != nil {
		return nil, fmt.Errorf("failed to unmarshal zone: %w", err)
	}

	// Backward compatibility: if Enabled field is not set in JSON, default to true
	if data != "" && !strings.Contains(data, `"enabled"`) {
		zone.Enabled = true
	}

	return &zone, nil
}

// saveRecord saves a record to storage
func (s *V1RecordService) saveRecord(ctx context.Context, domain string, record *models.DNSRecord) error {
	recordKey := recordKeyPrefix + domain + ":" + record.Name + ":" + record.Type
	recordData, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal record: %w", err)
	}

	return s.client.SetData(ctx, recordKey, string(recordData))
}

// validateRecord validates a DNS record
func (s *V1RecordService) validateRecord(record *models.DNSRecord) error {
	// Normalize type to uppercase
	record.Type = strings.ToUpper(record.Type)

	// Set default TTL if not specified
	if record.TTL == 0 {
		record.TTL = 300 // 5 minutes default
	}

	// Validate record type
	validTypes := map[string]bool{
		"A": true, "AAAA": true, "CNAME": true, "ALIAS": true, "MX": true,
		"NS": true, "TXT": true, "PTR": true, "SRV": true,
		"SOA": true, "CAA": true,
	}
	if !validTypes[record.Type] {
		return fmt.Errorf("invalid record type: %s", record.Type)
	}

	// Use the model's built-in validation for type-specific fields
	if err := record.Validate(); err != nil {
		return err
	}

	return nil
}
