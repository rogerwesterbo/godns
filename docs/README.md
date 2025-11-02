# GoDNS Documentation

Complete documentation for GoDNS DNS server and CLI tool.

## 📚 Documentation Index

### Getting Started

- **[Quick Start Guide](QUICK_START.md)** ⭐ START HERE
  - 5-minute setup walkthrough
  - First-time user guide
  - Basic commands and verification

### CLI Tool Documentation

- **[CLI Cheat Sheet](CLI_CHEAT_SHEET.md)** - One-page quick reference

  - Most common commands
  - Copy-paste examples
  - Perfect for keeping handy

- **[CLI Quick Reference](CLI_QUICK_REFERENCE.md)** - Fast command lookup

  - All commands with examples
  - All flags and options
  - Common workflows
  - Troubleshooting quick fixes

- **[CLI Guide](CLI_GUIDE.md)** - Complete documentation

  - Detailed command reference
  - Advanced usage examples
  - CI/CD integration
  - Kubernetes health checks
  - Monitoring and automation scripts

- **[Finding Domains Guide](FINDING_DOMAINS.md)** - Domain discovery and management
  - How to find what domains to query
  - Adding test zones
  - Listing and managing zones
  - Network discovery tips

### Server Configuration

- **[Valkey Authentication](VALKEY_AUTH.md)** - Authentication setup
  - ACL configuration
  - Username/password setup
  - Docker Compose integration
  - Security best practices

## 🚀 Which Document Should I Read?

### I'm brand new to GoDNS

👉 Start with the **[Quick Start Guide](QUICK_START.md)**

### I need a quick command reference

👉 Use the **[CLI Cheat Sheet](CLI_CHEAT_SHEET.md)**

### I want detailed CLI documentation

👉 Read the **[CLI Guide](CLI_GUIDE.md)**

### I need to set up authentication

👉 Follow the **[Valkey Authentication](VALKEY_AUTH.md)** guide

### I need fast lookup while working

👉 Keep the **[CLI Quick Reference](CLI_QUICK_REFERENCE.md)** open

### I don't know what domain to query

👉 Read the **[Finding Domains Guide](FINDING_DOMAINS.md)**

## 📖 Documentation by Size

| Document                                      | Size        | Best For           |
| --------------------------------------------- | ----------- | ------------------ |
| [CLI Cheat Sheet](CLI_CHEAT_SHEET.md)         | 1 page      | Quick lookup       |
| [CLI Quick Reference](CLI_QUICK_REFERENCE.md) | 5 min read  | Fast reference     |
| [Quick Start Guide](QUICK_START.md)           | 5 min read  | First-time setup   |
| [Valkey Authentication](VALKEY_AUTH.md)       | 10 min read | Auth configuration |
| [CLI Guide](CLI_GUIDE.md)                     | 20 min read | Complete reference |

## 🎯 Documentation by Use Case

### Testing DNS Server

1. [Quick Start Guide](QUICK_START.md) - Initial setup
2. [CLI Cheat Sheet](CLI_CHEAT_SHEET.md) - Common test commands

### Development

1. [CLI Guide](CLI_GUIDE.md) - Advanced usage
2. [CLI Quick Reference](CLI_QUICK_REFERENCE.md) - Command reference

### Production Deployment

1. [Valkey Authentication](VALKEY_AUTH.md) - Secure setup
2. [CLI Guide](CLI_GUIDE.md) - Health checks and monitoring

### Troubleshooting

1. [CLI Quick Reference](CLI_QUICK_REFERENCE.md) - Common errors
2. [CLI Guide](CLI_GUIDE.md) - Detailed troubleshooting

## 📝 Additional Resources

### In the Repository

- [Main README](../README.md) - Project overview
- [CLI Tool README](../cmd/godnscli/README.md) - Development info
- [License](../LICENSE) - License information

### Quick Links

```bash
# Build documentation is in Makefile
make help

# Environment variables in
.env.example

# Docker configuration in
docker-compose.yaml

# Valkey ACL examples in
hack/valkey/users.acl.example
```

## 🔍 Finding What You Need

### Search by Topic

- **Commands**: [CLI Quick Reference](CLI_QUICK_REFERENCE.md)
- **Examples**: [CLI Guide](CLI_GUIDE.md)
- **Setup**: [Quick Start Guide](QUICK_START.md)
- **Security**: [Valkey Authentication](VALKEY_AUTH.md)
- **Troubleshooting**: [CLI Guide](CLI_GUIDE.md#troubleshooting)

### Search by Command

| Command         | Documentation                             |
| --------------- | ----------------------------------------- |
| `query` / `q`   | [CLI Guide](CLI_GUIDE.md#query-alias-q)   |
| `health` / `h`  | [CLI Guide](CLI_GUIDE.md#health-alias-h)  |
| `test` / `t`    | [CLI Guide](CLI_GUIDE.md#test-alias-t)    |
| `version` / `v` | [CLI Guide](CLI_GUIDE.md#version-alias-v) |

## 💡 Tips

### For First-Time Users

1. Read [Quick Start Guide](QUICK_START.md)
2. Print [CLI Cheat Sheet](CLI_CHEAT_SHEET.md)
3. Bookmark [CLI Quick Reference](CLI_QUICK_REFERENCE.md)

### For Daily Use

1. Keep [CLI Cheat Sheet](CLI_CHEAT_SHEET.md) handy
2. Refer to [CLI Quick Reference](CLI_QUICK_REFERENCE.md) when needed
3. Deep dive into [CLI Guide](CLI_GUIDE.md) for complex tasks

### For Production

1. Configure with [Valkey Authentication](VALKEY_AUTH.md)
2. Set up monitoring using [CLI Guide](CLI_GUIDE.md#advanced-usage)
3. Create automation scripts from [CLI Guide](CLI_GUIDE.md#common-use-cases)

## 🆘 Getting Help

1. Check the relevant documentation above
2. Look at examples in [CLI Guide](CLI_GUIDE.md#common-use-cases)
3. Review troubleshooting in [CLI Guide](CLI_GUIDE.md#troubleshooting)
4. Check server logs: `docker-compose logs godns`

## 📌 Quick Access

Most frequently accessed documentation:

```bash
# In your browser, bookmark these:
docs/CLI_CHEAT_SHEET.md         # Daily reference
docs/CLI_QUICK_REFERENCE.md     # Detailed lookup
docs/QUICK_START.md             # Setup guide

# Print this for your desk:
docs/CLI_CHEAT_SHEET.md
```

## 🔄 Documentation Updates

All documentation is maintained in the `docs/` directory. If you make changes to the CLI or server, please update the relevant documentation.

---

**Need something specific?** Use the search above or browse the documents by use case.
