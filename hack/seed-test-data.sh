#!/bin/bash

# GoDNS Test Data Seeding Script
# Seeds the database with example zones and records for testing

set -e

API_URL="${GODNS_API_URL:-http://localhost:14000}"

echo "======================================"
echo "GoDNS Test Data Seeding Script"
echo "======================================"
echo "API URL: ${API_URL}"
echo ""

# Check if API is available
echo "Checking API health..."
if ! curl -s -f "${API_URL}/health" > /dev/null 2>&1; then
    echo "ERROR: GoDNS API is not available at ${API_URL}"
    echo "Please start the API server first: make run-api"
    exit 1
fi
echo "✓ API is available"
echo ""

# Seed Zone 1: home.lan (Home network)
echo "1. Creating zone: home.lan"
curl -s -X POST "${API_URL}/api/v1/zones" \
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
        "value": "192.168.1.10"
      },
      {
        "name": "server.home.lan.",
        "type": "A",
        "ttl": 300,
        "value": "192.168.1.100"
      },
      {
        "name": "printer.home.lan.",
        "type": "A",
        "ttl": 300,
        "value": "192.168.1.50"
      },
      {
        "name": "pi.home.lan.",
        "type": "A",
        "ttl": 300,
        "value": "192.168.1.200"
      }
    ]
  }' | jq .
echo "✓ Created home.lan"
echo ""

# Seed Zone 2: dev.local (Development environment)
echo "2. Creating zone: dev.local"
curl -s -X POST "${API_URL}/api/v1/zones" \
  -H "Content-Type: application/json" \
  -d '{
    "domain": "dev.local",
    "records": [
      {
        "name": "api.dev.local.",
        "type": "A",
        "ttl": 300,
        "value": "127.0.0.1"
      },
      {
        "name": "db.dev.local.",
        "type": "A",
        "ttl": 300,
        "value": "127.0.0.1"
      },
      {
        "name": "cache.dev.local.",
        "type": "A",
        "ttl": 300,
        "value": "127.0.0.1"
      },
      {
        "name": "web.dev.local.",
        "type": "A",
        "ttl": 300,
        "value": "127.0.0.1"
      },
      {
        "name": "mail.dev.local.",
        "type": "A",
        "ttl": 300,
        "value": "127.0.0.1"
      }
    ]
  }' | jq .
echo "✓ Created dev.local"
echo ""

# Seed Zone 3: k8s.local (Kubernetes cluster)
echo "3. Creating zone: k8s.local"
curl -s -X POST "${API_URL}/api/v1/zones" \
  -H "Content-Type: application/json" \
  -d '{
    "domain": "k8s.local",
    "records": [
      {
        "name": "master.k8s.local.",
        "type": "A",
        "ttl": 300,
        "value": "10.0.1.10"
      },
      {
        "name": "worker1.k8s.local.",
        "type": "A",
        "ttl": 300,
        "value": "10.0.1.20"
      },
      {
        "name": "worker2.k8s.local.",
        "type": "A",
        "ttl": 300,
        "value": "10.0.1.21"
      },
      {
        "name": "worker3.k8s.local.",
        "type": "A",
        "ttl": 300,
        "value": "10.0.1.22"
      },
      {
        "name": "ingress.k8s.local.",
        "type": "A",
        "ttl": 300,
        "value": "10.0.1.100"
      }
    ]
  }' | jq .
echo "✓ Created k8s.local"
echo ""

# Seed Zone 4: docker.local (Docker services)
echo "4. Creating zone: docker.local"
curl -s -X POST "${API_URL}/api/v1/zones" \
  -H "Content-Type: application/json" \
  -d '{
    "domain": "docker.local",
    "records": [
      {
        "name": "portainer.docker.local.",
        "type": "A",
        "ttl": 300,
        "value": "172.17.0.10"
      },
      {
        "name": "registry.docker.local.",
        "type": "A",
        "ttl": 300,
        "value": "172.17.0.20"
      },
      {
        "name": "traefik.docker.local.",
        "type": "A",
        "ttl": 300,
        "value": "172.17.0.30"
      },
      {
        "name": "grafana.docker.local.",
        "type": "A",
        "ttl": 300,
        "value": "172.17.0.40"
      },
      {
        "name": "prometheus.docker.local.",
        "type": "A",
        "ttl": 300,
        "value": "172.17.0.41"
      }
    ]
  }' | jq .
echo "✓ Created docker.local"
echo ""

# Seed Zone 5: lab.internal (Lab environment)
echo "5. Creating zone: lab.internal"
curl -s -X POST "${API_URL}/api/v1/zones" \
  -H "Content-Type: application/json" \
  -d '{
    "domain": "lab.internal",
    "records": [
      {
        "name": "test1.lab.internal.",
        "type": "A",
        "ttl": 60,
        "value": "10.10.10.1"
      },
      {
        "name": "test2.lab.internal.",
        "type": "A",
        "ttl": 60,
        "value": "10.10.10.2"
      },
      {
        "name": "test3.lab.internal.",
        "type": "A",
        "ttl": 60,
        "value": "10.10.10.3"
      },
      {
        "name": "vpn.lab.internal.",
        "type": "A",
        "ttl": 300,
        "value": "10.10.10.254"
      },
      {
        "name": "gateway.lab.internal.",
        "type": "A",
        "ttl": 300,
        "value": "10.10.10.1"
      }
    ]
  }' | jq .
echo "✓ Created lab.internal"
echo ""

# Seed Zone 6: example.lan (Example/demo zone)
echo "6. Creating zone: example.lan"
curl -s -X POST "${API_URL}/api/v1/zones" \
  -H "Content-Type: application/json" \
  -d '{
    "domain": "example.lan",
    "records": [
      {
        "name": "www.example.lan.",
        "type": "A",
        "ttl": 300,
        "value": "192.168.100.10"
      },
      {
        "name": "mail.example.lan.",
        "type": "A",
        "ttl": 300,
        "value": "192.168.100.20"
      },
      {
        "name": "ftp.example.lan.",
        "type": "A",
        "ttl": 300,
        "value": "192.168.100.30"
      },
      {
        "name": "db.example.lan.",
        "type": "A",
        "ttl": 300,
        "value": "192.168.100.40"
      },
      {
        "name": "ns1.example.lan.",
        "type": "A",
        "ttl": 3600,
        "value": "192.168.100.1"
      },
      {
        "name": "ns2.example.lan.",
        "type": "A",
        "ttl": 3600,
        "value": "192.168.100.2"
      }
    ]
  }' | jq .
echo "✓ Created example.lan"
echo ""

# Summary
echo "======================================"
echo "Test Data Seeding Complete!"
echo "======================================"
echo ""
echo "Created zones:"
echo "  1. home.lan      - 5 records (home network devices)"
echo "  2. dev.local     - 5 records (development services)"
echo "  3. k8s.local     - 5 records (Kubernetes cluster)"
echo "  4. docker.local  - 5 records (Docker services)"
echo "  5. lab.internal  - 5 records (lab environment)"
echo "  6. example.lan   - 6 records (example/demo)"
echo ""
echo "Total: 6 zones with 31 DNS records"
echo ""
echo "View all zones:"
echo "  curl ${API_URL}/api/v1/zones | jq ."
echo ""
echo "Test DNS queries (if DNS server is running):"
echo "  dig @localhost router.home.lan"
echo "  dig @localhost api.dev.local"
echo "  dig @localhost master.k8s.local"
echo ""
