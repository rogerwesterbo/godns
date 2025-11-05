#!/bin/bash
set -e

# Keycloak initialization script
# This script automatically configures Keycloak with:
# - GoDNS realm
# - API and CLI clients
# - Test user
# - Roles

KEYCLOAK_URL="${KEYCLOAK_URL:-https://keycloak:8443}"
ADMIN_USER="${KEYCLOAK_ADMIN:-admin}"
ADMIN_PASSWORD="${KEYCLOAK_ADMIN_PASSWORD:-admin}"
REALM_NAME="${KEYCLOAK_REALM:-godns}"
API_CLIENT_ID="${KEYCLOAK_API_CLIENT_ID:-godns-api}"
CLI_CLIENT_ID="${KEYCLOAK_CLI_CLIENT_ID:-godns-cli}"

echo "🔧 Initializing Keycloak for GoDNS..."
echo "Keycloak URL: $KEYCLOAK_URL"
echo "Realm: $REALM_NAME"

# Wait for Keycloak to be ready
echo "⏳ Waiting for Keycloak to be ready..."
MAX_RETRIES=30
RETRY_COUNT=0
until curl -sfk "$KEYCLOAK_URL/realms/master" > /dev/null 2>&1; do
    RETRY_COUNT=$((RETRY_COUNT + 1))
    if [ $RETRY_COUNT -ge $MAX_RETRIES ]; then
        echo "❌ Keycloak not ready after $MAX_RETRIES attempts"
        exit 1
    fi
    echo "Waiting... (attempt $RETRY_COUNT/$MAX_RETRIES)"
    sleep 2
done
echo "✅ Keycloak is ready"

# Get admin access token
echo "🔑 Getting admin access token..."
ADMIN_TOKEN=$(curl -sf -X POST "$KEYCLOAK_URL/realms/master/protocol/openid-connect/token" \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "username=$ADMIN_USER" \
    -d "password=$ADMIN_PASSWORD" \
    -d "grant_type=password" \
    -d "client_id=admin-cli" | jq -r '.access_token')

if [ -z "$ADMIN_TOKEN" ] || [ "$ADMIN_TOKEN" = "null" ]; then
    echo "❌ Failed to get admin token"
    exit 1
fi
echo "✅ Admin token obtained"

# Check if realm already exists
echo "🔍 Checking if realm '$REALM_NAME' exists..."
REALM_EXISTS=$(curl -sf -X GET "$KEYCLOAK_URL/admin/realms/$REALM_NAME" \
    -H "Authorization: Bearer $ADMIN_TOKEN" > /dev/null 2>&1 && echo "true" || echo "false")

if [ "$REALM_EXISTS" = "false" ]; then
    echo "📦 Creating realm '$REALM_NAME'..."
    curl -sf -X POST "$KEYCLOAK_URL/admin/realms" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d "{
            \"realm\": \"$REALM_NAME\",
            \"enabled\": true,
            \"sslRequired\": \"none\",
            \"displayName\": \"GoDNS\",
            \"displayNameHtml\": \"<b>GoDNS</b> DNS Management\",
            \"registrationAllowed\": false,
            \"loginWithEmailAllowed\": true,
            \"duplicateEmailsAllowed\": false,
            \"resetPasswordAllowed\": true,
            \"editUsernameAllowed\": false,
            \"bruteForceProtected\": true,
            \"accessTokenLifespan\": 3600,
            \"ssoSessionIdleTimeout\": 1800,
            \"ssoSessionMaxLifespan\": 36000
        }"
    echo "✅ Realm created"
else
    echo "ℹ️  Realm '$REALM_NAME' already exists"
    # Update SSL requirement to none for development
    echo "🔧 Updating realm SSL requirement..."
    curl -sf -X PUT "$KEYCLOAK_URL/admin/realms/$REALM_NAME" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d "{
            \"realm\": \"$REALM_NAME\",
            \"sslRequired\": \"none\"
        }"
    echo "✅ Realm updated"
fi

# Create or update API client
echo "🔧 Configuring API client '$API_CLIENT_ID'..."
API_CLIENT_JSON=$(cat <<EOF
{
    "clientId": "$API_CLIENT_ID",
    "name": "GoDNS API",
    "description": "GoDNS HTTP API Server",
    "enabled": true,
    "protocol": "openid-connect",
    "publicClient": false,
    "bearerOnly": true,
    "standardFlowEnabled": false,
    "directAccessGrantsEnabled": false,
    "serviceAccountsEnabled": false,
    "attributes": {
        "access.token.lifespan": "3600"
    }
}
EOF
)

# Check if client exists
API_CLIENT_ID_INTERNAL=$(curl -sf -X GET "$KEYCLOAK_URL/admin/realms/$REALM_NAME/clients?clientId=$API_CLIENT_ID" \
    -H "Authorization: Bearer $ADMIN_TOKEN" | jq -r '.[0].id // empty')

if [ -z "$API_CLIENT_ID_INTERNAL" ]; then
    echo "📦 Creating API client..."
    curl -sf -X POST "$KEYCLOAK_URL/admin/realms/$REALM_NAME/clients" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d "$API_CLIENT_JSON"
    echo "✅ API client created"
else
    echo "ℹ️  API client already exists"
fi

# Create or update CLI client (public client with device flow)
echo "🔧 Configuring CLI client '$CLI_CLIENT_ID'..."
CLI_CLIENT_JSON=$(cat <<EOF
{
    "clientId": "$CLI_CLIENT_ID",
    "name": "GoDNS CLI",
    "description": "GoDNS Command Line Interface",
    "enabled": true,
    "protocol": "openid-connect",
    "publicClient": true,
    "bearerOnly": false,
    "standardFlowEnabled": true,
    "directAccessGrantsEnabled": true,
    "serviceAccountsEnabled": false,
    "redirectUris": ["http://localhost:*", "urn:ietf:wg:oauth:2.0:oob"],
    "webOrigins": ["+"],
    "attributes": {
        "oauth2.device.authorization.grant.enabled": "true",
        "oidc.ciba.grant.enabled": "false"
    },
    "protocolMappers": [
        {
            "name": "audience-mapper",
            "protocol": "openid-connect",
            "protocolMapper": "oidc-audience-mapper",
            "config": {
                "included.client.audience": "$API_CLIENT_ID",
                "id.token.claim": "false",
                "access.token.claim": "true"
            }
        }
    ]
}
EOF
)

CLI_CLIENT_ID_INTERNAL=$(curl -sf -X GET "$KEYCLOAK_URL/admin/realms/$REALM_NAME/clients?clientId=$CLI_CLIENT_ID" \
    -H "Authorization: Bearer $ADMIN_TOKEN" | jq -r '.[0].id // empty')

if [ -z "$CLI_CLIENT_ID_INTERNAL" ]; then
    echo "📦 Creating CLI client..."
    curl -sf -X POST "$KEYCLOAK_URL/admin/realms/$REALM_NAME/clients" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d "$CLI_CLIENT_JSON"
    echo "✅ CLI client created"
    # Get the client ID after creation
    CLI_CLIENT_ID_INTERNAL=$(curl -sf -X GET "$KEYCLOAK_URL/admin/realms/$REALM_NAME/clients?clientId=$CLI_CLIENT_ID" \
        -H "Authorization: Bearer $ADMIN_TOKEN" | jq -r '.[0].id')
else
    echo "🔧 Updating CLI client with audience mapper..."
    curl -sf -X PUT "$KEYCLOAK_URL/admin/realms/$REALM_NAME/clients/$CLI_CLIENT_ID_INTERNAL" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d "$CLI_CLIENT_JSON"
    echo "✅ CLI client updated"
fi

# Remove default 'account' client scope from CLI client to prevent 'account' audience
echo "🔧 Removing 'account' client scope from CLI client..."
# Get account scope ID
ACCOUNT_SCOPE_ID=$(curl -sf -X GET "$KEYCLOAK_URL/admin/realms/$REALM_NAME/client-scopes" \
    -H "Authorization: Bearer $ADMIN_TOKEN" | jq -r '.[] | select(.name=="account") | .id // empty')

if [ -n "$ACCOUNT_SCOPE_ID" ]; then
    echo "Found account scope ID: $ACCOUNT_SCOPE_ID"
    # Remove from default scopes
    DELETE_RESULT=$(curl -s -w "%{http_code}" -o /dev/null -X DELETE "$KEYCLOAK_URL/admin/realms/$REALM_NAME/clients/$CLI_CLIENT_ID_INTERNAL/default-client-scopes/$ACCOUNT_SCOPE_ID" \
        -H "Authorization: Bearer $ADMIN_TOKEN")
    if [ "$DELETE_RESULT" = "204" ] || [ "$DELETE_RESULT" = "404" ]; then
        echo "✅ Removed 'account' from default scopes (or was not present)"
    else
        echo "⚠️  Failed to remove 'account' from default scopes (HTTP $DELETE_RESULT)"
    fi
    
    # Also try removing from optional scopes
    DELETE_RESULT=$(curl -s -w "%{http_code}" -o /dev/null -X DELETE "$KEYCLOAK_URL/admin/realms/$REALM_NAME/clients/$CLI_CLIENT_ID_INTERNAL/optional-client-scopes/$ACCOUNT_SCOPE_ID" \
        -H "Authorization: Bearer $ADMIN_TOKEN")
    if [ "$DELETE_RESULT" = "204" ] || [ "$DELETE_RESULT" = "404" ]; then
        echo "✅ Removed 'account' from optional scopes (or was not present)"
    fi
else
    echo "ℹ️  'account' scope not found in realm"
fi

# Create roles
echo "🔧 Creating roles..."
ROLES=("dns-admin" "dns-write" "dns-read")
ROLE_DESCRIPTIONS=(
    "Full DNS management access"
    "Create and update DNS records"
    "Read-only DNS access"
)

for i in "${!ROLES[@]}"; do
    ROLE_NAME="${ROLES[$i]}"
    ROLE_DESC="${ROLE_DESCRIPTIONS[$i]}"
    
    ROLE_EXISTS=$(curl -sf -X GET "$KEYCLOAK_URL/admin/realms/$REALM_NAME/roles/$ROLE_NAME" \
        -H "Authorization: Bearer $ADMIN_TOKEN" > /dev/null 2>&1 && echo "true" || echo "false")
    
    if [ "$ROLE_EXISTS" = "false" ]; then
        echo "📦 Creating role '$ROLE_NAME'..."
        curl -sf -X POST "$KEYCLOAK_URL/admin/realms/$REALM_NAME/roles" \
            -H "Authorization: Bearer $ADMIN_TOKEN" \
            -H "Content-Type: application/json" \
            -d "{
                \"name\": \"$ROLE_NAME\",
                \"description\": \"$ROLE_DESC\"
            }"
        echo "✅ Role '$ROLE_NAME' created"
    else
        echo "ℹ️  Role '$ROLE_NAME' already exists"
    fi
done

# Create test user
TEST_USERNAME="${TEST_USER:-testuser}"
TEST_PASSWORD="${TEST_PASSWORD:-password}"
TEST_EMAIL="${TEST_EMAIL:-testuser@godns.local}"

echo "🔧 Creating test user '$TEST_USERNAME'..."
USER_EXISTS=$(curl -sf -X GET "$KEYCLOAK_URL/admin/realms/$REALM_NAME/users?username=$TEST_USERNAME" \
    -H "Authorization: Bearer $ADMIN_TOKEN" | jq -r '.[0].id // empty')

if [ -z "$USER_EXISTS" ]; then
    echo "📦 Creating user..."
    USER_ID=$(curl -sf -X POST "$KEYCLOAK_URL/admin/realms/$REALM_NAME/users" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d "{
            \"username\": \"$TEST_USERNAME\",
            \"email\": \"$TEST_EMAIL\",
            \"firstName\": \"Test\",
            \"lastName\": \"User\",
            \"enabled\": true,
            \"emailVerified\": true,
            \"credentials\": [{
                \"type\": \"password\",
                \"value\": \"$TEST_PASSWORD\",
                \"temporary\": false
            }]
        }" -s -D - -o /dev/null | grep -i "^location:" | sed 's/.*\///g' | tr -d '\r\n')
    
    # Get user ID if creation succeeded
    if [ -z "$USER_ID" ]; then
        USER_ID=$(curl -sf -X GET "$KEYCLOAK_URL/admin/realms/$REALM_NAME/users?username=$TEST_USERNAME" \
            -H "Authorization: Bearer $ADMIN_TOKEN" | jq -r '.[0].id')
    fi
    
    echo "✅ User '$TEST_USERNAME' created with ID: $USER_ID"
    
    # Assign dns-admin role to test user
    if [ -n "$USER_ID" ]; then
        echo "🔧 Assigning dns-admin role to user..."
        ROLE_REPRESENTATION=$(curl -sf -X GET "$KEYCLOAK_URL/admin/realms/$REALM_NAME/roles/dns-admin" \
            -H "Authorization: Bearer $ADMIN_TOKEN")
        
        curl -sf -X POST "$KEYCLOAK_URL/admin/realms/$REALM_NAME/users/$USER_ID/role-mappings/realm" \
            -H "Authorization: Bearer $ADMIN_TOKEN" \
            -H "Content-Type: application/json" \
            -d "[$ROLE_REPRESENTATION]"
        echo "✅ Role assigned"
    fi
else
    echo "ℹ️  User '$TEST_USERNAME' already exists"
fi

echo ""
echo "🎉 Keycloak initialization complete!"
echo ""
echo "📋 Configuration Summary:"
echo "  Realm: $REALM_NAME"
echo "  API Client: $API_CLIENT_ID (bearer-only)"
echo "  CLI Client: $CLI_CLIENT_ID (public, device flow enabled)"
echo "  Test User: $TEST_USERNAME"
echo "  Test Password: $TEST_PASSWORD"
echo ""
echo "🔐 Test authentication:"
echo "  curl -X POST '$KEYCLOAK_URL/realms/$REALM_NAME/protocol/openid-connect/token' \\"
echo "    -d 'client_id=$CLI_CLIENT_ID' \\"
echo "    -d 'username=$TEST_USERNAME' \\"
echo "    -d 'password=$TEST_PASSWORD' \\"
echo "    -d 'grant_type=password'"
echo ""
