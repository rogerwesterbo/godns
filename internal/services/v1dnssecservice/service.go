package v1dnssecservice

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/miekg/dns"
	"github.com/rogerwesterbo/godns/internal/models"
)

type DNSSECService struct{}

func NewDNSSECService() *DNSSECService {
	return &DNSSECService{}
}

// GenerateKeys generates a new KSK and ZSK for the given domain
func (s *DNSSECService) GenerateKeys(domain string) (*dns.DNSKEY, *dns.DNSKEY, *ecdsa.PrivateKey, *ecdsa.PrivateKey, error) {
	// Generate KSK
	kskPriv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to generate KSK: %w", err)
	}

	ksk := &dns.DNSKEY{
		Hdr:       dns.RR_Header{Name: dns.Fqdn(domain), Rrtype: dns.TypeDNSKEY, Class: dns.ClassINET, Ttl: 3600},
		Algorithm: dns.ECDSAP256SHA256,
		Flags:     257, // KSK (SEP + Zone Key)
		Protocol:  3,
	}
	s.setPublicKey(ksk, &kskPriv.PublicKey)

	// Generate ZSK
	zskPriv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to generate ZSK: %w", err)
	}

	zsk := &dns.DNSKEY{
		Hdr:       dns.RR_Header{Name: dns.Fqdn(domain), Rrtype: dns.TypeDNSKEY, Class: dns.ClassINET, Ttl: 3600},
		Algorithm: dns.ECDSAP256SHA256,
		Flags:     256, // ZSK (Zone Key)
		Protocol:  3,
	}
	s.setPublicKey(zsk, &zskPriv.PublicKey)

	return ksk, zsk, kskPriv, zskPriv, nil
}

func (s *DNSSECService) setPublicKey(k *dns.DNSKEY, pub *ecdsa.PublicKey) {
	// ECDSA P-256 public key is 64 bytes (32 bytes X + 32 bytes Y)
	buf := make([]byte, 64)
	x := pub.X.Bytes()
	y := pub.Y.Bytes()

	// Pad X and Y to 32 bytes if needed
	copy(buf[32-len(x):32], x)
	copy(buf[64-len(y):64], y)

	k.PublicKey = base64.StdEncoding.EncodeToString(buf)
}

// SignZone signs the zone records and returns the signed zone including DNSKEY, RRSIG, and NSEC records
func (s *DNSSECService) SignZone(zone *models.DNSZone, ksk, zsk *dns.DNSKEY, kskPriv, zskPriv *ecdsa.PrivateKey) ([]dns.RR, error) {
	var records []dns.RR

	// Convert model records to dns.RR
	for _, r := range zone.Records {
		rr, err := s.convertToRR(r)
		if err != nil {
			continue
		}
		records = append(records, rr)
	}

	// Add DNSKEY records
	records = append(records, ksk, zsk)

	// Group records by Name and Type
	rrSets := make(map[string][]dns.RR)
	for _, r := range records {
		key := fmt.Sprintf("%s:%d", r.Header().Name, r.Header().Rrtype)
		rrSets[key] = append(rrSets[key], r)
	}

	// Sign each RRSet
	var signedRecords []dns.RR
	for _, set := range rrSets {
		signedRecords = append(signedRecords, set...)

		// Sign with ZSK
		rrsig, err := s.signRRSet(set, zsk, zskPriv)
		if err != nil {
			return nil, fmt.Errorf("failed to sign RRSet: %w", err)
		}
		signedRecords = append(signedRecords, rrsig)

		// If DNSKEY, also sign with KSK
		if set[0].Header().Rrtype == dns.TypeDNSKEY {
			rrsigKSK, err := s.signRRSet(set, ksk, kskPriv)
			if err != nil {
				return nil, fmt.Errorf("failed to sign DNSKEY with KSK: %w", err)
			}
			signedRecords = append(signedRecords, rrsigKSK)
		}
	}

	// Generate NSEC records
	nsecRecords := s.generateNSEC(records, zone.Domain)
	signedRecords = append(signedRecords, nsecRecords...)

	// Sign NSEC records
	nsecSets := make(map[string][]dns.RR)
	for _, r := range nsecRecords {
		key := fmt.Sprintf("%s:%d", r.Header().Name, r.Header().Rrtype)
		nsecSets[key] = append(nsecSets[key], r)
	}

	for _, set := range nsecSets {
		rrsig, err := s.signRRSet(set, zsk, zskPriv)
		if err != nil {
			return nil, fmt.Errorf("failed to sign NSEC RRSet: %w", err)
		}
		signedRecords = append(signedRecords, rrsig)
	}

	return signedRecords, nil
}

func (s *DNSSECService) signRRSet(rrSet []dns.RR, key *dns.DNSKEY, privKey *ecdsa.PrivateKey) (*dns.RRSIG, error) {
	rrsig := &dns.RRSIG{
		Hdr:         dns.RR_Header{Name: rrSet[0].Header().Name, Rrtype: dns.TypeRRSIG, Class: dns.ClassINET, Ttl: rrSet[0].Header().Ttl},
		TypeCovered: rrSet[0].Header().Rrtype,
		Algorithm:   key.Algorithm,
		Labels:      uint8(dns.CountLabel(rrSet[0].Header().Name)),
		OrigTtl:     rrSet[0].Header().Ttl,
		Expiration:  uint32(time.Now().Add(30 * 24 * time.Hour).Unix()), // 30 days
		Inception:   uint32(time.Now().Add(-1 * time.Hour).Unix()),      // 1 hour ago
		KeyTag:      key.KeyTag(),
		SignerName:  key.Header().Name,
	}

	if err := rrsig.Sign(privKey, rrSet); err != nil {
		return nil, err
	}

	return rrsig, nil
}

func (s *DNSSECService) generateNSEC(records []dns.RR, zoneDomain string) []dns.RR {
	// Extract all owner names
	ownerMap := make(map[string]bool)
	for _, r := range records {
		ownerMap[r.Header().Name] = true
	}

	var owners []string
	for name := range ownerMap {
		owners = append(owners, name)
	}
	sort.Strings(owners)

	var nsecs []dns.RR
	for i, owner := range owners {
		nextOwner := owners[(i+1)%len(owners)]

		// Find types for this owner
		typeMap := make(map[uint16]bool)
		typeMap[dns.TypeRRSIG] = true
		typeMap[dns.TypeNSEC] = true

		for _, r := range records {
			if r.Header().Name == owner {
				typeMap[r.Header().Rrtype] = true
			}
		}

		var types []uint16
		for t := range typeMap {
			types = append(types, t)
		}
		sort.Slice(types, func(i, j int) bool {
			return types[i] < types[j]
		})

		nsec := &dns.NSEC{
			Hdr:        dns.RR_Header{Name: owner, Rrtype: dns.TypeNSEC, Class: dns.ClassINET, Ttl: 3600},
			NextDomain: nextOwner,
			TypeBitMap: types,
		}
		nsecs = append(nsecs, nsec)
	}

	return nsecs
}

func (s *DNSSECService) convertToRR(r models.DNSRecord) (dns.RR, error) {
	// Simplified conversion - needs to be expanded for all types
	rrStr := fmt.Sprintf("%s %d IN %s %s", r.Name, r.TTL, r.Type, r.Value)
	return dns.NewRR(rrStr)
}

// ConvertRRToRecord converts a dns.RR to a models.DNSRecord
func (s *DNSSECService) ConvertRRToRecord(rr dns.RR) models.DNSRecord {
	header := rr.Header()
	record := models.DNSRecord{
		Name: header.Name,
		Type: dns.TypeToString[header.Rrtype],
		TTL:  header.Ttl,
	}

	// Get the string representation
	// dns.RR.String() returns "name ttl class type rdata" (tab separated)
	fullStr := rr.String()

	// Extract RDATA (everything after the type)
	// We assume standard formatting from miekg/dns
	parts := strings.Split(fullStr, "\t")
	if len(parts) >= 5 {
		// The last part is usually the RDATA
		// Find the index of the 4th tab
		tabCount := 0
		idx := -1
		for i, r := range fullStr {
			if r == '\t' {
				tabCount++
				if tabCount == 4 {
					idx = i
					break
				}
			}
		}

		if idx != -1 && idx+1 < len(fullStr) {
			record.Value = fullStr[idx+1:]
		}
	}

	return record
}
