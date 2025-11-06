package v1exportservice

import (
	"context"
	"fmt"

	"github.com/rogerwesterbo/godns/internal/models"
	"github.com/rogerwesterbo/godns/internal/services/v1zoneservice"
)

// ExportFormat represents the target DNS provider format
type ExportFormat string

const (
	FormatCoreDNS  ExportFormat = "coredns"
	FormatPowerDNS ExportFormat = "powerdns"
	FormatBIND     ExportFormat = "bind"
	FormatZoneFile ExportFormat = "zonefile" // Generic zone file format
)

// V1ExportService handles DNS zone export operations
type V1ExportService struct {
	zoneService *v1zoneservice.V1ZoneService
}

// NewV1ExportService creates a new export service
func NewV1ExportService(zoneService *v1zoneservice.V1ZoneService) *V1ExportService {
	return &V1ExportService{
		zoneService: zoneService,
	}
}

// ExportZone exports a single zone in the specified format
func (s *V1ExportService) ExportZone(ctx context.Context, domain string, format ExportFormat) (string, error) {
	zone, err := s.zoneService.GetZone(ctx, domain)
	if err != nil {
		return "", fmt.Errorf("failed to get zone: %w", err)
	}

	// Check if zone is enabled
	if !zone.Enabled {
		return "", fmt.Errorf("zone %s is disabled and cannot be exported", domain)
	}

	return s.formatZone(zone, format)
}

// ExportAllZones exports all zones in the specified format
// Only exports zones that are enabled
func (s *V1ExportService) ExportAllZones(ctx context.Context, format ExportFormat) (string, error) {
	zones, err := s.zoneService.ListZones(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list zones: %w", err)
	}

	result := ""
	exportedCount := 0
	for _, zone := range zones {
		// Skip disabled zones
		if !zone.Enabled {
			continue
		}

		formatted, err := s.formatZone(&zone, format)
		if err != nil {
			return "", fmt.Errorf("failed to format zone %s: %w", zone.Domain, err)
		}

		if exportedCount > 0 {
			result += "\n\n"
		}
		result += formatted
		exportedCount++
	}

	return result, nil
}

// formatZone formats a zone according to the specified format
func (s *V1ExportService) formatZone(zone *models.DNSZone, format ExportFormat) (string, error) {
	switch format {
	case FormatCoreDNS:
		return FormatCoreDNSZone(zone), nil
	case FormatPowerDNS:
		return FormatPowerDNSZone(zone), nil
	case FormatBIND, FormatZoneFile:
		return FormatBINDZone(zone), nil
	default:
		return "", fmt.Errorf("unsupported export format: %s", format)
	}
}

// ValidateFormat checks if the given format is supported
func ValidateFormat(format string) bool {
	switch ExportFormat(format) {
	case FormatCoreDNS, FormatPowerDNS, FormatBIND, FormatZoneFile:
		return true
	default:
		return false
	}
}
