#!/bin/bash
#
# Export DNS zones example script
# This script demonstrates how to use the GoDNS export API
#

set -e

API_URL="${GODNS_API_URL:-http://localhost:14082}"
OUTPUT_DIR="./exported-zones"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}GoDNS Zone Export Example${NC}"
echo "================================"
echo ""

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Function to export all zones in a specific format
export_all_zones() {
    local format=$1
    echo -e "${YELLOW}Exporting all zones in ${format} format...${NC}"
    
    curl -s -X GET "${API_URL}/api/v1/export?format=${format}" \
        -o "${OUTPUT_DIR}/all-zones-${format}.txt"
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ All zones exported to: ${OUTPUT_DIR}/all-zones-${format}.txt${NC}"
    else
        echo -e "${RED}✗ Failed to export zones in ${format} format${NC}"
        return 1
    fi
}

# Function to export a specific zone
export_zone() {
    local domain=$1
    local format=$2
    echo -e "${YELLOW}Exporting ${domain} in ${format} format...${NC}"
    
    curl -s -X GET "${API_URL}/api/v1/export/${domain}?format=${format}" \
        -o "${OUTPUT_DIR}/${domain}-${format}.txt"
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Zone ${domain} exported to: ${OUTPUT_DIR}/${domain}-${format}.txt${NC}"
    else
        echo -e "${RED}✗ Failed to export ${domain}${NC}"
        return 1
    fi
}

# Function to list available zones
list_zones() {
    echo -e "${YELLOW}Listing available zones...${NC}"
    local zones=$(curl -s -X GET "${API_URL}/api/v1/zones" | jq -r '.[].domain' 2>/dev/null)
    
    if [ -z "$zones" ]; then
        echo -e "${RED}No zones found or unable to connect to API${NC}"
        echo "Make sure the GoDNS API server is running at: ${API_URL}"
        exit 1
    fi
    
    echo -e "${GREEN}Available zones:${NC}"
    echo "$zones"
    echo ""
}

# Main script
echo "API URL: ${API_URL}"
echo "Output directory: ${OUTPUT_DIR}"
echo ""

# Check if API is reachable
if ! curl -s -f "${API_URL}/health" > /dev/null 2>&1; then
    echo -e "${RED}Error: Cannot reach GoDNS API at ${API_URL}${NC}"
    echo "Please ensure the API server is running."
    exit 1
fi

echo -e "${GREEN}✓ API is reachable${NC}"
echo ""

# List zones
list_zones

# Example 1: Export all zones in BIND format
echo "Example 1: Export all zones in BIND format"
echo "-------------------------------------------"
export_all_zones "bind"
echo ""

# Example 2: Export all zones in CoreDNS format
echo "Example 2: Export all zones in CoreDNS format"
echo "----------------------------------------------"
export_all_zones "coredns"
echo ""

# Example 3: Export all zones in PowerDNS format
echo "Example 3: Export all zones in PowerDNS format"
echo "-----------------------------------------------"
export_all_zones "powerdns"
echo ""

# Example 4: Export specific zones (if any exist)
zones_list=$(curl -s -X GET "${API_URL}/api/v1/zones" | jq -r '.[].domain' 2>/dev/null | head -n 2)

if [ ! -z "$zones_list" ]; then
    echo "Example 4: Export specific zones in different formats"
    echo "------------------------------------------------------"
    
    for zone in $zones_list; do
        echo ""
        export_zone "$zone" "bind"
        export_zone "$zone" "coredns"
        export_zone "$zone" "powerdns"
    done
fi

echo ""
echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}Export complete!${NC}"
echo ""
echo "Exported files are in: ${OUTPUT_DIR}/"
ls -lh "$OUTPUT_DIR/" | tail -n +2

echo ""
echo -e "${YELLOW}Tip: You can view the exported files with:${NC}"
echo "  cat ${OUTPUT_DIR}/all-zones-bind.txt"
echo "  cat ${OUTPUT_DIR}/all-zones-coredns.txt"
echo "  cat ${OUTPUT_DIR}/all-zones-powerdns.txt"
