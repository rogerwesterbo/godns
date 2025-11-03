# Release Checklist

Use this checklist when creating a new release.

## Pre-release

- [ ] All CI checks passing on main branch
- [ ] Update `CHANGELOG.md` with release notes
- [ ] Update version in `charts/godns/Chart.yaml` (if needed)
- [ ] Review and merge all PRs for this release
- [ ] Test locally with Docker Compose
- [ ] Test CLI commands
- [ ] Run security scan: `make go-security-scan`

## Creating the Release

1. **Create and push the tag:**

   ```bash
   VERSION=v1.0.0
   git tag -a $VERSION -m "Release $VERSION"
   git push origin $VERSION
   ```

2. **Create GitHub Release:**

   - Go to https://github.com/rogerwesterbo/godns/releases/new
   - Select the tag you just created
   - Release title: `v1.0.0` (use semantic version)
   - Click "Generate release notes" for automatic changelog
   - Add custom notes from CHANGELOG.md if needed
   - Check "Set as latest release" (unless it's a pre-release)
   - Click "Publish release"

3. **Automated Build:**
   - GitHub Actions will automatically build and publish:
     - Multi-arch Docker images
     - Helm charts
     - CLI binaries
   - Monitor the workflow: https://github.com/rogerwesterbo/godns/actions

## Post-release

- [ ] Verify Docker image is available: `docker pull ghcr.io/rogerwesterbo/godns:$VERSION`
- [ ] Test Helm installation: `helm install test oci://ghcr.io/rogerwesterbo/helm/godns --version $VERSION`
- [ ] Download and test CLI binaries from release page
- [ ] Verify checksums match
- [ ] Check security scan results in GitHub Security tab
- [ ] Update any dependent projects or documentation
- [ ] Announce release (Twitter, blog, Discord, etc.)

## Rollback

If issues are found after release:

1. **Mark release as pre-release:**

   - Edit the release in GitHub
   - Check "Set as a pre-release"
   - Add warning to release notes

2. **Create hotfix:**

   - Branch from the problematic tag
   - Fix the issue
   - Create a new patch release (e.g., v1.0.1)

3. **Delete problematic images (if critical):**
   ```bash
   # Contact GitHub support to delete container images if needed
   # Or use Docker Hub/GHCR UI to delete specific tags
   ```

## Release Types

### Major Release (X.0.0)

- Breaking changes
- Major new features
- API changes

### Minor Release (x.Y.0)

- New features
- Backwards compatible
- No breaking changes

### Patch Release (x.y.Z)

- Bug fixes
- Security patches
- No new features

## Version Numbering

We follow [Semantic Versioning](https://semver.org/):

```
MAJOR.MINOR.PATCH

1.2.3
│ │ └─── Patch: Bug fixes, security patches
│ └───── Minor: New features, backwards compatible
└─────── Major: Breaking changes
```

## Examples

### Patch Release (Bug Fix)

```bash
# Current version: v1.2.3
git tag -a v1.2.4 -m "Fix DNS query timeout issue"
git push origin v1.2.4
```

### Minor Release (New Feature)

```bash
# Current version: v1.2.4
git tag -a v1.3.0 -m "Add support for DNSSEC"
git push origin v1.3.0
```

### Major Release (Breaking Change)

```bash
# Current version: v1.3.0
git tag -a v2.0.0 -m "Redesign configuration API"
git push origin v2.0.0
```

## Pre-releases

For alpha, beta, or release candidates:

```bash
git tag -a v2.0.0-alpha.1 -m "Release v2.0.0 Alpha 1"
git push origin v2.0.0-alpha.1
```

Mark as "pre-release" in GitHub when creating the release.
