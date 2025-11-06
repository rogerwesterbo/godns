# CLI Record Management Guide

This guide shows how to use `godnscli` to manage DNS records with type-specific properties.

## Authentication

First, authenticate with the GoDNS API:

```bash
godnscli login
```

This will:

1. Start OAuth2 device flow with Keycloak
2. Display a URL and code to enter in your browser
3. Save the access token to `~/.godns/token.json`

## Zone Management

### List all zones

```bash
godnscli zone list
```

### Get zone details

```bash
godnscli zone get example.lan
```

### Delete a zone

```bash
godnscli zone delete example.lan
```

## Record Management

### List records in a zone

```bash
# List all records
godnscli record list example.lan

# Filter by type
godnscli record list example.lan --type-filter A
godnscli record list example.lan --type-filter MX
```

### Get specific record

```bash
godnscli record get example.lan www.example.lan. A
godnscli record get example.lan example.lan. MX
```

### Delete a record

```bash
godnscli record delete example.lan www.example.lan. A
```

## Creating Records

### Simple Record Types

These types use the `--value` flag:

#### A Record

```bash
godnscli record create example.lan \
  --name www.example.lan. \
  --type A \
  --value 192.168.1.100 \
  --ttl 300
```

#### AAAA Record

```bash
godnscli record create example.lan \
  --name www.example.lan. \
  --type AAAA \
  --value 2001:db8::1 \
  --ttl 300
```

#### CNAME Record

```bash
godnscli record create example.lan \
  --name alias.example.lan. \
  --type CNAME \
  --value www.example.lan. \
  --ttl 300
```

#### ALIAS Record

```bash
godnscli record create example.lan \
  --name example.lan. \
  --type ALIAS \
  --value www.example.lan. \
  --ttl 300
```

#### NS Record

```bash
godnscli record create example.lan \
  --name example.lan. \
  --type NS \
  --value ns1.example.lan. \
  --ttl 3600
```

#### TXT Record

```bash
godnscli record create example.lan \
  --name example.lan. \
  --type TXT \
  --value "v=spf1 mx ~all" \
  --ttl 300
```

#### PTR Record

```bash
godnscli record create 1.168.192.in-addr.arpa \
  --name 100.1.168.192.in-addr.arpa. \
  --type PTR \
  --value www.example.lan. \
  --ttl 300
```

### Complex Record Types

These types use type-specific flags:

#### MX Record

```bash
godnscli record create example.lan \
  --name example.lan. \
  --type MX \
  --mx-priority 10 \
  --mx-host mail.example.lan. \
  --ttl 300
```

**Flags:**

- `--mx-priority` (int, 0-65535): Priority (lower = preferred)
- `--mx-host` (string, required): Mail server FQDN

#### SRV Record

```bash
godnscli record create example.lan \
  --name _http._tcp.example.lan. \
  --type SRV \
  --srv-priority 10 \
  --srv-weight 60 \
  --srv-port 80 \
  --srv-target web.example.lan. \
  --ttl 300
```

**Flags:**

- `--srv-priority` (int, 0-65535): Priority (lower = preferred)
- `--srv-weight` (int, 0-65535): Weight for load balancing
- `--srv-port` (int, 0-65535): Service port number
- `--srv-target` (string, required): Target server FQDN

#### SOA Record

```bash
godnscli record create example.lan \
  --name example.lan. \
  --type SOA \
  --soa-mname ns1.example.lan. \
  --soa-rname admin.example.lan. \
  --soa-serial 2024010101 \
  --soa-refresh 3600 \
  --soa-retry 1800 \
  --soa-expire 604800 \
  --soa-minimum 300 \
  --ttl 3600
```

**Flags:**

- `--soa-mname` (string, required): Primary nameserver FQDN
- `--soa-rname` (string, required): Admin email (@ → .)
- `--soa-serial` (uint32): Serial number (auto-generated if omitted)
- `--soa-refresh` (uint32): Refresh interval (default: 3600)
- `--soa-retry` (uint32): Retry interval (default: 1800)
- `--soa-expire` (uint32): Expire time (default: 604800)
- `--soa-minimum` (uint32): Minimum TTL (default: 300)

**Note:** If `--soa-serial` is omitted or 0, it will be auto-generated in `YYYYMMDDnn` format (e.g., `2024010101`).

#### CAA Record

```bash
godnscli record create example.lan \
  --name example.lan. \
  --type CAA \
  --caa-flags 0 \
  --caa-tag issue \
  --caa-value letsencrypt.org \
  --ttl 300
```

**Flags:**

- `--caa-flags` (int): Flags (0 = non-critical, 128 = critical)
- `--caa-tag` (string, required): Tag (`issue`, `issuewild`, `iodef`)
- `--caa-value` (string, required): CA domain or contact URL

**Common CAA Tags:**

- `issue`: Authorize CA to issue certificates
- `issuewild`: Authorize CA to issue wildcard certificates
- `iodef`: URL for reporting policy violations

## Updating Records

Use the same flags as create, but provide domain, name, and type as arguments:

### Update simple record

```bash
godnscli record update example.lan www.example.lan. A \
  --value 192.168.1.200 \
  --ttl 600
```

### Update MX record

```bash
godnscli record update example.lan example.lan. MX \
  --mx-priority 20 \
  --mx-host mail2.example.lan. \
  --ttl 300
```

### Update SRV record

```bash
godnscli record update example.lan _http._tcp.example.lan. SRV \
  --srv-priority 5 \
  --srv-weight 80 \
  --srv-port 8080 \
  --srv-target newweb.example.lan. \
  --ttl 300
```

## Configuration

### API URL Configuration

By default, `godnscli` uses `http://localhost:14000` from `~/.godns/config.yaml`:

```yaml
api:
  url: http://localhost:14000
```

Override with flag:

```bash
godnscli --api-url https://godns.example.com zone list
```

### Token Storage

Access tokens are stored in `~/.godns/token.json`:

```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "expires_at": "2024-01-15T12:00:00Z",
  "token_type": "Bearer"
}
```

Tokens are automatically refreshed when expired.

## Output Formatting

### Record List Output

The `record list` command shows records in a formatted table:

```
NAME                          TYPE    VALUE                                             TTL
example.lan.                  SOA     ns1.example.lan. admin.example.lan. (Serial: ...  3600
example.lan.                  NS      ns1.example.lan.                                  3600
example.lan.                  A       192.168.1.1                                       300
example.lan.                  MX      10 → mail.example.lan.                           300
www.example.lan.              A       192.168.1.100                                     300
_http._tcp.example.lan.       SRV     Pri:10 Wgt:60 Port:80 → web.example.lan.         300
example.lan.                  CAA     [0] issue: letsencrypt.org                        300
```

### Record Get Output

The `record get` command shows full JSON details:

```json
{
  "name": "example.lan.",
  "type": "MX",
  "ttl": 300,
  "mx_priority": 10,
  "mx_host": "mail.example.lan."
}
```

## Examples by Use Case

### Setting up a new zone

```bash
# 1. Create SOA record (required)
godnscli record create example.lan \
  --name example.lan. \
  --type SOA \
  --soa-mname ns1.example.lan. \
  --soa-rname admin.example.lan. \
  --ttl 3600

# 2. Add NS records
godnscli record create example.lan \
  --name example.lan. \
  --type NS \
  --value ns1.example.lan. \
  --ttl 3600

# 3. Add A record for apex
godnscli record create example.lan \
  --name example.lan. \
  --type A \
  --value 192.168.1.1 \
  --ttl 300

# 4. Add www subdomain
godnscli record create example.lan \
  --name www.example.lan. \
  --type A \
  --value 192.168.1.100 \
  --ttl 300

# 5. Add mail server
godnscli record create example.lan \
  --name example.lan. \
  --type MX \
  --mx-priority 10 \
  --mx-host mail.example.lan. \
  --ttl 300
```

### Setting up mail infrastructure

```bash
# MX record for mail routing
godnscli record create example.lan \
  --name example.lan. \
  --type MX \
  --mx-priority 10 \
  --mx-host mail.example.lan. \
  --ttl 300

# A record for mail server
godnscli record create example.lan \
  --name mail.example.lan. \
  --type A \
  --value 192.168.1.50 \
  --ttl 300

# SPF record
godnscli record create example.lan \
  --name example.lan. \
  --type TXT \
  --value "v=spf1 mx ~all" \
  --ttl 300

# DKIM record
godnscli record create example.lan \
  --name default._domainkey.example.lan. \
  --type TXT \
  --value "v=DKIM1; k=rsa; p=MIGfMA0GCS..." \
  --ttl 300

# DMARC record
godnscli record create example.lan \
  --name _dmarc.example.lan. \
  --type TXT \
  --value "v=DMARC1; p=quarantine; rua=mailto:dmarc@example.lan" \
  --ttl 300
```

### Setting up service discovery

```bash
# HTTP service
godnscli record create example.lan \
  --name _http._tcp.example.lan. \
  --type SRV \
  --srv-priority 10 \
  --srv-weight 60 \
  --srv-port 80 \
  --srv-target web1.example.lan. \
  --ttl 300

# HTTPS service
godnscli record create example.lan \
  --name _https._tcp.example.lan. \
  --type SRV \
  --srv-priority 10 \
  --srv-weight 60 \
  --srv-port 443 \
  --srv-target web1.example.lan. \
  --ttl 300

# Add target A records
godnscli record create example.lan \
  --name web1.example.lan. \
  --type A \
  --value 192.168.1.100 \
  --ttl 300
```

### Setting up SSL/TLS with CAA

```bash
# Allow Let's Encrypt to issue certificates
godnscli record create example.lan \
  --name example.lan. \
  --type CAA \
  --caa-flags 0 \
  --caa-tag issue \
  --caa-value letsencrypt.org \
  --ttl 300

# Allow Let's Encrypt to issue wildcard certificates
godnscli record create example.lan \
  --name example.lan. \
  --type CAA \
  --caa-flags 0 \
  --caa-tag issuewild \
  --caa-value letsencrypt.org \
  --ttl 300

# Add iodef for reporting violations
godnscli record create example.lan \
  --name example.lan. \
  --type CAA \
  --caa-flags 0 \
  --caa-tag iodef \
  --caa-value "mailto:security@example.lan" \
  --ttl 300
```

## Error Handling

### Not authenticated

```
Error: not authenticated. Please run 'godnscli login' first
```

**Solution:** Run `godnscli login`

### Missing required flags

```
Error: --mx-host is required for MX records
```

**Solution:** Add the required type-specific flag

### API connection failed

```
Error: failed to connect to API: dial tcp [::1]:14000: connect: connection refused
```

**Solution:** Check that the GoDNS API server is running, or update `~/.godns/config.yaml` with correct API URL

### Record not found

```
Error: API request failed (404): Record not found
```

**Solution:** Verify the domain, name, and type are correct using `godnscli record list`

## Tips

1. **Always use FQDNs**: End domain names with a dot (`.`) to indicate they are fully qualified

   - ✅ `www.example.lan.`
   - ❌ `www.example.lan`

2. **Check existing records**: Use `godnscli record list` before creating to avoid duplicates

3. **SOA serial numbers**: Let the CLI auto-generate serial numbers for SOA records (omit `--soa-serial`)

4. **MX priorities**: Lower numbers have higher priority (10 is preferred over 20)

5. **SRV naming**: Follow the format `_service._protocol.domain.` (e.g., `_http._tcp.example.lan.`)

6. **CAA for security**: Add CAA records to restrict which CAs can issue certificates for your domain

7. **Type filtering**: Use `--type-filter` to quickly find specific record types in large zones

## See Also

- [DNS Record Types Documentation](./DNS_RECORD_TYPES.md) - Detailed record type specifications
- [API Documentation](./API_DOCUMENTATION.md) - HTTP API reference
- [Quick Start Guide](./QUICK_START.md) - Getting started with GoDNS
