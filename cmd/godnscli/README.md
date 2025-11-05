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
./bin/godnscli export       # Export DNS zones to different formats
./bin/godnscli version      # Show version
```

## Export Command

Export DNS zones to different DNS provider formats:

```bash
# Export all zones in BIND format
./bin/godnscli export --api-url http://localhost:14000

# Export all zones in CoreDNS format
./bin/godnscli export --format coredns --api-url http://localhost:14000

# Export specific zone in PowerDNS format
./bin/godnscli export example.lan --format powerdns --api-url http://localhost:14000

# Export to file
./bin/godnscli export example.lan --format bind --output example.lan.zone
```

**Supported Formats:**

- `bind` - Standard BIND zone file format (default)
- `coredns` - CoreDNS configuration format
- `powerdns` - PowerDNS JSON API format
- `zonefile` - Generic zone file (same as bind)

For more details, see [Export API Documentation](../../docs/EXPORT_API.md).

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
