# Port Configuration

This document describes all the ports used by GoDNS and its related services.

## Port Mapping

All ports have been configured to use the **14000+ range** to avoid conflicts with commonly used ports (especially 80xx range).

| Service                 | Port  | Environment Variable              | Description                   |
| ----------------------- | ----- | --------------------------------- | ----------------------------- |
| **DNS Server**          | 53    | `DNS_SERVER_PORT`                 | DNS query port (UDP/TCP)      |
| **HTTP API**            | 14000 | `HTTP_API_PORT`                   | REST API and Swagger UI       |
| **API Liveness Probe**  | 14001 | `HTTP_API_LIVENESS_PROBE_PORT`    | Kubernetes liveness checks    |
| **API Readiness Probe** | 14002 | `HTTP_API_READINESS_PROBE_PORT`   | Kubernetes readiness checks   |
| **DNS Liveness Probe**  | 14003 | `DNS_SERVER_LIVENESS_PROBE_PORT`  | DNS liveness checks           |
| **DNS Readiness Probe** | 14004 | `DNS_SERVER_READYNESS_PROBE_PORT` | DNS readiness checks          |
| **PostgreSQL**          | 14100 | `POSTGRES_PORT`                   | Keycloak database             |
| **Keycloak HTTP**       | 14101 | `KEYCLOAK_PORT_HTTP`              | OAuth2/OIDC server (HTTP)     |
| **Keycloak HTTPS**      | 14102 | `KEYCLOAK_PORT_HTTPS`             | OAuth2/OIDC server (HTTPS)    |
| **Valkey (Redis)**      | 14103 | `VALKEY_PORT`                     | DNS records storage           |
| **Frontend (Future)**   | 14200 | -                                 | Web UI (planned)              |

## Default Configuration

The default ports are defined in:

1. **Code defaults**: `internal/settings/settings.go`
2. **Environment template**: `.env.example`
3. **Docker Compose**: `docker-compose.yaml`

## Accessing Services

### DNS Server

```bash
# Query DNS
dig @localhost -p 53 example.lan

# Or using default port
dig @localhost example.lan
```

### Valkey (Redis)

```bash
# Connect to Valkey
redis-cli -h localhost -p 14103
```

### DNS Health Probes

```bash
# Liveness probe
curl http://localhost:14003/healthz

# Readiness probe
curl http://localhost:14004/readyz
```

### GoDNS HTTP API

```bash
# Get a token first
TOKEN=$(curl -s -X POST "http://localhost:14101/realms/godns/protocol/openid-connect/token" \
  -d "client_id=godns-cli" \
  -d "username=testuser" \
  -d "password=password" \
  -d "grant_type=password" | jq -r '.access_token')

# API endpoint (requires authentication)
curl -H "Authorization: Bearer $TOKEN" http://localhost:14000/api/v1/zones

# Swagger UI (includes OAuth2 login)
open http://localhost:14000/swagger/index.html
```

### CORS Configuration

The HTTP API supports Cross-Origin Resource Sharing (CORS) for web applications. Configure allowed origins in `.env`:

```bash
# Allow multiple origins (comma-separated)
HTTP_API_CORS_ALLOWED_ORIGINS=http://localhost:14000,http://localhost:14200

# This allows:
# - Swagger UI (http://localhost:14000)
# - Frontend web app (http://localhost:14200)
```

The API will respond with appropriate CORS headers for allowed origins:

- `Access-Control-Allow-Origin`
- `Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS`
- `Access-Control-Allow-Headers: Content-Type, Authorization, X-Requested-With`
- `Access-Control-Allow-Credentials: true`

### Keycloak Admin Console

```bash
# Open admin console (HTTP)
open http://localhost:14101

# Or HTTPS
open https://localhost:14102

# Login with default credentials:
# Username: admin
# Password: admin
```

### PostgreSQL

```bash
# Access from host
psql -h localhost -p 14100 -U keycloak

# Or via Docker
docker exec -it postgres psql -U keycloak
```

## Customizing Ports

### Using Environment Variables

Create a `.env` file in the project root:

```bash
# Copy the example
cp .env.example .env

# Edit the file
vim .env
```

Change any port values:

```bash
# Example: Change Keycloak HTTP port
KEYCLOAK_PORT_HTTP=14201

# Example: Change Valkey port
VALKEY_PORT=14203
```

### Using Command Line Flags

```bash
# Start GoDNS with custom ports
./bin/godns \
  --dns-port :5353 \
  --liveness-port :15080 \
  --readiness-port :15081

# Start API with custom port
./bin/godnsapi --port :15082
```

### Using Docker Compose

Edit `docker-compose.yaml` to change port mappings:

```yaml
services:
  valkey:
    ports:
      - "14203:6379" # Map host:14203 to container:6379

  keycloak:
    ports:
      - "14201:8080" # Map host:14201 to container:8080
```

## Port Conflicts

If you encounter port conflicts, check what's using the port:

```bash
# macOS/Linux
lsof -i :14000

# Or using netstat
netstat -an | grep 14000
```

### Common Conflicts

If you run into port conflicts, here are common culprits:

- **Port 53**: May require sudo on Unix systems
- **Port 14103**: Check for other Redis/Valkey instances  
- **Ports 14000-14004**: Check for other development services
- **Ports 14100-14103**: Check for other database/auth services

### Resolution

1. **Stop the conflicting service**:

   ```bash
   # Find the process
   lsof -i :14000

   # Kill it (replace PID)
   kill <PID>
   ```

2. **Change GoDNS port**:

   ```bash
   # Update .env file
   HTTP_API_PORT=:14050
   ```

3. **Use Docker network isolation**:
   Docker Compose services communicate internally without exposing ports to the host.

## Production Considerations

### Reverse Proxy

In production, use a reverse proxy (nginx, Traefik, Caddy) to:

- Map standard ports (80, 443) to GoDNS services
- Handle SSL/TLS termination
- Load balancing

Example nginx configuration:

```nginx
server {
    listen 80;
    server_name api.example.com;

    location / {
        proxy_pass http://localhost:14000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header Authorization $http_authorization;
    }
}

server {
    listen 80;
    server_name auth.example.com;

    location / {
        proxy_pass http://localhost:14101;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### Kubernetes

In Kubernetes, services use ClusterIP and aren't directly exposed:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: godns-api
spec:
  selector:
    app: godns
  ports:
    - name: http
      port: 80
      targetPort: 14000 # Container port
```

### Security

1. **Firewall Rules**: Only expose necessary ports
2. **Internal Services**: Keep PostgreSQL internal-only
3. **Authentication**: Always use authentication in production
4. **TLS/SSL**: Use HTTPS for all external services

## Troubleshooting

### Service Won't Start

```bash
# Check if port is in use
lsof -i :14000

# Check Docker containers
docker ps

# Check Docker Compose logs
docker-compose logs godnsapi
```

### Cannot Connect to Service

```bash
# Verify service is running
curl http://localhost:14000/health

# Check listening ports
netstat -tuln | grep 14000

# Verify Docker port mapping
docker port <container-name>
```

### Permission Denied (Port 53)

Port 53 requires elevated privileges on most systems:

```bash
# Run with sudo
sudo ./bin/godns

# Or use capability (Linux)
sudo setcap CAP_NET_BIND_SERVICE=+eip ./bin/godns

# Or bind to higher port
./bin/godns --dns-port :5353
```

## Quick Reference

```bash
# Start all services
docker-compose up -d

# Check all ports
docker-compose ps

# View logs
docker-compose logs -f

# Access services (with authentication)
godnscli login                                  # Login first
godnscli export --format bind                   # Use CLI (auto-auth)
open http://localhost:14000/swagger/index.html  # Swagger (OAuth2)
open http://localhost:14101                     # Keycloak admin
redis-cli -p 14103                              # Valkey
dig @localhost example.lan                      # DNS
```

## Summary

GoDNS uses the following port scheme:

- **DNS**: Standard port 53
- **All other services**: 14000+ range to avoid common port conflicts
- **Internal services**: Not exposed to host (PostgreSQL)

This configuration avoids conflicts with:

- Common web servers (80, 443, 8000, 8080, 8443)
- Development tools (3000, 5000, 9000)
- Database defaults (3306, 5432, 6379)
- Your specific range (80xx ports)

All ports are configurable via environment variables or command-line flags.
