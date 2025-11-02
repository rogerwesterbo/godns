# GoDNS CLI

A command-line tool to test and interact with your GoDNS server.

## Installation

Build the CLI tool:

```bash
make build-cli
```

The binary will be available at `./bin/godnscli`.

## Usage

### Global Flags

- `-s, --server`: DNS server address (default: `localhost:53`)
- `-v, --verbose`: Enable verbose output
- `-h, --help`: Show help

### Commands

#### 1. Query DNS Records (alias: `q`)

Query specific DNS records from your server:

```bash
# Query A record
./bin/godnscli query example.lan
# Or use the short alias
./bin/godnscli q example.lan

# Query AAAA record (IPv6)
./bin/godnscli q example.lan -t AAAA

# Query MX record
./bin/godnscli q example.lan -t MX

# Query from a different server
./bin/godnscli q example.lan -s 192.168.1.1:53

# Verbose output
./bin/godnscli q example.lan -v
```

**Supported Query Types:**

- `A` - IPv4 address
- `AAAA` - IPv6 address
- `MX` - Mail exchange
- `NS` - Name server
- `TXT` - Text record
- `CNAME` - Canonical name
- `SOA` - Start of authority
- `PTR` - Pointer record

#### 2. Health Check (alias: `h`)

Check the health status of your DNS server:

```bash
# Check default ports
./bin/godnscli health
# Or use the short alias
./bin/godnscli h

# Check custom ports
./bin/godnscli h --liveness-port 8080 --readiness-port 8081

# Check remote server
./bin/godnscli h -s dns.example.com:53
```

#### 3. Run Tests (alias: `t`)

Run a comprehensive test suite:

```bash
# Run all tests
./bin/godnscli test
# Or use the short alias
./bin/godnscli t

# Run tests with verbose output
./bin/godnscli t -v

# Test against a different server
./bin/godnscli t -s 192.168.1.1:53
```

The test suite checks:

- A records (IPv4)
- AAAA records (IPv6)
- MX records
- NS records
- External resolution (google.com)

#### 4. Version (alias: `v`)

Display version information:

```bash
./bin/godnscli version
# Or use the short alias
./bin/godnscli v
```

## Examples

### Quick Server Test

```bash
# Start your GoDNS server first
# Then run:
./bin/godnscli test
# Or use the short alias
./bin/godnscli t
```

### Query Local Domain

```bash
# Query your local DNS records
./bin/godnscli q myapp.lan
./bin/godnscli q myapp.lan -t AAAA
```

### Check Server Health

```bash
# Verify server is running and healthy
./bin/godnscli h
```

### Troubleshooting

```bash
# Verbose query to see all details
./bin/godnscli q example.lan -v

# Test with timeout
./bin/godnscli q example.lan --timeout 10
```

## Development

### Adding New Commands

1. Create a new file in `cmd/godnscli/cmd/`
2. Define your command using Cobra
3. Register it in `init()` by adding to `rootCmd`

Example:

```go
var myCmd = &cobra.Command{
    Use:   "mycommand",
    Short: "Description",
    RunE:  runMyCommand,
}

func init() {
    rootCmd.AddCommand(myCmd)
}

func runMyCommand(cmd *cobra.Command, args []string) error {
    // Implementation
    return nil
}
```

### Building for Release

```bash
# Build for current platform
make build-cli

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o bin/godnscli-linux-amd64 ./cmd/godnscli
GOOS=darwin GOARCH=arm64 go build -o bin/godnscli-darwin-arm64 ./cmd/godnscli
GOOS=windows GOARCH=amd64 go build -o bin/godnscli-windows-amd64.exe ./cmd/godnscli
```

## Shell Completion

Generate shell completion scripts:

```bash
# Bash
./bin/godnscli completion bash > /etc/bash_completion.d/godnscli

# Zsh
./bin/godnscli completion zsh > "${fpath[1]}/_godnscli"

# Fish
./bin/godnscli completion fish > ~/.config/fish/completions/godnscli.fish

# PowerShell
./bin/godnscli completion powershell > godnscli.ps1
```
