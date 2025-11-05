#!/bin/bash
# Quick script to disable SSL requirement for Keycloak realms

KEYCLOAK_URL="http://localhost:14101"
ADMIN_USER="admin"
ADMIN_PASSWORD="admin"

echo "Getting admin token..."
TOKEN=$(curl -s -X POST "$KEYCLOAK_URL/realms/master/protocol/openid-connect/token" \
  -d "username=$ADMIN_USER" \
  -d "password=$ADMIN_PASSWORD" \
  -d "grant_type=password" \
  -d "client_id=admin-cli" | jq -r '.access_token')

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
    echo "Failed to get admin token"
    exit 1
fi

echo "Updating master realm..."
curl -X PUT "$KEYCLOAK_URL/admin/realms/master" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"sslRequired":"none"}'

echo -e "\nUpdating godns realm..."
curl -X PUT "$KEYCLOAK_URL/admin/realms/godns" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"sslRequired":"none"}'

echo -e "\n✅ Done! Try accessing http://localhost:14101 now"
