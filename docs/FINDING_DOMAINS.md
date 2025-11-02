# Finding and Managing DNS Zones in GoDNS

This guide helps you discover what domains are available to query and how to add your own test zones.

## Quick Answer: What Should I Query?

### Option 1: Use the Discover Command

```bash
./bin/godnscli discover
# or use the short alias
./bin/godnscli d
```

This will show you:

- Server information
- How to find available zones
- Network discovery tips
- Example queries

### Option 2: List Zones in Valkey

```bash
# Connect to Valkey and list all zones
docker-compose exec valkey valkey-cli -a mysecretpassword KEYS 'dns:zone:*'
```

Output example:

```
1) "dns:zone:example.lan."
2) "dns:zone:myapp.lan."
3) "dns:zone:test.lan."
```

### Option 3: Add a Test Zone

```bash
# Use the helper script
./hack/add-test-zone.sh test.lan 192.168.1.100

# Then query it
./bin/godnscli q test.lan
```

## Detailed: Finding Available Domains

### Method 1: Check Valkey Directly

GoDNS stores all DNS zones in Valkey (Redis) with the key pattern `dns:zone:<domain>`.

```bash
# Start a Valkey CLI session
docker-compose exec valkey valkey-cli

# Authenticate (password from .env or docker-compose.yaml)
AUTH default mysecretpassword

# List all DNS zones
KEYS dns:zone:*

# View a specific zone
GET dns:zone:example.lan.

# Exit
exit
```

### Method 2: Use the Discover Command

The `discover` command (alias: `d`, `find`) provides comprehensive information:

```bash
# Basic discovery
./bin/godnscli d

# Verbose mode shows network interfaces
./bin/godnscli d -v

# Discover a remote server
./bin/godnscli d -s 192.168.1.100:53 -v
```

### Method 3: Watch Server Logs

```bash
# Watch GoDNS logs to see queries being processed
docker-compose logs -f godns

# In another terminal, make a query
./bin/godnscli q example.lan

# The logs will show if the zone exists
```

## Adding Test Zones

### Using the Helper Script (Recommended)

The easiest way to add a test zone:

```bash
# Basic usage (domain defaults to test.lan., IP to 192.168.1.100)
./hack/add-test-zone.sh

# Custom domain and IP
./hack/add-test-zone.sh myapp.lan 192.168.1.50

# Multiple zones
./hack/add-test-zone.sh app1.lan 192.168.1.10
./hack/add-test-zone.sh app2.lan 192.168.1.20
./hack/add-test-zone.sh app3.lan 192.168.1.30
```

### Manual Method via Valkey CLI

```bash
# Connect to Valkey
docker-compose exec valkey valkey-cli -a mysecretpassword

# Create a zone (example with JSON)
SET dns:zone:myapp.lan. '{
  "domain": "myapp.lan.",
  "records": [
    {
      "name": "myapp.lan.",
      "type": "A",
      "ttl": 300,
      "value": "192.168.1.100"
    },
    {
      "name": "myapp.lan.",
      "type": "NS",
      "ttl": 3600,
      "value": "ns1.myapp.lan."
    },
    {
      "name": "myapp.lan.",
      "type": "SOA",
      "ttl": 3600,
      "value": "ns1.myapp.lan. hostmaster.myapp.lan. 1 3600 600 604800 300"
    }
  ]
}'

# Verify it was saved
GET dns:zone:myapp.lan.

# Exit
exit
```

### Testing Your New Zone

```bash
# Query the new zone
./bin/godnscli q myapp.lan

# Verbose output
./bin/godnscli q myapp.lan -v

# Different record types
./bin/godnscli q myapp.lan -t A
./bin/godnscli q myapp.lan -t NS
./bin/godnscli q myapp.lan -t SOA
```

## Zone Name Requirements

### Important: Trailing Dot

DNS zone names **should end with a dot** (`.`):

- ✅ Correct: `example.lan.`
- ❌ Incorrect: `example.lan`

The helper script automatically adds the trailing dot if missing.

### Valid Zone Names

```bash
# Good examples
test.lan.
myapp.lan.
example.local.
server.home.
192.168.1.in-addr.arpa.  # For reverse DNS

# These work but are non-standard for local use
example.com.
myapp.org.
```

## Common Scenarios

### Scenario 1: Brand New Setup

```bash
# 1. Start GoDNS
docker-compose up -d

# 2. Discover what's there
./bin/godnscli d

# 3. Nothing there yet? Add a test zone
./hack/add-test-zone.sh test.lan 192.168.1.100

# 4. Query it
./bin/godnscli q test.lan

# 5. Test external resolution (uses upstream DNS)
./bin/godnscli q google.com
```

### Scenario 2: On Unknown Network

```bash
# 1. Find your local IP
./bin/godnscli d -v

# 2. Find the gateway/router
netstat -rn | grep default
# or
ip route | grep default

# 3. See if DNS is running locally
./bin/godnscli h

# 4. Add a zone for testing
./hack/add-test-zone.sh localtest.lan 192.168.1.50

# 5. Query it
./bin/godnscli q localtest.lan
```

### Scenario 3: Remote GoDNS Server

```bash
# 1. Discover the remote server
./bin/godnscli d -s 192.168.1.100:53 -v

# 2. SSH to the server and add a zone
ssh user@192.168.1.100
cd /path/to/godns
./hack/add-test-zone.sh remote.lan 10.0.0.5

# 3. Back on your machine, query it
./bin/godnscli q remote.lan -s 192.168.1.100:53
```

### Scenario 4: Kubernetes/Multi-Pod

```bash
# 1. Find a pod name
kubectl get pods -n godns

# 2. Exec into the pod
kubectl exec -it godns-pod-abc123 -n godns -- /bin/sh

# 3. Check what zones exist
# (Inside pod)
valkey-cli -a $VALKEY_PASSWORD KEYS 'dns:zone:*'

# 4. Add zone via Valkey (from your machine)
kubectl port-forward svc/valkey 6379:6379 -n godns
# Then use local valkey-cli or the helper script
```

## Listing All Zones

### One-liner to List All Zones

```bash
# List just the zone names
docker-compose exec valkey valkey-cli -a mysecretpassword KEYS 'dns:zone:*'

# List zones with details (requires jq)
for key in $(docker-compose exec -T valkey valkey-cli -a mysecretpassword KEYS 'dns:zone:*'); do
  echo "Zone: $key"
  docker-compose exec -T valkey valkey-cli -a mysecretpassword GET "$key" | jq '.'
  echo ""
done
```

### Create a Script to List Zones

```bash
#!/bin/bash
# list-zones.sh

echo "DNS Zones in GoDNS:"
echo "==================="
echo ""

docker-compose exec -T valkey valkey-cli -a mysecretpassword KEYS 'dns:zone:*' | while read zone; do
  domain=$(echo "$zone" | sed 's/dns:zone://')
  echo "✓ $domain"
done
```

## Deleting Zones

### Delete a Single Zone

```bash
# Via Valkey CLI
docker-compose exec valkey valkey-cli -a mysecretpassword DEL dns:zone:test.lan.

# Verify deletion
docker-compose exec valkey valkey-cli -a mysecretpassword KEYS 'dns:zone:*'
```

### Delete All Zones

```bash
# WARNING: This deletes ALL DNS zones!
docker-compose exec valkey valkey-cli -a mysecretpassword --eval "
  local keys = redis.call('KEYS', 'dns:zone:*')
  for i=1,#keys do
    redis.call('DEL', keys[i])
  end
  return #keys
"
```

## Troubleshooting

### "No answer section" When Querying

This usually means the zone doesn't exist:

```bash
# 1. Check if zone exists
docker-compose exec valkey valkey-cli -a mysecretpassword GET dns:zone:example.lan.

# 2. If empty, add the zone
./hack/add-test-zone.sh example.lan 192.168.1.100

# 3. Try query again
./bin/godnscli q example.lan
```

### Can't Connect to Valkey

```bash
# Check if Valkey is running
docker-compose ps

# Check Valkey logs
docker-compose logs valkey

# Test Valkey connection
docker-compose exec valkey valkey-cli ping
# Should return: PONG

# Test with password
docker-compose exec valkey valkey-cli -a mysecretpassword ping
```

### Zone Added But Still No Answer

```bash
# 1. Verify zone in Valkey
docker-compose exec valkey valkey-cli -a mysecretpassword GET dns:zone:test.lan.

# 2. Check GoDNS logs
docker-compose logs godns

# 3. Query with verbose
./bin/godnscli q test.lan -v

# 4. Restart GoDNS (shouldn't be needed, but try)
docker-compose restart godns
```

## Advanced: Zone JSON Structure

A complete zone JSON looks like this:

```json
{
  "domain": "example.lan.",
  "records": [
    {
      "name": "example.lan.",
      "type": "A",
      "ttl": 300,
      "value": "192.168.1.100"
    },
    {
      "name": "www.example.lan.",
      "type": "A",
      "ttl": 300,
      "value": "192.168.1.100"
    },
    {
      "name": "example.lan.",
      "type": "MX",
      "ttl": 600,
      "value": "10 mail.example.lan."
    },
    {
      "name": "example.lan.",
      "type": "NS",
      "ttl": 3600,
      "value": "ns1.example.lan."
    },
    {
      "name": "example.lan.",
      "type": "SOA",
      "ttl": 3600,
      "value": "ns1.example.lan. hostmaster.example.lan. 1 3600 600 604800 300"
    },
    {
      "name": "example.lan.",
      "type": "TXT",
      "ttl": 300,
      "value": "v=spf1 -all"
    }
  ]
}
```

## Quick Reference

```bash
# Discover server and find domains
./bin/godnscli d -v

# List all zones
docker-compose exec valkey valkey-cli -a mysecretpassword KEYS 'dns:zone:*'

# Add a test zone
./hack/add-test-zone.sh myapp.lan 192.168.1.50

# Query the zone
./bin/godnscli q myapp.lan

# View zone details
docker-compose exec valkey valkey-cli -a mysecretpassword GET dns:zone:myapp.lan.

# Delete a zone
docker-compose exec valkey valkey-cli -a mysecretpassword DEL dns:zone:myapp.lan.
```

## See Also

- [CLI Quick Reference](CLI_QUICK_REFERENCE.md) - All CLI commands
- [CLI Guide](CLI_GUIDE.md) - Complete CLI documentation
- [Quick Start](QUICK_START.md) - Getting started guide
