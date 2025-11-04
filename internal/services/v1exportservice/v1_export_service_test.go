package v1exportservice

import (
	"testing"

	"github.com/rogerwesterbo/godns/internal/models"
)

func TestFormatBINDZone(t *testing.T) {
	zone := &models.DNSZone{
		Domain: "example.lan.",
		Records: []models.DNSRecord{
			{Name: "example.lan.", Type: "A", TTL: 300, Value: "192.168.1.1"},
			{Name: "www.example.lan.", Type: "A", TTL: 300, Value: "192.168.1.2"},
			{Name: "mail.example.lan.", Type: "A", TTL: 300, Value: "192.168.1.3"},
			{Name: "example.lan.", Type: "MX", TTL: 300, Value: "10 mail.example.lan."},
		},
	}

	result := FormatBINDZone(zone)

	// Check that result contains expected elements
	if result == "" {
		t.Fatal("Expected non-empty BIND zone file")
	}

	// Should contain zone origin
	if !contains(result, "$ORIGIN example.lan.") {
		t.Error("Expected $ORIGIN directive in BIND zone file")
	}

	// Should contain TTL
	if !contains(result, "$TTL") {
		t.Error("Expected $TTL directive in BIND zone file")
	}

	// Should contain records
	if !contains(result, "A\t192.168.1.1") {
		t.Error("Expected A record in BIND zone file")
	}

	if !contains(result, "MX") {
		t.Error("Expected MX record in BIND zone file")
	}

	t.Logf("BIND Zone Output:\n%s", result)
}

func TestFormatCoreDNSZone(t *testing.T) {
	zone := &models.DNSZone{
		Domain: "example.lan.",
		Records: []models.DNSRecord{
			{Name: "example.lan.", Type: "A", TTL: 300, Value: "192.168.1.1"},
			{Name: "www.example.lan.", Type: "A", TTL: 300, Value: "192.168.1.2"},
		},
	}

	result := FormatCoreDNSZone(zone)

	// Check that result contains expected elements
	if result == "" {
		t.Fatal("Expected non-empty CoreDNS configuration")
	}

	// Should contain zone name
	if !contains(result, "example.lan {") {
		t.Error("Expected zone block in CoreDNS configuration")
	}

	// Should contain file plugin
	if !contains(result, "file /etc/coredns/zones/") {
		t.Error("Expected file plugin configuration")
	}

	// Should contain zone origin
	if !contains(result, "$ORIGIN example.lan.") {
		t.Error("Expected $ORIGIN directive")
	}

	t.Logf("CoreDNS Output:\n%s", result)
}

func TestFormatPowerDNSZone(t *testing.T) {
	zone := &models.DNSZone{
		Domain: "example.lan.",
		Records: []models.DNSRecord{
			{Name: "example.lan.", Type: "A", TTL: 300, Value: "192.168.1.1"},
			{Name: "www.example.lan.", Type: "A", TTL: 300, Value: "192.168.1.2"},
			{Name: "www.example.lan.", Type: "A", TTL: 300, Value: "192.168.1.3"}, // Multiple A records
		},
	}

	result := FormatPowerDNSZone(zone)

	// Check that result contains expected elements
	if result == "" {
		t.Fatal("Expected non-empty PowerDNS JSON")
	}

	// Should be valid JSON structure
	if !contains(result, "\"name\"") {
		t.Error("Expected JSON structure with 'name' field")
	}

	if !contains(result, "\"rrsets\"") {
		t.Error("Expected 'rrsets' field in PowerDNS JSON")
	}

	// Should contain zone name
	if !contains(result, "example.lan.") {
		t.Error("Expected zone domain in PowerDNS JSON")
	}

	t.Logf("PowerDNS Output:\n%s", result)
}

func TestValidateFormat(t *testing.T) {
	tests := []struct {
		format string
		valid  bool
	}{
		{"coredns", true},
		{"powerdns", true},
		{"bind", true},
		{"zonefile", true},
		{"invalid", false},
		{"", false},
		{"BIND", false}, // Case sensitive
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.format, func(t *testing.T) {
			result := ValidateFormat(tt.format)
			if result != tt.valid {
				t.Errorf("ValidateFormat(%q) = %v, want %v", tt.format, result, tt.valid)
			}
		})
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
