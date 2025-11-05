# Test Data Seeding

GoDNS automatically seeds test data when running in **development mode**.

## How It Works

When `DEVELOPMENT=true` is set in your environment, GoDNS will automatically create 5 test zones with sample DNS records on startup:

1. **home.lan** - 5 records (home network devices)
2. **dev.local** - 5 records (development services)
3. **k8s.local** - 5 records (Kubernetes cluster)
4. **docker.local** - 5 records (Docker services)
5. **example.lan** - 6 records (example/demo)

**Total**: 5 zones with 26 DNS records

## Enable Test Data Seeding

### Option 1: Environment Variable

Set `DEVELOPMENT=true` in your `.env` file:

```bash
# .env
DEVELOPMENT=true
```

### Option 2: Export Environment Variable

```bash
export DEVELOPMENT=true
```

### Option 3: Inline

```bash
DEVELOPMENT=true ./bin/godns
```

## Verify Test Data

Once GoDNS starts with `DEVELOPMENT=true`, you'll see this in the logs:

```
INFO  starting configuration seeding...
INFO  development mode detected - seeding test data...
INFO  seeding test zones...
INFO  successfully seeded 5 test zones with 26 DNS records
INFO  test data seeded successfully
INFO  configuration seeding completed successfully
```

### Check via API

```bash
# List all zones
curl http://localhost:14000/api/v1/zones | jq .

# Get specific zone
curl http://localhost:14000/api/v1/zones/home.lan | jq .
```

### Check via DNS

```bash
# Query DNS records
dig @localhost router.home.lan
dig @localhost api.dev.local
dig @localhost master.k8s.local
```

## Test Data Details

### home.lan (Home Network)

```
router.home.lan.   → 192.168.1.1
nas.home.lan.      → 192.168.1.10
server.home.lan.   → 192.168.1.100
printer.home.lan.  → 192.168.1.50
pi.home.lan.       → 192.168.1.200
```

### dev.local (Development Services)

```
api.dev.local.     → 127.0.0.1
db.dev.local.      → 127.0.0.1
cache.dev.local.   → 127.0.0.1
web.dev.local.     → 127.0.0.1
mail.dev.local.    → 127.0.0.1
```

### k8s.local (Kubernetes Cluster)

```
master.k8s.local.  → 10.0.1.10
worker1.k8s.local. → 10.0.1.20
worker2.k8s.local. → 10.0.1.21
worker3.k8s.local. → 10.0.1.22
ingress.k8s.local. → 10.0.1.100
```

### docker.local (Docker Services)

```
portainer.docker.local.  → 172.17.0.10
registry.docker.local.   → 172.17.0.20
traefik.docker.local.    → 172.17.0.30
grafana.docker.local.    → 172.17.0.40
prometheus.docker.local. → 172.17.0.41
```

### example.lan (Example/Demo)

```
www.example.lan.  → 192.168.100.10
mail.example.lan. → 192.168.100.20
ftp.example.lan.  → 192.168.100.30
db.example.lan.   → 192.168.100.40
ns1.example.lan.  → 192.168.100.1
ns2.example.lan.  → 192.168.100.2
```

## Disable Test Data Seeding

### Production Mode

In production, **never** set `DEVELOPMENT=true`. The seeding service will only create default configuration (allowed LANs, upstream DNS) but **no test zones**.

```bash
# .env (production)
# DEVELOPMENT=false  # or simply omit it
```

### Clear Test Data

If you want to start fresh:

```bash
# Stop services
docker-compose down

# Clear Valkey data
rm -rf hack/data/valkey/*

# Restart
docker-compose up -d
```

## Seeding Behavior

- **Idempotent**: Test data is only seeded if no zones exist
- **Safe**: Won't overwrite existing zones
- **Non-blocking**: If seeding fails, the application continues to start
- **Development only**: Never runs in production (when `DEVELOPMENT` is false or unset)

## Manual Seeding Script

For more control, you can also use the manual seeding script:

```bash
./hack/seed-test-data.sh
```

This script:

- Creates 6 zones (includes lab.internal)
- Provides detailed output
- Can be run anytime after GoDNS is running
- Useful for resetting test data

## Logs

Watch for these log messages:

✅ **Success**:

```
INFO  development mode detected - seeding test data...
INFO  successfully seeded 5 test zones with 26 DNS records
INFO  test data seeded successfully
```

ℹ️ **Already Seeded**:

```
INFO  skipping test data seeding - 5 zones already exist
```

⚠️ **Warning** (non-fatal):

```
WARN  failed to seed test data: <error message>
```

## Use Cases

Test data seeding is useful for:

- 🧪 **Development**: Immediate data to test with
- 📚 **Learning**: Example DNS configurations
- 🔍 **Testing**: Consistent test data across environments
- 🎯 **Demos**: Ready-to-use examples for presentations
- 🚀 **Quick Start**: No manual setup required

## Notes

- Test data uses common private IP ranges (RFC 1918)
- All records have sensible TTL values (60-3600 seconds)
- Zone names use `.lan`, `.local`, and `.internal` TLDs (safe for testing)
- Data is stored in Valkey just like production data
- Can be exported using the export API endpoints
