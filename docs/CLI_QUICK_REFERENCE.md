# GoDNS CLI Quick Reference

Fast reference for `godnscli` commands. For detailed guide, see [CLI_GUIDE.md](./CLI_GUIDE.md).

## Installation

```bash
make build-cli
./bin/godnscli v  # Verify installation
```

## Command Syntax

```bash
./bin/godnscli [command] [arguments] [flags]
```

## Quick Commands

| What You Want             | Command                        |
| ------------------------- | ------------------------------ |
| Discover server & domains | `./bin/godnscli d`             |
| Test everything           | `./bin/godnscli t`             |
| Check server health       | `./bin/godnscli h`             |
| Query a domain            | `./bin/godnscli q example.lan` |
| Show version              | `./bin/godnscli v`             |

## Global Flags

```bash
-s, --server    DNS server address (default: localhost:53)
-v, --verbose   Show detailed output
-h, --help      Show help
```

## Discover Command (`d`, `find`)

### Basic Usage

```bash
./bin/godnscli d [flags]
```

### Examples

```bash
# Basic discovery
./bin/godnscli d

# Verbose output with network interfaces
./bin/godnscli d -v

# Discover remote server
./bin/godnscli d -s 192.168.1.100:53
```

### What It Shows

- Server information and IP addresses
- Local network interfaces (with `-v`)
- How to find available DNS zones
- Network discovery commands
- Docker-specific tips (with `-v`)
- Example queries to try

## Query Command (`q`)

### Basic Usage

```bash
./bin/godnscli q <domain> [flags]
```

### Common Examples

```bash
# A record (IPv4)
./bin/godnscli q example.lan

# AAAA record (IPv6)
./bin/godnscli q example.lan -t AAAA

# MX records (mail)
./bin/godnscli q example.lan -t MX

# NS records (name servers)
./bin/godnscli q example.lan -t NS

# TXT records
./bin/godnscli q example.lan -t TXT

# Query different server
./bin/godnscli q example.lan -s 192.168.1.100:53

# Verbose output
./bin/godnscli q example.lan -v

# Custom timeout
./bin/godnscli q example.lan --timeout 10
```

### Query Flags

```bash
-t, --type      Record type (A, AAAA, MX, NS, TXT, CNAME, SOA, PTR, SRV)
--timeout       Timeout in seconds (default: 5)
```

## Health Command (`h`)

### Basic Usage

```bash
./bin/godnscli h [flags]
```

### Examples

```bash
# Check default ports
./bin/godnscli h

# Custom ports
./bin/godnscli h --liveness-port 9090 --readiness-port 9091

# Verbose output
./bin/godnscli h -v
```

### Health Flags

```bash
--liveness-port     Liveness probe port (default: 8080)
--readiness-port    Readiness probe port (default: 8081)
```

## Test Command (`t`)

### Basic Usage

```bash
./bin/godnscli t [flags]
```

### Examples

```bash
# Run all tests
./bin/godnscli t

# Verbose output
./bin/godnscli t -v

# Test different server
./bin/godnscli t -s 192.168.1.100:53
```

### What It Tests

- ✓ A records (IPv4)
- ✓ AAAA records (IPv6)
- ✓ MX records
- ✓ NS records
- ✓ External resolution

## Version Command (`v`)

```bash
./bin/godnscli v
```

## DNS Record Types

| Type    | Description        | Example Usage |
| ------- | ------------------ | ------------- |
| `A`     | IPv4 address       | `-t A`        |
| `AAAA`  | IPv6 address       | `-t AAAA`     |
| `MX`    | Mail exchange      | `-t MX`       |
| `NS`    | Name server        | `-t NS`       |
| `TXT`   | Text record        | `-t TXT`      |
| `CNAME` | Canonical name     | `-t CNAME`    |
| `SOA`   | Start of authority | `-t SOA`      |
| `PTR`   | Pointer record     | `-t PTR`      |
| `SRV`   | Service record     | `-t SRV`      |

## Common Workflows

### Finding What to Query

```bash
# 1. Discover server and get tips
./bin/godnscli d -v

# 2. List zones in Valkey
docker-compose exec valkey valkey-cli -a mysecretpassword KEYS 'dns:zone:*'

# 3. View a specific zone
docker-compose exec valkey valkey-cli -a mysecretpassword GET 'dns:zone:test.lan.'
```

### Adding a Test Zone

```bash
# Use the helper script
./hack/add-test-zone.sh test.lan 192.168.1.100

# Or manually via Valkey
docker-compose exec valkey valkey-cli -a mysecretpassword
AUTH default mysecretpassword
SET dns:zone:test.lan. '{"domain":"test.lan.","records":[...]}'
```

### Quick Test Workflow

```bash
# 1. Start server
docker-compose up -d

# 2. Add test zone
./hack/add-test-zone.sh myapp.lan 192.168.1.50

# 3. Test
./bin/godnscli t

# 4. Query
./bin/godnscli q myapp.lan
```

### Troubleshooting Workflow

```bash
# 1. Check health
./bin/godnscli h -v

# 2. Verbose query
./bin/godnscli q example.lan -v

# 3. Try different record types
./bin/godnscli q example.lan -t A
./bin/godnscli q example.lan -t AAAA
```

### Multi-Server Testing

```bash
# Local
./bin/godnscli q example.lan -s localhost:53

# Staging
./bin/godnscli q example.lan -s staging-dns:53

# Production
./bin/godnscli q example.lan -s prod-dns:53
```

## Exit Codes

| Code | Meaning |
| ---- | ------- |
| `0`  | Success |
| `1`  | Error   |

## Useful Shell Aliases

Add to `~/.zshrc` or `~/.bashrc`:

```bash
# GoDNS CLI shortcuts
alias dnsd='./bin/godnscli d'
alias dnsq='./bin/godnscli q'
alias dnst='./bin/godnscli t'
alias dnsh='./bin/godnscli h'
alias dnsv='./bin/godnscli v'
```

## Common Errors

### Connection Refused

```bash
# Check if server is running
./bin/godnscli h

# Verify address
./bin/godnscli q example.lan -s localhost:53 -v
```

### Timeout

```bash
# Increase timeout
./bin/godnscli q example.lan --timeout 30
```

### No Answer

```bash
# Try verbose mode
./bin/godnscli q example.lan -v

# Check different record types
./bin/godnscli q example.lan -t NS
```

## Examples for Copy-Paste

```bash
# Build CLI
make build-cli

# Test local server
./bin/godnscli t

# Query A record
./bin/godnscli q example.lan

# Query IPv6
./bin/godnscli q example.lan -t AAAA

# Query mail servers
./bin/godnscli q example.lan -t MX

# Check health
./bin/godnscli h

# Query with verbose output
./bin/godnscli q example.lan -v

# Query external DNS
./bin/godnscli q google.com -s 8.8.8.8:53

# Test remote server
./bin/godnscli t -s 192.168.1.100:53

# Show version
./bin/godnscli v
```

## Getting Help

```bash
# General help
./bin/godnscli --help

# Command help
./bin/godnscli query --help
./bin/godnscli health --help
./bin/godnscli test --help
```

## More Information

- **Detailed Guide**: [CLI_GUIDE.md](./CLI_GUIDE.md)
- **Development**: [../cmd/godnscli/README.md](../cmd/godnscli/README.md)
- **Authentication**: [VALKEY_AUTH.md](./VALKEY_AUTH.md)
- **Main README**: [../README.md](../README.md)
