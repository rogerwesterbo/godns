# Valkey Authentication Configuration

GoDNS uses ACL-based username/password authentication for Valkey, providing better security and access control.

## Quick Start

The default `docker-compose.yaml` uses ACL-based authentication with the `default` user.

**Environment Variables (.env):**

```bash
VALKEY_HOST=localhost
VALKEY_PORT=6379
VALKEY_USERNAME=default
VALKEY_PASSWORD=mysecretpassword
VALKEY_TOKEN=mysecretpassword  # Same as password
```

**Usage:**

```bash
# Copy example env file
cp .env.example .env

# Start Valkey
docker-compose up -d
```

## ACL Configuration

User accounts are configured in `hack/valkey/users.acl`. The default configuration includes:

```acl
# Default user with full permissions
user default on >mysecretpassword ~* &* +@all

# Example: Create a godns user with full permissions
# user godns on >godnspassword ~* &* +@all

# Example: Create a read-only user
# user readonly on >readonlypassword ~* &* +@read
```

### Adding New Users

1. **Edit ACL file** (`hack/valkey/users.acl`):

   ```acl
   # Add a new user
   user myuser on >mypassword ~* &* +@all
   ```

2. **Update Environment Variables** in `.env`:

   ```bash
   VALKEY_USERNAME=myuser
   VALKEY_PASSWORD=mypassword
   VALKEY_TOKEN=mypassword
   ```

3. **Restart Valkey**:
   ```bash
   docker-compose restart valkey
   ```

## Environment Variables Reference

| Variable                        | Required | Default   | Description                                 |
| ------------------------------- | -------- | --------- | ------------------------------------------- |
| `VALKEY_HOST`                   | Yes      | -         | Valkey server hostname                      |
| `VALKEY_PORT`                   | Yes      | `6379`    | Valkey server port                          |
| `VALKEY_USERNAME`               | No       | `default` | Username for ACL authentication             |
| `VALKEY_PASSWORD`               | No       | -         | Password for authentication                 |
| `VALKEY_TOKEN`                  | No       | -         | Alias for password (backward compatibility) |
| `VALKEY_MAX_RETRIES`            | No       | `3`       | Maximum retry attempts                      |
| `VALKEY_INITIAL_RETRY_DELAY_MS` | No       | `100`     | Initial retry delay in milliseconds         |

## ACL Syntax Reference

ACL rules in `hack/valkey/users.acl` use the following syntax:

```acl
user <username> <status> ><password> ~<keypattern> &<channelpattern> +<command>
```

- `<username>`: Username for authentication
- `<status>`: `on` (enabled) or `off` (disabled)
- `><password>`: Password (use `nopass` for no password)
- `~<keypattern>`: Key access pattern (`~*` = all keys)
- `&<channelpattern>`: Pub/Sub channel pattern (`&*` = all channels)
- `+<command>`: Allowed commands (`+@all` = all commands)

**Examples:**

```acl
# Full access
user admin on >adminpass ~* &* +@all

# Read-only access
user readonly on >readpass ~* &* +@read

# Specific key pattern
user appuser on >apppass ~app:* &* +@all

# Multiple passwords
user multipass on >pass1 >pass2 ~* &* +@all
```

## Migration Notes

GoDNS now uses ACL-based authentication by default. The system is backward compatible:

1. The `default` user is pre-configured with full permissions
2. Existing password-based setups work without changes (username defaults to `default`)
3. No code changes needed for migration

## Security Best Practices

1. **Change default passwords** - Never use `mysecretpassword` in production
2. **Use ACL for production** - Provides better audit trails and permission management
3. **Principle of least privilege** - Give users only the permissions they need
4. **Rotate credentials regularly** - Update passwords periodically
5. **Use environment variables** - Never commit credentials to version control

## Troubleshooting

### Connection refused

- Verify Valkey is running: `docker-compose ps`
- Check logs: `docker-compose logs valkey`

### Authentication failed

- Verify username/password in `.env` match your ACL configuration
- Check ACL file syntax: `docker-compose exec valkey valkey-cli ACL LIST`

### Permission denied

- Verify user has required permissions in ACL file
- Check current permissions: `docker-compose exec valkey valkey-cli ACL GETUSER <username>`
