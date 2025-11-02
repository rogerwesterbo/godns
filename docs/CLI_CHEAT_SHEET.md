# GoDNS CLI Cheat Sheet

## Most Common Commands

```bash
# Discover server and network info
./bin/godnscli d

# Test everything
./bin/godnscli t

# Query a domain
./bin/godnscli q example.lan

# Check health
./bin/godnscli h

# Show version
./bin/godnscli v
```

## Query Different Record Types

```bash
./bin/godnscli q example.lan           # A record (IPv4)
./bin/godnscli q example.lan -t AAAA   # IPv6
./bin/godnscli q example.lan -t MX     # Mail servers
./bin/godnscli q example.lan -t NS     # Name servers
./bin/godnscli q example.lan -t TXT    # Text records
```

## Query Different Servers

```bash
./bin/godnscli q example.lan -s localhost:53              # Local
./bin/godnscli q example.lan -s 192.168.1.100:53         # Custom
./bin/godnscli q google.com -s 8.8.8.8:53                # Google DNS
```

## Verbose Output

```bash
./bin/godnscli q example.lan -v        # Detailed query output
./bin/godnscli h -v                    # Detailed health check
./bin/godnscli t -v                    # Detailed test output
```

## Building

```bash
make build-cli     # Build CLI tool
make build-all     # Build everything
make lint          # Run linters
```

## Troubleshooting

```bash
# Server not responding?
./bin/godnscli h -v

# Query not working?
./bin/godnscli q example.lan -v --timeout 30

# Check different record types
./bin/godnscli q example.lan -t A
./bin/godnscli q example.lan -t AAAA
./bin/godnscli q example.lan -t NS
```

## All Available Commands

| Command    | Alias      | What It Does              |
| ---------- | ---------- | ------------------------- |
| `discover` | `d`,`find` | Discover server & domains |
| `query`    | `q`        | Query DNS records         |
| `health`   | `h`        | Check server health       |
| `test`     | `t`        | Run test suite            |
| `version`  | `v`        | Show version              |

## Record Types

| Type    | What It's For   |
| ------- | --------------- |
| `A`     | IPv4 addresses  |
| `AAAA`  | IPv6 addresses  |
| `MX`    | Mail servers    |
| `NS`    | Name servers    |
| `TXT`   | Text records    |
| `CNAME` | Aliases         |
| `SOA`   | Zone info       |
| `PTR`   | Reverse DNS     |
| `SRV`   | Service records |

## Quick Setup

```bash
# 1. Build
make build-cli

# 2. Start server
docker-compose up -d

# 3. Add a test zone
./hack/add-test-zone.sh test.lan 192.168.1.100

# 4. Test
./bin/godnscli t

# 5. Discover what's available
./bin/godnscli d

# 6. Query
./bin/godnscli q test.lan
```

## More Help

- Detailed guide: `docs/CLI_GUIDE.md`
- Quick reference: `docs/CLI_QUICK_REFERENCE.md`
- Command help: `./bin/godnscli --help`
- Command-specific: `./bin/godnscli query --help`
