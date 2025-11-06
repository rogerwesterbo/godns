# GoDNS Helper Scripts

This directory contains helper scripts for development, testing, and setup.

## Setup & Initialization

### `init-keycloak.sh`

Initializes Keycloak with the GoDNS realm, clients, and test users.

**Usage:**

```bash
./scripts/init-keycloak.sh
```

**What it does:**

- Creates the `godns` realm
- Sets up API client (`godns-api`)
- Sets up CLI client (`godns-cli`)
- Sets up web client (`godns-web`)
- Creates test user with roles
- Configures client credentials

**Environment variables:**

- `KEYCLOAK_URL` (default: `https://keycloak:8443`)
- `KEYCLOAK_ADMIN` (default: `admin`)
- `KEYCLOAK_ADMIN_PASSWORD` (default: `admin`)
- `KEYCLOAK_REALM` (default: `godns`)

---

### `init-web-client.sh`

Sets up the web client in Keycloak with PKCE configuration for the React app.

**Usage:**

```bash
./scripts/init-web-client.sh
```

**What it does:**

- Creates `godns-web` client with public access
- Configures PKCE (Proof Key for Code Exchange) for security
- Sets up redirect URIs for localhost development
- Enables CORS for the web application

**Environment variables:**

- `KEYCLOAK_URL` (default: `http://localhost:14101`)
- `REALM` (default: `godns`)
- `ADMIN_USER` (default: `admin`)
- `ADMIN_PASSWORD` (default: `admin`)

---

## Testing & Development

### `seed-test-data.sh`

Seeds the database with realistic DNS test data across multiple zones.

**Usage:**

```bash
./scripts/seed-test-data.sh
```

**What it creates:**

- **6 DNS zones** with 86 total records
- **home.lan** (13 records): Home network setup
- **dev.local** (13 records): Development environment
- **k8s.local** (13 records): Kubernetes cluster
- **docker.local** (12 records): Docker services
- **lab.internal** (13 records): Lab environment with load balancing
- **example.lan** (22 records): Complete demo with email security

**Record types included:**

- SOA (Start of Authority)
- NS (Name Servers)
- A (IPv4 addresses)
- AAAA (IPv6 addresses)
- CNAME (Aliases)
- MX (Mail Exchange)
- SRV (Service Discovery)
- TXT (SPF, DKIM, DMARC)
- Wildcards (_.apps, _.test)

**Prerequisites:**

- GoDNS API running on `http://localhost:14000`

---

### `test-api.sh`

Demonstrates basic CRUD operations using the GoDNS HTTP API.

**Usage:**

```bash
./scripts/test-api.sh
```

**What it tests:**

1. Health check endpoint
2. Creating a DNS zone
3. Listing all zones
4. Getting a specific zone
5. Adding records to a zone
6. Updating records
7. Deleting records
8. Deleting zones

**Output:**
JSON responses for each API operation using `jq` for formatting.

---

### `export-zones-example.sh`

Demonstrates how to export DNS zones in various formats.

**Usage:**

```bash
./scripts/export-zones-example.sh
```

**What it does:**

- Lists all available zones
- Exports all zones in BIND format
- Exports all zones in CoreDNS format
- Exports all zones in PowerDNS format
- Exports individual zones in different formats

**Output directory:**
`./exported-zones/` (created automatically)

**Environment variables:**

- `GODNS_API_URL` (default: `http://localhost:14000`)

---

## Prerequisites

All scripts assume:

1. **jq** is installed (JSON processor)
2. **curl** is installed (for API calls)
3. **docker-compose** is running (for Keycloak scripts)

Install dependencies on macOS:

```bash
brew install jq curl
```

Install dependencies on Ubuntu/Debian:

```bash
sudo apt-get install jq curl
```

---

## Script Organization

Scripts are organized by purpose:

| Script                    | Purpose           | When to Use                  |
| ------------------------- | ----------------- | ---------------------------- |
| `init-keycloak.sh`        | Keycloak setup    | First time setup, reset auth |
| `init-web-client.sh`      | Web client config | Web app development          |
| `seed-test-data.sh`       | Load test data    | Testing, demos, development  |
| `test-api.sh`             | API testing       | Verify API functionality     |
| `export-zones-example.sh` | Export demo       | Learn export API, migration  |

---

## Common Workflows

### First Time Setup

```bash
# 1. Start services
docker-compose up -d

# 2. Initialize Keycloak
./scripts/init-keycloak.sh

# 3. Setup web client
./scripts/init-web-client.sh

# 4. Seed test data
./scripts/seed-test-data.sh
```

### Testing Changes

```bash
# Run API tests
./scripts/test-api.sh

# Reseed test data
./scripts/seed-test-data.sh
```

### Export for Migration

```bash
# Export all zones in different formats
./scripts/export-zones-example.sh

# Check exported files
ls -la ./exported-zones/
```

---

## Troubleshooting

**API not reachable:**

```bash
# Check if API is running
curl http://localhost:14000/health

# Check docker-compose services
docker-compose ps
```

**Keycloak not ready:**

```bash
# Check Keycloak status
curl http://localhost:14101/health/ready

# Check Keycloak logs
docker-compose logs keycloak
```

**Permission denied:**

```bash
# Make scripts executable
chmod +x scripts/*.sh
```

---

## See Also

- [API Documentation](../docs/API_DOCUMENTATION.md) - Complete API reference
- [Authentication Guide](../docs/AUTHENTICATION.md) - Keycloak setup details
- [Testing Guide](../docs/TESTING_GUIDE.md) - Comprehensive testing guide
- [Quick Start](../docs/QUICK_START.md) - Get up and running quickly
