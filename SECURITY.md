# Security Policy

## Supported Versions

We release patches for security vulnerabilities. Currently supported versions:

| Version | Supported          |
| ------- | ------------------ |
| latest  | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

If you discover a security vulnerability in GoDNS, please report it by emailing the maintainers directly rather than opening a public issue.

**Please do NOT create a public GitHub issue for security vulnerabilities.**

### What to include in your report:

1. Description of the vulnerability
2. Steps to reproduce
3. Potential impact
4. Suggested fix (if available)

### What to expect:

- We will acknowledge receipt of your vulnerability report within 48 hours
- We will provide a detailed response within 7 days
- We will work with you to understand and resolve the issue
- We will credit you in the security advisory (unless you prefer to remain anonymous)

## Security Features

### Container Security

Our Docker images are built with security in mind:

- **Distroless Base Image**: Uses Google's distroless static image (Debian 12) with minimal attack surface
- **Non-root User**: Runs as user `nonroot` (UID 65532)
- **No Shell**: Distroless images don't include a shell, reducing attack vectors
- **Read-only Filesystem**: Container runs with read-only root filesystem
- **SBOM**: Each release includes Software Bill of Materials (SBOM) in SPDX format
- **Vulnerability Scanning**: All images are scanned with Trivy before release
- **Signed Images**: Images are signed with Cosign for supply chain security

### Build Security

- **Reproducible Builds**: All builds use `-trimpath` flag
- **Static Binaries**: CGO disabled for static linking
- **Stripped Binaries**: Debug info removed with `-s -w` ldflags
- **Multi-arch Support**: Native builds for amd64 and arm64

### Dependencies

- Regular dependency updates via Dependabot
- Security scanning with gosec in CI/CD
- Go module verification in builds

## Best Practices

When deploying GoDNS:

1. Always use specific version tags, never `latest` in production
2. Verify image digests before deployment
3. Use Kubernetes security contexts with restrictive settings
4. Enable network policies to restrict traffic
5. Run regular security scans on deployed containers
6. Keep your GoDNS installation up to date

## Kubernetes Security

Our Helm chart includes security best practices:

```yaml
securityContext:
  capabilities:
    drop:
      - ALL
  readOnlyRootFilesystem: true
  runAsNonRoot: true
  runAsUser: 1000
```

## Security Updates

Security updates are released as soon as possible after a vulnerability is confirmed. Subscribe to releases to be notified of security updates.
