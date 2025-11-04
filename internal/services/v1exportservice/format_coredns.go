package v1exportservice

import (
	"fmt"
	"strings"

	"github.com/rogerwesterbo/godns/internal/models"
)

// FormatCoreDNSZone formats a DNS zone for CoreDNS Corefile format
// CoreDNS uses a plugin-based configuration with the file plugin
func FormatCoreDNSZone(zone *models.DNSZone) string {
	var sb strings.Builder

	// CoreDNS Corefile block for this zone
	_, _ = sb.WriteString(fmt.Sprintf("# Zone: %s\n", zone.Domain))
	_, _ = sb.WriteString(fmt.Sprintf("%s {\n", strings.TrimSuffix(zone.Domain, ".")))
	_, _ = sb.WriteString("    file /etc/coredns/zones/db." + strings.TrimSuffix(zone.Domain, ".") + "\n")
	_, _ = sb.WriteString("    log\n")
	_, _ = sb.WriteString("    errors\n")
	_, _ = sb.WriteString("}\n\n")

	// Zone file content
	_, _ = sb.WriteString(fmt.Sprintf("# Zone file: db.%s\n", strings.TrimSuffix(zone.Domain, ".")))
	_, _ = sb.WriteString(fmt.Sprintf("$ORIGIN %s\n", zone.Domain))
	_, _ = sb.WriteString("$TTL 300\n\n")

	// Add SOA record if exists, otherwise create a default one
	hasSOA := false
	for _, record := range zone.Records {
		if record.Type == "SOA" {
			sb.WriteString(fmt.Sprintf("@\t%d\tIN\tSOA\t%s\n", record.TTL, record.Value))
			hasSOA = true
			break
		}
	}

	if !hasSOA {
		// Default SOA record
		_, _ = sb.WriteString(fmt.Sprintf("@\t300\tIN\tSOA\tns1.%s hostmaster.%s 1 3600 1800 604800 300\n", zone.Domain, zone.Domain))
	}
	_, _ = sb.WriteString("\n")

	// Add other records
	for _, record := range zone.Records {
		if record.Type == "SOA" {
			continue // Already handled
		}

		name := record.Name
		if name == zone.Domain {
			name = "@"
		} else if strings.HasSuffix(name, "."+zone.Domain) {
			// Make it relative to the zone
			name = strings.TrimSuffix(name, "."+zone.Domain)
		}

		_, _ = sb.WriteString(fmt.Sprintf("%s\t%d\tIN\t%s\t%s\n",
			name, record.TTL, record.Type, record.Value))
	}

	return sb.String()
}
