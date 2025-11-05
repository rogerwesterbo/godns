# 🔐 GoDNS Authentication - Quick Reference

## 🚀 Quick Start

```bash
# 1. Start everything
docker-compose up -d

# 2. Wait for init (~30 seconds)
docker-compose logs -f keycloak-init

# 3. Login from CLI
godnscli login
# Follow the browser instructions

# 4. Use the CLI
godnscli export --format bind
```

## 📋 Default Credentials

| Service        | URL                    | Username | Password |
| -------------- | ---------------------- | -------- | -------- |
| Keycloak Admin | http://localhost:14101 | admin    | admin    |
| Test User      | -                      | testuser | password |

## 🔑 Environment Variables

```bash
# Required
KEYCLOAK_URL=http://localhost:14101
KEYCLOAK_REALM=godns
KEYCLOAK_API_CLIENT_ID=godns-api
KEYCLOAK_CLI_CLIENT_ID=godns-cli

# Optional
AUTH_ENABLED=true                    # Set to false to disable auth
TEST_USER=testuser                   # Default test username
TEST_PASSWORD=password               # Default test password
```

## 💻 CLI Commands

```bash
# Authentication
godnscli login                       # Login with device flow
godnscli logout                      # Remove cached tokens

# The export command now requires authentication
godnscli export --format bind        # Uses cached token automatically
```

## 🌐 API Usage

```bash
# Get token (password grant - for testing)
TOKEN=$(curl -s -X POST "http://localhost:14101/realms/godns/protocol/openid-connect/token" \
  -d "client_id=godns-cli" \
  -d "username=testuser" \
  -d "password=password" \
  -d "grant_type=password" | jq -r '.access_token')

# Use token
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/zones
```

## 🛠️ Common Tasks

### Disable Authentication (Development)

```bash
echo "AUTH_ENABLED=false" >> .env
docker-compose restart godns godnsapi
```

### Reset Keycloak

```bash
docker-compose down -v
docker-compose up -d
```

### View Tokens

```bash
cat ~/.godns/token.json | jq
```

### Manual Keycloak Init

```bash
./hack/init-keycloak.sh
```

## 🐛 Troubleshooting

```bash
# Check Keycloak health
curl http://localhost:14101/health

# View logs
docker-compose logs keycloak
docker-compose logs keycloak-init

# Check token validity
cat ~/.godns/token.json | jq -r '.access_token' | cut -d'.' -f2 | base64 -d 2>/dev/null | jq

# Re-login if token issues
godnscli logout
godnscli login
```

## 📁 Token Storage

- **Location**: `~/.godns/token.json`
- **Permissions**: 0600 (owner read/write only)
- **Auto-refresh**: Yes, when expired
- **Auto-cleanup**: On logout

## 🎯 Architecture

```
┌─────────┐      Device Flow      ┌──────────┐
│   CLI   │ ──────────────────────>│ Keycloak │
└─────────┘                        └──────────┘
     │                                    │
     │ Bearer Token                       │ JWT
     v                                    │
┌─────────┐      Validate JWT      ┌─────v────┐
│   API   │<───────────────────────│   OIDC   │
└─────────┘                        └──────────┘
```

## 📚 More Information

- [AUTHENTICATION.md](docs/AUTHENTICATION.md) - Complete guide
- [KEYCLOAK_SETUP.md](docs/KEYCLOAK_SETUP.md) - Keycloak details
- [AUTHENTICATION_IMPLEMENTATION.md](AUTHENTICATION_IMPLEMENTATION.md) - Technical details
