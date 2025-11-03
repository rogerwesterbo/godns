# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Release automation with GitHub Actions
- Multi-arch Docker images (linux/amd64, linux/arm64)
- Helm chart packaging and publishing to GHCR
- CLI binaries for multiple platforms (Linux, macOS, Windows)
- Security scanning with Trivy
- SBOM generation for releases
- Distroless Docker image for enhanced security

### Changed

- Updated container base image to Google Distroless (Debian 12)
- Improved build security with hardening flags

### Deprecated

### Removed

### Fixed

### Security

- Container runs as non-root user (UID 65532)
- Read-only root filesystem in containers
- All capabilities dropped from containers

## [0.1.0] - YYYY-MM-DD

### Added

- Initial release
- DNS server implementation
- CLI tool for DNS queries
- Docker Compose setup with Valkey
- Comprehensive documentation

---

## Release Process

1. Update this CHANGELOG.md with release notes
2. Create a new release in GitHub with semantic version (e.g., v1.0.0)
3. GitHub Actions will automatically:
   - Build multi-arch Docker images
   - Push images to ghcr.io/rogerwesterbo/godns
   - Package and push Helm chart to ghcr.io/rogerwesterbo/helm/godns
   - Build CLI binaries for all platforms
   - Attach all artifacts to the release
   - Run security scans and generate SBOM

## Version Scheme

We use [Semantic Versioning](https://semver.org/):

- **MAJOR**: Incompatible API changes
- **MINOR**: Backwards-compatible functionality additions
- **PATCH**: Backwards-compatible bug fixes

## Links

- [GitHub Releases](https://github.com/rogerwesterbo/godns/releases)
- [Docker Images](https://github.com/rogerwesterbo/godns/pkgs/container/godns)
- [Helm Charts](https://github.com/rogerwesterbo/helm/pkgs/container/godns)
