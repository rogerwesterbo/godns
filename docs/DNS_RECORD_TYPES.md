# DNS Record Types Reference

GoDNS supports multiple DNS record types with both legacy (simple `value` field) and structured (type-specific fields) formats.

## Supported Record Types

| Type  | Purpose            | Structured Format | Legacy Format |
| ----- | ------------------ | ----------------- | ------------- |
| A     | IPv4 address       | ✓                 | ✓             |
| AAAA  | IPv6 address       | ✓                 | ✓             |
| CNAME | Canonical name     | ✓                 | ✓             |
| ALIAS | Zone apex alias    | ✓                 | ✓             |
| NS    | Name server        | ✓                 | ✓             |
| MX    | Mail exchange      | ✓                 | ✓             |
| TXT   | Text record        | ✓                 | ✓             |
| SRV   | Service discovery  | ✓                 | ✓             |
| SOA   | Start of authority | ✓                 | ✓             |
| CAA   | CA authorization   | ✓                 | ✓             |
| PTR   | Pointer (reverse)  | ✓                 | ✓             |

## Record Format Examples

### A Record (IPv4 Address)

**Structured Format:**

```json
{
  "name": "www.example.com.",
  "type": "A",
  "ttl": 300,
  "value": "192.168.1.10"
}
```

**DNS Output:**

```
www.example.com. 300 IN A 192.168.1.10
```

---

### AAAA Record (IPv6 Address)

**Structured Format:**

```json
{
  "name": "www.example.com.",
  "type": "AAAA",
  "ttl": 300,
  "value": "2001:db8::1"
}
```

**DNS Output:**

```
www.example.com. 300 IN AAAA 2001:db8::1
```

---

### CNAME Record (Canonical Name)

**Structured Format:**

```json
{
  "name": "www.example.com.",
  "type": "CNAME",
  "ttl": 300,
  "value": "web.example.com."
}
```

**DNS Output:**

```
www.example.com. 300 IN CNAME web.example.com.
```

**Note:** CNAME cannot be used at the zone apex (use ALIAS instead).

---

### ALIAS Record (Zone Apex Alias)

**Structured Format:**

```json
{
  "name": "example.com.",
  "type": "ALIAS",
  "ttl": 300,
  "value": "www.example.com."
}
```

**Use Case:** Unlike CNAME, ALIAS can be used at the zone apex (`example.com` vs `www.example.com`).

---

### NS Record (Name Server)

**Structured Format:**

```json
{
  "name": "example.com.",
  "type": "NS",
  "ttl": 3600,
  "value": "ns1.example.com."
}
```

**DNS Output:**

```
example.com. 3600 IN NS ns1.example.com.
```

---

### MX Record (Mail Exchange)

**Structured Format (Recommended):**

```json
{
  "name": "example.com.",
  "type": "MX",
  "ttl": 300,
  "mx_priority": 10,
  "mx_host": "mail.example.com."
}
```

**Legacy Format (Still Supported):**

```json
{
  "name": "example.com.",
  "type": "MX",
  "ttl": 300,
  "value": "10 mail.example.com."
}
```

**DNS Output:**

```
example.com. 300 IN MX 10 mail.example.com.
```

**Fields:**

- `mx_priority` (0-65535): Lower values have higher priority
- `mx_host`: Fully qualified domain name of mail server

---

### TXT Record (Text)

**Structured Format:**

```json
{
  "name": "example.com.",
  "type": "TXT",
  "ttl": 300,
  "value": "v=spf1 mx ip4:192.168.1.0/24 -all"
}
```

**Common Uses:**

- SPF records (`v=spf1 ...`)
- DKIM records (`v=DKIM1 ...`)
- DMARC records (`v=DMARC1 ...`)
- Domain verification

---

### SRV Record (Service Discovery)

**Structured Format (Recommended):**

```json
{
  "name": "_http._tcp.example.com.",
  "type": "SRV",
  "ttl": 300,
  "srv_priority": 10,
  "srv_weight": 60,
  "srv_port": 80,
  "srv_target": "web.example.com."
}
```

**Legacy Format (Still Supported):**

```json
{
  "name": "_http._tcp.example.com.",
  "type": "SRV",
  "ttl": 300,
  "value": "10 60 80 web.example.com."
}
```

**DNS Output:**

```
_http._tcp.example.com. 300 IN SRV 10 60 80 web.example.com.
```

**Fields:**

- `srv_priority` (0-65535): Priority of target host (lower = higher priority)
- `srv_weight` (0-65535): Relative weight for load balancing
- `srv_port` (0-65535): TCP/UDP port number
- `srv_target`: Target hostname

**Common Service Names:**

- `_http._tcp` - HTTP service
- `_https._tcp` - HTTPS service
- `_ldap._tcp` - LDAP service
- `_etcd-server._tcp` - etcd service (Kubernetes)

---

### SOA Record (Start of Authority)

**Structured Format (Recommended):**

```json
{
  "name": "example.com.",
  "type": "SOA",
  "ttl": 3600,
  "soa_mname": "ns1.example.com.",
  "soa_rname": "hostmaster.example.com.",
  "soa_serial": 2024110601,
  "soa_refresh": 3600,
  "soa_retry": 1800,
  "soa_expire": 604800,
  "soa_minimum": 300
}
```

**Legacy Format (Still Supported):**

```json
{
  "name": "example.com.",
  "type": "SOA",
  "ttl": 3600,
  "value": "ns1.example.com. hostmaster.example.com. 2024110601 3600 1800 604800 300"
}
```

**DNS Output:**

```
example.com. 3600 IN SOA ns1.example.com. hostmaster.example.com. 2024110601 3600 1800 604800 300
```

**Fields:**

- `soa_mname`: Primary name server for this zone
- `soa_rname`: Email address of zone administrator (dots replaced with @)
- `soa_serial`: Zone serial number (format: YYYYMMDDnn)
- `soa_refresh`: Seconds before secondary checks for updates
- `soa_retry`: Seconds before secondary retries after failure
- `soa_expire`: Seconds before secondary stops answering
- `soa_minimum`: Minimum TTL for negative caching

---

### CAA Record (Certificate Authority Authorization)

**Structured Format (Recommended):**

```json
{
  "name": "example.com.",
  "type": "CAA",
  "ttl": 300,
  "caa_flags": 0,
  "caa_tag": "issue",
  "caa_value": "letsencrypt.org"
}
```

**Legacy Format (Still Supported):**

```json
{
  "name": "example.com.",
  "type": "CAA",
  "ttl": 300,
  "value": "0 issue \"letsencrypt.org\""
}
```

**DNS Output:**

```
example.com. 300 IN CAA 0 issue "letsencrypt.org"
```

**Fields:**

- `caa_flags` (0 or 128): 0 = non-critical, 128 = critical
- `caa_tag`: Property tag (`issue`, `issuewild`, `iodef`)
- `caa_value`: Property value (CA domain or email)

**Common Tags:**

- `issue`: Authorize CA to issue certificates
- `issuewild`: Authorize CA to issue wildcard certificates
- `iodef`: URL for reporting policy violations

---

### PTR Record (Pointer/Reverse DNS)

**Structured Format:**

```json
{
  "name": "10.1.168.192.in-addr.arpa.",
  "type": "PTR",
  "ttl": 300,
  "value": "web.example.com."
}
```

**DNS Output:**

```
10.1.168.192.in-addr.arpa. 300 IN PTR web.example.com.
```

**Use Case:** Reverse DNS lookup (IP to hostname).

---

## Migration Guide

### Upgrading from Legacy to Structured Format

The API accepts both formats for backward compatibility. To migrate:

**Old MX Record:**

```json
{ "type": "MX", "value": "10 mail.example.com." }
```

**New MX Record:**

```json
{ "type": "MX", "mx_priority": 10, "mx_host": "mail.example.com." }
```

### Benefits of Structured Format

1. **Type Safety**: Fields are validated by type
2. **Clarity**: Explicit field names (no parsing required)
3. **API Consistency**: Same pattern across all complex types
4. **Tooling**: Better IDE autocomplete and validation

### Backward Compatibility

- ✅ Existing records with `value` field continue to work
- ✅ API accepts both formats on input
- ✅ No migration required for existing data
- ⚠️ New structured format is recommended for new records

---

## Common Patterns

### Complete Zone Example

```json
{
  "domain": "example.com",
  "records": [
    {
      "name": "example.com.",
      "type": "SOA",
      "ttl": 3600,
      "soa_mname": "ns1.example.com.",
      "soa_rname": "hostmaster.example.com.",
      "soa_serial": 2024110601,
      "soa_refresh": 3600,
      "soa_retry": 1800,
      "soa_expire": 604800,
      "soa_minimum": 300
    },
    {
      "name": "example.com.",
      "type": "NS",
      "ttl": 3600,
      "value": "ns1.example.com."
    },
    {
      "name": "example.com.",
      "type": "NS",
      "ttl": 3600,
      "value": "ns2.example.com."
    },
    {
      "name": "ns1.example.com.",
      "type": "A",
      "ttl": 3600,
      "value": "192.168.1.1"
    },
    {
      "name": "ns2.example.com.",
      "type": "A",
      "ttl": 3600,
      "value": "192.168.1.2"
    },
    {
      "name": "example.com.",
      "type": "ALIAS",
      "ttl": 300,
      "value": "www.example.com."
    },
    {
      "name": "www.example.com.",
      "type": "A",
      "ttl": 300,
      "value": "192.168.1.10"
    },
    {
      "name": "example.com.",
      "type": "MX",
      "ttl": 300,
      "mx_priority": 10,
      "mx_host": "mail.example.com."
    },
    {
      "name": "mail.example.com.",
      "type": "A",
      "ttl": 300,
      "value": "192.168.1.20"
    },
    {
      "name": "example.com.",
      "type": "TXT",
      "ttl": 300,
      "value": "v=spf1 mx ip4:192.168.1.0/24 -all"
    },
    {
      "name": "_http._tcp.example.com.",
      "type": "SRV",
      "ttl": 300,
      "srv_priority": 10,
      "srv_weight": 60,
      "srv_port": 80,
      "srv_target": "www.example.com."
    },
    {
      "name": "example.com.",
      "type": "CAA",
      "ttl": 300,
      "caa_flags": 0,
      "caa_tag": "issue",
      "caa_value": "letsencrypt.org"
    }
  ]
}
```

---

## API Endpoints

See [API Documentation](./API_DOCUMENTATION.md) for complete API reference.

**Create Zone with Records:**

```bash
POST /api/v1/zones
```

**Add Record to Zone:**

```bash
POST /api/v1/zones/{domain}/records
```

**Update Record:**

```bash
PUT /api/v1/zones/{domain}/records/{name}/{type}
```

---

## See Also

- [API Documentation](./API_DOCUMENTATION.md)
- [CLI Guide](./CLI_GUIDE.md)
- [Quick Start](./QUICK_START.md)
- [Testing Guide](./TESTING_GUIDE.md)
