# GoDNS HTTP API Server

This is the HTTP API server component of GoDNS, providing RESTful endpoints for managing DNS zones and records.

## Interactive Documentation

The API includes **Swagger/OpenAPI documentation** with an interactive UI. Once the server is running, access it at:

# GoDNS HTTP API Server

RESTful HTTP API for managing DNS zones and records.

## Quick Start

```bash
# Build
make build-api

# Run
./bin/godnsapi
```

The API server will start on `http://localhost:8082`

## Interactive Documentation

Access the **Swagger UI** for interactive API documentation:

**http://localhost:8082/swagger/index.html**

Features:

- 📖 Browse all API endpoints
- 🧪 Test endpoints directly in your browser
- 📋 View request/response schemas
- 💡 See example values and use cases

## Quick API Examples

```bash
# List all zones
curl http://localhost:8082/api/v1/zones

# Create a zone
curl -X POST http://localhost:8082/api/v1/zones \
  -H "Content-Type: application/json" \
  -d '{"domain":"example.lan","records":[{"name":"www.example.lan.","type":"A","ttl":300,"value":"192.168.1.100"}]}'

# Get a zone
curl http://localhost:8082/api/v1/zones/example.lan

# Delete a zone
curl -X DELETE http://localhost:8082/api/v1/zones/example.lan
```

## Full Documentation

For complete documentation, see:

- **[API Documentation](../../docs/API_DOCUMENTATION.md)** - Complete API reference with all endpoints
- **[API Release Workflow](../../docs/API_RELEASE_WORKFLOW.md)** - Build, release, and deployment guide
- **[Quick Start](../../docs/QUICK_START.md)** - Getting started with GoDNS

## Configuration

Environment variables:

- `HTTP_API_PORT` - Port to listen on (default: `:8082`)
- `VALKEY_ADDR` - Valkey server address (default: `localhost:6379`)
- `VALKEY_PASSWORD` - Valkey password (optional)
- `VALKEY_DB` - Valkey database number (default: `0`)

## Testing

```bash
# Test API endpoints
./hack/test-api.sh

# Test Swagger UI
./hack/test-swagger.sh
```

## Learn More

Visit the [main documentation](../../docs/README.md) for comprehensive guides and references.

Features:

- Interactive API exploration
- Try-it-out functionality for all endpoints
- Complete request/response schemas
- Example values and use cases

## Quick Start

### Running the API Server

#### Option 1: Integrated with DNS Server (Current Default)

The HTTP API is currently integrated into the main `godns` binary:

```bash
# Build
make build-dns

# Run (starts both DNS and HTTP API servers)
./bin/godns
```

The API will be available at `http://localhost:8082` (configurable via `HTTP_API_PORT`).

#### Option 2: Standalone API Server (Future/Recommended)

For production deployments, run the API server separately:

```bash
# Build
make build-api

# Run
./bin/godnsapi
```

This allows independent scaling and deployment of the API server.

## Configuration

Configure the API server using environment variables:

| Variable          | Default     | Description                              |
| ----------------- | ----------- | ---------------------------------------- |
| `HTTP_API_PORT`   | `:8082`     | Port for the HTTP API server             |
| `VALKEY_HOST`     | `localhost` | Valkey server host                       |
| `VALKEY_PORT`     | `6379`      | Valkey server port                       |
| `VALKEY_USERNAME` | -           | Valkey username (if auth enabled)        |
| `VALKEY_TOKEN`    | -           | Valkey password/token (if auth enabled)  |
| `LOG_LEVEL`       | `info`      | Logging level (debug, info, warn, error) |
| `LOG_JSON`        | `true`      | Enable JSON structured logging           |

## Example Usage

### Create a DNS Zone

```bash
curl -X POST http://localhost:8082/api/v1/zones \
  -H "Content-Type: application/json" \
  -d '{
    "domain": "example.lan",
    "records": [
      {
        "name": "www.example.lan.",
        "type": "A",
        "ttl": 300,
        "value": "192.168.1.100"
      }
    ]
  }'
```

### List All Zones

```bash
curl http://localhost:8082/api/v1/zones
```

### Get a Specific Zone

```bash
curl http://localhost:8082/api/v1/zones/example.lan
```

### Add a Record

```bash
curl -X POST http://localhost:8082/api/v1/zones/example.lan/records \
  -H "Content-Type: application/json" \
  -d '{
    "name": "api.example.lan.",
    "type": "A",
    "ttl": 300,
    "value": "192.168.1.150"
  }'
```

### Update a Record

```bash
curl -X PUT http://localhost:8082/api/v1/zones/example.lan/records/www.example.lan./A \
  -H "Content-Type: application/json" \
  -d '{
    "name": "www.example.lan.",
    "type": "A",
    "ttl": 600,
    "value": "192.168.1.200"
  }'
```

### Delete a Record

```bash
curl -X DELETE http://localhost:8082/api/v1/zones/example.lan/records/api.example.lan./A
```

### Delete a Zone

```bash
curl -X DELETE http://localhost:8082/api/v1/zones/example.lan
```

## Testing

Run the automated test script:

```bash
./hack/test-api.sh
```

This script demonstrates all CRUD operations on zones and records.

## Documentation

- [Interactive Swagger UI](http://localhost:8082/swagger/index.html) - Try the API in your browser
- [Full API Documentation](../docs/API_DOCUMENTATION.md) - Complete endpoint reference
- [Swagger Guide](../docs/SWAGGER_GUIDE.md) - How to use and customize Swagger docs
- [Implementation Details](../docs/HTTP_API_IMPLEMENTATION.md) - Architecture and design

## Development

### Building

```bash
# Build API server only
make build-api

# Build all components
make build-all

# Generate Swagger documentation
make swagger
```

### Running Locally

1. Start Valkey:

```bash
docker-compose up -d valkey
```

2. Start the API server:

```bash
./bin/godnsapi
```

3. Test the API:

```bash
./hack/test-api.sh
```

## Architecture

The API server uses:

- **Router**: Standard Go `net/http` with `http.ServeMux`
- **Storage**: Valkey (Redis-compatible) for persistent storage
- **Service Layer**: `v1zoneservice` for business logic
- **Models**: `DNSZone` and `DNSRecord` data structures

## Deployment

### Docker

Build and run as a Docker container:

```bash
# Build
docker build -t godnsapi .

# Run
docker run -p 8082:8082 \
  -e VALKEY_HOST=valkey \
  -e VALKEY_PORT=6379 \
  godnsapi
```

### Kubernetes

Deploy separately from the DNS server for better scalability:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: godnsapi
spec:
  replicas: 3 # Scale independently
  selector:
    matchLabels:
      app: godnsapi
  template:
    metadata:
      labels:
        app: godnsapi
    spec:
      containers:
        - name: godnsapi
          image: ghcr.io/rogerwesterbo/godns:latest
          command: ["/godnsapi"] # Use standalone binary
          ports:
            - containerPort: 8082
          env:
            - name: HTTP_API_PORT
              value: ":8082"
            - name: VALKEY_HOST
              value: "valkey"
```

## Security

⚠️ **Current State**: The API has no authentication or authorization.

For production use, consider:

- Adding API key authentication
- Implementing TLS/HTTPS
- Adding rate limiting
- Using network policies to restrict access

## Contributing

When adding new endpoints:

1. Add route in `internal/httpserver/httproutes/http_routes.go`
2. Add business logic in `internal/services/v1zoneservice/v1_zone_service.go`
3. Update API documentation in `docs/API_DOCUMENTATION.md`
4. Add tests

## License

See [LICENSE](../LICENSE) file.
