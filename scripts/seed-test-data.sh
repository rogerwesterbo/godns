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
        "name": "home.lan.",
        "type": "SOA",
        "ttl": 3600,
        "soa_mname": "ns1.home.lan.",
        "soa_rname": "hostmaster.home.lan.",
        "soa_serial": 2024110601,
        "soa_refresh": 3600,
        "soa_retry": 1800,
        "soa_expire": 604800,
        "soa_minimum": 300
      },
      {
        "name": "home.lan.",
        "type": "NS",
        "ttl": 3600,
        "value": "ns1.home.lan."
      },
      {
        "name": "home.lan.",
        "type": "NS",
        "ttl": 3600,
        "value": "ns2.home.lan."
      },
      {
        "name": "ns1.home.lan.",
        "type": "A",
        "ttl": 3600,
        "value": "192.168.1.1"
      },
      {
        "name": "ns2.home.lan.",
        "type": "A",
        "ttl": 3600,
        "value": "192.168.1.2"
      },
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
        "name": "nas.home.lan.",
        "type": "AAAA",
        "ttl": 300,
        "value": "fd00::10"
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
      },
      {
        "name": "www.home.lan.",
        "type": "CNAME",
        "ttl": 300,
        "value": "server.home.lan."
      },
      {
        "name": "home.lan.",
        "type": "TXT",
        "ttl": 300,
        "value": "v=spf1 ip4:192.168.1.0/24 -all"
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
        "name": "dev.local.",
        "type": "SOA",
        "ttl": 3600,
        "soa_mname": "ns1.dev.local.",
        "soa_rname": "hostmaster.dev.local.",
        "soa_serial": 2024110601,
        "soa_refresh": 7200,
        "soa_retry": 3600,
        "soa_expire": 1209600,
        "soa_minimum": 300
      },
      {
        "name": "dev.local.",
        "type": "NS",
        "ttl": 3600,
        "value": "ns1.dev.local."
      },
      {
        "name": "ns1.dev.local.",
        "type": "A",
        "ttl": 3600,
        "value": "127.0.0.1"
      },
      {
        "name": "api.dev.local.",
        "type": "A",
        "ttl": 60,
        "value": "127.0.0.1"
      },
      {
        "name": "db.dev.local.",
        "type": "A",
        "ttl": 60,
        "value": "127.0.0.1"
      },
      {
        "name": "cache.dev.local.",
        "type": "A",
        "ttl": 60,
        "value": "127.0.0.1"
      },
      {
        "name": "web.dev.local.",
        "type": "A",
        "ttl": 60,
        "value": "127.0.0.1"
      },
      {
        "name": "mail.dev.local.",
        "type": "A",
        "ttl": 60,
        "value": "127.0.0.1"
      },
      {
        "name": "dev.local.",
        "type": "MX",
        "ttl": 300,
        "mx_priority": 10,
        "mx_host": "mail.dev.local."
      },
      {
        "name": "www.dev.local.",
        "type": "CNAME",
        "ttl": 60,
        "value": "web.dev.local."
      },
      {
        "name": "app.dev.local.",
        "type": "CNAME",
        "ttl": 60,
        "value": "api.dev.local."
      },
      {
        "name": "_http._tcp.dev.local.",
        "type": "SRV",
        "ttl": 60,
        "srv_priority": 10,
        "srv_weight": 60,
        "srv_port": 80,
        "srv_target": "web.dev.local."
      },
      {
        "name": "dev.local.",
        "type": "TXT",
        "ttl": 300,
        "value": "v=spf1 mx -all"
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
        "name": "k8s.local.",
        "type": "SOA",
        "ttl": 3600,
        "soa_mname": "ns1.k8s.local.",
        "soa_rname": "hostmaster.k8s.local.",
        "soa_serial": 2024110601,
        "soa_refresh": 3600,
        "soa_retry": 1800,
        "soa_expire": 604800,
        "soa_minimum": 300
      },
      {
        "name": "k8s.local.",
        "type": "NS",
        "ttl": 3600,
        "value": "ns1.k8s.local."
      },
      {
        "name": "ns1.k8s.local.",
        "type": "A",
        "ttl": 3600,
        "value": "10.0.1.1"
      },
      {
        "name": "master.k8s.local.",
        "type": "A",
        "ttl": 300,
        "value": "10.0.1.10"
      },
      {
        "name": "master.k8s.local.",
        "type": "AAAA",
        "ttl": 300,
        "value": "fd00:10::10"
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
      },
      {
        "name": "api.k8s.local.",
        "type": "CNAME",
        "ttl": 300,
        "value": "master.k8s.local."
      },
      {
        "name": "*.apps.k8s.local.",
        "type": "A",
        "ttl": 300,
        "value": "10.0.1.100"
      },
      {
        "name": "_etcd-server._tcp.k8s.local.",
        "type": "SRV",
        "ttl": 300,
        "srv_priority": 10,
        "srv_weight": 100,
        "srv_port": 2380,
        "srv_target": "master.k8s.local."
      },
      {
        "name": "k8s.local.",
        "type": "TXT",
        "ttl": 300,
        "value": "kubernetes-cluster"
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
        "name": "docker.local.",
        "type": "SOA",
        "ttl": 3600,
        "soa_mname": "ns1.docker.local.",
        "soa_rname": "hostmaster.docker.local.",
        "soa_serial": 2024110601,
        "soa_refresh": 3600,
        "soa_retry": 1800,
        "soa_expire": 604800,
        "soa_minimum": 300
      },
      {
        "name": "docker.local.",
        "type": "NS",
        "ttl": 3600,
        "value": "ns1.docker.local."
      },
      {
        "name": "ns1.docker.local.",
        "type": "A",
        "ttl": 3600,
        "value": "172.17.0.1"
      },
      {
        "name": "portainer.docker.local.",
        "type": "A",
        "ttl": 60,
        "value": "172.17.0.10"
      },
      {
        "name": "registry.docker.local.",
        "type": "A",
        "ttl": 60,
        "value": "172.17.0.20"
      },
      {
        "name": "traefik.docker.local.",
        "type": "A",
        "ttl": 60,
        "value": "172.17.0.30"
      },
      {
        "name": "grafana.docker.local.",
        "type": "A",
        "ttl": 60,
        "value": "172.17.0.40"
      },
      {
        "name": "prometheus.docker.local.",
        "type": "A",
        "ttl": 60,
        "value": "172.17.0.41"
      },
      {
        "name": "monitor.docker.local.",
        "type": "CNAME",
        "ttl": 60,
        "value": "grafana.docker.local."
      },
      {
        "name": "metrics.docker.local.",
        "type": "CNAME",
        "ttl": 60,
        "value": "prometheus.docker.local."
      },
      {
        "name": "_metrics._tcp.docker.local.",
        "type": "SRV",
        "ttl": 60,
        "srv_priority": 10,
        "srv_weight": 100,
        "srv_port": 9090,
        "srv_target": "prometheus.docker.local."
      },
      {
        "name": "docker.local.",
        "type": "TXT",
        "ttl": 300,
        "value": "docker-network-services"
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
        "name": "lab.internal.",
        "type": "SOA",
        "ttl": 3600,
        "soa_mname": "ns1.lab.internal.",
        "soa_rname": "hostmaster.lab.internal.",
        "soa_serial": 2024110601,
        "soa_refresh": 3600,
        "soa_retry": 1800,
        "soa_expire": 604800,
        "soa_minimum": 60
      },
      {
        "name": "lab.internal.",
        "type": "NS",
        "ttl": 3600,
        "value": "ns1.lab.internal."
      },
      {
        "name": "ns1.lab.internal.",
        "type": "A",
        "ttl": 3600,
        "value": "10.10.10.1"
      },
      {
        "name": "test1.lab.internal.",
        "type": "A",
        "ttl": 60,
        "value": "10.10.10.10"
      },
      {
        "name": "test2.lab.internal.",
        "type": "A",
        "ttl": 60,
        "value": "10.10.10.11"
      },
      {
        "name": "test3.lab.internal.",
        "type": "A",
        "ttl": 60,
        "value": "10.10.10.12"
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
      },
      {
        "name": "lb.lab.internal.",
        "type": "A",
        "ttl": 60,
        "value": "10.10.10.10"
      },
      {
        "name": "lb.lab.internal.",
        "type": "A",
        "ttl": 60,
        "value": "10.10.10.11"
      },
      {
        "name": "lb.lab.internal.",
        "type": "A",
        "ttl": 60,
        "value": "10.10.10.12"
      },
      {
        "name": "*.test.lab.internal.",
        "type": "A",
        "ttl": 60,
        "value": "10.10.10.100"
      },
      {
        "name": "lab.internal.",
        "type": "TXT",
        "ttl": 300,
        "value": "testing-environment"
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
        "name": "example.lan.",
        "type": "SOA",
        "ttl": 3600,
        "soa_mname": "ns1.example.lan.",
        "soa_rname": "hostmaster.example.lan.",
        "soa_serial": 2024110601,
        "soa_refresh": 3600,
        "soa_retry": 1800,
        "soa_expire": 604800,
        "soa_minimum": 300
      },
      {
        "name": "example.lan.",
        "type": "NS",
        "ttl": 3600,
        "value": "ns1.example.lan."
      },
      {
        "name": "example.lan.",
        "type": "NS",
        "ttl": 3600,
        "value": "ns2.example.lan."
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
      },
      {
        "name": "example.lan.",
        "type": "A",
        "ttl": 300,
        "value": "192.168.100.10"
      },
      {
        "name": "example.lan.",
        "type": "AAAA",
        "ttl": 300,
        "value": "2001:db8::1"
      },
      {
        "name": "www.example.lan.",
        "type": "A",
        "ttl": 300,
        "value": "192.168.100.10"
      },
      {
        "name": "www.example.lan.",
        "type": "AAAA",
        "ttl": 300,
        "value": "2001:db8::1"
      },
      {
        "name": "mail.example.lan.",
        "type": "A",
        "ttl": 300,
        "value": "192.168.100.20"
      },
      {
        "name": "example.lan.",
        "type": "MX",
        "ttl": 300,
        "mx_priority": 10,
        "mx_host": "mail.example.lan."
      },
      {
        "name": "example.lan.",
        "type": "MX",
        "ttl": 300,
        "mx_priority": 20,
        "mx_host": "mail2.example.lan."
      },
      {
        "name": "mail2.example.lan.",
        "type": "A",
        "ttl": 300,
        "value": "192.168.100.21"
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
        "name": "api.example.lan.",
        "type": "CNAME",
        "ttl": 300,
        "value": "www.example.lan."
      },
      {
        "name": "blog.example.lan.",
        "type": "CNAME",
        "ttl": 300,
        "value": "www.example.lan."
      },
      {
        "name": "_ldap._tcp.example.lan.",
        "type": "SRV",
        "ttl": 300,
        "srv_priority": 10,
        "srv_weight": 100,
        "srv_port": 389,
        "srv_target": "ldap.example.lan."
      },
      {
        "name": "ldap.example.lan.",
        "type": "A",
        "ttl": 300,
        "value": "192.168.100.50"
      },
      {
        "name": "example.lan.",
        "type": "TXT",
        "ttl": 300,
        "value": "v=spf1 mx ip4:192.168.100.0/24 -all"
      },
      {
        "name": "_dmarc.example.lan.",
        "type": "TXT",
        "ttl": 300,
        "value": "v=DMARC1; p=quarantine; rua=mailto:dmarc@example.lan"
      },
      {
        "name": "default._domainkey.example.lan.",
        "type": "TXT",
        "ttl": 300,
        "value": "v=DKIM1; k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC"
      }
    ]
  }' | jq .
echo "✓ Created example.lan"
echo ""

# Summary
echo "✅ Test data seeding complete!"
echo ""
echo "Summary:"
echo "--------"
echo "Created 6 DNS zones with diverse record types:"
echo "  - home.lan (13 records): SOA, NS, A, AAAA, CNAME, TXT"
echo "  - dev.local (13 records): SOA, NS, A, MX, CNAME, SRV, TXT"
echo "  - k8s.local (13 records): SOA, NS, A, AAAA, CNAME, wildcard, SRV, TXT"
echo "  - docker.local (12 records): SOA, NS, A, CNAME, SRV, TXT"
echo "  - lab.internal (13 records): SOA, NS, A (with load balancing), wildcard, TXT"
echo "  - example.lan (22 records): SOA, NS, A, AAAA, MX, CNAME, SRV, TXT (SPF/DKIM/DMARC)"
echo ""
echo "Total: 86 DNS records across all zones"
echo ""
echo "Record types included:"
echo "  - SOA: Start of Authority (zone metadata)"
echo "  - NS: Name Server records"
echo "  - A: IPv4 addresses"
echo "  - AAAA: IPv6 addresses"
echo "  - CNAME: Canonical name aliases"
echo "  - MX: Mail exchange records"
echo "  - SRV: Service discovery records"
echo "  - TXT: Text records (SPF, DKIM, DMARC)"
echo "  - Wildcards: Catch-all domains (*.apps, *.test)"
echo ""
