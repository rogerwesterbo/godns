package models

// DNSRecord represents a DNS record stored in the system
type DNSRecord struct {
	Name  string `json:"name" example:"www.example.lan."` // Fully qualified domain name
	Type  string `json:"type" example:"A"`                // A, AAAA, MX, NS, TXT, etc.
	TTL   uint32 `json:"ttl" example:"300"`               // Time to live in seconds
	Value string `json:"value" example:"192.168.1.100"`   // The record value (IP, hostname, etc.)
}

// DNSZone represents a DNS zone configuration
type DNSZone struct {
	Domain  string      `json:"domain" example:"example.lan."` // e.g., "example.lan."
	Records []DNSRecord `json:"records"`                       // DNS records in this zone
}
