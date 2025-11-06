#!/bin/bash

# GoDNS HTTP API Test Script
# This script demonstrates basic CRUD operations using the GoDNS HTTP API

set -e

API_URL="http://localhost:14000"
ZONE="test.lan"

echo "======================================"
echo "GoDNS HTTP API Test Script"
echo "======================================"
echo ""

# Health check
echo "1. Testing health endpoint..."
curl -s "${API_URL}/health" | jq .
echo ""

# Create a zone
echo "2. Creating zone: ${ZONE}"
curl -s -X POST "${API_URL}/api/v1/zones" \
  -H "Content-Type: application/json" \
  -d "{
    \"domain\": \"${ZONE}\",
    \"records\": [
      {
        \"name\": \"www.${ZONE}.\",
        \"type\": \"A\",
        \"ttl\": 300,
        \"value\": \"192.168.1.100\"
      },
      {
        \"name\": \"mail.${ZONE}.\",
        \"type\": \"A\",
        \"ttl\": 300,
        \"value\": \"192.168.1.101\"
      }
    ]
  }" | jq .
echo ""

# List all zones
echo "3. Listing all zones..."
curl -s "${API_URL}/api/v1/zones" | jq .
echo ""

# Get specific zone
echo "4. Getting zone: ${ZONE}"
curl -s "${API_URL}/api/v1/zones/${ZONE}" | jq .
echo ""

# Add a record
echo "5. Adding a new record (api.${ZONE})"
curl -s -X POST "${API_URL}/api/v1/zones/${ZONE}/records" \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"api.${ZONE}.\",
    \"type\": \"A\",
    \"ttl\": 300,
    \"value\": \"192.168.1.150\"
  }" | jq .
echo ""

# Get the zone again to see the new record
echo "6. Getting updated zone..."
curl -s "${API_URL}/api/v1/zones/${ZONE}" | jq .
echo ""

# Get specific record
echo "7. Getting specific record (www.${ZONE}./A)"
curl -s "${API_URL}/api/v1/zones/${ZONE}/records/www.${ZONE}./A" | jq .
echo ""

# Update a record
echo "8. Updating record (www.${ZONE}./A)"
curl -s -X PUT "${API_URL}/api/v1/zones/${ZONE}/records/www.${ZONE}./A" \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"www.${ZONE}.\",
    \"type\": \"A\",
    \"ttl\": 600,
    \"value\": \"192.168.1.200\"
  }" | jq .
echo ""

# Verify the update
echo "9. Verifying updated record..."
curl -s "${API_URL}/api/v1/zones/${ZONE}/records/www.${ZONE}./A" | jq .
echo ""

# Delete a record
echo "10. Deleting record (api.${ZONE}./A)"
curl -s -X DELETE "${API_URL}/api/v1/zones/${ZONE}/records/api.${ZONE}./A"
echo "Record deleted"
echo ""

# Verify deletion
echo "11. Getting zone after record deletion..."
curl -s "${API_URL}/api/v1/zones/${ZONE}" | jq .
echo ""

# Update entire zone
echo "12. Updating entire zone..."
curl -s -X PUT "${API_URL}/api/v1/zones/${ZONE}" \
  -H "Content-Type: application/json" \
  -d "{
    \"domain\": \"${ZONE}\",
    \"records\": [
      {
        \"name\": \"www.${ZONE}.\",
        \"type\": \"A\",
        \"ttl\": 300,
        \"value\": \"192.168.1.100\"
      },
      {
        \"name\": \"ftp.${ZONE}.\",
        \"type\": \"A\",
        \"ttl\": 300,
        \"value\": \"192.168.1.102\"
      }
    ]
  }" | jq .
echo ""

# Delete the zone
echo "13. Deleting zone: ${ZONE}"
curl -s -X DELETE "${API_URL}/api/v1/zones/${ZONE}"
echo "Zone deleted"
echo ""

# Verify deletion
echo "14. Listing all zones (should be empty or without ${ZONE})..."
curl -s "${API_URL}/api/v1/zones" | jq .
echo ""

echo "======================================"
echo "Test completed successfully!"
echo "======================================"
