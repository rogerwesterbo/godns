package v1dnssecservice

import (
	"testing"

	"github.com/miekg/dns"
	"github.com/rogerwesterbo/godns/internal/models"
)

func TestDNSSECService_SignZone(t *testing.T) {
	service := NewDNSSECService()
	domain := "example.com."

	// 1. Generate Keys
	ksk, zsk, kskPriv, zskPriv, err := service.GenerateKeys(domain)
	if err != nil {
		t.Fatalf("GenerateKeys failed: %v", err)
	}

	if ksk == nil || zsk == nil {
		t.Fatal("Generated keys are nil")
	}

	// 2. Create a test zone
	zone := &models.DNSZone{
		Domain: domain,
		Records: []models.DNSRecord{
			{
				Name:  domain,
				Type:  "A",
				TTL:   3600,
				Value: "1.2.3.4",
			},
			{
				Name:  "www." + domain,
				Type:  "A",
				TTL:   3600,
				Value: "1.2.3.5",
			},
		},
	}

	// 3. Sign the zone
	signedRecords, err := service.SignZone(zone, ksk, zsk, kskPriv, zskPriv)
	if err != nil {
		t.Fatalf("SignZone failed: %v", err)
	}

	// 4. Verify results
	if len(signedRecords) == 0 {
		t.Fatal("No signed records returned")
	}

	// Check for presence of different record types
	hasDNSKEY := false
	hasRRSIG := false
	hasNSEC := false
	hasA := false

	for _, rr := range signedRecords {
		switch rr.Header().Rrtype {
		case dns.TypeDNSKEY:
			hasDNSKEY = true
		case dns.TypeRRSIG:
			hasRRSIG = true
		case dns.TypeNSEC:
			hasNSEC = true
		case dns.TypeA:
			hasA = true
		}
	}

	if !hasDNSKEY {
		t.Error("Missing DNSKEY records")
	}
	if !hasRRSIG {
		t.Error("Missing RRSIG records")
	}
	if !hasNSEC {
		t.Error("Missing NSEC records")
	}
	if !hasA {
		t.Error("Missing A records")
	}

	// Optional: Print records for visual inspection
	t.Logf("Generated %d records:", len(signedRecords))
	for _, rr := range signedRecords {
		t.Logf("%s", rr.String())
	}
}
