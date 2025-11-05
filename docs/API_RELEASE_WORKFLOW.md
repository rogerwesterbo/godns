# API Release Workflow

## Overview

The GoDNS project now includes automated Swagger documentation generation integrated into the build and release process. This ensures that API documentation is always up-to-date in released containers.

## Workflow

### Local Development

1. **Make code changes** to API handlers or models
2. **Regenerate Swagger docs**:
   ```bash
   make swagger
   ```
3. **Test locally**:
   ```bash
   make build-api
   ./bin/godnsapi
   ```
4. **View Swagger UI** at http://localhost:14000/swagger/index.html

### Docker Build

The Docker build automatically generates Swagger documentation:

```bash
make docker-build
```

This will:

1. ✓ Generate Swagger docs locally (via `swagger` target dependency)
2. ✓ Install swag in the build container
3. ✓ Generate Swagger docs inside the container
4. ✓ Build the Go binary with embedded docs
5. ✓ Create the final minimal container image

### Release Process

#### Single Platform Release

```bash
make release
```

This will:

1. Generate Swagger documentation
2. Build Docker image for current platform
3. Push to container registry (ghcr.io/rogerwesterbo/godns:latest)

#### Multi-Platform Release

```bash
make release-multiarch
```

This will:

1. Generate Swagger documentation
2. Build Docker image for multiple platforms (linux/amd64, linux/arm64)
3. Push to container registry with platform-specific tags

## Makefile Targets

### Documentation

- `make swagger` - Generate Swagger docs in `docs/` directory
  - Creates: `docs/docs.go`, `docs/swagger.json`, `docs/swagger.yaml`
  - Uses: `swag init -g cmd/godnsapi/swagger.go --output docs`

### Building

- `make build-dns` - Build DNS server (`cmd/godns`)
- `make build-api` - Build API server (`cmd/godnsapi`)
- `make build-cli` - Build CLI tool (`cmd/godnscli`)
- `make build-all` - Build all binaries

### Docker

- `make docker-build` - Build Docker image (single platform)
  - Depends on: `swagger` target
  - Generates docs both locally and in container
- `make docker-build-multiarch` - Build multi-platform Docker image

  - Platforms: linux/amd64, linux/arm64
  - Depends on: `swagger` target

- `make docker-push` - Push image to registry
- `make docker-run` - Run container locally

### Release

- `make release` - Full release (swagger → build → push)
- `make release-multiarch` - Multi-platform release

## File Locations

### Source Files

- `cmd/godnsapi/swagger.go` - Main Swagger annotations
- `internal/httpserver/httproutes/http_routes.go` - Route annotations
- `internal/models/dns_record.go` - Model definitions with examples

### Generated Files

- `docs/docs.go` - Embedded Go code for Swagger spec
- `docs/swagger.json` - OpenAPI 3.0 JSON spec
- `docs/swagger.yaml` - OpenAPI 3.0 YAML spec

### Build Files

- `Dockerfile` - Multi-stage build with Swagger generation
- `Makefile` - Build automation with colored output

## API Endpoints

All endpoints are documented in Swagger UI:

### Zones

- `GET /api/v1/zones` - List all zones
- `POST /api/v1/zones` - Create a zone
- `GET /api/v1/zones/{domain}` - Get a zone
- `PUT /api/v1/zones/{domain}` - Update a zone
- `DELETE /api/v1/zones/{domain}` - Delete a zone

### Records

- `POST /api/v1/zones/{domain}/records` - Create a record
- `GET /api/v1/zones/{domain}/records/{name}/{type}` - Get a record
- `PUT /api/v1/zones/{domain}/records/{name}/{type}` - Update a record
- `DELETE /api/v1/zones/{domain}/records/{name}/{type}` - Delete a record

### Health

- `GET /health` - Health check
- `GET /ready` - Readiness check

## Environment Variables

- `HTTP_API_PORT` - API server port (default: `:14000`)
- `VALKEY_ADDR` - Valkey server address (default: `localhost:14103`)
- `VALKEY_PASSWORD` - Valkey password (optional)
- `VALKEY_DB` - Valkey database number (default: `0`)

## Best Practices

1. **Always regenerate docs** after API changes:

   ```bash
   make swagger
   ```

2. **Update annotations** when modifying handlers:

   - Add `@Summary`, `@Description`, `@Tags`
   - Define `@Param` for all parameters
   - Specify `@Success` and `@Failure` responses
   - Include `@Router` path and method

3. **Test Swagger UI** before committing:

   - Run `make build-api && ./bin/godnsapi`
   - Visit http://localhost:14000/swagger/index.html
   - Verify all endpoints are documented correctly

4. **Use release targets** for production deployments:
   - Single platform: `make release`
   - Multi-platform: `make release-multiarch`

## Troubleshooting

### Swagger Generation Fails

```bash
# Clear docs and regenerate
rm -rf docs/*.go docs/*.json docs/*.yaml
make swagger
```

### Docker Build Fails

```bash
# Ensure local docs are generated first
make swagger

# Clean Docker build cache
docker builder prune -af

# Rebuild
make docker-build
```

### Missing Endpoints in Swagger UI

1. Check that handler has proper annotations
2. Regenerate docs: `make swagger`
3. Rebuild binary: `make build-api`
4. Restart server and refresh browser

## Version Information

- Swagger/OpenAPI: 3.0
- swaggo/swag: v1.16.6
- swaggo/http-swagger: v1.3.4
- Go: 1.25.3
- Docker base images:
  - Builder: golang:1.25.3-alpine
  - Runtime: gcr.io/distroless/static-debian12:nonroot
