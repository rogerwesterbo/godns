#!/bin/bash
# Migration script to set enabled=true for all existing zones

set -e

echo "Starting zone migration to add enabled field..."

# Get API URL from environment or use default
API_URL="${GODNS_API_URL:-http://localhost:8080}"
TOKEN="${GODNS_TOKEN}"

if [ -z "$TOKEN" ]; then
    echo "Error: GODNS_TOKEN environment variable must be set"
    echo "Usage: GODNS_TOKEN='your-token' ./migrate-zone-enabled.sh"
    exit 1
fi

# Get list of all zones
echo "Fetching zones..."
ZONES_JSON=$(curl -s -H "Authorization: Bearer $TOKEN" "$API_URL/api/v1/zones")

# Extract zone domains using jq
DOMAINS=$(echo "$ZONES_JSON" | jq -r '.[].domain')

if [ -z "$DOMAINS" ]; then
    echo "No zones found or failed to fetch zones"
    exit 1
fi

# Count zones
ZONE_COUNT=$(echo "$DOMAINS" | wc -l | tr -d ' ')
echo "Found $ZONE_COUNT zones to migrate"

# Migrate each zone
COUNTER=0
for DOMAIN in $DOMAINS; do
    COUNTER=$((COUNTER + 1))
    echo "[$COUNTER/$ZONE_COUNT] Enabling zone: $DOMAIN"
    
    # Set zone status to enabled
    RESPONSE=$(curl -s -w "\n%{http_code}" \
        -X PATCH \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d '{"enabled": true}' \
        "$API_URL/api/v1/zones/$DOMAIN/status")
    
    HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
    
    if [ "$HTTP_CODE" = "204" ]; then
        echo "  ✓ Successfully enabled $DOMAIN"
    else
        echo "  ✗ Failed to enable $DOMAIN (HTTP $HTTP_CODE)"
        echo "$RESPONSE" | head -n-1
    fi
done

echo ""
echo "Migration complete! Migrated $ZONE_COUNT zones."
echo "All zones have been set to enabled=true"
