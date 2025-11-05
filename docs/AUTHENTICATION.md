# GoDNS Authentication Guide

GoDNS uses **OAuth2/OpenID Connect** authentication powered by **Keycloak** to secure both the HTTP API and CLI access.

## 🚀 Quick Start

### 1. Start the Stack with Auto-Configuration

The stack now includes automatic Keycloak initialization:

```bash
# Start all services including Keycloak
docker-compose up -d

# The keycloak-init container will automatically:
# - Create the 'godns' realm
# - Configure API and CLI clients
# - Create a test user (testuser/password)
# - Set up default roles
```

Wait for all services to be healthy (~1-2 minutes on first start).

### 2. Login from CLI

```bash
# Authenticate using OAuth2 device flow
godnscli login

# Follow the instructions:
# 1. Visit the displayed URL in your browser
# 2. Enter the code shown
# 3. Login with your Keycloak credentials
# 4. The CLI will automatically save your token
```

### 3. Use the API

```bash
# After login, all CLI commands are authenticated
godnscli export example.lan --format bind

# Your access token is automatically included
```

## 📋 Architecture

### Components

1. **Keycloak** - Identity and access management server
2. **API Server** - JWT validation middleware on all endpoints
3. **CLI** - OAuth2 device flow for user authentication
4. **Init Container** - Automated Keycloak configuration

### Authentication Flow

#### Web Application (Future)

```
User → Browser → Keycloak → Authorization Code Flow → Access Token → API
```

#### CLI Application

```
User → CLI → Device Flow → Keycloak → User Approval → Access Token → API
```

#### API Validation

```
Request → Auth Middleware → JWT Validation → Extract User Context → Route Handler
```

## 🔐 Configuration

### Environment Variables

Configure authentication in your `.env` file or environment:

```bash
# Enable/disable authentication
AUTH_ENABLED=true

# Keycloak connection
KEYCLOAK_URL=http://localhost:14101
KEYCLOAK_REALM=godns
KEYCLOAK_API_CLIENT_ID=godns-api
KEYCLOAK_CLI_CLIENT_ID=godns-cli

# Admin credentials (for initialization)
KEYCLOAK_ADMIN_USER=admin
KEYCLOAK_ADMIN_PASSWORD=admin

# Test user (created by init script)
TEST_USER=testuser
TEST_PASSWORD=password
TEST_EMAIL=testuser@godns.local
```

### Keycloak Clients

Two clients are automatically configured:

#### 1. godns-api (Bearer-Only Client)

- Used by the API server to validate tokens
- No user interaction
- Validates JWT signatures and claims

#### 2. godns-cli (Public Client)

- Used by CLI for OAuth2 device flow
- No client secret (public client)
- Supports device authorization grant
- **PKCE enabled** (S256 code challenge method) for enhanced security

## 🔑 User Management

### Default Test User

The initialization script creates a test user:

- **Username**: testuser
- **Password**: password
- **Email**: testuser@godns.local
- **Role**: dns-admin

### Creating Additional Users

#### Via Keycloak Admin Console

1. Open http://localhost:14101
2. Login as admin (admin/admin)
3. Select the `godns` realm
4. Navigate to **Users** → **Add user**
5. Fill in user details
6. Go to **Credentials** tab and set password
7. Go to **Role mapping** tab and assign roles

#### Via CLI (Keycloak Admin CLI)

```bash
# Create user
docker exec -it keycloak /opt/keycloak/bin/kcadm.sh create users \
  -r godns \
  -s username=newuser \
  -s email=newuser@example.com \
  -s enabled=true

# Set password
docker exec -it keycloak /opt/keycloak/bin/kcadm.sh set-password \
  -r godns \
  --username newuser \
  --new-password secretpassword
```

### Roles

Three roles are automatically created:

- **dns-admin** - Full DNS management access
- **dns-write** - Create and update DNS records
- **dns-read** - Read-only DNS access

## 🖥️ CLI Usage

### Login

```bash
# Start device flow authentication
godnscli login

# You'll see:
# 🔐 GoDNS Authentication
# ========================
#
# Please visit: http://localhost:14101/realms/godns/device
# And enter code: ABCD-EFGH
#
# Or visit this URL directly:
# http://localhost:14101/realms/godns/device?user_code=ABCD-EFGH
#
# Waiting for authentication...
```

Open the URL in your browser, enter the code (or use the complete URL), and login.

**Security Note**: The CLI uses **PKCE (Proof Key for Code Exchange)** with SHA256 code challenge to prevent authorization code interception attacks. This is automatically handled by the CLI - no configuration needed.

### Token Storage

Tokens are securely stored at:

- **macOS/Linux**: `~/.godns/token.json`
- **Windows**: `%USERPROFILE%\.godns\token.json`

The CLI automatically:

- Uses cached tokens when valid
- Refreshes expired tokens
- Prompts for login when refresh fails

### Logout

```bash
# Remove cached credentials
godnscli logout
```

## 🌐 API Usage

### Getting Access Tokens

#### Using Password Grant (for testing)

```bash
curl -X POST "http://localhost:14101/realms/godns/protocol/openid-connect/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "client_id=godns-cli" \
  -d "username=testuser" \
  -d "password=password" \
  -d "grant_type=password"
```

Response:

```json
{
  "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI...",
  "expires_in": 3600,
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI...",
  "token_type": "Bearer"
}
```

### Using Tokens with API

```bash
# Store token
TOKEN="eyJhbGciOiJSUzI1NiIsInR5cCI..."

# Call API with Authorization header
curl -X GET "http://localhost:8080/api/v1/zones" \
  -H "Authorization: Bearer $TOKEN"
```

### Token Validation

The API server automatically:

1. Extracts the Bearer token from `Authorization` header
2. Validates JWT signature using Keycloak's public keys
3. Checks token expiration
4. Extracts user information (ID, email, username)
5. Adds user context to request

## 🔧 Advanced Configuration

### Disable Authentication

For development or testing, you can disable authentication:

```bash
# In .env or environment
AUTH_ENABLED=false
```

The API will accept all requests without authentication.

### Custom Token Lifetime

Configure in Keycloak:

1. Realm Settings → Tokens
2. Adjust **Access Token Lifespan** (default: 3600s / 1 hour)
3. Adjust **SSO Session Idle** (default: 1800s / 30 min)

### Social Login Providers

Add GitHub, Google, Facebook, etc.:

1. Keycloak Admin → Identity Providers
2. Select provider (GitHub, Google, etc.)
3. Configure client ID and secret
4. Users can now login with social accounts

### RBAC (Role-Based Access Control)

The authentication middleware extracts user roles. You can implement RBAC in handlers:

```go
userID, username, email := middleware.GetUserFromContext(r.Context())
// Check roles or permissions
// Implement authorization logic
```

## 🔒 Security Best Practices

### Production Checklist

- [ ] Change default admin password
- [ ] Use HTTPS for Keycloak (behind reverse proxy)
- [ ] Use managed PostgreSQL for Keycloak database
- [ ] Enable MFA for admin accounts
- [ ] Implement rate limiting
- [ ] Monitor authentication logs
- [ ] Rotate secrets regularly
- [ ] Use strong passwords
- [ ] Implement session timeout
- [ ] Review and audit permissions

### HTTPS Configuration

```yaml
# docker-compose.yaml
keycloak:
  environment:
    KC_HOSTNAME: auth.yourdomain.com
    KC_HOSTNAME_STRICT: "true"
    KC_HOSTNAME_STRICT_HTTPS: "true"
    KC_PROXY: edge
```

Set up nginx or Traefik for SSL termination.

## 🐛 Troubleshooting

### CLI Login Issues

**Problem**: "Failed to connect to Keycloak"

```bash
# Check Keycloak is running
docker-compose ps keycloak

# Check Keycloak logs
docker-compose logs keycloak

# Verify Keycloak URL
echo $KEYCLOAK_URL

# Test connectivity
curl http://localhost:14101/health
```

**Problem**: "Token verification failed"

- Check system time is synchronized
- Verify Keycloak realm and client configuration
- Check token hasn't expired

### API Authentication Issues

**Problem**: "Invalid or expired token"

- Check token expiration: decode JWT at https://jwt.io
- Verify `KEYCLOAK_URL` matches in API and Keycloak
- Check realm and client ID configuration

**Problem**: "Missing authorization header"

```bash
# Ensure Authorization header is set
curl -v -H "Authorization: Bearer TOKEN" ...
```

### Keycloak Init Issues

**Problem**: Init container fails

```bash
# Check init container logs
docker-compose logs keycloak-init

# Manually run init script
./hack/init-keycloak.sh

# Verify Keycloak is ready
curl http://localhost:14101/health/ready
```

## 📚 References

- [Keycloak Documentation](https://www.keycloak.org/documentation)
- [OAuth 2.0 Device Flow](https://oauth.net/2/device-flow/)
- [PKCE (RFC 7636)](https://datatracker.ietf.org/doc/html/rfc7636) - Proof Key for Code Exchange
- [OpenID Connect Specification](https://openid.net/connect/)
- [JWT.io](https://jwt.io/) - JWT decoder

## 🔄 Migration from Unauthenticated API

If you had an existing unauthenticated setup:

1. **Update environment variables** - Add Keycloak configuration
2. **Start Keycloak services** - `docker-compose up -d`
3. **Login from CLI** - `godnscli login`
4. **Update API clients** - Add Authorization header to requests

Authentication is enabled by default. To temporarily disable:

```bash
AUTH_ENABLED=false docker-compose up -d
```

## 💡 Examples

### Complete Authentication Flow

```bash
# 1. Start services
docker-compose up -d

# 2. Wait for initialization
docker-compose logs -f keycloak-init

# 3. Login from CLI
godnscli login

# 4. Use authenticated commands
godnscli export example.lan --format bind

# 5. Test with curl
TOKEN=$(cat ~/.godns/token.json | jq -r '.access_token')
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/zones
```

### Manual Token Testing

```bash
#!/bin/bash

# Get token
RESPONSE=$(curl -s -X POST "http://localhost:14101/realms/godns/protocol/openid-connect/token" \
  -d "client_id=godns-cli" \
  -d "username=testuser" \
  -d "password=password" \
  -d "grant_type=password")

TOKEN=$(echo $RESPONSE | jq -r '.access_token')

# Use token
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/zones | jq
```
