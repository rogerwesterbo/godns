# GoDNS Documentation

Complete documentation for GoDNS DNS server and CLI tools.

## 📚 Table of Contents

### Getting Started
- **[Quick Start Guide](QUICK_START.md)** - 5-minute setup walkthrough
- **[Quick Auth Reference](QUICK_AUTH_REFERENCE.md)** - Fast authentication commands
- **[Port Configuration](PORT_CONFIGURATION.md)** - Port mappings and configuration

### Authentication & Security
- **[Authentication Guide](AUTHENTICATION.md)** - Complete OAuth2/OIDC guide
- **[Authentication Implementation](AUTHENTICATION_IMPLEMENTATION.md)** - Technical details
- **[Keycloak Setup](KEYCLOAK_SETUP.md)** - Keycloak configuration
- **[Valkey Authentication](VALKEY_AUTH.md)** - Valkey ACL setup

### API Documentation
- **[API Documentation](API_DOCUMENTATION.md)** - REST API reference
- **[Export API](EXPORT_API.md)** - Zone export endpoints
- **[Export Feature Summary](EXPORT_FEATURE_SUMMARY.md)** - Export capabilities
- **[API Release Workflow](API_RELEASE_WORKFLOW.md)** - API versioning and releases

### CLI Tools
- **[CLI Guide](CLI_GUIDE.md)** - Complete CLI reference
- **[CLI Configuration](CLI_CONFIG.md)** - CLI config management

### Development & Operations
- **[Test Data Seeding](TEST_DATA_SEEDING.md)** - Test data setup
- **[Finding Domains](FINDING_DOMAINS.md)** - Domain discovery utilities

## 🚀 Quick Links

### Default Ports

| Service | Port | URL |
|---------|------|-----|
| DNS Server | 53 | - |
| HTTP API | 14000 | http://localhost:14000 |
| Swagger UI | 14000 | http://localhost:14000/swagger/index.html |
| Keycloak | 14101 | http://localhost:14101 |
| Valkey | 14103 | - |
| PostgreSQL | 14100 | - |

### Default Credentials

| Service | Username | Password |
|---------|----------|----------|
| Keycloak Admin | admin | admin |
| Test User | testuser | password |

### Common Commands

**Start the system:**
```bash
docker-compose up -d
```

**Login to CLI:**
```bash
godnscli login
```

**Check auth status:**
```bash
godnscli status
```

**Export zones:**
```bash
godnscli export --format bind
```

**View API documentation:**
```bash
open http://localhost:14000/swagger/index.html
```

**Access Keycloak admin:**
```bash
open http://localhost:14101
```

## 📖 Documentation by Use Case

### For Users
If you're using GoDNS:
1. [Quick Start Guide](QUICK_START.md) - Get started in 5 minutes
2. [Quick Auth Reference](QUICK_AUTH_REFERENCE.md) - Authentication commands
3. [CLI Guide](CLI_GUIDE.md) - Learn CLI commands
4. [Export API](EXPORT_API.md) - Export zones to different formats

### For Developers
If you're developing with/on GoDNS:
1. [Authentication](AUTHENTICATION.md) - OAuth2/OIDC implementation
2. [API Documentation](API_DOCUMENTATION.md) - REST API reference
3. [Test Data Seeding](TEST_DATA_SEEDING.md) - Set up test data
4. [Port Configuration](PORT_CONFIGURATION.md) - Configure ports

### For Operators
If you're deploying GoDNS:
1. [Port Configuration](PORT_CONFIGURATION.md) - Port mappings
2. [Keycloak Setup](KEYCLOAK_SETUP.md) - Configure authentication
3. [Valkey Authentication](VALKEY_AUTH.md) - Secure database access
4. [Authentication Implementation](AUTHENTICATION_IMPLEMENTATION.md) - Technical details

## 🔍 Documentation by Feature

### DNS Queries
- [Quick Start](QUICK_START.md)
- [CLI Guide](CLI_GUIDE.md)

### Zone Management
- [API Documentation](API_DOCUMENTATION.md)
- [Export API](EXPORT_API.md)
- [Export Feature Summary](EXPORT_FEATURE_SUMMARY.md)

### Authentication
- [Authentication Guide](AUTHENTICATION.md)
- [Quick Auth Reference](QUICK_AUTH_REFERENCE.md)
- [Keycloak Setup](KEYCLOAK_SETUP.md)

### CLI Tools
- [CLI Guide](CLI_GUIDE.md)
- [CLI Config](CLI_CONFIG.md)

### Configuration
- [Port Configuration](PORT_CONFIGURATION.md)
- [Valkey Authentication](VALKEY_AUTH.md)

## 🆘 Getting Help

### Troubleshooting

**Check service health:**
```bash
docker-compose ps
```

**View logs:**
```bash
docker-compose logs -f
```

**Check authentication:**
```bash
godnscli status
```

**Verify ports:**
See [Port Configuration](PORT_CONFIGURATION.md)

### Common Issues

| Issue | Solution |
|-------|----------|
| Can't connect to API | Check if port 14000 is available |
| Authentication fails | Run `godnscli login` again |
| Port conflicts | See [Port Configuration](PORT_CONFIGURATION.md) |
| Token expired | Use `godnscli status` to check (auto-refreshes) |

### Additional Resources

- Main README: `../README.md`
- Security Policy: `../SECURITY.md`
- Changelog: `../CHANGELOG.md`

## 📝 Contributing to Documentation

When adding documentation:

1. Place all docs in the `docs/` folder
2. Use descriptive filenames
3. Update this README with a link
4. Follow the existing structure and style
5. Include examples and code snippets
