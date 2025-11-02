package models

// DNSRecord represents a DNS record stored in the system
type DNSRecord struct {
	Name  string `json:"name"`  // Fully qualified domain name
	Type  string `json:"type"`  // A, AAAA, MX, NS, TXT, etc.
	TTL   uint32 `json:"ttl"`   // Time to live in seconds
	Value string `json:"value"` // The record value (IP, hostname, etc.)
}

// DNSZone represents a DNS zone configuration
type DNSZone struct {
	Domain  string      `json:"domain"`  // e.g., "example.lan."
	Records []DNSRecord `json:"records"` // DNS records in this zone
}
