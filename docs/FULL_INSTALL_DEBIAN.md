# Full Production Install Guide - Debian Trixie

Complete installation guide for GoDNS with Keycloak authentication, PostgreSQL, and Valkey on Debian Trixie.

## Architecture Overview

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│    Keycloak     │────▶│   PostgreSQL    │     │     Valkey      │
│  (Auth Server)  │     │  (Keycloak DB)  │     │  (DNS Storage)  │
└────────┬────────┘     └─────────────────┘     └────────┬────────┘
         │                                               │
         │              ┌─────────────────┐              │
         └─────────────▶│     GoDNS       │◀─────────────┘
                        │   (DNS + API)   │
                        └─────────────────┘
```

## Prerequisites

- Debian Trixie (13) or compatible
- Root/sudo access
- Domain name (optional, for production)
- At least 2GB RAM, 10GB disk space

## 1. System Preparation

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install required packages
sudo apt install -y curl wget gnupg2 \
    apt-transport-https ca-certificates lsb-release
```

## 2. Install PostgreSQL

```bash
# Install PostgreSQL
sudo apt install -y postgresql postgresql-contrib

# Start and enable PostgreSQL
sudo systemctl enable postgresql
sudo systemctl start postgresql

# Create Keycloak database and user
sudo -u postgres psql << EOF
CREATE USER keycloak WITH PASSWORD 'keycloak_secure_password';
CREATE DATABASE keycloak OWNER keycloak;
GRANT ALL PRIVILEGES ON DATABASE keycloak TO keycloak;
EOF

# Verify
sudo -u postgres psql -c "\l" | grep keycloak
```

## 3. Install Valkey

```bash
# Install Valkey
sudo apt install -y valkey

# Configure Valkey with password authentication
sudo nano /etc/valkey/valkey.conf
```

Add/modify these settings:

```conf
# Bind to localhost only (or your internal IP)
bind 127.0.0.1

# Set password
requirepass your_valkey_secure_password

# Enable persistence
appendonly yes
```

```bash
# Restart Valkey
sudo systemctl enable valkey
sudo systemctl restart valkey

# Test connection
valkey-cli -a your_valkey_secure_password ping
```

## 4. Install Java (for Keycloak)

```bash
# Install OpenJDK 17
sudo apt install -y openjdk-21-jdk (keycloak must have openjdk-21 ...)

# Verify
java -version
```

## 5. Install Keycloak

```bash
# Download Keycloak
KEYCLOAK_VERSION="26.4.7"
cd /opt
sudo wget "https://github.com/keycloak/keycloak/releases/download/${KEYCLOAK_VERSION}/keycloak-${KEYCLOAK_VERSION}.tar.gz"
sudo tar -xzf "keycloak-${KEYCLOAK_VERSION}.tar.gz"
sudo mv "keycloak-${KEYCLOAK_VERSION}" keycloak
sudo rm "keycloak-${KEYCLOAK_VERSION}.tar.gz"

# Create keycloak user
sudo useradd -r -s /bin/false keycloak
sudo chown -R keycloak:keycloak /opt/keycloak
```

### Configure Keycloak

```bash
sudo nano /opt/keycloak/conf/keycloak.conf
```

Add:

```properties
# Database
db=postgres
db-url=jdbc:postgresql://localhost:5432/keycloak
db-username=keycloak
db-password=keycloak_secure_password

# HTTP
hostname=<ip or hostname>
http-enabled=true
http-port=8180
hostname-strict=false
proxy-headers=xforwarded

# Production settings (uncomment for production)
# https-certificate-file=/path/to/cert.pem
# https-certificate-key-file=/path/to/key.pem
# hostname=auth.yourdomain.com
```

### Create Keycloak Admin User

```bash
cd /opt/keycloak
sudo -u keycloak KEYCLOAK_ADMIN=admin KEYCLOAK_ADMIN_PASSWORD=admin_secure_password \
    ./bin/kc.sh build
```

### Create Systemd Service

```bash
sudo nano /etc/systemd/system/keycloak.service
```

Add:

```ini
[Unit]
Description=Keycloak Authentication Server
After=network.target postgresql.service

[Service]
Type=simple
User=keycloak
Group=keycloak
Environment="KEYCLOAK_ADMIN=admin"
Environment="KEYCLOAK_ADMIN_PASSWORD=admin_secure_password"
ExecStart=/opt/keycloak/bin/kc.sh start --optimized
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

```bash
# Build optimized Keycloak
cd /opt/keycloak
sudo -u keycloak ./bin/kc.sh build

# Start Keycloak
sudo systemctl daemon-reload
sudo systemctl enable keycloak
sudo systemctl start keycloak

# Check status (may take 30-60 seconds to start)
sudo systemctl status keycloak
```

## 6. Configure Keycloak for GoDNS

Access Keycloak admin console at `http://your-server:8180`

### Create Realm

1. Log in with admin credentials
2. Click "Create realm"
3. Name: `godns`
4. Click "Create"

### Create API Client

1. Go to Clients → Create client
2. Client ID: `godns-api`
3. Client type: `OpenID Connect`
4. Click "Next"
5. Client authentication: `ON`
6. Authorization: `OFF`
7. Click "Next"
8. Valid redirect URIs: `http://localhost:8080/*`
9. Web origins: `*`
10. Click "Save"

### Create CLI Client

1. Go to Clients → Create client
2. Client ID: `godns-cli`
3. Client type: `OpenID Connect`
4. Click "Next"
5. Client authentication: `ON`
6. Standard flow: `ON`
7. Direct access grants: `ON`
8. Click "Next"
9. Valid redirect URIs: `http://localhost:*`
10. Click "Save"

### Create Web Client

This client is for the GoDNS Web Application frontend (SPA). It uses PKCE for secure authentication without exposing client secrets.

1. Go to Clients → Create client
2. Client ID: `godns-web`
3. Client type: `OpenID Connect`
4. Click "Next"
5. Client authentication: `OFF` (public client)
6. Standard flow: `ON`
7. Direct access grants: `OFF`
8. Click "Next"
9. Root URL: `http://localhost:14200`
10. Home URL: `http://localhost:14200`
11. Valid redirect URIs: `http://localhost:14200/*`
12. Valid post logout redirect URIs: `http://localhost:14200`
13. Web origins: `http://localhost:14200`
14. Click "Save"

After saving, configure PKCE:

1. Go to Clients → `godns-web` → Settings
2. Scroll to "Advanced" section (or click Advanced tab)
3. Under "Proof Key for Code Exchange Code Challenge Method", select `S256`
4. Click "Save"

> **Note**: For production, replace `http://localhost:14200` with your actual web application URL.

### Create User

1. Go to Users → Add user
2. Username: `dnsadmin`
3. Email: `dnsadmin@example.com`
4. First name: `DNS`
5. Last name: `Admin`
6. Email verified: `ON`
7. Click "Create"
8. Go to Credentials tab
9. Set password, disable "Temporary"

### Get Client Secrets

1. Go to Clients → `godns-api` → Credentials
2. Copy the Client secret
3. Go to Clients → `godns-cli` → Credentials
4. Copy the Client secret

## 7. Install GoDNS

```bash
# Download binary
VERSION="0.0.4"
curl -fL -O "https://github.com/rogerwesterbo/godns/releases/download/v${VERSION}/godns-${VERSION}-linux-amd64.tar.gz"

# Extract and install
tar -xzf "godns-${VERSION}-linux-amd64.tar.gz"
sudo mv "godns-${VERSION}-linux-amd64" /usr/local/bin/godns
sudo chmod +x /usr/local/bin/godns
```

### Configure GoDNS

```bash
sudo mkdir -p /etc/godns
sudo nano /etc/godns/godns.env
```

Add:

```bash
# Authentication (Keycloak)
AUTH_ENABLED=true
KEYCLOAK_URL=http://localhost:8180
KEYCLOAK_REALM=godns
KEYCLOAK_API_CLIENT_ID=godns-api
KEYCLOAK_CLI_CLIENT_ID=godns-cli

# Valkey connection
VALKEY_HOST=localhost
VALKEY_PORT=6379
VALKEY_USERNAME=default
VALKEY_TOKEN=your_valkey_secure_password

# DNS settings
DNS_SERVER_PORT=:53
DNS_UPSTREAM_SERVER=8.8.8.8:53

# HTTP API
HTTP_API_PORT=:8080

# Logging
LOG_LEVEL=info
LOG_JSON=true

# Production mode
DEVELOPMENT=false
```

### Create Systemd Service

```bash
sudo nano /etc/systemd/system/godns.service
```

Add:

```ini
[Unit]
Description=GoDNS DNS Server
After=network.target valkey.service keycloak.service
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

```bash
sudo systemctl daemon-reload
sudo systemctl enable godns
sudo systemctl start godns
```

## 8. Install GoDNS CLI

```bash
# Download and install
VERSION="0.0.4"
curl -fL -O "https://github.com/rogerwesterbo/godns/releases/download/v${VERSION}/godnscli-${VERSION}-linux-amd64.tar.gz"
tar -xzf "godnscli-${VERSION}-linux-amd64.tar.gz"
sudo mv "godnscli-${VERSION}-linux-amd64" /usr/local/bin/godnscli
sudo chmod +x /usr/local/bin/godnscli
```

### Configure CLI

```bash
mkdir -p ~/.config/godnscli
nano ~/.config/godnscli/config.yaml
```

Add:

```yaml
server: http://localhost:8080
keycloak:
  url: http://localhost:8180
  realm: godns
  client_id: godns-cli
  client_secret: YOUR_CLI_CLIENT_SECRET
```

### Login and Test

```bash
# Login (will prompt for username/password)
godnscli login

# Test
godnscli h
godnscli test
```

## 9. Firewall Configuration

```bash
# Install UFW if not present
sudo apt install -y ufw

# Allow SSH
sudo ufw allow ssh

# Allow DNS
sudo ufw allow 53/tcp
sudo ufw allow 53/udp

# Allow HTTP API (internal only, or restrict to specific IPs)
sudo ufw allow from 192.168.0.0/16 to any port 8080

# Allow Keycloak (internal only)
sudo ufw allow from 192.168.0.0/16 to any port 8180

# Enable firewall
sudo ufw enable
```

## 10. TLS/HTTPS (Production)

### Using Let's Encrypt

```bash
# Install certbot
sudo apt install -y certbot

# Get certificate (adjust domain)
sudo certbot certonly --standalone -d dns.yourdomain.com

# Certificates will be at:
# /etc/letsencrypt/live/dns.yourdomain.com/fullchain.pem
# /etc/letsencrypt/live/dns.yourdomain.com/privkey.pem
```

### Configure Keycloak with TLS

Edit `/opt/keycloak/conf/keycloak.conf`:

```properties
https-certificate-file=/etc/letsencrypt/live/dns.yourdomain.com/fullchain.pem
https-certificate-key-file=/etc/letsencrypt/live/dns.yourdomain.com/privkey.pem
hostname=auth.yourdomain.com
http-enabled=false
```

## 11. Service Management

```bash
# Start all services
sudo systemctl start postgresql valkey keycloak godns

# Check status
sudo systemctl status postgresql valkey keycloak godns

# View logs
sudo journalctl -u godns -f
sudo journalctl -u keycloak -f

# Restart services
sudo systemctl restart godns
```

## 12. Backup Strategy

### Backup PostgreSQL

```bash
# Backup Keycloak database
sudo -u postgres pg_dump keycloak > keycloak_backup.sql
```

### Backup Valkey

```bash
# Trigger save
valkey-cli -a your_valkey_secure_password BGSAVE

# Copy dump file
sudo cp /var/lib/valkey/dump.rdb /backup/valkey_backup.rdb
```

### Backup Configuration

```bash
# Backup configs
sudo tar -czf godns_config_backup.tar.gz \
    /etc/godns \
    /etc/valkey/valkey.conf \
    /opt/keycloak/conf
```

## Verification Checklist

- [ ] PostgreSQL running and accessible
- [ ] Valkey running with authentication
- [ ] Keycloak accessible at http://server:8180
- [ ] Keycloak realm `godns` created
- [ ] Keycloak clients configured
- [ ] Keycloak user created
- [ ] GoDNS service running
- [ ] DNS queries working: `dig @localhost google.com`
- [ ] API accessible: `curl http://localhost:8080/health`
- [ ] CLI authentication working: `godnscli login`

## Troubleshooting

### Keycloak won't start

```bash
# Check logs
sudo journalctl -u keycloak -n 100

# Verify PostgreSQL connection
sudo -u postgres psql -U keycloak -d keycloak -h localhost -c "SELECT 1"
```

### GoDNS authentication errors

```bash
# Verify Keycloak is reachable
curl http://localhost:8180/realms/godns/.well-known/openid_configuration

# Check GoDNS logs
sudo journalctl -u godns | grep -i auth
```

### DNS not resolving

```bash
# Check port 53
sudo ss -tlnp | grep :53

# Check GoDNS logs
sudo journalctl -u godns -f

# Test directly
dig @127.0.0.1 -p 53 google.com
```

## Next Steps

- Configure monitoring with Prometheus/Grafana
- Set up log aggregation
- Configure DNS-over-HTTPS (DoH)
- Set up high availability with multiple nodes
- Configure DNSSEC signing

## Related Documentation

- [Quick Install (Testing Only)](QUICK_INSTALL_DEBIAN.md)
- [API Documentation](API_DOCUMENTATION.md)
- [CLI Guide](CLI_GUIDE.md)
- [Keycloak Setup Details](KEYCLOAK_SETUP.md)
