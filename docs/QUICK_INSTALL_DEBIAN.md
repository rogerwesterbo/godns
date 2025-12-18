# Quick Install Guide - Debian Trixie (DNS Testing Only)

This guide covers a minimal installation for testing GoDNS DNS features on Debian Trixie.
**Authentication is disabled** - suitable for development and testing only.

## Prerequisites

- Debian Trixie (13) or compatible
- Root/sudo access
- Internet connection

## 1. Install Valkey

Valkey is a Redis-compatible in-memory data store.

```bash
# Install Valkey
sudo apt update
sudo apt install -y valkey

# Enable and start Valkey
sudo systemctl enable valkey
sudo systemctl start valkey

# Verify Valkey is running
valkey-cli ping
# Should return: PONG
```

### Configure Valkey (Optional - for password auth)

```bash
# Edit Valkey config
sudo nano /etc/valkey/valkey.conf

# Add or modify:
# requirepass your_secure_password

# Restart Valkey
sudo systemctl restart valkey

# Test with password
valkey-cli -a your_secure_password ping
```

## 2. Install GoDNS

### Option A: Download Pre-built Binary

```bash
# Get the latest version
VERSION="0.0.4"  # Check releases for latest

# Download the binary
curl -fL -O "https://github.com/rogerwesterbo/godns/releases/download/v${VERSION}/godns-${VERSION}-linux-amd64.tar.gz"

# Extract
tar -xzf "godns-${VERSION}-linux-amd64.tar.gz"

# Install
sudo mv "godns-${VERSION}-linux-amd64" /usr/local/bin/godns
sudo chmod +x /usr/local/bin/godns

# Verify
godns --help
```

### Option B: Build from Source

```bash
# Install Go (if not already installed)
sudo apt install -y golang-go

# Clone the repository
git clone https://github.com/rogerwesterbo/godns.git
cd godns

# Build
go build -o godns ./cmd/godns

# Install
sudo mv godns /usr/local/bin/
```

## 3. Install GoDNS CLI

```bash
# Download
VERSION="0.0.4"
curl -fL -O "https://github.com/rogerwesterbo/godns/releases/download/v${VERSION}/godnscli-${VERSION}-linux-amd64.tar.gz"

# Extract and install
tar -xzf "godnscli-${VERSION}-linux-amd64.tar.gz"
sudo mv "godnscli-${VERSION}-linux-amd64" /usr/local/bin/godnscli
sudo chmod +x /usr/local/bin/godnscli

# Verify
godnscli --help
```

## 4. Configure GoDNS

Create an environment file:

```bash
sudo mkdir -p /etc/godns
sudo nano /etc/godns/godns.env
```

Add the following configuration:

```bash
# Disable authentication for testing
AUTH_ENABLED=false

# Valkey connection
VALKEY_HOST=localhost
VALKEY_PORT=6379
VALKEY_USERNAME=
VALKEY_TOKEN=

# DNS settings
DNS_SERVER_PORT=:53
DNS_UPSTREAM_SERVER=8.8.8.8:53

# HTTP API
HTTP_API_PORT=:8080

# Logging
LOG_LEVEL=info
LOG_JSON=false

# Development mode (seeds test data)
DEVELOPMENT=true
```

## 5. Create Systemd Service

```bash
sudo nano /etc/systemd/system/godns.service
```

Add:

```ini
[Unit]
Description=GoDNS DNS Server
After=network.target valkey.service
Wants=valkey.service

[Service]
Type=simple
EnvironmentFile=/etc/godns/godns.env
ExecStart=/usr/local/bin/godns
Restart=always
RestartSec=5

# Security hardening
NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true
PrivateTmp=true

# DNS needs to bind to port 53
AmbientCapabilities=CAP_NET_BIND_SERVICE
CapabilityBoundingSet=CAP_NET_BIND_SERVICE

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl daemon-reload
sudo systemctl enable godns
sudo systemctl start godns

# Check status
sudo systemctl status godns

# View logs
sudo journalctl -u godns -f
```

## 6. Test the Installation

### Test DNS Server

```bash
# Query the DNS server
dig @localhost example.lan A

# Test upstream forwarding
dig @localhost google.com A
```

### Test API

```bash
# Health check
curl http://localhost:8080/health

# List zones (should be empty or have test data if DEVELOPMENT=true)
curl http://localhost:8080/api/v1/zones | jq .
```

### Test CLI

```bash
# Run built-in tests
godnscli test

# Query a record
godnscli q google.com

# Health check
godnscli h
```

## 7. Add DNS Records

### Using the API

```bash
# Create a zone with records
curl -X POST http://localhost:8080/api/v1/zones \
  -H "Content-Type: application/json" \
  -d '{
    "domain": "home.lan",
    "records": [
      {
        "name": "router.home.lan.",
        "type": "A",
        "ttl": 300,
        "value": "192.168.1.1"
      },
      {
        "name": "nas.home.lan.",
        "type": "A",
        "ttl": 300,
        "value": "192.168.1.100"
      }
    ]
  }' | jq .

# Test the new record
dig @localhost router.home.lan A
```

### Using the CLI

```bash
# Create a zone
godnscli zone create myzone.lan

# Add a record
godnscli record add myzone.lan www A 192.168.1.50 --ttl 300

# Query it
dig @localhost www.myzone.lan A
```

## 8. Configure System to Use GoDNS

### Option A: Point /etc/resolv.conf to GoDNS

```bash
# Backup original
sudo cp /etc/resolv.conf /etc/resolv.conf.backup

# Point to local DNS
echo "nameserver 127.0.0.1" | sudo tee /etc/resolv.conf
```

### Option B: Configure NetworkManager

```bash
# Edit connection
sudo nmcli con mod "Your Connection" ipv4.dns "127.0.0.1"
sudo nmcli con mod "Your Connection" ipv4.ignore-auto-dns yes
sudo nmcli con down "Your Connection" && sudo nmcli con up "Your Connection"
```

## Troubleshooting

### Port 53 already in use

```bash
# Check what's using port 53
sudo ss -tlnp | grep :53

# If systemd-resolved is using it:
sudo systemctl stop systemd-resolved
sudo systemctl disable systemd-resolved
```

### Valkey connection errors

```bash
# Check Valkey is running
sudo systemctl status valkey

# Test connection
valkey-cli ping
```

### View GoDNS logs

```bash
sudo journalctl -u godns -f
```

## Next Steps

- For production deployment with authentication, see [FULL_INSTALL_DEBIAN.md](FULL_INSTALL_DEBIAN.md)
- For testing guide, see [TESTING_GUIDE.md](TESTING_GUIDE.md)
- For API documentation, see [API_DOCUMENTATION.md](API_DOCUMENTATION.md)
