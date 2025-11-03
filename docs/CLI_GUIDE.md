# GoDNS CLI Guide

Complete guide for using the `godnscli` command-line tool to test and manage your GoDNS server.

## Quick Reference Cheat Sheet

### Most Common Commands

```bash
./bin/godnscli d                        # Discover server and network info
./bin/godnscli t                        # Test everything
./bin/godnscli q example.lan            # Query a domain
./bin/godnscli h                        # Check health
./bin/godnscli v                        # Show version
```

### Query Different Record Types

```bash
./bin/godnscli q example.lan            # A record (IPv4)
./bin/godnscli q example.lan -t AAAA    # IPv6
./bin/godnscli q example.lan -t MX      # Mail servers
./bin/godnscli q example.lan -t NS      # Name servers
./bin/godnscli q example.lan -t TXT     # Text records
```

### Query Different Servers

```bash
./bin/godnscli q example.lan -s localhost:53         # Local server
./bin/godnscli q example.lan -s 192.168.1.100:53     # Custom server
./bin/godnscli q google.com -s 8.8.8.8:53            # Google DNS
```

## Table of Contents

- [Quick Reference Cheat Sheet](#quick-reference-cheat-sheet)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Command Reference](#command-reference)
- [Common Use Cases](#common-use-cases)
- [Advanced Usage](#advanced-usage)
- [Troubleshooting](#troubleshooting)

## Installation

### Build from Source

```bash
# Build the CLI tool
make build-cli

# The binary will be at ./bin/godnscli
```

### Verify Installation

```bash
./bin/godnscli version
# or use the short alias
./bin/godnscli v
```

## Quick Start

### 1. Start Your GoDNS Server

```bash
# Start with docker-compose
docker-compose up -d

# Or build and run locally
make build
./bin/godns
```

### 2. Run a Quick Test

```bash
# Test all functionality
./bin/godnscli test

# Or use the short alias
./bin/godnscli t
```

### 3. Query a Domain

```bash
# Query an A record
./bin/godnscli query google.com

# Or use the short alias
./bin/godnscli q google.com
```

## Command Reference

### Global Flags

All commands support these global flags:

| Flag        | Short | Default        | Description           |
| ----------- | ----- | -------------- | --------------------- |
| `--server`  | `-s`  | `localhost:53` | DNS server address    |
| `--verbose` | `-v`  | `false`        | Enable verbose output |
| `--help`    | `-h`  | -              | Show help information |

### Commands Overview

| Command   | Alias | Description                    |
| --------- | ----- | ------------------------------ |
| `query`   | `q`   | Query DNS records for a domain |
| `health`  | `h`   | Check server health status     |
| `test`    | `t`   | Run comprehensive test suite   |
| `version` | `v`   | Display version information    |

---

### `query` (alias: `q`)

Query DNS records for a specific domain.

**Usage:**

```bash
./bin/godnscli query [domain] [flags]
```

**Flags:**

- `-t, --type`: Query type (default: "A")
- `--timeout`: Query timeout in seconds (default: 5)

**Supported Record Types:**

- `A` - IPv4 address
- `AAAA` - IPv6 address
- `MX` - Mail exchange
- `NS` - Name server
- `TXT` - Text records
- `CNAME` - Canonical name
- `SOA` - Start of authority
- `PTR` - Pointer record
- `SRV` - Service record

**Examples:**

```bash
# Basic A record query
./bin/godnscli q example.lan

# Query IPv6 address
./bin/godnscli q example.lan -t AAAA

# Query mail servers
./bin/godnscli q example.lan -t MX

# Query name servers
./bin/godnscli q example.lan -t NS

# Query with verbose output
./bin/godnscli q example.lan -v

# Query remote server
./bin/godnscli q example.lan -s 192.168.1.100:53

# Query with longer timeout
./bin/godnscli q example.lan --timeout 10

# Query public DNS (Google)
./bin/godnscli q google.com -s 8.8.8.8:53
```

**Output Example:**

```
;; Query time: 55.975292ms
;; SERVER: 8.8.8.8:53
;; WHEN: Mon, 03 Nov 2025 00:19:48 CET
;; MSG SIZE rcvd: 54

;; ANSWER SECTION:
google.com.     105     IN      A       142.250.74.46
```

---

### `health` (alias: `h`)

Check the health status of your GoDNS server (liveness and readiness probes).

**Usage:**

```bash
./bin/godnscli health [flags]
```

**Flags:**

- `--liveness-port`: Liveness probe port (default: "8080")
- `--readiness-port`: Readiness probe port (default: "8081")

**Examples:**

```bash
# Check default health endpoints
./bin/godnscli h

# Check with verbose output
./bin/godnscli h -v

# Check custom ports
./bin/godnscli h --liveness-port 9090 --readiness-port 9091

# Check remote server
./bin/godnscli h -s dns.example.com:53
```

**Output Example:**

```
Liveness:  ✓ OK
Readiness: ✓ OK

DNS server is healthy!
```

---

### `test` (alias: `t`)

Run a comprehensive test suite against your DNS server.

**Usage:**

```bash
./bin/godnscli test [flags]
```

**What It Tests:**

- A records (IPv4 resolution)
- AAAA records (IPv6 resolution)
- MX records (Mail exchange)
- NS records (Name servers)
- External resolution (google.com via upstream)

**Examples:**

```bash
# Run all tests
./bin/godnscli t

# Run with verbose output (see all query details)
./bin/godnscli t -v

# Test remote server
./bin/godnscli t -s 192.168.1.100:53

# Test production server
./bin/godnscli t -s dns.example.com:53
```

**Output Example:**

```
Testing DNS server at localhost:53

Running DNS server tests...

Test Results:
  ✓ A record test passed
  ✓ AAAA record test passed
  ✓ MX record test passed
  ✓ NS record test passed
  ✓ External resolution test passed

5 passed, 0 failed
```

---

### `version` (alias: `v`)

Display version, commit, and build information.

**Usage:**

```bash
./bin/godnscli version
```

**Examples:**

```bash
./bin/godnscli v
```

**Output Example:**

```
godnscli version v1.0.0
commit: abc123def
built:  2025-11-03T00:00:00Z
```

## Common Use Cases

### Testing Local Development

When developing locally, you'll want to test your DNS server:

```bash
# 1. Start your server
docker-compose up -d

# 2. Run the test suite
./bin/godnscli t

# 3. Query specific domains
./bin/godnscli q myapp.lan
./bin/godnscli q api.myapp.lan

# 4. Check health
./bin/godnscli h
```

### Troubleshooting DNS Issues

When DNS isn't working as expected:

```bash
# 1. Check server health
./bin/godnscli h -v

# 2. Query with verbose output
./bin/godnscli q example.lan -v

# 3. Test different record types
./bin/godnscli q example.lan -t A
./bin/godnscli q example.lan -t AAAA
./bin/godnscli q example.lan -t NS

# 4. Try with increased timeout
./bin/godnscli q example.lan --timeout 30
```

### Testing Multiple Servers

Compare responses from different DNS servers:

```bash
# Query your local server
./bin/godnscli q example.lan -s localhost:53

# Query staging server
./bin/godnscli q example.lan -s staging-dns.example.com:53

# Query production server
./bin/godnscli q example.lan -s prod-dns.example.com:53

# Compare with Google DNS
./bin/godnscli q google.com -s 8.8.8.8:53
```

### CI/CD Integration

Use in your CI/CD pipeline:

```bash
#!/bin/bash
# test-dns.sh

# Wait for DNS server to be ready
timeout 60 bash -c 'until ./bin/godnscli h; do sleep 2; done'

# Run test suite
./bin/godnscli t

if [ $? -eq 0 ]; then
    echo "✓ DNS tests passed"
    exit 0
else
    echo "✗ DNS tests failed"
    exit 1
fi
```

### Kubernetes Health Checks

Use in Kubernetes liveness/readiness probes:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: godns
spec:
  containers:
    - name: godns
      image: godns:latest
      livenessProbe:
        exec:
          command:
            - /usr/local/bin/godnscli
            - h
            - --liveness-port=8080
        initialDelaySeconds: 10
        periodSeconds: 30
      readinessProbe:
        exec:
          command:
            - /usr/local/bin/godnscli
            - h
            - --readiness-port=8081
        initialDelaySeconds: 5
        periodSeconds: 10
```

### Monitoring Scripts

Create monitoring scripts:

```bash
#!/bin/bash
# monitor-dns.sh

SERVER="dns.example.com:53"
DOMAINS=("app.lan" "api.lan" "web.lan")

for domain in "${DOMAINS[@]}"; do
    echo "Testing $domain..."
    if ./bin/godnscli q "$domain" -s "$SERVER" > /dev/null 2>&1; then
        echo "  ✓ $domain is resolving"
    else
        echo "  ✗ $domain failed to resolve"
        # Send alert here
    fi
done
```

## Advanced Usage

### Custom DNS Server Testing

```bash
# Test against custom port
./bin/godnscli q example.lan -s localhost:5353

# Test against IPv6 address
./bin/godnscli q example.lan -s [::1]:53

# Test with specific network interface
./bin/godnscli q example.lan -s 192.168.1.100:53
```

### Batch Queries

Create a script to query multiple domains:

```bash
#!/bin/bash
# batch-query.sh

DOMAINS_FILE="domains.txt"
SERVER="localhost:53"

while IFS= read -r domain; do
    echo "Querying $domain..."
    ./bin/godnscli q "$domain" -s "$SERVER"
    echo ""
done < "$DOMAINS_FILE"
```

### Performance Testing

Test query performance:

```bash
#!/bin/bash
# perf-test.sh

ITERATIONS=100
DOMAIN="example.lan"

echo "Running $ITERATIONS queries..."
for i in $(seq 1 $ITERATIONS); do
    ./bin/godnscli q "$DOMAIN" > /dev/null 2>&1
done

echo "Completed $ITERATIONS queries"
```

### Comparing Record Types

Check all record types for a domain:

```bash
#!/bin/bash
# check-all-records.sh

DOMAIN=$1
TYPES=("A" "AAAA" "MX" "NS" "TXT" "CNAME" "SOA")

echo "Checking all record types for $DOMAIN"
for type in "${TYPES[@]}"; do
    echo ""
    echo "=== $type Records ==="
    ./bin/godnscli q "$DOMAIN" -t "$type"
done
```

## Troubleshooting

### Command Not Found

```bash
# Make sure you've built the CLI
make build-cli

# Check if binary exists
ls -la ./bin/godnscli

# Make it executable if needed
chmod +x ./bin/godnscli
```

### Connection Refused

```bash
# Check if DNS server is running
./bin/godnscli h

# Verify the server address
./bin/godnscli q example.lan -s localhost:53 -v

# Check if port is open
nc -zv localhost 53
```

### Timeout Errors

```bash
# Increase timeout
./bin/godnscli q example.lan --timeout 30

# Check network connectivity
ping localhost

# Test with verbose output
./bin/godnscli q example.lan -v
```

### No Answer Section

If you get a response but no answer section:

```bash
# Check if domain exists in your DNS server
./bin/godnscli q example.lan -v

# Try different record types
./bin/godnscli q example.lan -t A
./bin/godnscli q example.lan -t AAAA
./bin/godnscli q example.lan -t NS

# Check server logs for errors
docker-compose logs godns
```

### Health Check Fails

```bash
# Check if health endpoints are running
curl http://localhost:8080/livez
curl http://localhost:8081/readyz

# Verify correct ports
./bin/godnscli h --liveness-port 8080 --readiness-port 8081 -v

# Check server logs
docker-compose logs godns
```

## Tips and Best Practices

### 1. Use Aliases for Faster Commands

Always use the short aliases to save typing:

- `q` instead of `query`
- `h` instead of `health`
- `t` instead of `test`
- `v` instead of `version`

### 2. Add to PATH

Make `godnscli` available everywhere:

```bash
# Add to your .zshrc or .bashrc
export PATH="$PATH:/path/to/godns/bin"

# Then use from anywhere
godnscli q example.lan
```

### 3. Create Shell Aliases

Add to your shell config:

```bash
# In ~/.zshrc or ~/.bashrc
alias dnsq='./bin/godnscli q'
alias dnst='./bin/godnscli t'
alias dnsh='./bin/godnscli h'

# Usage
dnsq example.lan
dnst
dnsh
```

### 4. Use Verbose Mode for Debugging

When debugging, always use `-v` flag to see full details:

```bash
./bin/godnscli q example.lan -v
```

### 5. Save Common Queries as Scripts

Create helper scripts for common operations:

```bash
#!/bin/bash
# check-dns.sh
./bin/godnscli t && echo "✓ All tests passed"
```

## Getting Help

```bash
# General help
./bin/godnscli --help

# Command-specific help
./bin/godnscli query --help
./bin/godnscli health --help
./bin/godnscli test --help
```

## Next Steps

- See [README.md](../cmd/godnscli/README.md) for development information
- See [VALKEY_AUTH.md](./VALKEY_AUTH.md) for authentication setup
- Check the [main README](../README.md) for GoDNS server configuration

## Support

For issues or questions:

1. Check this guide first
2. Review the troubleshooting section
3. Check server logs: `docker-compose logs godns`
4. Open an issue on GitHub
