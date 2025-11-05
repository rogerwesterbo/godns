#!/bin/bash
#
# Integration test for DNS Export API
# This script tests the export functionality end-to-end
#

set -e

API_URL="${GODNS_API_URL:-http://localhost:14000}"
TEST_ZONE="test-export.lan"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "DNS Export API Integration Test"
echo "================================"
echo ""

# Check if API is available
if ! curl -s -f "${API_URL}/health" > /dev/null 2>&1; then
    echo -e "${RED}Error: GoDNS API is not available at ${API_URL}${NC}"
    echo "Please start the API server first."
    exit 1
fi

echo -e "${GREEN}✓ API is available${NC}"
echo ""

# Clean up function
cleanup() {
    echo "Cleaning up test zone..."
    curl -s -X DELETE "${API_URL}/api/v1/zones/${TEST_ZONE}" > /dev/null 2>&1 || true
}

# Register cleanup on exit
trap cleanup EXIT

# Create a test zone
echo "Creating test zone: ${TEST_ZONE}"
CREATE_RESPONSE=$(curl -s -X POST "${API_URL}/api/v1/zones" \
    -H "Content-Type: application/json" \
    -d "{
        \"domain\": \"${TEST_ZONE}\",
        \"records\": [
            {\"name\": \"${TEST_ZONE}\", \"type\": \"A\", \"ttl\": 300, \"value\": \"192.168.1.1\"},
            {\"name\": \"www.${TEST_ZONE}\", \"type\": \"A\", \"ttl\": 300, \"value\": \"192.168.1.2\"},
            {\"name\": \"mail.${TEST_ZONE}\", \"type\": \"A\", \"ttl\": 300, \"value\": \"192.168.1.3\"},
            {\"name\": \"${TEST_ZONE}\", \"type\": \"MX\", \"ttl\": 300, \"value\": \"10 mail.${TEST_ZONE}\"},
            {\"name\": \"${TEST_ZONE}\", \"type\": \"TXT\", \"ttl\": 300, \"value\": \"v=spf1 mx -all\"}
        ]
    }")

if echo "$CREATE_RESPONSE" | grep -q "error"; then
    echo -e "${RED}✗ Failed to create test zone${NC}"
    echo "$CREATE_RESPONSE"
    exit 1
fi

echo -e "${GREEN}✓ Test zone created${NC}"
echo ""

# Test 1: Export in BIND format
echo "Test 1: Export zone in BIND format"
echo "-----------------------------------"
BIND_EXPORT=$(curl -s -X GET "${API_URL}/api/v1/export/${TEST_ZONE}?format=bind")

if echo "$BIND_EXPORT" | grep -q "\$ORIGIN ${TEST_ZONE}"; then
    echo -e "${GREEN}✓ BIND export contains \$ORIGIN${NC}"
else
    echo -e "${RED}✗ BIND export missing \$ORIGIN${NC}"
    exit 1
fi

if echo "$BIND_EXPORT" | grep -q "IN.*A.*192.168.1.1"; then
    echo -e "${GREEN}✓ BIND export contains A record${NC}"
else
    echo -e "${RED}✗ BIND export missing A record${NC}"
    exit 1
fi

if echo "$BIND_EXPORT" | grep -q "IN.*MX"; then
    echo -e "${GREEN}✓ BIND export contains MX record${NC}"
else
    echo -e "${RED}✗ BIND export missing MX record${NC}"
    exit 1
fi

echo ""

# Test 2: Export in CoreDNS format
echo "Test 2: Export zone in CoreDNS format"
echo "--------------------------------------"
COREDNS_EXPORT=$(curl -s -X GET "${API_URL}/api/v1/export/${TEST_ZONE}?format=coredns")

if echo "$COREDNS_EXPORT" | grep -q "test-export.lan {"; then
    echo -e "${GREEN}✓ CoreDNS export contains zone block${NC}"
else
    echo -e "${RED}✗ CoreDNS export missing zone block${NC}"
    exit 1
fi

if echo "$COREDNS_EXPORT" | grep -q "file /etc/coredns/zones/"; then
    echo -e "${GREEN}✓ CoreDNS export contains file plugin${NC}"
else
    echo -e "${RED}✗ CoreDNS export missing file plugin${NC}"
    exit 1
fi

echo ""

# Test 3: Export in PowerDNS format
echo "Test 3: Export zone in PowerDNS format"
echo "---------------------------------------"
POWERDNS_EXPORT=$(curl -s -X GET "${API_URL}/api/v1/export/${TEST_ZONE}?format=powerdns")

if echo "$POWERDNS_EXPORT" | grep -q '"name".*"test-export.lan."'; then
    echo -e "${GREEN}✓ PowerDNS export contains zone name${NC}"
else
    echo -e "${RED}✗ PowerDNS export missing zone name${NC}"
    exit 1
fi

if echo "$POWERDNS_EXPORT" | grep -q '"rrsets"'; then
    echo -e "${GREEN}✓ PowerDNS export contains rrsets${NC}"
else
    echo -e "${RED}✗ PowerDNS export missing rrsets${NC}"
    exit 1
fi

# Validate JSON
if echo "$POWERDNS_EXPORT" | grep -A 10000 "^{" | jq . > /dev/null 2>&1; then
    echo -e "${GREEN}✓ PowerDNS export is valid JSON${NC}"
else
    echo -e "${RED}✗ PowerDNS export is not valid JSON${NC}"
    exit 1
fi

echo ""

# Test 4: Export all zones
echo "Test 4: Export all zones"
echo "------------------------"
ALL_ZONES_EXPORT=$(curl -s -X GET "${API_URL}/api/v1/export?format=bind")

if echo "$ALL_ZONES_EXPORT" | grep -q "${TEST_ZONE}"; then
    echo -e "${GREEN}✓ Export all zones includes test zone${NC}"
else
    echo -e "${RED}✗ Export all zones missing test zone${NC}"
    exit 1
fi

echo ""

# Test 5: Invalid format handling
echo "Test 5: Invalid format handling"
echo "--------------------------------"
INVALID_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${API_URL}/api/v1/export/${TEST_ZONE}?format=invalid")
HTTP_CODE=$(echo "$INVALID_RESPONSE" | tail -n 1)

if [ "$HTTP_CODE" = "400" ]; then
    echo -e "${GREEN}✓ Invalid format returns 400 status${NC}"
else
    echo -e "${RED}✗ Invalid format should return 400, got ${HTTP_CODE}${NC}"
    exit 1
fi

echo ""

# Test 6: Non-existent zone
echo "Test 6: Non-existent zone handling"
echo "-----------------------------------"
NOTFOUND_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${API_URL}/api/v1/export/nonexistent.zone?format=bind")
HTTP_CODE=$(echo "$NOTFOUND_RESPONSE" | tail -n 1)

if [ "$HTTP_CODE" = "404" ]; then
    echo -e "${GREEN}✓ Non-existent zone returns 404 status${NC}"
else
    echo -e "${RED}✗ Non-existent zone should return 404, got ${HTTP_CODE}${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}All tests passed! ✓${NC}"
echo -e "${GREEN}================================${NC}"
