# GoDNS CLI Tool

Command-line tool for testing and interacting with GoDNS servers.

## Quick Start

```bash
# Build
make build-cli

# Run
./bin/godnscli --help
```

## Common Commands

```bash
./bin/godnscli discover     # Find available domains and server info
./bin/godnscli test         # Test DNS functionality
./bin/godnscli query example.lan  # Query a domain
./bin/godnscli health       # Check server health
./bin/godnscli version      # Show version
```

## Full Documentation

For complete documentation, see:

- **[CLI Guide](../../docs/CLI_GUIDE.md)** - Complete command reference and usage examples
- **[Quick Start](../../docs/QUICK_START.md)** - Getting started with GoDNS
- **[Finding Domains](../../docs/FINDING_DOMAINS.md)** - Discover what domains to query

## Global Flags

- `-s, --server` - DNS server address (default: `localhost:53`)
- `-v, --verbose` - Enable verbose output
- `-h, --help` - Show help

## Learn More

Visit the [main documentation](../../docs/README.md) for comprehensive guides and references.
