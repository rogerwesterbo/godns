# GoDNS

[![CI](https://github.com/rogerwesterbo/godns/actions/workflows/ci.yml/badge.svg)](https://github.com/rogerwesterbo/godns/actions/workflows/ci.yml)
[![Release](https://github.com/rogerwesterbo/godns/actions/workflows/release.yml/badge.svg)](https://github.com/rogerwesterbo/godns/actions/workflows/release.yml)
[![Security](https://img.shields.io/badge/security-distroless-blue)](https://github.com/GoogleContainerTools/distroless)
[![Go Version](https://img.shields.io/github/go-mod/go-version/rogerwesterbo/godns)](https://go.dev/)

A high-performance DNS server written in Go with Valkey (Redis) backend for dynamic configuration.

## Features

- 🚀 Fast DNS resolution with caching
- 🔧 Dynamic configuration via Valkey
- 🏥 Built-in health checks (liveness/readiness)
- 🔒 ACL-based authentication
- ☸️ Kubernetes-ready (multi-pod safe)
- 🛠️ CLI tool for testing and management

## Quick Start

**New to GoDNS?** See the [Quick Start Guide](docs/QUICK_START.md) for a 5-minute setup walkthrough.

### Installation

#### Using Docker

```bash
docker pull ghcr.io/rogerwesterbo/godns:latest
docker run -p 53:53/udp -p 53:53/tcp -p 8080:8080 ghcr.io/rogerwesterbo/godns:latest
```

#### Using Helm (Kubernetes)

```bash
helm install godns oci://ghcr.io/rogerwesterbo/helm/godns
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

# Run tests
./bin/godnscli test

# Query a domain
./bin/godnscli query example.lan
```

## Documentation

- **[Quick Start Guide](docs/QUICK_START.md)** - 5-minute setup walkthrough
- **[CLI Guide](docs/CLI_GUIDE.md)** - Complete guide for using godnscli
- **[CLI Quick Reference](docs/CLI_QUICK_REFERENCE.md)** - Fast command lookup
- **[CLI Cheat Sheet](docs/CLI_CHEAT_SHEET.md)** - One-page reference
- **[Valkey Authentication](docs/VALKEY_AUTH.md)** - Authentication setup guide

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

```bash
# Latest release
ghcr.io/rogerwesterbo/godns:latest

# Specific version
ghcr.io/rogerwesterbo/godns:1.0.0

# With digest for immutability
ghcr.io/rogerwesterbo/godns:1.0.0@sha256:abc123...
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

Key settings:

- `VALKEY_ADDR` - Valkey server address
- `VALKEY_USERNAME` - Valkey username
- `VALKEY_PASSWORD` - Valkey password
- `DNS_PORT` - DNS server port (default: 53)

## Development

### Prerequisites

- Go 1.25.3 or later
- Docker & Docker Compose
- Make

### Project Structure

```
godns/
├── cmd/
│   ├── godns/          # Main server
│   └── godnscli/       # CLI tool
├── internal/           # Internal packages
│   ├── handlers/       # DNS request handlers
│   └── services/       # Business logic
├── pkg/                # Public packages
│   ├── clients/        # External clients (Valkey)
│   └── options/        # Configuration options
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
   - Build multi-arch Docker images (linux/amd64, linux/arm64)
   - Push images to `ghcr.io/rogerwesterbo/godns`
   - Package Helm chart and push to `ghcr.io/rogerwesterbo/helm/godns`
   - Build CLI binaries for:
     - Linux (amd64, arm64)
     - macOS (Intel, Apple Silicon)
     - Windows (amd64)
   - Run security scans (Trivy)
   - Generate SBOM
   - Attach all artifacts to the release

### Release Artifacts

Each release includes:

- 🐳 Multi-arch Docker images
- 📦 Helm chart package
- 💻 CLI binaries for all platforms
- 🔒 SHA256 checksums
- 📋 SBOM (SPDX format)
- 🛡️ Security scan results

See the [CHANGELOG](CHANGELOG.md) for release history and the [Security Policy](SECURITY.md) for security information.

## License

See [LICENSE](LICENSE) file.
