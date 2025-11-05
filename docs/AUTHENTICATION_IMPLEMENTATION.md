# 🔐 OAuth2/OIDC Authentication Implementation

This implementation adds comprehensive authentication to GoDNS using **OAuth2/OpenID Connect** with **Keycloak** as the identity provider.

## ✨ Features Implemented

### 1. Keycloak Integration

- ✅ Docker Compose configuration with Keycloak service
- ✅ PostgreSQL backend for Keycloak data
- ✅ Automated initialization script
- ✅ Init container for zero-touch setup
- ✅ Realm, client, user, and role auto-configuration

### 2. API Authentication

- ✅ JWT validation middleware
- ✅ OIDC provider integration
- ✅ Bearer token authentication on all endpoints
- ✅ User context extraction (ID, email, username)
- ✅ Configurable enable/disable via `AUTH_ENABLED`

### 3. CLI Authentication

- ✅ OAuth2 device flow implementation
- ✅ `godnscli login` command
- ✅ `godnscli logout` command
- ✅ Automatic token caching (~/.godns/token.json)
- ✅ Automatic token refresh
- ✅ Authenticated HTTP client for all API calls

### 4. Documentation

- ✅ Complete authentication guide (AUTHENTICATION.md)
- ✅ Updated Keycloak setup documentation
- ✅ Troubleshooting section
- ✅ Security best practices

## 📦 New Files Created

```
hack/
  init-keycloak.sh                    # Keycloak auto-config script
Dockerfile.keycloak-init              # Init container Dockerfile
pkg/
  auth/
    auth.go                           # OAuth2 device flow & token management
  httpclient/
    client.go                         # Authenticated HTTP client for CLI
internal/
  httpserver/
    middleware/
      auth.go                         # JWT validation middleware
cmd/
  godnscli/
    cmd/
      login.go                        # CLI login command
      logout.go                       # CLI logout command
docs/
  AUTHENTICATION.md                   # Complete auth guide
```

## 🔧 Modified Files

```
docker-compose.yaml                   # Added keycloak-init service
go.mod                                # Added OIDC dependencies
pkg/consts/consts.go                  # Auth configuration constants
internal/settings/settings.go         # Auth default settings
cmd/godnscli/settings/cli_settings.go # CLI auth settings
internal/httpserver/http_server.go    # Auth middleware integration
internal/httpserver/httproutes/http_routes.go  # Apply auth to routes
cmd/godns/main.go                     # Handle auth middleware errors
cmd/godnsapi/main.go                  # Handle auth middleware errors
cmd/godnscli/cmd/export.go            # Use authenticated HTTP client
```

## 🚀 Usage

### Start the Stack

```bash
# Start all services (Keycloak will auto-configure)
docker-compose up -d

# Wait for initialization
docker-compose logs -f keycloak-init

# Should see:
# ✅ Keycloak initialization complete!
```

### CLI Authentication

```bash
# Login using device flow
godnscli login

# Follow the instructions to authenticate in browser
# Token is automatically saved and reused

# Use CLI commands (now authenticated)
godnscli export example.lan --format bind

# Logout when done
godnscli logout
```

### API Testing

```bash
# Get a token manually
curl -X POST "http://localhost:14101/realms/godns/protocol/openid-connect/token" \
  -d "client_id=godns-cli" \
  -d "username=testuser" \
  -d "password=password" \
  -d "grant_type=password"

# Use the token
curl -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/v1/zones
```

## 🔑 Default Credentials

**Keycloak Admin:**

- URL: http://localhost:14101
- Username: admin
- Password: admin

**Test User:**

- Username: testuser
- Password: password
- Role: dns-admin

## 🎯 Authentication Flow

### CLI Device Flow

1. User runs `godnscli login`
2. CLI requests device code from Keycloak
3. User visits URL and enters code in browser
4. User authenticates with Keycloak
5. CLI polls for token
6. Token is saved to `~/.godns/token.json`
7. Subsequent commands use cached token
8. Token auto-refreshes when expired

### API Request Flow

1. Client sends request with `Authorization: Bearer <token>`
2. Auth middleware extracts token
3. Middleware validates JWT with Keycloak JWKS
4. User info extracted and added to context
5. Request proceeds to handler

## 🔒 Security Features

- ✅ JWT signature validation
- ✅ Token expiration checking
- ✅ Refresh token support
- ✅ Secure token storage (0600 permissions)
- ✅ OAuth2 device flow (no client secrets)
- ✅ Role-based access control ready
- ✅ Configurable token lifetimes

## 🎨 Configuration

### Environment Variables

```bash
# Enable/disable authentication
AUTH_ENABLED=true

# Keycloak configuration
KEYCLOAK_URL=http://localhost:14101
KEYCLOAK_REALM=godns
KEYCLOAK_API_CLIENT_ID=godns-api
KEYCLOAK_CLI_CLIENT_ID=godns-cli

# Test user credentials
TEST_USER=testuser
TEST_PASSWORD=password
```

### Disable Authentication (Development)

```bash
# Disable for testing
AUTH_ENABLED=false

# Or in .env file
echo "AUTH_ENABLED=false" >> .env
docker-compose restart godns
```

## 🧪 Testing

### Manual Testing

```bash
# 1. Start services
docker-compose up -d

# 2. Login
godnscli login

# 3. Test authenticated request
godnscli export --format bind

# 4. Test without auth (should fail if AUTH_ENABLED=true)
curl http://localhost:8080/api/v1/zones
# Expected: {"error": "missing authorization header"}

# 5. Test with auth
TOKEN=$(cat ~/.godns/token.json | jq -r '.access_token')
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/zones
# Expected: Success (zone list)
```

## 📝 Next Steps

### For Future Development

1. **Web Application Login**

   - Implement authorization code flow
   - Add session management
   - Create login/logout pages

2. **Role-Based Access Control**

   - Implement role checking in handlers
   - Define granular permissions
   - Add role-specific endpoints

3. **Social Login**

   - Configure GitHub provider
   - Configure Google provider
   - Add other OAuth providers

4. **Advanced Features**
   - Multi-factor authentication
   - Account recovery
   - Password policies
   - Audit logging

## 🐛 Troubleshooting

### Keycloak won't start

```bash
docker-compose logs keycloak
# Check PostgreSQL is healthy
docker-compose ps postgres
```

### Init script fails

```bash
# Run manually to see detailed output
./hack/init-keycloak.sh
```

### CLI login fails

```bash
# Check Keycloak URL
echo $KEYCLOAK_URL

# Verify connectivity
curl http://localhost:14101/health

# Check realm exists
curl http://localhost:14101/realms/godns
```

### Token validation fails

```bash
# Verify API can reach Keycloak
docker exec godns curl http://keycloak:8080/health

# Check token at jwt.io
cat ~/.godns/token.json | jq -r '.access_token'
```

## 📚 Documentation

- **[AUTHENTICATION.md](docs/AUTHENTICATION.md)** - Complete authentication guide
- **[KEYCLOAK_SETUP.md](docs/KEYCLOAK_SETUP.md)** - Keycloak configuration details

## 🙏 Dependencies Added

```go
github.com/coreos/go-oidc/v3 v3.16.0      // OIDC client library
golang.org/x/oauth2                        // OAuth2 client library
```

## ✅ Implementation Checklist

- [x] Keycloak Docker setup
- [x] Auto-initialization script
- [x] Init container
- [x] JWT validation middleware
- [x] OAuth2 device flow
- [x] Token caching
- [x] Token refresh
- [x] CLI login/logout commands
- [x] Authenticated HTTP client
- [x] Configuration management
- [x] Documentation
- [x] Testing

## 🎉 Ready to Use!

The authentication system is fully implemented and ready for use. Simply start the stack and run `godnscli login` to get started!
