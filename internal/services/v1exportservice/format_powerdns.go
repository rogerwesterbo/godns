package v1exportservice

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rogerwesterbo/godns/internal/models"
)

// PowerDNSZone represents a PowerDNS zone structure
type PowerDNSZone struct {
	Name   string          `json:"name"`
	Type   string          `json:"type"`
	Kind   string          `json:"kind"`
	DNSsec bool            `json:"dnssec"`
	Serial int             `json:"serial"`
	RRsets []PowerDNSRRset `json:"rrsets"`
}

// PowerDNSRRset represents a PowerDNS resource record set
type PowerDNSRRset struct {
	Name    string           `json:"name"`
	Type    string           `json:"type"`
	TTL     uint32           `json:"ttl"`
	Records []PowerDNSRecord `json:"records"`
}

// PowerDNSRecord represents a PowerDNS record
type PowerDNSRecord struct {
	Content  string `json:"content"`
	Disabled bool   `json:"disabled"`
}

// FormatPowerDNSZone formats a DNS zone for PowerDNS JSON API format
func FormatPowerDNSZone(zone *models.DNSZone) string {
	pdnsZone := PowerDNSZone{
		Name:   zone.Domain,
		Type:   "Zone",
		Kind:   "Master",
		DNSsec: false,
		Serial: 1,
		RRsets: make([]PowerDNSRRset, 0),
	}

	// Group records by name and type
	rrsetMap := make(map[string]*PowerDNSRRset)

	for _, record := range zone.Records {
		key := record.Name + ":" + record.Type

		if rrset, exists := rrsetMap[key]; exists {
			// Add to existing RRset
			rrset.Records = append(rrset.Records, PowerDNSRecord{
				Content:  record.Value,
				Disabled: false,
			})
		} else {
			// Create new RRset
			rrsetMap[key] = &PowerDNSRRset{
				Name: record.Name,
				Type: record.Type,
				TTL:  record.TTL,
				Records: []PowerDNSRecord{
					{
						Content:  record.Value,
						Disabled: false,
					},
				},
			}
		}
	}

	// Convert map to slice
	for _, rrset := range rrsetMap {
		pdnsZone.RRsets = append(pdnsZone.RRsets, *rrset)
	}

	// Add default SOA if not present
	hasSOA := false
	for _, rrset := range pdnsZone.RRsets {
		if rrset.Type == "SOA" {
			hasSOA = true
			break
		}
	}

	if !hasSOA {
		soaContent := fmt.Sprintf("ns1.%s hostmaster.%s 1 3600 1800 604800 300", zone.Domain, zone.Domain)
		pdnsZone.RRsets = append([]PowerDNSRRset{
			{
				Name: zone.Domain,
				Type: "SOA",
				TTL:  300,
				Records: []PowerDNSRecord{
					{
						Content:  soaContent,
						Disabled: false,
					},
				},
			},
		}, pdnsZone.RRsets...)
	}

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(pdnsZone, "", "  ")
	if err != nil {
		return fmt.Sprintf("# Error formatting PowerDNS zone: %v", err)
	}

	var sb strings.Builder
	sb.WriteString("# PowerDNS API Zone Configuration\n")
	sb.WriteString(fmt.Sprintf("# Zone: %s\n", zone.Domain))
	sb.WriteString("# Use this with: pdnsutil load-zone <zone-name> <file>\n")
	sb.WriteString("# Or via API: POST /api/v1/servers/localhost/zones\n\n")
	sb.WriteString(string(jsonData))

	return sb.String()
}
