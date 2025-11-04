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

**URL:** http://localhost:14082/swagger/index.html

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
   # Standalone API server
   make build-api
   ./bin/godnsapi

   # Or integrated with DNS server
   make build-dns
   ./bin/godns
   ```

2. **Open Swagger UI:**
   Visit http://localhost:14082/swagger/index.html

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

By default, the API server runs on port `13082`. You can configure this using the `HTTP_API_PORT` environment variable.

```
http://localhost:14082
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

```json
{
  "name": "string (required)",    // Fully qualified domain name
  "type": "string (required)",    // Record type: A, AAAA, CNAME, MX, NS, TXT, PTR, SRV, SOA, CAA
  "ttl": number,                  // Time to live in seconds (default: 300)
  "value": "string (required)"    // Record value (IP address, hostname, text, etc.)
}
```

### Supported Record Types

- **A** - IPv4 address
- **AAAA** - IPv6 address
- **CNAME** - Canonical name (alias)
- **MX** - Mail exchange
- **NS** - Name server
- **TXT** - Text record
- **PTR** - Pointer record
- **SRV** - Service record
- **SOA** - Start of authority
- **CAA** - Certification authority authorization

---

## Example Usage

### Using curl

#### Create a zone with records:

```bash
curl -X POST http://localhost:14082/api/v1/zones \
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

#### List all zones:

```bash
curl http://localhost:14082/api/v1/zones
```

#### Get a specific zone:

```bash
curl http://localhost:14082/api/v1/zones/example.lan
```

#### Add a record to a zone:

```bash
curl -X POST http://localhost:14082/api/v1/zones/example.lan/records \
  -H "Content-Type: application/json" \
  -d '{
    "name": "api.example.lan.",
    "type": "A",
    "ttl": 300,
    "value": "192.168.1.150"
  }'
```

#### Update a record:

```bash
curl -X PUT http://localhost:14082/api/v1/zones/example.lan/records/www.example.lan./A \
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
curl -X DELETE http://localhost:14082/api/v1/zones/example.lan/records/api.example.lan./A
```

#### Delete a zone:

```bash
curl -X DELETE http://localhost:14082/api/v1/zones/example.lan
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

| Variable        | Default | Description                              |
| --------------- | ------- | ---------------------------------------- |
| `HTTP_API_PORT` | `:13082` | Port for the HTTP API server             |
| `LOG_LEVEL`     | `info`  | Logging level (debug, info, warn, error) |
| `LOG_JSON`      | `true`  | Enable JSON structured logging           |

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
- Check URL: http://localhost:14082/swagger/index.html
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
   make build-api
   ./bin/godnsapi
   ```

4. **Verify in Swagger UI:**
   Visit http://localhost:14082/swagger/index.html

### Testing Scripts

**Test API endpoints:**

```bash
./hack/test-api.sh
```

**Test Swagger UI:**

```bash
./hack/test-swagger.sh
```

---

## Additional Resources

- **API Release Workflow:** [docs/API_RELEASE_WORKFLOW.md](./API_RELEASE_WORKFLOW.md)
- **Quick Start Guide:** [docs/QUICK_START.md](./QUICK_START.md)
- **Swaggo Documentation:** https://github.com/swaggo/swag
- **OpenAPI Specification:** https://swagger.io/specification/

```

```
