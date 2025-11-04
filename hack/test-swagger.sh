#!/bin/bash

# Test script for Swagger UI endpoint

API_URL="http://localhost:14082"

echo "======================================"
echo "Swagger UI Test Script"
echo "======================================"
echo ""

# Check if server is running
echo "1. Checking if API server is running..."
if curl -s "${API_URL}/health" > /dev/null 2>&1; then
    echo "✅ API server is running"
else
    echo "❌ API server is not running. Please start it first:"
    echo "   ./bin/godnsapi"
    echo ""
    exit 1
fi

echo ""

# Check Swagger JSON endpoint
echo "2. Checking Swagger JSON endpoint..."
if curl -s "${API_URL}/swagger/doc.json" | grep -q '"swagger"'; then
    echo "✅ Swagger JSON is accessible"
else
    echo "⚠️  Swagger JSON endpoint not found (this might be expected)"
fi

echo ""

# Try to access Swagger UI
echo "3. Checking Swagger UI..."
echo "   Opening browser to: ${API_URL}/swagger/index.html"
echo ""
echo "   If the browser doesn't open automatically, please visit:"
echo "   ${API_URL}/swagger/index.html"
echo ""

# Try to open in browser (works on macOS)
if [[ "$OSTYPE" == "darwin"* ]]; then
    open "${API_URL}/swagger/index.html" 2>/dev/null || echo "   (Could not auto-open browser)"
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    xdg-open "${API_URL}/swagger/index.html" 2>/dev/null || echo "   (Could not auto-open browser)"
fi

echo ""
echo "======================================"
echo "Swagger UI should now be accessible!"
echo "======================================"
echo ""
echo "You can also access:"
echo "  - Swagger JSON: ${API_URL}/swagger/doc.json"
echo "  - Health Check: ${API_URL}/health"
echo "  - API Docs: ${API_URL}/swagger/index.html"
echo ""
