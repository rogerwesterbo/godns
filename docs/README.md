# GoDNS Documentation

Complete documentation for GoDNS DNS server, HTTP API, and CLI tool.

## 📚 Documentation Index

### Getting Started

- **[Quick Start Guide](QUICK_START.md)** ⭐ START HERE
  - 5-minute setup walkthrough
  - First-time user guide
  - Basic commands and verification

### HTTP API Documentation

- **[API Documentation](API_DOCUMENTATION.md)** - Complete API reference

  - All REST endpoints (zones, records)
  - Interactive Swagger UI guide
  - Request/response examples
  - Error handling
  - Data models and schemas

- **[API Release Workflow](API_RELEASE_WORKFLOW.md)** - Build and deployment
  - Docker build process
  - Swagger documentation generation
  - Release procedures
  - Multi-platform builds
  - CI/CD integration

### CLI Tool Documentation

- **[CLI Guide](CLI_GUIDE.md)** - Complete CLI reference

  - Quick reference cheat sheet
  - All commands with examples
  - Advanced usage patterns
  - Troubleshooting guide

- **[Finding Domains Guide](FINDING_DOMAINS.md)** - Domain discovery
  - How to find what domains to query
  - Adding test zones
  - Listing and managing zones
  - Network discovery tips

### Configuration & Setup

- **[Valkey Authentication](VALKEY_AUTH.md)** - Database authentication
  - ACL configuration
  - Username/password setup
  - Docker Compose integration
  - Security best practices

## 🚀 Quick Navigation

### I'm brand new to GoDNS

👉 Start with the **[Quick Start Guide](QUICK_START.md)**

### I want to use the REST API

👉 Read the **[API Documentation](API_DOCUMENTATION.md)** and use Swagger UI at `http://localhost:14082/swagger/index.html`

### I want to use the CLI

👉 Check the **[CLI Guide](CLI_GUIDE.md)** for the quick reference cheat sheet

### I need to build/deploy the API

👉 Follow the **[API Release Workflow](API_RELEASE_WORKFLOW.md)**

### I need to set up authentication

👉 Follow the **[Valkey Authentication](VALKEY_AUTH.md)** guide

### I don't know what domain to query

👉 Read the **[Finding Domains Guide](FINDING_DOMAINS.md)**

## 📖 Documentation by Component

### HTTP API Server

| Document                                        | Description                 | Best For               |
| ----------------------------------------------- | --------------------------- | ---------------------- |
| [API Documentation](API_DOCUMENTATION.md)       | Complete REST API reference | API users, integration |
| [API Release Workflow](API_RELEASE_WORKFLOW.md) | Build, release, deploy      | DevOps, deployment     |
| [Quick Start](QUICK_START.md)                   | Getting started             | First-time users       |

**Interactive:** Swagger UI at http://localhost:14082/swagger/index.html

### CLI Tool

| Document                              | Description            | Best For                 |
| ------------------------------------- | ---------------------- | ------------------------ |
| [CLI Guide](CLI_GUIDE.md)             | Complete CLI reference | All CLI users            |
| [Finding Domains](FINDING_DOMAINS.md) | Domain discovery       | Testing, troubleshooting |
| [Quick Start](QUICK_START.md)         | Getting started        | First-time users         |

### Server Configuration

| Document                                | Description         | Best For               |
| --------------------------------------- | ------------------- | ---------------------- |
| [Valkey Authentication](VALKEY_AUTH.md) | Database auth setup | Production deployments |
| [Quick Start](QUICK_START.md)           | Basic configuration | Development            |

## 🎯 Documentation by Use Case

### Testing the DNS Server

1. [Quick Start Guide](QUICK_START.md) - Initial setup
2. [CLI Guide](CLI_GUIDE.md) - Test commands

### Using the HTTP API

1. [API Documentation](API_DOCUMENTATION.md) - API reference
2. Swagger UI - Interactive testing at http://localhost:14082/swagger/index.html
3. Test scripts in `hack/test-api.sh`

### Building and Releasing

1. [API Release Workflow](API_RELEASE_WORKFLOW.md) - Complete release process
2. Makefile targets: `make swagger`, `make release`

### Production Deployment

1. [Valkey Authentication](VALKEY_AUTH.md) - Secure setup
2. [API Release Workflow](API_RELEASE_WORKFLOW.md) - Deployment procedures
3. [CLI Guide](CLI_GUIDE.md) - Health checks and monitoring

## 📝 Additional Resources

### In the Repository

- [Main README](../README.md) - Project overview
- [CLI Tool README](../cmd/godnscli/README.md) - CLI quick start
- [API Server README](../cmd/godnsapi/README.md) - API quick start
- [License](../LICENSE) - License information

### Quick Links

```bash
# Build documentation
make help

# Environment variables
.env.example

# Docker configuration
docker-compose.yaml

# Test scripts
hack/test-api.sh        # API testing
hack/test-swagger.sh    # Swagger UI testing
hack/add-test-zone.sh   # Add test data
```

## 🔍 Finding What You Need

### Search by Topic

- **REST API**: [API Documentation](API_DOCUMENTATION.md)
- **Swagger**: [API Documentation](API_DOCUMENTATION.md#swaggeropenapi-integration)
- **CLI Commands**: [CLI Guide](CLI_GUIDE.md)
- **Setup**: [Quick Start Guide](QUICK_START.md)
- **Security**: [Valkey Authentication](VALKEY_AUTH.md)
- **Deployment**: [API Release Workflow](API_RELEASE_WORKFLOW.md)

### Search by Command/Operation

| What                  | Where                                                         |
| --------------------- | ------------------------------------------------------------- |
| List zones (API)      | [API Documentation](API_DOCUMENTATION.md#list-all-zones)      |
| Create zone (API)     | [API Documentation](API_DOCUMENTATION.md#create-zone)         |
| Query DNS (CLI)       | [CLI Guide](CLI_GUIDE.md#quick-reference-cheat-sheet)         |
| Health check (CLI)    | [CLI Guide](CLI_GUIDE.md#quick-reference-cheat-sheet)         |
| Build Docker image    | [API Release Workflow](API_RELEASE_WORKFLOW.md#docker)        |
| Generate Swagger docs | [API Release Workflow](API_RELEASE_WORKFLOW.md#documentation) |

## 💡 Tips

### For First-Time Users

1. Read [Quick Start Guide](QUICK_START.md)
2. Try Swagger UI for API at http://localhost:14082/swagger/index.html
3. Bookmark [CLI Guide](CLI_GUIDE.md) for CLI reference

### For API Development

1. Use Swagger UI for interactive testing
2. Refer to [API Documentation](API_DOCUMENTATION.md) for details
3. Run `hack/test-api.sh` for automated testing
4. Run `make swagger` after code changes

### For Production

1. Configure with [Valkey Authentication](VALKEY_AUTH.md)
2. Follow [API Release Workflow](API_RELEASE_WORKFLOW.md) for deployments
3. Set up monitoring using [CLI Guide](CLI_GUIDE.md)

## 🆘 Getting Help

1. Check the relevant documentation above
2. Look at Swagger UI for API: http://localhost:14082/swagger/index.html
3. Review examples in [API Documentation](API_DOCUMENTATION.md) or [CLI Guide](CLI_GUIDE.md)
4. Run test scripts: `hack/test-api.sh` or `hack/test-swagger.sh`
5. Check server logs: `docker-compose logs godns`

## 🔄 Documentation Updates

All documentation is maintained in the `docs/` directory. If you make changes to the API or CLI, please:

1. Update relevant documentation
2. Regenerate Swagger docs: `make swagger`
3. Update examples in test scripts

---

**Need something specific?** Browse by component, use case, or search above.
