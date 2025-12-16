# Testing DNSSEC

This document describes how to test the DNSSEC functionality in GoDNS.

## Unit Testing

A unit test has been added to `internal/services/v1dnssecservice/service_test.go` to verify the key generation and zone signing logic.

To run the test:

```bash
go test -v internal/services/v1dnssecservice/service_test.go internal/services/v1dnssecservice/service.go
```

This test performs the following steps:

1. Generates KSK and ZSK keys for `example.com.`.
2. Creates a test zone with A records.
3. Signs the zone using the generated keys.
4. Verifies that the output contains `DNSKEY`, `RRSIG`, and `NSEC` records.

## Integration Testing (Manual)

To test DNSSEC in the running server, you need to integrate the signing process into your data seeding or API.

### 1. Update Seeding Service

You can modify `internal/services/seeding/seeding_service.go` to sign a test zone.

Add the `DNSSECService` to the `SeedingService` struct and constructor.

Then, in `seedTestData`, you can do something like this:

```go
// Generate keys
ksk, zsk, kskPriv, zskPriv, err := s.dnssecService.GenerateKeys("home.lan")
if err != nil {
    return err
}

// Create the zone object
zone := &models.DNSZone{
    Domain: "home.lan",
    Records: []models.DNSRecord{...}, // Add your records here
}

// Sign the zone
signedRRs, err := s.dnssecService.SignZone(zone, ksk, zsk, kskPriv, zskPriv)
if err != nil {
    return err
}

// Convert signed RRs back to models.DNSRecord and add to zone
var signedRecords []models.DNSRecord
for _, rr := range signedRRs {
    signedRecords = append(signedRecords, s.dnssecService.ConvertRRToRecord(rr))
}
zone.Records = signedRecords
zone.DNSSECEnabled = true
// Store keys (PEM format) if needed for future re-signing
// zone.KSK = ...

// Save the zone
if err := s.zoneService.CreateZone(ctx, zone); err != nil {
    return err
}
```

### 2. Verify with `dig`

Once the server is running with a signed zone, you can verify it using `dig`:

```bash
# Request DNSKEY records
dig @localhost -p 5354 home.lan DNSKEY +dnssec

# Request A record with signatures
dig @localhost -p 5354 ns1.home.lan A +dnssec
```

You should see `RRSIG` records in the answer section.
