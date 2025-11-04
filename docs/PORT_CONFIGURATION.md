# Port Configuration

This document describes all the ports used by GoDNS and its related services.

## Port Mapping

All ports have been configured to use the **14000+ range** to avoid conflicts with commonly used ports (especially 80xx range).

| Service                 | Port  | Environment Variable              | Description                       |
| ----------------------- | ----- | --------------------------------- | --------------------------------- |
| **DNS Server**          | 53    | `DNS_SERVER_PORT`                 | DNS query port (UDP/TCP)          |
| **Valkey (Redis)**      | 14379 | `VALKEY_PORT`                     | Key-value store                   |
| **DNS Liveness Probe**  | 14080 | `DNS_SERVER_LIVENESS_PROBE_PORT`  | Kubernetes liveness checks        |
| **DNS Readiness Probe** | 14081 | `DNS_SERVER_READYNESS_PROBE_PORT` | Kubernetes readiness checks       |
| **GoDNS HTTP API**      | 14082 | `HTTP_API_PORT`                   | REST API and Swagger UI           |
| **Keycloak**            | 14083 | `KEYCLOAK_PORT`                   | Authentication server             |
| **PostgreSQL**          | 5432  | -                                 | Keycloak database (internal only) |

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
redis-cli -h localhost -p 14379
```

### DNS Health Probes

```bash
# Liveness probe
curl http://localhost:14080/healthz

# Readiness probe
curl http://localhost:14081/readyz
```

### GoDNS HTTP API

```bash
# API endpoint
curl http://localhost:14082/api/v1/zones

# Swagger UI
open http://localhost:14082/swagger/index.html
```

### Keycloak Admin Console

```bash
# Open admin console
open http://localhost:14083

# Login with default credentials:
# Username: admin
# Password: admin
```

### PostgreSQL (Internal)

```bash
# Access from host (only via Docker)
docker exec -it postgres psql -U keycloak

# Not exposed to host by default
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
# Example: Change API port to 15000
HTTP_API_PORT=:15000

# Example: Change Keycloak port to 14090
KEYCLOAK_PORT=14090
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
      - "15379:6379" # Map host:15379 to container:6379

  keycloak:
    ports:
      - "15083:8080" # Map host:15083 to container:8080
```

## Port Conflicts

If you encounter port conflicts, check what's using the port:

```bash
# macOS/Linux
lsof -i :14082

# Or using netstat
netstat -an | grep 14082
```

### Common Conflicts

If you run into port conflicts, here are common culprits:

- **Port 53**: May require sudo on Unix systems
- **Port 14379**: Check for other Redis/Valkey instances
- **Ports 14080-14083**: Check for other development services

### Resolution

1. **Stop the conflicting service**:

   ```bash
   # Find the process
   lsof -i :14082

   # Kill it (replace PID)
   kill <PID>
   ```

2. **Change GoDNS port**:

   ```bash
   # Update .env file
   HTTP_API_PORT=:14092
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
        proxy_pass http://localhost:14082;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}

server {
    listen 80;
    server_name auth.example.com;

    location / {
        proxy_pass http://localhost:14083;
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
      targetPort: 14082 # Container port
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
lsof -i :14082

# Check Docker containers
docker ps

# Check Docker Compose logs
docker-compose logs godns
```

### Cannot Connect to Service

```bash
# Verify service is running
curl http://localhost:14082/health

# Check listening ports
netstat -tuln | grep 14082

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

# Access services
curl http://localhost:14082/api/v1/zones        # API
open http://localhost:14082/swagger/index.html  # Swagger
open http://localhost:14083                     # Keycloak
redis-cli -p 14379                              # Valkey
dig @localhost example.lan                       # DNS
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
