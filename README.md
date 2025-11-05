# GoDNS

[![CI](https://github.com/rogerwesterbo/godns/actions/workflows/ci.yml/badge.svg)](https://github.com/rogerwesterbo/godns/actions/workflows/ci.yml)
[![Release](https://github.com/rogerwesterbo/godns/actions/workflows/release.yml/badge.svg)](https://github.com/rogerwesterbo/godns/actions/workflows/release.yml)
[![Security](https://img.shields.io/badge/security-distroless-blue)](https://github.com/GoogleContainerTools/distroless)
[![Go Version](https://img.shields.io/github/go-mod/go-version/rogerwesterbo/godns)](https://go.dev/)

A high-performance DNS server written in Go with Valkey (Redis) backend for dynamic configuration.

## Features

- 🚀 Fast DNS resolution with caching
- 🔧 Dynamic configuration via Valkey
- 🔐 OAuth2/OIDC authentication with Keycloak
- 🌐 REST API with Swagger documentation
- 💻 **Web UI** - Modern React-based management interface
- 🏥 Built-in health checks (liveness/readiness)
- 🔒 ACL-based Valkey authentication
- ☸️ Kubernetes-ready (multi-pod safe)
- 🛠️ CLI tool for testing and management

## Quick Start

**New to GoDNS?** See the [Quick Start Guide](docs/QUICK_START.md) for a 5-minute setup walkthrough.

### Installation

#### Using Docker

```bash
docker pull ghcr.io/rogerwesterbo/godns:latest
docker run -p 53:53/udp -p 53:53/tcp ghcr.io/rogerwesterbo/godns:latest
```

#### Using Helm (Kubernetes)

```bash
# Install DNS server
helm install godns oci://ghcr.io/rogerwesterbo/helm/godns

# Install Web UI
helm install godnsweb oci://ghcr.io/rogerwesterbo/helm/godnsweb \
  --set ingress.enabled=true \
  --set ingress.hosts[0].host=godns.example.com
```

#### Download CLI Binary

Download the latest release for your platform from the [releases page](https://github.com/rogerwesterbo/godns/releases/latest).

**Linux/macOS:**

```bash
# Download (replace VERSION and PLATFORM)
curl -LO https://github.com/rogerwesterbo/godns/releases/download/v1.0.0/godnscli-1.0.0-linux-amd64.tar.gz
tar -xzf godnscli-1.0.0-linux-amd64.tar.gz
sudo mv godnscli-1.0.0-linux-amd64 /usr/local/bin/godnscli
chmod +x /usr/local/bin/godnscli
```

**Windows:**

```powershell
# Download from releases page and add to PATH
```

### Start the Server

```bash
# Using Docker Compose (recommended)
docker-compose up -d

# Or build and run locally
make build
./bin/godns
```

### Test with CLI

```bash
# Build the CLI tool
make build-cli

# Login to API
./bin/godnscli login

# Check authentication status
./bin/godnscli status

# Export zones
./bin/godnscli export --format bind

# Query DNS directly
./bin/godnscli query example.lan

# Run tests
./bin/godnscli test
```

## Authentication

GoDNS uses **OAuth2/OIDC** authentication via Keycloak for API access:

```bash
# Start all services (includes Keycloak)
docker-compose up -d

# Login with CLI
./bin/godnscli login

# Or get token for direct API access
TOKEN=$(curl -s -X POST "http://localhost:14101/realms/godns/protocol/openid-connect/token" \
  -d "client_id=godns-cli" \
  -d "username=testuser" \
  -d "password=password" \
  -d "grant_type=password" | jq -r '.access_token')

# Use token with API
curl -H "Authorization: Bearer $TOKEN" http://localhost:14000/api/v1/zones
```

**Default credentials:**

- Keycloak Admin: `admin` / `admin` (http://localhost:14101)
- Test User: `testuser` / `password`

See [Authentication Guide](docs/AUTHENTICATION.md) for complete details.

## Documentation

### Getting Started

- **[Quick Start Guide](docs/QUICK_START.md)** - 5-minute setup walkthrough
- **[Authentication Guide](docs/AUTHENTICATION.md)** - OAuth2/OIDC setup
- **[Quick Auth Reference](docs/QUICK_AUTH_REFERENCE.md)** - Fast auth commands

### API & CLI

- **[API Documentation](docs/API_DOCUMENTATION.md)** - REST API reference
- **[CLI Guide](docs/CLI_GUIDE.md)** - Complete guide for using godnscli
- **[CLI Config](docs/CLI_CONFIG.md)** - CLI configuration management

### Web UI

- **[Web UI Documentation](web/godnsweb/README.md)** - React-based management interface
- **[Web UI Deployment](web/godnsweb/docs/DOCKER.md)** - Docker and Kubernetes setup

### Configuration

- **[Port Configuration](docs/PORT_CONFIGURATION.md)** - Port mappings and setup
- **[Valkey Authentication](docs/VALKEY_AUTH.md)** - Valkey auth setup guide

## Building

```bash
# Build server only
make build

# Build CLI only
make build-cli

# Build everything
make build-all

# Build Docker image
make docker-build

# Build multi-arch Docker image (requires buildx)
make docker-build-multiarch

# Run linting
make lint
```

## Docker & Container Images

GoDNS uses **Google Distroless** base images for maximum security:

- ✅ Minimal attack surface (no shell, no package manager)
- ✅ Runs as non-root user (UID 65532)
- ✅ Read-only root filesystem
- ✅ Multi-arch support (linux/amd64, linux/arm64)
- ✅ Signed with Cosign
- ✅ Includes SBOM (Software Bill of Materials)

### Available Images

**DNS Server:**

```bash
# Latest release
ghcr.io/rogerwesterbo/godns:latest

# Specific version
ghcr.io/rogerwesterbo/godns:1.0.0
```

**Web UI:**

```bash
# Latest release
ghcr.io/rogerwesterbo/godns-web:latest

# Specific version
ghcr.io/rogerwesterbo/godns-web:1.0.0

# Run Web UI
docker run -p 8080:8080 \
  -e VITE_KEYCLOAK_URL=http://keycloak:8080 \
  -e VITE_API_BASE_URL=http://godns-api:8080 \
  ghcr.io/rogerwesterbo/godns-web:latest
```

### Security Scanning

All images are scanned with Trivy before release. View scan results in the [Security tab](https://github.com/rogerwesterbo/godns/security).

## CLI Usage

The `godnscli` tool provides easy testing and management:

```bash
# Quick test
./bin/godnscli t

# Query a domain
./bin/godnscli q example.lan

# Check health
./bin/godnscli h

# Show version
./bin/godnscli v
```

See the [CLI Guide](docs/CLI_GUIDE.md) for complete documentation.

## Configuration

Configuration is managed via environment variables. See `.env.example` for available options.

### Key Settings

**Valkey:**

- `VALKEY_ADDR` - Valkey server address (default: `localhost:14103`)
- `VALKEY_USERNAME` - Valkey username
- `VALKEY_PASSWORD` - Valkey password

**DNS Server:**

- `DNS_SERVER_PORT` - DNS server port (default: `53`)

**HTTP API:**

- `HTTP_API_PORT` - API server port (default: `:14000`)
- `HTTP_API_CORS_ALLOWED_ORIGINS` - CORS allowed origins

**Authentication:**

- `AUTH_ENABLED` - Enable/disable authentication (default: `true`)
- `KEYCLOAK_URL` - Keycloak server URL (default: `http://localhost:14101`)
- `KEYCLOAK_REALM` - OAuth2 realm (default: `godns`)

**Ports:**

- DNS: `53`
- HTTP API: `14000`
- Keycloak: `14101` (HTTP), `14102` (HTTPS)
- Valkey: `14103`
- PostgreSQL: `14100`

See [Port Configuration](docs/PORT_CONFIGURATION.md) for details.

## Development

### Prerequisites

- Go 1.25.3 or later
- Docker & Docker Compose
- Make

### Project Structure

```
godns/
├── cmd/
│   ├── godns/          # Main DNS server
│   ├── godnsapi/       # HTTP API server
│   └── godnscli/       # CLI tool
├── internal/           # Internal packages
│   ├── dnsserver/      # DNS server logic
│   ├── httpserver/     # HTTP API server
│   └── services/       # Business logic
├── pkg/                # Public packages
│   ├── auth/           # Authentication
│   ├── clients/        # External clients (Valkey)
│   └── options/        # Configuration options
├── web/                # Web UI
│   └── godnsweb/       # React application
├── charts/             # Helm charts
│   ├── godns/          # DNS server chart
│   └── godnsweb/       # Web UI chart
├── docs/               # Documentation
└── hack/               # Development utilities
```

## Releases

Releases are automated via GitHub Actions. To create a new release:

1. Update `CHANGELOG.md` with release notes
2. Create and push a new tag:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```
3. Create a [new release](https://github.com/rogerwesterbo/godns/releases/new) in GitHub
4. GitHub Actions will automatically:
   - **Detect changes** - Build only what changed (DNS server, Web UI, or both)
   - **Build DNS Server** (if backend changed):
     - Multi-arch Docker images (linux/amd64, linux/arm64, linux/arm/v7)
     - CLI binaries for Linux, macOS, Windows
     - Helm chart package
   - **Build Web UI** (if web files changed):
     - Multi-arch Web UI Docker image
     - Web UI Helm chart
   - **Security & Quality**:
     - Run security scans (Trivy)
     - Generate SBOM
     - Sign images with Cosign
   - Attach all artifacts to the release

### Release Artifacts

Each release may include (depending on what changed):

**DNS Server:**

- 🐳 Multi-arch Docker images (linux/amd64, linux/arm64, linux/arm/v7)
- 📦 Helm chart package (godns)
- 💻 CLI binaries (Linux, macOS, Windows)

**Web UI:**

- 🌐 Web UI Docker image (linux/amd64, linux/arm64)
- 📦 Helm chart package (godnsweb)

**Common:**

- 🔒 SHA256 checksums
- 📋 SBOM (SPDX format)
- 🛡️ Security scan results

See the [CHANGELOG](CHANGELOG.md) for release history and the [Security Policy](SECURITY.md) for security information.

## License

See [LICENSE](LICENSE) file.
