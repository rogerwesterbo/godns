# GoDNS Quick Start Guide

Get up and running with GoDNS in 5 minutes.

## Prerequisites

- Docker & Docker Compose
- Make
- Go 1.25.3+ (for building)

## Step 1: Clone and Build

```bash
# Navigate to the godns directory
cd godns

# Build everything
make build-all
```

**Expected output:**

```
Building godns...
godns built at /path/to/godns/bin/godns
Building godnscli...
godnscli built at /path/to/godns/bin/godnscli
```

## Step 2: Start the Server

```bash
# Start GoDNS and Valkey with Docker Compose
docker-compose up -d
```

**Verify it's running:**

```bash
# Check containers
docker-compose ps

# Check logs
docker-compose logs godns
```

## Step 3: Test the Server

```bash
# Run the test suite
./bin/godnscli test
```

**Expected output:**

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

## Step 4: Query Domains

```bash
# Query a domain
./bin/godnscli query example.lan

# Or use the short alias
./bin/godnscli q example.lan
```

## Step 5: Check Health

```bash
# Check server health
./bin/godnscli health

# Or use the short alias
./bin/godnscli h
```

## What's Next?

### Learn the CLI

- **[Cheat Sheet](CLI_CHEAT_SHEET.md)** - Quick command reference
- **[Quick Reference](CLI_QUICK_REFERENCE.md)** - Fast lookup guide
- **[Complete Guide](CLI_GUIDE.md)** - Detailed documentation

### Configure the Server

1. Copy the example environment file:

   ```bash
   cp .env.example .env
   ```

2. Edit `.env` with your settings:

   ```bash
   VALKEY_ADDR=localhost:14379
   VALKEY_USERNAME=default
   VALKEY_PASSWORD=mysecretpassword
   DNS_PORT=53
   ```

3. Restart the server:
   ```bash
   docker-compose restart
   ```

### Add DNS Records

The server uses Valkey (Redis) for dynamic configuration. Records are stored with keys:

- Allowed LANs: `dns:config:allowedlans`
- Upstream DNS: `dns:config:upstream`

See the [Valkey Authentication Guide](VALKEY_AUTH.md) for more details.

## Common Commands

```bash
# Build
make build-all              # Build everything
make build-cli              # Build CLI only
make lint                   # Run linters

# Docker
docker-compose up -d        # Start server
docker-compose down         # Stop server
docker-compose logs godns   # View logs
docker-compose restart      # Restart server

# CLI
./bin/godnscli t            # Run tests
./bin/godnscli q <domain>   # Query domain
./bin/godnscli h            # Check health
./bin/godnscli v            # Show version
```

## Troubleshooting

### Server won't start

```bash
# Check logs
docker-compose logs godns

# Check Valkey
docker-compose logs valkey

# Restart everything
docker-compose down
docker-compose up -d
```

### CLI can't connect

```bash
# Check if server is running
docker-compose ps

# Check health
./bin/godnscli h -v

# Try with verbose output
./bin/godnscli q example.lan -v
```

### Tests failing

```bash
# Run tests with verbose output
./bin/godnscli t -v

# Check server logs
docker-compose logs godns

# Verify Valkey connection
docker-compose exec valkey valkey-cli ping
```

## Getting Help

```bash
# CLI help
./bin/godnscli --help

# Command-specific help
./bin/godnscli query --help
./bin/godnscli health --help
./bin/godnscli test --help
```

## Next Steps

1. **Read the Documentation**

   - [CLI Guide](CLI_GUIDE.md) - Complete CLI documentation
   - [Valkey Auth](VALKEY_AUTH.md) - Authentication setup

2. **Customize Configuration**

   - Edit `.env` file
   - Configure Valkey ACLs
   - Add your DNS zones

3. **Deploy to Production**
   - Set up Kubernetes deployment
   - Configure health checks
   - Set up monitoring

## Quick Reference

| Task         | Command                     |
| ------------ | --------------------------- |
| Build all    | `make build-all`            |
| Start server | `docker-compose up -d`      |
| Stop server  | `docker-compose down`       |
| Run tests    | `./bin/godnscli t`          |
| Query domain | `./bin/godnscli q <domain>` |
| Check health | `./bin/godnscli h`          |
| View logs    | `docker-compose logs godns` |

## Support

- Check [CLI Guide](CLI_GUIDE.md) for detailed help
- Review [Troubleshooting](CLI_GUIDE.md#troubleshooting) section
- Check server logs: `docker-compose logs godns`

---

**Congratulations!** 🎉 You now have GoDNS running and ready to use.
