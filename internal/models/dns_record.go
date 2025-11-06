package models

import "fmt"

// DNSRecord represents a DNS record stored in the system
// Each record type has common fields (Name, Type, TTL) plus type-specific fields
type DNSRecord struct {
	// Common fields for all record types
	Name     string `json:"name" example:"www.example.lan."` // Fully qualified domain name
	Type     string `json:"type" example:"A"`                // A, AAAA, MX, NS, TXT, SRV, SOA, CAA, CNAME, ALIAS, PTR
	TTL      uint32 `json:"ttl" example:"300"`               // Time to live in seconds
	Disabled bool   `json:"disabled"`                        // Whether the record is disabled (not served by DNS)

	// Simple value field (for A, AAAA, CNAME, ALIAS, NS, PTR, TXT)
	Value string `json:"value,omitempty" example:"192.168.1.100"` // The record value (IP, hostname, text, etc.)

	// MX record specific fields
	MXPriority *uint16 `json:"mx_priority,omitempty" example:"10"`            // Mail exchange priority (0-65535)
	MXHost     *string `json:"mx_host,omitempty" example:"mail.example.com."` // Mail server hostname

	// SRV record specific fields
	SRVPriority *uint16 `json:"srv_priority,omitempty" example:"10"`                // Priority (0-65535, lower is higher priority)
	SRVWeight   *uint16 `json:"srv_weight,omitempty" example:"60"`                  // Weight for load balancing (0-65535)
	SRVPort     *uint16 `json:"srv_port,omitempty" example:"443"`                   // Port number (0-65535)
	SRVTarget   *string `json:"srv_target,omitempty" example:"server.example.com."` // Target hostname

	// SOA record specific fields
	SOAMName   *string `json:"soa_mname,omitempty" example:"ns1.example.com."`        // Primary name server
	SOARName   *string `json:"soa_rname,omitempty" example:"hostmaster.example.com."` // Responsible person email
	SOASerial  *uint32 `json:"soa_serial,omitempty" example:"2024110601"`             // Serial number (YYYYMMDDnn)
	SOARefresh *uint32 `json:"soa_refresh,omitempty" example:"3600"`                  // Refresh interval (seconds)
	SOARetry   *uint32 `json:"soa_retry,omitempty" example:"1800"`                    // Retry interval (seconds)
	SOAExpire  *uint32 `json:"soa_expire,omitempty" example:"604800"`                 // Expire time (seconds)
	SOAMinimum *uint32 `json:"soa_minimum,omitempty" example:"300"`                   // Minimum TTL (seconds)

	// CAA record specific fields
	CAAFlags *uint8  `json:"caa_flags,omitempty" example:"0"`               // Flags (usually 0 or 128 for critical)
	CAATag   *string `json:"caa_tag,omitempty" example:"issue"`             // Property tag (issue, issuewild, iodef)
	CAAValue *string `json:"caa_value,omitempty" example:"letsencrypt.org"` // Property value
}

// DNSZone represents a DNS zone configuration
type DNSZone struct {
	Domain  string      `json:"domain" example:"example.lan."` // e.g., "example.lan."
	Records []DNSRecord `json:"records"`                       // DNS records in this zone
	Enabled bool        `json:"enabled"`                       // Whether the zone is enabled/active
}

// GetRData returns the RDATA (resource data) string for the DNS record
// This converts type-specific fields into the wire format string
func (r *DNSRecord) GetRData() string {
	switch r.Type {
	case "A", "AAAA", "CNAME", "ALIAS", "NS", "PTR", "TXT":
		return r.Value
	case "MX":
		if r.MXPriority != nil && r.MXHost != nil {
			return fmt.Sprintf("%d %s", *r.MXPriority, *r.MXHost)
		}
		return r.Value // Fallback for backward compatibility
	case "SRV":
		if r.SRVPriority != nil && r.SRVWeight != nil && r.SRVPort != nil && r.SRVTarget != nil {
			return fmt.Sprintf("%d %d %d %s", *r.SRVPriority, *r.SRVWeight, *r.SRVPort, *r.SRVTarget)
		}
		return r.Value // Fallback for backward compatibility
	case "SOA":
		if r.SOAMName != nil && r.SOARName != nil && r.SOASerial != nil &&
			r.SOARefresh != nil && r.SOARetry != nil && r.SOAExpire != nil && r.SOAMinimum != nil {
			return fmt.Sprintf("%s %s %d %d %d %d %d",
				*r.SOAMName, *r.SOARName, *r.SOASerial,
				*r.SOARefresh, *r.SOARetry, *r.SOAExpire, *r.SOAMinimum)
		}
		return r.Value // Fallback for backward compatibility
	case "CAA":
		if r.CAAFlags != nil && r.CAATag != nil && r.CAAValue != nil {
			return fmt.Sprintf("%d %s %q", *r.CAAFlags, *r.CAATag, *r.CAAValue)
		}
		return r.Value // Fallback for backward compatibility
	default:
		return r.Value
	}
}

// Validate checks if the DNS record has all required fields for its type
func (r *DNSRecord) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("record name is required")
	}
	if r.Type == "" {
		return fmt.Errorf("record type is required")
	}

	switch r.Type {
	case "A", "AAAA", "CNAME", "ALIAS", "NS", "PTR", "TXT":
		if r.Value == "" {
			return fmt.Errorf("%s record requires a value", r.Type)
		}
	case "MX":
		// Accept either new format or legacy value field
		if r.MXPriority == nil || r.MXHost == nil {
			if r.Value == "" {
				return fmt.Errorf("MX record requires either mx_priority+mx_host or value field")
			}
		}
	case "SRV":
		if r.SRVPriority == nil || r.SRVWeight == nil || r.SRVPort == nil || r.SRVTarget == nil {
			if r.Value == "" {
				return fmt.Errorf("SRV record requires either srv_priority+srv_weight+srv_port+srv_target or value field")
			}
		}
	case "SOA":
		if r.SOAMName == nil || r.SOARName == nil || r.SOASerial == nil ||
			r.SOARefresh == nil || r.SOARetry == nil || r.SOAExpire == nil || r.SOAMinimum == nil {
			if r.Value == "" {
				return fmt.Errorf("SOA record requires all soa_* fields or value field")
			}
		}
	case "CAA":
		if r.CAAFlags == nil || r.CAATag == nil || r.CAAValue == nil {
			if r.Value == "" {
				return fmt.Errorf("CAA record requires either caa_flags+caa_tag+caa_value or value field")
			}
		}
	}

	return nil
}

// Helper functions to create common record types

// NewARecord creates an A record
func NewARecord(name, ipv4 string, ttl uint32) DNSRecord {
	return DNSRecord{
		Name:  name,
		Type:  "A",
		TTL:   ttl,
		Value: ipv4,
	}
}

// NewAAAARecord creates an AAAA record
func NewAAAARecord(name, ipv6 string, ttl uint32) DNSRecord {
	return DNSRecord{
		Name:  name,
		Type:  "AAAA",
		TTL:   ttl,
		Value: ipv6,
	}
}

// NewCNAMERecord creates a CNAME record
func NewCNAMERecord(name, target string, ttl uint32) DNSRecord {
	return DNSRecord{
		Name:  name,
		Type:  "CNAME",
		TTL:   ttl,
		Value: target,
	}
}

// NewALIASRecord creates an ALIAS record
// ALIAS records are like CNAME but can be used at the zone apex
func NewALIASRecord(name, target string, ttl uint32) DNSRecord {
	return DNSRecord{
		Name:  name,
		Type:  "ALIAS",
		TTL:   ttl,
		Value: target,
	}
}

// NewNSRecord creates an NS record
func NewNSRecord(name, nameserver string, ttl uint32) DNSRecord {
	return DNSRecord{
		Name:  name,
		Type:  "NS",
		TTL:   ttl,
		Value: nameserver,
	}
}

// NewMXRecord creates an MX record
func NewMXRecord(name string, priority uint16, mailserver string, ttl uint32) DNSRecord {
	return DNSRecord{
		Name:       name,
		Type:       "MX",
		TTL:        ttl,
		MXPriority: &priority,
		MXHost:     &mailserver,
	}
}

// NewSRVRecord creates an SRV record
func NewSRVRecord(name string, priority, weight, port uint16, target string, ttl uint32) DNSRecord {
	return DNSRecord{
		Name:        name,
		Type:        "SRV",
		TTL:         ttl,
		SRVPriority: &priority,
		SRVWeight:   &weight,
		SRVPort:     &port,
		SRVTarget:   &target,
	}
}

// NewSOARecord creates an SOA record
func NewSOARecord(name, mname, rname string, serial, refresh, retry, expire, minimum uint32, ttl uint32) DNSRecord {
	return DNSRecord{
		Name:       name,
		Type:       "SOA",
		TTL:        ttl,
		SOAMName:   &mname,
		SOARName:   &rname,
		SOASerial:  &serial,
		SOARefresh: &refresh,
		SOARetry:   &retry,
		SOAExpire:  &expire,
		SOAMinimum: &minimum,
	}
}

// NewTXTRecord creates a TXT record
func NewTXTRecord(name, text string, ttl uint32) DNSRecord {
	return DNSRecord{
		Name:  name,
		Type:  "TXT",
		TTL:   ttl,
		Value: text,
	}
}

// NewCAARecord creates a CAA record
func NewCAARecord(name string, flags uint8, tag, value string, ttl uint32) DNSRecord {
	return DNSRecord{
		Name:     name,
		Type:     "CAA",
		TTL:      ttl,
		CAAFlags: &flags,
		CAATag:   &tag,
		CAAValue: &value,
	}
}
