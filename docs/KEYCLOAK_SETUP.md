# Keycloak Setup for GoDNS Authentication

This guide explains how to set up and configure Keycloak for authentication with the GoDNS API.

## Quick Start

### 1. Start Keycloak

Start Keycloak and its database using docker-compose:

```bash
docker-compose up -d postgres keycloak
```

Wait for Keycloak to be ready (about 1-2 minutes on first start):

```bash
# Check logs
docker-compose logs -f keycloak

# Wait for: "Keycloak ... started"
```

### 2. Access Keycloak Admin Console

Open your browser and navigate to:

```
http://localhost:14083
```

**Default Admin Credentials:**

- Username: `admin`
- Password: `admin`

> ⚠️ **Security Warning**: Change the default admin password in production!

## Configuration

### Environment Variables

Configure Keycloak using environment variables in your `.env` file:

```bash
# Keycloak Port
KEYCLOAK_PORT=8080

# Admin credentials
KEYCLOAK_ADMIN_USER=admin
KEYCLOAK_ADMIN_PASSWORD=admin

# Hostname configuration
KEYCLOAK_HOSTNAME=localhost
KEYCLOAK_PROXY=edge

# Database configuration
KEYCLOAK_DB_NAME=keycloak
KEYCLOAK_DB_USER=keycloak
KEYCLOAK_DB_PASSWORD=keycloak_password
```

### Port Configuration

If port 14083 conflicts with other services, you can change Keycloak's port:

```bash
# .env file
KEYCLOAK_PORT=14084
```

Then access Keycloak at `http://localhost:14084`

## Setting Up GoDNS Realm and Client

### 1. Create a Realm

1. Log in to Keycloak Admin Console
2. Click on the dropdown in the top-left (shows "master")
3. Click **"Create Realm"**
4. Enter realm name: `godns`
5. Click **"Create"**

### 2. Create a Client

1. In the `godns` realm, go to **Clients** → **Create client**
2. Configure the client:

   - **Client type**: OpenID Connect
   - **Client ID**: `godns-api`
   - Click **Next**

3. **Capability config**:

   - ✅ Client authentication: ON
   - ✅ Authorization: ON (if you need fine-grained permissions)
   - ✅ Standard flow: ON
   - ✅ Direct access grants: ON
   - Click **Next**

4. **Login settings**:

   - Valid redirect URIs: `http://localhost:14082/*`
   - Valid post logout redirect URIs: `http://localhost:14082/*`
   - Web origins: `http://localhost:14082`
   - Click **Save**

5. Get the **Client Secret**:
   - Go to **Clients** → `godns-api` → **Credentials** tab
   - Copy the **Client Secret** (you'll need this for API configuration)

### 3. Create Users

1. Go to **Users** → **Add user**
2. Fill in user details:

   - **Username**: `testuser`
   - **Email**: `test@example.com`
   - **First name**: `Test`
   - **Last name**: `User`
   - Click **Create**

3. Set password:
   - Go to **Credentials** tab
   - Click **Set password**
   - Enter password (e.g., `password`)
   - Toggle **Temporary** to OFF (so user doesn't need to change it)
   - Click **Save**

### 4. Create Roles (Optional)

Create roles for managing permissions:

1. Go to **Realm roles** → **Create role**
2. Create roles like:

   - `dns-admin` - Full DNS management access
   - `dns-read` - Read-only DNS access
   - `dns-write` - Create/update DNS records

3. Assign roles to users:
   - Go to **Users** → Select user → **Role mapping**
   - Click **Assign role**
   - Select the roles and click **Assign**

## Integration with GoDNS API

### Getting Access Tokens

#### Using Password Grant (Direct Access)

```bash
# Get access token
curl -X POST "http://localhost:14083/realms/godns/protocol/openid-connect/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "client_id=godns-api" \
  -d "client_secret=YOUR_CLIENT_SECRET" \
  -d "username=testuser" \
  -d "password=password" \
  -d "grant_type=password"
```

Response:

```json
{
  "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI...",
  "expires_in": 300,
  "refresh_expires_in": 1800,
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI...",
  "token_type": "Bearer"
}
```

#### Using the Access Token

```bash
# Store the token
TOKEN="eyJhbGciOiJSUzI1NiIsInR5cCI..."

# Use it to call GoDNS API
curl -X GET "http://localhost:14082/api/v1/zones" \
  -H "Authorization: Bearer $TOKEN"
```

### Token Validation

To validate tokens in your GoDNS API, you'll need:

1. **Keycloak URL**: `http://keycloak:14083` (internal Docker network) or `http://localhost:14083` (from host)
2. **Realm**: `godns`
3. **Client ID**: `godns-api`
4. **Client Secret**: (from Keycloak admin console)

### JWKS Endpoint

Get public keys for token verification:

```bash
curl http://localhost:14083/realms/godns/protocol/openid-connect/certs
```

## Production Considerations

### 1. Use HTTPS

In production, always use HTTPS:

```yaml
# docker-compose.yaml
environment:
  KC_HOSTNAME_STRICT_HTTPS: "true"
  KC_PROXY: edge
```

Configure a reverse proxy (nginx, traefik) to handle SSL termination.

### 2. Use Production Database

Replace the PostgreSQL container with a managed database service:

```yaml
environment:
  KC_DB_URL: jdbc:postgresql://your-db-host:5432/keycloak
```

### 3. Change Admin Password

```bash
# Set strong admin password
KEYCLOAK_ADMIN_PASSWORD=<strong-random-password>
```

### 4. Use Production Mode

Replace `start-dev` with production configuration:

```yaml
command:
  - start
  - --optimized
```

### 5. Configure Hostname

```yaml
environment:
  KC_HOSTNAME: auth.yourdomain.com
  KC_HOSTNAME_STRICT: "true"
```

## Docker Compose Management

### Start Services

```bash
# Start all services
docker-compose up -d

# Start only Keycloak services
docker-compose up -d postgres keycloak
```

### View Logs

```bash
# Keycloak logs
docker-compose logs -f keycloak

# Database logs
docker-compose logs -f keycloak-db
```

### Stop Services

```bash
# Stop all
docker-compose down

# Stop but keep data
docker-compose stop

# Stop and remove volumes (⚠️ deletes all data)
docker-compose down -v
```

### Restart Keycloak

```bash
docker-compose restart keycloak
```

## Troubleshooting

### Keycloak Won't Start

1. Check database is healthy:

   ```bash
   docker-compose ps
   docker-compose logs postgres
   ```

2. Check Keycloak logs:

   ```bash
   docker-compose logs keycloak
   ```

3. Ensure ports aren't in use:
   ```bash
   lsof -i :14083
   ```

### Can't Access Admin Console

1. Verify Keycloak is running:

   ```bash
   docker-compose ps keycloak
   ```

2. Check health endpoint:

   ```bash
   curl http://localhost:14083/health
   ```

3. Verify port mapping:
   ```bash
   docker-compose ps
   ```

### Database Connection Issues

1. Check database is running:

   ```bash
   docker-compose ps postgres
   ```

2. Verify connection from Keycloak container:
   ```bash
   docker exec -it keycloak bash
   pg_isready -h postgres -U keycloak
   ```

### Reset Everything

To start fresh (⚠️ deletes all Keycloak data):

```bash
docker-compose down -v
docker-compose up -d keycloak-db keycloak
```

## Backup and Restore

### Backup Keycloak Data

```bash
# Export realm configuration
docker exec -it keycloak /opt/keycloak/bin/kc.sh export \
  --dir /tmp/export \
  --realm godns

# Copy to host
docker cp keycloak:/tmp/export ./keycloak-backup
```

### Backup Database

```bash
docker exec postgres pg_dump -U keycloak keycloak > keycloak-db-backup.sql
```

### Restore Database

```bash
cat keycloak-db-backup.sql | docker exec -i postgres psql -U keycloak
```

## Next Steps

1. **Integrate with GoDNS API**: Add JWT validation middleware to your API endpoints
2. **Configure RBAC**: Use Keycloak roles to control access to DNS operations
3. **Set up SSO**: Configure other applications to use Keycloak for single sign-on
4. **Enable MFA**: Configure multi-factor authentication for enhanced security

## Useful Resources

- [Keycloak Documentation](https://www.keycloak.org/documentation)
- [Keycloak Server Administration Guide](https://www.keycloak.org/docs/latest/server_admin/)
- [Securing Applications and Services Guide](https://www.keycloak.org/docs/latest/securing_apps/)
- [OpenID Connect Specification](https://openid.net/connect/)

## Example: Complete Authentication Flow

Here's a complete example of authenticating and using the GoDNS API:

```bash
#!/bin/bash

# Configuration
KEYCLOAK_URL="http://localhost:14083"
REALM="godns"
CLIENT_ID="godns-api"
CLIENT_SECRET="your-client-secret"
USERNAME="testuser"
PASSWORD="password"
API_URL="http://localhost:14082"

# Get access token
TOKEN_RESPONSE=$(curl -s -X POST "${KEYCLOAK_URL}/realms/${REALM}/protocol/openid-connect/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "client_id=${CLIENT_ID}" \
  -d "client_secret=${CLIENT_SECRET}" \
  -d "username=${USERNAME}" \
  -d "password=${PASSWORD}" \
  -d "grant_type=password")

# Extract access token
ACCESS_TOKEN=$(echo $TOKEN_RESPONSE | jq -r '.access_token')

echo "Access Token: ${ACCESS_TOKEN:0:50}..."

# Use token to call API
echo "Fetching zones..."
curl -X GET "${API_URL}/api/v1/zones" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  | jq .

# Create a zone (if authorized)
echo "Creating test zone..."
curl -X POST "${API_URL}/api/v1/zones" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "domain": "test.lan",
    "records": [
      {"name": "test.lan.", "type": "A", "ttl": 300, "value": "192.168.1.1"}
    ]
  }' | jq .
```

Save this as `keycloak-auth-example.sh` and run it to test the complete flow.
