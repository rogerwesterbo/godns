#!/bin/bash
# Script to add a test DNS zone to GoDNS via Valkey
# Usage: ./add-test-zone.sh [domain] [ip-address]

set -e

DOMAIN=${1:-"test.lan."}
IP=${2:-"192.168.1.100"}

# Ensure domain ends with a dot
if [[ ! "$DOMAIN" =~ \.$ ]]; then
    DOMAIN="${DOMAIN}."
fi

echo "🔧 Adding test DNS zone to GoDNS"
echo "================================="
echo ""
echo "Domain: $DOMAIN"
echo "IP:     $IP"
echo ""

# Create zone JSON
ZONE_JSON=$(cat <<EOF
{
  "domain": "$DOMAIN",
  "records": [
    {
      "name": "$DOMAIN",
      "type": "A",
      "ttl": 300,
      "value": "$IP"
    },
    {
      "name": "$DOMAIN",
      "type": "NS",
      "ttl": 3600,
      "value": "ns1.$DOMAIN"
    },
    {
      "name": "$DOMAIN",
      "type": "SOA",
      "ttl": 3600,
      "value": "ns1.$DOMAIN hostmaster.$DOMAIN 1 3600 600 604800 300"
    }
  ]
}
EOF
)

echo "📝 Zone configuration:"
echo "$ZONE_JSON" | jq '.' 2>/dev/null || echo "$ZONE_JSON"
echo ""

# Check if docker-compose is running
if ! docker-compose ps | grep -q "valkey.*Up"; then
    echo "❌ Error: Valkey container is not running"
    echo "   Start it with: docker-compose up -d"
    exit 1
fi

# Add to Valkey
echo "💾 Saving to Valkey..."
echo "SET dns:zone:$DOMAIN '$ZONE_JSON'" | docker-compose exec -T valkey valkey-cli -a mysecretpassword > /dev/null 2>&1

if [ $? -eq 0 ]; then
    echo "✅ Zone added successfully!"
    echo ""
    echo "🧪 Test with:"
    echo "   ./bin/godnscli q $DOMAIN"
    echo ""
    echo "📋 List all zones:"
    echo "   docker-compose exec valkey valkey-cli -a mysecretpassword KEYS 'dns:zone:*'"
else
    echo "❌ Failed to add zone"
    exit 1
fi
