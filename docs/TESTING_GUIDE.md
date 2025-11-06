# GoDNS Testing Guide

Complete guide for testing your GoDNS server.

## Quick Start

```bash
# Build everything
make build-all

# Start services
docker-compose up -d

# Run built-in tests
./bin/godnscli test

# Query a domain
dig @localhost -p 53 www.example.lan A
```

## Test Categories

### 1. Basic Functionality Tests

```bash
# Test different record types
./bin/godnscli q www.example.lan          # A record (IPv4)
./bin/godnscli q www.example.lan -t AAAA  # AAAA record (IPv6)
./bin/godnscli q example.lan -t MX        # Mail exchange
./bin/godnscli q example.lan -t NS        # Name servers
./bin/godnscli q example.lan -t TXT       # Text records
./bin/godnscli q example.lan -t SOA       # Start of authority

# Health checks
./bin/godnscli h                          # Check health
./bin/godnscli h -v                       # Verbose health check
```

### 2. Using dig for Advanced Testing

```bash
# Basic query
dig @localhost -p 53 www.example.lan A

# TCP instead of UDP
dig @localhost -p 53 +tcp www.example.lan A

# Short output (just the IP)
dig @localhost -p 53 www.example.lan A +short

# Detailed output
dig @localhost -p 53 www.example.lan A +trace

# EDNS support
dig @localhost -p 53 +edns=0 www.example.lan A

# Query all record types
dig @localhost -p 53 example.lan ANY

# Reverse DNS lookup
dig @localhost -p 53 -x 192.168.100.10
```

### 3. Upstream Forwarding Test

```bash
# Should forward to upstream DNS
dig @localhost -p 53 google.com A
dig @localhost -p 53 github.com A

# Should get NXDOMAIN for non-existent
dig @localhost -p 53 thisdoesnotexist.invalid A
```

### 4. API Testing

```bash
# Health check
curl http://localhost:14000/health

# List all zones
curl http://localhost:14000/api/v1/zones | jq .

# Get specific zone
curl http://localhost:14000/api/v1/zones/example.lan | jq .

# Create a new zone
curl -X POST http://localhost:14000/api/v1/zones \
  -H "Content-Type: application/json" \
  -d '{
    "domain": "mytest.lan",
    "records": [
      {
        "name": "www.mytest.lan.",
        "type": "A",
        "ttl": 300,
        "value": "192.168.1.100"
      }
    ]
  }' | jq .

# Test the new zone
dig @localhost -p 53 www.mytest.lan A
```

### 5. Performance Testing

#### Simple Performance Test

```bash
# Time 100 queries
time for i in {1..100}; do
    dig @localhost -p 53 www.example.lan A +short > /dev/null
done
```

#### Load Test with dnsperf

Install dnsperf:

```bash
# Ubuntu/Debian
sudo apt-get install dnsperf

# macOS
brew install dnsperf
```

Create query file:

```bash
cat > /tmp/queries.txt << EOF
www.example.lan A
mail.example.lan A
ftp.example.lan A
db.example.lan A
router.home.lan A
nas.home.lan A
server.home.lan A
api.dev.local A
db.dev.local A
cache.dev.local A
EOF
```

Run load test:

```bash
# 10 concurrent clients for 30 seconds
dnsperf -s localhost -p 53 -d /tmp/queries.txt -c 10 -l 30

# More aggressive test
dnsperf -s localhost -p 53 -d /tmp/queries.txt -c 100 -l 60
```

### 6. Comprehensive Test Script

Save as `comprehensive-test.sh`:

```bash
#!/bin/bash
set -e

echo "=== GoDNS Comprehensive Test Suite ==="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

PASSED=0
FAILED=0

test_query() {
    local name=$1
    local type=$2
    local description=$3

    echo -n "Testing $description... "

    if dig @localhost -p 53 "$name" "$type" +short > /dev/null 2>&1; then
        echo -e "${GREEN}✓ PASS${NC}"
        ((PASSED++))
    else
        echo -e "${RED}✗ FAIL${NC}"
        ((FAILED++))
    fi
}

# Test health
echo "1. Health Checks"
test_query "www.example.lan" "A" "Health check"
echo ""

# Test record types
echo "2. Record Types"
test_query "www.example.lan" "A" "A record"
test_query "mail.example.lan" "A" "A record (mail)"
test_query "example.lan" "NS" "NS record"
test_query "example.lan" "SOA" "SOA record"
echo ""

# Test zones
echo "3. Multiple Zones"
test_query "router.home.lan" "A" "home.lan zone"
test_query "api.dev.local" "A" "dev.local zone"
test_query "master.k8s.local" "A" "k8s.local zone"
test_query "portainer.docker.local" "A" "docker.local zone"
echo ""

# Test upstream
echo "4. Upstream Forwarding"
test_query "google.com" "A" "Upstream DNS"
test_query "github.com" "A" "Upstream DNS"
echo ""

# Test protocols
echo "5. Protocol Support"
if dig @localhost -p 53 +tcp www.example.lan A +short > /dev/null 2>&1; then
    echo -e "TCP support: ${GREEN}✓ PASS${NC}"
    ((PASSED++))
else
    echo -e "TCP support: ${RED}✗ FAIL${NC}"
    ((FAILED++))
fi

if dig @localhost -p 53 +edns=0 www.example.lan A +short > /dev/null 2>&1; then
    echo -e "EDNS support: ${GREEN}✓ PASS${NC}"
    ((PASSED++))
else
    echo -e "EDNS support: ${RED}✗ FAIL${NC}"
    ((FAILED++))
fi
echo ""

# Summary
echo "=== Test Summary ==="
echo -e "Passed: ${GREEN}$PASSED${NC}"
echo -e "Failed: ${RED}$FAILED${NC}"

if [ $FAILED -eq 0 ]; then
    echo -e "\n${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "\n${RED}Some tests failed!${NC}"
    exit 1
fi
```

Make it executable and run:

```bash
chmod +x comprehensive-test.sh
./comprehensive-test.sh
```

### 7. Debugging Tests

When tests fail, use verbose mode:

```bash
# Verbose query
./bin/godnscli q www.example.lan -v

# Detailed dig output
dig @localhost -p 53 www.example.lan A +trace +stats

# Check server logs
docker-compose logs godns | tail -50

# Check Valkey data
docker-compose exec valkey valkey-cli -a mysecretpassword
> KEYS zone:*
> KEYS record:*
> GET zone:example.lan
```

### 8. Monitoring in Production

```bash
# Watch query logs
docker-compose logs -f godns

# Check metrics endpoint (if implemented)
curl http://localhost:14000/metrics

# Health check in production
watch -n 5 './bin/godnscli h'
```

## Test Data

Use the seeding scripts to populate test data:

```bash
# Automatic seeding (if DEVELOPMENT=true)
# Data is seeded on startup

# Manual seeding via API
### Seeding Test Data

The project includes a comprehensive test data seeding script:

```bash
./scripts/seed-test-data.sh
```
```

## Continuous Testing

For CI/CD pipelines:

```yaml
# .github/workflows/test.yml
name: Test
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2

      - name: Start services
        run: docker-compose up -d

      - name: Wait for services
        run: |
          timeout 60 bash -c 'until curl -f http://localhost:14000/health; do sleep 2; done'

      - name: Run tests
        run: |
          make build-cli
          ./bin/godnscli test

      - name: Run comprehensive tests
        run: ./comprehensive-test.sh
```

## Troubleshooting

### Server won't start

```bash
docker-compose logs godns
docker-compose logs valkey
```

### Queries not resolving

```bash
# Check if zone exists
curl http://localhost:14000/api/v1/zones | jq .

# Check specific record
dig @localhost -p 53 www.example.lan A +trace

# Verbose CLI query
./bin/godnscli q www.example.lan -v
```

### Upstream not working

```bash
# Check upstream configuration
docker-compose exec valkey valkey-cli -a mysecretpassword
> GET dns:config:upstream

# Test upstream directly
dig @8.8.8.8 google.com A
```

## Next Steps

- Set up automated testing in CI/CD
- Add performance benchmarks
- Implement monitoring and alerting
- Create custom test scenarios for your use case
