#!/bin/bash

# Initialize Keycloak web client for GoDNS Web Application
# This script creates the godns-web client in Keycloak with proper PKCE configuration

set -e

KEYCLOAK_URL="${KEYCLOAK_URL:-http://localhost:14101}"
REALM="${REALM:-godns}"
ADMIN_USER="${ADMIN_USER:-admin}"
ADMIN_PASSWORD="${ADMIN_PASSWORD:-admin}"
CLIENT_ID="godns-web"
REDIRECT_URI="http://localhost:14200/callback"
WEB_ORIGIN="http://localhost:14200"

echo "🔧 Setting up Keycloak Web Client..."
echo "Keycloak URL: $KEYCLOAK_URL"
echo "Realm: $REALM"
echo "Client ID: $CLIENT_ID"
echo ""

# Get admin access token
echo "📝 Getting admin access token..."
ADMIN_TOKEN=$(curl -s -X POST "$KEYCLOAK_URL/realms/master/protocol/openid-connect/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$ADMIN_USER" \
  -d "password=$ADMIN_PASSWORD" \
  -d "grant_type=password" \
  -d "client_id=admin-cli" | jq -r '.access_token')

if [ -z "$ADMIN_TOKEN" ] || [ "$ADMIN_TOKEN" == "null" ]; then
  echo "❌ Failed to get admin token. Is Keycloak running?"
  exit 1
fi

echo "✅ Admin token obtained"

# Check if client already exists
echo "🔍 Checking if client already exists..."
EXISTING_CLIENT=$(curl -s -X GET "$KEYCLOAK_URL/admin/realms/$REALM/clients?clientId=$CLIENT_ID" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json")

if [ "$(echo "$EXISTING_CLIENT" | jq '. | length')" -gt 0 ]; then
  echo "⚠️  Client $CLIENT_ID already exists. Deleting..."
  CLIENT_UUID=$(echo "$EXISTING_CLIENT" | jq -r '.[0].id')
  curl -s -X DELETE "$KEYCLOAK_URL/admin/realms/$REALM/clients/$CLIENT_UUID" \
    -H "Authorization: Bearer $ADMIN_TOKEN"
  echo "✅ Existing client deleted"
fi

# Create the web client
echo "🚀 Creating web client..."
curl -s -X POST "$KEYCLOAK_URL/admin/realms/$REALM/clients" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "clientId": "'"$CLIENT_ID"'",
    "name": "GoDNS Web Application",
    "description": "Web frontend for GoDNS DNS Management",
    "enabled": true,
    "publicClient": true,
    "protocol": "openid-connect",
    "standardFlowEnabled": true,
    "directAccessGrantsEnabled": false,
    "implicitFlowEnabled": false,
    "serviceAccountsEnabled": false,
    "redirectUris": [
      "'"$REDIRECT_URI"'",
      "http://localhost:14200/*"
    ],
    "webOrigins": [
      "'"$WEB_ORIGIN"'"
    ],
    "rootUrl": "'"$WEB_ORIGIN"'",
    "baseUrl": "'"$WEB_ORIGIN"'",
    "adminUrl": "",
    "attributes": {
      "pkce.code.challenge.method": "S256",
      "post.logout.redirect.uris": "'"$WEB_ORIGIN"'",
      "access.token.lifespan": "3600"
    }
  }'

echo ""
echo "✅ Web client created successfully!"
echo ""
echo "Client configuration:"
echo "  Client ID: $CLIENT_ID"
echo "  Client Type: Public (PKCE required)"
echo "  Redirect URI: $REDIRECT_URI"
echo "  Web Origins: $WEB_ORIGIN"
echo "  PKCE Method: S256"
echo ""
echo "🎉 You can now login to the web application at http://localhost:14200"
echo ""
