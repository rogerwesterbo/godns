# GoDNS

[![CI](https://github.com/rogerwesterbo/godns/actions/workflows/ci.yml/badge.svg)](https://github.com/rogerwesterbo/godns/actions/workflows/ci.yml)

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

# Run linting
make lint
```

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

## License

See [LICENSE](LICENSE) file.
