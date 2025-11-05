# GoDNS CLI Configuration

The GoDNS CLI uses a YAML configuration file located at `~/.godns/config.yaml` to store settings and preferences.

## Configuration File Location

- **Config file**: `~/.godns/config.yaml`
- **Token cache**: `~/.godns/token.json`

The directory is created automatically on first use with secure permissions (0700).

## Configuration Commands

### View configuration

```bash
# Show all configuration values
godnscli config show

# Get config file path
godnscli config path

# Get a specific value
godnscli config get api.url
```

### Set configuration

```bash
# Set a configuration value
godnscli config set api.url http://localhost:14000

# Set Keycloak URL
godnscli config set keycloak_url http://localhost:14101
```

## Authentication Commands

### Login

Authenticate with the GoDNS API using OAuth2 device flow:

```bash
godnscli login
```

This will:

1. Request a device code from Keycloak
2. Display a URL and code for authorization
3. Wait for you to complete authentication in your browser
4. Save the access token to `~/.godns/token.json`

### Check status

View your current authentication status:

```bash
godnscli status
```

Shows:

- Login status
- Token expiration time
- Whether a refresh token is available

### Logout

Remove stored credentials:

```bash
godnscli logout
```

## Default Configuration

The following defaults are set automatically:

```yaml
keycloak_url: http://localhost:14101
keycloak_realm: godns
keycloak_cli_client_id: godns-cli
http_api_port: 14000
api:
  url: http://localhost:14000
development: true
```

## Environment Variables

Environment variables can override configuration values:

```bash
export KEYCLOAK_URL=https://keycloak.example.com
export HTTP_API_PORT=8080
```

## Token Management

Access tokens are automatically:

- Stored securely in `~/.godns/token.json`
- Used for all API requests
- Refreshed when expired (if refresh token available)
- Validated before each request

If a token is invalid and can't be refreshed, you'll be prompted to run `godnscli login` again.

## Using the API

Once logged in, all API commands automatically use your stored token:

```bash
# Export zones (uses stored token)
godnscli export

# Export specific zone
godnscli export example.lan --format bind

# The token is automatically added as Bearer token in Authorization header
```

## Security Notes

- Config directory has permissions 0700 (owner read/write/execute only)
- Token file has permissions 0600 (owner read/write only)
- Tokens are stored in JSON format with expiration timestamps
- Refresh tokens enable automatic token renewal
