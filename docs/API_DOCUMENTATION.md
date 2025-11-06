# GoDNS HTTP API Documentation

The GoDNS HTTP API provides RESTful endpoints for managing DNS zones and records.

## Table of Contents

- [Interactive Documentation (Swagger UI)](#interactive-documentation-swagger-ui)
- [Base URL](#base-url)
- [Quick Start](#quick-start)
- [Health Endpoints](#health-endpoints)
- [DNS Zone Endpoints](#dns-zone-endpoints)
- [DNS Record Endpoints](#dns-record-endpoints)
- [Data Models](#data-models)
- [Example Usage](#example-usage)
- [Error Responses](#error-responses)
- [Configuration](#configuration)
- [Swagger/OpenAPI Integration](#swaggeropenapi-integration)
- [Development Workflow](#development-workflow)

---

## Interactive Documentation (Swagger UI)

GoDNS includes comprehensive Swagger/OpenAPI documentation with an interactive UI for exploring and testing the API.

### Accessing Swagger UI

**URL:** http://localhost:14000/swagger/index.html

**Features:**

- ✅ Interactive API exploration - Browse all endpoints with descriptions
- ✅ Try It Out - Test API endpoints directly from the browser
- ✅ Schema Validation - See request/response data structures
- ✅ Example Values - Pre-filled examples for all models
- ✅ Auto-Generated - Always in sync with the code
- ✅ Standard Format - OpenAPI 3.0 specification

### Quick Start with Swagger

1. **Start the API Server:**

   ```bash
   # The API is integrated with the DNS server
   make build-dns
   ./bin/godns

   # Enable HTTP API in .env if not already enabled
   # DNS_ENABLE_HTTP_API=true (default)
   ```

2. **Open Swagger UI:**
   Visit http://localhost:14000/swagger/index.html

3. **Explore & Test:**
   - Browse endpoints organized by tags (Health, Zones, Records)
   - Click any endpoint to see details
   - Click "Try it out" to test
   - Fill in parameters and click "Execute"
   - View the response

### Generated Files

Swagger documentation generates these files in `docs/`:

- `docs.go` - Embedded Go code for Swagger spec
- `swagger.json` - OpenAPI specification (JSON format)
- `swagger.yaml` - OpenAPI specification (YAML format)

## Base URL

By default, the API server runs on port `14000`. You can configure this using the `HTTP_API_PORT` environment variable.

```
http://localhost:14000
```

## Health Endpoints

### Health Check

Check if the API server is healthy.

**Endpoint:** `GET /health`

**Response:**

```json
{
  "status": "healthy"
}
```

### Readiness Check

Check if the API server is ready to accept requests.

**Endpoint:** `GET /ready`

**Response:**

```json
{
  "status": "ready"
}
```

---

## DNS Zone Endpoints

### List All Zones

Get a list of all DNS zones.

**Endpoint:** `GET /api/v1/zones`

**Response:**

```json
[
  {
    "domain": "example.lan.",
    "records": [
      {
        "name": "www.example.lan.",
        "type": "A",
        "ttl": 300,
        "value": "192.168.1.100"
      }
    ]
  }
]
```

### Create Zone

Create a new DNS zone with optional records.

**Endpoint:** `POST /api/v1/zones`

**Request Body:**

```json
{
  "domain": "example.lan",
  "records": [
    {
      "name": "www.example.lan.",
      "type": "A",
      "ttl": 300,
      "value": "192.168.1.100"
    },
    {
      "name": "mail.example.lan.",
      "type": "A",
      "ttl": 300,
      "value": "192.168.1.101"
    }
  ]
}
```

**Response:** `201 Created`

```json
{
  "domain": "example.lan.",
  "records": [...]
}
```

**Errors:**

- `400 Bad Request` - Invalid zone data
- `409 Conflict` - Zone already exists

### Get Zone

Retrieve a specific DNS zone by domain.

**Endpoint:** `GET /api/v1/zones/{domain}`

**Example:** `GET /api/v1/zones/example.lan`

**Response:**

```json
{
  "domain": "example.lan.",
  "records": [
    {
      "name": "www.example.lan.",
      "type": "A",
      "ttl": 300,
      "value": "192.168.1.100"
    }
  ]
}
```

**Errors:**

- `404 Not Found` - Zone does not exist

### Update Zone

Update an existing DNS zone. This replaces all records in the zone.

**Endpoint:** `PUT /api/v1/zones/{domain}`

**Request Body:**

```json
{
  "domain": "example.lan",
  "records": [
    {
      "name": "www.example.lan.",
      "type": "A",
      "ttl": 600,
      "value": "192.168.1.200"
    }
  ]
}
```

**Response:** `200 OK`

**Errors:**

- `400 Bad Request` - Invalid zone data
- `404 Not Found` - Zone does not exist

### Delete Zone

Delete a DNS zone and all its records.

**Endpoint:** `DELETE /api/v1/zones/{domain}`

**Response:** `204 No Content`

**Errors:**

- `404 Not Found` - Zone does not exist

---

## DNS Record Endpoints

### Create Record

Add a new record to an existing zone.

**Endpoint:** `POST /api/v1/zones/{domain}/records`

**Request Body:**

```json
{
  "name": "api.example.lan.",
  "type": "A",
  "ttl": 300,
  "value": "192.168.1.150"
}
```

**Response:** `201 Created`

**Errors:**

- `400 Bad Request` - Invalid record data
- `404 Not Found` - Zone does not exist
- `409 Conflict` - Record already exists

### Get Record

Retrieve a specific DNS record.

**Endpoint:** `GET /api/v1/zones/{domain}/records/{name}/{type}`

**Example:** `GET /api/v1/zones/example.lan/records/www.example.lan./A`

**Response:**

```json
{
  "name": "www.example.lan.",
  "type": "A",
  "ttl": 300,
  "value": "192.168.1.100"
}
```

**Errors:**

- `404 Not Found` - Zone or record does not exist

### Update Record

Update an existing DNS record.

**Endpoint:** `PUT /api/v1/zones/{domain}/records/{name}/{type}`

**Request Body:**

```json
{
  "name": "www.example.lan.",
  "type": "A",
  "ttl": 600,
  "value": "192.168.1.200"
}
```

**Response:** `200 OK`

**Errors:**

- `400 Bad Request` - Invalid record data
- `404 Not Found` - Zone or record does not exist

### Delete Record

Delete a specific DNS record from a zone.

**Endpoint:** `DELETE /api/v1/zones/{domain}/records/{name}/{type}`

**Response:** `204 No Content`

**Errors:**

- `404 Not Found` - Zone or record does not exist

---

## Data Models

### DNSZone

```json
{
  "domain": "string (required)",  // Domain name, will be normalized with trailing dot
  "records": [DNSRecord]          // Array of DNS records
}
```

### DNSRecord

The DNSRecord model supports both simple records (A, AAAA, CNAME, etc.) and complex records with type-specific fields (MX, SRV, SOA, CAA).

**Basic Record Structure:**

```json
{
  "name": "string (required)",    // Fully qualified domain name
  "type": "string (required)",    // Record type: A, AAAA, CNAME, ALIAS, MX, NS, TXT, PTR, SRV, SOA, CAA
  "ttl": number,                  // Time to live in seconds (default: 300)
  "value": "string"               // Record value - used for simple types (A, AAAA, CNAME, NS, TXT, PTR)
}
```

**Type-Specific Fields:**

For certain record types, you can use structured fields instead of (or in addition to) the `value` field:

**MX Records:**

```json
{
  "name": "example.lan.",
  "type": "MX",
  "ttl": 300,
  "mx_priority": 10, // Mail server priority (0-65535)
  "mx_host": "mail.example.lan." // Mail server hostname
}
```

**SRV Records:**

```json
{
  "name": "_http._tcp.example.lan.",
  "type": "SRV",
  "ttl": 300,
  "srv_priority": 10, // Priority (0-65535)
  "srv_weight": 60, // Weight for load balancing (0-65535)
  "srv_port": 80, // Service port (0-65535)
  "srv_target": "web.example.lan." // Target hostname
}
```

**SOA Records:**

```json
{
  "name": "example.lan.",
  "type": "SOA",
  "ttl": 3600,
  "soa_mname": "ns1.example.lan.", // Primary nameserver
  "soa_rname": "hostmaster.example.lan.", // Admin email (@ replaced with .)
  "soa_serial": 2024110601, // Serial number (YYYYMMDDnn)
  "soa_refresh": 3600, // Refresh interval (seconds)
  "soa_retry": 1800, // Retry interval (seconds)
  "soa_expire": 604800, // Expire time (seconds)
  "soa_minimum": 300 // Minimum TTL (seconds)
}
```

**CAA Records:**

```json
{
  "name": "example.lan.",
  "type": "CAA",
  "ttl": 300,
  "caa_flags": 0, // Flags (0 or 128 for critical)
  "caa_tag": "issue", // Tag: issue, issuewild, iodef
  "caa_value": "letsencrypt.org" // Value (CA domain or URL)
}
```

### Supported Record Types

- **A** - IPv4 address (use `value` field)
- **AAAA** - IPv6 address (use `value` field)
- **CNAME** - Canonical name alias (use `value` field)
- **ALIAS** - Zone apex alias (use `value` field)
- **MX** - Mail exchange (use `mx_priority` + `mx_host` or `value`)
- **NS** - Name server (use `value` field)
- **TXT** - Text record (use `value` field)
- **PTR** - Pointer record (use `value` field)
- **SRV** - Service record (use `srv_*` fields or `value`)
- **SOA** - Start of authority (use `soa_*` fields or `value`)
- **CAA** - Certification authority authorization (use `caa_*` fields or `value`)

**Note:** For records with type-specific fields, you can use either the structured fields OR the `value` field with space-separated values. The structured approach is recommended for clarity and validation.

---

## Example Usage

### Using curl

#### Create a zone with simple records:

```bash
curl -X POST http://localhost:14000/api/v1/zones \
  -H "Content-Type: application/json" \
  -d '{
    "domain": "example.lan",
    "records": [
      {
        "name": "www.example.lan.",
        "type": "A",
        "ttl": 300,
        "value": "192.168.1.100"
      },
      {
        "name": "mail.example.lan.",
        "type": "A",
        "ttl": 300,
        "value": "192.168.1.101"
      }
    ]
  }'
```

#### Create a zone with type-specific records:

```bash
curl -X POST http://localhost:14000/api/v1/zones \
  -H "Content-Type: application/json" \
  -d '{
    "domain": "example.lan",
    "records": [
      {
        "name": "example.lan.",
        "type": "SOA",
        "ttl": 3600,
        "soa_mname": "ns1.example.lan.",
        "soa_rname": "hostmaster.example.lan.",
        "soa_serial": 2024110601,
        "soa_refresh": 3600,
        "soa_retry": 1800,
        "soa_expire": 604800,
        "soa_minimum": 300
      },
      {
        "name": "example.lan.",
        "type": "NS",
        "ttl": 3600,
        "value": "ns1.example.lan."
      },
      {
        "name": "ns1.example.lan.",
        "type": "A",
        "ttl": 3600,
        "value": "192.168.1.1"
      },
      {
        "name": "example.lan.",
        "type": "MX",
        "ttl": 300,
        "mx_priority": 10,
        "mx_host": "mail.example.lan."
      },
      {
        "name": "mail.example.lan.",
        "type": "A",
        "ttl": 300,
        "value": "192.168.1.20"
      },
      {
        "name": "_http._tcp.example.lan.",
        "type": "SRV",
        "ttl": 300,
        "srv_priority": 10,
        "srv_weight": 60,
        "srv_port": 80,
        "srv_target": "web.example.lan."
      },
      {
        "name": "web.example.lan.",
        "type": "A",
        "ttl": 300,
        "value": "192.168.1.100"
      }
    ]
  }'
```

#### List all zones:

```bash
curl http://localhost:14000/api/v1/zones
```

#### Get a specific zone:

```bash
curl http://localhost:14000/api/v1/zones/example.lan
```

#### Add a simple A record to a zone:

```bash
curl -X POST http://localhost:14000/api/v1/zones/example.lan/records \
  -H "Content-Type: application/json" \
  -d '{
    "name": "api.example.lan.",
    "type": "A",
    "ttl": 300,
    "value": "192.168.1.150"
  }'
```

#### Add an MX record with priority:

```bash
curl -X POST http://localhost:14000/api/v1/zones/example.lan/records \
  -H "Content-Type: application/json" \
  -d '{
    "name": "example.lan.",
    "type": "MX",
    "ttl": 300,
    "mx_priority": 10,
    "mx_host": "mail.example.lan."
  }'
```

#### Add an SRV record for service discovery:

```bash
curl -X POST http://localhost:14000/api/v1/zones/example.lan/records \
  -H "Content-Type: application/json" \
  -d '{
    "name": "_ldap._tcp.example.lan.",
    "type": "SRV",
    "ttl": 300,
    "srv_priority": 10,
    "srv_weight": 100,
    "srv_port": 389,
    "srv_target": "ldap.example.lan."
  }'
```

#### Add a CAA record for certificate authority:

```bash
curl -X POST http://localhost:14000/api/v1/zones/example.lan/records \
  -H "Content-Type: application/json" \
  -d '{
    "name": "example.lan.",
    "type": "CAA",
    "ttl": 300,
    "caa_flags": 0,
    "caa_tag": "issue",
    "caa_value": "letsencrypt.org"
  }'
```

#### Update a record:

```bash
curl -X PUT http://localhost:14000/api/v1/zones/example.lan/records/www.example.lan./A \
  -H "Content-Type: application/json" \
  -d '{
    "name": "www.example.lan.",
    "type": "A",
    "ttl": 600,
    "value": "192.168.1.200"
  }'
```

#### Delete a record:

```bash
curl -X DELETE http://localhost:14000/api/v1/zones/example.lan/records/api.example.lan./A
```

#### Delete a zone:

```bash
curl -X DELETE http://localhost:14000/api/v1/zones/example.lan
```

---

## DNS Record Type Examples

This section provides complete examples for all supported DNS record types.

### Simple Record Types

These record types use only the `value` field:

#### A Record (IPv4)

```json
{
  "name": "www.example.lan.",
  "type": "A",
  "ttl": 300,
  "value": "192.168.1.100"
}
```

#### AAAA Record (IPv6)

```json
{
  "name": "www.example.lan.",
  "type": "AAAA",
  "ttl": 300,
  "value": "2001:db8::1"
}
```

#### CNAME Record (Alias)

```json
{
  "name": "www.example.lan.",
  "type": "CNAME",
  "ttl": 300,
  "value": "web.example.lan."
}
```

#### ALIAS Record (Zone Apex Alias)

```json
{
  "name": "example.lan.",
  "type": "ALIAS",
  "ttl": 300,
  "value": "lb.example.lan."
}
```

#### NS Record (Name Server)

```json
{
  "name": "example.lan.",
  "type": "NS",
  "ttl": 3600,
  "value": "ns1.example.lan."
}
```

#### TXT Record (Text)

```json
{
  "name": "example.lan.",
  "type": "TXT",
  "ttl": 300,
  "value": "v=spf1 mx -all"
}
```

#### PTR Record (Reverse DNS)

```json
{
  "name": "100.1.168.192.in-addr.arpa.",
  "type": "PTR",
  "ttl": 300,
  "value": "www.example.lan."
}
```

### Complex Record Types

These record types support type-specific fields for better structure and validation:

#### MX Record (Mail Exchange)

**Structured format (recommended):**

```json
{
  "name": "example.lan.",
  "type": "MX",
  "ttl": 300,
  "mx_priority": 10,
  "mx_host": "mail.example.lan."
}
```

**Legacy format (also supported):**

```json
{
  "name": "example.lan.",
  "type": "MX",
  "ttl": 300,
  "value": "10 mail.example.lan."
}
```

**Multiple MX records for redundancy:**

```bash
# Primary mail server
curl -X POST http://localhost:14000/api/v1/zones/example.lan/records \
  -H "Content-Type: application/json" \
  -d '{
    "name": "example.lan.",
    "type": "MX",
    "ttl": 300,
    "mx_priority": 10,
    "mx_host": "mail1.example.lan."
  }'

# Backup mail server
curl -X POST http://localhost:14000/api/v1/zones/example.lan/records \
  -H "Content-Type: application/json" \
  -d '{
    "name": "example.lan.",
    "type": "MX",
    "ttl": 300,
    "mx_priority": 20,
    "mx_host": "mail2.example.lan."
  }'
```

#### SRV Record (Service Discovery)

**Structured format (recommended):**

```json
{
  "name": "_http._tcp.example.lan.",
  "type": "SRV",
  "ttl": 300,
  "srv_priority": 10,
  "srv_weight": 60,
  "srv_port": 80,
  "srv_target": "web.example.lan."
}
```

**Legacy format (also supported):**

```json
{
  "name": "_http._tcp.example.lan.",
  "type": "SRV",
  "ttl": 300,
  "value": "10 60 80 web.example.lan."
}
```

**Common SRV record examples:**

```bash
# LDAP service
curl -X POST http://localhost:14000/api/v1/zones/example.lan/records \
  -H "Content-Type: application/json" \
  -d '{
    "name": "_ldap._tcp.example.lan.",
    "type": "SRV",
    "ttl": 300,
    "srv_priority": 10,
    "srv_weight": 100,
    "srv_port": 389,
    "srv_target": "ldap.example.lan."
  }'

# HTTPS service with load balancing
curl -X POST http://localhost:14000/api/v1/zones/example.lan/records \
  -H "Content-Type: application/json" \
  -d '{
    "name": "_https._tcp.example.lan.",
    "type": "SRV",
    "ttl": 300,
    "srv_priority": 10,
    "srv_weight": 50,
    "srv_port": 443,
    "srv_target": "web1.example.lan."
  }'

# Kubernetes etcd service
curl -X POST http://localhost:14000/api/v1/zones/k8s.local/records \
  -H "Content-Type: application/json" \
  -d '{
    "name": "_etcd-server._tcp.k8s.local.",
    "type": "SRV",
    "ttl": 300,
    "srv_priority": 10,
    "srv_weight": 100,
    "srv_port": 2380,
    "srv_target": "master.k8s.local."
  }'
```

#### SOA Record (Start of Authority)

**Structured format (recommended):**

```json
{
  "name": "example.lan.",
  "type": "SOA",
  "ttl": 3600,
  "soa_mname": "ns1.example.lan.",
  "soa_rname": "hostmaster.example.lan.",
  "soa_serial": 2024110601,
  "soa_refresh": 3600,
  "soa_retry": 1800,
  "soa_expire": 604800,
  "soa_minimum": 300
}
```

**Legacy format (also supported):**

```json
{
  "name": "example.lan.",
  "type": "SOA",
  "ttl": 3600,
  "value": "ns1.example.lan. hostmaster.example.lan. 2024110601 3600 1800 604800 300"
}
```

**Field descriptions:**

- `soa_mname`: Primary nameserver for the zone
- `soa_rname`: Email of zone administrator (@ replaced with .)
- `soa_serial`: Zone serial number (format: YYYYMMDDnn)
- `soa_refresh`: Secondary nameserver refresh interval (seconds)
- `soa_retry`: Retry interval if refresh fails (seconds)
- `soa_expire`: When zone data expires (seconds)
- `soa_minimum`: Minimum TTL for negative caching (seconds)

#### CAA Record (Certificate Authority Authorization)

**Structured format (recommended):**

```json
{
  "name": "example.lan.",
  "type": "CAA",
  "ttl": 300,
  "caa_flags": 0,
  "caa_tag": "issue",
  "caa_value": "letsencrypt.org"
}
```

**Legacy format (also supported):**

```json
{
  "name": "example.lan.",
  "type": "CAA",
  "ttl": 300,
  "value": "0 issue \"letsencrypt.org\""
}
```

**Common CAA examples:**

```bash
# Allow Let's Encrypt to issue certificates
curl -X POST http://localhost:14000/api/v1/zones/example.lan/records \
  -H "Content-Type: application/json" \
  -d '{
    "name": "example.lan.",
    "type": "CAA",
    "ttl": 300,
    "caa_flags": 0,
    "caa_tag": "issue",
    "caa_value": "letsencrypt.org"
  }'

# Allow wildcard certificates
curl -X POST http://localhost:14000/api/v1/zones/example.lan/records \
  -H "Content-Type: application/json" \
  -d '{
    "name": "example.lan.",
    "type": "CAA",
    "ttl": 300,
    "caa_flags": 0,
    "caa_tag": "issuewild",
    "caa_value": "letsencrypt.org"
  }'

# Incident reporting
curl -X POST http://localhost:14000/api/v1/zones/example.lan/records \
  -H "Content-Type: application/json" \
  -d '{
    "name": "example.lan.",
    "type": "CAA",
    "ttl": 300,
    "caa_flags": 0,
    "caa_tag": "iodef",
    "caa_value": "mailto:security@example.lan"
  }'
```

**CAA tag values:**

- `issue`: Authorize CA to issue certificates for this domain
- `issuewild`: Authorize CA to issue wildcard certificates
- `iodef`: URL/email for reporting policy violations

### Complete Zone Example

Here's a complete zone with various record types:

```bash
curl -X POST http://localhost:14000/api/v1/zones \
  -H "Content-Type: application/json" \
  -d '{
    "domain": "example.lan",
    "records": [
      {
        "name": "example.lan.",
        "type": "SOA",
        "ttl": 3600,
        "soa_mname": "ns1.example.lan.",
        "soa_rname": "hostmaster.example.lan.",
        "soa_serial": 2024110601,
        "soa_refresh": 3600,
        "soa_retry": 1800,
        "soa_expire": 604800,
        "soa_minimum": 300
      },
      {
        "name": "example.lan.",
        "type": "NS",
        "ttl": 3600,
        "value": "ns1.example.lan."
      },
      {
        "name": "ns1.example.lan.",
        "type": "A",
        "ttl": 3600,
        "value": "192.168.1.1"
      },
      {
        "name": "example.lan.",
        "type": "A",
        "ttl": 300,
        "value": "192.168.1.100"
      },
      {
        "name": "www.example.lan.",
        "type": "CNAME",
        "ttl": 300,
        "value": "example.lan."
      },
      {
        "name": "example.lan.",
        "type": "MX",
        "ttl": 300,
        "mx_priority": 10,
        "mx_host": "mail.example.lan."
      },
      {
        "name": "mail.example.lan.",
        "type": "A",
        "ttl": 300,
        "value": "192.168.1.20"
      },
      {
        "name": "example.lan.",
        "type": "TXT",
        "ttl": 300,
        "value": "v=spf1 mx -all"
      },
      {
        "name": "_http._tcp.example.lan.",
        "type": "SRV",
        "ttl": 300,
        "srv_priority": 10,
        "srv_weight": 60,
        "srv_port": 80,
        "srv_target": "www.example.lan."
      },
      {
        "name": "example.lan.",
        "type": "CAA",
        "ttl": 300,
        "caa_flags": 0,
        "caa_tag": "issue",
        "caa_value": "letsencrypt.org"
      }
    ]
  }'
```

---

## Error Responses

All error responses follow this format:

```json
{
  "error": "Error message describing what went wrong"
}
```

### HTTP Status Codes

- **200 OK** - Request succeeded
- **201 Created** - Resource created successfully
- **204 No Content** - Request succeeded with no response body
- **400 Bad Request** - Invalid request data
- **404 Not Found** - Resource not found
- **405 Method Not Allowed** - HTTP method not supported for this endpoint
- **409 Conflict** - Resource already exists
- **500 Internal Server Error** - Server error

---

## Configuration

The HTTP API server can be configured using environment variables:

| Variable        | Default  | Description                              |
| --------------- | -------- | ---------------------------------------- |
| `HTTP_API_PORT` | `:14000` | Port for the HTTP API server             |
| `LOG_LEVEL`     | `info`   | Logging level (debug, info, warn, error) |
| `LOG_JSON`      | `true`   | Enable JSON structured logging           |

---

## Notes

- Domain names are automatically normalized to end with a trailing dot (`.`)
- Record types are case-insensitive but stored as uppercase
- Default TTL is 300 seconds (5 minutes) if not specified
- All zone and record operations are immediately persisted to Valkey
- The API uses JSON for all request and response bodies

---

## Swagger/OpenAPI Integration

### Generating Swagger Documentation

The Swagger documentation is auto-generated from code annotations using [swaggo/swag](https://github.com/swaggo/swag).

**Regenerate documentation:**

```bash
make swagger
```

This generates the OpenAPI specification files in the `docs/` directory.

### Using Swagger Annotations

Endpoints are documented with Swagger annotations in the code:

```go
// @Summary Create a new DNS zone
// @Description Create a new DNS zone with optional records
// @Tags Zones
// @Accept json
// @Produce json
// @Param zone body models.DNSZone true "Zone to create"
// @Success 201 {object} models.DNSZone "Zone created"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Router /api/v1/zones [post]
func (r *Router) createZone(w http.ResponseWriter, req *http.Request) {
    // implementation
}
```

### Exporting API Specification

The OpenAPI specification can be exported for use with other tools:

**Import to Postman:**

1. Open Postman
2. File → Import
3. Select `docs/swagger.json`
4. All endpoints imported as a collection

**Generate Client SDK:**

```bash
# Install OpenAPI Generator
npm install -g @openapitools/openapi-generator-cli

# Generate Python client
openapi-generator-cli generate \
  -i docs/swagger.json \
  -g python \
  -o clients/python

# Generate JavaScript client
openapi-generator-cli generate \
  -i docs/swagger.json \
  -g javascript \
  -o clients/javascript
```

### Swagger Troubleshooting

**Swagger UI not loading?**

- Ensure API server is running
- Check URL: http://localhost:14000/swagger/index.html
- Verify `docs/docs.go` exists

**Documentation not updating?**

1. Run `make swagger` to regenerate
2. Rebuild binary with `make build-api`
3. Restart the server

**Missing endpoints in Swagger UI?**

1. Check that handler has proper annotations
2. Regenerate docs: `make swagger`
3. Rebuild binary: `make build-api`
4. Restart server and refresh browser

---

## Development Workflow

### Adding New Endpoints

1. **Create handler with Swagger annotations:**

   ```go
   // @Summary Your endpoint summary
   // @Description Detailed description
   // @Tags YourTag
   // @Accept json
   // @Produce json
   // @Param name path string true "Parameter description"
   // @Success 200 {object} YourModel "Success description"
   // @Failure 400 {object} map[string]string "Error description"
   // @Router /your/path [get]
   func (r *Router) yourHandler(w http.ResponseWriter, req *http.Request) {
       // implementation
   }
   ```

2. **Regenerate Swagger docs:**

   ```bash
   make swagger
   ```

3. **Rebuild and test:**

   ```bash
   make build-dns
   ./bin/godns
   ```

4. **Verify in Swagger UI:**
   Visit http://localhost:14000/swagger/index.html

### Testing Scripts

**Test API endpoints:**

```bash
./scripts/test-api.sh
```

This script demonstrates basic CRUD operations using the HTTP API.

---

## Additional Resources

- **API Release Workflow:** [docs/API_RELEASE_WORKFLOW.md](./API_RELEASE_WORKFLOW.md)
- **Quick Start Guide:** [docs/QUICK_START.md](./QUICK_START.md)
- **Swaggo Documentation:** https://github.com/swaggo/swag
- **OpenAPI Specification:** https://swagger.io/specification/

```

```
